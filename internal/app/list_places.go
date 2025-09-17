package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterListPlaces struct {
	Classes        []string
	Statuses       []string
	CityIDs        []uuid.UUID
	DistributorIDs []uuid.UUID
	Verified       *bool
	Name           *string
	Address        *string

	Location *GeoFilterListPlaces
}

type GeoFilterListPlaces struct {
	Point   orb.Point
	RadiusM uint64
}

func (a App) ListPlaces(
	ctx context.Context,
	locale string,
	filter FilterListPlaces,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.PlaceWithDetails, pagi.Response, error) {
	ents := place.FilterListParams{}
	if filter.Classes != nil && len(filter.Classes) > 0 {
		ents.Classes = filter.Classes
	}
	if filter.Statuses != nil && len(filter.Statuses) > 0 {
		ents.Statuses = filter.Statuses
	}
	if filter.CityIDs != nil && len(filter.CityIDs) > 0 {
		ents.CityIDs = filter.CityIDs
	}
	if filter.DistributorIDs != nil && len(filter.DistributorIDs) > 0 {
		ents.DistributorIDs = filter.DistributorIDs
	}

	if filter.Verified != nil {
		ents.Verified = filter.Verified
	}
	if filter.Name != nil {
		ents.Name = filter.Name
	}
	if filter.Address != nil {
		ents.Address = filter.Address
	}

	if filter.Location != nil {
		ents.Location.Point = filter.Location.Point
		ents.Location.RadiusM = filter.Location.RadiusM
	}

	return a.place.List(ctx, locale, ents, pag, sort)
}
