package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type FilterListParams struct {
	Class         []string
	Status        []string
	CityID        []uuid.UUID
	DistributorID []uuid.UUID
	Verified      *bool
	Name          *string
	Address       *string

	Time     *models.TimeInterval
	Location *FilterListDistance
}

type FilterListDistance struct {
	Point   orb.Point
	RadiusM uint64
}

func (p Place) List(
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

	query := p.query.New()

	if len(filter.Class) > 0 && filter.Class != nil {
		query = query.FilterClass(filter.Class...)
	}
	if len(filter.Status) > 0 && filter.Status != nil {
		query = query.FilterStatus(filter.Status...)
	}
	if filter.Verified != nil {
		query = query.FilterVerified(*filter.Verified)
	}
	if len(filter.CityID) > 0 && filter.CityID != nil {
		query = query.FilterCityID(filter.CityID...)
	}
	if len(filter.DistributorID) > 0 && filter.DistributorID != nil {
		query = query.FilterDistributorID(filter.DistributorID...)
	}
	if filter.Name != nil {
		query = query.FilterNameLike(*filter.Name)
	}
	if filter.Address != nil {
		query = query.FilterAddressLike(*filter.Address)
	}
	if filter.Location.RadiusM > 0 && filter.Location != nil {
		query = query.FilterWithinRadiusMeters(filter.Location.Point, filter.Location.RadiusM)
	}
	if filter.Time != nil {
		query = query.FilterTimetableBetween(filter.Time.ToNumberMinutes())
	}

	count, err := p.query.Count(ctx)

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

	loc := constant.LocaleEN
	err = constant.IsValidLocaleSupported(locale)

	rows, err := query.Page(limit, offset).SelectWithDetails(ctx, loc)
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
