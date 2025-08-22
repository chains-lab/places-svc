package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const placesTypesTable = "places_types"

type PlaceType struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Category  string    `db:"category"` // соответствует ENUM place_category
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type PlacesTypesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesTypesQ(db *sql.DB) PlacesTypesQ {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return PlacesTypesQ{
		db:       db,
		selector: builder.Select("*").From(placesTypesTable),
		updater:  builder.Update(placesTypesTable),
		inserter: builder.Insert(placesTypesTable),
		deleter:  builder.Delete(placesTypesTable),
		counter:  builder.Select("COUNT(*) AS count").From(placesTypesTable),
	}
}

func (q PlacesTypesQ) New() PlacesTypesQ {
	return NewPlacesTypesQ(q.db)
}

func (q PlacesTypesQ) applyConditions(conditions ...sq.Sqlizer) PlacesTypesQ {
	q.selector = q.selector.Where(conditions)
	q.counter = q.counter.Where(conditions)
	q.updater = q.updater.Where(conditions)
	q.deleter = q.deleter.Where(conditions)
	return q
}

// CRUD

func (q PlacesTypesQ) Insert(ctx context.Context, input PlaceType) error {
	values := map[string]interface{}{
		"id":         input.ID,
		"name":       input.Name,
		"category":   input.Category,
		"updated_at": input.UpdatedAt,
		"created_at": input.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building inserter query for table: %s: %w", placesTypesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q PlacesTypesQ) Get(ctx context.Context) (PlaceType, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceType{}, fmt.Errorf("building selector query for table: %s: %w", placesTypesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var pt PlaceType
	err = row.Scan(
		&pt.ID,
		&pt.Name,
		&pt.Category,
		&pt.UpdatedAt,
		&pt.CreatedAt,
	)
	if err != nil {
		return PlaceType{}, fmt.Errorf("scanning row for table: %s: %w", placesTypesTable, err)
	}
	return pt, nil
}

func (q PlacesTypesQ) Select(ctx context.Context) ([]PlaceType, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building selector query for table: %s: %w", placesTypesTable, err)
	}

	var rows *sql.Rows
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, query, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []PlaceType
	for rows.Next() {
		var pt PlaceType
		if err := rows.Scan(
			&pt.ID,
			&pt.Name,
			&pt.Category,
			&pt.UpdatedAt,
			&pt.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning row for table: %s: %w", placesTypesTable, err)
		}
		res = append(res, pt)
	}
	return res, nil
}

func (q PlacesTypesQ) Update(ctx context.Context, input map[string]any) error {
	values := map[string]any{}

	if name, ok := input["name"]; ok {
		values["name"] = name
	}
	if category, ok := input["category"]; ok {
		values["category"] = category
	}
	if updatedAt, ok := input["updated_at"]; ok {
		values["updated_at"] = updatedAt
	} else {
		values["updated_at"] = time.Now().UTC()
	}

	query, args, err := q.updater.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building updater query for table: %s: %w", placesTypesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q PlacesTypesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building deleter query for table: %s: %w", placesTypesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

// Фильтры

func (q PlacesTypesQ) FilterID(id uuid.UUID) PlacesTypesQ {
	return q.applyConditions(sq.Eq{"id": id})
}

func (q PlacesTypesQ) FilterName(name string) PlacesTypesQ {
	return q.applyConditions(sq.Eq{"name": name})
}

func (q PlacesTypesQ) FilterCategory(category string) PlacesTypesQ {
	return q.applyConditions(sq.Eq{"category": category})
}

func (q PlacesTypesQ) LikeName(name string) PlacesTypesQ {
	p := fmt.Sprintf("%%%s%%", name)
	return q.applyConditions(sq.Like{"name": p})
}

func (q PlacesTypesQ) Page(limit, offset uint64) PlacesTypesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	q.counter = q.counter.Limit(limit).Offset(offset)
	return q
}

func (q PlacesTypesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building counter query for table: %s: %w", placesTypesTable, err)
	}

	var count uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, query, args...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(ctx, query, args...).Scan(&count)
	}
	if err != nil {
		return 0, fmt.Errorf("scanning count for table: %s: %w", placesTypesTable, err)
	}
	return count, nil
}

func (q PlacesTypesQ) Transaction(fn func(ctx context.Context) error) error {
	ctx := context.Background()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctxWithTx := context.WithValue(ctx, txKey, tx)

	if err := fn(ctxWithTx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
