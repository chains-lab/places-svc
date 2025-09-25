package domaintest

import (
	"database/sql"
	"testing"

	"github.com/chains-lab/places-svc/cmd/migrations"
	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
)

// TEST DATABASE CONNECTION
const testDatabaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}

type services struct {
	class class.Service
	place place.Service
}

type Setup struct {
	domain services
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
	database := data.NewDatabase(testDatabaseURL)
	classMod := class.NewService(database)
	placeMod := place.NewService(database)

	return Setup{
		domain: services{
			class: classMod,
			place: placeMod,
		},
	}, nil
}
