package place

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place/geo"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	CityID        uuid.UUID
	DistributorID *uuid.UUID
	Class         string
	Website       *string
	Address       string
	Status        string
	Phone         *string
	Point         orb.Point
	Locale        string
	Name          string
	Description   string
}

func (m Service) Create(
	ctx context.Context,
	params CreateParams,
) (models.PlaceWithDetails, error) {
	now := time.Now().UTC()

	placeID := uuid.New()

	stmt := schemas.Place{
		ID:        placeID,
		CityID:    params.CityID,
		Class:     params.Class,
		Status:    enum.PlaceStatusActive,
		Verified:  false,
		Point:     params.Point,
		Address:   params.Address,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if params.DistributorID != nil {
		stmt.DistributorID = uuid.NullUUID{UUID: *params.DistributorID, Valid: true}
	}
	if params.Website != nil {
		stmt.Website = sql.NullString{String: *params.Website, Valid: true}
	}
	if params.Phone != nil {
		stmt.Phone = sql.NullString{String: *params.Phone, Valid: true}
	}

	var addr geo.Address
	trxErr := m.db.Transaction(ctx, func(ctx context.Context) error {
		err := m.db.Places().Insert(ctx, stmt)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not create place, cause %w", err),
			)
		}

		stmtLocale := schemas.PlaceLocale{

			PlaceID:     placeID,
			Locale:      params.Locale,
			Name:        params.Name,
			Description: params.Description,
		}
		err = m.db.PlaceLocales().Insert(ctx, stmtLocale)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not create place locale, cause %w", err),
			)
		}

		addr, err = m.geo.Guess(ctx, orb.Point{30.5234, 50.4501}) // Киев
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not guess address for Location, cause %w", err),
			)
		}

		return nil
	})
	if trxErr != nil {
		return models.PlaceWithDetails{}, trxErr
	}

	res := models.PlaceWithDetails{
		ID:          placeID,
		CityID:      params.CityID,
		Class:       params.Class,
		Status:      params.Status,
		Verified:    false,
		Point:       params.Point,
		CreatedAt:   now,
		UpdatedAt:   now,
		Address:     fmt.Sprintf("%+v\n", addr),
		Locale:      params.Locale,
		Name:        params.Name,
		Description: params.Description,
		Timetable:   models.Timetable{},
	}
	if params.DistributorID != nil {
		res.DistributorID = params.DistributorID
	}
	if params.Website != nil {
		res.Website = params.Website
	}
	if params.Phone != nil {
		res.Phone = params.Phone
	}

	return res, nil
}
