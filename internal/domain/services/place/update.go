package place

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateParams struct {
	Class   *string
	Status  *string
	Point   *orb.Point
	Website *string
	Phone   *string
	Address *string
}

func (m Service) Update(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdateParams,
) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, locale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	stmt := schemas.UpdatePlaceParams{
		UpdatedAt: time.Now().UTC(),
	}

	if params.Class != nil {
		_, err = m.db.Classes().FilterCode(*params.Class).Get(ctx)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return models.PlaceWithDetails{}, errx.ErrorClassNotFound.Raise(
					fmt.Errorf("class with code %s not found, cause: %w", *params.Class, err),
				)
			default:
				return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
					fmt.Errorf("failed to get class with code %s, cause: %w", *params.Class, err),
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
