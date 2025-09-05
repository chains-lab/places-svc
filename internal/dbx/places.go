package dbx

import (
	"context"
	"database/sql"
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

	Status    string    `db:"status"`
	Verified  bool      `db:"verified"`
	Ownership string    `db:"ownership"`
	Point     orb.Point `db:"point"`

	Locale LocaleForPlace

	Website sql.NullString `db:"website"`
	Phone   sql.NullString `db:"phone"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type LocaleForPlace struct {
	Locale      string         `db:"locale"`
	Name        string         `db:"name"`
	Address     string         `db:"address"`
	Description sql.NullString `db:"description"`
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

	cols := []string{
		"p.id",
		"p.city_id",
		"p.distributor_id",
		"p.class",
		"p.status",
		"p.verified",
		"p.ownership",
		"ST_X(p.point::geometry) AS point_lon",
		"ST_Y(p.point::geometry) AS point_lat",
		"p.website",
		"p.phone",
		"p.created_at",
		"p.updated_at",
		"NULL AS loc_locale",
		"NULL AS loc_name",
		"NULL AS loc_address",
		"NULL AS loc_description",
	}

	return PlacesQ{
		db:       db,
		selector: b.Select(cols...).From(placesTable + " p"),
		inserter: b.Insert(placesTable),
		updater:  b.Update(placesTable),
		deleter:  b.Delete(placesTable),
		counter:  b.Select("COUNT(*) AS count").From(placesTable + " p"),
	}
}

func (q PlacesQ) New() PlacesQ { return NewPlacesQ(q.db) }

func scanPlaceRow(scanner interface{ Scan(dest ...any) error }) (Place, error) {
	var (
		p                                              Place
		lon, lat                                       float64
		locLocale, locName, locAddress, locDescription sql.NullString
	)
	if err := scanner.Scan(
		&p.ID,
		&p.CityID,
		&p.DistributorID,
		&p.Class, // <— было Class
		&p.Status,
		&p.Verified,
		&p.Ownership,
		&lon,
		&lat,
		&p.Website,
		&p.Phone,
		&p.CreatedAt,
		&p.UpdatedAt,
		&locLocale,
		&locName,
		&locAddress,
		&locDescription,
	); err != nil {
		return Place{}, err
	}

	if locLocale.Valid {
		p.Locale.Locale = locLocale.String
	}
	if locName.Valid {
		p.Locale.Name = locName.String
	}
	if locAddress.Valid {
		p.Locale.Address = locAddress.String
	}
	if locDescription.Valid {
		p.Locale.Description.Valid = true
		p.Locale.Description.String = locDescription.String
	}

	p.Point = orb.Point{lon, lat}

	return p, nil
}

func (q PlacesQ) Insert(ctx context.Context, in Place) error {
	values := map[string]interface{}{
		"id":             in.ID,
		"city_id":        in.CityID,
		"distributor_id": in.DistributorID,
		"class":          in.Class,
		"status":         in.Status,
		"verified":       in.Verified,
		"ownership":      in.Ownership,
		"point":          sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", in.Point[0], in.Point[1]),
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
	Class     *string
	Status    *string
	Verified  *bool
	Ownership *string
	Point     *orb.Point // [lon, lat]

	Website *sql.NullString
	Phone   *sql.NullString

	UpdatedAt time.Time
}

func (q PlacesQ) Update(ctx context.Context, p UpdatePlaceParams) error {
	values := map[string]interface{}{
		"updated_at": p.UpdatedAt,
	}
	if p.Class != nil {
		values["class"] = *p.Class
	}
	if p.Status != nil {
		values["status"] = *p.Status
	}
	if p.Verified != nil {
		values["verified"] = *p.Verified
	}
	if p.Ownership != nil {
		values["ownership"] = *p.Ownership
	}
	if p.Point != nil {
		values["point"] = sq.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography", (*p.Point)[0], (*p.Point)[1])
	}
	if p.Website != nil {
		if p.Website.Valid {
			values["website"] = p.Website.String
		} else {
			values["website"] = nil
		}
	}
	if p.Phone != nil {
		if p.Phone.Valid {
			values["phone"] = p.Phone.String
		} else {
			values["phone"] = nil
		}
	}

	if len(values) == 0 {
		return nil
	}

	query, args, err := q.updater.SetMap(values).ToSql()
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

func (q PlacesQ) FilterByID(id uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.id": id})
	q.counter = q.counter.Where(sq.Eq{"p.id": id})
	q.updater = q.updater.Where(sq.Eq{"id": id})
	q.deleter = q.deleter.Where(sq.Eq{"id": id})

	return q
}

func (q PlacesQ) FilterByCityID(cityID uuid.UUID) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.city_id": cityID})
	q.counter = q.counter.Where(sq.Eq{"p.city_id": cityID})
	q.updater = q.updater.Where(sq.Eq{"city_id": cityID})
	q.deleter = q.deleter.Where(sq.Eq{"city_id": cityID})

	return q
}

func (q PlacesQ) FilterByDistributorID(distributorID uuid.NullUUID) PlacesQ {
	if !distributorID.Valid {
		q.selector = q.selector.Where("p.distributor_id IS NULL")
		q.counter = q.counter.Where("p.distributor_id IS NULL")
		q.updater = q.updater.Where("distributor_id IS NULL")
		q.deleter = q.deleter.Where("distributor_id IS NULL")
	} else {
		q.selector = q.selector.Where(sq.Eq{"p.distributor_id": distributorID})
		q.counter = q.counter.Where(sq.Eq{"p.distributor_id": distributorID})
		q.updater = q.updater.Where(sq.Eq{"distributor_id": distributorID})
		q.deleter = q.deleter.Where(sq.Eq{"distributor_id": distributorID})
	}

	return q
}

func (q PlacesQ) FilterClass(Class ...string) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.class": Class})
	q.counter = q.counter.Where(sq.Eq{"p.class": Class})
	q.updater = q.updater.Where(sq.Eq{"class": Class})
	q.deleter = q.deleter.Where(sq.Eq{"class": Class})

	return q
}

func (q PlacesQ) FilterByStatus(status ...string) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.status": status})
	q.counter = q.counter.Where(sq.Eq{"p.status": status})
	q.updater = q.updater.Where(sq.Eq{"status": status})
	q.deleter = q.deleter.Where(sq.Eq{"status": status})

	return q
}

func (q PlacesQ) FilterByVerified(verified bool) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.verified": verified})
	q.counter = q.counter.Where(sq.Eq{"p.verified": verified})
	q.updater = q.updater.Where(sq.Eq{"verified": verified})
	q.deleter = q.deleter.Where(sq.Eq{"verified": verified})

	return q
}

func (q PlacesQ) FilterByOwnership(ownership ...string) PlacesQ {
	q.selector = q.selector.Where(sq.Eq{"p.ownership": ownership})
	q.counter = q.counter.Where(sq.Eq{"p.ownership": ownership})
	q.updater = q.updater.Where(sq.Eq{"ownership": ownership})
	q.deleter = q.deleter.Where(sq.Eq{"ownership": ownership})

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

	join := fmt.Sprintf("%s pd ON pd.place_id = p.id", placeLocalizationTable)

	q.selector = q.selector.LeftJoin(join).Where("pd.name ILIKE ?", pattern).Distinct()
	q.counter = q.counter.LeftJoin(join).Where("pd.name ILIKE ?", pattern).Distinct()

	sub := sq.Select("1").
		From(placeLocalizationTable+" pd").
		Where("pd.place_id = places.id").
		Where("pd.name ILIKE ?", pattern)

	subSQL, subArgs, _ := sub.ToSql()

	expr := sq.Expr("EXISTS ("+subSQL+")", subArgs...)

	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q PlacesQ) FilterAddressLike(addr string) PlacesQ {
	pattern := "%" + addr + "%"

	join := fmt.Sprintf("%s pd ON pd.place_id = p.id", placeLocalizationTable)

	q.selector = q.selector.LeftJoin(join).Where("pd.address ILIKE ?", pattern).Distinct()
	q.counter = q.counter.LeftJoin(join).Where("pd.address ILIKE ?", pattern).Distinct()

	sub := sq.Select("1").
		From(placeLocalizationTable+" pd").
		Where("pd.place_id = places.id").
		Where("pd.address ILIKE ?", pattern)

	subSQL, subArgs, _ := sub.ToSql()

	expr := sq.Expr("EXISTS ("+subSQL+")", subArgs...)

	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q PlacesQ) FilterTimetableBetween(start, end int) PlacesQ {
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

	buildOverlap := func(alias string) sq.Sqlizer {
		colS := alias + ".start_min"
		colE := alias + ".end_min"
		if s < e {
			return sq.And{
				sq.Lt{colS: e},
				sq.Gt{colE: s},
			}
		}
		return sq.Or{
			sq.Gt{colE: s},
			sq.Lt{colS: e},
		}
	}

	join := fmt.Sprintf("%s pt ON pt.place_id = p.id", placeTimetablesTable)

	q.selector = q.selector.LeftJoin(join).Where(buildOverlap("pt")).Distinct()
	q.counter = q.counter.LeftJoin(join).Where(buildOverlap("pt")).Distinct()

	sub := sq.
		Select("1").
		From(placeTimetablesTable + " pt").
		Where("pt.place_id = places.id").
		Where(buildOverlap("pt"))

	subSQL, subArgs, _ := sub.ToSql()
	expr := sq.Expr("EXISTS ("+subSQL+")", subArgs...)

	q.updater = q.updater.Where(expr)
	q.deleter = q.deleter.Where(expr)

	return q
}

func (q PlacesQ) WithLocale(locale string) PlacesQ {
	base := placesTable
	i18n := placeLocalizationTable

	l := sanitizeLocale(locale)

	col := func(field, alias string) sq.Sqlizer {
		return sq.Expr(
			`CASE
                WHEN EXISTS (
                    SELECT 1 FROM `+i18n+` i
                    WHERE i.place_id = p.id AND i.locale = ?
                )
                THEN (SELECT i.`+field+`  FROM `+i18n+` i  WHERE i.place_id = p.id AND i.locale = ?)
                ELSE (SELECT i2.`+field+` FROM `+i18n+` i2 WHERE i2.place_id = p.id AND i2.locale = 'en')
             END AS `+alias,
			l, l, // ← первый ? для EXISTS, второй ? для THEN-подзапроса
		)
	}

	q.selector = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"p.id",
			"p.city_id",
			"p.distributor_id",
			"p.class",
			"p.status",
			"p.verified",
			"p.ownership",
			"ST_X(p.point::geometry) AS point_lon",
			"ST_Y(p.point::geometry) AS point_lat",
			"p.website",
			"p.phone",
			"p.created_at",
			"p.updated_at",
		).
		// порядок под scanPlaceRow: loc_locale, loc_name, loc_address, loc_description
		Column(col("locale", "loc_locale")).
		Column(col("name", "loc_name")).
		Column(col("address", "loc_address")).
		Column(col("description", "loc_description")).
		From(base + " AS p")

	return q
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

	var count uint64
	var row *sql.Row
	if tx, ok := ctx.Value(TxKey).(*sql.Tx); ok {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = q.db.QueryRowContext(ctx, query, args...)
	}

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("scanning count for %s: %w", placesTable, err)
	}

	return count, nil
}

func (q PlacesQ) Page(limit, offset uint64) PlacesQ {
	q.selector = q.selector.Limit(limit).Offset(offset)

	return q
}
