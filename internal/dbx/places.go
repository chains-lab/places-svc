package dbx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/paulmach/orb"
)

const placesTable = "places"

type Place struct {
	ID            uuid.UUID     `db:"id"`
	CityID        uuid.UUID     `db:"city_id"`
	DistributorID uuid.NullUUID `db:"distributor_id"`
	Class         string        `db:"class"`

	Status   string    `db:"status"`
	Verified bool      `db:"verified"`
	Point    orb.Point `db:"point"`
	Address  string    `db:"address"`

	Website sql.NullString `db:"website"`
	Phone   sql.NullString `db:"phone"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PlaceWithDetails struct {
	Place
	Locale      string
	Name        string
	Description string
	Timetable   []PlaceTimetable
}

type PlacesQ struct {
	db       *sql.DB
	selector sq.SelectBuilder
	inserter sq.InsertBuilder
	updater  sq.UpdateBuilder
	deleter  sq.DeleteBuilder
	counter  sq.SelectBuilder
}

func NewPlacesQ(db *sql.DB) PlacesQ {
	b := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return PlacesQ{
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
		//counter: b.Select("COUNT(*) AS count").From(placesTable + " AS p"), // DISTINCT не нужен, т.к. нет JOIN'ов по таблицам с множеством записей на place
	}
}

func (q PlacesQ) New() PlacesQ { return NewPlacesQ(q.db) }

func scanPlaceRow(scanner interface{ Scan(dest ...any) error }) (Place, error) {
	var (
		p        Place
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
		return Place{}, err
	}

	p.Point = orb.Point{lon, lat}

	return p, nil
}

func scanPlaceWihDetails(scanner interface{ Scan(dest ...any) error }) (PlaceWithDetails, error) {
	var (
		p         Place
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
		return PlaceWithDetails{}, err
	}

	p.Point = orb.Point{lon, lat}

	var tt []PlaceTimetable
	if len(ttJSON) > 0 {
		if err := json.Unmarshal(ttJSON, &tt); err != nil {
			return PlaceWithDetails{}, fmt.Errorf("unmarshal timetable: %w", err)
		}
	}

	return PlaceWithDetails{
		Place:       p,
		Locale:      locLocale,
		Name:        locName,
		Description: locDesc,
		Timetable:   tt,
	}, nil
}

func (q PlacesQ) Insert(ctx context.Context, in Place) error {
	values := map[string]interface{}{
		"id":       in.ID,
		"city_id":  in.CityID,
		"class":    in.Class,
		"status":   in.Status,
		"verified": in.Verified,
		"point":    sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", in.Point[0], in.Point[1]),
		"address":  in.Address,
	}
	if in.Website.Valid {
		values["website"] = in.Website.String
	}
	if in.Phone.Valid {
		values["phone"] = in.Phone.String
	}
	if !in.CreatedAt.IsZero() {
		values["created_at"] = in.CreatedAt
	}
	if !in.UpdatedAt.IsZero() {
		values["updated_at"] = in.UpdatedAt
	}
	if in.DistributorID.Valid {
		values["distributor_id"] = in.DistributorID.UUID
	}

	query, args, err := q.inserter.SetMap(values).ToSql()
	if err != nil {
		return fmt.Errorf("building insert query for %s: %w", placesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlacesQ) Get(ctx context.Context) (Place, error) {
	query, args, err := q.selector.Limit(1).ToSql()
	if err != nil {
		return Place{}, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	return scanPlaceRow(row)
}

func (q PlacesQ) Select(ctx context.Context) ([]Place, error) {
	query, args, err := q.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placesTable, err)
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

	var out []Place
	for rows.Next() {
		m, err := scanPlaceRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, m)
	}

	return out, rows.Err()
}

type UpdatePlaceParams struct {
	Class    *string
	Status   *string
	Verified *bool
	Point    *orb.Point // [lon, lat]
	Address  *string

	Website *sql.NullString
	Phone   *sql.NullString

	UpdatedAt time.Time
}

func (q PlacesQ) Update(ctx context.Context, p UpdatePlaceParams) error {
	upd := q.updater.Set("updated_at", p.UpdatedAt)

	if p.Class != nil {
		upd = upd.Set("class", *p.Class)
	}
	if p.Status != nil {
		upd = upd.Set("status", *p.Status)
	}
	if p.Verified != nil {
		upd = upd.Set("verified", *p.Verified)
	}
	if p.Point != nil {
		lon, lat := (*p.Point)[0], (*p.Point)[1] // lon,lat
		upd = upd.Set("point", sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", lon, lat))
	}
	if p.Address != nil {
		upd = upd.Set("address", *p.Address)
	}
	if p.Website != nil {
		if p.Website.Valid {
			upd = upd.Set("website", p.Website.String)
		} else {
			upd = upd.Set("website", nil)
		}
	}
	if p.Phone != nil {
		if p.Phone.Valid {
			upd = upd.Set("phone", p.Phone.String)
		} else {
			upd = upd.Set("phone", nil)
		}
	}

	query, args, err := upd.ToSql()
	if err != nil {
		return fmt.Errorf("building update query for %s: %w", placesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}
	return err
}

func (q PlacesQ) Delete(ctx context.Context) error {
	query, args, err := q.deleter.ToSql()
	if err != nil {
		return fmt.Errorf("building delete query for %s: %w", placesTable, err)
	}

	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = q.db.ExecContext(ctx, query, args...)
	}

	return err
}

func (q PlacesQ) FilterID(id uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.id": id})
	q.counter = q.counter.Where(sq.Eq{"p.id": id})
	q.updater = q.updater.Where(sq.Eq{"p.id": id})
	q.deleter = q.deleter.Where(sq.Eq{"p.id": id})

	return q
}

func (q PlacesQ) FilterCityID(cityID ...uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"p.city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"p.city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"p.city_id": cityID})

	return q
}

func (q PlacesQ) FilterDistributorID(distributorID ...uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.distributor_id": distributorID})
	q.counter = q.counter.Where(sq.Eq{"p.distributor_id": distributorID})
	q.updater = q.updater.Where(sq.Eq{"p.distributor_id": distributorID})
	q.deleter = q.deleter.Where(sq.Eq{"p.distributor_id": distributorID})

	return q
}

func (q PlacesQ) FilterClass(codes ...string) PlacesQ {
	if len(codes) == 0 {
		return q
	}

	// placeholders для IN (?, ?, ...)
	ph := make([]byte, 0, len(codes)*2-1)
	for i := range codes {
		if i > 0 {
			ph = append(ph, ',')
		}
		ph = append(ph, '?')
	}

	// Аргументы как []any
	args := make([]any, len(codes))
	for i, v := range codes {
		args[i] = v
	}

	// Рекурсивно строим множество: seed-ы (codes) + все их потомки
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

	// ВАЖНО: без лишних JOIN-ов, чтобы не раздувать выборку и Total
	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)

	return q
}

//func (q PlacesQ) FilterClass(codes ...string) PlacesQ {
//	if len(codes) == 0 {
//		return q
//	}
//
//	join := placeClassesTable + " pc ON pc.code = p.class"
//
//	cond := sq.Or{
//		sq.Eq{"pc.code": codes},
//		sq.Eq{"pc.parent": codes},
//	}
//
//	q.selector = q.selector.LeftJoin(join).Where(cond)
//	q.counter = q.counter.LeftJoin(join).Where(cond)
//
//	sub := sq.Select("1").
//		From(placeClassesTable + " pc").
//		Where("pc.code = p.class").
//		Where(cond)
//
//	subSQL, subArgs, _ := sub.ToSql()
//
//	q.updater = q.updater.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))
//	q.deleter = q.deleter.Where(sq.Expr("EXISTS ("+subSQL+")", subArgs...))
//	return q
//}

func (q PlacesQ) FilterStatus(status ...string) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.status": status})
	q.counter = q.counter.Where(sq.Eq{"p.status": status})
	q.updater = q.updater.Where(sq.Eq{"p.status": status})
	q.deleter = q.deleter.Where(sq.Eq{"p.status": status})

	return q
}

func (q PlacesQ) FilterVerified(verified bool) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.verified": verified})
	q.counter = q.counter.Where(sq.Eq{"p.verified": verified})
	q.updater = q.updater.Where(sq.Eq{"p.verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"p.verified": verified})

	return q
}

func (q PlacesQ) FilterWithinRadiusMeters(point orb.Point, radiusM uint64) PlacesQ {
	p := sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", point[0], point[1])
	cond := sq.Expr("ST_DWithin(p.point, ?, ?)", p, radiusM)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)

	return q
}

func (q PlacesQ) FilterWithinBBox(minLon, minLat, maxLon, maxLat float64) PlacesQ {
	env := sq.Expr("ST_MakeEnvelope(?, ?, ?, ?, 4326)", minLon, minLat, maxLon, maxLat)
	cond := sq.Expr("ST_Within(p.point::geometry, ?)", env)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q PlacesQ) FilterWithinPolygonWKT(polyWKT string) PlacesQ {
	poly := sq.Expr("ST_SetSRID(ST_GeomFromText(?), 4326)", polyWKT)
	cond := sq.Expr("ST_Within(p.point::geometry, ?)", poly)

	q.selector = q.selector.Where(cond)
	q.counter = q.counter.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	return q
}

func (q PlacesQ) FilterNameLike(name string) PlacesQ {
	pattern := "%" + name + "%"
	sub := sq.Select("1").
		From(placeLocalizationTable+" pd").
		Where("pd.place_id = p.id").
		Where("pd.name ILIKE ?", pattern)

	q.selector = q.selector.Where(sq.Expr("EXISTS (?)", sub))
	q.counter = q.counter.Where(sq.Expr("EXISTS (?)", sub))
	// updater/deleter аналогично, если нужно
	q.updater = q.updater.Where(sq.Expr("EXISTS (?)", sub))
	q.deleter = q.deleter.Where(sq.Expr("EXISTS (?)", sub))
	return q
}

//func (q PlacesQ) FilterNameLike(name string) PlacesQ {
//	pattern := "%" + name + "%"
//
//	join := fmt.Sprintf("%s pd ON pd.place_id = p.id", placeLocalizationTable)
//	q.selector = q.selector.LeftJoin(join).Where("pd.name ILIKE ?", pattern).Distinct()
//	q.counter = q.counter.LeftJoin(join).Where("pd.name ILIKE ?", pattern).Distinct()
//
//	sub := sq.Select("1").
//		From(placeLocalizationTable+" pd").
//		Where("pd.place_id = places.id").
//		Where("pd.name ILIKE ?", pattern)
//
//	q.updater = q.updater.Where(sq.Expr("EXISTS (?)", sub))
//	q.deleter = q.deleter.Where(sq.Expr("EXISTS (?)", sub))
//
//	return q
//}

func (q PlacesQ) FilterAddressLike(addr string) PlacesQ {
	q.selector = q.selector.Where("p.address ILIKE ?", "%"+addr+"%")
	q.counter = q.counter.Where("p.address ILIKE ?", "%"+addr+"%")

	return q
}

//func (q PlacesQ) FilterTimetableBetween(start, end int) PlacesQ {
//	const week = 7 * 24 * 60 // 10080
//
//	norm := func(x int) int {
//		x %= week
//		if x < 0 {
//			x += week
//		}
//		return x
//	}
//	s := norm(start)
//	e := norm(end)
//
//	if s == e {
//		return q
//	}
//
//	buildOverlap := func(alias string) sq.Sqlizer {
//		colS := alias + ".start_min"
//		colE := alias + ".end_min"
//		if s < e {
//			// [s, e) обычный случай
//			return sq.And{
//				sq.Lt{colS: e}, // start < e
//				sq.Gt{colE: s}, // end   > s
//			}
//		}
//		// Перелом недели: [s, 10080) ∪ [0, e)
//		return sq.Or{
//			sq.Gt{colE: s}, // кусок до конца недели
//			sq.Lt{colS: e}, // кусок с начала недели
//		}
//	}
//
//	// Для selector/counter — JOIN + DISTINCT
//	join := fmt.Sprintf("%s pt ON pt.place_id = p.id", placeTimetablesTable)
//	q.selector = q.selector.LeftJoin(join).Where(buildOverlap("pt")).Distinct()
//	q.counter = q.counter.LeftJoin(join).Where(buildOverlap("pt")).Distinct()
//
//	// Для updater/deleter — EXISTS (подзапрос)
//	sub := sq.Select("1").
//		From(placeTimetablesTable + " pt").
//		Where("pt.place_id = places.id").
//		Where(buildOverlap("pt"))
//
//	q.updater = q.updater.Where(sq.Expr("EXISTS (?)", sub))
//	q.deleter = q.deleter.Where(sq.Expr("EXISTS (?)", sub))
//
//	return q
//}

func (q PlacesQ) FilterTimetableBetween(start, end int) PlacesQ {
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

func (q PlacesQ) withLocale(locale string) PlacesQ {
	l := sanitizeLocale(locale)

	// выбираем запись i18n по приоритету: нужная → 'en' → любая
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

func (q PlacesQ) withTimetable() PlacesQ {
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

func (q PlacesQ) GetWithDetails(ctx context.Context, locale string) (PlaceWithDetails, error) {
	qq := q.withLocale(locale).withTimetable()

	query, args, err := qq.selector.Limit(1).ToSql()
	if err != nil {
		return PlaceWithDetails{}, fmt.Errorf("building select query for %s: %w", placesTable, err)
	}

	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}
	return scanPlaceWihDetails(row)
}

func (q PlacesQ) SelectWithDetails(ctx context.Context, locale string) ([]PlaceWithDetails, error) {
	qq := q.withLocale(locale).withTimetable()

	query, args, err := qq.selector.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query for %s: %w", placesTable, err)
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

	var out []PlaceWithDetails
	for rows.Next() {
		item, err := scanPlaceWihDetails(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (q PlacesQ) OrderByCreatedAt(asc bool) PlacesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}

	q.selector = q.selector.OrderBy("p.created_at " + dir)

	return q
}

func (q PlacesQ) OrderByDistance(point orb.Point, asc bool) PlacesQ {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}

	geog := sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", point[0], point[1])

	q.selector = q.selector.OrderByClause("ST_Distance(p.point, ?) "+dir, geog)

	return q
}

func (q PlacesQ) Count(ctx context.Context) (uint64, error) {
	query, args, err := q.counter.ToSql()
	if err != nil {
		return 0, fmt.Errorf("building count query for %s: %w", placesTable, err)
	}
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
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

func (q PlacesQ) Page(limit, offset uint64) PlacesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
