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

func modelFromDB(dbClass schemas.Class) models.Class {
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
	resp.Name = dbClass.Name

	return resp
}
