package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
)

func (a App) ListClassLocales(
	ctx context.Context,
	class string,
	pag pagi.Request,
) ([]models.ClassLocale, pagi.Response, error) {
	return a.classificator.LocalesList(ctx, class, pag)
}
