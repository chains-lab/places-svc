package schemas

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type PlacesQ interface {
	Insert(ctx context.Context, input Place) error
	Get(ctx context.Context) (Place, error)
	Select(ctx context.Context) ([]Place, error)
	Delete(ctx context.Context) error

	Update(ctx context.Context, params UpdatePlaceParams) error

	FilterID(id uuid.UUID) PlacesQ
	FilterClass(class ...string) PlacesQ
	FilterStatus(status ...string) PlacesQ
	FilterCityID(cityID ...uuid.UUID) PlacesQ
	FilterDistributorID(distributorID ...uuid.UUID) PlacesQ
	FilterVerified(verified bool) PlacesQ
	FilterNameLike(name string) PlacesQ
	FilterAddressLike(address string) PlacesQ
	FilterWithinRadiusMeters(point orb.Point, radiusM uint64) PlacesQ
	FilterWithinBBox(minLon, minLat, maxLon, maxLat float64) PlacesQ
	FilterWithinPolygonWKT(polyWKT string) PlacesQ
	FilterTimetableBetween(start, end int) PlacesQ

	GetWithDetails(ctx context.Context, locale string) (PlaceWithDetails, error)
	SelectWithDetails(ctx context.Context, locale string) ([]PlaceWithDetails, error)
	WithTimetable() PlacesQ
	WithLocale(locale string) PlacesQ

	OrderByCreatedAt(ascend bool) PlacesQ
	OrderByDistance(point orb.Point, ascend bool) PlacesQ

	Page(limit, offset uint) PlacesQ
	Count(ctx context.Context) (uint, error)
}

type Place struct {
	ID            uuid.UUID     `storage:"id"`
	CityID        uuid.UUID     `storage:"city_id"`
	DistributorID uuid.NullUUID `storage:"distributor_id"`
	Class         string        `storage:"Class"`

	Status   string    `storage:"Status"`
	Verified bool      `storage:"Verified"`
	Point    orb.Point `storage:"Point"`
	Address  string    `storage:"Address"`

	Website sql.NullString `storage:"Website"`
	Phone   sql.NullString `storage:"Phone"`

	CreatedAt time.Time `storage:"created_at"`
	UpdatedAt time.Time `storage:"updated_at"`
}

type PlaceWithDetails struct {
	Place
	Locale      string
	Name        string
	Description string
	Timetable   []PlaceTimetable
}

type UpdatePlaceParams struct {
	Class     *string
	Status    *string
	Verified  *bool
	Point     *orb.Point
	Address   *string
	Website   *sql.NullString
	Phone     *sql.NullString
	UpdatedAt time.Time
}
