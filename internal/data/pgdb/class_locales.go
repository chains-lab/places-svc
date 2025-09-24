package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data/schemas"
)

const classLocalesTable = "place_class_i18n"

var reLocale = regexp.MustCompile(`^[a-z]{2}$`)

func sanitizeLocale(l string) string {
	if reLocale.MatchString(l) {
		return l
	}

	return enum.LocaleEN
}

type classLocalesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewClassLocalesQ(db *sql.DB) schemas.ClassLocalesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &classLocalesQ{
		db:       db,
		selector: b.Select("*").From(classLocalesTable),
		inserter: b.Insert(classLocalesTable),
		updater:  b.Update(classLocalesTable),
		deleter:  b.Delete(classLocalesTable),
		counter:  b.Select("COUNT(*) AS count").From(classLocalesTable),
	}
}

func (q *classLocalesQ) Insert(ctx context.Context, in ...schemas.ClassLocale) error {
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
		return fmt.Errorf("building insert query for %s: %w", classLocalesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *classLocalesQ) Upsert(ctx context.Context, in ...schemas.ClassLocale) error {
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
	`, classLocalesTable, strings.Join(ph, ","))

	var err error
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *classLocalesQ) Get(ctx context.Context) (schemas.ClassLocale, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.ClassLocale{}, fmt.Errorf("building select query for %s: %w", classLocalesTable, err)
	}

	var out schemas.ClassLocale
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
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

func (q *classLocalesQ) Select(ctx context.Context) ([]schemas.ClassLocale, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", classLocalesTable, err)
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

	var out []schemas.ClassLocale
	for rows.Next() {
		var item schemas.ClassLocale
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

func (q *classLocalesQ) Update(ctx context.Context, in schemas.UpdateClassLocaleParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Name != nil {
		values["name"] = *in.Name
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", classLocalesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q *classLocalesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", classLocalesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q *classLocalesQ) FilterClass(code string) schemas.ClassLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"class": code})
	q.counter = q.counter.Where(sq.Eq{"class": code})
	q.updater = q.updater.Where(sq.Eq{"class": code})
	q.deleter = q.deleter.Where(sq.Eq{"class": code})

	return q
}

func (q *classLocalesQ) FilterLocale(locale string) schemas.ClassLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"locale": locale})
	q.counter = q.counter.Where(sq.Eq{"locale": locale})
	q.updater = q.updater.Where(sq.Eq{"locale": locale})
	q.deleter = q.deleter.Where(sq.Eq{"locale": locale})

	return q
}

func (q *classLocalesQ) FilterNameLike(name string) schemas.ClassLocalesQ {
	q.selector = q.selector.Where(sq.Like{"name": name})
	q.counter = q.counter.Where(sq.Like{"name": name})
	q.updater = q.updater.Where(sq.Like{"name": name})
	q.deleter = q.deleter.Where(sq.Like{"name": name})

	return q
}

func (q *classLocalesQ) FilterName(name string) schemas.ClassLocalesQ {
	q.selector = q.selector.Where(sq.Eq{"name": name})
	q.counter = q.counter.Where(sq.Eq{"name": name})
	q.updater = q.updater.Where(sq.Eq{"name": name})
	q.deleter = q.deleter.Where(sq.Eq{"name": name})

	return q
}

func (q *classLocalesQ) OrderByLocale(asc bool) schemas.ClassLocalesQ {
	dir := "DESC"
	if asc {
		dir = "ASC"
	}

	q.selector = q.selector.OrderBy("locale " + dir)

	return q
}

func (q *classLocalesQ) Page(limit, offset uint) schemas.ClassLocalesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))

	return q
}

func (q *classLocalesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", classLocalesTable, err)
	}

	var count uint
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	if err = row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count from %s: %w", classLocalesTable, err)
	}

	return count, nil
}
