package place

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type UpdateParams struct {
	Class   *string
	Point   *orb.Point
	Website *string
	Phone   *string
	Address *string
}

func (s Service) Update(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdateParams,
) (models.Place, error) {
	place, err := s.Get(ctx, placeID, locale)
	if err != nil {
		return models.Place{}, err
	}

	if params.Class != nil {
		place.Class = *params.Class
	}
	if params.Point != nil {
		place.Point = *params.Point
	}
	if params.Website != nil {
		place.Website = params.Website
	}
	if params.Phone != nil {
		place.Phone = params.Phone
	}
	if params.Address != nil {
		place.Address = *params.Address
	}
	place.UpdatedAt = time.Now().UTC()

	err = s.db.UpdatePlace(ctx, placeID, params, place.UpdatedAt)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update place, cause: %w", err),
		)
	}

	return place, nil
}

func (s Service) UpdateStatus(
	ctx context.Context,
	placeID uuid.UUID,
	status string,
	locale string,
) (models.Place, error) {
	place, err := s.Get(ctx, placeID, locale)
	if err != nil {
		return models.Place{}, err
	}

	err = enum.CheckPlaceStatus(status)
	if err != nil {
		return models.Place{}, err
	}

	place.UpdatedAt = time.Now().UTC()
	place.Status = status

	err = s.db.UpdatePlaceStatus(ctx, placeID, status, place.UpdatedAt)
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}
