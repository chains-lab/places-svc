package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/class"
	"github.com/chains-lab/places-svc/internal/app/models"
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
