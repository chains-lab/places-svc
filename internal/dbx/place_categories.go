package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const placeCategoriesTable = "place_categories"

type PlaceCategory struct {
	Code      string            `db:"code"`
	Status    string            `db:"status"`
	Icon      string            `db:"icon"`
	Locale    LocaleForCategory `db:"locales"`
	UpdatedAt time.Time         `db:"updated_at"`
	CreatedAt time.Time         `db:"created_at"`
}

type LocaleForCategory struct {
	Locale string `db:"locale"`
	Name   string `db:"name"`
}

type CategoryQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewCategoryQ(db *sql.DB) CategoryQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return CategoryQ{
		db:       db,
		selector: b.Select("*").From(placeCategoriesTable),
		inserter: b.Insert(placeCategoriesTable),
		updater:  b.Update(placeCategoriesTable),
		deleter:  b.Delete(placeCategoriesTable),
		counter:  b.Select("COUNT(*) AS count").From(placeCategoriesTable),
	}
}

func (q CategoryQ) New() CategoryQ { return NewCategoryQ(q.db) }

func (q CategoryQ) Insert(ctx context.Context, in PlaceCategory) error {
	values := map[string]interface{}{
		"code":       in.Code,
		"status":     in.Status,
		"icon":       in.Icon,
		"created_at": in.CreatedAt,
		"updated_at": in.UpdatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert query for %s: %w", placeCategoriesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func scanPlaceCategory(scanner interface{ Scan(dest ...any) error }) (PlaceCategory, error) {
	var pc PlaceCategory
	var locName, locLocale sql.NullString

	err := scanner.Scan(
		&pc.Code,
		&pc.Status,
		&pc.Icon,
		&pc.CreatedAt,
		&pc.UpdatedAt,
		&locName,
		&locLocale,
	)
	if err != nil {
		return PlaceCategory{}, err
	}

	if locName.Valid && locLocale.Valid {
		pc.Locale = LocaleForCategory{
			Locale: locLocale.String,
			Name:   locName.String,
		}
	} else {
		pc.Locale = LocaleForCategory{}
	}

	return pc, nil
}

func (q CategoryQ) Get(ctx context.Context) (PlaceCategory, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceCategory{}, fmt.Errorf("build select query for %s: %w", placeCategoriesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	return scanPlaceCategory(row)
}

func (q CategoryQ) Select(ctx context.Context) ([]PlaceCategory, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query for %s: %w", placeCategoriesTable, err)
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

	var out []PlaceCategory
	for rows.Next() {
		pc, err := scanPlaceCategory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pc)
	}

	return out, rows.Err()
}

type UpdatePlaceCategoryParams struct {
	Status    *string
	Icon      *string
	UpdatedAt time.Time
}

func (q CategoryQ) Update(ctx context.Context, in UpdatePlaceCategoryParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Status != nil {
		values["status"] = *in.Status
	}
	if in.Icon != nil {
		values["icon"] = *in.Icon
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build update query for %s: %w", placeCategoriesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CategoryQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete query for %s: %w", placeCategoriesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q CategoryQ) WithLocale(locale string) CategoryQ {
	base := placeCategoriesTable
	i18n := PlaceCategoryLocalesTable

	subq := fmt.Sprintf(`
		LATERAL (
			SELECT i.name, i.locale
			FROM %s i
			WHERE i.category_code = %s.code
			  AND i.locale IN ($1, 'en')
			ORDER BY CASE
				WHEN i.locale = $1 THEN 1
				WHEN i.locale = 'en' THEN 2
				ELSE 3
			END
			LIMIT 1
		) loc
	`, i18n, base)

	q.selector = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			base+".code",
			base+".status",
			base+".created_at",
			base+".updated_at",
			"loc.name AS loc_name",
			"loc.locale AS loc_locale",
		).
		From(base).
		LeftJoin(subq + " ON TRUE")

	q.selector = q.selector.PlaceholderFormat(sq.Dollar).RunWith(q.db).Suffix("", locale)

	return q
}

func (q CategoryQ) FilterCode(code string) CategoryQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})

	return q
}

func (q CategoryQ) FilterStatus(status string) CategoryQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})

	return q
}

func (q CategoryQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count query for %s: %w", placeCategoriesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (q CategoryQ) Page(limit, offset uint64) CategoryQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
