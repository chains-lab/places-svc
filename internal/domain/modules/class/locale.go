package class

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (m Module) LocalesList(
	ctx context.Context,
	class string,
	pag pagi.Request,
) ([]models.ClassLocale, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	_, err := m.Get(ctx, class, enum.LocaleEN)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	rows, err := m.db.ClassLocales().FilterClass(class).Page(limit, offset).OrderByLocale(true).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for class %m, cause: %w", class, err),
		)
	}

	count, err := m.db.ClassLocales().FilterClass(class).Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for class %m, cause: %w", class, err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	result := make([]models.ClassLocale, 0, len(rows))
	for _, loc := range rows {
		result = append(result, classLocaleModelFromDB(loc))
	}

	return result, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}

type SetLocaleParams struct {
	Locale string
	Name   string
}

func (m Module) SetLocales(
	ctx context.Context,
	code string,
	locales ...SetLocaleParams,
) error {
	for _, param := range locales {
		err := enum.IsValidLocaleSupported(param.Locale)
		if err != nil {
			return errx.ErrorInvalidLocale.Raise(
				fmt.Errorf("invalid locale provided: %m, cause %w", param.Locale, err),
			)
		}
	}

	_, err := m.Get(ctx, code, enum.DefaultLocale)
	if err != nil {
		return err
	}

	stmts := make([]data.ClassLocale, 0, len(locales))
	for _, locale := range locales {
		stmts = append(stmts, data.ClassLocale{
			Class:  code,
			Locale: locale.Locale,
			Name:   locale.Name,
		})
	}

	if len(stmts) == 0 { //TODO remove all locales before this or not
		return errx.ErrorNedAtLeastOneLocaleForClass.Raise(
			fmt.Errorf("need at least one locale for class"),
		)
	}

	err = m.db.ClassLocales().Upsert(ctx, stmts...)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locales for class %m, cause: %w", code, err),
		)
	}

	return nil
}

func (m Module) DeleteLocale(
	ctx context.Context,
	class, locale string,
) error {
	_, err := m.Get(ctx, class, locale)
	if err != nil {
		return err
	}

	locs, err := m.db.ClassLocales().FilterClass(class).FilterLocale(locale).Select(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get locale %m for class %m, cause: %w", locale, class, err),
		)
	}

	if len(locs) == 0 {
		return errx.ErrorClassLocaleNotFound.Raise(
			fmt.Errorf("locale %m for class %m not found", locale, class),
		)
	}

	count, err := m.db.ClassLocales().FilterClass(class).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for class %m, cause: %w", class, err),
		)
	}
	if count <= 1 {
		return errx.ErrorNedAtLeastOneLocaleForClass.Raise(
			fmt.Errorf("need at least one locale for class"),
		)
	}

	err = m.db.ClassLocales().FilterClass(class).FilterLocale(locale).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locale %m for class %m, cause: %w", locale, class, err),
		)
	}

	return nil
}
