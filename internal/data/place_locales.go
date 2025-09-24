package data

import (
	"context"

	"github.com/google/uuid"
)

type PlaceLocalesQ interface {
	New() PlaceLocalesQ

	Insert(ctx context.Context, in PlaceLocale) error
	Upsert(ctx context.Context, in ...PlaceLocale) error
	Get(ctx context.Context) (PlaceLocale, error)
	Select(ctx context.Context) ([]PlaceLocale, error)

	Delete(ctx context.Context) error

	FilterPlaceID(id uuid.UUID) PlaceLocalesQ
	FilterByLocale(locale string) PlaceLocalesQ
	FilterByName(name string) PlaceLocalesQ

	OrderByLocale(asc bool) PlaceLocalesQ
	Page(limit, offset uint64) PlaceLocalesQ
	Count(ctx context.Context) (uint64, error)
}

type PlaceLocale struct {
	PlaceID     uuid.UUID `db:"place_id"`
	Locale      string    `db:"locale"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}

type UpdatePlaceLocaleParams struct {
	Name        *string
	Description *string
}
