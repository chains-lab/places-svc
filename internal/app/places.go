package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreatePlaceParams struct {
	ID            uuid.UUID
	CityID        uuid.UUID
	DistributorID *uuid.UUID
	Class         string
	Status        string
	Website       *string
	Phone         *string
	Ownership     string
	Point         orb.Point
}

type CreatePlaceLocalParams struct {
	Locale      string
	Name        string
	Address     string
	Description *string
}

func (a App) CreatePlace(
	ctx context.Context,
	params CreatePlaceParams,
	locale CreatePlaceLocalParams,
) (models.PlaceWithLocale, error) {
	ents := entities.CreatePlaceParams{
		ID:            params.ID,
		CityID:        params.CityID,
		DistributorID: params.DistributorID,
		Class:         params.Class,
		Status:        params.Status,
		Ownership:     params.Ownership,
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
		Address:     locale.Address,
		Description: locale.Description,
	})
}

func (a App) AddPlaceLocales(
	ctx context.Context,
	placeID uuid.UUID,
	locales ...CreatePlaceLocalParams,
) (models.LocaleForPlace, error) {
	_, err := a.place.GetPlaceByID(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return models.LocaleForPlace{}, err
	}

	out := make([]entities.AddLocaleParams, 0, len(locales))
	for _, locale := range locales {
		err := constant.IsValidLocaleSupported(locale.Locale)
		if err != nil {
			return models.LocaleForPlace{}, err
		}

		s := entities.AddLocaleParams{
			Locale:  locale.Locale,
			Name:    locale.Name,
			Address: locale.Address,
		}
		if locale.Description != nil {
			s.Description = locale.Description
		}

		out = append(out, s)
	}

	return a.place.AddPlaceLocales(ctx, placeID, out...)
}

type UpdatePlaceParams struct {
	Class     *string
	Ownership *string
	Point     *orb.Point
	Website   *string
	Phone     *string
}

func (a App) UpdatePlace(
	ctx context.Context,
	placeID uuid.UUID,
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
	if params.Ownership != nil {
		err := constant.IsValidOwnership(*params.Ownership)
		if err != nil {
			return models.PlaceWithLocale{}, err
		}
		input.Ownership = params.Ownership
	}
	_, err := a.place.UpdatePlace(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return models.PlaceWithLocale{}, err
	}
}
