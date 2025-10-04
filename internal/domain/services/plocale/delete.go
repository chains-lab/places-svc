package plocale

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/google/uuid"
)

func (s Service) Delete(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) error {
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

	total, err := s.db.CountPlaceLocales(ctx, placeID)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get locales for place %s, cause: %w", placeID, err),
		)
	}
	if total <= 1 {
		return errx.ErrorNeedAtLeastOneLocaleForPlace.Raise(
			fmt.Errorf("cannot delete the last locale for place %s", placeID),
		)
	}

	err = s.db.DeletePlaceLocale(ctx, placeID, locale)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locale %s for place %s, cause: %w", locale, placeID, err),
		)
	}

	return nil
}
