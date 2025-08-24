package app

import (
	"database/sql"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/config"
	"github.com/chains-lab/places-svc/internal/dbx"
)

type App struct {
	place dbx.PlacesQ
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		place: dbx.NewPlacesQ(pg),
	}, err
}

func placeModelFromDBX(dbxPlace dbx.Place) models.Place {
	res := models.Place{
		ID:          dbxPlace.ID,
		Type:        dbxPlace.Type,
		Status:      dbxPlace.Status,
		Ownership:   dbxPlace.Ownership,
		Name:        dbxPlace.Name,
		Description: dbxPlace.Description,
		Coords: models.Coords{
			Lon: dbxPlace.Lon,
			Lat: dbxPlace.Lat,
		},
		Address:   dbxPlace.Address,
		Website:   dbxPlace.Website,
		Phone:     dbxPlace.Phone,
		UpdatedAt: dbxPlace.UpdatedAt,
		CreatedAt: dbxPlace.CreatedAt,
	}
	if dbxPlace.DistributorID != nil {
		res.DistributorID = dbxPlace.DistributorID
	}

	return res
}
