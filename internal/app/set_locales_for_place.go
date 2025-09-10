package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/google/uuid"
)

type SetPlaceLocalParams struct {
	Locale      string
	Name        string
	Description string
}

func (a App) SetPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	locales ...SetPlaceLocalParams,
) error {
	out := make([]place.SetLocaleParams, 0, len(locales))
	for _, locale := range locales {
		err := constant.IsValidLocaleSupported(locale.Locale)
		if err != nil {
			return err
		}

		s := place.SetLocaleParams{
			Locale:      locale.Locale,
			Name:        locale.Name,
			Description: locale.Description,
		}

		out = append(out, s)
	}

	err := a.place.SetLocales(ctx, placeID, out...)
	if err != nil {
		return err
	}

	return nil
}
