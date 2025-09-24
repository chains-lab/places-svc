package domaintest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
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
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Data.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Data.Code)

	distributorFirstID := uuid.New()

	cityFirstID := uuid.New()
	citySecondID := uuid.New()

	food := CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         FoodClass.Data.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        "en",
		Name:          "Food Place",
		Address:       "456 Market St",
		Description:   "A big Food place",
	})

	restaurant := CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorFirstID,
		Class:         RestaurantClass.Data.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        "en",
		Name:          "Restaurant Place",
		Address:       "123 Main St",
		Description:   "A nice restaurant place",
	})

	clothes := CreatePlace(s, t, place.CreateParams{
		CityID:        citySecondID,
		DistributorID: &distributorFirstID,
		Class:         SuperMarketClass.Data.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleEN,
		Name:          "SuperMarket Place Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket place in second city",
	})

	err = s.domain.place.SetLocales(ctx, food.ID, place.SetLocaleParams{
		Locale:      enum.LocaleUK,
		Name:        "Food Place UK",
		Description: "A big Food place UK",
	}, place.SetLocaleParams{
		Locale:      enum.LocaleRU,
		Name:        "Food Place RU",
		Description: "A big Food place RU",
	})
	if err != nil {
		t.Fatalf("SetPlaceLocales: %v", err)
	}

	err = s.domain.place.SetLocales(ctx, restaurant.ID, place.SetLocaleParams{
		Locale:      enum.LocaleUK,
		Name:        "Restaurant Place UK",
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

	shops, pag, err := s.domain.place.List(ctx, enum.LocaleUK, place.FilterListParams{}, pagi.Request{}, []pagi.SortField{})
	if err != nil {
		t.Fatalf("ListPlaces: %v", err)
	}
	if len(shops) != 3 {
		t.Fatalf("expected 3 places, got %d", len(shops))
	}
	if pag.Total != 3 {
		t.Fatalf("expected pag.Total 3, got %d", pag.Total)
	}

	t.Logf("ListPlaces: got %d places", len(shops))

	for i, _ := range shops {
		t.Logf("Place: %+v", shops[i])
		switch shops[i].ID {
		case food.ID:
			if shops[i].Locale != enum.LocaleUK {
				t.Fatalf("expected locale %s, got %s", enum.LocaleUK, shops[i].Locale)
			}
		case restaurant.ID:
			if shops[i].Locale != enum.LocaleUK {
				t.Fatalf("expected locale %s, got %s", enum.LocaleUK, shops[i].Locale)
			}
		case clothes.ID:
			if shops[i].Locale != enum.LocaleEN {
				t.Fatalf("expected locale %s, got %s", enum.LocaleEN, shops[i].Locale)
			}
		}
	}
}
