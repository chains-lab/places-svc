package class

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
)

func (s Service) Delete(
	ctx context.Context,
	code string,
) error {
	_, err := s.Get(ctx, code)
	if err != nil {
		return err
	}

	count, err := s.db.CountClassChildren(ctx, code)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check class children existence, cause: %w", err),
		)
	}
	if count > 0 {
		return errx.ErrorCannotDeleteClassWithChildren.Raise(
			fmt.Errorf("failed to delete class %s with active children", code),
		)
	}

	count, err = s.db.CountPlacesByClass(ctx, code)
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

	err = s.db.DeleteClass(ctx, code)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete class %s, cause: %w", code, err),
		)
	}

	return nil
}
