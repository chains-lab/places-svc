package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) ListPlaceTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	return a.place.ListTimetable(ctx, placeID)
}
