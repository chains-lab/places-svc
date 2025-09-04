package app

import (
	"database/sql"

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
