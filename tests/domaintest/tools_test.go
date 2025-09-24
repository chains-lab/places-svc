package domaintest

import (
	"context"
	"testing"

	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
)

func CreateClass(s Setup, t *testing.T, name, code string, parent *string) models.ClassWithLocale {
	t.Helper()
	c, err := s.domain.class.Create(context.Background(), class.CreateParams{
		Name:   name,
		Code:   code,
		Icon:   "icon",
		Parent: parent,
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	c, err = s.domain.class.Activate(context.Background(), code, "en")
	if err != nil {
		t.Fatalf("ActivateClass: %v", err)
	}

	return c
}

func CreatePlace(s Setup, t *testing.T, params place.CreateParams) models.PlaceWithDetails {
	t.Helper()
	p, err := s.domain.place.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("CreatePlace: %v", err)
	}

	return p
}
