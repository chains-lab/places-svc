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

type UpdatePlaceParams struct {
	Class    *string
	Status   *string
	Verified *bool
	Point    *orb.Point
	Website  *string
	Phone    *string
	Address  *string
}

func (p Place) UpdatePlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdatePlaceParams,
) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, locale) //TODO locale
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	stmt := dbx.UpdatePlaceParams{
		UpdatedAt: time.Now().UTC(),
	}

	if params.Class != nil {
		stmt.Class = params.Class
		place.Class = *params.Class
	}

	if params.Status != nil {
		stmt.Status = params.Status
		place.Status = *params.Status
	}

	if params.Verified != nil {
		stmt.Verified = params.Verified
		place.Verified = *params.Verified
	}

	if params.Point != nil {
		stmt.Point = params.Point
		place.Point = *params.Point
	}

	if params.Website != nil {
		switch *params.Website {
		case "":
			stmt.Website = &sql.NullString{Valid: false}
			place.Website = nil
		default:
			stmt.Website = &sql.NullString{String: *params.Website, Valid: true}
			place.Website = params.Website
		}
	}

	if params.Phone != nil {
		switch *params.Phone {
		case "":
			stmt.Phone = &sql.NullString{Valid: false}
			place.Phone = nil
		default:
			stmt.Phone = &sql.NullString{String: *params.Phone, Valid: true}
			place.Phone = params.Phone
		}
	}

	err = p.query.New().FilterID(placeID).Update(ctx, stmt)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update place with id %s, cause: %w", placeID, err),
		)
	}

	return place, nil
}

type UpdatePlacesFilter struct {
	Class         *string
	CityID        *uuid.UUID
	DistributorID *uuid.UUID
}

type UpdatePlacesParams struct {
	Class    *string
	Status   *string
	Verified *bool
}

func (p Place) UpdatePlaces(
	ctx context.Context,
	filter UpdatePlacesFilter,
	params UpdatePlaceParams,
) error {
	query := p.query.New()

	if filter.Class != nil {
		query = query.FilterClass(*filter.Class)
	}
	if filter.CityID != nil {
		query = query.FilterCityID(*filter.CityID)
	}
	if filter.DistributorID != nil {
		query = query.FilterDistributorID(*filter.DistributorID)
	}

	stmt := dbx.UpdatePlaceParams{}
	if params.Class != nil {
		stmt.Class = params.Class
	}
	if params.Status != nil {
		err := enum.IsValidPlaceStatus(*params.Status)
		if err != nil {
			return errx.ErrorPlaceStatusInvalid.Raise(
				fmt.Errorf("invalid place status, cause: %w", err),
			)
		}

		stmt.Status = params.Status
	}
	if params.Verified != nil {
		stmt.Verified = params.Verified
	}

	err := query.Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update locos, cause: %w", err),
		)
	}

	return nil
}
