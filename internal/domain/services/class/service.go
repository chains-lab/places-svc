package class

import (
	"context"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/models"
)

type Service struct {
	db database
}

func NewService(db database) Service {
	return Service{db: db}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreateClass(ctx context.Context, class models.Class) (models.Class, error)

	GetClassByCode(ctx context.Context, code string) (models.Class, error)

	ClassIsExistByCode(ctx context.Context, code string) (bool, error)
	ClassIsExistByName(ctx context.Context, name string) (bool, error)

	CountClassChildren(ctx context.Context, parentCode string) (uint64, error)
	CountPlacesByClass(ctx context.Context, classCode string) (uint64, error)

	UpdateClass(ctx context.Context, code string, params UpdateParams, updateAt time.Time) error
	UpdateClassStatus(ctx context.Context, code string, status string, updateAt time.Time) error
	ReplaceClassInPlaces(ctx context.Context, oldCode, newCode string, updateAt time.Time) error

	FilterClasses(ctx context.Context, filter FilterParams, page, size uint64) (models.ClassesCollection, error)

	CheckParentCycle(ctx context.Context, classCode, parentCode string) (bool, error)

	DeleteClass(ctx context.Context, code string) error
}

// TODO check for parent cycle
//_, err = m.db.Classes().FilterParentCycle(class.Code).FilterCode(*params.Parent).Get(ctx)
//if err != nil && !errors.Is(err, sql.ErrNoRows) {
//return models.Class{}, errx.ErrorInternal.Raise(
//fmt.Errorf("failed to check parent cycle for class with code %s, cause: %w", code, err),
//)
//}
//if err == nil {
//return models.Class{}, errx.ErrorClassParentCycle.Raise(
//fmt.Errorf("parent cycle detected for class with code %s", code),
//)
//}
