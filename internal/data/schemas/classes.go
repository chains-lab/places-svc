package schemas

import (
	"context"
	"database/sql"
	"time"
)

type ClassesQ interface {
	Insert(ctx context.Context, in Class) error
	Get(ctx context.Context) (Class, error)
	Select(ctx context.Context) ([]Class, error)
	Update(ctx context.Context, in UpdateClassParams) error
	Delete(ctx context.Context) error

	FilterCode(code string) ClassesQ
	FilterParent(parent sql.NullString) ClassesQ
	FilterParentCycle(code string) ClassesQ
	FilterStatus(status string) ClassesQ
	FilterName(name string) ClassesQ
	FilterNameLike(name string) ClassesQ

	OrderBy(orderBy string) ClassesQ
	Page(limit, offset uint) ClassesQ
	Count(ctx context.Context) (uint, error)
}

type Class struct {
	Code      string         `storage:"code"`
	Parent    sql.NullString `storage:"parent"` // NULL для корней
	Status    string         `storage:"status"`
	Icon      string         `storage:"icon"`
	Name      string         `storage:"name"`
	Path      string         `storage:"path"` // ltree как text
	CreatedAt time.Time      `storage:"created_at"`
	UpdatedAt time.Time      `storage:"updated_at"`
}

type UpdateClassParams struct {
	Parent    *sql.NullString
	Status    *string
	Icon      *string
	Name      *string
	UpdatedAt time.Time
}
