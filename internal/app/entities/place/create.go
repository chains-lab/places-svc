package place

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
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
	Status        string
	Website       *string
	Phone         *string
	Point         orb.Point
}

type CreateLocalParams struct {
	Locale      string
	Name        string
	Description string
}

func (p Place) Create(
	ctx context.Context,
	params CreateParams,
	locale CreateLocalParams,
) (models.PlaceWithDetails, error) {
	now := time.Now().UTC()

	stmt := dbx.Place{
		ID:        params.ID,
		CityID:    params.CityID,
		Class:     params.Class,
		Status:    constant.PlaceStatusActive,
		Verified:  false,
		Point:     params.Point,
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
		Locale:      locale.Locale,
		Name:        locale.Name,
		Description: locale.Description,
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

	res := models.Place{
		ID:        params.ID,
		CityID:    params.CityID,
		Class:     params.Class,
		Status:    params.Status,
		Verified:  false,
		Point:     params.Point,
		CreatedAt: now,
		UpdatedAt: now,
		Address:   fmt.Sprintf("%+v\n", addr),
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

	paramsLocale := models.PlaceLocale{
		PlaceID:     params.ID,
		Locale:      locale.Locale,
		Name:        locale.Name,
		Description: locale.Description,
	}

	return models.PlaceWithDetails{
		Place:  res,
		Locale: paramsLocale,
	}, nil
}
