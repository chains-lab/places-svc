package class

import (
	"context"

	"github.com/chains-lab/places-svc/internal/domain/models"
)

type FilterParams struct {
	Status      *string
	Parent      *string
	ParentCycle bool
}

func (s Service) List(
	ctx context.Context,
	filter FilterParams,
	page, size uint64,
) (models.ClassesCollection, error) {
	classes, err := s.db.FilterClasses(ctx, filter, page, size)
	if err != nil {
		return models.ClassesCollection{}, err
	}

	return classes, nil
}
