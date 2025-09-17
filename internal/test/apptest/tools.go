package apptest

import (
	"context"
	"testing"

	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/app/models"
)

func CreateClass(s Setup, t *testing.T, name, code string, parent *string) models.ClassWithLocale {
	t.Helper()
	class, err := s.app.CreateClass(context.Background(), app.CreateClassParams{
		Name:   name,
		Code:   code,
		Icon:   "icon",
		Parent: parent,
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	class, err = s.app.ActivateClass(context.Background(), code, "en")
	if err != nil {
		t.Fatalf("ActivateClass: %v", err)
	}

	return class
}

func CreatePlace(s Setup, t *testing.T, params app.CreatePlaceParams) models.PlaceWithDetails {
	t.Helper()
	place, err := s.app.CreatePlace(context.Background(), params)
	if err != nil {
		t.Fatalf("CreatePlace: %v", err)
	}

	return place
}
