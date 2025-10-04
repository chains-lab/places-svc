package plocale

import (
	"context"

	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

type Service struct {
	db database
}

func NewService(db database) Service {
	return Service{db: db}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	PlaceExists(ctx context.Context, placeID uuid.UUID) (bool, error)
	GetPlaceByID(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error)

	UpsertLocaleForPlace(ctx context.Context, locales ...SetParams) error

	GetPlaceLocales(
		ctx context.Context,
		placeID uuid.UUID,
		page, size uint64,
	) (models.PlaceLocaleCollection, error)
	CountPlaceLocales(ctx context.Context, placeID uuid.UUID) (uint64, error)

	DeletePlaceLocale(ctx context.Context, placeID uuid.UUID, locale string) error
}
