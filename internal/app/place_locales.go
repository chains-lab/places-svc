package app

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/domain/place"
	"github.com/chains-lab/places-svc/internal/app/models"
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
		err := enum.IsValidLocaleSupported(locale.Locale)
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

func (a App) ListPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	pag pagi.Request,
) ([]models.PlaceLocale, pagi.Response, error) {
	return a.place.ListLocales(ctx, placeID, pag)
}
