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

func modelFromDB(dbClass schemas.PlaceClassWithLocale) models.Class {
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
	resp.Locale = dbClass.Locale
	resp.Name = dbClass.Name

	return resp
}

func localeFromDB(dbLoc schemas.ClassLocale) models.ClassLocale {
	return models.ClassLocale{
		Class:  dbLoc.Class,
		Locale: dbLoc.Locale,
		Name:   dbLoc.Name,
	}
}

func detailsFromDB(dbClass schemas.PlaceClass) models.ClassDetails {
	resp := models.ClassDetails{
		Code:      dbClass.Code,
		Status:    dbClass.Status,
		Icon:      dbClass.Icon,
		CreatedAt: dbClass.CreatedAt,
		UpdatedAt: dbClass.UpdatedAt,
	}
	if dbClass.Parent.Valid {
		resp.Parent = &dbClass.Parent.String
	}
	return resp
}
