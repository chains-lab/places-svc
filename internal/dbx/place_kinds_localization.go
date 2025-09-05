package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceKindLocalesTable = "place_kind_i18n"

type PlaceKindLocale struct {
	KindCode string `db:"kind_code"`
	Locale   string `db:"locale"`
	Name     string `db:"name"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type KindLocaleQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewKindLocaleQ(db *sql.DB) KindLocaleQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return KindLocaleQ{
		db:       db,
		selector: b.Select("*").From(PlaceKindLocalesTable),
		inserter: b.Insert(PlaceKindLocalesTable),
		updater:  b.Update(PlaceKindLocalesTable),
		deleter:  b.Delete(PlaceKindLocalesTable),
	}
}

func (q KindLocaleQ) New() KindLocaleQ { return NewKindLocaleQ(q.db) }

func (q KindLocaleQ) Insert(ctx context.Context, in PlaceKindLocale) error {
	values := map[string]interface{}{
		"kind_code":  in.KindCode,
		"locale":     in.Locale,
		"name":       in.Name,
		"created_at": in.CreatedAt,
		"updated_at": in.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceKindLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q KindLocaleQ) Upsert(ctx context.Context, in PlaceKindLocale) error {
	query := fmt.Sprintf(`
        INSERT INTO %s (kind_code, locale, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (kind_code, locale) DO UPDATE
        SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
    `, PlaceKindLocalesTable)

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err := tx.ExecContext(ctx, query, in.KindCode, in.Locale, in.Name, in.CreatedAt, in.UpdatedAt)
		return err
	}
	_, err := q.db.ExecContext(ctx, query, in.KindCode, in.Locale, in.Name, in.CreatedAt, in.UpdatedAt)
	return err
}

func (q KindLocaleQ) Get(ctx context.Context) (PlaceKindLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceKindLocale{}, fmt.Errorf("building select query for %s: %w", PlaceKindLocalesTable, err)
	}

	var out PlaceKindLocale
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.KindCode,
		&out.Locale,
		&out.Name,
		&out.CreatedAt,
		&out.UpdatedAt,
	)

	return out, err
}

func (q KindLocaleQ) Select(ctx context.Context) ([]PlaceKindLocale, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceKindLocalesTable, err)
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

	var out []PlaceKindLocale
	for rows.Next() {
		var item PlaceKindLocale
		err = rows.Scan(
			&item.KindCode,
			&item.Locale,
			&item.Name,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}

	return out, err
}

func (q KindLocaleQ) FilterKindCode(code string) KindLocaleQ {
	q.selector = q.selector.Where(sq.Eq{"kind_code": code})
	q.counter = q.counter.Where(sq.Eq{"kind_code": code})
	q.updater = q.updater.Where(sq.Eq{"kind_code": code})
	q.deleter = q.deleter.Where(sq.Eq{"kind_code": code})

	return q
}

func (q KindLocaleQ) FilterLocale(locale string) KindLocaleQ {
	q.selector = q.selector.Where(sq.Eq{"locale": locale})
	q.counter = q.counter.Where(sq.Eq{"locale": locale})
	q.updater = q.updater.Where(sq.Eq{"locale": locale})
	q.deleter = q.deleter.Where(sq.Eq{"locale": locale})

	return q
}

func (q KindLocaleQ) FilterNameLike(name string) KindLocaleQ {
	q.selector = q.selector.Where(sq.Like{"name": name})
	q.counter = q.counter.Where(sq.Like{"name": name})
	q.updater = q.updater.Where(sq.Like{"name": name})
	q.deleter = q.deleter.Where(sq.Like{"name": name})

	return q
}

type UpdateKindLocaleParams struct {
	Name      *string
	UpdatedAt time.Time
}

func (q KindLocaleQ) Update(ctx context.Context, in UpdateKindLocaleParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Name != nil {
		values["name"] = *in.Name
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", PlaceKindLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q KindLocaleQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceKindLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}
