package class

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
)

func (m Service) Delete(
	ctx context.Context,
	code string,
) error {
	_, err := m.Get(ctx, code)
	if err != nil {
		return err
	}

	count, err := m.db.Classes().FilterParent(sql.NullString{String: code, Valid: true}).Count(ctx)
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

	count, err = m.db.Places().FilterClass(code).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check places with class existence, cause: %w", err),
		)
	}
	if count > 0 {
		return errx.ErrorCantDeleteClassWithPlaces.Raise(
			fmt.Errorf("failed to delete class %s with active places", code),
		)
	}

	err = m.db.Classes().FilterCode(code).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete class locales, cause: %w", err),
		)
	}

	return nil
}
