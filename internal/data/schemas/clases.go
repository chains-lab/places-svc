package schemas

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
	Page(limit, offset uint) ClassesQ
	Count(ctx context.Context) (uint, error)
}

type PlaceClass struct {
	Code      string         `storage:"code"`
	Parent    sql.NullString `storage:"parent"` // NULL для корней
	Status    string         `storage:"status"`
	Icon      string         `storage:"icon"`
	Path      string         `storage:"path"` // ltree как text
	CreatedAt time.Time      `storage:"created_at"`
	UpdatedAt time.Time      `storage:"updated_at"`
}

type PlaceClassWithLocale struct {
	Code      string         `storage:"code"`
	Parent    sql.NullString `storage:"parent"`
	Status    string         `storage:"status"`
	Icon      string         `storage:"icon"`
	Path      string         `storage:"path"`
	Locale    string         `storage:"locale"`
	Name      string         `storage:"name"`
	CreatedAt time.Time      `storage:"created_at"`
	UpdatedAt time.Time      `storage:"updated_at"`
}

type UpdateClassParams struct {
	Parent    *sql.NullString
	Status    *string
	Icon      *string
	UpdatedAt time.Time
}
