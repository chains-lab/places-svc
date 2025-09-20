package place

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (p Place) VerifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Verified {
		return place, nil
	}

	verified := true
	updated, err := p.UpdatePlace(ctx, placeID, enum.LocaleEN, UpdatePlaceParams{
		Verified: &verified,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}

func (p Place) UnverifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if !place.Verified {
		return place, nil
	}

	verified := false
	updated, err := p.UpdatePlace(ctx, placeID, enum.LocaleEN, UpdatePlaceParams{
		Verified: &verified,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}
