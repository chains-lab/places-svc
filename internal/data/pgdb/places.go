package pgdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/google/uuid"

	"github.com/paulmach/orb"
)

const placesTable = "places"

type placesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesQ(db *sql.DB) schemas.PlacesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &placesQ{
		db: db,
		selector: b.Select(
			"p.id",
			"p.city_id",
			"p.distributor_id",
			"p.class",
			"p.status",
			"p.verified",
			"ST_X(p.point::geometry) AS point_lon",
			"ST_Y(p.point::geometry) AS point_lat",
			"p.address",
			"p.website",
			"p.phone",
			"p.created_at",
			"p.updated_at",
		).From(placesTable + " AS p"),
		inserter: b.Insert(placesTable),
		updater:  b.Update(placesTable + " AS p"),
		deleter:  b.Delete(placesTable + " AS p"),
		counter:  b.Select("COUNT(DISTINCT p.id) AS count").From(placesTable + " AS p"),
	}
}

func scanPlaceRow(scanner interface{ Scan(dest ...any) error }) (schemas.Place, error) {
	var (
		p        schemas.Place
		lon, lat float64
	)
	if err := scanner.Scan(
		&p.ID,
		&p.CityID,
		&p.DistributorID,
		&p.Class,
		&p.Status,
		&p.Verified,
		&lon,
		&lat,
		&p.Address,
		&p.Website,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		return schemas.Place{}, err
	}

	p.Point = orb.Point{lon, lat}

	return p, nil
}

func scanPlaceWihDetails(scanner interface{ Scan(dest ...any) error }) (schemas.PlaceWithDetails, error) {
	var (
		p         schemas.Place
		lon, lat  float64
		locLocale string
		locName   string
		locDesc   string
		ttJSON    []byte
	)

	if err := scanner.Scan(
		&p.ID,
		&p.CityID,
		&p.DistributorID,
		&p.Class,
		&p.Status,
		&p.Verified,
		&lon,
		&lat,
		&p.Address,
		&p.Website,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
		&locLocale,
		&locName,
		&locDesc,
		&ttJSON, // ← агрегированное расписание
	); err != nil {
		return schemas.PlaceWithDetails{}, err
	}

	p.Point = orb.Point{lon, lat}

	var tt []schemas.PlaceTimetable
	if len(ttJSON) > 0 {
		if err := json.Unmarshal(ttJSON, &tt); err != nil {
			return schemas.PlaceWithDetails{}, fmt.Errorf("unmarshal timetable: %w", err)
		}
	}

	return schemas.PlaceWithDetails{
		Place:       p,
		Locale:      locLocale,
		Name:        locName,
		Description: locDesc,
		Timetable:   tt,
	}, nil
}

func (q *placesQ) Insert(ctx context.Context, in schemas.Place) error {
	stmt := map[string]interface{}{
		"id":             in.ID,
		"city_id":        in.CityID,
		"class":          in.Class,
		"status":         in.Status,
		"verified":       in.Verified,
		"point":          sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", in.Point[0], in.Point[1]),
		"address":        in.Address,
		"created_at":     in.CreatedAt,
		"updated_at":     in.UpdatedAt,
		"distributor_id": in.DistributorID,
	}

	if in.DistributorID.Valid {
		stmt["distributor_id"] = in.DistributorID.UUID
	} else {
		stmt["distributor_id"] = nil
	}
	if in.Website.Valid {
		stmt["website"] = in.Website.String
	} else {
		stmt["website"] = nil
	}
	if in.Phone.Valid {
		stmt["phone"] = in.Phone.String
	} else {
		stmt["phone"] = nil
	}

	query, args, err := q.inserter.SetMap(stmt).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", placesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q *placesQ) Get(ctx context.Context) (schemas.Place, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.Place{}, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	return scanPlaceRow(row)
}

func (q *placesQ) Select(ctx context.Context) ([]schemas.Place, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placesTable, err)
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

	var out []schemas.Place
	for rows.Next() {
		m, err := scanPlaceRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	return out, rows.Err()
}

func (q *placesQ) Update(ctx context.Context, params schemas.UpdatePlaceParams) error {
	upd := q.updater.Set("updated_at", params.UpdatedAt)

	if params.Class != nil {
		upd = upd.Set("class", *params.Class)
	}
	if params.Status != nil {
		upd = upd.Set("status", *params.Status)
	}
	if params.Verified != nil {
		upd = upd.Set("verified", *params.Verified)
	}
	if params.Point != nil {
		lon, lat := (*params.Point)[0], (*params.Point)[1]
		upd = upd.Set("point", sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", lon, lat))
	}
	if params.Address != nil {
		upd = upd.Set("address", *params.Address)
	}
	if params.Website != nil {
		if params.Website.Valid {
			upd = upd.Set("website", params.Website.String)
		} else {
			upd = upd.Set("website", nil)
		}
	}
	if params.Phone != nil {
		if params.Phone.Valid {
			upd = upd.Set("phone", params.Phone.String)
		} else {
			upd = upd.Set("phone", nil)
		}
	}

	query, args, err := upd.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", placesTable, err)
	}
	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q *placesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", placesTable, err)
	}

	if tx, ok := TxFromCtx(ctx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q *placesQ) FilterID(id uuid.UUID) schemas.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.id": id})
	q.counter = q.counter.Where(sq.Eq{"p.id": id})
	q.updater = q.updater.Where(sq.Eq{"p.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"p.id": id})

	return q
}

func (q *placesQ) FilterCityID(cityID ...uuid.UUID) schemas.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"p.city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"p.city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"p.city_id": cityID})

	return q
}

func (q *placesQ) FilterDistributorID(distributorID ...uuid.UUID) schemas.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.distributor_id": distributorID})
	q.counter = q.counter.Where(sq.Eq{"p.distributor_id": distributorID})
	q.updater = q.updater.Where(sq.Eq{"p.distributor_id": distributorID})
	q.deleter = q.deleter.Where(sq.Eq{"p.distributor_id": distributorID})

	return q
}

func (q *placesQ) FilterClass(codes ...string) schemas.PlacesQ {
	if len(codes) == 0 {
		return q
	}

	ph := make([]byte, 0, len(codes)*2-1)
	for i := range codes {
		if i > 0 {
			ph = append(ph, ',')
		}
		ph = append(ph, '?')
	}

	args := make([]any, len(codes))
	for i, v := range codes {
		args[i] = v
	}

	cte := `
		WITH RECURSIVE cls(code) AS (
		    SELECT pc.code
		    FROM ` + placeClassesTable + ` pc
		    WHERE pc.code IN (` + string(ph) + `)
		  UNION ALL
		    SELECT pc2.code
		    FROM ` + placeClassesTable + ` pc2
		    JOIN cls ON pc2.parent = cls.code
		)
		SELECT 1 FROM cls WHERE cls.code = p.class
	`

	cond := sq.Expr("EXISTS ("+cte+")", args...)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)

	return q
}

func (q *placesQ) FilterStatus(status ...string) schemas.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.status": status})
	q.counter = q.counter.Where(sq.Eq{"p.status": status})
	q.updater = q.updater.Where(sq.Eq{"p.status": status})
	q.deleter = q.deleter.Where(sq.Eq{"p.status": status})

	return q
}

func (q *placesQ) FilterVerified(verified bool) schemas.PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.verified": verified})
	q.counter = q.counter.Where(sq.Eq{"p.verified": verified})
	q.updater = q.updater.Where(sq.Eq{"p.verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"p.verified": verified})

	return q
}

func (q *placesQ) FilterWithinRadiusMeters(point orb.Point, radiusM uint64) schemas.PlacesQ {
	p := sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", point[0], point[1])
	cond := sq.Expr("ST_DWithin(p.point, ?, ?)", p, radiusM)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)

	return q
}

func (q *placesQ) FilterWithinBBox(minLon, minLat, maxLon, maxLat float64) schemas.PlacesQ {
	env := sq.Expr("ST_MakeEnvelope(?, ?, ?, ?, 4326)", minLon, minLat, maxLon, maxLat)
	cond := sq.Expr("ST_Within(p.point::geometry, ?)", env)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q *placesQ) FilterWithinPolygonWKT(polyWKT string) schemas.PlacesQ {
	poly := sq.Expr("ST_SetSRID(ST_GeomFromText(?), 4326)", polyWKT)
	cond := sq.Expr("ST_Within(p.point::geometry, ?)", poly)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q *placesQ) FilterNameLike(name string) schemas.PlacesQ {
	pattern := "%" + name + "%"
	sub := sq.Select("1").
		From(placeLocalizationTable+" pd").
		Where("pd.place_id = p.id").
		Where("pd.name ILIKE ?", pattern)

	q.selector = q.selector.Where(sq.Expr("EXISTS (?)", sub))
	q.counter = q.counter.Where(sq.Expr("EXISTS (?)", sub))

	q.updater = q.updater.Where(sq.Expr("EXISTS (?)", sub))
	q.deleter = q.deleter.Where(sq.Expr("EXISTS (?)", sub))
	return q
}

func (q *placesQ) FilterAddressLike(addr string) schemas.PlacesQ {
	q.selector = q.selector.Where("p.address ILIKE ?", "%"+addr+"%")
	q.counter = q.counter.Where("p.address ILIKE ?", "%"+addr+"%")

	return q
}

func (q *placesQ) FilterTimetableBetween(start, end int) schemas.PlacesQ {
	const week = 7 * 24 * 60
	norm := func(x int) int {
		x %= week
		if x < 0 {
			x += week
		}
		return x
	}
	s, e := norm(start), norm(end)
	if s == e {
		return q
	}

	buildOverlap := func(alias string) sq.Sqlizer {
		colS := alias + ".start_min"
		colE := alias + ".end_min"
		if s < e {
			return sq.And{sq.Lt{colS: e}, sq.Gt{colE: s}}
		}
		return sq.Or{sq.Gt{colE: s}, sq.Lt{colS: e}}
	}

	sub := sq.Select("1").
		From(placeTimetablesTable + " pt").
		Where("pt.place_id = p.id").
		Where(buildOverlap("pt"))

	q.selector = q.selector.Where(sq.Expr("EXISTS (?)", sub))
	q.counter = q.counter.Where(sq.Expr("EXISTS (?)", sub))
	q.updater = q.updater.Where(sq.Expr("EXISTS (?)", sub))
	q.deleter = q.deleter.Where(sq.Expr("EXISTS (?)", sub))
	return q
}

func (q *placesQ) WithLocale(locale string) schemas.PlacesQ {
	l := sanitizeLocale(locale)

	col := func(field, alias string, def any) sq.Sqlizer {
		return sq.Expr(
			`COALESCE(
			   (SELECT i.`+field+`
			      FROM `+placeLocalizationTable+` i
			     WHERE i.place_id = p.id
			     ORDER BY CASE
			       WHEN i.locale = ?     THEN 0
			       WHEN i.locale = 'en'  THEN 1
			       ELSE 2
			     END
			     LIMIT 1),
			   ?
			 ) AS `+alias,
			l, def,
		)
	}

	q.selector = q.selector.
		Column(col("locale", "loc_locale", "en")).
		Column(col("name", "loc_name", "")).
		Column(col("description", "loc_description", ""))

	return q
}

func (q *placesQ) WithTimetable() schemas.PlacesQ {
	q.selector = q.selector.
		LeftJoin("LATERAL (" +
			"SELECT json_agg(json_build_object(" +
			" 'id', pt.id, 'place_id', pt.place_id, 'start_min', pt.start_min, 'end_min', pt.end_min" +
			") ORDER BY pt.start_min) AS tt_json " +
			"FROM " + placeTimetablesTable + " pt WHERE pt.place_id = p.id" +
			") tt ON TRUE").
		Column("COALESCE(tt.tt_json, '[]'::json) AS tt_json")
	return q
}

func (q *placesQ) GetWithDetails(ctx context.Context, locale string) (schemas.PlaceWithDetails, error) {
	qq := *q
	qq.WithLocale(locale)
	qq.WithTimetable()

	query, args, err := qq.selector.Limit(1).ToSql()
	if err != nil {
		return schemas.PlaceWithDetails{}, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	var row *sql.Row
	if tx, ok := TxFromCtx(ctx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	return scanPlaceWihDetails(row)
}

func (q *placesQ) SelectWithDetails(ctx context.Context, locale string) ([]schemas.PlaceWithDetails, error) {
	qq := *q
	qq.WithLocale(locale)
	qq.WithTimetable()

	query, args, err := qq.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placesTable, err)
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

	var out []schemas.PlaceWithDetails
	for rows.Next() {
		item, err := scanPlaceWihDetails(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (q *placesQ) OrderByCreatedAt(asc bool) schemas.PlacesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}

	q.selector = q.selector.OrderBy("p.created_at " + dir)

	return q
}

func (q *placesQ) OrderByDistance(point orb.Point, asc bool) schemas.PlacesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}

	q.selector = q.selector.OrderByClause(
		"ST_Distance(p.point, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) "+dir,
		point[0], point[1],
	)
	return q
}

func (q *placesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", placesTable, err)
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

func (q *placesQ) Page(limit, offset uint64) schemas.PlacesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
