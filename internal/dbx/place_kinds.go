package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceKindsTable = "place_kinds"

type PlaceKind struct {
	Code         string        `db:"code"`
	CategoryCode string        `db:"category_code"`
	Status       string        `db:"status"`
	Icon         string        `db:"icon"`
	Locale       LocaleForKind `db:"locale"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}

type LocaleForKind struct {
	Locale string `db:"locale"`
	Name   string `db:"name"`
}

type KindsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceKindsQ(db *sql.DB) KindsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return KindsQ{
		db: db,
		selector: b.
			Select(
				PlaceKindsTable+".code",
				PlaceKindsTable+".category_code",
				PlaceKindsTable+".status",
				PlaceKindsTable+".icon",
				PlaceKindsTable+".created_at",
				PlaceKindsTable+".updated_at",
				"NULL AS loc_name",
				"NULL AS loc_locale",
			).
			From(PlaceKindsTable),
		inserter: b.Insert(PlaceKindsTable),
		updater:  b.Update(PlaceKindsTable),
		deleter:  b.Delete(PlaceKindsTable),
		counter:  b.Select("COUNT(*) AS count").From(PlaceKindsTable),
	}
}

func (q KindsQ) New() KindsQ { return NewPlaceKindsQ(q.db) }

func (q KindsQ) Insert(ctx context.Context, in PlaceKind) error {
	values := map[string]interface{}{
		"code":          in.Code,
		"category_code": in.CategoryCode,
		"status":        in.Status,
		"icon":          in.Icon,
		"updated_at":    in.UpdatedAt,
		"created_at":    in.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceKindsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func scanPlaceKind(scanner interface{ Scan(dest ...any) error }) (PlaceKind, error) {
	var pc PlaceKind
	var locName, locLocale sql.NullString

	err := scanner.Scan(
		&pc.Code,
		&pc.CategoryCode,
		&pc.Status,
		&pc.Icon,
		&pc.CreatedAt,
		&pc.UpdatedAt,
		&locName,
		&locLocale,
	)
	if err != nil {
		return PlaceKind{}, err
	}

	if locName.Valid && locLocale.Valid {
		pc.Locale = LocaleForKind{
			Locale: locLocale.String,
			Name:   locName.String,
		}
	} else {
		pc.Locale = LocaleForKind{}
	}

	return pc, nil
}

func (q KindsQ) Get(ctx context.Context) (PlaceKind, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceKind{}, fmt.Errorf("building select query for %s: %w", PlaceKindsTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	return scanPlaceKind(row)
}

func (q KindsQ) Select(ctx context.Context) ([]PlaceKind, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", PlaceKindsTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, fmt.Errorf("querying select for %s: %w", PlaceKindsTable, err)
	}
	defer rows.Close()

	var out []PlaceKind
	for rows.Next() {
		pk, err := scanPlaceKind(rows)
		if err != nil {
			return out, err
		}
		out = append(out, pk)
	}

	if err = rows.Err(); err != nil {
		return out, fmt.Errorf("iterating rows for %s: %w", PlaceKindsTable, err)
	}

	return out, nil
}

type PlaceUpdateParams struct {
	CategoryCode *string   `db:"category_code"`
	Status       *string   `db:"status"`
	Icon         *string   `db:"icon"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (q KindsQ) Update(ctx context.Context, params PlaceUpdateParams) error {
	values := map[string]interface{}{
		"updated_at": params.UpdatedAt,
	}
	if params.CategoryCode != nil {
		values["category_code"] = *params.CategoryCode
	}
	if params.Status != nil {
		values["status"] = *params.Status
	}
	if params.Icon != nil {
		values["icon"] = *params.Icon
	}

	if len(values) == 0 {
		return nil // nothing to update
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", PlaceKindsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q KindsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceKindsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q KindsQ) FilterCode(code string) KindsQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})

	return q
}

func (q KindsQ) FilterCategoryCode(categoryCode string) KindsQ {
	q.selector = q.selector.Where(sq.Eq{"category_code": categoryCode})
	q.updater = q.updater.Where(sq.Eq{"category_code": categoryCode})
	q.deleter = q.deleter.Where(sq.Eq{"category_code": categoryCode})
	q.counter = q.counter.Where(sq.Eq{"category_code": categoryCode})

	return q
}

func (q KindsQ) FilterStatus(status string) KindsQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})

	return q
}

func (q KindsQ) WithLocale(locale string) KindsQ {
	base := PlaceKindsTable       // "place_kinds"
	i18n := PlaceKindLocalesTable // "place_kind_i18n"

	l := sanitizeLocale(locale)

	col := func(field, alias string) sq.Sqlizer {
		return sq.Expr(
			`CASE
				WHEN EXISTS (
					SELECT 1 FROM `+i18n+` i
					WHERE i.kind_code = k.code AND i.locale = ?
				)
				THEN (SELECT i.`+field+`  FROM `+i18n+` i  WHERE i.kind_code = k.code AND i.locale = ?)
				ELSE (SELECT i2.`+field+` FROM `+i18n+` i2 WHERE i2.kind_code = k.code AND i2.locale = 'en')
			END AS `+alias,
			l, l,
		)
	}

	q.selector = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"k.code",
			"k.category_code",
			"k.status",
			"k.icon",
			"k.created_at",
			"k.updated_at",
		).
		// порядок для скана: loc_locale, loc_name
		Column(col("name", "loc_name")).
		Column(col("locale", "loc_locale")).
		From(base + " AS k")

	return q
}

func (q KindsQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceKindsTable, err)
	}

	var count int
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlaceKindsTable, err)
	}

	return count, nil
}

func (q KindsQ) Page(limit, offset uint64) KindsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
