package data

import (
	"context"
	"database/sql"
	"time"
)

type ClassesQ interface {
	Insert(ctx context.Context, in PlaceClass) error
	Get(ctx context.Context) (PlaceClass, error)
	Select(ctx context.Context) ([]PlaceClass, error)
	Update(ctx context.Context, in UpdateClassParams) error
	Delete(ctx context.Context) error

	FilterCode(code string) ClassesQ
	FilterParent(parent sql.NullString) ClassesQ
	FilterParentCycle(code string) ClassesQ
	FilterStatus(status string) ClassesQ

	WithLocale(locale string) ClassesQ
	GetWithLocale(ctx context.Context, locale string) (PlaceClassWithLocale, error)
	SelectWithLocale(ctx context.Context, locale string) ([]PlaceClassWithLocale, error)

	OrderBy(orderBy string) ClassesQ
	Page(limit, offset uint64) ClassesQ
	Count(ctx context.Context) (uint64, error)
}

type PlaceClass struct {
	Code      string         `db:"code"`
	Parent    sql.NullString `db:"parent"` // NULL для корней
	Status    string         `db:"status"`
	Icon      string         `db:"icon"`
	Path      string         `db:"path"` // ltree как text
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type PlaceClassWithLocale struct {
	Code      string         `db:"code"`
	Parent    sql.NullString `db:"parent"`
	Status    string         `db:"status"`
	Icon      string         `db:"icon"`
	Path      string         `db:"path"`
	Locale    string         `db:"locale"`
	Name      string         `db:"name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type UpdateClassParams struct {
	Parent    *sql.NullString
	Status    *string
	Icon      *string
	UpdatedAt time.Time
}
