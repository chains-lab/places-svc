package entities

import (
	"database/sql"

	"github.com/chains-lab/places-svc/internal/dbx"
)

type Place struct {
	query   dbx.PlacesQ
	localeQ dbx.PlaceLocalesQ
}

func NewPlace(db *sql.DB) Place {
	return Place{
		query:   dbx.NewPlacesQ(db),
		localeQ: dbx.NewPlaceLocalesQ(db),
	}
}

func (c Place) CreatePlace() {

}
