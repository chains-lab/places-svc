package data

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

	Page(limit, offset uint64) PlacesQ
	Count(ctx context.Context) (uint64, error)
}

type Place struct {
	ID            uuid.UUID     `db:"id"`
	CityID        uuid.UUID     `db:"city_id"`
	DistributorID uuid.NullUUID `db:"distributor_id"`
	Class         string        `db:"Class"`

	Status   string    `db:"Status"`
	Verified bool      `db:"Verified"`
	Point    orb.Point `db:"Point"`
	Address  string    `db:"Address"`

	Website sql.NullString `db:"Website"`
	Phone   sql.NullString `db:"Phone"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
