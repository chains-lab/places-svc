package timetable

import (
	"context"

	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db database
}

func NewService(db database) Service {
	return Service{db: db}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	PlaceExists(ctx context.Context, placeID uuid.UUID) (bool, error)

	GetPlaceByID(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error)

	SetTimetable(ctx context.Context, placeID uuid.UUID, intervals models.Timetable) error
	GetTimetableByPlaceID(ctx context.Context, placeID uuid.UUID) (models.Timetable, error)
	DeleteTimetableByPlaceID(ctx context.Context, placeID uuid.UUID) error
}
