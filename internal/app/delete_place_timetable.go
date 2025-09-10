package app

import (
	"context"

	"github.com/google/uuid"
)

func (a App) DeletePlaceTimetable(ctx context.Context, placeID uuid.UUID) error {
	return a.place.DeleteTimetable(ctx, placeID)
}
