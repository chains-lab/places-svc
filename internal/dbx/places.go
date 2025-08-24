package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const placesTable = "places"

type Place struct {
	ID            uuid.UUID  `db:"id"`
	DistributorID *uuid.UUID `db:"distributor_id"`
	Type          string     `db:"type"`
	Status        string     `db:"status"`
	Ownership     string     `db:"ownership"`
	Name          string     `db:"name"`
	Description   string     `db:"description"`
	Lon           float64    `db:"lon"`
	Lat           float64    `db:"lat"`
	Address       string     `db:"address"`
	Website       string     `db:"website"`
	Phone         string     `db:"phone"`
	UpdatedAt     time.Time  `db:"updated_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type PlacesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	updater  sq.UpdateBuilder
	inserter sq.InsertBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesQ(db *sql.DB) PlacesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return PlacesQ{
		db: db,
		selector: b.Select(
			"id",
			"type",
			"status",
			"name",
			"description",
			"ST_X(coords::geometry) AS lon",
			"ST_Y(coords::geometry) AS lat",
			"address",
			"website",
			"phone",
			"updated_at",
			"created_at",
		).From(placesTable),
		updater:  b.Update(placesTable),
		inserter: b.Insert(placesTable),
		deleter:  b.Delete(placesTable),
		counter:  b.Select("COUNT(*) AS count").From(placesTable),
	}
}

func scanPlaceRow(scanner interface{ Scan(dest ...any) error }) (Place, error) {
	var (
		p             Place
		distributorID *uuid.UUID
	)
	if err := scanner.Scan(
		&p.ID,
		distributorID,
		&p.Type,
		&p.Status,
		&p.Ownership,
		&p.Name,
		&p.Description,
		&p.Lon,
		&p.Lat,
		&p.Address,
		&p.Website,
		&p.Phone,
		&p.UpdatedAt,
		&p.CreatedAt,
	); err != nil {
		return Place{}, err
	}
	return p, nil
}

func (q PlacesQ) applyConditions(conditions ...sq.Sqlizer) PlacesQ {
	q.selector = q.selector.Where(conditions)
	q.counter = q.counter.Where(conditions)
	q.updater = q.updater.Where(conditions)
	q.deleter = q.deleter.Where(conditions)
	return q
}

func (q PlacesQ) New() PlacesQ { return NewPlacesQ(q.db) }

func (q PlacesQ) Insert(ctx context.Context, in Place) error {
	vals := map[string]any{
		"id":          in.ID,
		"type":        in.Type,
		"status":      in.Status,
		"ownership":   in.Ownership,
		"name":        in.Name,
		"description": in.Description,
		"coords":      sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", in.Lon, in.Lat),
		"address":     in.Address,
		"website":     in.Website,
		"phone":       in.Phone,
		"updated_at":  in.UpdatedAt,
		"created_at":  in.CreatedAt,
	}
	if in.DistributorID != nil {
		vals["distributor_id"] = in.DistributorID
	}
	qry, args, err := q.inserter.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build insert %s: %w", placesTable, err)
	}
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}
	return err
}

func (q PlacesQ) Get(ctx context.Context) (Place, error) {
	qry, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Place{}, fmt.Errorf("build select %s: %w", placesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, qry, args...)
	} else {
		row = q.db.QueryRowContext(ctx, qry, args...)
	}
	return scanPlaceRow(row)
}

func (q PlacesQ) Select(ctx context.Context) ([]Place, error) {
	qry, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select %s: %w", placesTable, err)
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

	var out []Place
	for rows.Next() {
		p, err := scanPlaceRow(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", placesTable, err)
		}
		out = append(out, p)
	}
	return out, nil
}

func (q PlacesQ) Update(ctx context.Context, in map[string]any) error {
	vals := map[string]any{}
	if v, ok := in["distributor_id"]; ok {
		vals["distributor_id"] = v
	}
	if v, ok := in["type"]; ok {
		vals["type"] = v
	}
	if v, ok := in["status"]; ok {
		vals["status"] = v
	}
	if v, ok := in["ownership"]; ok {
		vals["ownership"] = v
	}
	if v, ok := in["name"]; ok {
		vals["name"] = v
	}
	if v, ok := in["description"]; ok {
		vals["description"] = v
	}
	if v, ok := in["address"]; ok {
		vals["address"] = v
	}
	if v, ok := in["website"]; ok {
		vals["website"] = v
	}
	if v, ok := in["phone"]; ok {
		vals["phone"] = v
	}
	if lon, ok := in["lon"]; ok {
		if lat, ok2 := in["lat"]; ok2 {
			vals["coords"] = sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", lon, lat)
		}
	}
	if _, ok := in["updated_at"]; ok {
		vals["updated_at"] = in["updated_at"]
	} else {
		vals["updated_at"] = time.Now().UTC()
	}

	qry, args, err := q.updater.SetMap(vals).ToSql()
	if err != nil {
		return fmt.Errorf("build update %s: %w", placesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}

	return err
}

func (q PlacesQ) Delete(ctx context.Context) error {
	qry, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("build delete %s: %w", placesTable, err)
	}

	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, qry, args...)
	} else {
		_, err = q.db.ExecContext(ctx, qry, args...)
	}

	return err
}

func (q PlacesQ) FilterID(id uuid.UUID) PlacesQ {
	return q.applyConditions(sq.Eq{"id": id})
}

func (q PlacesQ) FilterDistributorID(distributorID uuid.UUID) PlacesQ {
	return q.applyConditions(sq.Eq{"distributor_id": distributorID})
}

func (q PlacesQ) FilterType(placeType string) PlacesQ {
	return q.applyConditions(sq.Eq{"type": placeType})
}

func (q PlacesQ) FilterStatus(v string) PlacesQ {
	return q.applyConditions(sq.Eq{"status": v})
}

func (q PlacesQ) FilterOwnership(v string) PlacesQ {
	return q.applyConditions(sq.Eq{"ownership": v})
}

func (q PlacesQ) LikeName(name string) PlacesQ {
	return q.applyConditions(sq.Expr("name ILIKE ?", fmt.Sprintf("%%%s%%", name)))
}

func (q PlacesQ) WithinRadius(lon, lat float64, meters float64) PlacesQ {
	cond := sq.Expr("ST_DWithin(coords, ST_SetSRID(ST_MakePoint(?, ?),4326)::geography, ?)", lon, lat, meters)
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}

func (q PlacesQ) OrderByDistanceFrom(lon, lat float64) PlacesQ {
	q.selector = q.selector.OrderByClause(
		sq.Expr("ST_Distance(coords, ST_SetSRID(ST_MakePoint(?, ?),4326)::geography) ASC", lon, lat),
	)
	return q
}

func (q PlacesQ) Page(limit, offset uint64) PlacesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)
	return q
}

func (q PlacesQ) Count(ctx context.Context) (uint64, error) {
	qry, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build count %s: %w", placesTable, err)
	}
	var n uint64
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		err = tx.QueryRowContext(ctx, qry, args...).Scan(&n)
	} else {
		err = q.db.QueryRowContext(ctx, qry, args...).Scan(&n)
	}
	if err != nil {
		return 0, fmt.Errorf("scan count %s: %w", placesTable, err)
	}
	return n, nil
}
