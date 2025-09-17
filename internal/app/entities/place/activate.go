package place

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (p Place) Deactivate(ctx context.Context, locale string, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, locale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Status == enum.PlaceStatusInactive {
		return place, nil
	}

	status := enum.PlaceStatusInactive
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

	if place.Status == enum.PlaceStatusActive {
		return place, nil
	}

	status := enum.PlaceStatusActive
	updated, err := p.UpdatePlace(ctx, placeID, locale, UpdatePlaceParams{
		Status: &status,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}
