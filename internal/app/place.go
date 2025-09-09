package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/entities"
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
) (models.PlaceWithLocale, error) {
	ID := uuid.New()

	ents := entities.CreatePlaceParams{
		ID:            ID,
		CityID:        params.CityID,
		DistributorID: params.DistributorID,
		Class:         params.Class,
		Point:         params.Point,
	}
	if params.Website != nil {
		ents.Website = params.Website
	}
	if params.Phone != nil {
		ents.Phone = params.Phone
	}

	_, err := a.classificator.GetClass(ctx, params.Class, constant.LocaleEN)
	if err != nil {
		return models.PlaceWithLocale{}, err
	}

	return a.place.Create(ctx, ents, entities.CreatePlaceLocalParams{
		Locale:      locale.Locale,
		Name:        locale.Name,
		Description: locale.Description,
	})
}

func (a App) GetPlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithLocale, error) {
	return a.place.GetPlaceByID(ctx, placeID, locale)
}

func (a App) GetPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	pag pagi.Request,
) ([]models.PlaceLocale, pagi.Response, error) {
	return a.place.ListLocalesForPlace(ctx, placeID, pag)
}

type SearchPlacesFilter struct {
	Class          []string
	Status         []string
	CityIDs        []uuid.UUID
	DistributorIDs []uuid.UUID
	Verified       *bool
	Name           *string
	Address        *string

	Location *SearchPlaceDistanceFilter
}

type SearchPlaceDistanceFilter struct {
	Point   orb.Point
	RadiusM uint64
}

func (a App) SearchPlaces(
	ctx context.Context,
	locale string,
	filter SearchPlacesFilter,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.PlaceWithLocale, pagi.Response, error) {
	ents := entities.SearchPlacesFilter{}
	if len(filter.Class) > 0 && filter.Class != nil {
		ents.Class = filter.Class
	}
	if len(filter.Status) > 0 && filter.Status != nil {
		ents.Status = filter.Status
	}
	if len(filter.CityIDs) > 0 && filter.CityIDs != nil {
		ents.CityID = filter.CityIDs
	}
	if len(filter.DistributorIDs) > 0 && filter.DistributorIDs != nil {
		ents.DistributorID = filter.DistributorIDs
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

	return a.place.SearchPlaces(ctx, locale, ents, pag, sort)
}

type UpdatePlaceParams struct {
	Class   *string
	Website *string
	Phone   *string
}

func (a App) AddPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	locales ...CreatePlaceLocalParams,
) error {
	_, err := a.place.GetPlaceByID(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return err
	}

	out := make([]entities.AddLocaleParams, 0, len(locales))
	for _, locale := range locales {
		err := constant.IsValidLocaleSupported(locale.Locale)
		if err != nil {
			return err
		}

		s := entities.AddLocaleParams{
			Locale:      locale.Locale,
			Name:        locale.Name,
			Description: locale.Description,
		}

		out = append(out, s)
	}

	err = a.place.AddPlaceLocales(ctx, placeID, out...)
	if err != nil {
		return err
	}

	return nil
}

func (a App) UpdatePlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdatePlaceParams,
) (models.PlaceWithLocale, error) {
	input := entities.UpdatePlaceParams{}
	if params.Class != nil {
		_, err := a.classificator.GetClass(ctx, *params.Class, constant.LocaleEN)
		if err != nil {
			return models.PlaceWithLocale{}, err
		}
		input.Class = params.Class
	}

	place, err := a.place.UpdatePlace(ctx, placeID, locale, entities.UpdatePlaceParams{
		Class:   input.Class,
		Website: params.Website,
		Phone:   params.Phone,
	})
	if err != nil {
		return models.PlaceWithLocale{}, err
	}

	return place, nil
}

func (a App) ReactivatePlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithLocale, error) {
	return a.place.ReactivatePlace(ctx, locale, placeID)
}

func (a App) DeactivatePlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithLocale, error) {
	return a.place.DeactivatePlace(ctx, locale, placeID)
}

func (a App) VerifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithLocale, error) {
	return a.place.VerifyPlace(ctx, placeID)
}

func (a App) DeletePlace(ctx context.Context, placeID uuid.UUID) error {
	return a.place.DeletePlace(ctx, placeID)
}
