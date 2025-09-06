package entities

import (
	"context"
	"database/sql"
	"time"

	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
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

type CreatePlaceParams struct {
	ID            uuid.UUID
	CityID        uuid.UUID
	DistributorID *uuid.UUID
	Class         string
	Status        string
	Ownership     string
	Point         orb.Point
	Locale        string
	Name          string
	Address       string
	Description   *string
	Website       *string
	Phone         *string
}

func (c Place) CreatePlace(ctx context.Context, params CreatePlaceParams) (Place, error) {
	now := time.Now().UTC()

	stmt := dbx.InsertPlace{
		ID:            params.ID,
		CityID:        params.CityID,
		Class:         params.Class,
		Status:        params.Status,
		Ownership:     params.Ownership,
		Point:         params.Point,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err := c.query.New().Insert(ctx,
}
