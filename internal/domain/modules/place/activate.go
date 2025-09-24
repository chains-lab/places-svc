package place

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (m Module) Deactivate(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, locale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Status == enum.PlaceStatusInactive {
		return place, nil
	}

	status := enum.PlaceStatusInactive
	updated, err := m.Update(ctx, placeID, locale, UpdateParams{
		Status: &status,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}

func (m Module) Activate(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, locale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Status == enum.PlaceStatusActive {
		return place, nil
	}

	status := enum.PlaceStatusActive
	updated, err := m.Update(ctx, placeID, locale, UpdateParams{
		Status: &status,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}
