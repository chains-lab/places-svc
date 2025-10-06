package timetable

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetForPlace(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	rows, err := s.db.GetTimetableByPlaceID(ctx, placeID)
	if err != nil {
		return models.Timetable{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not list timetable, cause: %w", err),
		)
	}

	return rows, nil
}
