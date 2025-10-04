package place

import (
	"context"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Verify(ctx context.Context, placeID uuid.UUID, value bool) (models.Place, error) {
	place, err := s.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return models.Place{}, err
	}

	if place.Verified == value {
		return place, nil
	}

	now := time.Now().UTC()
	err = s.db.UpdateVerifiedPlace(ctx, placeID, value, now)
	if err != nil {
		return models.Place{}, err
	}

	place.Verified = true
	place.UpdatedAt = now

	return place, nil
}
