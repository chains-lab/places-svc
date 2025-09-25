package domaintest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/google/uuid"
)

func TestCreatingPlace(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	classFirst, err := s.domain.class.Create(ctx, class.CreateParams{
		Name: "Classes First",
		Code: "class_first",
		Icon: "icon_1",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	classSecond, err := s.domain.class.Create(ctx, class.CreateParams{
		Name: "Classes Second",
		Code: "class_second",
		Icon: "icon_2",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	classFirst, err = s.domain.class.Activate(ctx, classFirst.Code)
	if err != nil {
		t.Fatalf("ActivateClass: %v", err)
	}

	distributorID := uuid.New()
	cityID := uuid.New()

	p, err := s.domain.place.Create(ctx, place.CreateParams{
		CityID:        cityID,
		DistributorID: &distributorID,
		Class:         classFirst.Code,
		Point:         [2]float64{30.0, 50.0},
		Locale:        "en",
		Name:          "place Name",
		Address:       "123 Main St",
		Description:   "A nice p",
	})

	if err != nil {
		t.Fatalf("CreatePlace: %v", err)
	}
	if p.Class != classFirst.Code {
		t.Fatalf("CreatePlace: expected classFirst %s, got %s", classFirst.Code, p.Class)
	}
	if p.CityID != cityID {
		t.Fatalf("CreatePlace: expected city ID %s, got %s", cityID, p.CityID)
	}

	placeID := p.ID

	p, err = s.domain.place.Update(ctx, placeID, p.Locale, place.UpdateParams{})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if p.Class != classFirst.Code {
		t.Fatalf("Update: expected classFirst %s, got %s", classFirst.Code, p.Class)
	}
	if p.CityID != cityID {
		t.Fatalf("Update: expected city ID %s, got \n %+v", cityID, p)
	}

	noneID := "none"
	p, err = s.domain.place.Update(ctx, placeID, p.Locale, place.UpdateParams{
		Class: &noneID,
	})
	if !errors.Is(err, errx.ErrorClassNotFound) {
		t.Fatalf("Update with none classFirst: expected error %v, got %v", errx.ErrorClassNotFound, err)
	}

	ws := "website"
	ph := "+1234567890"
	p, err = s.domain.place.Update(ctx, placeID, p.Locale, place.UpdateParams{
		Class:   &classSecond.Code,
		Website: &ws,
		Phone:   &ph,
	})
	if err != nil {
		t.Fatalf("Update with classSecond: %v", err)
	}
	if p.Class != classSecond.Code {
		t.Fatalf("Update with classSecond: expected classSecond %s, got %s", classSecond.Code, p.Class)
	}
	if p.Website == nil || *p.Website != ws {
		t.Fatalf("Update with classSecond: expected website %s, got %v", ws, p.Website)
	}
	if p.Phone == nil || *p.Phone != ph {
		t.Fatalf("Update with classSecond: expected phone %s, got %v", ph, p.Phone)
	}

	emtp := ""
	p, err = s.domain.place.Update(ctx, placeID, p.Locale, place.UpdateParams{
		Class:   &classSecond.Code,
		Website: &emtp,
		Phone:   &emtp,
	})
	if err != nil {
		t.Fatalf("Update with empty website and phone: %v", err)
	}
	if p.Website != nil {
		t.Fatalf("Update with empty website: expected nil, got %v", p.Website)
	}
	if p.Phone != nil {
		t.Fatalf("Update with empty phone: expected nil, got %v", p.Phone)
	}
}
