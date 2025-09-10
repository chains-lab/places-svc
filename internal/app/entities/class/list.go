package class

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
)

type FilterListParams struct {
	Parent      *string
	ParentCycle *bool
	Status      *string
}

func (c Classificator) List(
	ctx context.Context,
	locale string,
	filter FilterListParams,
	pag pagi.Request,
) ([]models.ClassWithLocale, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	query := c.query.New()

	if filter.Parent != nil {
		if filter.ParentCycle != nil && *filter.ParentCycle {
			query = query.FilterParentCycle(*filter.Parent)
		}
		query = query.FilterParent(sql.NullString{
			String: *filter.Parent,
			Valid:  true,
		})
	}
	if filter.Status != nil {
		query = query.FilterStatus(*filter.Status)
	}

	query = query.Page(limit, offset)

	l := locale
	err := constant.IsValidLocaleSupported(l)
	if err != nil {
		l = constant.DefaultLocale
	}

	rows, err := query.SelectWithLocale(ctx, l)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to select classes, cause: %w", err),
		)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	classes := make([]models.ClassWithLocale, 0, len(rows))
	for _, r := range rows {
		classes = append(classes, classWithLocaleModelFromDB(r))
	}

	return classes, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}
