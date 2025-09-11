package class

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (c Classificator) Activate(
	ctx context.Context,
	code, locale string,
) (models.ClassWithLocale, error) {
	class, err := c.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	if class.Data.Status == constant.PlaceClassStatusesActive {
		return class, nil
	}

	status := constant.PlaceClassStatusesActive
	now := time.Now().UTC()
	err = c.query.New().FilterCode(code).Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &status,
		UpdatedAt: now,
	})
	if err != nil {
		return models.ClassWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to activate class with code %s, cause: %w", code, err),
		)
	}

	class.Data.Status = status
	class.Data.UpdatedAt = now
	return models.ClassWithLocale{
		Data:   class.Data,
		Locale: class.Locale,
	}, nil
}

func (c Classificator) Deactivate(
	ctx context.Context,
	locale, code string,
	replaceClasses string,
) (models.ClassWithLocale, error) {
	class, err := c.Get(ctx, code, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}

	replaceClass, err := c.Get(ctx, replaceClasses, locale)
	if err != nil {
		return models.ClassWithLocale{}, err
	}
	if replaceClass.Data.Code == replaceClass.Data.Code {
		return models.ClassWithLocale{}, errx.ErrorClassDeactivateReplaceSame.Raise(
			fmt.Errorf("cannot replace class %s with itself", replaceClasses),
		)
	}

	if replaceClass.Data.Status == constant.PlaceClassStatusesInactive {
		return models.ClassWithLocale{}, errx.ErrorClassDeactivateReplaceInactive.Raise(
			fmt.Errorf("cannot replace with inactive class %s", replaceClasses),
		)
	}

	if class.Data.Status == constant.PlaceClassStatusesInactive {
		return class, nil
	}

	status := constant.PlaceClassStatusesInactive
	now := time.Now().UTC()
	err = c.query.New().FilterCode(code).Update(ctx, dbx.UpdatePlaceClassParams{
		Status:    &status,
		UpdatedAt: now,
	})
	if err != nil {
		return models.ClassWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to deactivate class with code %s, cause: %w", code, err),
		)
	}

	class.Data.Status = status
	class.Data.UpdatedAt = now
	return models.ClassWithLocale{
		Data:   class.Data,
		Locale: class.Locale,
	}, nil
}
