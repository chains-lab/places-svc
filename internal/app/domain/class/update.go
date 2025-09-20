package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
)

type UpdateParams struct {
	Name   *string
	Icon   *string
	Parent *string
}

func (c Classificator) Update(
	ctx context.Context,
	code string,
	params UpdateParams,
) error {
	class, err := c.Get(ctx, code, enum.LocaleEN)
	if err != nil {
		return err
	}

	stmt := dbx.UpdatePlaceClassParams{
		UpdatedAt: time.Now().UTC(),
	}

	if params.Parent != nil {
		if *params.Parent == code {
			return errx.ErrorClassParentEqualCode.Raise(
				fmt.Errorf("parent cycle detected for class with code %s", code),
			)
		}
		_, err = c.query.New().FilterParentCycle(class.Data.Code).FilterCode(*params.Parent).Get(ctx)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to check parent cycle for class with code %s, cause: %w", code, err),
			)
		}

		if err == nil {
			return errx.ErrorClassParentCycle.Raise(
				fmt.Errorf("parent cycle detected for class with code %s", code),
			)
		}
		stmt.Parent = &sql.NullString{String: *params.Parent, Valid: true}
	}
	
	if params.Icon != nil {
		stmt.Icon = params.Icon
	}

	err = c.query.New().FilterCode(code).Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update class with code %s, cause: %w", code, err),
		)
	}

	return nil
}
