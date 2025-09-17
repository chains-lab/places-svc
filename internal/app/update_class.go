package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/class"
	"github.com/chains-lab/places-svc/internal/app/models"
)

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
