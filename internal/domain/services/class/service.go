package class

import (
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

type Service struct {
	db data.Database
}

func NewService(db data.Database) Service {
	return Service{db: db}
}

func classWithLocaleModelFromDB(dbClass schemas.PlaceClassWithLocale) models.ClassWithLocale {
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

func classLocaleModelFromDB(dbLoc schemas.ClassLocale) models.ClassLocale {
	return models.ClassLocale{
		Class:  dbLoc.Class,
		Locale: dbLoc.Locale,
		Name:   dbLoc.Name,
	}
}
