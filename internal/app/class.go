package app

import (
	"context"
	"errors"

	"github.com/chains-lab/places-svc/internal/app/entities"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
)

type CreateClassParams struct {
	Code   string
	Father *string
	Icon   string
	Name   string
}

func (a App) CreateClass(
	ctx context.Context,
	params CreateClassParams,
) (models.PlaceClassWithLocale, error) {
	_, err := a.classificator.GetClass(ctx, params.Code, constant.LocaleEN)
	if err != nil && !errors.Is(err, errx.ErrorClassNotFound) {
		return models.PlaceClassWithLocale{}, err
	}
	if err == nil {
		return models.PlaceClassWithLocale{}, errx.ErrorClassAlreadyExists.Raise(
			err,
		)
	}

	ent := entities.CreateClassParams{
		Code: params.Code,
		Icon: params.Icon,
		Name: params.Name,
	}
	if params.Father != nil {
		_, err = a.classificator.GetClass(ctx, *params.Father, constant.LocaleEN)
		if err != nil {
			return models.PlaceClassWithLocale{}, err
		}
	}

	class, err := a.classificator.CreateClass(ctx, ent)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	return class, err
}

func (a App) GetClass(
	ctx context.Context,
	code, locale string,
) (models.PlaceClassWithLocale, error) {
	return a.classificator.GetClass(ctx, code, locale)
}

func (a App) DeactivateClass(
	ctx context.Context,
	code, locale string,
) (models.PlaceClassWithLocale, error) {
	class, err := a.classificator.GetClass(ctx, code, locale)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	if class.Data.Status == constant.PlaceClassStatusesInactive {
		return class, nil
	}

	updated, err := a.classificator.DeactivateClass(ctx, code, locale)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	return updated, nil
}

type UpdateClassParams struct {
	Icon   *string
	Father *string
}

func (a App) UpdateClass(
	ctx context.Context,
	code, locale string,
	params UpdateClassParams,
) (models.PlaceClassWithLocale, error) {
	class, err := a.classificator.GetClass(ctx, code, locale)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	ent := entities.UpdateClassParams{}
	if params.Icon != nil {
		ent.Icon = params.Icon
	}
	if params.Father != nil {
		ent.Father = params.Father
	}

	err = a.classificator.UpdateClass(ctx, *params.Father, ent)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	class, err = a.classificator.GetClass(ctx, *params.Father, locale)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	return class, nil
}
