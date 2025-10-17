package data

import (
	"context"
	"database/sql"

	"github.com/chains-lab/places-svc/internal/data/pgdb"
	"github.com/chains-lab/places-svc/internal/domain/models"
	_ "github.com/lib/pq" // postgres driver don`t delete
)

type Database struct {
	sql SqlDB
}

func New(pg *sql.DB) Database {
	return Database{
		sql: SqlDB{
			classes:    pgdb.NewClassesQ(pg),
			places:     pgdb.NewPlacesQ(pg),
			pLocales:   pgdb.NewPlaceLocalesQ(pg),
			timetables: pgdb.NewPlaceTimetablesQ(pg),
		},
	}
}

func (d Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.sql.classes.New().Transaction(ctx, fn)
}

type SqlDB struct {
	classes    pgdb.ClassesQ
	places     pgdb.PlacesQ
	pLocales   pgdb.PlaceLocalesQ
	timetables pgdb.PlaceTimetablesQ
}

func modelFromDB(in pgdb.Place) models.Place {
	p := detailsFromDB(in.PlaceRow)
	t := timetableFromDB(in.Timetable)

	out := models.Place{
		ID:        p.ID,
		CityID:    p.CityID,
		CompanyID: p.CompanyID,
		Class:     p.Class,

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

func detailsFromDB(dbPlace pgdb.PlaceRow) models.PlaceDetails {
	place := models.PlaceDetails{
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
	if dbPlace.CompanyID.Valid {
		place.CompanyID = &dbPlace.CompanyID.UUID
	}
	if dbPlace.Website.Valid {
		place.Website = &dbPlace.Website.String
	}
	if dbPlace.Phone.Valid {
		place.Phone = &dbPlace.Phone.String
	}

	return place
}

func localeFromDB(dbLoc pgdb.PlaceLocale) models.PlaceLocale {
	return models.PlaceLocale{
		PlaceID:     dbLoc.PlaceID,
		Locale:      dbLoc.Locale,
		Name:        dbLoc.Name,
		Description: dbLoc.Description,
	}

}

func timetableFromDB(dbTI []pgdb.PlaceTimetableRow) models.Timetable {
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
