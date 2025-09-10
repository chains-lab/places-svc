package dbx

import (
	"database/sql"
	"embed"
	"regexp"

	"github.com/chains-lab/places-svc/internal/config"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

type TxKeyType struct{}

var TxKey = TxKeyType{}

var reLocale = regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`)

func sanitizeLocale(l string) string {
	if reLocale.MatchString(l) {
		return l
	}
	return "en"
}

func (q ClassesQ) SelectorToSql() (string, []any, error) {
	return q.selector.ToSql()
}

func (q PlacesQ) SelectorToSql() (string, []any, error) {
	return q.selector.ToSql()
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

//go:embed migrations/*.sql
var Migrations embed.FS

var migrations = &migrate.EmbedFileSystemMigrationSource{
	FileSystem: Migrations,
	Root:       "migrations",
}

func MigrateUp(cfg config.Config) error {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to applyConditions migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")

	return nil
}

func MigrateDown(cfg config.Config) error {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to applyConditions migrations")
	}
	logrus.WithField("applied", applied).Info("migrations applied")

	return nil
}
