package place

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	ID            uuid.UUID
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

func (p Place) Create(
	ctx context.Context,
	params CreateParams,
) (models.PlaceWithDetails, error) {
	now := time.Now().UTC()

	stmt := dbx.Place{
		ID:        params.ID,
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

	err := p.query.New().Insert(ctx, stmt)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not create Location, cause %w", err),
		)
	}

	stmtLocale := dbx.PlaceLocale{
		PlaceID:     params.ID,
		Locale:      params.Locale,
		Name:        params.Name,
		Description: params.Description,
	}
	err = p.locale.Insert(ctx, stmtLocale)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not create Location locale, cause %w", err),
		)
	}

	addr, err := p.geo.Guess(ctx, orb.Point{30.5234, 50.4501}) // Киев
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not guess address for Location, cause %w", err),
		)
	}

	res := models.PlaceWithDetails{
		ID:          params.ID,
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
