package place

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

func (m Module) Get(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithDetails, error) {
	err := enum.IsValidLocaleSupported(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	place, err := m.db.Places().FilterID(placeID).GetWithDetails(ctx, locale)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.PlaceWithDetails{}, errx.ErrorPlaceNotFound.Raise(
				fmt.Errorf("place with id %m not found, cause %w", placeID, err),
			)
		default:
			return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get place with id %m: %w", placeID, err),
			)
		}
	}

	return placeWithDetailsModelFromDB(place), nil
}
