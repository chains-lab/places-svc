package domain_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/chains-lab/places-svc/internal/data"
	"github.com/chains-lab/places-svc/internal/domain/infra/geo"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
	"github.com/chains-lab/places-svc/internal/domain/services/timetable"
	"github.com/chains-lab/places-svc/test"
	"github.com/google/uuid"
)

type Class interface {
	Create(
		ctx context.Context,
		params class.CreateParams,
	) (models.Class, error)

	Filter(
		ctx context.Context,
		filter class.FilterParams,
		page, size uint64,
	) (models.ClassesCollection, error)
	Get(ctx context.Context, code string) (models.Class, error)

	Activate(ctx context.Context, code string) (models.Class, error)
	Deactivate(ctx context.Context, code, replaceCode string) (models.Class, error)

	Update(ctx context.Context, code string, params class.UpdateParams) (models.Class, error)

	Delete(
		ctx context.Context,
		code string,
	) error
}

type Place interface {
	Create(
		ctx context.Context,
		params place.CreateParams,
	) (models.Place, error)

	Filter(
		ctx context.Context,
		locale string,
		filter place.FilterParams,
		sort place.SortParams,
		page, size uint64,
	) (models.PlacesCollection, error)
	Get(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error)

	Update(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
		params place.UpdateParams,
	) (models.Place, error)
	UpdateStatus(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
		status string,
	) (models.Place, error)
	Block(ctx context.Context, placeID uuid.UUID, locale string, block bool) (models.Place, error)
	Verify(ctx context.Context, placeID uuid.UUID, locale string, value bool) (models.Place, error)

	Delete(ctx context.Context, placeID uuid.UUID) error
}

type PlaceLocales interface {
	SetForPlace(
		ctx context.Context,
		placeID uuid.UUID,
		locales ...plocale.SetParams,
	) error

	GetForPlace(
		ctx context.Context,
		placeID uuid.UUID,
		page uint64,
		size uint64,
	) (models.PlaceLocaleCollection, error)

	Delete(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
	) error
}

type Timetable interface {
	SetForPlace(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
		intervals models.Timetable,
	) (models.Place, error)

	GetForPlace(ctx context.Context, placeID uuid.UUID) (models.Timetable, error)

	DeleteForPlace(ctx context.Context, placeID uuid.UUID) error
}

type domain struct {
	class     Class
	place     Place
	plocale   PlaceLocales
	timetable Timetable
}

type Setup struct {
	domain domain
}

func newSetup(t *testing.T) (Setup, error) {
	pg, err := sql.Open("postgres", test.TestDatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", "error", err)
	}

	database := data.New(pg)

	geoGuesser := geo.NewGuesser()

	classSvc := class.NewService(database)
	placeSvc := place.NewService(database, geoGuesser)
	pLocalesSvc := plocale.NewService(database)
	timetableSvc := timetable.NewService(database)

	return Setup{
		domain: domain{
			class:     classSvc,
			place:     placeSvc,
			plocale:   pLocalesSvc,
			timetable: timetableSvc,
		},
	}, nil
}

func CreateClass(s Setup, t *testing.T, name, code string, parent *string) models.Class {
	t.Helper()
	c, err := s.domain.class.Create(context.Background(), class.CreateParams{
		Name:   name,
		Code:   code,
		Icon:   "icon",
		Parent: parent,
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}

	c, err = s.domain.class.Activate(context.Background(), code)
	if err != nil {
		t.Fatalf("ActivateClass: %v", err)
	}

	return c
}

func CreatePlace(s Setup, t *testing.T, params place.CreateParams) models.Place {
	t.Helper()
	p, err := s.domain.place.Create(context.Background(), params)
	if err != nil {
		t.Fatalf("CreatePlace: %v", err)
	}

	return p
}
