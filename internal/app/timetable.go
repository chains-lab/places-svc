package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (a App) SetPlaceTimeTable(ctx context.Context, placeID uuid.UUID, intervals models.Timetable) (models.PlaceWithDetails, error) {
	return a.place.SetTimetable(ctx, placeID, intervals)
}

func (a App) GetTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	place, err := a.place.GetTimetable(ctx, placeID)
	if err != nil {
		return models.Timetable{}, err
	}

	return place, nil
}

func (a App) DeleteTimetable(ctx context.Context, placeID uuid.UUID) error {
	return a.place.DeleteTimetable(ctx, placeID)
}
