package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/places-svc/internal/data/pgdb"
	"github.com/chains-lab/places-svc/internal/domain/models"
	_ "github.com/lib/pq" // postgres driver don`t delete
)

type Database struct {
	db *sql.DB
}

func NewDatabase(url string) Database {
	db, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	return Database{db}
}

func (d *Database) Places() pgdb.PlacesQ {
	return pgdb.NewPlacesQ(d.db)
}

func (d *Database) PlaceLocales() pgdb.PlaceLocalesQ {
	return pgdb.NewPlaceLocalesQ(d.db)
}

func (d *Database) PlaceTimetables() pgdb.PlaceTimetablesQ {
	return pgdb.NewPlaceTimetablesQ(d.db)
}

func (d *Database) Classes() pgdb.ClassesQ {
	return pgdb.NewClassesQ(d.db)
}

func (d *Database) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := pgdb.TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, pgdb.TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func modelFromDB(in pgdb.PlaceWithDetails) models.Place {
	p := detailsFromDB(in.Place)
	t := timetableFromDB(in.Timetable)

	out := models.Place{
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

func detailsFromDB(dbPlace pgdb.Place) models.PlaceDetails {
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

func localeFromDB(dbLoc pgdb.PlaceLocale) models.PlaceLocale {
	return models.PlaceLocale{
		PlaceID:     dbLoc.PlaceID,
		Locale:      dbLoc.Locale,
		Name:        dbLoc.Name,
		Description: dbLoc.Description,
	}

}

func timetableFromDB(dbTI []pgdb.PlaceTimetable) models.Timetable {
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
