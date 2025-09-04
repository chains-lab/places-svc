package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const PlaceTypesTable = "place_types"

type PlaceType struct {
	ID         string    `db:"id"`
	CategoryID string    `db:"category_id"`
	Name       string    `db:"name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type TypesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceTypesQ(db *sql.DB) TypesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return TypesQ{
		db:       db,
		selector: b.Select("*").From(PlaceTypesTable),
		inserter: b.Insert(PlaceTypesTable),
		updater:  b.Update(PlaceTypesTable),
		deleter:  b.Delete(PlaceTypesTable),
		counter:  b.Select("COUNT(*) AS count").From(PlaceTypesTable),
	}
}

func (q TypesQ) New() TypesQ { return NewPlaceTypesQ(q.db) }

func (q TypesQ) Insert(ctx context.Context, in PlaceType) error {
	values := map[string]interface{}{
		"id":          in.ID,
		"category_id": in.CategoryID,
		"name":        in.Name,
		"updated_at":  in.UpdatedAt,
		"created_at":  in.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", PlaceTypesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q TypesQ) Get(ctx context.Context) (PlaceType, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceType{}, fmt.Errorf("building select query for %s: %w", PlaceTypesTable, err)
	}

	var out PlaceType
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	err = row.Scan(
		&out.ID,
		&out.CategoryID,
		&out.Name,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	return out, err
}

func (q TypesQ) Select(ctx context.Context) ([]PlaceType, error) {
	var out []PlaceType

	query, args, err := q.selector.ToSql()
	if err != nil {
		return out, fmt.Errorf("building select query for %s: %w", PlaceTypesTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return out, fmt.Errorf("querying select for %s: %w", PlaceTypesTable, err)
	}
	defer rows.Close()

	for rows.Next() {
		var t PlaceType
		if err = rows.Scan(
			&t.ID,
			&t.CategoryID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return out, fmt.Errorf("scanning row for %s: %w", PlaceTypesTable, err)
		}
		out = append(out, t)
	}

	if err = rows.Err(); err != nil {
		return out, fmt.Errorf("iterating rows for %s: %w", PlaceTypesTable, err)
	}

	return out, nil
}

type PlaceUpdateParams struct {
	Name      *string   `db:"name"`
	Category  *string   `db:"category"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (q TypesQ) Update(ctx context.Context, params PlaceUpdateParams) error {
	values := map[string]interface{}{
		"updated_at": params.UpdatedAt,
	}
	if params.Name != nil {
		values["name"] = *params.Name
	}
	if params.Category != nil {
		values["category"] = *params.Category
	}

	if len(values) == 0 {
		return nil // nothing to update
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", PlaceTypesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q TypesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", PlaceTypesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q TypesQ) FilterByID(id string) TypesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	return q
}

func (q TypesQ) FilterByCategoryID(category string) TypesQ {
	q.selector = q.selector.Where(sq.Eq{"category_id": category})
	q.updater = q.updater.Where(sq.Eq{"category_id": category})
	q.deleter = q.deleter.Where(sq.Eq{"category_id": category})
	q.counter = q.counter.Where(sq.Eq{"category_id": category})
	return q
}

func (q TypesQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", PlaceTypesTable, err)
	}

	var count int
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", PlaceTypesTable, err)
	}

	return count, nil
}

func (q TypesQ) Page(limit, offset uint64) TypesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}
