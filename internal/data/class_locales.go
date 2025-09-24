package data

import (
	"context"
	"time"
)

type ClassLocalesQ interface {
	Insert(ctx context.Context, in ...ClassLocale) error
	Upsert(ctx context.Context, in ...ClassLocale) error
	Update(ctx context.Context, params UpdateClassLocaleParams) error
	Get(ctx context.Context) (ClassLocale, error)
	Select(ctx context.Context) ([]ClassLocale, error)
	Delete(ctx context.Context) error

	FilterClass(class string) ClassLocalesQ
	FilterLocale(locale string) ClassLocalesQ
	FilterNameLike(name string) ClassLocalesQ

	OrderByLocale(asc bool) ClassLocalesQ

	Page(limit, offset uint64) ClassLocalesQ
	Count(ctx context.Context) (uint64, error)
}

type ClassLocale struct {
	Class  string `db:"class"`
	Locale string `db:"locale"`
	Name   string `db:"name"`
}

type UpdateClassLocaleParams struct {
	Name      *string
	UpdatedAt time.Time
}
