package class

import (
	"context"
	"database/sql"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/dbx"
)

type classQ interface {
	New() dbx.ClassesQ

	Insert(ctx context.Context, input dbx.PlaceClass) error
	Get(ctx context.Context) (dbx.PlaceClass, error)
	Select(ctx context.Context) ([]dbx.PlaceClass, error)
	Update(ctx context.Context, input dbx.UpdatePlaceClassParams) error
	Delete(ctx context.Context) error

	FilterCode(code string) dbx.ClassesQ
	FilterParent(parent sql.NullString) dbx.ClassesQ
	FilterParentCycle(parent string) dbx.ClassesQ
	FilterStatus(status string) dbx.ClassesQ

	WithLocale(locale string) dbx.ClassesQ
	SelectWithLocale(ctx context.Context, locale string) ([]dbx.PlaceClassWithLocale, error)
	GetWithLocale(ctx context.Context, locale string) (dbx.PlaceClassWithLocale, error)

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.ClassesQ
}

type classLocaleQ interface {
	New() dbx.ClassLocaleQ

	Insert(ctx context.Context, input ...dbx.PlaceClassLocale) error
	Upsert(ctx context.Context, input ...dbx.PlaceClassLocale) error
	Get(ctx context.Context) (dbx.PlaceClassLocale, error)
	Select(ctx context.Context) ([]dbx.PlaceClassLocale, error)
	Update(ctx context.Context, input dbx.UpdateClassLocaleParams) error
	Delete(ctx context.Context) error

	FilterClass(class string) dbx.ClassLocaleQ
	FilterLocale(locale string) dbx.ClassLocaleQ
	FilterNameLike(name string) dbx.ClassLocaleQ

	OrderByLocale(asc bool) dbx.ClassLocaleQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.ClassLocaleQ
}

type Classificator struct {
	query   classQ
	localeQ classLocaleQ
}

func NewClassificator(db *sql.DB) Classificator {
	return Classificator{
		query:   dbx.NewClassesQ(db),
		localeQ: dbx.NewClassLocaleQ(db),
	}
}

func classWithLocaleModelFromDB(dbClass dbx.PlaceClassWithLocale) models.ClassWithLocale {
	resp := models.Class{
		Code:      dbClass.Code,
		Status:    dbClass.Status,
		Icon:      dbClass.Icon,
		CreatedAt: dbClass.CreatedAt,
		UpdatedAt: dbClass.UpdatedAt,
	}
	if dbClass.Parent.Valid {
		resp.Parent = &dbClass.Parent.String
	}

	return models.ClassWithLocale{
		Data: resp,
		Locale: models.ClassLocale{
			Class:  dbClass.Code,
			Locale: dbClass.Locale,
			Name:   dbClass.Name,
		},
	}
}

func classLocaleModelFromDB(dbLoc dbx.PlaceClassLocale) models.ClassLocale {
	return models.ClassLocale{
		Class:  dbLoc.Class,
		Locale: dbLoc.Locale,
		Name:   dbLoc.Name,
	}
}
