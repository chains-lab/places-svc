package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) ActivatePlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithDetails, error) {
	return a.place.Activate(ctx, placeID, locale)
}

func (a App) DeactivatePlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithDetails, error) {
	return a.place.Deactivate(ctx, placeID, locale)
}
