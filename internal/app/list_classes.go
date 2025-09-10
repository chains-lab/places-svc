package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/entities/class"
	"github.com/chains-lab/places-svc/internal/app/models"
)

type FilterListClassesParams struct {
	Parent      *string
	ParentCycle *bool
	Status      *string
}

func (a App) ListClasses(
	ctx context.Context,
	locale string,
	filter FilterListClassesParams,
	pag pagi.Request,
) ([]models.ClassWithLocale, pagi.Response, error) {
	ent := class.FilterListParams{}
	if filter.Parent != nil {
		ent.Parent = filter.Parent
	}
	if filter.ParentCycle != nil {
		ent.ParentCycle = filter.ParentCycle
	}
	if filter.Status != nil {
		ent.Status = filter.Status
	}

	return a.classificator.List(ctx, locale, ent, pag)
}
