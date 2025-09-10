package place

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/google/uuid"
)

func (p Place) Deactivate(ctx context.Context, locale string, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, locale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Place.Status == constant.PlaceStatusInactive {
		return place, nil
	}

	status := constant.PlaceStatusInactive
	updated, err := p.UpdatePlace(ctx, placeID, locale, UpdatePlaceParams{
		Status: &status,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}

func (p Place) Activate(ctx context.Context, locale string, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, locale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Place.Status == constant.PlaceStatusActive {
		return place, nil
	}

	status := constant.PlaceStatusActive
	updated, err := p.UpdatePlace(ctx, placeID, locale, UpdatePlaceParams{
		Status: &status,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}
