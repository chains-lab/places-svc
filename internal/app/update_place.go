package app

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/google/uuid"
)

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
		_, err := a.classificator.Get(ctx, *params.Class, constant.LocaleEN)
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
