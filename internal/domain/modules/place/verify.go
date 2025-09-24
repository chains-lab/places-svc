package place

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (m Module) Verify(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Verified {
		return place, nil
	}

	verified := true
	updated, err := m.Update(ctx, placeID, enum.LocaleEN, UpdateParams{
		Verified: &verified,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}

func (m Module) Unverify(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if !place.Verified {
		return place, nil
	}

	verified := false
	updated, err := m.Update(ctx, placeID, enum.LocaleEN, UpdateParams{
		Verified: &verified,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	return updated, nil
}
