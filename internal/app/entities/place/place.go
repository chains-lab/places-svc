package place

import (
	"context"
	"database/sql"

	"github.com/chains-lab/places-svc/internal/app/geo"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type placeQ interface {
	New() dbx.PlacesQ

	Insert(ctx context.Context, input dbx.Place) error
	Get(ctx context.Context) (dbx.Place, error)
	Select(ctx context.Context) ([]dbx.Place, error)
	Update(ctx context.Context, input dbx.UpdatePlaceParams) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.PlacesQ
	FilterClass(class ...string) dbx.PlacesQ
	FilterStatus(status ...string) dbx.PlacesQ
	FilterCityID(cityID ...uuid.UUID) dbx.PlacesQ
	FilterDistributorID(distributorID ...uuid.UUID) dbx.PlacesQ
	FilterVerified(verified bool) dbx.PlacesQ
	FilterNameLike(name string) dbx.PlacesQ
	FilterAddressLike(address string) dbx.PlacesQ
	FilterWithinRadiusMeters(point orb.Point, radiusM uint64) dbx.PlacesQ
	FilterWithinBBox(minLon, minLat, maxLon, maxLat float64) dbx.PlacesQ
	FilterWithinPolygonWKT(polyWKT string) dbx.PlacesQ
	FilterTimetableBetween(start, end int) dbx.PlacesQ

	GetWithDetails(ctx context.Context, locale string) (dbx.PlaceWithDetails, error)
	SelectWithDetails(ctx context.Context, locale string) ([]dbx.PlaceWithDetails, error)

	OrderByCreatedAt(ascend bool) dbx.PlacesQ
	OrderByDistance(point orb.Point, ascend bool) dbx.PlacesQ

	Page(limit, offset uint64) dbx.PlacesQ
	Count(ctx context.Context) (uint64, error)
}

type Place struct {
	query     placeQ
	locale    placeLocaleQ
	timetable timetableQ

	geo *geo.Guesser
}

func NewPlace(db *sql.DB) Place {
	return Place{
		query:     dbx.NewPlacesQ(db),
		locale:    dbx.NewPlaceLocalesQ(db),
		timetable: dbx.NewPlaceTimetablesQ(db),
		geo:       geo.NewGuesser(),
	}
}

func placeWithDetailsModelFromDB(in dbx.PlaceWithDetails) models.PlaceWithDetails {
	p := placeModelFromDB(in.Place)
	t := placeTimeTableModelFromDB(in.Timetable)

	out := models.PlaceWithDetails{
		ID:            p.ID,
		CityID:        p.CityID,
		DistributorID: p.DistributorID,
		Class:         p.Class,

		Status:    p.Status,
		Verified:  p.Verified,
		Ownership: p.Ownership,
		Point:     p.Point,
		Address:   p.Address,

		Locale:      in.Locale,
		Name:        in.Name,
		Description: in.Description,

		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,

		Timetable: t,
	}
	if p.Website != nil {
		out.Website = p.Website
	}
	if p.Phone != nil {
		out.Phone = p.Phone
	}

	return out
}

func placeModelFromDB(dbPlace dbx.Place) models.Place {
	place := models.Place{
		ID:        dbPlace.ID,
		CityID:    dbPlace.CityID,
		Class:     dbPlace.Class,
		Status:    dbPlace.Status,
		Verified:  dbPlace.Verified,
		Point:     dbPlace.Point,
		Address:   dbPlace.Address,
		CreatedAt: dbPlace.CreatedAt,
		UpdatedAt: dbPlace.UpdatedAt,
	}
	if dbPlace.DistributorID.Valid {
		place.DistributorID = &dbPlace.DistributorID.UUID
	}
	if dbPlace.Website.Valid {
		place.Website = &dbPlace.Website.String
	}
	if dbPlace.Phone.Valid {
		place.Phone = &dbPlace.Phone.String
	}

	return place
}

func placeLocaleModelFromDB(dbLoc dbx.PlaceLocale) models.PlaceLocale {
	return models.PlaceLocale{
		PlaceID:     dbLoc.PlaceID,
		Locale:      dbLoc.Locale,
		Name:        dbLoc.Name,
		Description: dbLoc.Description,
	}

}

func placeTimeTableModelFromDB(dbTI []dbx.PlaceTimetable) models.Timetable {
	res := models.Timetable{
		Table: make([]models.TimeInterval, 0, len(dbTI)),
	}
	for _, ti := range dbTI {
		res.Table = append(res.Table, models.TimeInterval{
			From: models.NumberMinutesToMoment(ti.StartMin),
			To:   models.NumberMinutesToMoment(ti.EndMin),
		})
	}

	return res
}
