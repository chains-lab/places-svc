package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
)

type CreateParams struct {
	Code   string
	Parent *string
	Icon   string
	Name   string
}

func (c Classificator) Create(
	ctx context.Context,
	params CreateParams,
) (models.ClassWithLocale, error) {
	_, err := c.query.New().FilterCode(params.Code).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.ClassWithLocale{}, errx.ErrorInternal.Raise(fmt.Errorf("failed to check class existence, cause: %w", err))
	}
	if err == nil {
		return models.ClassWithLocale{}, errx.ErrorClassAlreadyExists.Raise(
			fmt.Errorf("class with code %s already exists", params.Code),
		)
	}

	parentValue := sql.NullString{}
	if params.Parent != nil {
		parentValue = sql.NullString{String: *params.Parent, Valid: true}
	}

	now := time.Now().UTC()

	err = c.query.New().Insert(ctx, dbx.PlaceClass{
		Code:      params.Code,
		Parent:    parentValue,
		Status:    constant.PlaceClassStatusesInactive,
		Icon:      params.Icon,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return models.ClassWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create class, cause: %w", err),
		)
	}

	err = c.localeQ.Insert(ctx, dbx.PlaceClassLocale{
		Class:  params.Code,
		Locale: constant.PlaceClassStatusesInactive,
		Name:   params.Name,
	})
	if err != nil {
		return models.ClassWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create class locale, cause: %w", err),
		)
	}

	return c.Get(ctx, params.Code, constant.PlaceClassStatusesInactive)
}
