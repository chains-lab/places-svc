package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain"
	"github.com/chains-lab/places-svc/internal/domain/models"
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
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Data.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Data.Code)

	ShopsClass := CreateClass(s, t, "Shops", "shops", nil)
	ElectronicsClass := CreateClass(s, t, "Electronics", "electronics", &ShopsClass.Data.Code)
	ClothesClass := CreateClass(s, t, "Clothes", "clothes", &ShopsClass.Data.Code)
	SHosesShopClass := CreateClass(s, t, "Shoes", "shoes", &ClothesClass.Data.Code)

	distributorFirstID := uuid.New()
	distributorSecondID := uuid.New()

	cityFirstID := uuid.New()
	citySecondID := uuid.New()

	_ = CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         RestaurantClass.Data.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        "en",
		Name:          "Restaurant Place",
		Address:       "123 Main St",
		Description:   "A nice restaurant place",
	})

	placeSuperMarket := CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        cityFirstID,
		DistributorID: &distributorSecondID,
		Class:         SuperMarketClass.Data.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        "en",
		Name:          "SuperMarket Place",
		Address:       "456 Market St",
		Description:   "A big supermarket place",
	})

	_ = CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        citySecondID,
		DistributorID: &distributorFirstID,
		Class:         SuperMarketClass.Data.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleUK,
		Name:          "SuperMarket Place Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket place in second city",
	})

	places, pag, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
		Classes: []string{
			FoodClass.Data.Code,
		},
		DistributorIDs: []uuid.UUID{
			distributorFirstID,
			distributorSecondID,
		},
	}, pagi.Request{}, []pagi.SortField{})
	if err != nil {
		t.Fatalf("ListPlaces by FoodClass: %v", err)
	}

	t.Logf("ListPlaces by FoodClass: got %d places", len(places))
	for _, place := range places {
		t.Logf("Place: %+v", place)
	}

	if len(places) != 3 {
		t.Fatalf("ListPlaces by FoodClass: expected 2 places, got %d", len(places))
	}
	if pag.Total != 3 {
		t.Fatalf("ListPlaces by FoodClass: expected total 2 places, got %d", pag.Total)
	}

	placeTxt := "Place"
	places, pag, err = s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
		Name: &placeTxt,
		Classes: []string{
			SuperMarketClass.Data.Code,
		},
	}, pagi.Request{}, []pagi.SortField{})
	if err != nil {
		t.Fatalf("ListPlaces by SuperMarketClass: %v", err)
	}

	t.Logf("ListPlaces by SuperMarketClass: got %d places", len(places))
	for _, place := range places {
		t.Logf("Place: %+v", place)
	}

	if len(places) != 2 {
		t.Fatalf("ListPlaces by SuperMarketClass: expected 2 places, got %d", len(places))
	}
	if pag.Total != 2 {
		t.Fatalf("ListPlaces by SuperMarketClass: expected total 2 places, got %d", pag.Total)
	}

	placeSuperMarket, err = s.app.GetPlace(ctx, placeSuperMarket.ID, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetPlace: %v", err)
	}
	if placeSuperMarket.Name != "SuperMarket Place" {
		t.Errorf("GetPlace: expected name 'SuperMarket Place', got '%s'", placeSuperMarket.Name)
	}

	places, pag, err = s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
		CityIDs: []uuid.UUID{cityFirstID},
	}, pagi.Request{}, []pagi.SortField{})
	if err != nil {
		t.Fatalf("ListPlaces by CityID: %v", err)
	}

	t.Logf("ListPlaces by CityID: got %d places", len(places))
	for _, place := range places {
		t.Logf("Place: %+v", place)
	}

	_ = CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        cityFirstID,
		DistributorID: nil,
		Class:         ShopsClass.Data.Code,
		Point:         [2]float64{30.4, 50.4},
		Locale:        "en",
		Name:          "Food Place",
		Address:       "303 Food St",
		Description:   "A delicious food place",
	})

	_ = CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         ElectronicsClass.Data.Code,
		Point:         [2]float64{30.2, 50.2},
		Locale:        "en",
		Name:          "Electronics Shop",
		Address:       "101 Electronics St",
		Description:   "A cool electronics shop",
	})

	_ = CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        cityFirstID,
		DistributorID: &distributorSecondID,
		Class:         ClothesClass.Data.Code,
		Point:         [2]float64{30.3, 50.3},
		Locale:        "ru",
		Name:          "Clothes Shop",
		Address:       "202 Clothes St",
		Description:   "A trendy clothes shop",
	})

	_ = CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        cityFirstID,
		DistributorID: nil,
		Class:         SHosesShopClass.Data.Code,
		Point:         [2]float64{30.5, 50.5},
		Locale:        "uk",
		Name:          "Clothes Shop UK",
		Address:       "204 Clothes St",
		Description:   "A trendy clothes shop UK",
	})

	places, pag, err = s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
		Classes: []string{
			ShopsClass.Data.Code,
		},
	}, pagi.Request{}, []pagi.SortField{})
	if err != nil {
		t.Fatalf("ListPlaces by CityID after adding shops: %v", err)
	}

	t.Logf("ListPlaces by codes after adding shops: got %d places", len(places))
	for _, place := range places {
		t.Logf("Place: %+v", place)
	}

	if pag.Total != 4 {
		t.Fatalf("ListPlaces by CityID after adding shops: expected total 4 places, got %d", pag.Total)
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
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Data.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Data.Code)

	// --- Идентификаторы
	distributorA := uuid.New()
	distributorB := uuid.New()
	city1 := uuid.New()
	city2 := uuid.New()

	// --- Тестовые точки (geography, метры)
	// center: (30.000, 50.000)
	// near:   (30.010, 50.000) ≈ 1.1 км восточнее
	// far:    (30.200, 50.000) ≈ 17–18 км
	pFoodCenter := CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        city1,
		DistributorID: &distributorA,
		Class:         FoodClass.Data.Code,
		Point:         [2]float64{30.000, 50.000},
		Locale:        "en",
		Name:          "Food Center",
		Address:       "123 Main St",
		Description:   "Center food",
	})
	pMarketNear := CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        city1,
		DistributorID: &distributorA,
		Class:         SuperMarketClass.Data.Code,
		Point:         [2]float64{30.010, 50.000},
		Locale:        "en",
		Name:          "SuperMarket Near",
		Address:       "456 Market St",
		Description:   "Near supermarket",
	})
	pRestFar := CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        city1,
		DistributorID: &distributorB,
		Class:         RestaurantClass.Data.Code,
		Point:         [2]float64{30.200, 50.000},
		Locale:        "en",
		Name:          "Restaurant Far",
		Address:       "789 Far St",
		Description:   "Far restaurant",
	})
	pMarketCity2 := CreatePlace(s, t, app.CreatePlaceParams{
		CityID:        city2,
		DistributorID: &distributorB,
		Class:         SuperMarketClass.Data.Code,
		Point:         [2]float64{31.100, 51.100},
		Locale:        "en",
		Name:          "SuperMarket City2",
		Address:       "111 Another St",
		Description:   "Other city",
	})

	if pRestFar.ID != pMarketCity2.ID {
	}

	// Сделаем один из них verified, чтобы проверить фильтр
	if _, err := s.app.VerifyPlace(ctx, pMarketNear.ID); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	// sanity-check: без фильтров должно быть 4
	all, pag, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{}, pagi.Request{}, nil)
	if err != nil {
		t.Fatalf("ListPlaces (no filters): %v", err)
	}
	if got := len(all); got != 4 || pag.Total != 4 {
		t.Fatalf("ListPlaces (no filters): len=%d total=%d; want 4/4", got, pag.Total)
	}

	t.Run("Filter by Class=Food (должно включать детей)", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Classes: []string{FoodClass.Data.Code}, // ждём Food + SuperMarket + Restaurant
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if pr.Total != 4 || len(places) != 4 {
			t.Fatalf("by class Food: got len=%d total=%d; want 4/4", len(places), pr.Total)
		}
	})

	t.Run("Filter by Class=SuperMarket", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Classes: []string{SuperMarketClass.Data.Code},
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if pr.Total != 2 || len(places) != 2 {
			t.Fatalf("by class SuperMarket: got len=%d total=%d; want 2/2", len(places), pr.Total)
		}
	})

	t.Run("Filter by CityIDs=city1", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			CityIDs: []uuid.UUID{city1},
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// в city1 три места
		if pr.Total != 3 || len(places) != 3 {
			t.Fatalf("by city1: got len=%d total=%d; want 3/3", len(places), pr.Total)
		}
	})

	t.Run("Filter by DistributorIDs", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			DistributorIDs: []uuid.UUID{distributorA},
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// у distributorA: pFoodCenter, pMarketNear
		if pr.Total != 2 || len(places) != 2 {
			t.Fatalf("by distributorA: got len=%d total=%d; want 2/2", len(places), pr.Total)
		}
	})

	t.Run("Filter by Verified=true", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Verified: boolPtr(true),
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// verified — только pMarketNear
		if pr.Total != 1 || len(places) != 1 || places[0].ID != pMarketNear.ID {
			t.Fatalf("by verified: got len=%d total=%d firstID=%v; want 1/1 %v",
				len(places), pr.Total, idOrZero(places), pMarketNear.ID)
		}
	})

	t.Run("Filter by Name LIKE 'Market'", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Name: strptr("Market"),
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// В названии 'Market' у pMarketNear и pMarketCity2
		if pr.Total != 2 || len(places) != 2 {
			t.Fatalf("by name like: got len=%d total=%d; want 2/2", len(places), pr.Total)
		}
	})

	t.Run("Filter by Address LIKE 'Main St'", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Address: strptr("Main St"),
		}, pagi.Request{}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// 'Main St' у pFoodCenter
		if pr.Total != 1 || len(places) != 1 || places[0].ID != pFoodCenter.ID {
			t.Fatalf("by address like: got len=%d total=%d firstID=%v; want 1/1 %v",
				len(places), pr.Total, idOrZero(places), pFoodCenter.ID)
		}
	})

	t.Run("Geo filter (radius ~2km from center) + sort by distance ASC", func(t *testing.T) {
		places, pr, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Location: &app.GeoFilterListPlaces{
				Point:   orb.Point{30.000, 50.000},
				RadiusM: 2000,
			},
		}, pagi.Request{}, []pagi.SortField{{Field: "distance", Ascend: true}})
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		// В радиусе 2км только pFoodCenter (0 м) и pMarketNear (~1.1 км)
		if pr.Total != 2 || len(places) != 2 {
			t.Fatalf("geo: got len=%d total=%d; want 2/2", len(places), pr.Total)
		}
		// порядок: сначала центр, потом near
		if places[0].ID != pFoodCenter.ID || places[1].ID != pMarketNear.ID {
			t.Fatalf("geo order: got [%v, %v]; want [%v, %v]",
				places[0].ID, places[1].ID, pFoodCenter.ID, pMarketNear.ID)
		}
	})

	t.Run("Pagination with class=Food, page size=2", func(t *testing.T) {
		// Всего 4 по классу Food (включая детей)
		page1, pr1, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Classes: []string{FoodClass.Data.Code},
		}, pagi.Request{Page: 1, Size: 2}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if len(page1) != 2 || pr1.Total != 4 {
			t.Fatalf("page1: len=%d total=%d; want 2/4", len(page1), pr1.Total)
		}

		page2, pr2, err := s.app.ListPlaces(ctx, enum.LocaleUK, app.FilterListPlaces{
			Classes: []string{FoodClass.Data.Code},
		}, pagi.Request{Page: 2, Size: 2}, nil)
		if err != nil {
			t.Fatalf("ListPlaces: %v", err)
		}
		if len(page2) != 2 || pr2.Total != 4 {
			t.Fatalf("page2: len=%d total=%d; want 2/4", len(page2), pr2.Total)
		}
	})
}

func boolPtr(b bool) *bool { return &b }

func idOrZero(pl []models.PlaceWithDetails) any {
	if len(pl) == 0 {
		return "none"
	}
	return pl[0].ID
}
