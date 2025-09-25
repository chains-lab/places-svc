package class

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

type FilterListParams struct {
	Status      *string
	Parent      *string
	ParentCycle bool
	Page        uint
	Size        uint
}

func (m Service) List(ctx context.Context, filter FilterListParams) (models.ClassesCollection, error) {
	limit, offset := pagi.PagConvert(filter.Page, filter.Size)

	query := m.db.Classes()

	if filter.Parent != nil {
		_, err := m.Get(ctx, *filter.Parent)
		if errors.Is(err, errx.ErrorClassNotFound) {
			return models.ClassesCollection{}, errx.ErrorParentClassNotFound.Raise(
				fmt.Errorf("parent class not found: %s", *filter.Parent),
			)
		}
		if err != nil {
			return models.ClassesCollection{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get parent class: %w", err),
			)
		}

		if filter.ParentCycle {
			query = query.FilterParentCycle(*filter.Parent)
		} else {
			query = query.FilterParent(sql.NullString{
				String: *filter.Parent,
				Valid:  true,
			})
		}

	}

	if filter.Status != nil {
		err := enum.IsValidPlaceStatus(*filter.Status)
		if err != nil {
			return models.ClassesCollection{}, errx.ErrorClassStatusInvalid.Raise(
				fmt.Errorf("invalid status filter: %s", *filter.Status),
			)
		}
		query = query.FilterStatus(*filter.Status)
	}

	query = query.Page(limit, offset)

	rows, err := query.Select(ctx)
	if err != nil {
		return models.ClassesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to select classes, cause: %w", err),
		)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return models.ClassesCollection{}, errx.ErrorInternal.Raise(
			fmt.Errorf("internal error, cause: %w", err),
		)
	}

	classes := make([]models.Class, 0, len(rows))
	for _, r := range rows {
		classes = append(classes, modelFromDB(r))
	}

	return models.ClassesCollection{
		Data:  classes,
		Page:  filter.Page,
		Size:  filter.Size,
		Total: count,
	}, nil
}
