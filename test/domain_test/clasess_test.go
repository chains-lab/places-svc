package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/test"
)

func TestClasses(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	test.CleanDB(t)

	ctx := context.Background()

	firstClass, err := s.domain.class.Create(ctx, class.CreateParams{
		Name: "Classes",
		Code: "class_first",
		Icon: "icon_1",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	classChild, err := s.domain.class.Create(ctx, class.CreateParams{
		Name:   "Classes Child",
		Code:   "class_child",
		Icon:   "icon_child",
		Parent: &firstClass.Code,
	})
	if err != nil {
		t.Fatalf("CreateClass Child: %v", err)
	}
	if *classChild.Parent != firstClass.Code {
		t.Fatalf("CreateClass Child: expected parent %s, got %v", firstClass.Code, classChild.Parent)
	}

	_, err = s.domain.class.Update(ctx, firstClass.Code, class.UpdateParams{
		Parent: &firstClass.Code,
	})
	if !errors.Is(err, errx.ErrorClassParentCycle) {
		t.Fatalf("UpdateClass: expected error %v, got %v", errx.ErrorClassParentCycle, err)
	}

	_, err = s.domain.class.Update(ctx, firstClass.Code, class.UpdateParams{
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

	firstClass, err = s.domain.class.Update(ctx, firstClass.Code, class.UpdateParams{
		Parent: &classParent.Code,
	})
	if err != nil {
		t.Fatalf("UpdateClass: %v", err)
	}
	if *firstClass.Parent != classParent.Code {
		t.Fatalf("UpdateClass: expected parent %s, got %v", classParent.Code, firstClass.Parent)
	}

	t.Logf("Parent: %v", classParent.Parent)
	t.Logf("Classes: %s", *firstClass.Parent)
	t.Logf("Child: %s", *classChild.Parent)

	classes, err := s.domain.class.Filter(ctx, class.FilterParams{
		Parent:      &classParent.Code,
		ParentCycle: true,
	}, 1, 10)
	if err != nil {
		t.Fatalf("FilterClasses: %v", err)
	}
	if len(classes.Data) != 3 {
		t.Fatalf("FilterClasses: expected 3 classes, got %d", len(classes.Data))
	}
}
