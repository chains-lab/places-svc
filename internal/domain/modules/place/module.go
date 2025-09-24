package place

import (
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/data/fabric"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/modules/place/geo"
)

type Module struct {
	db  fabric.Database
	geo *geo.Guesser
}

func NewModule(db fabric.Database) Module {
	return Module{
		db:  db,
		geo: geo.NewGuesser(),
	}
}

func placeWithDetailsModelFromDB(in data.PlaceWithDetails) models.PlaceWithDetails {
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

func placeModelFromDB(dbPlace data.Place) models.Place {
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

func placeLocaleModelFromDB(dbLoc data.PlaceLocale) models.PlaceLocale {
	return models.PlaceLocale{
		PlaceID:     dbLoc.PlaceID,
		Locale:      dbLoc.Locale,
		Name:        dbLoc.Name,
		Description: dbLoc.Description,
	}

}

func placeTimeTableModelFromDB(dbTI []data.PlaceTimetable) models.Timetable {
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
