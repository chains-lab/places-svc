package app

import (
	"context"
	"errors"

	"github.com/chains-lab/places-svc/internal/app/entities/class"
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
	_, err := a.classificator.Get(ctx, params.Code, constant.LocaleEN)
	if err != nil && !errors.Is(err, errx.ErrorClassNotFound) {
		return models.ClassWithLocale{}, err
	}
	if err == nil {
		return models.ClassWithLocale{}, errx.ErrorClassAlreadyExists.Raise(
			err,
		)
	}

	ent := class.CreateParams{
		Code: params.Code,
		Icon: params.Icon,
		Name: params.Name,
	}
	if params.Parent != nil {
		ent.Parent = params.Parent
	}

	c, err := a.classificator.Create(ctx, ent)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	return c, err
}
