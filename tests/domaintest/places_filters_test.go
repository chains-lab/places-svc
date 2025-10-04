package domaintest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func TestPlacesFilters(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	FoodClass := CreateClass(s, t, "Food", "food", nil)
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Code)

	ShopsClass := CreateClass(s, t, "Shops", "shops", nil)
	ElectronicsClass := CreateClass(s, t, "Electronics", "electronics", &ShopsClass.Code)
	ClothesClass := CreateClass(s, t, "Clothes", "clothes", &ShopsClass.Code)
	SHosesShopClass := CreateClass(s, t, "Shoes", "shoes", &ClothesClass.Code)

	distributorFirstID := uuid.New()
	distributorSecondID := uuid.New()

	cityFirstID := uuid.New()
	citySecondID := uuid.New()

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        "en",
		Name:          "Restaurant PlaceDetails",
		Address:       "123 Main St",
		Description:   "A nice restaurant p",
	})

	placeSuperMarket := CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorSecondID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        "en",
		Name:          "SuperMarket PlaceDetails",
		Address:       "456 Market St",
		Description:   "A big supermarket p",
	})

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        citySecondID,
		DistributorID: &distributorFirstID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleUK,
		Name:          "SuperMarket PlaceDetails Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket p in second city",
	})

	p, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
		Classes: []string{
			FoodClass.Code,
		},
		DistributorIDs: []uuid.UUID{
			distributorFirstID,
			distributorSecondID,
		},
	}, place.SortParams{})
	if err != nil {
		t.Fatalf("ListPlaces by FoodClass: %v", err)
	}

	t.Logf("ListPlaces by FoodClass: got %d p", len(p.Data))
	for _, p := range p.Data {
		t.Logf("PlaceDetails: %+v", p)
	}

	if len(p.Data) != 3 {
		t.Fatalf("ListPlaces by FoodClass: expected 2 p, got %d", len(p.Data))
	}
	if p.Total != 3 {
		t.Fatalf("ListPlaces by FoodClass: expected total 2 p, got %d", p.Total)
	}

	placeTxt := "PlaceDetails"
	p, err = s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
		Name: &placeTxt,
		Classes: []string{
			SuperMarketClass.Code,
		},
	}, place.SortParams{})
	if err != nil {
		t.Fatalf("ListPlaces by SuperMarketClass: %v", err)
	}

	t.Logf("ListPlaces by SuperMarketClass: got %d p", len(p.Data))
	for _, place := range p.Data {
		t.Logf("PlaceDetails: %+v", place)
	}

	if len(p.Data) != 2 {
		t.Fatalf("ListPlaces by SuperMarketClass: expected 2 p, got %d", len(p.Data))
	}
	if p.Total != 2 {
		t.Fatalf("ListPlaces by SuperMarketClass: expected total 2 p, got %d", p.Total)
	}

	placeSuperMarket, err = s.domain.place.Get(ctx, placeSuperMarket.ID, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetPlace: %v", err)
	}
	if placeSuperMarket.Name != "SuperMarket PlaceDetails" {
		t.Errorf("GetPlace: expected name 'SuperMarket PlaceDetails', got '%s'", placeSuperMarket.Name)
	}

	p, err = s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
		CityIDs: []uuid.UUID{cityFirstID},
	}, place.SortParams{})
	if err != nil {
		t.Fatalf("ListPlaces by CityID: %v", err)
	}

	t.Logf("ListPlaces by CityID: got %d p", len(p.Data))
	for _, place := range p.Data {
		t.Logf("PlaceDetails: %+v", place)
	}

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: nil,
		Class:         ShopsClass.Code,
		Point:         [2]float64{30.4, 50.4},
		Locale:        "en",
		Name:          "Food PlaceDetails",
		Address:       "303 Food St",
		Description:   "A delicious food p",
	})

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         ElectronicsClass.Code,
		Point:         [2]float64{30.2, 50.2},
		Locale:        "en",
		Name:          "Electronics Shop",
		Address:       "101 Electronics St",
		Description:   "A cool electronics shop",
	})

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorSecondID,
		Class:         ClothesClass.Code,
		Point:         [2]float64{30.3, 50.3},
		Locale:        "ru",
		Name:          "Clothes Shop",
		Address:       "202 Clothes St",
		Description:   "A trendy clothes shop",
	})

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: nil,
		Class:         SHosesShopClass.Code,
		Point:         [2]float64{30.5, 50.5},
		Locale:        "uk",
		Name:          "Clothes Shop UK",
		Address:       "204 Clothes St",
		Description:   "A trendy clothes shop UK",
	})

	p, err = s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
		Classes: []string{
			ShopsClass.Code,
		},
	}, place.SortParams{})
	if err != nil {
		t.Fatalf("ListPlaces by CityID after adding shops: %v", err)
	}

	t.Logf("ListPlaces by codes after adding shops: got %d p", len(p.Data))
	for _, place := range p.Data {
		t.Logf("PlaceDetails: %+v", place)
	}

	if p.Total != 4 {
		t.Fatalf("ListPlaces by CityID after adding shops: expected total 4 p, got %d", p.Total)
	}
}

func strptr(s string) *string { return &s }

func TestListPlaces_FiltersAndSorting(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	cleanDb(t)

	ctx := context.Background()

	// --- Классы: дерево Food -> (SuperMarket, Restaurant)
	FoodClass := CreateClass(s, t, "Food", "food", nil)
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Code)

	// --- Идентификаторы
	distributorA := uuid.New()
	distributorB := uuid.New()
	city1 := uuid.New()
	city2 := uuid.New()

	// --- Тестовые точки (geography, метры)
	// center: (30.000, 50.000)
	// near:   (30.010, 50.000) ≈ 1.1 км восточнее
	// far:    (30.200, 50.000) ≈ 17–18 км
	pFoodCenter := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorA,
		Class:         FoodClass.Code,
		Point:         [2]float64{30.000, 50.000},
		Locale:        "en",
		Name:          "Food Center",
		Address:       "123 Main St",
		Description:   "Center food",
	})
	pMarketNear := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorA,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{30.010, 50.000},
		Locale:        "en",
		Name:          "SuperMarket Near",
		Address:       "456 Market St",
		Description:   "Near supermarket",
	})
	pRestFar := CreatePlace(s, t, place.CreateParams{
		CityID:        city1,
		DistributorID: &distributorB,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.200, 50.000},
		Locale:        "en",
		Name:          "Restaurant Far",
		Address:       "789 Far St",
		Description:   "Far restaurant",
	})
	pMarketCity2 := CreatePlace(s, t, place.CreateParams{
		CityID:        city2,
		DistributorID: &distributorB,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{31.100, 51.100},
		Locale:        "en",
		Name:          "SuperMarket City2",
		Address:       "111 Another St",
		Description:   "Other city",
	})

	if pRestFar.ID != pMarketCity2.ID {
	}

	// Сделаем один из них verified, чтобы проверить фильтр
	if _, err := s.domain.place.Verify(ctx, pMarketNear.ID); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	// sanity-check: без фильтров должно быть 4
	all, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{}, place.SortParams{})
	if err != nil {
		t.Fatalf("ListPlaces (no filters): %v", err)
	}
	if got := len(all.Data); got != 4 || all.Total != 4 {
		t.Fatalf("ListPlaces (no filters): len=%d total=%d; want 4/4", got, all.Total)
	}

	t.Run("Filter by classData=Food (должно включать детей)", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Classes: []string{FoodClass.Code}, // ждём Food + SuperMarket + Restaurant
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if places.Total != 4 || len(places.Data) != 4 {
			t.Fatalf("by class Food: got len=%d total=%d; want 4/4", len(places.Data), places.Total)
		}
	})

	t.Run("Filter by classData=SuperMarket", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Classes: []string{SuperMarketClass.Code},
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if places.Total != 2 || len(places.Data) != 2 {
			t.Fatalf("by class SuperMarket: got len=%d total=%d; want 2/2", len(places.Data), places.Total)
		}
	})

	t.Run("Filter by CityIDs=city1", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			CityIDs: []uuid.UUID{city1},
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// в city1 три места
		if places.Total != 3 || len(places.Data) != 3 {
			t.Fatalf("by city1: got len=%d total=%d; want 3/3", len(places.Data), places.Total)
		}
	})

	t.Run("Filter by DistributorIDs", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			DistributorIDs: []uuid.UUID{distributorA},
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// у distributorA: pFoodCenter, pMarketNear
		if places.Total != 2 || len(places.Data) != 2 {
			t.Fatalf("by distributorA: got len=%d total=%d; want 2/2", len(places.Data), places.Total)
		}
	})

	t.Run("Filter by Verified=true", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Verified: boolPtr(true),
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// verified — только pMarketNear
		if places.Total != 1 || len(places.Data) != 1 || places.Data[0].ID != pMarketNear.ID {
			t.Fatalf("by verified: got len=%d total=%d firstID=%v; want 1/1 %v",
				len(places.Data), places.Total, idOrZero(places.Data), pMarketNear.ID)
		}
	})

	t.Run("Filter by Name LIKE 'Market'", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Name: strptr("Market"),
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// В названии 'Market' у pMarketNear и pMarketCity2
		if places.Total != 2 || len(places.Data) != 2 {
			t.Fatalf("by name like: got len=%d total=%d; want 2/2", len(places.Data), places.Total)
		}
	})

	t.Run("Filter by Address LIKE 'Main St'", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Address: strptr("Main St"),
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// 'Main St' у pFoodCenter
		if places.Total != 1 || len(places.Data) != 1 || places.Data[0].ID != pFoodCenter.ID {
			t.Fatalf("by address like: got len=%d total=%d firstID=%v; want 1/1 %v",
				len(places.Data), places.Total, idOrZero(places.Data), pFoodCenter.ID)
		}
	})

	t.Run("Geo filter (radius ~2km from center) + sort by distance ASC", func(t *testing.T) {
		places, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Location: &place.FilterDistance{
				Point:   orb.Point{30.000, 50.000},
				RadiusM: 2000,
			},
		}, place.SortParams{
			ByDistance: func(b bool) *bool { return &b }(true),
		})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// В радиусе 2км только pFoodCenter (0 м) и pMarketNear (~1.1 км)
		if places.Total != 2 || len(places.Data) != 2 {
			t.Fatalf("geo: got len=%d total=%d; want 2/2", len(places.Data), places.Total)
		}
		// порядок: сначала центр, потом near
		if places.Data[0].ID != pFoodCenter.ID || places.Data[1].ID != pMarketNear.ID {
			t.Fatalf("geo order: got [%v, %v]; want [%v, %v]",
				places.Data[0].ID, places.Data[1].ID, pFoodCenter.ID, pMarketNear.ID)
		}
	})

	t.Run("PagConvert with class=Food, page size=2", func(t *testing.T) {
		// Всего 4 по классу Food (включая детей)
		page1, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Classes: []string{FoodClass.Code},
			Page:    1, Size: 2,
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if len(page1.Data) != 2 || page1.Total != 4 {
			t.Fatalf("page1: len=%d total=%d; want 2/4", len(page1.Data), page1.Total)
		}

		page2, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{
			Classes: []string{FoodClass.Code},
			Page:    2, Size: 2,
		}, place.SortParams{})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if len(page2.Data) != 2 || page2.Total != 4 {
			t.Fatalf("page2: len=%d total=%d; want 2/4", len(page2.Data), page2.Total)
		}
	})
}

func boolPtr(b bool) *bool { return &b }

func idOrZero(pl []models.Place) any {
	if len(pl) == 0 {
		return "none"
	}
	return pl[0].ID
}
