package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/test"
	"github.com/google/uuid"
)

func TestPlaceTimetable(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	test.CleanDB(t)

	ctx := context.Background()

	// --- Классы (одно дерево для всего теста) ---
	FoodClass := CreateClass(s, t, "Food", "food", nil)
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Code)

	// --- Базовые сущности ---
	distributorID := uuid.New()
	city1 := uuid.New()
	city2 := uuid.New()

	// --- Места для сценариев Set/Get/Delete/Replace ---
	food := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         FoodClass.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        enum.LocaleEN,
		Name:          "Food PlaceRow",
		Address:       "456 Market St",
		Description:   "A big Food place",
	})
	restaurant := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        enum.LocaleEN,
		Name:          "Restaurant PlaceRow",
		Address:       "123 Main St",
		Description:   "A nice restaurant place",
	})
	clothes := CreatePlace(s, t, place.CreateParams{
		CityID:        city2,
		DistributorID: &distributorID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleEN,
		Name:          "SuperMarket PlaceRow Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket place in second city",
	})

	// --- Хелперы для расписаний ---
	moment := func(w time.Weekday, h, m int) models.Moment {
		return models.Moment{Weekday: w, Time: time.Duration(h)*time.Hour + time.Duration(m)*time.Minute}
	}
	ti := func(wFrom time.Weekday, hFrom, mFrom int, wTo time.Weekday, hTo, mTo int) models.TimeInterval {
		return models.TimeInterval{
			From: moment(wFrom, hFrom, mFrom),
			To:   moment(wTo, hTo, mTo),
		}
	}
	mins := func(mm models.Moment) int { return mm.ToNumberMinutes() }

	// === Часть A: Set/Get/Delete/Replace ===

	t.Run("set timetable for restaurant and verify", func(t *testing.T) {
		ttRestaurant := models.Timetable{
			Table: []models.TimeInterval{
				ti(time.Monday, 9, 0, time.Monday, 18, 0),
				ti(time.Wednesday, 12, 0, time.Wednesday, 14, 0),
			},
		}
		got, err := s.domain.timetable.SetForPlace(ctx, restaurant.ID, enum.LocaleUK, ttRestaurant)
		if err != nil {
			t.Fatalf("SetPlaceTimeTable(restaurant): %v", err)
		}
		if len(got.Timetable.Table) != 2 {
			t.Fatalf("restaurant timetable expected 2 intervals, got %d", len(got.Timetable.Table))
		}

		rFromDB, err := s.domain.place.Get(ctx, restaurant.ID, enum.LocaleEN)
		if err != nil {
			t.Fatalf("GetPlace(restaurant): %v", err)
		}
		if len(rFromDB.Timetable.Table) != 2 {
			t.Fatalf("GetPlace(restaurant) timetable expected 2 intervals, got %d", len(rFromDB.Timetable.Table))
		}

		exp0s, exp0e := ttRestaurant.Table[0].ToNumberMinutes()
		exp1s, exp1e := ttRestaurant.Table[1].ToNumberMinutes()
		got0s, got0e := rFromDB.Timetable.Table[0].ToNumberMinutes()
		got1s, got1e := rFromDB.Timetable.Table[1].ToNumberMinutes()

		if got0s != exp0s || got0e != exp0e {
			t.Fatalf("restaurant interval[0] mismatch: want (%d,%d), got (%d,%d)", exp0s, exp0e, got0s, got0e)
		}
		if got1s != exp1s || got1e != exp1e {
			t.Fatalf("restaurant interval[1] mismatch: want (%d,%d), got (%d,%d)", exp1s, exp1e, got1s, got1e)
		}
		if got0s != mins(moment(time.Monday, 9, 0)) || got0e != mins(moment(time.Monday, 18, 0)) {
			t.Fatalf("restaurant interval[0] wrong minutes (Mon 09-18)")
		}
	})

	t.Run("set timetable for food and verify", func(t *testing.T) {
		ttFood := models.Timetable{
			Table: []models.TimeInterval{
				ti(time.Tuesday, 10, 0, time.Tuesday, 12, 0),
			},
		}
		got, err := s.domain.timetable.SetForPlace(ctx, food.ID, enum.LocaleUK, ttFood)
		if err != nil {
			t.Fatalf("SetPlaceTimeTable(food): %v", err)
		}
		if len(got.Timetable.Table) != 1 {
			t.Fatalf("food timetable expected 1 interval, got %d", len(got.Timetable.Table))
		}

		fFromDB, err := s.domain.place.Get(ctx, food.ID, enum.LocaleEN)
		if err != nil {
			t.Fatalf("GetPlace(food): %v", err)
		}
		if len(fFromDB.Timetable.Table) != 1 {
			t.Fatalf("GetPlace(food) timetable expected 1 interval, got %d", len(fFromDB.Timetable.Table))
		}
		fs, fe := fFromDB.Timetable.Table[0].ToNumberMinutes()
		es, ee := ttFood.Table[0].ToNumberMinutes()
		if fs != es || fe != ee {
			t.Fatalf("food interval mismatch: want (%d,%d), got (%d,%d)", es, ee, fs, fe)
		}
	})

	t.Run("delete restaurant timetable and verify empty", func(t *testing.T) {
		if err := s.domain.timetable.DeleteForPlace(ctx, restaurant.ID); err != nil {
			t.Fatalf("DeleteForPlace(restaurant): %v", err)
		}
		rAfterDel, err := s.domain.place.Get(ctx, restaurant.ID, enum.LocaleEN)
		if err != nil {
			t.Fatalf("GetPlace(restaurant after delete): %v", err)
		}
		if len(rAfterDel.Timetable.Table) != 0 {
			t.Fatalf("restaurant timetable expected empty after delete, got %d", len(rAfterDel.Timetable.Table))
		}
	})

	t.Run("replace food timetable and verify order/content", func(t *testing.T) {
		ttFoodReplace := models.Timetable{
			Table: []models.TimeInterval{
				ti(time.Tuesday, 14, 0, time.Tuesday, 16, 0),
				ti(time.Thursday, 9, 0, time.Thursday, 10, 0),
			},
		}
		if _, err := s.domain.timetable.SetForPlace(ctx, food.ID, enum.LocaleEN, ttFoodReplace); err != nil {
			t.Fatalf("SetPlaceTimeTable(food replace): %v", err)
		}

		fAfterReplace, err := s.domain.place.Get(ctx, food.ID, enum.LocaleEN)
		if err != nil {
			t.Fatalf("GetPlace(food after replace): %v", err)
		}
		if len(fAfterReplace.Timetable.Table) != 2 {
			t.Fatalf("food timetable expected 2 intervals after replace, got %d", len(fAfterReplace.Timetable.Table))
		}

		s0, e0 := fAfterReplace.Timetable.Table[0].ToNumberMinutes()
		s1, e1 := fAfterReplace.Timetable.Table[1].ToNumberMinutes()
		rs0, re0 := ttFoodReplace.Table[0].ToNumberMinutes() // Tue 14–16
		rs1, re1 := ttFoodReplace.Table[1].ToNumberMinutes() // Thu 09–10
		if s0 != rs0 || e0 != re0 || s1 != rs1 || e1 != re1 {
			t.Fatalf("food intervals order/content mismatch after replace")
		}
	})

	t.Run("clothes has empty timetable", func(t *testing.T) {
		cFromDB, err := s.domain.place.Get(ctx, clothes.ID, enum.LocaleEN)
		if err != nil {
			t.Fatalf("GetPlace(clothes): %v", err)
		}
		if len(cFromDB.Timetable.Table) != 0 {
			t.Fatalf("clothes timetable expected empty, got %d", len(cFromDB.Timetable.Table))
		}
	})

	// === Часть B: фильтрация по окну времени (p1–p4) ===

	// отдельные места для фильтрации
	p1 := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        enum.LocaleEN,
		Name:          "P1 Restaurant",
		Address:       "Addr 1",
		Description:   "Desc 1",
	})
	p2 := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        enum.LocaleEN,
		Name:          "P2 Market",
		Address:       "Addr 2",
		Description:   "Desc 2",
	})
	p3 := CreatePlace(s, t, place.CreateParams{
		CityID:        city2,
		DistributorID: &distributorID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleEN,
		Name:          "P3 Restaurant Tue",
		Address:       "Addr 3",
		Description:   "Desc 3",
	})
	// p4 для «перелома недели» (вс—пн окно)
	p4 := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.05, 50.05},
		Locale:        enum.LocaleEN,
		Name:          "P4 Late Night",
		Address:       "Addr 4",
		Description:   "Desc 4",
	})

	// расписания под фильтрацию:
	// p1: Mon 09:00–13:00
	if _, err := s.domain.timetable.SetForPlace(ctx, p1.ID, enum.LocaleEN, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Monday, Time: 9 * time.Hour},
			To:   models.Moment{Weekday: time.Monday, Time: 13 * time.Hour},
		}},
	}); err != nil {
		t.Fatalf("SetPlaceTimeTable p1: %v", err)
	}

	// p2: Mon 15:00–20:00
	if _, err := s.domain.timetable.SetForPlace(ctx, p2.ID, enum.LocaleEN, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Monday, Time: 15 * time.Hour},
			To:   models.Moment{Weekday: time.Monday, Time: 20 * time.Hour},
		}},
	}); err != nil {
		t.Fatalf("SetPlaceTimeTable p2: %v", err)
	}

	// p3: Tue 09:00–12:00
	if _, err := s.domain.timetable.SetForPlace(ctx, p3.ID, enum.LocaleEN, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Tuesday, Time: 9 * time.Hour},
			To:   models.Moment{Weekday: time.Tuesday, Time: 12 * time.Hour},
		}},
	}); err != nil {
		t.Fatalf("SetPlaceTimeTable p3: %v", err)
	}

	// p4: Mon 00:15–01:00 (для окна Вс 23:30 → Пн 00:30)
	if _, err := s.domain.timetable.SetForPlace(ctx, p4.ID, enum.LocaleEN, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Monday, Time: 15 * time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 1 * time.Hour},
		}},
	}); err != nil {
		t.Fatalf("SetPlaceTimeTable p4: %v", err)
	}

	// helper для фильтра
	call := func(win models.TimeInterval) (models.PlacesCollection, int) {
		res, err := s.domain.place.Filter(
			ctx, enum.LocaleEN,
			place.FilterParams{Time: &win},
			place.SortParams{},
			0, 10,
		)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		return res, int(res.Total)
	}

	t.Run("Mon 10:30–10:31 -> only p1", func(t *testing.T) {
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Monday, Time: 10*time.Hour + 30*time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 10*time.Hour + 31*time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p1.ID {
			t.Fatalf("Mon 10:30: want only p1; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	})

	t.Run("Mon 16:00–16:01 -> only p2", func(t *testing.T) {
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Monday, Time: 16 * time.Hour},
			To:   models.Moment{Weekday: time.Monday, Time: 16*time.Hour + 1*time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p2.ID {
			t.Fatalf("Mon 16:00: want only p2; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	})

	t.Run("Tue 09:30–09:31 -> only p3", func(t *testing.T) {
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Tuesday, Time: 9*time.Hour + 30*time.Minute},
			To:   models.Moment{Weekday: time.Tuesday, Time: 9*time.Hour + 31*time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p3.ID {
			t.Fatalf("Tue 09:30: want only p3; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	})

	t.Run("Sun 23:30 -> Mon 00:30 (wrap) -> only p4", func(t *testing.T) {
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Sunday, Time: 23*time.Hour + 30*time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 30 * time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p4.ID {
			t.Fatalf("Sun->Mon wrap: want only p4; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	})
}

// idsOf — для удобных сообщений об ошибках
func idsOf(pp []models.Place) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(pp))
	for _, p := range pp {
		out = append(out, p.ID)
	}
	return out
}
