package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

type SetLocaleParams struct {
	Locale      string
	Name        string
	Description string
}

func (m Service) SetLocales(
	ctx context.Context,
	placeID uuid.UUID,
	locales ...SetLocaleParams,
) error {
	for _, param := range locales {
		err := enum.IsValidLocaleSupported(param.Locale)
		if err != nil {
			return errx.ErrorInvalidLocale.Raise(
				fmt.Errorf("invalid locale provided: %s, cause %w", param.Locale, err),
			)
		}
	}

	_, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return err
	}

	stmts := make([]schemas.PlaceLocale, 0, len(locales))
	for _, param := range locales {
		stmt := schemas.PlaceLocale{
			PlaceID:     placeID,
			Locale:      param.Locale,
			Name:        param.Name,
			Description: param.Description,
		}

		stmts = append(stmts, stmt)
	}

	if len(stmts) == 0 {
		return errx.ErrorNeedAtLeastOneLocaleForPlace.Raise(
			fmt.Errorf("at least one locale must be provided for place %s", placeID),
		)
	}

	err = m.db.PlaceLocales().Upsert(ctx, stmts...)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert locales for place %s, cause: %w", placeID, err),
		)
	}

	return nil
}

func (m Service) ListLocales(
	ctx context.Context,
	placeID uuid.UUID,
	page uint,
	size uint,
) (models.PlaceLocaleCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	_, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceLocaleCollection{}, err
	}

	rows, err := m.db.PlaceLocales().FilterPlaceID(placeID).Page(limit, offset).OrderByLocale(true).Select(ctx)
	if err != nil {
		return models.PlaceLocaleCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list locales for place %s, cause: %w", placeID, err),
		)
	}

	count, err := m.db.PlaceLocales().FilterPlaceID(placeID).Count(ctx)
	if err != nil {
		return models.PlaceLocaleCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to count locales for place %s, cause: %w", placeID, err),
		)
	}

	result := make([]models.PlaceLocale, 0, len(rows))
	for _, loc := range rows {
		result = append(result, localeFromDB(loc))
	}

	return models.PlaceLocaleCollection{
		Data:  result,
		Page:  page,
		Size:  size,
		Total: count,
	}, nil
}

func (m Service) DeleteLocale(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) error {
	_, err := m.Get(ctx, placeID, locale)
	if err != nil {
		return err
	}

	locs, err := m.db.PlaceLocales().FilterPlaceID(placeID).Select(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get locales for place %s, cause: %w", placeID, err),
		)
	}

	if len(locs) == 0 {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("no locales found for place %s", placeID),
		)
	}

	count, err := m.db.PlaceLocales().FilterPlaceID(placeID).FilterByLocale(locale).Count(ctx)
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

	err = m.db.PlaceLocales().FilterPlaceID(placeID).FilterByLocale(locale).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locale %s for place %s, cause: %w", locale, placeID, err),
		)
	}

	return nil
}
