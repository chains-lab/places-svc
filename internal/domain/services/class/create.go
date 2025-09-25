package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

type CreateParams struct {
	Code   string
	Parent *string
	Icon   string
	Name   string
}

func (m Service) Create(
	ctx context.Context,
	params CreateParams,
) (models.Class, error) {
	_, err := m.db.Classes().FilterCode(params.Code).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Class{}, errx.ErrorInternal.Raise(fmt.Errorf("failed to check class existence, cause: %w", err))
	}
	if err == nil {
		return models.Class{}, errx.ErrorClassCodeAlreadyTaken.Raise(
			fmt.Errorf("class with code %s already exists", params.Code),
		)
	}

	parentValue := sql.NullString{}
	if params.Parent != nil {
		parentValue = sql.NullString{String: *params.Parent, Valid: true}
	}

	_, err = m.db.Classes().FilterName(params.Name).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Class{}, errx.ErrorInternal.Raise(fmt.Errorf("failed to check class locale existence, cause: %w", err))
	}
	if err == nil {
		return models.Class{}, errx.ErrorClassNameAlreadyTaken.Raise(
			fmt.Errorf("class locale with name %s already exists", params.Name),
		)
	}

	now := time.Now().UTC()

	trxErr := m.db.Transaction(ctx, func(ctx context.Context) error {
		err = m.db.Classes().Insert(ctx, schemas.Class{
			Code:      params.Code,
			Parent:    parentValue,
			Status:    enum.PlaceClassStatusesInactive,
			Icon:      params.Icon,
			Name:      params.Name,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create class, cause: %w", err),
			)
		}

		return nil
	})
	if trxErr != nil {
		return models.Class{}, trxErr
	}

	return m.Get(ctx, params.Code)
}
