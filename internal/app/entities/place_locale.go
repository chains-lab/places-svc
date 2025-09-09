package entities

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

type placeLocaleQ interface {
	New() dbx.PlaceLocalesQ

	Insert(ctx context.Context, in dbx.PlaceLocale) error
	Update(ctx context.Context, params dbx.UpdatePlaceLocaleParams) error
	Upsert(ctx context.Context, in ...dbx.PlaceLocale) error
	Get(ctx context.Context) (dbx.PlaceLocale, error)
	Select(ctx context.Context) ([]dbx.PlaceLocale, error)
	Delete(ctx context.Context) error

	FilterPlaceID(placeID uuid.UUID) dbx.PlaceLocalesQ
	FilterByLocale(locale string) dbx.PlaceLocalesQ
	FilterByName(name string) dbx.PlaceLocalesQ

	OrderByLocale(asc bool) dbx.PlaceLocalesQ

	Page(limit, offset uint64) dbx.PlaceLocalesQ
	Count(ctx context.Context) (uint64, error)
}

type AddLocaleParams struct {
	Locale      string
	Name        string
	Description string
}

func (p Place) AddPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	locales ...AddLocaleParams,
) error {
	for _, param := range locales {
		err := constant.IsValidLocaleSupported(param.Locale)
		if err != nil {
			return errx.ErrorNeedAtLeastOneLocaleForPlace.Raise(
				fmt.Errorf("invalid locale provided: %s, cause %w", param.Locale, err),
			)
		}
	}

	_, err := p.GetPlaceByID(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return err
	}

	stmts := make([]dbx.PlaceLocale, 0, len(locales))
	for _, param := range locales {
		stmt := dbx.PlaceLocale{
			PlaceID:     placeID,
			Locale:      param.Locale,
			Name:        param.Name,
			Description: param.Description,
		}

		stmts = append(stmts, stmt)
	}

	err = p.localeQ.Upsert(ctx, stmts...)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locales for place %s, cause: %w", placeID, err),
		)
	}

	return nil
}

func (p Place) ListLocalesForPlace(
	ctx context.Context,
	placeID uuid.UUID,
	pag pagi.Request,
) ([]models.PlaceLocale, pagi.Response, error) {
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

	_, err := p.GetPlaceByID(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return nil, pagi.Response{}, err
	}

	rows, err := p.localeQ.New().FilterPlaceID(placeID).Page(limit, offset).OrderByLocale(true).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for place %s, cause: %w", placeID, err),
		)
	}

	count, err := p.localeQ.New().FilterPlaceID(placeID).Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for place %s, cause: %w", placeID, err),
		)
	}

	if len(rows) > int(limit) {
		rows = rows[:limit]
	}

	result := make([]models.PlaceLocale, 0, len(rows))
	for _, loc := range rows {
		result = append(result, placeLocaleModelFromDB(loc))
	}

	return result, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}

func (p Place) DeleteLocale(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) error {
	_, err := p.GetPlaceByID(ctx, placeID, locale)
	if err != nil {
		return err
	}

	locs, err := p.localeQ.New().FilterPlaceID(placeID).Select(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get locales for place %s, cause: %w", placeID, err),
		)
	}

	if len(locs) == 0 {
		return errx.ErrorPlaceLocaleNotFound.Raise(
			fmt.Errorf("no locales found for place %s", placeID),
		)
	}

	count, err := p.localeQ.New().FilterPlaceID(placeID).FilterByLocale(locale).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for place %s, cause: %w", placeID, err),
		)
	}
	if count <= 1 {
		return errx.ErrorNeedAtLeastOneLocaleForPlace.Raise(
			fmt.Errorf("cannot delete the last locale for place %s", placeID),
		)
	}

	err = p.localeQ.New().FilterPlaceID(placeID).FilterByLocale(locale).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locale %s for place %s, cause: %w", locale, placeID, err),
		)
	}

	return nil
}

func placeLocaleModelFromDB(dbLoc dbx.PlaceLocale) models.PlaceLocale {
	return models.PlaceLocale{
		PlaceID:     dbLoc.PlaceID,
		Locale:      dbLoc.Locale,
		Name:        dbLoc.Name,
		Description: dbLoc.Description,
	}

}
