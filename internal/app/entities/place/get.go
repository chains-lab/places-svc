package place

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

func (p Place) Get(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithDetails, error) {
	err := enum.IsValidLocaleSupported(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	place, err := p.query.New().FilterID(placeID).GetWithDetails(ctx, locale)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.PlaceWithDetails{}, errx.ErrorPlaceNotFound.Raise(
				fmt.Errorf("location with id %s not found, cause %w", placeID, err),
			)
		default:
			return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get place with id %s: %w", placeID, err),
			)
		}
	}

	return placeWithDetailsModelFromDB(place), nil
}
