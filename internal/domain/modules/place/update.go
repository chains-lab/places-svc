package place

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateParams struct {
	Class    *string
	Status   *string
	Verified *bool
	Point    *orb.Point
	Website  *string
	Phone    *string
	Address  *string
}

func (m Module) Update(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdateParams,
) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, locale) //TODO locale
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	stmt := data.UpdatePlaceParams{
		UpdatedAt: time.Now().UTC(),
	}

	if params.Class != nil {
		_, err = m.db.Classes().FilterCode(*params.Class).Get(ctx)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return models.PlaceWithDetails{}, errx.ErrorClassNotFound.Raise(
					fmt.Errorf("class with code %m not found, cause: %w", *params.Class, err),
				)
			default:
				return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
					fmt.Errorf("failed to get class with code %m, cause: %w", *params.Class, err),
				)
			}
		}
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

	err = m.db.Places().FilterID(placeID).Update(ctx, stmt)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update place with id %m, cause: %w", placeID, err),
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

func (m Module) UpdatePlaces(
	ctx context.Context,
	filter UpdatePlacesFilter,
	params UpdateParams,
) error {
	query := m.db.Places()

	if filter.Class != nil {
		query = query.FilterClass(*filter.Class)
	}
	if filter.CityID != nil {
		query = query.FilterCityID(*filter.CityID)
	}
	if filter.DistributorID != nil {
		query = query.FilterDistributorID(*filter.DistributorID)
	}

	stmt := data.UpdatePlaceParams{}
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
