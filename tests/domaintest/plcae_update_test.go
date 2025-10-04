package domaintest

import (
	"context"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/google/uuid"
)

func TestPlaceUpdate(t *testing.T) {
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

	clothesNew, err := s.domain.place.Update(ctx, clothes.ID, clothes.Locale, place.UpdateParams{
		Class:   &ShoesShopClass.Code,
		Website: func(s string) *string { return &s }("https://new-website.com"),
		Phone:   func(s string) *string { return &s }("+1234567890"),
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	if clothesNew.Class != ShoesShopClass.Code {
		t.Errorf("expected updated place class to be %s, got %s", ShoesShopClass.Code, clothesNew.Class)
	}
	if clothesNew.Website == nil || *clothesNew.Website != "https://new-website.com" {
		t.Errorf("expected updated place website to be 'https://new-website.com', got %v", clothesNew.Website)
	}
	if clothesNew.Phone == nil || *clothesNew.Phone != "+1234567890" {
		t.Errorf("expected updated place phone to be '+1234567890', got %v", clothesNew.Phone)
	}

	restaurant, err = s.domain.place.Verify(ctx, restaurant.ID)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !restaurant.Verified {
		t.Errorf("expected place to be verified, got not verified")
	}

	restaurant, err = s.domain.place.UpdateStatus(ctx, restaurant.ID, enum.LocaleUK)
	if err != nil {
		t.Fatalf("ActivatePlace: %v", err)
	}
	if restaurant.Status != enum.PlaceStatusActive {
		t.Errorf("expected place status to be 'active', got '%s'", restaurant.Status)
	}

	restaurant, err = s.domain.place.Deactivate(ctx, restaurant.ID, enum.LocaleUK)
	if err != nil {
		t.Fatalf("DeactivatePlace: %v", err)
	}
	if restaurant.Status != enum.PlaceStatusInactive {
		t.Errorf("expected place status to be 'inactive', got '%s'", restaurant.Status)
	}
}
