package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const placeTypesTable = "place_types"

type PlaceType struct {
	Name      string    `db:"name"`
	Category  string    `db:"category"` // ENUM place_category в БД → string в Go
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}

type TypesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewTypesQ(db *sql.DB) TypesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return TypesQ{
		db: db,
		selector: b.Select(
			"name",
			"category",
			"updated_at",
			"created_at",
		).From(placeTypesTable),
		updater:  b.Update(placeTypesTable),
		inserter: b.Insert(placeTypesTable),
		deleter:  b.Delete(placeTypesTable),
		counter:  b.Select("COUNT(*) AS count").From(placeTypesTable),
	}
}

func (q TypesQ) New() TypesQ { return NewTypesQ(q.db) }

func (q TypesQ) applyConditions(conds ...sq.Sqlizer) TypesQ {
	q.selector = q.selector.Where(conds)
	q.counter = q.counter.Where(conds)
	q.updater = q.updater.Where(conds)
	q.deleter = q.deleter.Where(conds)
	return q
}

func scanPlaceTypeRow(scanner interface{ Scan(dest ...any) error }) (PlaceType, error) {
	var t PlaceType
	if err := scanner.Scan(
		&t.Name,
		&t.Category,
		&t.UpdatedAt,
		&t.CreatedAt,
	); err != nil {
		return PlaceType{}, err
	}
	return t, nil
}

func (q TypesQ) Insert(ctx context.Context, in PlaceType) error {
	vals := map[string]any{
		"name":     in.Name,
		"category": in.Category,
	}
	if !in.UpdatedAt.IsZero() {
		vals["updated_at"] = in.UpdatedAt
	}
	if !in.CreatedAt.IsZero() {
		vals["created_at"] = in.CreatedAt
	}

	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", placeTypesTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q TypesQ) Get(ctx context.Context) (PlaceType, error) {
	qry, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceType{}, fmt.Errorf("build select %s: %w", placeTypesTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanPlaceTypeRow(row)
}

func (q TypesQ) Select(ctx context.Context) ([]PlaceType, error) {
	qry, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", placeTypesTable, err)
	}
	var rows *sql.Rows
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		rows, err = tx.QueryContext(ctx, qry, args...)
	} else {
		rows, err = q.db.QueryContext(ctx, qry, args...)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PlaceType
	for rows.Next() {
		t, err := scanPlaceTypeRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", placeTypesTable, err)
		}
		out = append(out, t)
	}
	return out, nil
}

func (q TypesQ) Update(ctx context.Context, in map[string]any) error {
	vals := map[string]any{}
	if v, ok := in["name"]; ok { // переименование типа (аккуратно!)
		vals["name"] = v
	}
	if v, ok := in["category"]; ok {
		vals["category"] = v
	}
	if v, ok := in["updated_at"]; ok {
		vals["updated_at"] = v
	} else {
		vals["updated_at"] = time.Now().UTC()
	}

	qry, args, err := q.updater.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", placeTypesTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q TypesQ) Delete(ctx context.Context) error {
	qry, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", placeTypesTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q TypesQ) FilterName(name string) TypesQ {
	return q.applyConditions(sq.Eq{"name": name})
}

func (q TypesQ) FilterCategory(cat string) TypesQ {
	return q.applyConditions(sq.Eq{"category": cat})
}

func (q TypesQ) LikeName(substr string) TypesQ {
	return q.applyConditions(sq.Expr("name ILIKE ?", "%"+substr+"%"))
}

func (q TypesQ) Page(limit, offset uint64) TypesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	// counter намеренно без Limit/Offset
	return q
}

func (q TypesQ) Count(ctx context.Context) (uint64, error) {
	qry, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placeTypesTable, err)
	}
	var n uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	if err != nil {
		return 0, fmt.Errorf("scan count %s: %w", placeTypesTable, err)
	}
	return n, nil
}
