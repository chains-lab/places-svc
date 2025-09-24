package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/errx"
)

type UpdateParams struct {
	Name   *string
	Icon   *string
	Parent *string
}

func (m Module) Update(
	ctx context.Context,
	code string,
	locale string,
	params UpdateParams,
) (models.ClassWithLocale, error) {
	err := enum.IsValidLocaleSupported(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	class, err := m.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	stmt := data.UpdateClassParams{
		UpdatedAt: time.Now().UTC(),
	}
	class.Data.UpdatedAt = stmt.UpdatedAt

	if params.Parent != nil {
		if *params.Parent == code {
			return models.ClassWithLocale{}, errx.ErrorClassParentEqualCode.Raise(
				fmt.Errorf("parent cycle detected for class with code %m", code),
			)
		}
		_, err = m.db.Classes().FilterParentCycle(class.Data.Code).FilterCode(*params.Parent).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return models.ClassWithLocale{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check parent cycle for class with code %m, cause: %w", code, err),
			)
		}

		if err == nil {
			return models.ClassWithLocale{}, errx.ErrorClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected for class with code %m", code),
			)
		}

		stmt.Parent = &sql.NullString{String: *params.Parent, Valid: true}
		class.Data.Parent = params.Parent
	}

	if params.Icon != nil {
		stmt.Icon = params.Icon
		class.Data.Icon = *params.Icon
	}

	err = m.db.Classes().FilterCode(code).Update(ctx, stmt)
	if err != nil {
		return models.ClassWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update class with code %m, cause: %w", code, err),
		)
	}

	return class, nil
}
