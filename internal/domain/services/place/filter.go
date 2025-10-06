package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterParams struct {
	Classes       []string
	Statuses      []string
	CityID        *uuid.UUID
	DistributorID *uuid.UUID
	Verified      *bool
	Name          *string
	Address       *string

	Time     *models.TimeInterval
	Location *FilterDistance
}

type FilterDistance struct {
	Point   orb.Point
	RadiusM uint64
}

type SortParams struct {
	ByCreatedAt *bool
	ByDistance  *bool
}

func (s Service) Filter(
	ctx context.Context,
	locale string,
	filter FilterParams,
	sort SortParams,
	page, size uint64,
) (models.PlacesCollection, error) {
	rows, err := s.db.FilterPlaces(ctx, locale, filter, sort, page, size)
	if err != nil {
		return models.PlacesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to list places, cause: %w", err),
		)
	}

	return rows, nil
}
