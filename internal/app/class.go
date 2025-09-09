package app

import (
	"context"
	"errors"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/entities"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
)

type CreateClassParams struct {
	Code   string
	Parent *string
	Icon   string
	Name   string
}

func (a App) CreateClass(
	ctx context.Context,
	params CreateClassParams,
) (models.ClassWithLocale, error) {
	_, err := a.classificator.GetClass(ctx, params.Code, constant.LocaleEN)
	if err != nil && !errors.Is(err, errx.ErrorClassNotFound) {
		return models.ClassWithLocale{}, err
	}
	if err == nil {
		return models.ClassWithLocale{}, errx.ErrorClassAlreadyExists.Raise(
			err,
		)
	}

	ent := entities.CreateClassParams{
		Code: params.Code,
		Icon: params.Icon,
		Name: params.Name,
	}
	if params.Parent != nil {
		ent.Parent = params.Parent
	}

	class, err := a.classificator.CreateClass(ctx, ent)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return class, err
}

func (a App) GetClass(
	ctx context.Context,
	code, locale string,
) (models.ClassWithLocale, error) {
	return a.classificator.GetClass(ctx, code, locale)
}

type FilterClassesParams struct {
	Parent      *string
	ParentCycle *bool
	Status      *string
	Name        *string
}

func (a App) ListClasses(
	ctx context.Context,
	filter FilterClassesParams,
	pag pagi.Request,
) ([]models.ClassWithLocale, pagi.Response, error) {
	ent := entities.FilterClassesParams{}
	if filter.Parent != nil {
		ent.Parent = filter.Parent
	}
	if filter.ParentCycle != nil {
		ent.ParentCycle = filter.ParentCycle
	}
	if filter.Status != nil {
		ent.Status = filter.Status
	}

	return a.classificator.ListClasses(ctx, ent, pag)
}

func (a App) DeactivateClass(
	ctx context.Context,
	code, locale string,
	replaceClasses string,
) (models.ClassWithLocale, error) {
	class, err := a.classificator.GetClass(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	replaceClass, err := a.classificator.GetClass(ctx, replaceClasses, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}
	if class.Data.Code == replaceClass.Data.Code {
		return models.ClassWithLocale{}, errx.ErrorClassDeactivateReplaceSame.Raise(
			err,
		)
	}

	if class.Data.Status == constant.PlaceClassStatusesInactive {
		return class, nil
	}

	var updated models.ClassWithLocale
	txErr := a.transaction(func(txCtx context.Context) error {
		err = a.place.UpdatePlaces(ctx,
			entities.UpdatePlacesFilter{
				Class: []string{
					class.Data.Code,
				},
			},
			entities.UpdatePlaceParams{
				Class: &replaceClasses,
			},
		)

		updated, err = a.classificator.DeactivateClass(ctx, code, locale)
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

func (a App) DeleteClass(
	ctx context.Context,
	code string,
) error {
	class, err := a.classificator.GetClass(ctx, code, constant.DefaultLocale)
	if err != nil {
		return err
	}

	if class.Data.Status == constant.PlaceClassStatusesActive {
		return errx.ErrorCannotDeleteActiveClass.Raise(err)
	}

	places, _, err := a.place.SearchPlaces(ctx, constant.DefaultLocale,
		entities.SearchPlacesFilter{Class: []string{class.Data.Code}}, pagi.Request{}, nil)

	if len(places) != 0 {
		return errx.ErrorClassHasPlaces.Raise(err)
	}

	return a.classificator.DeleteClass(ctx, code)
}

func (a App) ActivateClass(
	ctx context.Context,
	code, locale string,
) (models.ClassWithLocale, error) {
	class, err := a.classificator.GetClass(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	if class.Data.Status == constant.PlaceClassStatusesActive {
		return class, nil
	}

	updated, err := a.classificator.ActivateClass(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return updated, nil
}

type UpdateClassParams struct {
	Name   *string
	Icon   *string
	Parent *string
}

func (a App) UpdateClass(
	ctx context.Context,
	code, locale string,
	params UpdateClassParams,
) (models.ClassWithLocale, error) {
	class, err := a.classificator.GetClass(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	ent := entities.UpdateClassParams{}
	if params.Icon != nil {
		ent.Icon = params.Icon
	}
	if params.Parent != nil {
		ent.Parent = params.Parent
	}

	err = a.classificator.UpdateClass(ctx, *params.Parent, ent)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	class, err = a.classificator.GetClass(ctx, *params.Parent, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return class, nil
}
