package class

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (c Classificator) LocalesList(
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

	_, err := c.Get(ctx, class, enum.LocaleEN)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	rows, err := c.localeQ.New().FilterClass(class).Page(limit, offset).OrderByLocale(true).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for class %s, cause: %w", class, err),
		)
	}

	count, err := c.localeQ.New().FilterClass(class).Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for class %s, cause: %w", class, err),
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

type SetClassLocaleParams struct {
	Locale string
	Name   string
}

func (c Classificator) SetLocales(
	ctx context.Context,
	code string,
	locales ...SetClassLocaleParams,
) error {
	for _, param := range locales {
		err := enum.IsValidLocaleSupported(param.Locale)
		if err != nil {
			return errx.ErrorInvalidLocale.Raise(
				fmt.Errorf("invalid locale provided: %s, cause %w", param.Locale, err),
			)
		}
	}

	_, err := c.Get(ctx, code, enum.DefaultLocale)
	if err != nil {
		return err
	}

	stmts := make([]dbx.PlaceClassLocale, 0, len(locales))
	for _, locale := range locales {
		stmts = append(stmts, dbx.PlaceClassLocale{
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

	err = c.localeQ.Upsert(ctx, stmts...)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locales for class %s, cause: %w", code, err),
		)
	}

	return nil
}

func (c Classificator) DeleteLocale(
	ctx context.Context,
	class, locale string,
) error {
	_, err := c.Get(ctx, class, locale)
	if err != nil {
		return err
	}

	locs, err := c.localeQ.New().FilterClass(class).FilterLocale(locale).Select(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get locale %s for class %s, cause: %w", locale, class, err),
		)
	}

	if len(locs) == 0 {
		return errx.ErrorClassLocaleNotFound.Raise(
			fmt.Errorf("locale %s for class %s not found", locale, class),
		)
	}

	count, err := c.localeQ.New().FilterClass(class).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for class %s, cause: %w", class, err),
		)
	}
	if count <= 1 {
		return errx.ErrorNedAtLeastOneLocaleForClass.Raise(
			fmt.Errorf("need at least one locale for class"),
		)
	}

	err = c.localeQ.New().FilterClass(class).FilterLocale(locale).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locale %s for class %s, cause: %w", locale, class, err),
		)
	}

	return nil
}
