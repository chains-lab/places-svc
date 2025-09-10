package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) DeletePlace(ctx context.Context, placeID uuid.UUID) error {
	return a.place.DeletePlace(ctx, placeID)
}
