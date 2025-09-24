package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

func (m Service) Get(
	ctx context.Context,
	code, locale string,
) (models.Class, error) {
	err := enum.IsValidLocaleSupported(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	class, err := m.db.Classes().FilterCode(code).GetWithLocale(ctx, locale)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Class{}, errx.ErrorClassNotFound.Raise(
				fmt.Errorf("class with code %s not found, cause: %w", code, err),
			)
		default:
			return models.Class{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get class with code %s, cause: %w", code, err),
			)
		}
	}

	return modelFromDB(class), nil
}
