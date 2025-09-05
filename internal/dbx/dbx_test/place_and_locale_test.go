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

func insertBaseKindInfra(t *testing.T) {
	t.Helper()
	db := openDB(t)
	ctx := context.Background()
	now := time.Now().UTC()

	// root class: food
	if err := dbx.NewClassesQ(db).Insert(ctx, dbx.InsertPlaceClass{
		Code:      "food",
		Father:    sql.NullString{Valid: false},
		Status:    "active",
		Icon:      "🍔",
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("insert class food: %v", err)
	}

	// child class: restaurant -> food
	parent := sql.NullString{String: "food", Valid: true}
	if err := dbx.NewClassesQ(db).Insert(ctx, dbx.InsertPlaceClass{
		Code:      "restaurant",
		Father:    parent,
		Status:    "active",
		Icon:      "🍽️",
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("insert class restaurant: %v", err)
	}

	// i18n (en fallback)
	if err := dbx.NewClassLocaleQ(db).Insert(ctx, dbx.PlaceClassLocale{
		Class: "restaurant", Locale: "en", Name: "Restaurant",
	}); err != nil {
		t.Fatalf("insert class_i18n restaurant en: %v", err)
	}
}

func insertPlace(t *testing.T, id uuid.UUID) {
	t.Helper()
	db := openDB(t)
	now := time.Now().UTC()
	err := dbx.NewPlacesQ(db).Insert(context.Background(), dbx.Place{
		ID:            id,
		CityID:        uuid.New(),
		DistributorID: uuid.NullUUID{}, // NULL
		Class:         "restaurant",
		Status:        "active",
		Verified:      true,
		Ownership:     "private",
		Point:         orb.Point{30.5234, 50.4501}, // Kyiv
		Website:       sql.NullString{Valid: true, String: "https://example.test"},
		Phone:         sql.NullString{}, // NULL
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	if err != nil {
		t.Fatalf("insert place: %v", err)
	}
}

func insertPlaceLocale(t *testing.T, placeID uuid.UUID, locale, name, addr string, desc sql.NullString) {
	t.Helper()
	db := openDB(t)
	err := dbx.NewPlaceLocalesQ(db).Insert(context.Background(), dbx.PlaceLocale{
		PlaceID:     placeID,
		Locale:      locale,
		Name:        name,
		Address:     addr,
		Description: desc,
	})
	if err != nil {
		t.Fatalf("insert place_i18n %s: %v", locale, err)
	}
}

func setupPlaceWithInfra(t *testing.T) (db *sql.DB, ctx context.Context, placeID uuid.UUID) {
	t.Helper()
	setupClean(t)
	insertBaseKindInfra(t)
	db = openDB(t)
	ctx = context.Background()
	placeID = uuid.New()
	insertPlace(t, placeID)
	return
}

// ---------- tests ----------

func TestPlaces_WithLocale_CRUD_Fallback(t *testing.T) {
	setupClean(t)
	insertBaseKindInfra(t)

	db := openDB(t)
	ctx := context.Background()

	placeID := uuid.New()
	insertPlace(t, placeID)

	// en базовая локаль (fallback)
	insertPlaceLocale(t, placeID, "en", "Coffee House", "Main St 1", sql.NullString{Valid: true, String: "Cozy place"})
	// uk локаль
	insertPlaceLocale(t, placeID, "uk", "Кавʼярня", "вул. Головна, 1", sql.NullString{Valid: false})

	// 1) exact: WithLocale("uk")
	got, err := dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with uk: %v", err)
	}
	if got.Locale == nil || *got.Locale != "uk" {
		t.Fatalf("expected locale uk, got: %#v", got.Locale)
	}
	if got.Name == nil || *got.Name != "Кавʼярня" {
		t.Fatalf("expected uk name, got: %#v", got.Name)
	}
	if got.Address == nil || *got.Address != "вул. Головна, 1" {
		t.Fatalf("expected uk address, got: %#v", got.Address)
	}
	if got.Description != nil && got.Description.Valid {
		t.Fatalf("expected uk description = NULL, got: %+v", *got.Description)
	}

	// 2) fallback: WithLocale("fr") → en
	got, err = dbx.NewPlacesQ(db).WithLocale("fr").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with fr (fallback): %v", err)
	}
	if got.Locale == nil || *got.Locale != "en" {
		t.Fatalf("expected fallback locale en, got: %#v", got.Locale)
	}
	if got.Name == nil || *got.Name != "Coffee House" {
		t.Fatalf("expected name Coffee House, got: %#v", got.Name)
	}
	if got.Address == nil || *got.Address != "Main St 1" {
		t.Fatalf("expected address Main St 1, got: %#v", got.Address)
	}
	if got.Description == nil || !got.Description.Valid || got.Description.String != "Cozy place" {
		t.Fatalf("expected description 'Cozy place', got: %#v", got.Description)
	}

	// 3) Update uk через PlaceLocalesQ.Update
	newName := "Кавʼярня (оновлено)"
	newAddr := "вул. Оновлена, 5"
	newDesc := sql.NullString{Valid: true, String: "Оновлений опис"}

	if err := dbx.NewPlaceLocalesQ(db).
		FilterPlaceID(placeID).
		FilterByLocale("uk").
		Update(ctx, dbx.UpdatePlaceLocaleParams{
			Name:        &newName,
			Address:     &newAddr,
			Description: &newDesc,
		}); err != nil {
		t.Fatalf("update uk: %v", err)
	}

	got, err = dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get after update uk: %v", err)
	}
	if got.Name == nil || *got.Name != newName {
		t.Fatalf("expected updated uk name, got: %#v", got.Name)
	}
	if got.Address == nil || *got.Address != newAddr {
		t.Fatalf("expected updated uk address, got: %#v", got.Address)
	}
	if got.Description == nil || !got.Description.Valid || got.Description.String != newDesc.String {
		t.Fatalf("expected updated uk description, got: %#v", got.Description)
	}

	// 4) Delete uk → fallback на en
	if err := dbx.NewPlaceLocalesQ(db).FilterPlaceID(placeID).FilterByLocale("uk").Delete(ctx); err != nil {
		t.Fatalf("delete uk: %v", err)
	}
	got, err = dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get after delete uk: %v", err)
	}
	if got.Locale == nil || *got.Locale != "en" || got.Name == nil || *got.Name != "Coffee House" {
		t.Fatalf("expected fallback en after delete uk, got: locale=%#v name=%#v", got.Locale, got.Name)
	}

	// 5) Delete en тоже → теперь нет ни точного, ни fallback → пустые локали
	if err := dbx.NewPlaceLocalesQ(db).FilterPlaceID(placeID).FilterByLocale("en").Delete(ctx); err != nil {
		t.Fatalf("delete en: %v", err)
	}
	got, err = dbx.NewPlacesQ(db).WithLocale("fr").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get after delete en: %v", err)
	}
	if got.Locale != nil || got.Name != nil || got.Address != nil || (got.Description != nil && got.Description.Valid) {
		t.Fatalf("expected empty locale after all i18n deleted, got: locale=%#v name=%#v addr=%#v desc=%#v",
			got.Locale, got.Name, got.Address, got.Description)
	}

	// 6) Без WithLocale вообще → локальные поля пустые (NULL AS loc_*)
	got, err = dbx.NewPlacesQ(db).FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get without WithLocale: %v", err)
	}
	if got.Locale != nil || got.Name != nil || got.Address != nil || (got.Description != nil && got.Description.Valid) {
		t.Fatalf("expected empty locale without WithLocale, got: locale=%#v name=%#v addr=%#v desc=%#v",
			got.Locale, got.Name, got.Address, got.Description)
	}
}

// Проверка, что WithLocale использует плейсхолдеры (locale в args, не в тексте SQL).
func TestPlaces_WithLocale_SQL_IsParameterized(t *testing.T) {
	setupClean(t)
	insertBaseKindInfra(t)

	db := openDB(t)
	placeID := uuid.New()
	insertPlace(t, placeID)

	// en локаль
	insertPlaceLocale(t, placeID, "en", "Coffee House", "Main St 1", sql.NullString{Valid: true, String: "Cozy"})

	// Требуется метод в dbx:
	// func (q PlacesQ) SelectorToSql() (string, []any, error) { return q.selector.ToSql() }
	sqlStr, args, err := dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).SelectorToSql()
	if err != nil {
		t.Fatalf("ToSql(): %v", err)
	}

	t.Logf("SQL: %s", sqlStr)
	t.Logf("Args: %#v", args)

	// как минимум два placeholder-а для locale (EXISTS и THEN), плюс id
	if len(args) < 3 {
		t.Fatalf("expected at least 3 args, got %#v", args)
	}
	hasUK := false
	for _, a := range args {
		if s, ok := a.(string); ok && s == "uk" {
			hasUK = true
			break
		}
	}
	if !hasUK {
		t.Fatalf(`expected "uk" among args, got %#v`, args)
	}
}

// 1) Невалидная локаль → sanitize до en. Проверяем, что WithLocale("xx-!!") вернёт EN.
func TestPlaces_InvalidLocale_SanitizedToEN(t *testing.T) {
	db, ctx, placeID := setupPlaceWithInfra(t)

	// только EN
	insertPlaceLocale(t, placeID, "en", "Coffee House", "Main St 1", sql.NullString{Valid: true, String: "Cozy place"})

	got, err := dbx.NewPlacesQ(db).WithLocale("xx-!!").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with invalid locale: %v", err)
	}
	if got.Locale == nil || *got.Locale != "en" ||
		got.Name == nil || *got.Name != "Coffee House" ||
		got.Address == nil || *got.Address != "Main St 1" ||
		got.Description == nil || !got.Description.Valid || got.Description.String != "Cozy place" {
		t.Fatalf("expected sanitized to en fallback, got: locale=%#v name=%#v addr=%#v desc=%#v",
			got.Locale, got.Name, got.Address, got.Description)
	}
}

//  2. Частично заполнённая локаль: uk есть, но description = NULL.
//     Должны получить uk.name/uk.address и NULL description (НЕ подтягивать en.description!)
func TestPlaces_PartialLocale_NoFieldMixing(t *testing.T) {
	db, ctx, placeID := setupPlaceWithInfra(t)

	// EN (полная)
	insertPlaceLocale(t, placeID, "en", "Coffee House", "Main St 1", sql.NullString{Valid: true, String: "Cozy place"})
	// UK (без description)
	insertPlaceLocale(t, placeID, "uk", "Кавʼярня", "вул. Головна, 1", sql.NullString{Valid: false})

	got, err := dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with uk: %v", err)
	}
	if got.Locale == nil || *got.Locale != "uk" {
		t.Fatalf("expected uk locale, got: %#v", got.Locale)
	}
	if got.Name == nil || *got.Name != "Кавʼярня" || got.Address == nil || *got.Address != "вул. Головна, 1" {
		t.Fatalf("expected uk name/address, got: name=%#v addr=%#v", got.Name, got.Address)
	}
	if got.Description != nil && got.Description.Valid {
		t.Fatalf("expected uk.description = NULL (no mixing from en), got: %#v", got.Description)
	}
}

//  3. Поиск по имени (FilterNameLike) c JOIN и DISTINCT: два плейса,
//     проверяем, что поиск по подстроке даёт корректный набор без дублей.
func TestPlaces_FilterNameLike_Distinct(t *testing.T) {
	db, ctx, p1 := setupPlaceWithInfra(t)
	p2 := uuid.New()
	insertPlace(t, p2)

	// p1: en "Coffee House"
	insertPlaceLocale(t, p1, "en", "Coffee House", "Main St 1", sql.NullString{})
	// p2: en "Coffee Corner"
	insertPlaceLocale(t, p2, "en", "Coffee Corner", "Second St 2", sql.NullString{})

	// найдём по "Coffee"
	list, err := dbx.NewPlacesQ(db).
		WithLocale("en").
		FilterNameLike("Coffee").
		Select(ctx)
	if err != nil {
		t.Fatalf("select with name like: %v", err)
	}
	// без дублей и оба попали
	found := map[uuid.UUID]bool{}
	for _, it := range list {
		found[it.ID] = true
	}
	if !found[p1] || !found[p2] || len(list) != 2 {
		t.Fatalf("expected both places (no duplicates). got=%v, ids=%v", len(list), found)
	}
}

// 4) Нет EN и нет exact → возвращаем пустую локаль (не фолбэчим на «какую-то» другую).
func TestPlaces_Fallback_NoEN_NoExact_ReturnsEmptyLocale(t *testing.T) {
	db, ctx, placeID := setupPlaceWithInfra(t)

	// есть только FR
	insertPlaceLocale(t, placeID, "fr", "Maison du Café", "Rue Principale 1", sql.NullString{})

	// запросим DE → ни exact (de), ни en нет → ожидаем пустые локализованные поля
	got, err := dbx.NewPlacesQ(db).WithLocale("de").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with de (no en, no exact): %v", err)
	}
	if got.Locale != nil || got.Name != nil || got.Address != nil || (got.Description != nil && got.Description.Valid) {
		t.Fatalf("expected empty locale (no en & no exact), got: locale=%#v name=%#v addr=%#v desc=%#v",
			got.Locale, got.Name, got.Address, got.Description)
	}
}

// 5) Есть только другая локаль: exact работает, а вот произвольный запрос без en — пусто.
func TestPlaces_OnlyOtherLocale_ExactVsEmpty(t *testing.T) {
	db, ctx, placeID := setupPlaceWithInfra(t)

	// есть только FR
	insertPlaceLocale(t, placeID, "fr", "Maison du Café", "Rue Principale 1", sql.NullString{Valid: true, String: "Sympa"})

	// exact: WithLocale("fr") → вернётся fr
	got, err := dbx.NewPlacesQ(db).WithLocale("fr").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with fr: %v", err)
	}
	if got.Locale == nil || *got.Locale != "fr" ||
		got.Name == nil || *got.Name != "Maison du Café" ||
		got.Address == nil || *got.Address != "Rue Principale 1" ||
		got.Description == nil || !got.Description.Valid || got.Description.String != "Sympa" {
		t.Fatalf("expected exact fr, got: locale=%#v name=%#v addr=%#v desc=%#v",
			got.Locale, got.Name, got.Address, got.Description)
	}

	// другая: WithLocale("uk") → ни exact (uk), ни en → пустые локали
	got, err = dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with uk (no en): %v", err)
	}
	if got.Locale != nil || got.Name != nil || got.Address != nil || (got.Description != nil && got.Description.Valid) {
		t.Fatalf("expected empty locale (no en & no exact), got: locale=%#v name=%#v addr=%#v desc=%#v",
			got.Locale, got.Name, got.Address, got.Description)
	}
}
