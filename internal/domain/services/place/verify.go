package place

import (
	"context"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (m Service) Verify(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if place.Verified {
		return place, nil
	}

	verified := true
	now := time.Now().UTC()
	err = m.db.Places().FilterID(placeID).Update(ctx, schemas.UpdatePlaceParams{
		Verified:  &verified,
		UpdatedAt: now,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	place.Verified = verified
	place.UpdatedAt = now

	return place, nil
}

func (m Service) Unverify(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	if !place.Verified {
		return place, nil
	}

	verified := false
	now := time.Now().UTC()
	err = m.db.Places().FilterID(placeID).Update(ctx, schemas.UpdatePlaceParams{
		Verified:  &verified,
		UpdatedAt: now,
	})
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	place.Verified = verified
	place.UpdatedAt = now

	return place, nil
}
