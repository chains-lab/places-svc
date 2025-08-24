package dbx

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const placeDetailsTable = "place_details"

type PlaceDetail struct {
	PlaceID     uuid.UUID `db:"place_id"`
	Language    string    `db:"language"` // RFC 5646 (e.g., "uk", "uk-UA")
	Name        string    `db:"name"`
	Address     string    `db:"address"`
	Description *string   `db:"description"` // nullable
}

type DetailsQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewDetailsQ(db *sql.DB) DetailsQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return DetailsQ{
		db: db,
		selector: b.Select(
			"place_id",
			"language",
			"name",
			"address",
			"description",
		).From(placeDetailsTable),
		updater:  b.Update(placeDetailsTable),
		inserter: b.Insert(placeDetailsTable),
		deleter:  b.Delete(placeDetailsTable),
		counter:  b.Select("COUNT(*) AS count").From(placeDetailsTable),
	}
}

func (q DetailsQ) New() DetailsQ { return NewDetailsQ(q.db) }

func (q DetailsQ) applyConditions(conds ...sq.Sqlizer) DetailsQ {
	q.selector = q.selector.Where(conds)
	q.counter = q.counter.Where(conds)
	q.updater = q.updater.Where(conds)
	q.deleter = q.deleter.Where(conds)
	return q
}

func scanPlaceDetailRow(scanner interface{ Scan(dest ...any) error }) (PlaceDetail, error) {
	var d PlaceDetail
	if err := scanner.Scan(
		&d.PlaceID,
		&d.Language,
		&d.Name,
		&d.Address,
		&d.Description,
	); err != nil {
		return PlaceDetail{}, err
	}
	return d, nil
}

func (q DetailsQ) Insert(ctx context.Context, in PlaceDetail) error {
	vals := map[string]any{
		"place_id": in.PlaceID,
		"language": in.Language,
		"name":     in.Name,
		"address":  in.Address,
	}
	if in.Description != nil {
		vals["description"] = *in.Description
	}

	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", placeDetailsTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q DetailsQ) Get(ctx context.Context) (PlaceDetail, error) {
	qry, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceDetail{}, fmt.Errorf("build select %s: %w", placeDetailsTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanPlaceDetailRow(row)
}

func (q DetailsQ) Select(ctx context.Context) ([]PlaceDetail, error) {
	qry, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", placeDetailsTable, err)
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

	var out []PlaceDetail
	for rows.Next() {
		d, err := scanPlaceDetailRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", placeDetailsTable, err)
		}
		out = append(out, d)
	}
	return out, nil
}

func (q DetailsQ) Update(ctx context.Context, in map[string]any) error {
	vals := map[string]any{}
	if v, ok := in["name"]; ok {
		vals["name"] = v
	}
	if v, ok := in["address"]; ok {
		vals["address"] = v
	}
	if v, ok := in["description"]; ok {
		vals["description"] = v
	}

	if len(vals) == 0 {
		return nil
	}

	qry, args, err := q.updater.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", placeDetailsTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q DetailsQ) Delete(ctx context.Context) error {
	qry, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", placeDetailsTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q DetailsQ) FilterPlaceID(id uuid.UUID) DetailsQ {
	return q.applyConditions(sq.Eq{"place_id": id})
}

func (q DetailsQ) FilterLanguage(lang ...string) DetailsQ {
	return q.applyConditions(sq.Eq{"language": lang})
}

func (q DetailsQ) LikeName(name string) DetailsQ {
	return q.applyConditions(sq.Expr("name ILIKE ?", fmt.Sprintf("%%%s%%", name)))
}

func (q DetailsQ) LikeAddress(address string) DetailsQ {
	return q.applyConditions(sq.Expr("address ILIKE ?", fmt.Sprintf("%%%s%%", address)))
}

func (q DetailsQ) LikeDescription(descr string) DetailsQ {
	return q.applyConditions(sq.Expr("description ILIKE ?", fmt.Sprintf("%%%s%%", descr)))
}

func (q DetailsQ) Page(limit, offset uint64) DetailsQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q DetailsQ) Count(ctx context.Context) (uint64, error) {
	qry, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placeDetailsTable, err)
	}
	var n uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	if err != nil {
		return 0, fmt.Errorf("scan count %s: %w", placeDetailsTable, err)
	}
	return n, nil
}
