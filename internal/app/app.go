package app

import (
	"database/sql"

	"github.com/chains-lab/places-svc/internal/config"
)

type App struct {
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{}, nil
}
