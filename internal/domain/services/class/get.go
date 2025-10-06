package class

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
)

func (s Service) Get(ctx context.Context, code string) (models.Class, error) {
	class, err := s.db.GetClassByCode(ctx, code)
	if err != nil {
		return models.Class{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get class with code %s, cause: %w", code, err),
		)
	}

	if class.IsNil() {
		return models.Class{}, errx.ErrorClassNotFound.Raise(
			fmt.Errorf("class with code %s not found", code),
		)
	}

	return class, nil
}
