package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/places-svc/internal/data/schemas"
)

const placeClassesTable = "place_classes"

type classesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewClassesQ(db *sql.DB) schemas.ClassesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &classesQ{
		db: db,
		selector: b.Select(
			"pc.code",
			"pc.parent",
			"pc.status",
			"pc.icon",
			"pc.name",
			"pc.path",
			"pc.created_at",
			"pc.updated_at",
		).From(placeClassesTable + " AS pc"),
		inserter: b.Insert(placeClassesTable),
		updater:  b.Update(placeClassesTable + " AS pc"),
		deleter:  b.Delete(placeClassesTable + " AS pc"),
		counter:  b.Select("COUNT(*) AS count").From(placeClassesTable + " AS pc"),
	}
}

func scanPlaceClass(scanner interface{ Scan(dest ...any) error }) (schemas.Class, error) {
	var pc schemas.Class
	if err := scanner.Scan(
		&pc.Code,
		&pc.Parent,
		&pc.Status,
		&pc.Icon,
		&pc.Name,
		&pc.Path,
		&pc.CreatedAt,
		&pc.UpdatedAt,
	); err != nil {
		return schemas.Class{}, err
	}
	return pc, nil
}

func scanPlaceClassWithLocale(scanner interface{ Scan(dest ...any) error }) (schemas.Class, error) {
	var pc schemas.Class
	if err := scanner.Scan(
		&pc.Code,
		&pc.Parent,
		&pc.Status,
		&pc.Icon,
		&pc.Path,
		&pc.CreatedAt,
		&pc.UpdatedAt,
		&pc.Name,
	); err != nil {
		return schemas.Class{}, err
	}
	return pc, nil
}

func (q *classesQ) Insert(ctx context.Context, in schemas.Class) error {
	values := map[string]any{
		"code":   in.Code,
		"status": in.Status,
		"icon":   in.Icon,
		"path":   in.Path,
		"name":   in.Name,
	}
	if in.Parent.Valid {
		values["parent"] = in.Parent.String
	} else {
		values["parent"] = nil
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", placeClassesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *classesQ) Get(ctx context.Context) (schemas.Class, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.Class{}, fmt.Errorf("build select %s: %w", placeClassesTable, err)
	}
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	return scanPlaceClass(row)
}

func (q *classesQ) Select(ctx context.Context) ([]schemas.Class, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", placeClassesTable, err)
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

	var out []schemas.Class
	for rows.Next() {
		pc, err := scanPlaceClass(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pc)
	}
	return out, rows.Err()
}

func (q *classesQ) Update(ctx context.Context, params schemas.UpdateClassParams) error {
	values := map[string]any{
		"updated_at": params.UpdatedAt,
	}
	if params.Parent != nil {
		if params.Parent.Valid {
			values["parent"] = params.Parent.String
		} else {
			values["parent"] = nil
		}
	}
	if params.Status != nil {
		values["status"] = *params.Status
	}
	if params.Icon != nil {
		values["icon"] = *params.Icon
	}
	if params.Name != nil {
		values["name"] = *params.Name
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", placeClassesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *classesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", placeClassesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *classesQ) FilterCode(code string) schemas.ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.code": code})
	q.updater = q.updater.Where(sq.Eq{"pc.code": code})
	q.deleter = q.deleter.Where(sq.Eq{"pc.code": code})
	q.counter = q.counter.Where(sq.Eq{"pc.code": code})
	return q
}

func (q *classesQ) FilterParent(parent sql.NullString) schemas.ClassesQ {
	if !parent.Valid {
		q.selector = q.selector.Where("pc.parent IS NULL")
		q.updater = q.updater.Where("pc.parent IS NULL")
		q.deleter = q.deleter.Where("pc.parent IS NULL")
		q.counter = q.counter.Where("pc.parent IS NULL")
		return q
	}
	q.selector = q.selector.Where(sq.Eq{"pc.parent": parent.String})
	q.updater = q.updater.Where(sq.Eq{"pc.parent": parent.String})
	q.deleter = q.deleter.Where(sq.Eq{"pc.parent": parent.String})
	q.counter = q.counter.Where(sq.Eq{"pc.parent": parent.String})
	return q
}

func (q *classesQ) FilterStatus(status string) schemas.ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.status": status})
	q.updater = q.updater.Where(sq.Eq{"pc.status": status})
	q.deleter = q.deleter.Where(sq.Eq{"pc.status": status})
	q.counter = q.counter.Where(sq.Eq{"pc.status": status})
	return q
}

func (q *classesQ) FilterName(name string) schemas.ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.name": name})
	q.updater = q.updater.Where(sq.Eq{"pc.name": name})
	q.deleter = q.deleter.Where(sq.Eq{"pc.name": name})
	q.counter = q.counter.Where(sq.Eq{"pc.name": name})
	return q
}

func (q *classesQ) FilterNameLike(name string) schemas.ClassesQ {
	likePattern := fmt.Sprintf("%%%s%%", name)
	q.selector = q.selector.Where(sq.Like{"pc.name": likePattern})
	q.updater = q.updater.Where(sq.Like{"pc.name": likePattern})
	q.deleter = q.deleter.Where(sq.Like{"pc.name": likePattern})
	q.counter = q.counter.Where(sq.Like{"pc.name": likePattern})
	return q
}

func (q *classesQ) FilterParentCycle(code string) schemas.ClassesQ {
	cond := sq.Expr(
		"pc.path <@ (SELECT path FROM "+placeClassesTable+" WHERE code = ?)",
		code,
	)
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q *classesQ) OrderBy(orderBy string) schemas.ClassesQ {
	q.selector = q.selector.OrderBy(orderBy)
	return q
}

func (q *classesQ) Page(limit, offset uint) schemas.ClassesQ {
	q.selector = q.selector.Limit(uint64(limit)).Offset(uint64(offset))
	return q
}

func (q *classesQ) Count(ctx context.Context) (uint, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placeClassesTable, err)
	}
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var cnt uint
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
