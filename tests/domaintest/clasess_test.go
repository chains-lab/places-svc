package domaintest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
)

func TestCreatingClassAndDetails(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	c, err := s.domain.class.Create(ctx, class.CreateParams{
		Name: "Classes",
		Code: "class_first",
		Icon: "icon_1",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	err = s.domain.class.SetLocale(ctx, c.Code, class.SetLocaleParams{
		Locale: enum.LocaleEN,
		Name:   "Classes EN",
	},

		class.SetLocaleParams{
			Locale: enum.LocaleRU,
			Name:   "Classes RU",
		},

		class.SetLocaleParams{
			Locale: enum.LocaleUK,
			Name:   "Classes UK",
		})
	if err != nil {
		t.Fatalf("SetClassLocales: %v", err)
	}

	classEN, err := s.domain.class.Get(ctx, c.Code, enum.LocaleEN)
	if err != nil {
		t.Fatalf("GetClass EN: %v", err)
	}
	if classEN.Name != "Classes EN" {
		t.Fatalf("GetClass EN: expected name %s, got %s", "Classes EN", classEN.Name)
	}

	classRU, err := s.domain.class.Get(ctx, c.Code, enum.LocaleRU)
	if err != nil {
		t.Fatalf("GetClass RU: %v", err)
	}
	if classRU.Name != "Classes RU" {
		t.Fatalf("GetClass RU: expected name %s, got %s", "Classes RU", classRU.Name)
	}

	classChild, err := s.domain.class.Create(ctx, class.CreateParams{
		Name:   "Classes Child",
		Code:   "class_child",
		Icon:   "icon_child",
		Parent: &c.Code,
	})
	if err != nil {
		t.Fatalf("CreateClass Child: %v", err)
	}
	if *classChild.Parent != c.Code {
		t.Fatalf("CreateClass Child: expected parent %s, got %v", c.Code, classChild.Parent)
	}

	_, err = s.domain.class.Update(ctx, c.Code, c.Locale, class.UpdateParams{
		Parent: &c.Code,
	})
	if !errors.Is(err, errx.ErrorClassParentEqualCode) {
		t.Fatalf("UpdateClass: expected error %v, got %v", errx.ErrorClassParentEqualCode, err)
	}

	_, err = s.domain.class.Update(ctx, c.Code, c.Locale, class.UpdateParams{
		Parent: &classChild.Code,
	})
	if !errors.Is(err, errx.ErrorClassParentCycle) {
		t.Fatalf("UpdateClass: expected error %v, got %v", errx.ErrorClassParentCycle, err)
	}

	classParent, err := s.domain.class.Create(ctx, class.CreateParams{
		Name: "Classes Parent",
		Code: "class_parent",
		Icon: "icon_parent",
	})
	if err != nil {
		t.Fatalf("CreateClass Parent: %v", err)
	}

	c, err = s.domain.class.Update(ctx, c.Code, c.Locale, class.UpdateParams{
		Parent: &classParent.Code,
	})
	if err != nil {
		t.Fatalf("UpdateClass: %v", err)
	}
	if *c.Parent != classParent.Code {
		t.Fatalf("UpdateClass: expected parent %s, got %v", classParent.Code, c.Parent)
	}

	t.Logf("Parent: %v", classParent.Parent)
	t.Logf("Classes: %s", *c.Parent)
	t.Logf("Child: %s", *classChild.Parent)

	classes, err := s.domain.class.List(ctx, enum.LocaleUK, class.FilterListParams{
		Parent:      &classParent.Code,
		ParentCycle: true,
	})
	if err != nil {
		t.Fatalf("ListClasses: %v", err)
	}
	if len(classes.Data) != 3 {
		t.Fatalf("ListClasses: expected 3 classes, got %d", len(classes.Data))
	}
}
