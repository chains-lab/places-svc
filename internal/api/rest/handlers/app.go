package handlers

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

type App interface {
	//CLASS
	CreateClass(ctx context.Context, params app.CreateClassParams) (models.ClassWithLocale, error)

	GetClass(ctx context.Context, code, locale string) (models.ClassWithLocale, error)

	ActivateClass(ctx context.Context, code, locale string) (models.ClassWithLocale, error)
	DeactivateClass(ctx context.Context, code, locale, replace string) (models.ClassWithLocale, error)

	SetClassLocales(ctx context.Context, code string, locales ...app.SetClassLocaleParams) error

	DeleteClass(ctx context.Context, code string) error

	ListClasses(
		ctx context.Context,
		locale string,
		filter app.FilterListClassesParams,
		pag pagi.Request,
	) ([]models.ClassWithLocale, pagi.Response, error)

	ListClassLocales(
		ctx context.Context,
		class string,
		pag pagi.Request,
	) ([]models.ClassLocale, pagi.Response, error)

	UpdateClass(
		ctx context.Context,
		code, locale string,
		params app.UpdateClassParams,
	) (models.ClassWithLocale, error)

	//PLACE

	CreatePlace(ctx context.Context, params app.CreatePlaceParams) (models.PlaceWithDetails, error)

	GetPlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithDetails, error)

	ActivatePlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithDetails, error)
	DeactivatePlace(ctx context.Context, placeID uuid.UUID, locale string) (models.PlaceWithDetails, error)
	VerifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error)
	UnverifyPlace(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error)

	DeleteTimetable(ctx context.Context, placeID uuid.UUID) error
	DeletePlace(ctx context.Context, placeID uuid.UUID) error

	ListPlaceLocales(ctx context.Context, placeID uuid.UUID, pag pagi.Request) ([]models.PlaceLocale, pagi.Response, error)

	ListPlaces(
		ctx context.Context,
		locale string,
		filter app.FilterListPlaces,
		pag pagi.Request,
		sort []pagi.SortField,
	) ([]models.PlaceWithDetails, pagi.Response, error)

	SetPlaceTimeTable(
		ctx context.Context,
		placeID uuid.UUID,
		intervals models.Timetable,
	) (models.PlaceWithDetails, error)

	SetPlaceLocales(
		ctx context.Context,
		placeID uuid.UUID,
		locales ...app.SetPlaceLocalParams,
	) error

	UpdatePlace(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
		params app.UpdatePlaceParams,
	) (models.PlaceWithDetails, error)

	GetTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error)
}
