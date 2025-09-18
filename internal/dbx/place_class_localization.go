package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceClassLocalesTable = "place_class_i18n"

type PlaceClassLocale struct {
	Class  string `db:"class"`
	Locale string `db:"locale"`
	Name   string `db:"name"`
}

type ClassLocaleQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewClassLocaleQ(db *sql.DB) ClassLocaleQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ClassLocaleQ{
		db:       db,
		selector: b.Select("*").From(PlaceClassLocalesTable),
		inserter: b.Insert(PlaceClassLocalesTable),
		updater:  b.Update(PlaceClassLocalesTable),
		deleter:  b.Delete(PlaceClassLocalesTable),
		counter:  b.Select("COUNT(*) AS count").From(PlaceClassLocalesTable),
	}
}

func (q ClassLocaleQ) New() ClassLocaleQ { return NewClassLocaleQ(q.db) }

func (q ClassLocaleQ) Insert(ctx context.Context, in ...PlaceClassLocale) error {
	if len(in) == 0 {
		return nil
	}

	ins := q.inserter.Columns(
		"class",
		"locale",
		"name",
	)
	for _, item := range in {
		ins = ins.Values(item.Class, item.Locale, item.Name)
	}

	query, args, err := ins.ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceClassLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassLocaleQ) Upsert(ctx context.Context, in ...PlaceClassLocale) error {
	if len(in) == 0 {
		return nil
	}

	args := make([]any, 0, len(in)*3)
	ph := make([]string, 0, len(in))
	for i, row := range in {
		base := i*3 + 1
		ph = append(ph, fmt.Sprintf("($%d,$%d,$%d)", base, base+1, base+2))
		args = append(args, row.Class, row.Locale, row.Name)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (class, locale, name)
		VALUES %s
		ON CONFLICT (class, locale) DO UPDATE
		SET name = EXCLUDED.name
	`, PlaceClassLocalesTable, strings.Join(ph, ","))

	var err error
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassLocaleQ) Get(ctx context.Context) (PlaceClassLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceClassLocale{}, fmt.Errorf("building select query for %s: %w", PlaceClassLocalesTable, err)
	}

	var out PlaceClassLocale
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.Class,
		&out.Locale,
		&out.Name,
	)

	return out, err
}

func (q ClassLocaleQ) Select(ctx context.Context) ([]PlaceClassLocale, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceClassLocalesTable, err)
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

	var out []PlaceClassLocale
	for rows.Next() {
		var item PlaceClassLocale
		err = rows.Scan(
			&item.Class,
			&item.Locale,
			&item.Name,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}

	return out, err
}

type UpdateClassLocaleParams struct {
	Name      *string
	UpdatedAt time.Time
}

func (q ClassLocaleQ) Update(ctx context.Context, in UpdateClassLocaleParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Name != nil {
		values["name"] = *in.Name
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", PlaceClassLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q ClassLocaleQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceClassLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q ClassLocaleQ) FilterClass(code string) ClassLocaleQ {
	q.selector = q.selector.Where(sq.Eq{"class": code})
	q.counter = q.counter.Where(sq.Eq{"class": code})
	q.updater = q.updater.Where(sq.Eq{"class": code})
	q.deleter = q.deleter.Where(sq.Eq{"class": code})

	return q
}

func (q ClassLocaleQ) FilterLocale(locale string) ClassLocaleQ {
	q.selector = q.selector.Where(sq.Eq{"locale": locale})
	q.counter = q.counter.Where(sq.Eq{"locale": locale})
	q.updater = q.updater.Where(sq.Eq{"locale": locale})
	q.deleter = q.deleter.Where(sq.Eq{"locale": locale})

	return q
}

func (q ClassLocaleQ) FilterNameLike(name string) ClassLocaleQ {
	q.selector = q.selector.Where(sq.Like{"name": name})
	q.counter = q.counter.Where(sq.Like{"name": name})
	q.updater = q.updater.Where(sq.Like{"name": name})
	q.deleter = q.deleter.Where(sq.Like{"name": name})

	return q
}

func (q ClassLocaleQ) OrderByLocale(asc bool) ClassLocaleQ {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}

	q.selector = q.selector.OrderBy("locale " + dir)

	return q
}

func (q ClassLocaleQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceClassLocalesTable, err)
	}

	var count uint64
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	if err = row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count from %s: %w", PlaceClassLocalesTable, err)
	}

	return count, nil
}

func (q ClassLocaleQ) Page(limit, offset uint64) ClassLocaleQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
