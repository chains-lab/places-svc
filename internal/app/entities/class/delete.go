package class

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (c Classificator) DeleteClass(
	ctx context.Context,
	code string,
) error {
	_, err := c.Get(ctx, code, constant.LocaleEN)
	if err != nil {
		return err
	}

	count, err := c.query.New().FilterParent(sql.NullString{String: code, Valid: true}).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check class children existence, cause: %w", err),
		)
	}
	if count > 0 {
		return errx.ErrorClassHasChildren.Raise(
			fmt.Errorf("class with code %s has children, cannot be deleted", code),
		)
	}

	err = c.localeQ.New().FilterClass(code).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete class locales, cause: %w", err),
		)
	}

	return nil
}
