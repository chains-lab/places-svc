package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (a App) DeactivateClass(
	ctx context.Context,
	code, locale string,
	replaceClasses string,
) (models.ClassWithLocale, error) {
	c, err := a.classificator.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	replaceClass, err := a.classificator.Get(ctx, replaceClasses, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}
	if c.Data.Code == replaceClass.Data.Code {
		return models.ClassWithLocale{}, errx.ErrorClassDeactivateReplaceSame.Raise(
			err,
		)
	}

	if c.Data.Status == constant.PlaceClassStatusesInactive {
		return c, nil
	}

	var updated models.ClassWithLocale
	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.place.UpdatePlaces(ctx,
			place.UpdatePlacesFilter{
				Class: []string{
					c.Data.Code,
				},
			},
			place.UpdatePlaceParams{
				Class: &replaceClasses,
			},
		)

		updated, err = a.classificator.Deactivate(ctx, code, locale)
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
	c, err := a.classificator.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	if c.Data.Status == constant.PlaceClassStatusesActive {
		return c, nil
	}

	updated, err := a.classificator.Activate(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return updated, nil
}
