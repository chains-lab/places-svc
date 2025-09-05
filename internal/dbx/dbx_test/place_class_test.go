package dbx_test

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/chains-lab/places-svc/internal/dbx"
)

func TestPlaceClasses_Integration(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := openDB(t)

	// Чистим таблицы
	mustExec(t, db, "DELETE FROM place_class_i18n")
	mustExec(t, db, "DELETE FROM place_classes")

	q := dbx.NewClassesQ(db)

	// root: food
	root := dbx.InsertPlaceClass{
		Code:   "food",
		Status: "active",
		Icon:   "🍔",
	}
	if err := q.Insert(ctx, root); err != nil {
		t.Fatalf("insert root: %v", err)
	}
	// i18n для root
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)",
		"food", "en", "Food")
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)",
		"food", "uk", "Їжа")

	// child: restaurant -> food
	parent := sql.NullString{String: "food", Valid: true}
	child := dbx.InsertPlaceClass{
		Code:   "restaurant",
		Father: parent,
		Status: "active",
		Icon:   "🍽️",
	}
	if err := q.Insert(ctx, child); err != nil {
		t.Fatalf("insert child: %v", err)
	}
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)",
		"restaurant", "en", "Restaurant")
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)",
		"restaurant", "uk", "Ресторан")

	// grandchild: cafe -> restaurant
	parent2 := sql.NullString{String: "restaurant", Valid: true}
	grand := dbx.InsertPlaceClass{
		Code:   "cafe",
		Father: parent2,
		Status: "active",
		Icon:   "☕",
	}
	if err := q.Insert(ctx, grand); err != nil {
		t.Fatalf("insert grandchild: %v", err)
	}
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)",
		"cafe", "en", "Cafe")

	// WithLocale: uk
	got, err := q.New().WithLocale("uk").FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get with uk: %v", err)
	}
	if got.Name == nil || *got.Name != "Ресторан" {
		t.Errorf("want 'Ресторан', got %#v", got.Name)
	}
	if got.Locale == nil || *got.Locale != "uk" {
		t.Errorf("want locale 'uk', got %#v", got.Locale)
	}

	// WithLocale: fallback to en
	got, err = q.New().WithLocale("fr").FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get with fr fallback: %v", err)
	}
	if got.Name == nil || *got.Name != "Restaurant" {
		t.Errorf("want fallback 'Restaurant', got %#v", got.Name)
	}
	if got.Locale == nil || *got.Locale != "en" {
		t.Errorf("want fallback locale 'en', got %#v", got.Locale)
	}

	// Descendants of food (without food itself)
	desc, err := q.New().FilterFatherCycle("food").OrderBy("pc.code ASC").Select(ctx)
	if err != nil {
		t.Fatalf("select descendants: %v", err)
	}
	if len(desc) != 2 {
		t.Fatalf("want 2 descendants, got %d", len(desc))
	}
	if desc[0].Code != "cafe" || desc[1].Code != "restaurant" {
		t.Errorf("unexpected descendants: %v, %v", desc[0].Code, desc[1].Code)
	}

	// Count active
	activeCount, err := q.New().FilterStatus("active").Count(ctx)
	if err != nil {
		t.Fatalf("count active: %v", err)
	}
	if activeCount != 3 {
		t.Errorf("want active 3, got %d", activeCount)
	}

	// Pagination sanity
	page, err := q.New().OrderBy("pc.code ASC").Page(2, 0).Select(ctx)
	if err != nil {
		t.Fatalf("paginate: %v", err)
	}
	if len(page) != 2 {
		t.Errorf("want 2 on page, got %d", len(page))
	}

	// Cascade deprecate: food -> deprecated
	depr := "deprecated"
	if err := q.New().FilterCode("food").Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &depr,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("deprecate root: %v", err)
	}

	rs, err := q.New().FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get child after cascade: %v", err)
	}
	if rs.Status != "deprecated" {
		t.Errorf("child should be deprecated, got %s", rs.Status)
	}
	cf, err := q.New().FilterCode("cafe").Get(ctx)
	if err != nil {
		t.Fatalf("get grand after cascade: %v", err)
	}
	if cf.Status != "deprecated" {
		t.Errorf("grandchild should be deprecated, got %s", cf.Status)
	}

	// Forbid activation under deprecated ancestor
	act := "active"
	err = q.New().FilterCode("restaurant").Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &act,
		UpdatedAt: time.Now().UTC(),
	})
	if err == nil {
		t.Fatalf("expected error when activating under deprecated ancestor")
	}
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "cannot activate node") {
		t.Logf("activation blocked (driver msg): %v", err)
	}

	// Anti-cycle: try to move food under cafe
	newParent := "cafe"
	err = q.New().FilterCode("food").Update(ctx, dbx.UpdatePlaceClassParams{
		Father:    &newParent,
		UpdatedAt: time.Now().UTC(),
	})
	if err == nil {
		t.Fatalf("expected cycle error moving root under its descendant")
	}
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "cycle") {
		t.Logf("cycle blocked (driver msg): %v", err)
	}

	// Delete in order: cafe -> restaurant -> food
	if err := q.New().FilterCode("cafe").Delete(ctx); err != nil {
		t.Fatalf("delete cafe: %v", err)
	}
	if err := q.New().FilterCode("restaurant").Delete(ctx); err != nil {
		t.Fatalf("delete restaurant: %v", err)
	}
	if err := q.New().FilterCode("food").Delete(ctx); err != nil {
		t.Fatalf("delete food: %v", err)
	}

	// Should be empty
	finalCount, err := q.New().Count(ctx)
	if err != nil {
		t.Fatalf("count final: %v", err)
	}
	if finalCount != 0 {
		t.Errorf("want final 0, got %d", finalCount)
	}
}

func TestPlaceClasses_RepathAndRoots(t *testing.T) {
	t.Parallel()
	setupClean(t)

	ctx := context.Background()
	db := openDB(t)

	q := dbx.NewClassesQ(db)

	// roots: food, services
	food := dbx.InsertPlaceClass{Code: "food", Status: "active", Icon: "🍔"}
	if err := q.Insert(ctx, food); err != nil {
		t.Fatalf("insert food: %v", err)
	}
	services := dbx.InsertPlaceClass{Code: "services", Status: "active", Icon: "🧰"}
	if err := q.Insert(ctx, services); err != nil {
		t.Fatalf("insert services: %v", err)
	}

	// i18n
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)", "food", "en", "Food")
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)", "services", "en", "Services")

	// child: restaurant -> food
	parentFood := sql.NullString{String: "food", Valid: true}
	restaurant := dbx.InsertPlaceClass{Code: "restaurant", Father: parentFood, Status: "active", Icon: "🍽️"}
	if err := q.Insert(ctx, restaurant); err != nil {
		t.Fatalf("insert restaurant: %v", err)
	}
	mustExec(t, db, "INSERT INTO place_class_i18n(class, locale, name) VALUES ($1,$2,$3)", "restaurant", "en", "Restaurant")

	// grandchild: cafe -> restaurant
	parentRest := sql.NullString{String: "restaurant", Valid: true}
	cafe := dbx.InsertPlaceClass{Code: "cafe", Father: parentRest, Status: "active", Icon: "☕"}
	if err := q.Insert(ctx, cafe); err != nil {
		t.Fatalf("insert cafe: %v", err)
	}
	// uk для кафе через Upsert
	if err := dbx.NewClassLocaleQ(db).Upsert(ctx, dbx.PlaceClassLocale{
		Class:  "cafe",
		Locale: "uk",
		Name:   "Кав'ярня",
	}); err != nil {
		t.Fatalf("upsert cafe uk: %v", err)
	}

	// sanity paths
	gotFood, err := q.New().FilterCode("food").Get(ctx)
	if err != nil {
		t.Fatalf("get food: %v", err)
	}
	if gotFood.Path != "food" {
		t.Fatalf("want path food, got %q", gotFood.Path)
	}
	gotRest, err := q.New().FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get restaurant: %v", err)
	}
	if gotRest.Path != "food.restaurant" {
		t.Fatalf("want path food.restaurant, got %q", gotRest.Path)
	}
	gotCafe, err := q.New().FilterCode("cafe").Get(ctx)
	if err != nil {
		t.Fatalf("get cafe: %v", err)
	}
	if gotCafe.Path != "food.restaurant.cafe" {
		t.Fatalf("want path food.restaurant.cafe, got %q", gotCafe.Path)
	}

	// 1) reparent: restaurant -> services
	newParent := "services"
	if err := q.New().FilterCode("restaurant").Update(ctx, dbx.UpdatePlaceClassParams{
		Father:    &newParent,
		UpdatedAt: time.Now().UTC(),
	}); err != nil {
		t.Fatalf("reparent restaurant under services: %v", err)
	}

	// subtree paths updated
	gotRest, err = q.New().FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get restaurant after reparent: %v", err)
	}
	if gotRest.Path != "services.restaurant" {
		t.Fatalf("want path services.restaurant, got %q", gotRest.Path)
	}
	gotCafe, err = q.New().FilterCode("cafe").Get(ctx)
	if err != nil {
		t.Fatalf("get cafe after reparent: %v", err)
	}
	if gotCafe.Path != "services.restaurant.cafe" {
		t.Fatalf("want path services.restaurant.cafe, got %q", gotCafe.Path)
	}

	// 2) roots via FilterFather(nil)
	roots, err := q.New().FilterFather(sql.NullString{
		Valid: false,
	}).OrderBy("pc.code ASC").Select(ctx)
	if err != nil {
		t.Fatalf("select roots: %v", err)
	}
	if len(roots) != 2 {
		t.Fatalf("want 2 roots, got %d", len(roots))
	}
	if roots[0].Code != "food" || roots[1].Code != "services" {
		t.Fatalf("unexpected roots: %s, %s", roots[0].Code, roots[1].Code)
	}

	// 3) cannot delete parent with children (RESTRICT)
	err = q.New().FilterCode("services").Delete(ctx)
	if err == nil {
		t.Fatalf("expected FK restriction when deleting parent with children")
	}
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "foreign key") &&
		!strings.Contains(strings.ToLower(err.Error()), "restrict") {
		t.Logf("delete parent failed as expected (driver msg): %v", err)
	}

	// 4) WithLocale('uk') for cafe
	locCafe, err := q.New().WithLocale("uk").FilterCode("cafe").Get(ctx)
	if err != nil {
		t.Fatalf("get cafe with uk locale: %v", err)
	}
	if locCafe.Name == nil || *locCafe.Name != "Кав'ярня" {
		t.Fatalf("want uk name `Кав'ярня`, got %#v", locCafe.Name)
	}
	if locCafe.Locale == nil || *locCafe.Locale != "uk" {
		t.Fatalf("want locale `uk`, got %#v", locCafe.Locale)
	}
}
