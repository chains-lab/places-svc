package schemas

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
	Page(limit, offset uint) PlaceLocalesQ
	Count(ctx context.Context) (uint, error)
}

type PlaceLocale struct {
	PlaceID     uuid.UUID `storage:"place_id"`
	Locale      string    `storage:"locale"`
	Name        string    `storage:"name"`
	Description string    `storage:"description"`
}

type UpdatePlaceLocaleParams struct {
	Name        *string
	Description *string
}
