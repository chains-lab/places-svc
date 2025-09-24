package apptest

import (
	"database/sql"
	"testing"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/cmd/config"
	"github.com/chains-lab/places-svc/internal/domain"
	"github.com/chains-lab/places-svc/internal/data/migrations"
)

// TEST DATABASE CONNECTION
const testDatabaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}

type Setup struct {
	app app.App
	Log logium.Logger
	Cfg config.Config
}

func cleanDb(t *testing.T) {
	err := migrations.MigrateDown(testDatabaseURL)
	if err != nil {
		t.Fatalf("migrate down: %v", err)
	}
	err = migrations.MigrateUp(testDatabaseURL)
	if err != nil {
		t.Fatalf("migrate up: %v", err)
	}
}

func newSetup(t *testing.T) (Setup, error) {
	cfg := config.Config{
		Database: config.DatabaseConfig{
			SQL: struct {
				URL string `mapstructure:"url"`
			}{
				URL: testDatabaseURL,
			},
		},
	}

	log := logium.NewLogger("debug", "text")

	a, err := app.NewApp(cfg)
	if err != nil {
		t.Fatal(err)
	}

	return Setup{
		app: a,
		Log: logium.NewWithBase(log),
		Cfg: cfg,
	}, nil
}
