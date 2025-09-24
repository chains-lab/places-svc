package fabric

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/data/pgdb"
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

func (d *Database) Places() data.PlacesQ {
	return pgdb.NewPlacesQ(d.db)
}

func (d *Database) PlaceLocales() data.PlaceLocalesQ {
	return pgdb.NewPlaceLocalesQ(d.db)
}

func (d *Database) PlaceTimetables() data.PlaceTimetablesQ {
	return pgdb.NewPlaceTimetablesQ(d.db)
}

func (d *Database) Classes() data.ClassesQ {
	return pgdb.NewClassesQ(d.db)
}

func (d *Database) ClassLocales() data.ClassLocalesQ {
	return pgdb.NewClassLocalesQ(d.db)
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
