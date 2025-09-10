package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
)

func (a App) GetClass(
	ctx context.Context,
	code, locale string,
) (models.ClassWithLocale, error) {
	return a.classificator.Get(ctx, code, locale)
}
