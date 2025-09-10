package place

import (
	"context"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/google/uuid"
)

func (p Place) VerifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Place.Verified {
		return place, nil
	}

	verified := true
	updated, err := p.UpdatePlace(ctx, placeID, constant.LocaleEN, UpdatePlaceParams{
		Verified: &verified,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}
