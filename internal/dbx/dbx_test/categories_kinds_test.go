package dbx_test

import (
	"context"
	"testing"
	"time"

	"github.com/chains-lab/places-svc/internal/dbx"
	_ "github.com/lib/pq"
)

func TestCategoryLocale_CRUD_AndFallback(t *testing.T) {
	setupClean(t)
	ctx := context.Background()
	db := openDB(t)

	insertBaseCategory(t, "food")

	cLocQ := dbx.NewCategoryLocaleQ(db)

	// Insert en
	if err := cLocQ.Insert(ctx, dbx.PlaceCategoryLocale{
		CategoryCode: "food",
		Locale:       "en",
		Name:         "Food",
	}); err != nil {
		t.Fatalf("insert en: %v", err)
	}

	// Insert uk
	if err := cLocQ.Insert(ctx, dbx.PlaceCategoryLocale{
		CategoryCode: "food",
		Locale:       "uk",
		Name:         "Їжа",
	}); err != nil {
		t.Fatalf("insert uk: %v", err)
	}

	// WithLocale uk → exact
	got, err := dbx.NewCategoryQ(db).WithLocale("uk").FilterCode("food").Get(ctx)
	if err != nil {
		t.Fatalf("get with uk: %v", err)
	}
	if got.Locale.Name != "Їжа" || got.Locale.Locale != "uk" {
		t.Fatalf("expected uk/Їжа, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}

	// WithLocale fr → fallback en
	got, err = dbx.NewCategoryQ(db).WithLocale("fr").FilterCode("food").Get(ctx)
	if err != nil {
		t.Fatalf("get with fr (fallback): %v", err)
	}
	if got.Locale.Name != "Food" || got.Locale.Locale != "en" {
		t.Fatalf("expected fallback en/Food, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}

	// Upsert uk → rename
	if err := cLocQ.Upsert(ctx, dbx.PlaceCategoryLocale{
		CategoryCode: "food",
		Locale:       "uk",
		Name:         "Їжа (оновлено)",
	}); err != nil {
		t.Fatalf("upsert uk rename: %v", err)
	}

	got, err = dbx.NewCategoryQ(db).WithLocale("uk").FilterCode("food").Get(ctx)
	if err != nil {
		t.Fatalf("get after upsert: %v", err)
	}
	if got.Locale.Name != "Їжа (оновлено)" {
		t.Fatalf("expected updated uk name, got %q", got.Locale.Name)
	}

	// Delete uk → должен остаться fallback en
	if err := cLocQ.New().FilterCategoryCode("food").FilterLocale("uk").Delete(ctx); err != nil {
		t.Fatalf("delete uk: %v", err)
	}
	got, err = dbx.NewCategoryQ(db).WithLocale("uk").FilterCode("food").Get(ctx)
	if err != nil {
		t.Fatalf("get after delete uk: %v", err)
	}
	if got.Locale.Locale != "en" || got.Locale.Name != "Food" {
		t.Fatalf("expected fallback en/Food after delete, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}

	// Delete en → теперь локали нет вовсе → пустой Locale{}
	if err := cLocQ.New().FilterCategoryCode("food").FilterLocale("en").Delete(ctx); err != nil {
		t.Fatalf("delete en: %v", err)
	}
	// Без WithLocale: Locale должен быть пустой (мы выбираем NULL AS loc_*)
	got, err = dbx.NewCategoryQ(db).FilterCode("food").Get(ctx)
	if err != nil {
		t.Fatalf("get without WithLocale: %v", err)
	}
	if got.Locale.Name != "" || got.Locale.Locale != "" {
		t.Fatalf("expected empty locale after all i18n deleted, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}
}

// ===== KINDS: i18n CRUD + fallback =====

func TestKindLocale_CRUD_AndFallback(t *testing.T) {
	setupClean(t)
	ctx := context.Background()
	db := openDB(t)

	insertBaseCategory(t, "food")
	insertBaseKind(t, "restaurant", "food")

	kLocQ := dbx.NewKindLocaleQ(db)

	// Insert en
	if err := kLocQ.Insert(ctx, dbx.PlaceKindLocale{
		KindCode: "restaurant",
		Locale:   "en",
		Name:     "Restaurant",
	}); err != nil {
		t.Fatalf("insert en: %v", err)
	}

	// Upsert uk
	if err := kLocQ.Upsert(ctx, dbx.PlaceKindLocale{
		KindCode: "restaurant",
		Locale:   "uk",
		Name:     "Ресторан",
	}); err != nil {
		t.Fatalf("upsert uk: %v", err)
	}

	// WithLocale uk → exact
	got, err := dbx.NewPlaceKindsQ(db).WithLocale("uk").FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get uk: %v", err)
	}
	if got.Locale.Locale != "uk" || got.Locale.Name != "Ресторан" {
		t.Fatalf("expected uk/Ресторан, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}

	// WithLocale fr → fallback en
	got, err = dbx.NewPlaceKindsQ(db).WithLocale("fr").FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get fr (fallback): %v", err)
	}
	if got.Locale.Locale != "en" || got.Locale.Name != "Restaurant" {
		t.Fatalf("expected fallback en/Restaurant, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}

	// Update via Update() (not upsert)
	newName := "Ресторан (оновлено)"
	if err := kLocQ.New().
		FilterKindCode("restaurant").
		FilterLocale("uk").
		Update(ctx, dbx.UpdateKindLocaleParams{Name: &newName}); err != nil {
		t.Fatalf("update uk via Update: %v", err)
	}

	got, err = dbx.NewPlaceKindsQ(db).WithLocale("uk").FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get after Update: %v", err)
	}
	if got.Locale.Name != newName {
		t.Fatalf("expected updated uk name, got %q", got.Locale.Name)
	}

	// Delete uk → fallback на en
	if err := kLocQ.New().FilterKindCode("restaurant").FilterLocale("uk").Delete(ctx); err != nil {
		t.Fatalf("delete uk: %v", err)
	}
	got, err = dbx.NewPlaceKindsQ(db).WithLocale("uk").FilterCode("restaurant").Get(ctx)
	if err != nil {
		t.Fatalf("get after delete uk: %v", err)
	}
	if got.Locale.Locale != "en" || got.Locale.Name != "Restaurant" {
		t.Fatalf("expected fallback en/Restaurant after delete, got %q/%q", got.Locale.Locale, got.Locale.Name)
	}
}

// ===== Проверка параметризации WithLocale (плейсхолдеры) =====

func TestWithLocale_SQL_IsParameterized(t *testing.T) {
	setupClean(t)
	db := openDB(t)
	ctx := context.Background()

	insertBaseCategory(t, "food")

	// накидаем локали
	if err := dbx.NewCategoryLocaleQ(db).Insert(ctx, dbx.PlaceCategoryLocale{
		CategoryCode: "food",
		Locale:       "en",
		Name:         "Food",
	}); err != nil {
		t.Fatalf("insert en: %v", err)
	}
	if err := dbx.NewCategoryLocaleQ(db).Insert(ctx, dbx.PlaceCategoryLocale{
		CategoryCode: "food",
		Locale:       "uk",
		Name:         "Їжа",
	}); err != nil {
		t.Fatalf("insert uk: %v", err)
	}

	// строим SQL отдельно (без выполнения) и убеждаемся, что локаль — в args
	q := dbx.NewCategoryQ(db).WithLocale("uk").FilterCode("food")
	sqlStr, args, err := q.SelectorToSql() // небольшой хелпер в пакете dbx (см. ниже)
	if err != nil {
		t.Fatalf("ToSql(): %v", err)
	}
	if len(args) == 0 {
		t.Fatalf("expected args, got none; sql: %s", sqlStr)
	}
	found := false
	for _, a := range args {
		if s, ok := a.(string); ok && s == "uk" {
			found = true
		}
	}
	if !found {
		t.Fatalf(`expected "uk" among args, got %#v`, args)
	}
}

// ===== Select без WithLocale (NULL AS loc_*) и с пагинацией =====

func TestCategory_Select_NoLocaleAndPaginate(t *testing.T) {
	setupClean(t)
	ctx := context.Background()
	db := openDB(t)
	now := time.Now().UTC()

	cq := dbx.NewCategoryQ(db)
	for _, code := range []string{"food", "drinks", "shops"} {
		if err := cq.Insert(ctx, dbx.PlaceCategory{
			Code:      code,
			Status:    "active",
			Icon:      "🧩",
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			t.Fatalf("insert %s: %v", code, err)
		}
	}

	// без WithLocale локальные поля должны быть пустыми
	items, err := cq.New().Page(2, 0).Select(ctx)
	if err != nil {
		t.Fatalf("select: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	for _, it := range items {
		if it.Locale.Name != "" || it.Locale.Locale != "" {
			t.Fatalf("expected empty locale without WithLocale, got %q/%q", it.Locale.Locale, it.Locale.Name)
		}
	}
}
