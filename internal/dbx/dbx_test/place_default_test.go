package dbx_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func Test_Full_E2E_Using_Qs(t *testing.T) {
	ctx := context.Background()
	db := openDB(t)

	mustExec(t, db, `DELETE FROM place_timetables`)
	mustExec(t, db, `DELETE FROM place_details`)
	mustExec(t, db, `DELETE FROM places`)
	mustExec(t, db, `DELETE FROM place_kinds`)
	mustExec(t, db, `DELETE FROM place_categories`)

	now := time.Now().UTC()

	// ---------- seed: category ----------
	cat := dbx.PlaceCategory{ID: "test_food", Name: "Test Food", CreatedAt: now, UpdatedAt: now}
	catQ := dbx.NewCategoryQ(db)
	if err := catQ.Insert(ctx, cat); err != nil {
		t.Fatalf("insert category: %v", err)
	}
	t.Cleanup(func() { _ = dbx.NewCategoryQ(db).FilterByID(cat.ID).Delete(ctx) })

	// verify select + count + name like
	{
		got, err := catQ.New().FilterByID(cat.ID).Get(ctx)
		if err != nil {
			t.Fatalf("get category: %v", err)
		}
		if got.Name != cat.Name {
			t.Fatalf("category name mismatch: got %q want %q", got.Name, cat.Name)
		}

		n, err := catQ.New().FilterNameLike("Food").Count(ctx)
		if err != nil || n != 1 {
			t.Fatalf("count name like: n=%d err=%v", n, err)
		}
	}

	// ---------- seed: type ----------
	typ := dbx.PlaceType{ID: "test_cafe", CategoryID: cat.ID, Name: "Test Cafe", CreatedAt: now, UpdatedAt: now}
	typQ := dbx.NewPlaceTypesQ(db)
	if err := typQ.Insert(ctx, typ); err != nil {
		t.Fatalf("insert type: %v", err)
	}
	t.Cleanup(func() { _ = typQ.New().FilterByID(typ.ID).Delete(ctx) })

	// verify select by category
	{
		types, err := typQ.New().FilterByCategoryID(cat.ID).Select(ctx)
		if err != nil {
			t.Fatalf("select types by category: %v", err)
		}
		if len(types) != 1 || types[0].ID != typ.ID {
			t.Fatalf("unexpected types result: %+v", types)
		}
	}

	// ---------- seed: place ----------
	pl := dbx.Place{
		ID:            uuid.New(),
		CityID:        uuid.New(),
		DistributorID: uuid.NullUUID{}, // NULL
		TypeID:        typ.ID,
		Status:        "active",
		Verified:      true,
		Ownership:     "private",
		Point:         orb.Point{30.5234, 50.4501}, // Киев
		Website:       sql.NullString{String: "https://ex.am", Valid: true},
		Phone:         sql.NullString{String: "+380", Valid: true},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	plQ := dbx.NewPlacesQ(db)
	if err := plQ.Insert(ctx, pl); err != nil {
		t.Fatalf("insert place: %v", err)
	}
	t.Cleanup(func() { _ = plQ.New().FilterByID(pl.ID).Delete(ctx) })

	// verify Get + Select basic
	{
		got, err := plQ.New().FilterByID(pl.ID).Get(ctx)
		if err != nil {
			t.Fatalf("get place: %v", err)
		}
		if got.TypeID != pl.TypeID || got.Verified != true || got.Ownership != "private" {
			t.Fatalf("place fields mismatch: %+v", got)
		}
		list, err := plQ.New().FilterByCityID(pl.CityID).Select(ctx)
		if err != nil || len(list) != 1 {
			t.Fatalf("select by city: len=%d err=%v", len(list), err)
		}
	}

	// ---------- seed: place_details (2 языка) ----------
	pdQ := dbx.NewPlaceDetailsQ(db)
	if err := pdQ.Insert(ctx, dbx.PlaceDetails{
		PlaceID:   pl.ID,
		Language:  "en",
		Name:      "Coffee Bar",
		Address:   "Main st 1",
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("insert details en: %v", err)
	}
	if err := pdQ.Insert(ctx, dbx.PlaceDetails{
		PlaceID:   pl.ID,
		Language:  "uk",
		Name:      "Кав'ярня",
		Address:   "Головна 1",
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("insert details uk: %v", err)
	}
	t.Cleanup(func() { _ = pdQ.New().FilterPlaceID(pl.ID).Delete(ctx) })

	// ---------- seed: timetable (2 слота) ----------
	ptQ := dbx.NewPlaceTimetablesQ(db)
	ins1, err := ptQ.Insert(dbx.PlaceTimetable{
		ID: uuid.New(), PlaceID: pl.ID, StartMin: 9 * 60, EndMin: 18 * 60,
	})
	if err != nil {
		t.Fatalf("build insert timetable 1: %v", err)
	}
	ins2, err := ptQ.Insert(dbx.PlaceTimetable{
		ID: uuid.New(), PlaceID: pl.ID, StartMin: 20 * 60, EndMin: 23 * 60,
	})
	if err != nil {
		t.Fatalf("build insert timetable 2: %v", err)
	}
	// Вставляем билдеры
	{
		q1, a1, _ := ins1.ToSql()
		if _, err := db.ExecContext(ctx, q1, a1...); err != nil {
			t.Fatalf("exec timetable 1: %v", err)
		}
		q2, a2, _ := ins2.ToSql()
		if _, err := db.ExecContext(ctx, q2, a2...); err != nil {
			t.Fatalf("exec timetable 2: %v", err)
		}
	}
	t.Cleanup(func() { _ = ptQ.New().FilterByPlaceID(pl.ID).Delete(ctx) })

	// ---------- проверки фильтров НА ТВОИХ Q ----------

	// 1) NameLike (JOIN place_details) + Count
	{
		rows, err := plQ.New().FilterNameLike("Coffee").Select(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatalf("FilterNameLike Select: len=%d err=%v", len(rows), err)
		}
		n, err := plQ.New().FilterNameLike("Coffee").Count(ctx)
		if err != nil || n != 1 {
			t.Fatalf("FilterNameLike Count: n=%d err=%v", n, err)
		}
	}

	// 2) AddressLike (JOIN place_details)
	{
		rows, err := plQ.New().FilterAddressLike("Main").Select(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatalf("FilterAddressLike Select: len=%d err=%v", len(rows), err)
		}
	}

	// 3) CategoryID (JOIN place_kinds)
	{
		rows, err := plQ.New().FilterCategoryID(cat.ID).Select(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatalf("FilterCategoryID Select: len=%d err=%v", len(rows), err)
		}
	}

	// 4) TimetableBetween — окно 08:00..10:00 пересекает 09:00..18:00
	{
		rows, err := plQ.New().FilterTimetableBetween(8*60, 10*60).Select(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatalf("FilterTimetableBetween Select: len=%d err=%v", len(rows), err)
		}
		// wrap-around окно: 22:00..06:00 — пересечёт 20:00..23:00
		rows2, err := plQ.New().FilterTimetableBetween(22*60, 6*60).Select(ctx)
		if err != nil || len(rows2) != 1 {
			t.Fatalf("FilterTimetableBetween wrap Select: len=%d err=%v", len(rows2), err)
		}
	}

	// 5) WithinRadiusMeters (1 км)
	{
		rows, err := plQ.New().FilterWithinRadiusMeters(pl.Point, 1000).Select(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatalf("FilterWithinRadiusMeters: len=%d err=%v", len(rows), err)
		}
	}

	// 6) DistributorID NULL и non-NULL
	{
		// NULL
		rows, err := plQ.New().FilterByDistributorID(uuid.NullUUID{}).Select(ctx)
		if err != nil || len(rows) != 1 {
			t.Fatalf("FilterByDistributorID(NULL): len=%d err=%v", len(rows), err)
		}
		// non-NULL (должно быть 0)
		rows2, err := plQ.New().FilterByDistributorID(uuid.NullUUID{UUID: uuid.New(), Valid: true}).Select(ctx)
		if err != nil {
			t.Fatalf("FilterByDistributorID(non-null) err: %v", err)
		}
		if len(rows2) != 0 {
			t.Fatalf("expected 0 with non-matching distributor_id, got %d", len(rows2))
		}
	}

	// 7) Verified/Ownership/Status/TypeID/CityID
	{
		if n, err := plQ.New().FilterByVerified(true).Count(ctx); err != nil || n != 1 {
			t.Fatalf("FilterByVerified count: n=%d err=%v", n, err)
		}
		if n, err := plQ.New().FilterByOwnership("private").Count(ctx); err != nil || n != 1 {
			t.Fatalf("FilterByOwnership count: n=%d err=%v", n, err)
		}
		if n, err := plQ.New().FilterByStatus("active").Count(ctx); err != nil || n != 1 {
			t.Fatalf("FilterByStatus count: n=%d err=%v", n, err)
		}
		if n, err := plQ.New().FilterByTypeID(typ.ID).Count(ctx); err != nil || n != 1 {
			t.Fatalf("FilterByTypeID count: n=%d err=%v", n, err)
		}
		if n, err := plQ.New().FilterByCityID(pl.CityID).Count(ctx); err != nil || n != 1 {
			t.Fatalf("FilterByCityID count: n=%d err=%v", n, err)
		}
	}

	// 8) Update на place (обнулим website, поменяем статус)
	{
		u := dbx.UpdatePlaceParams{
			Status:    strPtr("inactive"),
			Website:   &sql.NullString{Valid: false},
			UpdatedAt: time.Now().UTC(),
		}
		if err := plQ.New().FilterByID(pl.ID).Update(ctx, u); err != nil {
			t.Fatalf("update place: %v", err)
		}
		got, err := plQ.New().FilterByID(pl.ID).Get(ctx)
		if err != nil {
			t.Fatalf("get after update: %v", err)
		}
		if got.Status != "inactive" || got.Website.Valid {
			t.Fatalf("update not applied: %+v", got)
		}
	}

	// 9) Delete place (каскадом должны уйти details и timetable)
	{
		if err := plQ.New().FilterByID(pl.ID).Delete(ctx); err != nil {
			t.Fatalf("delete place: %v", err)
		}
		// place
		if n, err := plQ.New().FilterByID(pl.ID).Count(ctx); err != nil || n != 0 {
			t.Fatalf("place still exists: n=%d err=%v", n, err)
		}
		// details
		if n, err := pdQ.New().FilterPlaceID(pl.ID).Count(ctx); err != nil || n != 0 {
			t.Fatalf("details still exist: n=%d err=%v", n, err)
		}
		// timetables
		if n, err := ptQ.New().FilterByPlaceID(pl.ID).Count(ctx); err != nil || n != 0 {
			t.Fatalf("timetables still exist: n=%d err=%v", n, err)
		}
	}
}
