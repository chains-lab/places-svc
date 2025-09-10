package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) SetPlaceTimeTable(ctx context.Context, placeID uuid.UUID, intervals ...models.TimeInterval) error {
	return a.place.SetTimetable(ctx, placeID, intervals...)
}
