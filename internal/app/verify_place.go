package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) VerifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	return a.place.VerifyPlace(ctx, placeID)
}

func (a App) UnverifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	return a.place.UnverifyPlace(ctx, placeID)
}
