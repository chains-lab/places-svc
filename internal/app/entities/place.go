package entities

import (
	"github.com/chains-lab/places-svc/internal/dbx"
)

type Place struct {
	placeQ   dbx.PlacesQ
	typesQ   dbx.TypesQ
	detailsQ dbx.DetailsQ
}
