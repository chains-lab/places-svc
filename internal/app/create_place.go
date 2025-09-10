package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreatePlaceParams struct {
	CityID        uuid.UUID
	DistributorID *uuid.UUID
	Class         string
	Website       *string
	Phone         *string
	Point         orb.Point
}

type CreatePlaceLocalParams struct {
	Locale      string
	Name        string
	Address     string
	Description string
}

func (a App) CreatePlace(
	ctx context.Context,
	params CreatePlaceParams,
	locale CreatePlaceLocalParams,
) (models.PlaceWithDetails, error) {
	p := place.CreateParams{
		ID:            uuid.New(),
		CityID:        params.CityID,
		DistributorID: params.DistributorID,
		Class:         params.Class,
		Point:         params.Point,
	}
	if params.Website != nil {
		p.Website = params.Website
	}
	if params.Phone != nil {
		p.Phone = params.Phone
	}

	_, err := a.classificator.Get(ctx, params.Class, constant.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return a.place.Create(ctx, p, place.CreateLocalParams{
		Locale:      locale.Locale,
		Name:        locale.Name,
		Description: locale.Description,
	})
}
