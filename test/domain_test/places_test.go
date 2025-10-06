package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
	"github.com/chains-lab/places-svc/test"
	"github.com/google/uuid"
)

func TestPlace(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}
	test.CleanDB(t)

	ctx := context.Background()

	FoodClass := CreateClass(s, t, "Food", "food", nil)
	SuperMarketClass := CreateClass(s, t, "SuperMarket", "supermarket", &FoodClass.Code)
	RestaurantClass := CreateClass(s, t, "Restaurant", "restaurant", &FoodClass.Code)

	ShopsClass := CreateClass(s, t, "Shops", "shops", nil)
	_ = CreateClass(s, t, "Electronics", "electronics", &ShopsClass.Code)
	ClothesClass := CreateClass(s, t, "Clothes", "clothes", &ShopsClass.Code)
	ShoesShopClass := CreateClass(s, t, "Shoes", "shoes", &ClothesClass.Code)

	distributorFirstID := uuid.New()
	distributorSecondID := uuid.New()
	cityFirstID := uuid.New()
	citySecondID := uuid.New()

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

	_ = CreatePlace(s, t, place.CreateParams{
		CityID:        cityFirstID,
		DistributorID: &distributorSecondID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{30.1, 50.1},
		Locale:        "en",
		Name:          "SuperMarket place",
		Address:       "456 Market St",
		Description:   "A big supermarket place",
	})

	clothes := CreatePlace(s, t, place.CreateParams{
		CityID:        citySecondID,
		DistributorID: &distributorFirstID,
		Class:         SuperMarketClass.Code,
		Point:         [2]float64{31.1, 51.1},
		Locale:        enum.LocaleUK,
		Name:          "SuperMarket place Second City",
		Address:       "789 Market St",
		Description:   "A big supermarket place in second city",
	})

	t.Run("Update_noop_keeps_values", func(t *testing.T) {
		got, err := s.domain.place.Update(ctx, restaurant.ID, restaurant.Locale, place.UpdateParams{})
		if err != nil {
			t.Fatalf("Update: %v", err)
		}
		if got.CityID != restaurant.CityID {
			t.Fatalf("Update: expected city ID %s, got %s", restaurant.CityID, got.CityID)
		}
		if got.Class != restaurant.Class {
			t.Fatalf("Update: expected class %s, got %s", restaurant.Class, got.Class)
		}
		restaurant = got
	})

	t.Run("Update_with_non_existing_class", func(t *testing.T) {
		noneID := "none"
		_, err := s.domain.place.Update(ctx, restaurant.ID, restaurant.Locale, place.UpdateParams{
			Class: &noneID,
		})
		if !errors.Is(err, errx.ErrorClassNotFound) {
			t.Fatalf("expected %v, got %v", errx.ErrorClassNotFound, err)
		}
	})

	t.Run("Update_with_valid_class_and_contacts", func(t *testing.T) {
		ws := "website"
		ph := "+1234567890"
		got, err := s.domain.place.Update(ctx, restaurant.ID, restaurant.Locale, place.UpdateParams{
			Class:   &SuperMarketClass.Code, // меняем класс на другой из того же дерева
			Website: &ws,
			Phone:   &ph,
		})
		if err != nil {
			t.Fatalf("Update with class+contacts: %v", err)
		}
		if got.Class != SuperMarketClass.Code {
			t.Fatalf("expected class %s, got %s", SuperMarketClass.Code, got.Class)
		}
		if got.Website == nil || *got.Website != ws {
			t.Fatalf("expected website %q, got %v", ws, got.Website)
		}
		if got.Phone == nil || *got.Phone != ph {
			t.Fatalf("expected phone %q, got %v", ph, got.Phone)
		}
		restaurant = got
	})

	t.Run("Update_clear_contacts_with_empty_strings", func(t *testing.T) {
		emtp := ""
		got, err := s.domain.place.Update(ctx, restaurant.ID, restaurant.Locale, place.UpdateParams{
			Class:   &restaurant.Class, // оставляем текущий класс
			Website: &emtp,
			Phone:   &emtp,
		})
		if err != nil {
			t.Fatalf("Update clear contacts: %v", err)
		}
		if got.Website != nil {
			t.Fatalf("expected website nil, got %v", *got.Website)
		}
		if got.Phone != nil {
			t.Fatalf("expected phone nil, got %v", *got.Phone)
		}
		restaurant = got
	})

	t.Run("Update_details_clothes_change_class_and_contacts", func(t *testing.T) {
		clothesNew, err := s.domain.place.Update(ctx, clothes.ID, clothes.Locale, place.UpdateParams{
			Class:   &ShoesShopClass.Code,
			Website: func(s string) *string { return &s }("https://new-website.com"),
			Phone:   func(s string) *string { return &s }("+1234567890"),
		})
		if err != nil {
			t.Fatalf("Update clothes: %v", err)
		}
		if clothesNew.Class != ShoesShopClass.Code {
			t.Errorf("expected class %s, got %s", ShoesShopClass.Code, clothesNew.Class)
		}
		if clothesNew.Website == nil || *clothesNew.Website != "https://new-website.com" {
			t.Errorf("expected website 'https://new-website.com', got %v", clothesNew.Website)
		}
		if clothesNew.Phone == nil || *clothesNew.Phone != "+1234567890" {
			t.Errorf("expected phone '+1234567890', got %v", clothesNew.Phone)
		}
	})

	t.Run("Verify_restaurant", func(t *testing.T) {
		got, err := s.domain.place.Verify(ctx, restaurant.ID, enum.LocaleUK, true)
		if err != nil {
			t.Fatalf("Verify: %v", err)
		}
		if !got.Verified {
			t.Errorf("expected verified=true")
		}
		restaurant = got
	})

	t.Run("Update_status_active", func(t *testing.T) {
		got, err := s.domain.place.UpdateStatus(ctx, restaurant.ID, enum.LocaleUK, enum.PlaceStatusActive)
		if err != nil {
			t.Fatalf("UpdateStatus active: %v", err)
		}
		if got.Status != enum.PlaceStatusActive {
			t.Errorf("expected status 'active', got %q", got.Status)
		}
		restaurant = got
	})

	t.Run("Update_status_inactive", func(t *testing.T) {
		got, err := s.domain.place.UpdateStatus(ctx, restaurant.ID, enum.LocaleUK, enum.PlaceStatusInactive)
		if err != nil {
			t.Fatalf("UpdateStatus inactive: %v", err)
		}
		if got.Status != enum.PlaceStatusInactive {
			t.Errorf("expected status 'inactive', got %q", got.Status)
		}
		restaurant = got
	})

}

func TestPlaceLocales(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDB(t)

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

	err = s.domain.plocale.SetForPlace(ctx, food.ID, plocale.SetParams{
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

	err = s.domain.plocale.SetForPlace(ctx, restaurant.ID, plocale.SetParams{
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

	shops, err := s.domain.place.Filter(ctx, enum.LocaleUK, place.FilterParams{}, place.SortParams{}, 1, 10)
	if err != nil {
		t.Fatalf("FilterPlaces: %v", err)
	}
	if len(shops.Data) != 3 {
		t.Fatalf("expected 3 places, got %d", len(shops.Data))
	}
	if shops.Total != 3 {
		t.Fatalf("expected pag.Total 3, got %d", shops.Total)
	}

	for i, _ := range shops.Data {
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
