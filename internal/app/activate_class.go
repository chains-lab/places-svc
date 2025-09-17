package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
)

func (a App) DeactivateClass(
	ctx context.Context,
	code, locale string,
	replace string,
) (models.ClassWithLocale, error) {
	var err error
	var updated models.ClassWithLocale
	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.place.UpdatePlaces(ctx,
			place.UpdatePlacesFilter{
				Class: &code,
			},
			place.UpdatePlaceParams{
				Class: &replace,
			},
		)

		updated, err = a.classificator.Deactivate(ctx, locale, code, replace)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.ClassWithLocale{}, txErr
	}

	return updated, nil
}

func (a App) ActivateClass(
	ctx context.Context,
	code, locale string,
) (models.ClassWithLocale, error) {
	updated, err := a.classificator.Activate(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return updated, nil
}
