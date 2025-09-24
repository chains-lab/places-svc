package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/google/uuid"
)

const placeLocalizationTable = "place_i18n"

type placeLocalesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceLocalesQ(db *sql.DB) schemas.PlaceLocalesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &placeLocalesQ{
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

func (q *placeLocalesQ) New() schemas.PlaceLocalesQ { return NewPlaceLocalesQ(q.db) }

func (q *placeLocalesQ) Insert(ctx context.Context, in schemas.PlaceLocale) error {
	values := map[string]any{
		"place_id":    in.PlaceID,
		"locale":      sanitizeLocale(in.Locale),
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

func (q *placeLocalesQ) Upsert(ctx context.Context, in ...schemas.PlaceLocale) error {
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
		args = append(args, row.PlaceID, sanitizeLocale(row.Locale), row.Name, row.Description)
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

func (q *placeLocalesQ) Get(ctx context.Context) (schemas.PlaceLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.PlaceLocale{}, fmt.Errorf("build select %s: %w", placeLocalizationTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var out schemas.PlaceLocale
	if err := row.Scan(&out.PlaceID, &out.Locale, &out.Name, &out.Description); err != nil {
		return out, err
	}
	return out, nil
}

func (q *placeLocalesQ) Select(ctx context.Context) ([]schemas.PlaceLocale, error) {
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

	var out []schemas.PlaceLocale
	for rows.Next() {
		var pl schemas.PlaceLocale
		if err := rows.Scan(&pl.PlaceID, &pl.Locale, &pl.Name, &pl.Description); err != nil {
			return nil, fmt.Errorf("scan %s: %w", placeLocalizationTable, err)
		}
		out = append(out, pl)
	}
	return out, rows.Err()
}

func (q *placeLocalesQ) Update(ctx context.Context, params schemas.UpdatePlaceLocaleParams) error {
	updates := map[string]any{}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.Description != nil {
		updates["description"] = *params.Description
	}
	if len(updates) == 0 {
		return nil // no-op
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", placeLocalizationTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *placeLocalesQ) Delete(ctx context.Context) error {
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

func (q *placeLocalesQ) FilterPlaceID(id uuid.UUID) schemas.PlaceLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": id})
	q.updater = q.updater.Where(sq.Eq{"place_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": id})
	q.counter = q.counter.Where(sq.Eq{"place_id": id})
	return q
}

func (q *placeLocalesQ) FilterByLocale(locale string) schemas.PlaceLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"locale": locale})
	q.updater = q.updater.Where(sq.Eq{"locale": locale})
	q.deleter = q.deleter.Where(sq.Eq{"locale": locale})
	q.counter = q.counter.Where(sq.Eq{"locale": locale})
	return q
}

func (q *placeLocalesQ) FilterByName(name string) schemas.PlaceLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"name": name})
	q.updater = q.updater.Where(sq.Eq{"name": name})
	q.deleter = q.deleter.Where(sq.Eq{"name": name})
	q.counter = q.counter.Where(sq.Eq{"name": name})
	return q
}

func (q *placeLocalesQ) OrderByLocale(asc bool) schemas.PlaceLocalesQ {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}
	q.selector = q.selector.OrderBy("locale " + dir)
	return q
}

func (q *placeLocalesQ) Page(limit, offset uint64) schemas.PlaceLocalesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q *placeLocalesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placeLocalizationTable, err)
	}

	var cnt uint64
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
