package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/places-svc/internal/app/entities/class"
	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/config"
	"github.com/chains-lab/places-svc/internal/dbx"
)

type App struct {
	place         place.Place
	classificator class.Classificator

	db *sql.DB
}

func NewApp(cfg config.Config) (App, error) {
	pg, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return App{}, err
	}

	return App{
		place:         place.NewPlace(pg),
		classificator: class.NewClassificator(pg),

		db: pg,
	}, err
}

func (a App) transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, dbx.TxKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
