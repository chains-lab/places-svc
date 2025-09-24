package pgdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/google/uuid"
)

const placeTimetablesTable = "place_timetables"

type placeTimetablesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlaceTimetablesQ(db *sql.DB) schemas.PlaceTimetablesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &placeTimetablesQ{
		db: db,
		selector: b.Select(
			"id",
			"place_id",
			"start_min",
			"end_min",
		).From(placeTimetablesTable),
		inserter: b.Insert(placeTimetablesTable),
		updater:  b.Update(placeTimetablesTable),
		deleter:  b.Delete(placeTimetablesTable),
		counter:  b.Select("COUNT(*) AS count").From(placeTimetablesTable),
	}
}

func (q *placeTimetablesQ) New() schemas.PlaceTimetablesQ { return NewPlaceTimetablesQ(q.db) }

// ---------- CRUD ----------

func (q *placeTimetablesQ) Insert(ctx context.Context, in ...schemas.PlaceTimetable) error {
	if len(in) == 0 {
		return nil
	}

	ins := q.inserter.Columns("id", "place_id", "start_min", "end_min")
	for _, t := range in {
		ins = ins.Values(t.ID, t.PlaceID, t.StartMin, t.EndMin)
	}

	query, args, err := ins.ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", placeTimetablesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *placeTimetablesQ) Upsert(ctx context.Context, in ...schemas.PlaceTimetable) error {
	if len(in) == 0 {
		return nil
	}

	const cols = "(id, place_id, start_min, end_min)"
	var (
		args []any
		ph   []string
		i    = 1
	)
	for _, r := range in {
		ph = append(ph, fmt.Sprintf("($%d,$%d,$%d,$%d)", i, i+1, i+2, i+3))
		args = append(args, r.ID, r.PlaceID, r.StartMin, r.EndMin)
		i += 4
	}

	query := fmt.Sprintf(`
		INSERT INTO %s %s VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET place_id = EXCLUDED.place_id,
		    start_min = EXCLUDED.start_min,
		    end_min   = EXCLUDED.end_min
	`, placeTimetablesTable, cols, strings.Join(ph, ","))

	if tx, ok := TxFromCtx(ctx); ok {
		_, err := tx.ExecContext(ctx, query, args...)
		return err
	}
	_, err := q.db.ExecContext(ctx, query, args...)
	return err
}

func (q *placeTimetablesQ) Get(ctx context.Context) (schemas.PlaceTimetable, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.PlaceTimetable{}, fmt.Errorf("build select %s: %w", placeTimetablesTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	var out schemas.PlaceTimetable
	if err := row.Scan(&out.ID, &out.PlaceID, &out.StartMin, &out.EndMin); err != nil {
		return out, err
	}
	return out, nil
}

func (q *placeTimetablesQ) Select(ctx context.Context) ([]schemas.PlaceTimetable, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", placeTimetablesTable, err)
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

	var out []schemas.PlaceTimetable
	for rows.Next() {
		var t schemas.PlaceTimetable
		if err := rows.Scan(&t.ID, &t.PlaceID, &t.StartMin, &t.EndMin); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (q *placeTimetablesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", placeTimetablesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

// ---------- filters ----------

func (q *placeTimetablesQ) FilterByID(id uuid.UUID) schemas.PlaceTimetablesQ {
	q.selector = q.selector.Where(sq.Eq{"id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})
	q.counter = q.counter.Where(sq.Eq{"id": id})
	return q
}

func (q *placeTimetablesQ) FilterPlaceID(placeID uuid.UUID) schemas.PlaceTimetablesQ {
	q.selector = q.selector.Where(sq.Eq{"place_id": placeID})
	q.updater = q.updater.Where(sq.Eq{"place_id": placeID})
	q.deleter = q.deleter.Where(sq.Eq{"place_id": placeID})
	q.counter = q.counter.Where(sq.Eq{"place_id": placeID})
	return q
}

func (q *placeTimetablesQ) FilterBetween(start, end int) schemas.PlaceTimetablesQ {
	const week = 7 * 24 * 60 // 10080

	norm := func(x int) int {
		x %= week
		if x < 0 {
			x += week
		}
		return x
	}
	s, e := norm(start), norm(end)

	if s == e {
		// пустое окно — вернуть пустую выборку
		q.selector = q.selector.Where("1=0")
		q.updater = q.updater.Where("1=0")
		q.deleter = q.deleter.Where("1=0")
		q.counter = q.counter.Where("1=0")
		return q
	}

	var cond any
	if s < e {
		// [s, e): start < e AND end > s
		cond = sq.And{
			sq.Lt{"start_min": e},
			sq.Gt{"end_min": s},
		}
	} else {
		// перелом недели: (end > s) OR (start < e)
		cond = sq.Or{
			sq.Gt{"end_min": s},
			sq.Lt{"start_min": e},
		}
	}

	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

// ---------- page/count ----------

func (q *placeTimetablesQ) Page(limit, offset uint64) schemas.PlaceTimetablesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q *placeTimetablesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placeTimetablesTable, err)
	}

	var cnt uint64
	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
