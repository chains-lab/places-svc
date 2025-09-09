package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const placeLocalizationTable = "place_i18n"

type PlaceLocale struct {
	PlaceID     uuid.UUID `db:"place_id"`
	Locale      string    `db:"locale"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
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
		db:       db,
		selector: b.Select("*").From(placeLocalizationTable),
		inserter: b.Insert(placeLocalizationTable),
		updater:  b.Update(placeLocalizationTable),
		deleter:  b.Delete(placeLocalizationTable),
		counter:  b.Select("COUNT(*) AS count").From(placeLocalizationTable),
	}
}

func (q PlaceLocalesQ) New() PlaceLocalesQ { return NewPlaceLocalesQ(q.db) }

func (q PlaceLocalesQ) Insert(ctx context.Context, in PlaceLocale) error {
	values := map[string]interface{}{
		"place_id":    in.PlaceID,
		"locale":      in.Locale,
		"name":        in.Name,
		"description": in.Description,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query for %s: %w", placeLocalizationTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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

	const cols = "(place_id, locale, name, address, description)"
	var (
		args []any
		ph   []string
		i    = 1
	)

	for _, row := range in {
		// на каждый ряд по 5 плейсхолдеров: ($1,$2,$3,$4,$5), ($6,$7,$8,$9,$10), ...
		ph = append(ph, fmt.Sprintf("($%d,$%d,$%d,$%d)", i, i+1, i+2, i+3))
		i += 5
		args = append(args, row.PlaceID, row.Locale, row.Name, row.Description)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s %s VALUES %s
		ON CONFLICT (place_id, locale) DO UPDATE
		SET name = EXCLUDED.name,
		    address = EXCLUDED.address,
		    description = EXCLUDED.description
	`, placeLocalizationTable, cols, strings.Join(ph, ","))

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}
	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}

type UpdatePlaceLocaleParams struct {
	Name        *string
	Address     *string
	Description *string
}

func (q PlaceLocalesQ) Update(ctx context.Context, params UpdatePlaceLocaleParams) error {
	updates := map[string]interface{}{}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.Address != nil {
		updates["address"] = *params.Address
	}
	if params.Description != nil {
		updates["description"] = *params.Description
	}

	if len(updates) == 1 { // только updated_at
		return nil
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query for %s: %w", placeLocalizationTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlaceLocalesQ) Get(ctx context.Context) (PlaceLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceLocale{}, fmt.Errorf("failed to build select query for %s: %w", placeLocalizationTable, err)
	}

	var out PlaceLocale
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.PlaceID,
		&out.Locale,
		&out.Name,
		&out.Description,
	)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (q PlaceLocalesQ) Select(ctx context.Context) ([]PlaceLocale, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query for %s: %w", placeLocalizationTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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
		var pd PlaceLocale
		err := rows.Scan(
			&pd.PlaceID,
			&pd.Locale,
			&pd.Name,
			&pd.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row for %s: %w", placeLocalizationTable, err)
		}
		out = append(out, pd)
	}

	return out, rows.Err()
}

func (q PlaceLocalesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query for %s: %w", placeLocalizationTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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

func (q PlaceLocalesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query for %s: %w", placeLocalizationTable, err)
	}

	var count uint64
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to scan count for %s: %w", placeLocalizationTable, err)
	}

	return count, nil
}

func (q PlaceLocalesQ) Page(offset, limit uint64) PlaceLocalesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
