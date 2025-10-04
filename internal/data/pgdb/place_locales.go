package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/enum"
	"github.com/google/uuid"
)

const placeLocalizationTable = "place_i18n"

var reLocale = regexp.MustCompile(`^[a-z]{2}$`)

func SanitizeLocale(l string) string {
	if reLocale.MatchString(l) {
		return l
	}

	return enum.LocaleEN
}

type PlaceLocale struct {
	PlaceID     uuid.UUID `storage:"place_id"`
	Locale      string    `storage:"locale"`
	Name        string    `storage:"name"`
	Description string    `storage:"description"`
}

type PlaceLocalesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceLocalesQ(db *sql.DB) PlaceLocalesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlaceLocalesQ{
		db: db,
		selector: b.Select(
			"place_id",
			"locale",
			"name",
			"description",
		).From(placeLocalizationTable),
		inserter: b.Insert(placeLocalizationTable),
		updater:  b.Update(placeLocalizationTable),
		deleter:  b.Delete(placeLocalizationTable),
		counter:  b.Select("COUNT(*) AS count").From(placeLocalizationTable),
	}
}

func (q PlaceLocalesQ) New() PlaceLocalesQ { return NewPlaceLocalesQ(q.db) }

func (q PlaceLocalesQ) Insert(ctx context.Context, in PlaceLocale) error {
	values := map[string]any{
		"place_id":    in.PlaceID,
		"locale":      SanitizeLocale(in.Locale),
		"name":        in.Name,
		"description": in.Description,
	}
	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", placeLocalizationTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q PlaceLocalesQ) Upsert(ctx context.Context, in ...PlaceLocale) error {
	if len(in) == 0 {
		return nil
	}

	const cols = "(place_id, locale, name, description)"
	var (
		args []any
		ph   []string
		i    = 1
	)
	for _, row := range in {
		ph = append(ph, fmt.Sprintf("($%d,$%d,$%d,$%d)", i, i+1, i+2, i+3))
		i += 4
		args = append(args, row.PlaceID, SanitizeLocale(row.Locale), row.Name, row.Description)
	}
	query := fmt.Sprintf(`
		INSERT INTO %s %s VALUES %s
		ON CONFLICT (place_id, locale) DO UPDATE
		SET name = EXCLUDED.name,
		    description = EXCLUDED.description
	`, placeLocalizationTable, cols, strings.Join(ph, ","))

	if tx, ok := TxFromCtx(ctx); ok {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}
	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}

func (q PlaceLocalesQ) Get(ctx context.Context) (PlaceLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceLocale{}, fmt.Errorf("build select %s: %w", placeLocalizationTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var out PlaceLocale
	if err := row.Scan(&out.PlaceID, &out.Locale, &out.Name, &out.Description); err != nil {
		return out, err
	}
	return out, nil
}

func (q PlaceLocalesQ) Select(ctx context.Context) ([]PlaceLocale, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", placeLocalizationTable, err)
	}

	var rows *sql.Rows
	if tx, ok := TxFromCtx(ctx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PlaceLocale
	for rows.Next() {
		var pl PlaceLocale
		if err := rows.Scan(&pl.PlaceID, &pl.Locale, &pl.Name, &pl.Description); err != nil {
			return nil, fmt.Errorf("scan %s: %w", placeLocalizationTable, err)
		}
		out = append(out, pl)
	}
	return out, rows.Err()
}

func (q PlaceLocalesQ) Update(ctx context.Context) error {
	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", placeLocalizationTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q PlaceLocalesQ) UpdateName(name string) PlaceLocalesQ {
	q.updater = q.updater.Set("name", name)
	return q
}

func (q PlaceLocalesQ) UpdateDescription(description string) PlaceLocalesQ {
	q.updater = q.updater.Set("description", description)
	return q
}

func (q PlaceLocalesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", placeLocalizationTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q PlaceLocalesQ) FilterPlaceID(id uuid.UUID) PlaceLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": id})
	q.updater = q.updater.Where(sq.Eq{"place_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": id})
	q.counter = q.counter.Where(sq.Eq{"place_id": id})
	return q
}

func (q PlaceLocalesQ) FilterByLocale(locale string) PlaceLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"locale": locale})
	q.updater = q.updater.Where(sq.Eq{"locale": locale})
	q.deleter = q.deleter.Where(sq.Eq{"locale": locale})
	q.counter = q.counter.Where(sq.Eq{"locale": locale})
	return q
}

func (q PlaceLocalesQ) FilterByName(name string) PlaceLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"name": name})
	q.updater = q.updater.Where(sq.Eq{"name": name})
	q.deleter = q.deleter.Where(sq.Eq{"name": name})
	q.counter = q.counter.Where(sq.Eq{"name": name})
	return q
}

func (q PlaceLocalesQ) OrderByLocale(asc bool) PlaceLocalesQ {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}
	q.selector = q.selector.OrderBy("locale " + dir)
	return q
}

func (q PlaceLocalesQ) Page(limit, offset uint) PlaceLocalesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q PlaceLocalesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placeLocalizationTable, err)
	}

	var cnt uint
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
