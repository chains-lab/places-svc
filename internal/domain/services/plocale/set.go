package plocale

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/google/uuid"
)

type SetParams struct {
	Locale      string
	Name        string
	Description string
}

func (s Service) SetForPlace(
	ctx context.Context,
	placeID uuid.UUID,
	locales ...SetParams,
) error {
	if len(locales) == 0 {
		return nil
	}

	for _, param := range locales {
		err := enum.CheckLocale(param.Locale)
		if err != nil {
			return errx.ErrorInvalidLocale.Raise(
				fmt.Errorf("invalid locale provided: %s, cause %w", param.Locale, err),
			)
		}
	}

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

	err = s.db.UpsertLocaleForPlace(ctx, placeID, locales...)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locales for place %s, cause: %w", placeID, err),
		)
	}

	return nil
}
