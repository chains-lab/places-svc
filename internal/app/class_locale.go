package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/domain/class"
	"github.com/chains-lab/places-svc/internal/app/models"
)

type SetClassLocaleParams struct {
	Locale string
	Name   string
}

func (a App) SetClassLocales(ctx context.Context, code string, locales ...SetClassLocaleParams) error {
	out := make([]class.SetClassLocaleParams, 0, len(locales))
	for _, locale := range locales {
		s := class.SetClassLocaleParams{
			Locale: locale.Locale,
			Name:   locale.Name,
		}

		out = append(out, s)
	}

	err := a.classificator.SetLocales(ctx, code, out...)
	if err != nil {
		return err
	}

	return nil
}

func (a App) ListClassLocales(
	ctx context.Context,
	class string,
	pag pagi.Request,
) ([]models.ClassLocale, pagi.Response, error) {
	return a.classificator.LocalesList(ctx, class, pag)
}
