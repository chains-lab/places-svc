package entities

import (
	"context"
	"database/sql"
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
	Upsert(ctx context.Context, in dbx.PlaceLocale) error
	Get(ctx context.Context) (dbx.PlaceLocale, error)
	Select(ctx context.Context) ([]dbx.PlaceLocale, error)
	Delete(ctx context.Context) error

	FilterPlaceID(placeID uuid.UUID) dbx.PlaceLocalesQ
	FilterByLocale(locale string) dbx.PlaceLocalesQ
	FilterByName(name string) dbx.PlaceLocalesQ
	FilterByAddress(address string) dbx.PlaceLocalesQ

	OrderByLocale(asc bool) dbx.PlaceLocalesQ

	Page(limit, offset uint64) dbx.PlaceLocalesQ
	Count(ctx context.Context) (uint64, error)
}

type UpsetLocaleToPlaceParams struct {
	Name        string
	Address     string
	Description *string
}

func (p Place) UpsertLocaleToPlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpsetLocaleToPlaceParams,
) (models.LocaleForPlace, error) {
	err := constant.IsValidLocaleSupported(locale)
	if err != nil {
		return models.LocaleForPlace{}, errx.ErrorNeedAtLeastOneLocaleForPlace.Raise(
			fmt.Errorf("invalid locale provided: %s, cause %w", locale, err),
		)
	}

	place, err := p.GetPlaceByID(ctx, placeID, locale)
	if err != nil {
		return models.LocaleForPlace{}, err
	}

	if place.Locale.Locale != locale {
		return models.LocaleForPlace{}, errx.ErrorCurrentLocaleNotFoundForPlace.Raise(
			fmt.Errorf("current locale %s not found for place %s", locale, placeID),
		)
	}

	stmt := dbx.PlaceLocale{
		PlaceID: placeID,
		Locale:  locale,
		Name:    params.Name,
		Address: params.Address,
	}
	if params.Description != nil {
		switch *params.Description {
		case "":
			stmt.Description = sql.NullString{String: "", Valid: false}
		default:
			stmt.Description = sql.NullString{String: *params.Description, Valid: true}
		}
	}

	err = p.localeQ.Upsert(ctx, stmt)
	if err != nil {
		return models.LocaleForPlace{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locale %s for place %s, cause: %w", locale, placeID, err),
		)
	}

	return models.LocaleForPlace{
		PlaceID:     placeID,
		Locale:      locale,
		Name:        params.Name,
		Address:     params.Address,
		Description: params.Description,
	}, nil
}

func (p Place) ListForPlace(
	ctx context.Context,
	placeID uuid.UUID,
	pag pagi.Request,
) ([]models.LocaleForPlace, pagi.Response, error) {
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

	result := make([]models.LocaleForPlace, 0, len(rows))
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

func placeLocaleModelFromDB(dbLoc dbx.PlaceLocale) models.LocaleForPlace {
	resp := models.LocaleForPlace{
		PlaceID: dbLoc.PlaceID,
		Locale:  dbLoc.Locale,
		Name:    dbLoc.Name,
		Address: dbLoc.Address,
	}
	if dbLoc.Description.Valid {
		resp.Description = &dbLoc.Description.String
	}
	return resp
}
