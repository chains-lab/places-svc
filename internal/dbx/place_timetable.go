package dbx

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const placeTimetablesTable = "place_timetables"

type PlaceTimetable struct {
	ID       uuid.UUID `db:"id"`
	PlaceID  uuid.UUID `db:"place_id"`
	StartMin int       `db:"start_min"`
	EndMin   int       `db:"end_min"`
}

type PlaceTimetablesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceTimetablesQ(db *sql.DB) PlaceTimetablesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlaceTimetablesQ{
		db:       db,
		selector: b.Select("*").From(placeTimetablesTable),
		inserter: b.Insert(placeTimetablesTable),
		updater:  b.Update(placeTimetablesTable),
		deleter:  b.Delete(placeTimetablesTable),
		counter:  b.Select("COUNT(*) AS count").From(placeTimetablesTable),
	}
}

func (q PlaceTimetablesQ) New() PlaceTimetablesQ { return NewPlaceTimetablesQ(q.db) }

func (q PlaceTimetablesQ) Insert(in PlaceTimetable) (sq.InsertBuilder, error) {
	values := map[string]interface{}{
		"id":        in.ID,
		"place_id":  in.PlaceID,
		"start_min": in.StartMin,
		"end_min":   in.EndMin,
	}

	query := q.inserter.SetMap(values)
	return query, nil
}

func (q PlaceTimetablesQ) Get(ctx context.Context) (PlaceTimetable, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceTimetable{}, err
	}

	var out PlaceTimetable
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	err = row.Scan(
		&out.ID,
		&out.PlaceID,
		&out.StartMin,
		&out.EndMin,
	)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (q PlaceTimetablesQ) Select(ctx context.Context) ([]PlaceTimetable, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, err
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

	var out []PlaceTimetable
	for rows.Next() {
		var c PlaceTimetable
		if err = rows.Scan(
			&c.ID,
			&c.PlaceID,
			&c.StartMin,
			&c.EndMin,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (q PlaceTimetablesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build delete %s: %w", placeTimetablesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlaceTimetablesQ) FilterByID(id uuid.UUID) PlaceTimetablesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})

	return q
}

func (q PlaceTimetablesQ) FilterByPlaceID(placeID uuid.UUID) PlaceTimetablesQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": placeID})
	q.updater = q.updater.Where(sq.Eq{"place_id": placeID})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": placeID})
	q.counter = q.counter.Where(sq.Eq{"place_id": placeID})

	return q
}

func (q PlaceTimetablesQ) FilterBetween(start, end int) PlaceTimetablesQ {
	const week = 7 * 24 * 60 // 10080

	norm := func(x int) int {
		x %= week
		if x < 0 {
			x += week
		}
		return x
	}
	s := norm(start)
	e := norm(end)

	if s == e {
		return q
	}

	var cond any
	if s < e {
		cond = sq.And{
			sq.GtOrEq{"start_min": s},
			sq.LtOrEq{"end_min": e},
		}
	} else {
		cond = sq.Or{
			sq.And{sq.GtOrEq{"start_min": s}},
			sq.And{sq.LtOrEq{"end_min": e}},
		}
	}

	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)

	return q
}

func (q PlaceTimetablesQ) Count(ctx context.Context) (int, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query for %s: %w", placeTimetablesTable, err)
	}

	var count int
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count result for %s: %w", placeTimetablesTable, err)
	}

	return count, nil
}

func (q PlaceTimetablesQ) Page(offset, limit uint64) PlaceTimetablesQ {
	q.selector = q.selector.Offset(offset).Limit(limit)

	return q
}
