package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceClassesTable = "place_classes"

type PlaceClass struct {
	Code      string         `db:"code"`
	Father    sql.NullString `db:"father"` // NULL для корней
	Status    string         `db:"status"`
	Icon      string         `db:"icon"`
	Path      string         `db:"path"` // ltree как text
	Locale    *string        `db:"locale"`
	Name      *string        `db:"name"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type ClassesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewClassesQ(db *sql.DB) ClassesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ClassesQ{
		db: db,
		selector: b.Select(
			"pc.code",
			"pc.father",
			"pc.status",
			"pc.icon",
			"pc.path",
			"pc.created_at",
			"pc.updated_at",
			"NULL AS loc_name",
			"NULL AS loc_locale",
		).From(PlaceClassesTable + " AS pc"),
		inserter: b.Insert(PlaceClassesTable),
		updater:  b.Update(PlaceClassesTable + " AS pc"),
		deleter:  b.Delete(PlaceClassesTable + " AS pc"),
		counter:  b.Select("COUNT(*) AS count").From(PlaceClassesTable + " AS pc"),
	}
}

func (q ClassesQ) New() ClassesQ { return NewClassesQ(q.db) }

type InsertPlaceClass struct {
	Code      string         `db:"code"`
	Father    sql.NullString `db:"father"` // NULL для корней
	Status    string         `db:"status"`
	Icon      string         `db:"icon"`
	Path      string         `db:"path"` // ltree как text
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func (q ClassesQ) Insert(ctx context.Context, in InsertPlaceClass) error {
	values := map[string]interface{}{
		"code":   in.Code,
		"status": in.Status,
		"icon":   in.Icon,
	}
	if in.Father.Valid {
		values["father"] = in.Father.String
	} else {
		values["father"] = nil
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert query for %s: %w", PlaceClassesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func scanPlaceClass(scanner interface{ Scan(dest ...any) error }) (PlaceClass, error) {
	var pc PlaceClass
	var locName, locLocale *string

	if err := scanner.Scan(
		&pc.Code,
		&pc.Father,
		&pc.Status,
		&pc.Icon,
		&pc.Path,
		&pc.CreatedAt,
		&pc.UpdatedAt,
		&locName,
		&locLocale,
	); err != nil {
		return PlaceClass{}, err
	}

	// Locale (пустая структура, если не найдено ни запрошенной, ни 'en')
	pc.Name = locName
	pc.Locale = locLocale

	return pc, nil
}

func (q ClassesQ) Get(ctx context.Context) (PlaceClass, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("build select query for %s: %w", PlaceClassesTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	return scanPlaceClass(row)
}

func (q ClassesQ) Select(ctx context.Context) ([]PlaceClass, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query for %s: %w", PlaceClassesTable, err)
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

	var out []PlaceClass
	for rows.Next() {
		pc, err := scanPlaceClass(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pc)
	}
	return out, rows.Err()
}

type UpdatePlaceClassParams struct {
	Father    *string
	Status    *string
	Icon      *string
	UpdatedAt time.Time
}

func (q ClassesQ) Update(ctx context.Context, in UpdatePlaceClassParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Father != nil {
		values["father"] = *in.Father
	}
	if in.Status != nil {
		values["status"] = *in.Status
	}
	if in.Icon != nil {
		values["icon"] = *in.Icon
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build update query for %s: %w", PlaceClassesTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete query for %s: %w", PlaceClassesTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassesQ) FilterCode(code string) ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.code": code})
	q.updater = q.updater.Where(sq.Eq{"pc.code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	return q
}

func (q ClassesQ) FilterFather(code sql.NullString) ClassesQ {
	if code.Valid == false {
		q.selector = q.selector.Where("pc.father IS NULL")
		q.updater = q.updater.Where("pc.father IS NULL")
		q.deleter = q.deleter.Where("father IS NULL")
		q.counter = q.counter.Where("father IS NULL")
		return q
	}
	q.selector = q.selector.Where(sq.Eq{"pc.father": code.Valid})
	q.updater = q.updater.Where(sq.Eq{"pc.father": code.Valid})
	q.deleter = q.deleter.Where(sq.Eq{"father": code.Valid})
	q.counter = q.counter.Where(sq.Eq{"father": code.Valid})
	return q
}

func (q ClassesQ) FilterStatus(status string) ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.status": status})
	q.updater = q.updater.Where(sq.Eq{"pc.status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	return q
}

func (q ClassesQ) FilterFatherCycle(code string) ClassesQ {
	cond := sq.Expr(
		"pc.path <@ (SELECT path FROM "+PlaceClassesTable+" WHERE code = ?) AND pc.code <> ?",
		code, code,
	)
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q ClassesQ) WithLocale(locale string) ClassesQ {
	base := PlaceClassesTable
	i18n := PlaceClassLocalesTable
	l := sanitizeLocale(locale)

	col := func(field, alias string) sq.Sqlizer {
		return sq.Expr(
			`CASE
                WHEN EXISTS (
                    SELECT 1 FROM `+i18n+` i
                    WHERE i.class = pc.code AND i.locale = ?
                )
                THEN (SELECT i.`+field+` FROM `+i18n+` i  WHERE i.class = pc.code AND i.locale = ?)
                ELSE (SELECT i2.`+field+` FROM `+i18n+` i2 WHERE i2.class = pc.code AND i2.locale = 'en')
            END AS `+alias,
			l, l,
		)
	}

	q.selector = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"pc.code",
			"pc.father",
			"pc.status",
			"pc.icon",
			"pc.path",
			"pc.created_at",
			"pc.updated_at",
		).
		Column(col("name", "loc_name")).
		Column(col("locale", "loc_locale")).
		From(base + " AS pc")
	return q
}

func (q ClassesQ) OrderBy(orderBy string) ClassesQ {
	q.selector = q.selector.OrderBy(orderBy)
	return q
}

func (q ClassesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count query for %s: %w", PlaceClassesTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var cnt uint64
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}

func (q ClassesQ) Page(limit, offset uint64) ClassesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
