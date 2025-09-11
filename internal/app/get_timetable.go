package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) GetTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	place, err := a.place.GetTimetable(ctx, placeID)
	if err != nil {
		return models.Timetable{}, err
	}

	return place, nil
}
