package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(ctx context.Context, placeID uuid.UUID) error {
	place, err := s.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return err
	}

	if place.Status != enum.PlaceStatusInactive {
		return errx.ErrorPlaceForDeleteMustBeInactive.Raise(
			fmt.Errorf("place %s is not in inactive status", place.ID.String()),
		)
	}

	err = s.db.DeletePlace(ctx, placeID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete Location locale with id %s, cause: %w", placeID, err),
		)
	}

	return nil
}
