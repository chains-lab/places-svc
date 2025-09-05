package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceCategoryLocalesTable = "place_category_i18n"

type PlaceCategoryLocale struct {
	CategoryCode string `db:"category_code"`
	Locale       string `db:"locale"`
	Name         string `db:"name"`
}

type CategoryLocaleQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCategoryLocaleQ(db *sql.DB) CategoryLocaleQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return CategoryLocaleQ{
		db:       db,
		selector: b.Select("*").From(PlaceCategoryLocalesTable),
		inserter: b.Insert(PlaceCategoryLocalesTable),
		updater:  b.Update(PlaceCategoryLocalesTable),
		deleter:  b.Delete(PlaceCategoryLocalesTable),
	}
}

func (q CategoryLocaleQ) New() CategoryLocaleQ { return NewCategoryLocaleQ(q.db) }

func (q CategoryLocaleQ) Insert(ctx context.Context, in PlaceCategoryLocale) error {
	values := map[string]interface{}{
		"category_code": in.CategoryCode,
		"locale":        in.Locale,
		"name":          in.Name,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceCategoryLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CategoryLocaleQ) Upsert(ctx context.Context, in PlaceCategoryLocale) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (category_code, locale, name)
		VALUES ($1, $2, $3)
		ON CONFLICT (category_code, locale) DO UPDATE
		SET name = EXCLUDED.name
    `, PlaceCategoryLocalesTable)

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err := tx.ExecContext(ctx, query, in.CategoryCode, in.Locale, in.Name)
		return err
	}
	_, err := q.db.ExecContext(ctx, query, in.CategoryCode, in.Locale, in.Name)
	return err
}

func (q CategoryLocaleQ) Get(ctx context.Context) (PlaceCategoryLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceCategoryLocale{}, fmt.Errorf("building select query for %s: %w", PlaceCategoryLocalesTable, err)
	}

	var out PlaceCategoryLocale
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.CategoryCode,
		&out.Locale,
		&out.Name,
	)

	return out, err
}

func (q CategoryLocaleQ) Select(ctx context.Context) ([]PlaceCategoryLocale, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceCategoryLocalesTable, err)
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

	var out []PlaceCategoryLocale
	for rows.Next() {
		var item PlaceCategoryLocale
		err = rows.Scan(
			&item.CategoryCode,
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

func (q CategoryLocaleQ) FilterCategoryCode(code string) CategoryLocaleQ {
	q.selector = q.selector.Where(sq.Eq{"category_code": code})
	q.counter = q.counter.Where(sq.Eq{"category_code": code})
	q.updater = q.updater.Where(sq.Eq{"category_code": code})
	q.deleter = q.deleter.Where(sq.Eq{"category_code": code})

	return q
}

func (q CategoryLocaleQ) FilterLocale(locale string) CategoryLocaleQ {
	q.selector = q.selector.Where(sq.Eq{"locale": locale})
	q.counter = q.counter.Where(sq.Eq{"locale": locale})
	q.updater = q.updater.Where(sq.Eq{"locale": locale})
	q.deleter = q.deleter.Where(sq.Eq{"locale": locale})

	return q
}

func (q CategoryLocaleQ) FilterNameLike(name string) CategoryLocaleQ {
	q.selector = q.selector.Where(sq.Like{"name": name})
	q.counter = q.counter.Where(sq.Like{"name": name})
	q.updater = q.updater.Where(sq.Like{"name": name})
	q.deleter = q.deleter.Where(sq.Like{"name": name})

	return q
}

type UpdateCategoryLocaleParams struct {
	Name      *string
	UpdatedAt time.Time
}

func (q CategoryLocaleQ) Update(ctx context.Context, in UpdateCategoryLocaleParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Name != nil {
		values["name"] = *in.Name
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", PlaceCategoryLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CategoryLocaleQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceCategoryLocalesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}
