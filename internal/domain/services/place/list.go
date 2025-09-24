package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterListParams struct {
	Classes        []string
	Statuses       []string
	CityIDs        []uuid.UUID
	DistributorIDs []uuid.UUID
	Verified       *bool
	Name           *string
	Address        *string

	Time     *models.TimeInterval
	Location *FilterListDistance

	Page uint
	Size uint
}

type FilterListDistance struct {
	Point   orb.Point
	RadiusM uint64
}

type SortListField struct {
	ByCreatedAt *bool
	ByDistance  *bool
}

func (m Service) List(
	ctx context.Context,
	locale string,
	filter FilterListParams,
	sort SortListField,
) (models.PlacesCollection, error) {
	limit, offset := pagi.PagConvert(filter.Page, filter.Size)

	query := m.db.Places()

	if filter.Classes != nil && len(filter.Classes) > 0 {
		query = query.FilterClass(filter.Classes...)
	}
	if filter.Statuses != nil && len(filter.Statuses) > 0 {
		query = query.FilterStatus(filter.Statuses...)
	}
	if filter.CityIDs != nil && len(filter.CityIDs) > 0 {
		query = query.FilterCityID(filter.CityIDs...)
	}
	if filter.DistributorIDs != nil && len(filter.DistributorIDs) > 0 {
		query = query.FilterDistributorID(filter.DistributorIDs...)
	}
	if filter.Verified != nil {
		query = query.FilterVerified(*filter.Verified)
	}
	if filter.Name != nil {
		query = query.FilterNameLike(*filter.Name)
	}
	if filter.Address != nil {
		query = query.FilterAddressLike(*filter.Address)
	}
	if filter.Location != nil && filter.Location.RadiusM > 0 {
		query = query.FilterWithinRadiusMeters(filter.Location.Point, filter.Location.RadiusM)
	}
	if filter.Time != nil {
		query = query.FilterTimetableBetween(filter.Time.ToNumberMinutes())
	}
	count, err := query.Count(ctx)

	if sort.ByCreatedAt != nil {
		query = query.OrderByCreatedAt(*sort.ByCreatedAt)
	}
	if sort.ByDistance != nil && filter.Location != nil {
		query = query.OrderByDistance(filter.Location.Point, *sort.ByDistance)
	}

	err = enum.IsValidLocaleSupported(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	rows, err := query.Page(limit, offset).SelectWithDetails(ctx, locale)
	if err != nil {
		return models.PlacesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to search locos, cause: %w", err),
		)
	}

	places := make([]models.Place, 0, len(rows))
	for _, row := range rows {
		places = append(places, modelFromDB(row))
	}

	return models.PlacesCollection{
		Data:  places,
		Page:  filter.Page,
		Size:  filter.Size,
		Total: count,
	}, nil
}
