package class

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

type CreateParams struct {
	Code   string
	Parent *string
	Icon   string
	Name   string
}

func (s Service) Create(
	ctx context.Context,
	params CreateParams,
) (models.Class, error) {
	exist, err := s.db.ClassIsExistByCode(ctx, params.Code)
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check class existence, cause: %w", err),
		)
	}

	if exist {
		return models.Class{}, errx.ErrorClassCodeAlreadyTaken.Raise(
			fmt.Errorf("class with code %s already exists", params.Code),
		)
	}

	exist, err = s.db.ClassIsExistByName(ctx, params.Name)
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check class name existence, cause: %w", err),
		)
	}

	if exist {
		return models.Class{}, errx.ErrorClassNameAlreadyTaken.Raise(
			fmt.Errorf("class with name %s already exists", params.Name),
		)
	}

	now := time.Now().UTC()
	class := models.Class{
		Code:      params.Code,
		Parent:    params.Parent,
		Status:    enum.PlaceClassStatusesInactive,
		Icon:      params.Icon,
		Name:      params.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CreateClass(ctx, class)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create class, cause: %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Class{}, err
	}

	return class, nil
}
