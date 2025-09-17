package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/google/uuid"
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
