package app

import (
	"database/sql"

	"github.com/chains-lab/places-svc/internal/app/entities"
	"github.com/chains-lab/places-svc/internal/config"
)

type App struct {
	place         entities.Place
	classificator entities.Classificator
	timetable     entities.Timetable
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		place:         entities.NewPlace(pg),
		classificator: entities.NewClassificator(pg),
		timetable:     entities.NewTimetable(pg),
	}, err
}
