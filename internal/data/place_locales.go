package data

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/data/pgdb"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
	"github.com/google/uuid"
)

func (d Database) UpsertLocaleForPlace(ctx context.Context, placeID uuid.UUID, locales ...plocale.SetParams) error {
	if len(locales) == 0 {
		return nil
	}

	schemas := make([]pgdb.PlaceLocale, 0, len(locales))
	for _, locale := range locales {
		schemas = append(schemas, placeLocaleSetParamsToSchema(models.PlaceLocale{
			PlaceID:     placeID,
			Locale:      locale.Locale,
			Name:        locale.Name,
			Description: locale.Description,
		}))
	}

	return d.sql.pLocales.Upsert(ctx, schemas...)
}

func (d Database) GetPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	page, size uint64,
) (models.PlaceLocaleCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	schemas, err := d.sql.pLocales.New().FilterPlaceID(placeID).Page(limit, offset).Select(ctx)
	if err != nil {
		return models.PlaceLocaleCollection{}, err
	}

	total, err := d.sql.pLocales.New().FilterPlaceID(placeID).Count(ctx)
	if err != nil {
		return models.PlaceLocaleCollection{}, err
	}

	locales := make([]models.PlaceLocale, 0, len(schemas))
	for _, schema := range schemas {
		locales = append(locales, models.PlaceLocale{
			PlaceID:     schema.PlaceID,
			Locale:      schema.Locale,
			Name:        schema.Name,
			Description: schema.Description,
		})
	}

	return models.PlaceLocaleCollection{
		Data:  locales,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d Database) CountPlaceLocales(ctx context.Context, placeID uuid.UUID) (uint64, error) {
	return d.sql.pLocales.New().FilterPlaceID(placeID).Count(ctx)
}

func (d Database) DeletePlaceLocale(ctx context.Context, placeID uuid.UUID, locale string) error {
	return d.sql.pLocales.New().FilterPlaceID(placeID).FilterLocale(locale).Delete(ctx)
}

func (d Database) CreatePlaceLocale(ctx context.Context, input models.PlaceLocale) error {
	return d.sql.pLocales.Insert(ctx, placeLocaleSetParamsToSchema(input))
}

func placeLocaleSetParamsToSchema(input models.PlaceLocale) pgdb.PlaceLocale {
	return pgdb.PlaceLocale{
		PlaceID:     input.PlaceID,
		Locale:      input.Locale,
		Name:        input.Name,
		Description: input.Description,
	}
}
