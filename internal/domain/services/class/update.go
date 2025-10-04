package class

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

type UpdateParams struct {
	Name   *string
	Icon   *string
	Parent *string
}

func (s Service) Update(ctx context.Context, code string, params UpdateParams) (models.Class, error) {
	class, err := s.Get(ctx, code)
	if err != nil {
		return models.Class{}, err
	}

	if params.Parent != nil {
		if *params.Parent == code {
			return models.Class{}, errx.ErrorClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected for class with code %s", code),
			)
		}

		exist, err := s.db.ClassIsExistByCode(ctx, *params.Parent)
		if err != nil {
			return models.Class{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check parent class existence, cause: %w", err),
			)
		}
		if !exist {
			return models.Class{}, errx.ErrorParentClassNotFound.Raise(
				fmt.Errorf("parent class with code %s not found", *params.Parent),
			)
		}

		cycle, err := s.db.CheckParentCycle(ctx, code, *params.Parent)
		if err != nil {
			return models.Class{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check parent cycle for class with code %s, cause: %w", code, err),
			)
		}
		if cycle {
			return models.Class{}, errx.ErrorClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected for class with code %s", code),
			)
		}

		class.Parent = params.Parent
	}

	if params.Icon != nil {
		class.Icon = *params.Icon
	}

	if params.Name != nil {
		exist, err := s.db.ClassIsExistByName(ctx, *params.Name)
		if err != nil {
			return models.Class{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check class name uniqueness for name %s, cause: %w", *params.Name, err),
			)
		}
		if exist {
			return models.Class{}, errx.ErrorClassNameExists.Raise(
				fmt.Errorf("class with name %s already exists", *params.Name),
			)
		}

		class.Name = *params.Name
	}

	now := time.Now().UTC()

	err = s.db.UpdateClass(ctx, code, params, now)
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update class with code %s, cause: %w", code, err),
		)
	}

	return class, nil
}
