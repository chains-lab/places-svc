package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

type UpdateParams struct {
	Name   *string
	Icon   *string
	Parent *string
}

func (m Service) Update(ctx context.Context, code string, params UpdateParams) (models.Class, error) {
	class, err := m.Get(ctx, code)
	if err != nil {
		return models.Class{}, err
	}

	stmt := schemas.UpdateClassParams{
		UpdatedAt: time.Now().UTC(),
	}
	class.UpdatedAt = stmt.UpdatedAt

	if params.Parent != nil {
		if *params.Parent == code {
			return models.Class{}, errx.ErrorClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected for class with code %s", code),
			)
		}

		_, err = m.db.Classes().FilterCode(*params.Parent).Get(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.Class{}, errx.ErrorParentClassNotFound.Raise(
					fmt.Errorf("parent class with code %s not found", *params.Parent),
				)
			}
		}

		_, err = m.db.Classes().FilterParentCycle(class.Code).FilterCode(*params.Parent).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return models.Class{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check parent cycle for class with code %s, cause: %w", code, err),
			)
		}
		if err == nil {
			return models.Class{}, errx.ErrorClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected for class with code %s", code),
			)
		}

		stmt.Parent = &sql.NullString{String: *params.Parent, Valid: true}
		class.Parent = params.Parent
	}

	if params.Icon != nil {
		stmt.Icon = params.Icon
		class.Icon = *params.Icon
	}

	if params.Name != nil {
		_, err = m.db.Classes().FilterName(*params.Name).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return models.Class{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check class name uniqueness for name %s, cause: %w", *params.Name, err),
			)
		}
		if err == nil {
			return models.Class{}, errx.ErrorClassNameExists.Raise(
				fmt.Errorf("class with name %s already exists", *params.Name),
			)
		}

		stmt.Name = params.Name
		class.Name = *params.Name
	}

	err = m.db.Classes().FilterCode(code).Update(ctx, stmt)
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update class with code %s, cause: %w", code, err),
		)
	}

	return class, nil
}
