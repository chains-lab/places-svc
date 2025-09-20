package app

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/domain/place"
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

func (a App) GetPlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithDetails, error) {
	return a.place.Get(ctx, placeID, locale)
}

type FilterListPlaces struct {
	Classes        []string
	Statuses       []string
	CityIDs        []uuid.UUID
	DistributorIDs []uuid.UUID
	Verified       *bool
	Name           *string
	Address        *string

	Location *GeoFilterListPlaces
	Time     *models.TimeInterval
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

	if filter.Location != nil && filter.Location.RadiusM > 0 {
		ents.Location = &place.FilterListDistance{
			Point:   filter.Location.Point,
			RadiusM: filter.Location.RadiusM,
		}
	}

	if filter.Time != nil {
		ents.Time = filter.Time
	}

	return a.place.List(ctx, locale, ents, pag, sort)
}

type UpdatePlaceParams struct {
	Class   *string
	Website *string
	Phone   *string
}

func (a App) UpdatePlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdatePlaceParams,
) (models.PlaceWithDetails, error) {
	input := place.UpdatePlaceParams{}
	if params.Class != nil {
		_, err := a.classificator.Get(ctx, *params.Class, enum.LocaleEN)
		if err != nil {
			return models.PlaceWithDetails{}, err
		}
		input.Class = params.Class
	}

	p, err := a.place.UpdatePlace(ctx, placeID, locale, place.UpdatePlaceParams{
		Class:   input.Class,
		Website: params.Website,
		Phone:   params.Phone,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return p, nil
}

func (a App) DeletePlace(ctx context.Context, placeID uuid.UUID) error {
	return a.place.DeletePlace(ctx, placeID)
}
