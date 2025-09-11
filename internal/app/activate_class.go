package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
)

func (a App) DeactivateClass(
	ctx context.Context,
	code, locale string,
	replaceClasses string,
) (models.ClassWithLocale, error) {
	var err error
	var updated models.ClassWithLocale
	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.place.UpdatePlaces(ctx,
			place.UpdatePlacesFilter{
				Class: []string{
					code,
				},
			},
			place.UpdatePlaceParams{
				Class: &replaceClasses,
			},
		)

		updated, err = a.classificator.Deactivate(ctx, locale, code, replaceClasses)
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
