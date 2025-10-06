package timetable

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) DeleteForPlace(ctx context.Context, placeID uuid.UUID) error {
	exist, err := s.db.PlaceExists(ctx, placeID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check existence of place %s, cause: %w", placeID, err),
		)
	}
	if !exist {
		return errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
	}

	err = s.db.DeleteTimetableByPlaceID(ctx, placeID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}
