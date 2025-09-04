package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const placeDetailsTable = "place_details"

type PlaceDetails struct {
	PlaceID     uuid.UUID      `db:"place_id"`
	Language    string         `db:"language"`
	Name        string         `db:"name"`
	Address     string         `db:"address"`
	Description sql.NullString `db:"description"`

	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type PlaceDetailsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceDetailsQ(db *sql.DB) PlaceDetailsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlaceDetailsQ{
		db:       db,
		selector: b.Select("*").From(placeDetailsTable),
		inserter: b.Insert(placeDetailsTable),
		updater:  b.Update(placeDetailsTable),
		deleter:  b.Delete(placeDetailsTable),
		counter:  b.Select("COUNT(*) AS count").From(placeDetailsTable),
	}
}

func (q PlaceDetailsQ) New() PlaceDetailsQ { return NewPlaceDetailsQ(q.db) }

func (q PlaceDetailsQ) Insert(ctx context.Context, in PlaceDetails) error {
	values := map[string]interface{}{
		"place_id":    in.PlaceID,
		"language":    in.Language,
		"name":        in.Name,
		"address":     in.Address,
		"description": in.Description,
		"updated_at":  in.UpdatedAt,
		"created_at":  in.CreatedAt,
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build insert query for %s: %w", placeDetailsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlaceDetailsQ) Get(ctx context.Context) (PlaceDetails, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceDetails{}, fmt.Errorf("failed to build select query for %s: %w", placeDetailsTable, err)
	}

	var out PlaceDetails
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.PlaceID,
		&out.Language,
		&out.Name,
		&out.Address,
		&out.Description,
		&out.UpdatedAt,
		&out.CreatedAt,
	)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (q PlaceDetailsQ) Select(ctx context.Context) ([]PlaceDetails, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query for %s: %w", placeDetailsTable, err)
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

	var out []PlaceDetails
	for rows.Next() {
		var pd PlaceDetails
		err := rows.Scan(
			&pd.PlaceID,
			&pd.Language,
			&pd.Name,
			&pd.Address,
			&pd.Description,
			&pd.CreatedAt,
			&pd.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row for %s: %w", placeDetailsTable, err)
		}
		out = append(out, pd)
	}
	return out, rows.Err()
}

type UpdatePlaceDetailsParams struct {
	Name        *string
	Address     *string
	Description *sql.NullString
	UpdatedAt   time.Time
}

func (q PlaceDetailsQ) Update(ctx context.Context, params UpdatePlaceDetailsParams) error {
	updates := map[string]interface{}{
		"updated_at": params.UpdatedAt,
	}
	if params.Name != nil {
		updates["name"] = params.Name
	}
	if params.Address != nil {
		updates["address"] = params.Address
	}
	if params.Description != nil {
		if params.Description.Valid {
			updates["description"] = params.Description
		} else {
			updates["description"] = nil
		}
	}

	if len(updates) == 1 { // только updated_at
		return nil
	}

	query, args, err := q.updater.SetMap(updates).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query for %s: %w", placeDetailsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlaceDetailsQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete query for %s: %w", placeDetailsTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlaceDetailsQ) FilterPlaceID(id uuid.UUID) PlaceDetailsQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": id})
	q.updater = q.updater.Where(sq.Eq{"place_id": id})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": id})
	q.counter = q.counter.Where(sq.Eq{"place_id": id})
	return q
}

func (q PlaceDetailsQ) FilterByLanguage(language string) PlaceDetailsQ {
	q.selector = q.selector.Where(sq.Eq{"language": language})
	q.updater = q.updater.Where(sq.Eq{"language": language})
	q.deleter = q.deleter.Where(sq.Eq{"language": language})
	q.counter = q.counter.Where(sq.Eq{"language": language})
	return q
}

func (q PlaceDetailsQ) FilterByName(name string) PlaceDetailsQ {
	q.selector = q.selector.Where(sq.Eq{"name": name})
	q.updater = q.updater.Where(sq.Eq{"name": name})
	q.deleter = q.deleter.Where(sq.Eq{"name": name})
	q.counter = q.counter.Where(sq.Eq{"name": name})
	return q
}

func (q PlaceDetailsQ) FilterByAddress(address string) PlaceDetailsQ {
	q.selector = q.selector.Where(sq.Eq{"address": address})
	q.updater = q.updater.Where(sq.Eq{"address": address})
	q.deleter = q.deleter.Where(sq.Eq{"address": address})
	q.counter = q.counter.Where(sq.Eq{"address": address})
	return q
}

func (q PlaceDetailsQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query for %s: %w", placeDetailsTable, err)
	}

	var count int
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to scan count for %s: %w", placeDetailsTable, err)
	}

	return count, nil
}

func (q PlaceDetailsQ) Page(offset, limit uint64) PlaceDetailsQ {
	q.selector = q.selector.Offset(offset).Limit(limit)
	return q
}
