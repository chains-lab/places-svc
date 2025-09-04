package dbx_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/paulmach/orb"
)

const databaseURL = "postgresql://postgres:postgres@localhost:7777/postgres?sslmode=disable"

func strPtr(s string) *string { return &s }

func openDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}
	return db
}

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := db.Exec(q, args...); err != nil {
		t.Fatalf("exec failed: %v", err)
	}
}

// корректно вставляет интервал расписания с учётом перелива через границу недели
func insertTT(t *testing.T, db *sql.DB, placeID uuid.UUID, start, end int) {
	t.Helper()
	const week = 7 * 24 * 60
	norm := func(x int) int {
		x %= week
		if x < 0 {
			x += week
		}
		return x
	}
	s := norm(start)
	e := norm(end)

	makeIns := func(smin, emin int) {
		ttQ := dbx.NewPlaceTimetablesQ(db)
		ins, _ := ttQ.Insert(dbx.PlaceTimetable{
			ID:       uuid.New(),
			PlaceID:  placeID,
			StartMin: smin,
			EndMin:   emin,
		})
		sqlStr, args, _ := ins.ToSql()
		mustExec(t, db, sqlStr, args...)
	}

	switch {
	case s == e:
		// нулевое окно — ничего не вставляем (можно изменить под бизнес-логику)
		return
	case s < e:
		makeIns(s, e)
	default:
		// перелив: разбиваем на два допустимых интервала
		makeIns(s, week)
		makeIns(0, e)
	}
}

func Test_Search_ByRadius_And_Joins(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := openDB(t)

	// ЧИСТКА (детерминизм)
	mustExec(t, db, `DELETE FROM place_timetables`)
	mustExec(t, db, `DELETE FROM place_details`)
	mustExec(t, db, `DELETE FROM places`)
	mustExec(t, db, `DELETE FROM place_kinds`)
	mustExec(t, db, `DELETE FROM place_categories`)

	now := time.Now().UTC()

	// 1) Категория
	cat := dbx.PlaceCategory{
		ID:        "food",
		Name:      "Food & Drinks",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := dbx.NewCategoryQ(db).Insert(ctx, cat); err != nil {
		t.Fatalf("insert category: %v", err)
	}

	// 2) Тип
	pt := dbx.PlaceType{
		ID:         "cafe",
		CategoryID: cat.ID,
		Name:       "Cafe",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := dbx.NewPlaceTypesQ(db).Insert(ctx, pt); err != nil {
		t.Fatalf("insert type: %v", err)
	}

	// 3) Place (Киев, Майдан)
	p := dbx.Place{
		ID:            uuid.New(),
		CityID:        uuid.New(),
		DistributorID: uuid.NullUUID{}, // NULL
		TypeID:        pt.ID,
		Status:        "active",
		Verified:      true,
		Ownership:     "private",
		Point:         orb.Point{30.5234, 50.4501}, // lon, lat
		Website:       sql.NullString{String: "https://coffee.example", Valid: true},
		Phone:         sql.NullString{String: "+380441234567", Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := dbx.NewPlacesQ(db).Insert(ctx, p); err != nil {
		t.Fatalf("insert place: %v", err)
	}

	// 4) Детали
	d := dbx.PlaceDetails{
		PlaceID:     p.ID,
		Language:    "en",
		Name:        "Coffee Point Maidan",
		Address:     "Maidan Nezalezhnosti, 1",
		Description: sql.NullString{String: "Specialty coffee", Valid: true},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := dbx.NewPlaceDetailsQ(db).Insert(ctx, d); err != nil {
		t.Fatalf("insert details: %v", err)
	}

	// 5) Расписание:
	//    — обычное окно 10:00–18:00
	insertTT(t, db, p.ID, 10*60, 18*60)
	//    — «перелив» 23:00–02:00 (вставится двумя строками)
	insertTT(t, db, p.ID, 23*60, 2*60)

	plQ := dbx.NewPlacesQ(db)

	// A) Радиус ~500 м
	{
		rows, err := plQ.New().
			FilterWithinRadiusMeters(orb.Point{30.5234, 50.4501}, 500).
			Select(ctx)
		if err != nil {
			t.Fatalf("radius select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("radius select expected 1, got %d", len(rows))
		}
		n, err := plQ.New().
			FilterWithinRadiusMeters(orb.Point{30.5234, 50.4501}, 500).
			Count(ctx)
		if err != nil || n != 1 {
			t.Fatalf("radius count n=%d err=%v", n, err)
		}
	}

	// B) ILIKE name
	{
		rows, err := plQ.New().
			FilterNameLike("Coffee").
			Select(ctx)
		if err != nil {
			t.Fatalf("name like select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("name like expected 1, got %d", len(rows))
		}
	}

	// C) ILIKE address
	{
		rows, err := plQ.New().
			FilterAddressLike("Maidan").
			Select(ctx)
		if err != nil {
			t.Fatalf("address like select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("address like expected 1, got %d", len(rows))
		}
	}

	// D) Категория
	{
		rows, err := plQ.New().
			FilterCategoryID(cat.ID).
			Select(ctx)
		if err != nil {
			t.Fatalf("category select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("category expected 1, got %d", len(rows))
		}
	}

	// E) Пересечение (12:00–13:00 внутри 10–18)
	{
		rows, err := plQ.New().
			FilterTimetableBetween(12*60, 13*60).
			Select(ctx)
		if err != nil {
			t.Fatalf("timetable day select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("timetable day expected 1, got %d", len(rows))
		}
	}

	// F) Перелив (01:00–01:30 пересекает 23:00–02:00)
	{
		rows, err := plQ.New().
			FilterTimetableBetween(1*60, 1*60+30).
			Select(ctx)
		if err != nil {
			t.Fatalf("timetable wrap select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("timetable wrap expected 1, got %d", len(rows))
		}
	}

	// G) Комбо: радиус + имя + категория + verified + сорт + пагинация
	{
		rows, err := plQ.New().
			FilterWithinRadiusMeters(orb.Point{30.5234, 50.4501}, 700).
			FilterNameLike("Coffee").
			FilterCategoryID("food").
			FilterByVerified(true).
			OrderByCreatedAt(true).
			Page(10, 0).
			Select(ctx)
		if err != nil {
			t.Fatalf("combo select err: %v", err)
		}
		if len(rows) != 1 {
			t.Fatalf("combo expected 1, got %d", len(rows))
		}
	}

	// H) DistributorID IS NULL
	{
		n, err := plQ.New().
			FilterByDistributorID(uuid.NullUUID{}).
			Count(ctx)
		if err != nil {
			t.Fatalf("distributor null count err: %v", err)
		}
		if n != 1 {
			t.Fatalf("distributor null expected 1, got %d", n)
		}
	}

	// I) Update + Get
	{
		up := dbx.UpdatePlaceParams{
			Status:    strPtr("inactive"),
			Website:   &sql.NullString{String: "https://new.example", Valid: true},
			UpdatedAt: time.Now().UTC(),
		}
		if err := plQ.New().FilterByID(p.ID).Update(ctx, up); err != nil {
			t.Fatalf("update place err: %v", err)
		}
		got, err := plQ.New().FilterByID(p.ID).Get(ctx)
		if err != nil {
			t.Fatalf("get after update err: %v", err)
		}
		if got.Status != "inactive" || !got.Website.Valid || got.Website.String != "https://new.example" {
			t.Fatalf("unexpected updated place: %+v", got)
		}
	}
}
