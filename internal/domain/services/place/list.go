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
}

type FilterListDistance struct {
	Point   orb.Point
	RadiusM uint64
}

func (m Service) List(
	ctx context.Context,
	locale string,
	filter FilterListParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.PlaceWithDetails, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1 // +1 чтобы определить наличие next
	offset := (pag.Page - 1) * pag.Size

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

	for _, s := range sort {
		switch s.Field {
		case "created_at":
			query = query.OrderByCreatedAt(s.Ascend)
		case "distance":
			if filter.Location != nil {
				query = query.OrderByDistance(filter.Location.Point, s.Ascend)
			}
		}
	}

	err = enum.IsValidLocaleSupported(locale)
	if err != nil {
		locale = enum.LocaleEN
	}

	rows, err := query.Page(limit, offset).SelectWithDetails(ctx, locale)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to search locos, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	places := make([]models.PlaceWithDetails, 0, len(rows))
	for _, row := range rows {
		places = append(places, placeWithDetailsModelFromDB(row))
	}

	return places, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}
