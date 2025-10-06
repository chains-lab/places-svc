package place

import (
	"context"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (s Service) Block(ctx context.Context, placeID uuid.UUID, locale string, block bool) (models.Place, error) {
	place, err := s.Get(ctx, placeID, locale)
	if err != nil {
		return models.Place{}, err
	}

	place.UpdatedAt = time.Now().UTC()
	var status string

	if block {
		status = enum.PlaceStatusBlocked
	} else {
		status = enum.PlaceStatusInactive
	}

	err = s.db.UpdatePlaceStatus(ctx, placeID, status, place.UpdatedAt)
	if err != nil {
		return models.Place{}, err
	}

	return place, nil
}
