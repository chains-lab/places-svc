package timetable

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) SetTimetable(
	ctx context.Context,
	placeID uuid.UUID,
	intervals models.Timetable,
) (models.Place, error) {
	place, err := s.db.GetPlaceByID(ctx, placeID)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get place %s, cause: %w", placeID, err),
		)
	}

	if place.IsNil() {
		return models.Place{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
	}

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.DeleteTimetableByPlaceID(ctx, placeID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not upsert timetable, cause: %w", err),
			)
		}

		err = s.db.SetTimetable(ctx, placeID, intervals)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not upsert timetable, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Place{}, err
	}

	place.Timetable = intervals

	return place, nil
}
