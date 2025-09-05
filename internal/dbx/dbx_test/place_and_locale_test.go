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
	now := time.Now().UTC()

	// category
	if err := dbx.NewCategoryQ(db).Insert(context.Background(), dbx.PlaceCategory{
		Code:      "food",
		Status:    "active",
		Icon:      "🍔",
		CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("insert category: %v", err)
	}
	// kind
	if err := dbx.NewPlaceKindsQ(db).Insert(context.Background(), dbx.PlaceKind{
		Code:         "restaurant",
		CategoryCode: "food",
		Status:       "active",
		Icon:         "🍽️",
		CreatedAt:    now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("insert kind: %v", err)
	}
	// kind i18n (en fallback)
	if err := dbx.NewKindLocaleQ(db).Insert(context.Background(), dbx.PlaceKindLocale{
		KindCode: "restaurant", Locale: "en", Name: "Restaurant",
	}); err != nil {
		t.Fatalf("insert kind en: %v", err)
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
		TypeCode:      "restaurant",
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
	err := dbx.NewPlaceDetailsQ(db).Insert(context.Background(), dbx.PlaceLocale{
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
	setupClean(t)          // из предыдущих тестов
	insertBaseKindInfra(t) // создаёт category/kind: category "infra", kind "power_station" (пример)
	db = openDB(t)
	ctx = context.Background()
	placeID = uuid.New()
	insertPlace(t, placeID) // вставляет базовый place с type_code на существующий kind
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
	if got.Locale.Locale != "uk" || got.Locale.Name != "Кавʼярня" || got.Locale.Address != "вул. Головна, 1" || got.Locale.Description.Valid {
		t.Fatalf("expected uk locale with name/address and empty description, got: %+v", got.Locale)
	}

	// 2) fallback: WithLocale("fr") → en
	got, err = dbx.NewPlacesQ(db).WithLocale("fr").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with fr (fallback): %v", err)
	}
	if got.Locale.Locale != "en" || got.Locale.Name != "Coffee House" || got.Locale.Address != "Main St 1" || !got.Locale.Description.Valid || got.Locale.Description.String != "Cozy place" {
		t.Fatalf("expected fallback to en/\"Coffee House\", got: %+v", got.Locale)
	}

	// 3) Update uk через PlaceLocalesQ.Update
	newName := "Кавʼярня (оновлено)"
	newAddr := "вул. Оновлена, 5"
	newDesc := sql.NullString{Valid: true, String: "Оновлений опис"}

	if err := dbx.NewPlaceDetailsQ(db).
		FilterPlaceID(placeID).
		FilterByLocale("uk").
		Update(ctx, dbx.UpdatePlaceLocaleParams{
			Name:        &newName,
			Address:     &newAddr,
			Description: &newDesc,
			UpdatedAt:   time.Now().UTC(),
		}); err != nil {
		t.Fatalf("update uk: %v", err)
	}

	got, err = dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get after update uk: %v", err)
	}
	if got.Locale.Name != newName || got.Locale.Address != newAddr || !got.Locale.Description.Valid || got.Locale.Description.String != newDesc.String {
		t.Fatalf("expected updated uk locale, got: %+v", got.Locale)
	}

	// 4) Delete uk → fallback на en
	if err := dbx.NewPlaceDetailsQ(db).FilterPlaceID(placeID).FilterByLocale("uk").Delete(ctx); err != nil {
		t.Fatalf("delete uk: %v", err)
	}
	got, err = dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get after delete uk: %v", err)
	}
	if got.Locale.Locale != "en" || got.Locale.Name != "Coffee House" {
		t.Fatalf("expected fallback en after delete uk, got: %+v", got.Locale)
	}

	// 5) Delete en тоже → теперь нет ни точного, ни fallback → пустые локали
	if err := dbx.NewPlaceDetailsQ(db).FilterPlaceID(placeID).FilterByLocale("en").Delete(ctx); err != nil {
		t.Fatalf("delete en: %v", err)
	}
	got, err = dbx.NewPlacesQ(db).WithLocale("fr").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get after delete en: %v", err)
	}
	if got.Locale.Locale != "" || got.Locale.Name != "" || got.Locale.Address != "" || got.Locale.Description.Valid {
		t.Fatalf("expected empty locale after all i18n deleted, got: %+v", got.Locale)
	}

	// 6) Без WithLocale вообще → локальные поля пустые (заглушки NULL AS loc_*)
	got, err = dbx.NewPlacesQ(db).FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get without WithLocale: %v", err)
	}
	if got.Locale.Locale != "" || got.Locale.Name != "" || got.Locale.Address != "" || got.Locale.Description.Valid {
		t.Fatalf("expected empty locale without WithLocale, got: %+v", got.Locale)
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

	// ХЕЛПЕР: добавь в пакет dbx метод:
	// func (q PlacesQ) SelectorToSql() (string, []any, error) { return q.selector.ToSql() }
	sqlStr, args, err := dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).SelectorToSql()
	if err != nil {
		t.Fatalf("ToSql(): %v", err)
	}

	t.Logf("SQL: %s", sqlStr)
	t.Logf("Args: %#v", args)

	if len(args) < 3 {
		t.Fatalf("expected at least 3 args, got %#v", args)
	}
	hasUK := false
	for _, a := range args {
		if s, ok := a.(string); ok && s == "uk" {
			hasUK = true
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
	if got.Locale.Locale != "en" || got.Locale.Name != "Coffee House" || got.Locale.Address != "Main St 1" || !got.Locale.Description.Valid || got.Locale.Description.String != "Cozy place" {
		t.Fatalf("expected sanitized to en fallback, got: %+v", got.Locale)
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
	if got.Locale.Locale != "uk" {
		t.Fatalf("expected uk locale, got: %q", got.Locale.Locale)
	}
	if got.Locale.Name != "Кавʼярня" || got.Locale.Address != "вул. Головна, 1" {
		t.Fatalf("expected uk name/address, got: %+v", got.Locale)
	}
	if got.Locale.Description.Valid {
		t.Fatalf("expected uk.description = NULL (no mixing from en), got: %+v", got.Locale.Description)
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
	if got.Locale.Locale != "" || got.Locale.Name != "" || got.Locale.Address != "" || got.Locale.Description.Valid {
		t.Fatalf("expected empty locale (no en & no exact), got: %+v", got.Locale)
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
	if got.Locale.Locale != "fr" || got.Locale.Name != "Maison du Café" || got.Locale.Address != "Rue Principale 1" || !got.Locale.Description.Valid || got.Locale.Description.String != "Sympa" {
		t.Fatalf("expected exact fr, got: %+v", got.Locale)
	}

	// другая: WithLocale("uk") → ни exact (uk), ни en → пустые локали
	got, err = dbx.NewPlacesQ(db).WithLocale("uk").FilterByID(placeID).Get(ctx)
	if err != nil {
		t.Fatalf("get with uk (no en): %v", err)
	}
	if got.Locale.Locale != "" || got.Locale.Name != "" || got.Locale.Address != "" || got.Locale.Description.Valid {
		t.Fatalf("expected empty locale (no en & no exact), got: %+v", got.Locale)
	}
}
