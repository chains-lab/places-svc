package dbxtest

import (
	"database/sql"
	"testing"
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

func setupClean(t *testing.T) {
	t.Helper()
	db := openDB(t)
	// порядок важен из-за FK
	mustExec(t, db, "DELETE FROM place_timetables")
	mustExec(t, db, "DELETE FROM place_i18n")
	mustExec(t, db, "DELETE FROM places")
	mustExec(t, db, "DELETE FROM place_class_i18n")
	mustExec(t, db, "DELETE FROM place_classes")
}

//func insertBaseCategory(t *testing.T, code string) {
//	t.Helper()
//	db := openDB(t)
//	now := time.Now().UTC()
//	err := dbx.NewCategoryQ(db).Insert(context.Background(), dbx.PlaceCategory{
//		Code:      code,
//		Statuses:    "active",
//		Icon:      "🧩",
//		CreatedAt: now,
//		UpdatedAt: now,
//	})
//	if err != nil {
//		t.Fatalf("insertBaseCategory(%s): %v", code, err)
//	}
//}
//
//func insertBaseKind(t *testing.T, code, catCode string) {
//	t.Helper()
//	db := openDB(t)
//	now := time.Now().UTC()
//	err := dbx.NewPlaceKindsQ(db).Insert(context.Background(), dbx.PlaceKind{
//		Code:         code,
//		CategoryCode: catCode,
//		Statuses:       "active",
//		Icon:         "🏷️",
//		CreatedAt:    now,
//		UpdatedAt:    now,
//	})
//	if err != nil {
//		t.Fatalf("insertBaseKind(%s): %v", code, err)
//	}
//}
