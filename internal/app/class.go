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
	if params.Parent != nil {
		_, err = a.classificator.GetClass(ctx, *params.Parent, constant.LocaleEN)
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

type FilterClassesParams struct {
	Parent      *string
	ParentCycle *bool
	Status      *string
	Name        *string
}

func (a App) ListClasses(
	ctx context.Context,
	filter FilterClassesParams,
	pag pagi.Response,
) ([]models.PlaceClassWithLocale, pagi.Response, error) {
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
	Parent *string
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
	if params.Parent != nil {
		ent.Parent = params.Parent
	}

	err = a.classificator.UpdateClass(ctx, *params.Parent, ent)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	class, err = a.classificator.GetClass(ctx, *params.Parent, locale)
	if err != nil {
		return models.PlaceClassWithLocale{}, err
	}

	return class, nil
}
