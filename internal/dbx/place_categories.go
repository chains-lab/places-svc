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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
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
		"id":         in.ID,
		"name":       in.Name,
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

func (q CategoryQ) Get(ctx context.Context) (PlaceCategory, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceCategory{}, fmt.Errorf("build select query for %s: %w", placeCategoriesTable, err)
	}

	var out PlaceCategory
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.ID,
		&out.Name,
		&out.CreatedAt,
		&out.UpdatedAt,
	)

	return out, err
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
		var m PlaceCategory
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	return out, rows.Err()
}

type UpdatePlaceCategoryParams struct {
	Name      *string
	UpdatedAt time.Time
}

func (q CategoryQ) Update(ctx context.Context, in UpdatePlaceCategoryParams) error {
	values := map[string]interface{}{
		"updated_at": in.UpdatedAt,
	}
	if in.Name != nil {
		values["name"] = *in.Name
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

func (q CategoryQ) FilterByID(id string) CategoryQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})

	return q
}

func (q CategoryQ) FilterNameLike(name string) CategoryQ {
	q.selector = q.selector.Where("name ILIKE ?", "%"+name+"%")
	q.counter = q.counter.Where("name ILIKE ?", "%"+name+"%")

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

func (q CategoryQ) Paginate(limit, offset uint64) CategoryQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
