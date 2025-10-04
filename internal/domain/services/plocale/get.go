package plocale

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) GetForPlace(
	ctx context.Context,
	placeID uuid.UUID,
	page uint64,
	size uint64,
) (models.PlaceLocaleCollection, error) {
	exist, err := s.db.PlaceExists(ctx, placeID)
	if err != nil {
		return models.PlaceLocaleCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check existence of place %s, cause: %w", placeID, err),
		)
	}
	if !exist {
		return models.PlaceLocaleCollection{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place %s not found", placeID),
		)
	}

	locales, err := s.db.GetPlaceLocales(ctx, placeID, page, size)
	if err != nil {
		return models.PlaceLocaleCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for place %s, cause: %w", placeID, err),
		)
	}

	return locales, nil
}
