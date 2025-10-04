package domaintest

import (
	"context"
	"testing"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func TestFilterPlaceByTime(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	cleanDb(t)
	ctx := context.Background()

	// классы
	Food := CreateClass(s, t, "Food", "food", nil)
	Restaurant := CreateClass(s, t, "Restaurant", "restaurant", &Food.Code)
	SuperMarket := CreateClass(s, t, "SuperMarket", "supermarket", &Food.Code)

	// базовые сущности
	city1 := uuid.New()
	city2 := uuid.New()
	dist := uuid.New()

	// места
	p1 := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &dist,
		Class:         Restaurant.Code,
		Point:         orb.Point{30.0, 50.0},
		Locale:        enum.LocaleEN,
		Name:          "P1 Restaurant",
		Address:       "Addr 1",
		Description:   "Desc 1",
	})
	p2 := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &dist,
		Class:         SuperMarket.Code,
		Point:         orb.Point{30.1, 50.1},
		Locale:        enum.LocaleEN,
		Name:          "P2 Market",
		Address:       "Addr 2",
		Description:   "Desc 2",
	})
	p3 := CreatePlace(s, t, place.CreateParams{
		CityID:        city2,
		DistributorID: &dist,
		Class:         Restaurant.Code,
		Point:         orb.Point{31.1, 51.1},
		Locale:        enum.LocaleEN,
		Name:          "P3 Restaurant Tue",
		Address:       "Addr 3",
		Description:   "Desc 3",
	})
	// p4 для «перелома недели» (вс—пн окно)
	p4 := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &dist,
		Class:         Restaurant.Code,
		Point:         orb.Point{30.05, 50.05},
		Locale:        enum.LocaleEN,
		Name:          "P4 Late Night",
		Address:       "Addr 4",
		Description:   "Desc 4",
	})

	// расписания
	// p1: Mon 09:00–13:00
	_, err = s.domain.place.SetTimetable(ctx, p1.ID, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Monday, Time: 9 * time.Hour},
			To:   models.Moment{Weekday: time.Monday, Time: 13 * time.Hour},
		}},
	})
	if err != nil {
		t.Fatalf("SetPlaceTimeTable p1: %v", err)
	}

	// p2: Mon 15:00–20:00
	_, err = s.domain.place.SetTimetable(ctx, p2.ID, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Monday, Time: 15 * time.Hour},
			To:   models.Moment{Weekday: time.Monday, Time: 20 * time.Hour},
		}},
	})
	if err != nil {
		t.Fatalf("SetPlaceTimeTable p2: %v", err)
	}

	// p3: Tue 09:00–12:00
	_, err = s.domain.place.SetTimetable(ctx, p3.ID, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Tuesday, Time: 9 * time.Hour},
			To:   models.Moment{Weekday: time.Tuesday, Time: 12 * time.Hour},
		}},
	})
	if err != nil {
		t.Fatalf("SetPlaceTimeTable p3: %v", err)
	}

	// p4: Mon 00:15–01:00 (для теста с окном Вс 23:30 → Пн 00:30)
	_, err = s.domain.place.SetTimetable(ctx, p4.ID, models.Timetable{
		Table: []models.TimeInterval{{
			From: models.Moment{Weekday: time.Monday, Time: 15 * time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 1*time.Hour + 0*time.Minute},
		}},
	})
	if err != nil {
		t.Fatalf("SetPlaceTimeTable p4: %v", err)
	}

	// helper
	call := func(win models.TimeInterval) (models.PlacesCollection, int) {
		res, err := s.domain.place.List(
			ctx, enum.LocaleEN,
			place.FilterParams{Time: &win},
			place.SortParams{},
		)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		return res, int(res.Total)
	}

	// 1) Понедельник 10:30–10:31 → только p1
	{
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Monday, Time: 10*time.Hour + 30*time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 10*time.Hour + 31*time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p1.ID {
			t.Fatalf("Mon 10:30: want only p1; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	}

	// 2) Понедельник 16:00–16:01 → только p2
	{
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Monday, Time: 16*time.Hour + 0*time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 16*time.Hour + 1*time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p2.ID {
			t.Fatalf("Mon 16:00: want only p2; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	}

	// 3) Вторник 09:30–09:31 → только p3
	{
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Tuesday, Time: 9*time.Hour + 30*time.Minute},
			To:   models.Moment{Weekday: time.Tuesday, Time: 9*time.Hour + 31*time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p3.ID {
			t.Fatalf("Tue 09:30: want only p3; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	}

	// 4) Перелом недели: Вс 23:30 → Пн 00:30 → попадает p4 (Mon 00:15–01:00)
	{
		win := models.TimeInterval{
			From: models.Moment{Weekday: time.Sunday, Time: 23*time.Hour + 30*time.Minute},
			To:   models.Moment{Weekday: time.Monday, Time: 30 * time.Minute},
		}
		res, total := call(win)
		if total != 1 || len(res.Data) != 1 || res.Data[0].ID != p4.ID {
			t.Fatalf("Sun->Mon wrap: want only p4; got total=%d len=%d ids=%v", total, len(res.Data), idsOf(res.Data))
		}
	}
}

func idsOf(xs []models.Place) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(xs))
	for _, x := range xs {
		out = append(out, x.ID)
	}
	return out
}
