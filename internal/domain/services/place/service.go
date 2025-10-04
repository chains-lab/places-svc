package place

import (
	"context"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place/geo"
	"github.com/google/uuid"
)

type Service struct {
	db  database
	geo *geo.Guesser
}

func NewService(db database) Service {
	return Service{db: db}
}

type database interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	CreatePlace(ctx context.Context, input models.Place) error

	UpdatePlace(ctx context.Context, placeID uuid.UUID, params UpdateParams, updatedAt time.Time) error
	UpdateVerifiedPlace(ctx context.Context, placeID uuid.UUID, verified bool, updatedAt time.Time) error
	UpdatePlaceStatus(ctx context.Context, placeID uuid.UUID, status string, updatedAt time.Time) error

	FilterPlaces(ctx context.Context, locale string, filter FilterParams, sort SortParams, page, size uint64) (models.PlacesCollection, error)
	GetPlaceByID(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error)

	DeletePlace(ctx context.Context, placeID uuid.UUID) error

	CreatePlaceLocale(ctx context.Context, input models.PlaceLocale) error
}
