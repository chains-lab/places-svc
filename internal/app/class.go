package app

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/domain/class"
	"github.com/chains-lab/places-svc/internal/app/domain/place"
	"github.com/chains-lab/places-svc/internal/app/models"
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
	var c models.ClassWithLocale
	var err error
	err = a.transaction(func(ctx context.Context) error {
		ent := class.CreateParams{
			Code: params.Code,
			Icon: params.Icon,
			Name: params.Name,
		}
		if params.Parent != nil {
			ent.Parent = params.Parent
		}

		c, err = a.classificator.Create(ctx, ent)
		return err
	})

	return c, err
}

func (a App) GetClass(
	ctx context.Context,
	code, locale string,
) (models.ClassWithLocale, error) {
	return a.classificator.Get(ctx, code, locale)
}

type FilterListClassesParams struct {
	Parent      *string
	ParentCycle bool
	Status      *string
}

func (a App) ListClasses(
	ctx context.Context,
	locale string,
	filter FilterListClassesParams,
	pag pagi.Request,
) ([]models.ClassWithLocale, pagi.Response, error) {
	ent := class.FilterListParams{
		ParentCycle: filter.ParentCycle,
	}
	if filter.Parent != nil {
		ent.Parent = filter.Parent
	}
	ent.ParentCycle = filter.ParentCycle
	if filter.Status != nil {
		ent.Status = filter.Status
	}

	return a.classificator.List(ctx, locale, ent, pag)
}

type UpdateClassParams struct {
	Icon   *string
	Parent *string
}

func (a App) UpdateClass(
	ctx context.Context,
	code, locale string,
	params UpdateClassParams,
) (models.ClassWithLocale, error) {
	c, err := a.classificator.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	ent := class.UpdateParams{}
	if params.Icon != nil {
		ent.Icon = params.Icon
	}
	if params.Parent != nil {
		ent.Parent = params.Parent
	}

	err = a.classificator.Update(ctx, code, ent)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	c, err = a.classificator.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return c, nil
}

func (a App) DeleteClass(
	ctx context.Context,
	code string,
) error {
	c, err := a.classificator.Get(ctx, code, enum.DefaultLocale)
	if err != nil {
		return err
	}

	if c.Data.Status == enum.PlaceClassStatusesActive {
		return errx.ErrorCannotDeleteActiveClass.Raise(err)
	}

	places, _, err := a.place.List(ctx, enum.DefaultLocale,
		place.FilterListParams{Classes: []string{c.Data.Code}}, pagi.Request{}, nil)

	if len(places) != 0 {
		return errx.ErrorCantDeleteClassWithPlaces.Raise(err)
	}

	return a.classificator.DeleteClass(ctx, code)
}
