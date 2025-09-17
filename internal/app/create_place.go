package app

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/errx"
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
	Locale        string
	Name          string
	Address       string
	Description   string
}

func (a App) CreatePlace(
	ctx context.Context,
	params CreatePlaceParams,
) (models.PlaceWithDetails, error) {
	p := place.CreateParams{
		ID:            uuid.New(),
		CityID:        params.CityID,
		DistributorID: params.DistributorID,
		Class:         params.Class,
		Point:         params.Point,
		Status:        enum.PlaceStatusActive,
		Address:       params.Address,
		Locale:        params.Locale,
		Name:          params.Name,
		Description:   params.Description,
	}
	if params.Website != nil {
		p.Website = params.Website
	}
	if params.Phone != nil {
		p.Phone = params.Phone
	}

	class, err := a.classificator.Get(ctx, params.Class, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}
	if class.Data.Status != enum.PlaceClassStatusesActive {
		return models.PlaceWithDetails{}, errx.ErrorClassStatusIsNotActive
	}

	var res models.PlaceWithDetails
	txErr := a.transaction(func(txCtx context.Context) error {
		res, err = a.place.Create(ctx, p)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return models.PlaceWithDetails{}, txErr
	}

	return res, nil
}
