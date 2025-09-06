package dbx_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/google/uuid"
)

// --- helpers ---

func ensurePlace(t *testing.T, db *sql.DB) uuid.UUID {
	t.Helper()
	insertBaseKindInfra(t) // уже есть в твоих тестах
	pid := uuid.New()
	insertPlace(t, pid) // уже есть в твоих тестах
	return pid
}

func insertTT(t *testing.T, db *sql.DB, placeID uuid.UUID, s, e int) uuid.UUID {
	t.Helper()
	ctx := context.Background()
	id := uuid.New()
	err := dbx.NewPlaceTimetablesQ(db).Insert(ctx, dbx.PlaceTimetable{
		ID:       id,
		PlaceID:  placeID,
		StartMin: s,
		EndMin:   e,
	})
	if err != nil {
		t.Fatalf("insert timetable (%d,%d): %v", s, e, err)
	}
	return id
}

// --- tests ---

// 1) CRUD + EXCLUDE: пересекающиеся интервалы для одного place не должны вставляться.
func TestTimetables_CRUD_ExcludeOverlap(t *testing.T) {
	setupClean(t)
	db := openDB(t)
	ctx := context.Background()
	pid := ensurePlace(t, db)

	// вставляем 10:00-12:00 и 13:00-15:00
	_ = insertTT(t, db, pid, 10*60, 12*60)
	_ = insertTT(t, db, pid, 13*60, 15*60)

	// пробуем перекрывающийся 11:00-14:00 → ДОЛЖНО упасть на EXCLUDE
	err := dbx.NewPlaceTimetablesQ(db).Insert(ctx, dbx.PlaceTimetable{
		ID:       uuid.New(),
		PlaceID:  pid,
		StartMin: 11 * 60,
		EndMin:   14 * 60,
	})
	if err == nil {
		t.Fatalf("expected EXCLUDE violation on overlapping insert, got nil")
	}
	// для наглядности можно проверить код ошибки, но в разных драйверах сообщения отличаются
}

// 2) FilterBetween (прямой интервал): выбирает ЛЮБОЕ ПЕРЕСЕЧЕНИЕ, а не только полное вхождение.
func TestTimetables_FilterBetween_StraightOverlap(t *testing.T) {
	setupClean(t)
	db := openDB(t)
	ctx := context.Background()
	pid := ensurePlace(t, db)

	// интервалы:
	a := insertTT(t, db, pid, 600, 720)  // [10:00, 12:00)
	b := insertTT(t, db, pid, 721, 800)  // [11:40, 13:20) -- пересекает [11:00, 12:30)
	c := insertTT(t, db, pid, 900, 1020) // [15:00, 17:00) -- НЕ пересечёт

	// фильтруем окно [11:00, 12:30) => [660, 750)
	list, err := dbx.NewPlaceTimetablesQ(db).
		FilterPlaceID(pid).
		FilterBetween(660, 750).
		Select(ctx)
	if err != nil {
		t.Fatalf("select FilterBetween: %v", err)
	}

	gotIDs := map[uuid.UUID]bool{}
	for _, it := range list {
		gotIDs[it.ID] = true
	}
	// ожидаем a и b, но НЕ c
	if !gotIDs[a] || !gotIDs[b] || gotIDs[c] {
		t.Fatalf("expected a&b overlap only, got=%v", gotIDs)
	}
}

// 3) FilterBetween (перелом недели): s > e, например [10000, 200).
func TestTimetables_FilterBetween_WrapWeek(t *testing.T) {
	setupClean(t)
	ctx := context.Background()
	db := openDB(t)
	pid := ensurePlace(t, db)

	// Вставим 2 окна:
	x := insertTT(t, db, pid, 9950, 10070) // пересечёт хвост [10000, 200) (10000..10070)
	y := insertTT(t, db, pid, 30, 120)     // пересечёт начало [10000, 200) (0..120)

	// Плюс окно далеко
	_ = insertTT(t, db, pid, 500, 700)

	list, err := dbx.NewPlaceTimetablesQ(db).
		FilterPlaceID(pid).
		FilterBetween(10000, 200).
		Select(ctx)
	if err != nil {
		t.Fatalf("select FilterBetween wrap: %v", err)
	}

	got := map[uuid.UUID]bool{}
	for _, it := range list {
		got[it.ID] = true
	}
	if !got[x] || !got[y] || len(list) != 2 {
		t.Fatalf("expected exactly x&y for wrap, got=%v", got)
	}
}

// 4) Count и Page (limit/offset)
func TestTimetables_Count_And_Page(t *testing.T) {
	setupClean(t)
	db := openDB(t)
	ctx := context.Background()
	pid := ensurePlace(t, db)

	for i := 0; i < 5; i++ {
		insertTT(t, db, pid, 100*i+10, 100*i+50) // непересекающиеся
	}

	cnt, err := dbx.NewPlaceTimetablesQ(db).FilterPlaceID(pid).Count(ctx)
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if cnt != 5 {
		t.Fatalf("expected count=5, got=%d", cnt)
	}

	list, err := dbx.NewPlaceTimetablesQ(db).FilterPlaceID(pid).Page(0, 2).Select(ctx)
	if err != nil {
		t.Fatalf("page select: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected page of 2, got=%d", len(list))
	}
}

// 5) Get / FilterID / Delete
func TestTimetables_Get_FilterByID_Delete(t *testing.T) {
	setupClean(t)
	db := openDB(t)
	ctx := context.Background()
	pid := ensurePlace(t, db)

	id := insertTT(t, db, pid, 600, 700)

	got, err := dbx.NewPlaceTimetablesQ(db).FilterByID(id).Get(ctx)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.ID != id || got.PlaceID != pid || got.StartMin != 600 || got.EndMin != 700 {
		t.Fatalf("mismatch row: %+v", got)
	}

	if err := dbx.NewPlaceTimetablesQ(db).FilterByID(id).Delete(ctx); err != nil {
		t.Fatalf("delete by id: %v", err)
	}

	_, err = dbx.NewPlaceTimetablesQ(db).FilterByID(id).Get(ctx)
	if err == nil {
		t.Fatalf("expected no rows after delete")
	}
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}
