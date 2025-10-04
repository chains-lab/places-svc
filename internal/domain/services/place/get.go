package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Get(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error) {
	err := enum.CheckLocale(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	place, err := s.db.GetPlaceByID(ctx, placeID, locale)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get place with id %s: %w", placeID, err),
		)
	}

	if place.IsNil() {
		return models.Place{}, errx.ErrorPlaceNotFound.Raise(
			fmt.Errorf("place with id %s not found", placeID),
		)
	}

	return place, nil
}
