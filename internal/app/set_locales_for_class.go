package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/class"
)

type SetClassLocaleParams struct {
	Locale string
	Name   string
}

func (a App) SetLocalesForClass(ctx context.Context, code string, locales ...SetClassLocaleParams) error {
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
