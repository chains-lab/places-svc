package class

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

func (s Service) Activate(ctx context.Context, code string) (models.Class, error) {
	class, err := s.Get(ctx, code)
	if err != nil {
		return models.Class{}, err
	}

	if class.Status == enum.PlaceClassStatusesActive {
		return class, nil
	}

	now := time.Now().UTC()
	err = s.db.UpdateClassStatus(ctx, code, enum.PlaceClassStatusesActive, now)
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to activate class with code %s, cause: %w", code, err),
		)
	}

	class.Status = enum.PlaceClassStatusesActive
	class.UpdatedAt = now

	return class, nil
}

func (s Service) Deactivate(ctx context.Context, code, replaceCode string) (models.Class, error) {
	if code == replaceCode {
		return models.Class{}, errx.ErrorClassDeactivateReplaceSame.Raise(
			fmt.Errorf("cannot replace class %s with itself", replaceCode),
		)
	}

	current, err := s.Get(ctx, code)
	if err != nil {
		return models.Class{}, err
	}

	replaceClass, err := s.Get(ctx, replaceCode)
	if err != nil {
		if errors.Is(err, errx.ErrorClassCodeAlreadyTaken) {
			return models.Class{}, errx.ErrorReplaceClassNotFound.Raise(
				fmt.Errorf("replace class with code %s not found", replaceCode),
			)
		}
		return models.Class{}, err
	}

	if replaceClass.Status == enum.PlaceClassStatusesInactive {
		return models.Class{}, errx.ErrorClassDeactivateReplaceInactive.Raise(
			fmt.Errorf("cannot replace with inactive class %s", replaceCode),
		)
	}

	if current.Status == enum.PlaceClassStatusesInactive {
		return current, nil
	}

	now := time.Now().UTC()

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.UpdateClassStatus(ctx, code, enum.PlaceClassStatusesInactive, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to deactivate class with code %s, cause: %w", code, err),
			)
		}

		err = s.db.ReplaceClassInPlaces(ctx, code, replaceCode, now)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update places with class %s to class %s, cause: %w", code, replaceCode, err),
			)
		}

		return nil
	}); err != nil {
		return models.Class{}, err
	}

	current.Status = enum.PlaceClassStatusesInactive
	current.UpdatedAt = now

	return current, nil
}
