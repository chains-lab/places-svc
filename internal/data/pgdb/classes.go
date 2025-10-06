package pgdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const classesTable = "place_classes"

type Class struct {
	Code      string         `storage:"code"`
	Parent    sql.NullString `storage:"parent"` // NULL для корней
	Status    string         `storage:"status"`
	Icon      string         `storage:"icon"`
	Name      string         `storage:"name"`
	Path      string         `storage:"path"` // ltree как text
	CreatedAt time.Time      `storage:"created_at"`
	UpdatedAt time.Time      `storage:"updated_at"`
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
			"pc.parent",
			"pc.status",
			"pc.icon",
			"pc.name",
			"pc.path",
			"pc.created_at",
			"pc.updated_at",
		).From(classesTable + " AS pc"),
		inserter: b.Insert(classesTable),
		updater:  b.Update(classesTable + " AS pc"),
		deleter:  b.Delete(classesTable + " AS pc"),
		counter:  b.Select("COUNT(*) AS count").From(classesTable + " AS pc"),
	}
}

func scanPlaceClass(scanner interface{ Scan(dest ...any) error }) (Class, error) {
	var pc Class
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
		return Class{}, err
	}
	return pc, nil
}

func scanPlaceClassWithLocale(scanner interface{ Scan(dest ...any) error }) (Class, error) {
	var pc Class
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
		return Class{}, err
	}
	return pc, nil
}

func (q ClassesQ) New() ClassesQ {
	return NewClassesQ(q.db)
}

func (q ClassesQ) Insert(ctx context.Context, in Class) error {
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
		return fmt.Errorf("build insert %s: %w", classesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassesQ) Exists(ctx context.Context) (bool, error) {
	query, args, err := q.counter.Limit(1).ToSql()
	if err != nil {
		return false, fmt.Errorf("build exists %s: %w", classesTable, err)
	}
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	var cnt uint
	if err := row.Scan(&cnt); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (q ClassesQ) Get(ctx context.Context) (Class, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Class{}, fmt.Errorf("build select %s: %w", classesTable, err)
	}
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	return scanPlaceClass(row)
}

func (q ClassesQ) Select(ctx context.Context) ([]Class, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", classesTable, err)
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

	var out []Class
	for rows.Next() {
		pc, err := scanPlaceClass(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pc)
	}
	return out, rows.Err()
}

func (q ClassesQ) Update(ctx context.Context, updatedAt time.Time) error {
	q.updater = q.updater.Set("updated_at", updatedAt)

	query, args, err := q.updater.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", classesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassesQ) UpdateParent(parent sql.NullString) ClassesQ {
	if parent.Valid {
		q.updater = q.updater.Set("pc.parent", parent.String)
	} else {
		q.updater = q.updater.Set("pc.parent", nil)
	}
	return q
}

func (q ClassesQ) UpdateStatus(status string) ClassesQ {
	q.updater = q.updater.Set("pc.status", status)
	return q
}

func (q ClassesQ) UpdateIcon(icon string) ClassesQ {
	q.updater = q.updater.Set("pc.icon", icon)
	return q
}

func (q ClassesQ) UpdateName(name string) ClassesQ {
	q.updater = q.updater.Set("pc.name", name)
	return q
}

func (q ClassesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", classesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q ClassesQ) FilterCode(code string) ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.code": code})
	q.updater = q.updater.Where(sq.Eq{"pc.code": code})
	q.deleter = q.deleter.Where(sq.Eq{"pc.code": code})
	q.counter = q.counter.Where(sq.Eq{"pc.code": code})
	return q
}

func (q ClassesQ) FilterParent(parent sql.NullString) ClassesQ {
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

func (q ClassesQ) FilterStatus(status string) ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.status": status})
	q.updater = q.updater.Where(sq.Eq{"pc.status": status})
	q.deleter = q.deleter.Where(sq.Eq{"pc.status": status})
	q.counter = q.counter.Where(sq.Eq{"pc.status": status})
	return q
}

func (q ClassesQ) FilterName(name string) ClassesQ {
	q.selector = q.selector.Where(sq.Eq{"pc.name": name})
	q.updater = q.updater.Where(sq.Eq{"pc.name": name})
	q.deleter = q.deleter.Where(sq.Eq{"pc.name": name})
	q.counter = q.counter.Where(sq.Eq{"pc.name": name})
	return q
}

func (q ClassesQ) FilterNameLike(name string) ClassesQ {
	likePattern := fmt.Sprintf("%%%s%%", name)
	q.selector = q.selector.Where(sq.Like{"pc.name": likePattern})
	q.updater = q.updater.Where(sq.Like{"pc.name": likePattern})
	q.deleter = q.deleter.Where(sq.Like{"pc.name": likePattern})
	q.counter = q.counter.Where(sq.Like{"pc.name": likePattern})
	return q
}

func (q ClassesQ) FilterParentCycle(code string) ClassesQ {
	cond := sq.Expr(
		"pc.path <@ (SELECT path FROM "+classesTable+" WHERE code = ?)",
		code,
	)
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q ClassesQ) OrderBy(orderBy string) ClassesQ {
	q.selector = q.selector.OrderBy(orderBy)
	return q
}

func (q ClassesQ) Page(limit, offset uint64) ClassesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q ClassesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", classesTable, err)
	}
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
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

func (q ClassesQ) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	_, ok := TxFromCtx(ctx)
	if ok {
		return fn(ctx)
	}

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
				err = fmt.Errorf("tx err: %v; rollback err: %v", err, rbErr)
			}
		}
	}()

	ctxWithTx := context.WithValue(ctx, TxKey, tx)

	if err = fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
