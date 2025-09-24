package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

func (m Service) LocalesList(
	ctx context.Context,
	class string,
	page uint,
	size uint,
) (models.ClassLocaleCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	_, err := m.Get(ctx, class, enum.LocaleEN)
	if err != nil {
		return models.ClassLocaleCollection{}, err
	}

	rows, err := m.db.ClassLocales().FilterClass(class).Page(limit, offset).OrderByLocale(true).Select(ctx)
	if err != nil {
		return models.ClassLocaleCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for class %s, cause: %w", class, err),
		)
	}

	count, err := m.db.ClassLocales().FilterClass(class).Count(ctx)
	if err != nil {
		return models.ClassLocaleCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for class %s, cause: %w", class, err),
		)
	}

	result := make([]models.ClassLocale, 0, len(rows))
	for _, loc := range rows {
		result = append(result, localeFromDB(loc))
	}

	return models.ClassLocaleCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: count,
	}, nil
}

type SetLocaleParams struct {
	Locale string
	Name   string
}

func (m Service) SetLocale(
	ctx context.Context,
	code string,
	loc SetLocaleParams,
) error {
	err := enum.IsValidLocaleSupported(loc.Locale)
	if err != nil {
		return errx.ErrorInvalidLocale.Raise(
			fmt.Errorf("invalid locale provided: %s, cause %w", loc.Locale, err),
		)
	}

	_, err = m.Get(ctx, code, enum.DefaultLocale)
	if err != nil {
		return err
	}

	_, err = m.db.ClassLocales().FilterLocale(loc.Locale).FilterName(loc.Name).Count(ctx)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check locale name uniqueness for class %s, cause: %w", code, err),
			)
		}
		if err == nil {
			return errx.ErrorClassNameAlreadyTaken.Raise(
				fmt.Errorf("locale name %s already taken for locale %s", loc.Name, loc.Locale),
			)
		}
	}

	stmts := schemas.ClassLocale{
		Class:  code,
		Locale: loc.Locale,
		Name:   loc.Name,
	}

	err = m.db.ClassLocales().Upsert(ctx, stmts)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert loc for class %s, cause: %w", code, err),
		)
	}

	return nil
}

func (m Service) DeleteLocale(
	ctx context.Context,
	class, locale string,
) error {
	_, err := m.Get(ctx, class, locale)
	if err != nil {
		return err
	}

	err = enum.IsValidLocaleSupported(locale)
	if err != nil {
		return errx.ErrorInvalidLocale.Raise(
			fmt.Errorf("invalid locale provided: %s, cause %w", locale, err),
		)
	}

	if locale == enum.DefaultLocale {
		return errx.ErrorCannotDeleteDefaultLocaleForClass.Raise(
			fmt.Errorf("cannot delete default locale for class"),
		)
	}

	err = m.db.ClassLocales().FilterClass(class).FilterLocale(locale).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locale %s for class %s, cause: %w", locale, class, err),
		)
	}

	return nil
}
