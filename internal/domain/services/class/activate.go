package class

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

func (m Service) Activate(
	ctx context.Context,
	code, locale string,
) (models.Class, error) {
	class, err := m.Get(ctx, code, locale)
	if err != nil {
		return models.Class{}, err
	}

	if class.Status == enum.PlaceClassStatusesActive {
		return class, nil
	}

	status := enum.PlaceClassStatusesActive
	now := time.Now().UTC()
	err = m.db.Classes().FilterCode(code).Update(ctx, schemas.UpdateClassParams{
		Status:    &status,
		UpdatedAt: now,
	})
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to activate class with code %s, cause: %w", code, err),
		)
	}

	class.Status = status
	class.UpdatedAt = now
	return class, nil
}

func (m Service) Deactivate(
	ctx context.Context,
	code, locale string,
	replaceClasses string,
) (models.Class, error) {
	class, err := m.Get(ctx, code, locale)
	if err != nil {
		return models.Class{}, err
	}

	replaceClass, err := m.Get(ctx, replaceClasses, locale)
	if err != nil {
		return models.Class{}, err
	}

	if replaceClass.Code == replaceClass.Code {
		return models.Class{}, errx.ErrorClassDeactivateReplaceSame.Raise(
			fmt.Errorf("cannot replace class %s with itself", replaceClasses),
		)
	}

	if replaceClass.Status == enum.PlaceClassStatusesInactive {
		return models.Class{}, errx.ErrorClassDeactivateReplaceInactive.Raise(
			fmt.Errorf("cannot replace with inactive class %s", replaceClasses),
		)
	}

	if class.Status == enum.PlaceClassStatusesInactive {
		return class, nil
	}

	status := enum.PlaceClassStatusesInactive
	now := time.Now().UTC()
	trxErr := m.db.Transaction(ctx, func(ctx context.Context) error {
		err = m.db.Classes().FilterCode(code).Update(ctx, schemas.UpdateClassParams{
			Status:    &status,
			UpdatedAt: now,
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to deactivate class with code %s, cause: %w", code, err),
			)
		}

		err = m.db.Places().FilterClass(code).Update(ctx, schemas.UpdatePlaceParams{
			Class:     &replaceClasses,
			UpdatedAt: now,
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to update places with class %s to class %s, cause: %w", code, replaceClasses, err),
			)
		}

		return nil
	})
	if trxErr != nil {
		return models.Class{}, trxErr
	}

	class.Status = status
	class.UpdatedAt = now
	return class, nil
}
