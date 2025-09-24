package domaintest

import (
	"context"
	"testing"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/google/uuid"
)

func TestPlaceTimetable(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	cleanDb(t)

	ctx := context.Background()

	FoodClass := CreateClass(s, t, "Food", "food", nil)
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Code)

	distributorID := uuid.New()
	city1 := uuid.New()
	city2 := uuid.New()

	food := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         FoodClass.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        enum.LocaleEN,
		Name:          "Food PlaceDetails",
		Address:       "456 Market St",
		Description:   "A big Food place",
	})
	restaurant := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        enum.LocaleEN,
		Name:          "Restaurant PlaceDetails",
		Address:       "123 Main St",
		Description:   "A nice restaurant place",
	})
	clothes := CreatePlace(s, t, place.CreateParams{
		CityID:        city2,
		DistributorID: &distributorID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleEN,
		Name:          "SuperMarket PlaceDetails Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket place in second city",
	})

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

	ttRestaurant := models.Timetable{
		Table: []models.TimeInterval{
			ti(time.Monday, 9, 0, time.Monday, 18, 0),
			ti(time.Wednesday, 12, 0, time.Wednesday, 14, 0),
		},
	}

	got, err := s.domain.place.SetTimetable(ctx, restaurant.ID, ttRestaurant)
	if err != nil {
		t.Fatalf("SetPlaceTimeTable(restaurant): %v", err)
	}
	if len(got.Timetable.Table) != 2 {
		t.Fatalf("restaurant timetable expected 2 intervals, got %d", len(got.Timetable.Table))
	}

	ttFood := models.Timetable{
		Table: []models.TimeInterval{
			ti(time.Tuesday, 10, 0, time.Tuesday, 12, 0),
		},
	}

	got, err = s.domain.place.SetTimetable(ctx, food.ID, ttFood)
	if err != nil {
		t.Fatalf("SetPlaceTimeTable(food): %v", err)
	}
	if len(got.Timetable.Table) != 1 {
		t.Fatalf("food timetable expected 1 interval, got %d", len(got.Timetable.Table))
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

	if err = s.domain.place.DeleteTimetable(ctx, restaurant.ID); err != nil {
		t.Fatalf("DeleteTimetable(restaurant): %v", err)
	}
	rAfterDel, err := s.domain.place.Get(ctx, restaurant.ID, enum.LocaleEN)
	if err != nil {
		t.Fatalf("GetPlace(restaurant after delete): %v", err)
	}
	if len(rAfterDel.Timetable.Table) != 0 {
		t.Fatalf("restaurant timetable expected empty after delete, got %d", len(rAfterDel.Timetable.Table))
	}

	ttFoodReplace := models.Timetable{
		Table: []models.TimeInterval{
			ti(time.Tuesday, 14, 0, time.Tuesday, 16, 0),
			ti(time.Thursday, 9, 0, time.Thursday, 10, 0),
		},
	}
	if _, err := s.domain.place.SetTimetable(ctx, food.ID, ttFoodReplace); err != nil {
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

	cFromDB, err := s.domain.place.Get(ctx, clothes.ID, enum.LocaleEN)
	if err != nil {
		t.Fatalf("GetPlace(clothes): %v", err)
	}
	if len(cFromDB.Timetable.Table) != 0 {
		t.Fatalf("clothes timetable expected empty, got %d", len(cFromDB.Timetable.Table))
	}
}
