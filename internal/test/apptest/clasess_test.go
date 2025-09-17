package apptest

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/errx"
)

func TestCreatingClassAndDetails(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	class, err := s.app.CreateClass(ctx, app.CreateClassParams{
		Name: "Classes",
		Code: "class_first",
		Icon: "icon_1",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	err = s.app.SetClassLocales(ctx, class.Data.Code, app.SetClassLocaleParams{
		Locale: enum.LocaleEN,
		Name:   "Classes EN",
	},

		app.SetClassLocaleParams{
			Locale: enum.LocaleRU,
			Name:   "Classes RU",
		},

		app.SetClassLocaleParams{
			Locale: enum.LocaleUK,
			Name:   "Classes UK",
		})
	if err != nil {
		t.Fatalf("SetClassLocales: %v", err)
	}

	classEN, err := s.app.GetClass(ctx, class.Data.Code, enum.LocaleEN)
	if err != nil {
		t.Fatalf("GetClass EN: %v", err)
	}
	if classEN.Locale.Name != "Classes EN" {
		t.Fatalf("GetClass EN: expected name %s, got %s", "Classes EN", classEN.Locale.Name)
	}

	classRU, err := s.app.GetClass(ctx, class.Data.Code, enum.LocaleRU)
	if err != nil {
		t.Fatalf("GetClass RU: %v", err)
	}
	if classRU.Locale.Name != "Classes RU" {
		t.Fatalf("GetClass RU: expected name %s, got %s", "Classes RU", classRU.Locale.Name)
	}

	classChild, err := s.app.CreateClass(ctx, app.CreateClassParams{
		Name:   "Classes Child",
		Code:   "class_child",
		Icon:   "icon_child",
		Parent: &class.Data.Code,
	})
	if err != nil {
		t.Fatalf("CreateClass Child: %v", err)
	}
	if *classChild.Data.Parent != class.Data.Code {
		t.Fatalf("CreateClass Child: expected parent %s, got %v", class.Data.Code, classChild.Data.Parent)
	}

	_, err = s.app.UpdateClass(ctx, class.Data.Code, class.Locale.Locale, app.UpdateClassParams{
		Parent: &class.Data.Code,
	})
	if !errors.Is(err, errx.ErrorClassParentEqualCode) {
		t.Fatalf("UpdateClass: expected error %v, got %v", errx.ErrorClassParentEqualCode, err)
	}

	_, err = s.app.UpdateClass(ctx, class.Data.Code, class.Locale.Locale, app.UpdateClassParams{
		Parent: &classChild.Data.Code,
	})
	if !errors.Is(err, errx.ErrorClassParentCycle) {
		t.Fatalf("UpdateClass: expected error %v, got %v", errx.ErrorClassParentCycle, err)
	}

	classParent, err := s.app.CreateClass(ctx, app.CreateClassParams{
		Name: "Classes Parent",
		Code: "class_parent",
		Icon: "icon_parent",
	})
	if err != nil {
		t.Fatalf("CreateClass Parent: %v", err)
	}

	class, err = s.app.UpdateClass(ctx, class.Data.Code, class.Locale.Locale, app.UpdateClassParams{
		Parent: &classParent.Data.Code,
	})
	if err != nil {
		t.Fatalf("UpdateClass: %v", err)
	}
	if *class.Data.Parent != classParent.Data.Code {
		t.Fatalf("UpdateClass: expected parent %s, got %v", classParent.Data.Code, class.Data.Parent)
	}

	t.Logf("Parent: %v", classParent.Data.Parent)
	t.Logf("Classes: %s", *class.Data.Parent)
	t.Logf("Child: %s", *classChild.Data.Parent)

	classes, _, err := s.app.ListClasses(ctx, enum.LocaleUK, app.FilterListClassesParams{
		Parent:      &classParent.Data.Code,
		ParentCycle: true,
	}, pagi.Request{})
	if err != nil {
		t.Fatalf("ListClasses: %v", err)
	}
	if len(classes) != 3 {
		t.Fatalf("ListClasses: expected 3 classes, got %d", len(classes))
	}
}
