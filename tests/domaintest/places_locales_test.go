package domaintest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
	"github.com/google/uuid"
)

func TestPlaceLocales(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	FoodClass := CreateClass(s, t, "Food", "food", nil)
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Code)

	distributorFirstID := uuid.New()

	cityFirstID := uuid.New()
	citySecondID := uuid.New()

	food := CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         FoodClass.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        "en",
		Name:          "Food place",
		Address:       "456 Market St",
		Description:   "A big Food place",
	})

	restaurant := CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         RestaurantClass.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        "en",
		Name:          "Restaurant place",
		Address:       "123 Main St",
		Description:   "A nice restaurant place",
	})

	clothes := CreatePlace(s, t, place.CreateParams{
		CityID:        citySecondID,
		DistributorID: &distributorFirstID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleEN,
		Name:          "SuperMarket place Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket place in second city",
	})

	err = s.domain.place.SetLocales(ctx, food.ID, plocale.SetParams{
		Locale:      enum.LocaleUK,
		Name:        "Food place UK",
		Description: "A big Food place UK",
	}, plocale.SetParams{
		Locale:      enum.LocaleRU,
		Name:        "Food place RU",
		Description: "A big Food place RU",
	})
	if err != nil {
		t.Fatalf("SetPlaceLocales: %v", err)
	}

	err = s.domain.place.SetLocales(ctx, restaurant.ID, plocale.SetParams{
		Locale:      enum.LocaleUK,
		Name:        "Restaurant place UK",
		Description: "A nice restaurant place UK",
	})
	if err != nil {
		t.Fatalf("SetPlaceLocales: %v", err)
	}

	clothesEn, err := s.domain.place.Get(ctx, clothes.ID, enum.LocaleUK)
	if err != nil {
		t.Fatalf("GetPlace: %v", err)
	}
	if clothesEn.Locale != enum.LocaleEN {
		t.Fatalf("expected locale %s, got %s", enum.LocaleEN, clothesEn.Locale)
	}

	shops, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterParams{}, place.SortParams{})
	if err != nil {
		t.Fatalf("ListPlaces: %v", err)
	}
	if len(shops.Data) != 3 {
		t.Fatalf("expected 3 places, got %d", len(shops.Data))
	}
	if shops.Total != 3 {
		t.Fatalf("expected pag.Total 3, got %d", shops.Total)
	}

	t.Logf("ListPlaces: got %d places", len(shops.Data))

	for i, _ := range shops.Data {
		t.Logf("place: %+v", shops.Data[i])
		switch shops.Data[i].ID {
		case food.ID:
			if shops.Data[i].Locale != enum.LocaleUK {
				t.Fatalf("expected locale %s, got %s", enum.LocaleUK, shops.Data[i].Locale)
			}
		case restaurant.ID:
			if shops.Data[i].Locale != enum.LocaleUK {
				t.Fatalf("expected locale %s, got %s", enum.LocaleUK, shops.Data[i].Locale)
			}
		case clothes.ID:
			if shops.Data[i].Locale != enum.LocaleEN {
				t.Fatalf("expected locale %s, got %s", enum.LocaleEN, shops.Data[i].Locale)
			}
		}
	}
}
