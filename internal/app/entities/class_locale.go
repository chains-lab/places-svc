package entities

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
)

type ClassLocaleQ interface {
	New() dbx.ClassLocaleQ

	Insert(ctx context.Context, input ...dbx.PlaceClassLocale) error
	Upsert(ctx context.Context, input dbx.PlaceClassLocale) error
	Get(ctx context.Context) (dbx.PlaceClassLocale, error)
	Select(ctx context.Context) ([]dbx.PlaceClassLocale, error)
	Update(ctx context.Context, input dbx.UpdateClassLocaleParams) error
	Delete(ctx context.Context) error

	FilterClass(class string) dbx.ClassLocaleQ
	FilterLocale(locale string) dbx.ClassLocaleQ
	FilterNameLike(name string) dbx.ClassLocaleQ

	OrderByLocale(asc bool) dbx.ClassLocaleQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.ClassLocaleQ
}

func (c Classificator) UpsetLocaleToClass(
	ctx context.Context,
	codeClass, locale, name string,
) (models.ClassLocale, error) {
	class, err := c.GetClass(ctx, codeClass, locale)
	if err != nil {
		return models.ClassLocale{}, err
	}

	if class.Locale.Locale != locale {
		return models.ClassLocale{}, errx.ErrorCurrentLocaleNotFoundForClass.Raise(
			fmt.Errorf("current locale %s not found for class %s", locale, codeClass),
		)
	}

	err = c.localeQ.Upsert(ctx, dbx.PlaceClassLocale{
		Class:  codeClass,
		Locale: locale,
		Name:   name,
	})
	if err != nil {
		return models.ClassLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locale %s for codeClass %s, cause: %w", locale, codeClass, err),
		)
	}

	return models.ClassLocale{
		Class:  codeClass,
		Locale: locale,
		Name:   name,
	}, nil
}

func (c Classificator) ListForClass(
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

	_, err := c.GetClass(ctx, class, constant.LocaleEN)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	rows, err := c.localeQ.New().FilterClass(class).Page(limit, offset).OrderByLocale(true).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for class %s, cause: %w", class, err),
		)
	}

	count, err := c.localeQ.New().FilterLocale(class).Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for class %s, cause: %w", class, err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:limit]
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

func (c Classificator) DeleteLocaleFromClass(
	ctx context.Context,
	class, locale string,
) error {
	_, err := c.GetClass(ctx, class, locale)
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

func classLocaleModelFromDB(dbLoc dbx.PlaceClassLocale) models.ClassLocale {
	return models.ClassLocale{
		Class:  dbLoc.Class,
		Locale: dbLoc.Locale,
		Name:   dbLoc.Name,
	}
}
