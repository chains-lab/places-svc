package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

func TestCreatingPlace(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	classFirst, err := s.app.CreateClass(ctx, app.CreateClassParams{
		Name: "Classes First",
		Code: "class_first",
		Icon: "icon_1",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	classSecond, err := s.app.CreateClass(ctx, app.CreateClassParams{
		Name: "Classes Second",
		Code: "class_second",
		Icon: "icon_2",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	classFirst, err = s.app.ActivateClass(ctx, classFirst.Data.Code, enum.LocaleUK)
	if err != nil {
		t.Fatalf("ActivateClass: %v", err)
	}

	distributorID := uuid.New()
	cityID := uuid.New()

	place, err := s.app.CreatePlace(ctx, app.CreatePlaceParams{
		CityID:        cityID,
		DistributorID: &distributorID,
		Class:         classFirst.Data.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        "en",
		Name:          "Place Name",
		Address:       "123 Main St",
		Description:   "A nice place",
	})

	if err != nil {
		t.Fatalf("CreatePlace: %v", err)
	}
	if place.Class != classFirst.Data.Code {
		t.Fatalf("CreatePlace: expected classFirst %s, got %s", classFirst.Data.Code, place.Class)
	}
	if place.CityID != cityID {
		t.Fatalf("CreatePlace: expected city ID %s, got %s", cityID, place.CityID)
	}

	placeID := place.ID

	place, err = s.app.UpdatePlace(ctx, placeID, place.Locale, app.UpdatePlaceParams{})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if place.Class != classFirst.Data.Code {
		t.Fatalf("Update: expected classFirst %s, got %s", classFirst.Data.Code, place.Class)
	}
	if place.CityID != cityID {
		t.Fatalf("Update: expected city ID %s, got \n %+v", cityID, place)
	}

	noneID := "none"
	place, err = s.app.UpdatePlace(ctx, placeID, place.Locale, app.UpdatePlaceParams{
		Class: &noneID,
	})
	if !errors.Is(err, errx.ErrorClassNotFound) {
		t.Fatalf("Update with none classFirst: expected error %v, got %v", errx.ErrorClassNotFound, err)
	}

	ws := "website"
	ph := "+1234567890"
	place, err = s.app.UpdatePlace(ctx, placeID, place.Locale, app.UpdatePlaceParams{
		Class:   &classSecond.Data.Code,
		Website: &ws,
		Phone:   &ph,
	})
	if err != nil {
		t.Fatalf("Update with classSecond: %v", err)
	}
	if place.Class != classSecond.Data.Code {
		t.Fatalf("Update with classSecond: expected classSecond %s, got %s", classSecond.Data.Code, place.Class)
	}
	if place.Website == nil || *place.Website != ws {
		t.Fatalf("Update with classSecond: expected website %s, got %v", ws, place.Website)
	}
	if place.Phone == nil || *place.Phone != ph {
		t.Fatalf("Update with classSecond: expected phone %s, got %v", ph, place.Phone)
	}

	emtp := ""
	place, err = s.app.UpdatePlace(ctx, placeID, place.Locale, app.UpdatePlaceParams{
		Class:   &classSecond.Data.Code,
		Website: &emtp,
		Phone:   &emtp,
	})
	if err != nil {
		t.Fatalf("Update with empty website and phone: %v", err)
	}
	if place.Website != nil {
		t.Fatalf("Update with empty website: expected nil, got %v", place.Website)
	}
	if place.Phone != nil {
		t.Fatalf("Update with empty phone: expected nil, got %v", place.Phone)
	}
}
