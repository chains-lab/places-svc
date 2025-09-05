package entities

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (c Classificator) UpsetLocaleToClass(ctx context.Context, class, locale, name string) (models.LocaleForClass, error) {
	_, err := c.GetClass(ctx, class, locale)
	if err != nil {
		return models.LocaleForClass{}, err
	}

	err = c.localeQ.Upsert(ctx, dbx.PlaceClassLocale{
		Class:  class,
		Locale: locale,
		Name:   name,
	})
	if err != nil {
		return models.LocaleForClass{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locale %s for class %s, cause: %w", locale, class, err),
		)
	}

	return models.LocaleForClass{
		Class:  class,
		Locale: locale,
		Name:   name,
	}, nil
}

func (c Classificator) ListForClass(ctx context.Context, class string) ([]models.LocaleForClass, error) {
	_, err := c.GetClass(ctx, class, constant.LocaleEN)
	if err != nil {
		return nil, err
	}

	locales, err := c.localeQ.New().FilterClass(class).Select(ctx)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for class %s, cause: %w", class, err),
		)
	}

	result := make([]models.LocaleForClass, 0, len(locales))
	for _, loc := range locales {
		result = append(result, classLocaleModelFromDB(loc))
	}

	return result, nil
}

func classLocaleModelFromDB(dbLoc dbx.PlaceClassLocale) models.LocaleForClass {
	return models.LocaleForClass{
		Class:  dbLoc.Class,
		Locale: dbLoc.Locale,
		Name:   dbLoc.Name,
	}
}
