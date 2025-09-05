package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceClassTable = "place_class"

type PlaceClass struct {
	Code       string          `db:"code"`
	FatherCode *sql.NullString `db:"father_code"` // NULL для корней
	Status     string          `db:"status"`
	Icon       string          `db:"icon"`
	Path       string          `db:"path"` // ltree как text
	Locale     LocaleForClass  `db:"locales"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}

type LocaleForClass struct {
	Locale string `db:"locale"`
	Name   string `db:"name"`
}

type ClassQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewClassQ(db *sql.DB) ClassQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return ClassQ{
		db:       db,
		selector: b.Select("code", "father_code", "status", "icon", "path", "created_at", "updated_at").From(PlaceClassTable),
		inserter: b.Insert(PlaceClassTable),
		updater:  b.Update(PlaceClassTable),
		deleter:  b.Delete(PlaceClassTable),
		counter:  b.Select("COUNT(*) AS count").From(PlaceClassTable),
	}
}

func (q ClassQ) New() ClassQ { return NewClassQ(q.db) }

func (q ClassQ) Insert(ctx context.Context, in PlaceClass) error {
	values := map[string]interface{}{
		"code":        in.Code,
		"father_code": in.FatherCode,
		"status":      in.Status,
		"icon":        in.Icon,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert query for %s: %w", PlaceClassTable, err)
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
	var father sql.NullString
	var locName, locLocale sql.NullString

	if err := scanner.Scan(
		&pc.Code,
		&father,
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

	// father_code → *sql.NullString (nil для корня)
	if father.Valid {
		pc.FatherCode = &father
	} else {
		pc.FatherCode = nil
	}

	// Locale (пустая структура, если не найдено ни запрошенной, ни 'en')
	if locName.Valid && locLocale.Valid {
		pc.Locale = LocaleForClass{
			Locale: locLocale.String,
			Name:   locName.String,
		}
	} else {
		pc.Locale = LocaleForClass{}
	}

	return pc, nil
}

func (q ClassQ) Get(ctx context.Context) (PlaceClass, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceClass{}, fmt.Errorf("build select query for %s: %w", PlaceClassTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	return scanPlaceClass(row)
}

func (q ClassQ) Select(ctx context.Context) ([]PlaceClass, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query for %s: %w", PlaceClassTable, err)
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
	FatherCode *string
	Status     *string
	Icon       *string
	UpdatedAt  time.Time
}

func (q ClassQ) Update(ctx context.Context, in UpdatePlaceClassParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.FatherCode != nil {
		values["father_code"] = *in.FatherCode
	}
	if in.Status != nil {
		values["status"] = *in.Status
	}
	if in.Icon != nil {
		values["icon"] = *in.Icon
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build update query for %s: %w", PlaceClassTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete query for %s: %w", PlaceClassTable, err)
	}
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassQ) FilterCode(code string) ClassQ {
	q.selector = q.selector.Where(sq.Eq{"code": code})
	q.updater = q.updater.Where(sq.Eq{"code": code})
	q.deleter = q.deleter.Where(sq.Eq{"code": code})
	q.counter = q.counter.Where(sq.Eq{"code": code})
	return q
}

func (q ClassQ) FilterFatherCode(code *string) ClassQ {
	if code == nil {
		q.selector = q.selector.Where("father_code IS NULL")
		q.updater = q.updater.Where("father_code IS NULL")
		q.deleter = q.deleter.Where("father_code IS NULL")
		q.counter = q.counter.Where("father_code IS NULL")
		return q
	}
	q.selector = q.selector.Where(sq.Eq{"father_code": *code})
	q.updater = q.updater.Where(sq.Eq{"father_code": *code})
	q.deleter = q.deleter.Where(sq.Eq{"father_code": *code})
	q.counter = q.counter.Where(sq.Eq{"father_code": *code})
	return q
}

func (q ClassQ) FilterFatherCodeCycle(code string) ClassQ {
	cond := sq.Expr(
		fmt.Sprintf("path <@ (SELECT path FROM %s WHERE code = ?) AND code <> ?", PlaceClassTable),
		code, code,
	)

	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)

	return q
}

func (q ClassQ) FilterStatus(status string) ClassQ {
	q.selector = q.selector.Where(sq.Eq{"status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})
	q.counter = q.counter.Where(sq.Eq{"status": status})
	return q
}

func (q ClassQ) WithLocale(locale string) ClassQ {
	base := PlaceClassTable
	i18n := PlaceClassLocalesTable

	subq := fmt.Sprintf(`
		LATERAL (
			SELECT i.name, i.locale
			FROM %s i
			WHERE i.class_code = %s.code
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
			base+".father_code",
			base+".status",
			base+".icon",
			base+".path",
			base+".created_at",
			base+".updated_at",
			"loc.name   AS loc_name",
			"loc.locale AS loc_locale",
		).
		From(base).
		LeftJoin(subq + " ON TRUE")

	q.selector = q.selector.PlaceholderFormat(sq.Dollar).RunWith(q.db).Suffix("", locale)

	return q
}

func (q ClassQ) OrderBy(orderBy string) ClassQ {
	q.selector = q.selector.OrderBy(orderBy)
	return q
}

func (q ClassQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count query for %s: %w", PlaceClassTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}

func (q ClassQ) Paginate(limit, offset uint64) ClassQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
