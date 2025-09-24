package controller

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/cmd/config"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/modules/class"
	"github.com/chains-lab/places-svc/internal/domain/modules/place"
	"github.com/google/uuid"
)

type Class interface {
	Activate(
		ctx context.Context,
		code, locale string,
	) (models.ClassWithLocale, error)

	Deactivate(
		ctx context.Context,
		code, locale string,
		replaceClasses string,
	) (models.ClassWithLocale, error)

	Create(
		ctx context.Context,
		params class.CreateParams,
	) (models.ClassWithLocale, error)

	Delete(
		ctx context.Context,
		code string,
	) error

	Get(
		ctx context.Context,
		code, locale string,
	) (models.ClassWithLocale, error)

	List(
		ctx context.Context,
		locale string,
		filter class.FilterListParams,
		pag pagi.Request,
	) ([]models.ClassWithLocale, pagi.Response, error)

	LocalesList(
		ctx context.Context,
		class string,
		pag pagi.Request,
	) ([]models.ClassLocale, pagi.Response, error)

	SetLocales(
		ctx context.Context,
		code string,
		locales ...class.SetLocaleParams,
	) error

	DeleteLocale(
		ctx context.Context,
		class, locale string,
	) error

	Update(
		ctx context.Context,
		code string,
		locale string,
		params class.UpdateParams,
	) (models.ClassWithLocale, error)
}

type Place interface {
	Deactivate(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
	) (models.PlaceWithDetails, error)

	Activate(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
	) (models.PlaceWithDetails, error)

	Create(
		ctx context.Context,
		params place.CreateParams,
	) (models.PlaceWithDetails, error)

	DeleteOne(ctx context.Context, placeID uuid.UUID) error

	DeleteMany(ctx context.Context, filter place.DeleteFilter) error

	Get(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
	) (models.PlaceWithDetails, error)

	List(
		ctx context.Context,
		locale string,
		filter place.FilterListParams,
		pag pagi.Request,
		sort []pagi.SortField,
	) ([]models.PlaceWithDetails, pagi.Response, error)

	SetLocales(
		ctx context.Context,
		placeID uuid.UUID,
		locales ...place.SetLocaleParams,
	) error

	ListLocales(
		ctx context.Context,
		placeID uuid.UUID,
		pag pagi.Request,
	) ([]models.PlaceLocale, pagi.Response, error)

	DeleteLocale(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
	) error

	SetTimetable(
		ctx context.Context,
		placeID uuid.UUID,
		intervals models.Timetable,
	) (models.PlaceWithDetails, error)

	GetTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error)

	DeleteTimetable(ctx context.Context, placeID uuid.UUID) error

	Update(
		ctx context.Context,
		placeID uuid.UUID,
		locale string,
		params place.UpdateParams,
	) (models.PlaceWithDetails, error)

	UpdatePlaces(
		ctx context.Context,
		filter place.UpdatePlacesFilter,
		params place.UpdateParams,
	) error

	Verify(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error)

	Unverify(ctx context.Context, placeID uuid.UUID) (models.PlaceWithDetails, error)
}

type modules struct {
	Class
	Place
}

type Service struct {
	domain modules
	log    logium.Logger
	cfg    config.Config
}

func NewService(cfg config.Config, log logium.Logger, class class.Module, place place.Module) Service {
	return Service{
		domain: modules{
			Class: class,
			Place: place,
		},

		log: log,
		cfg: cfg,
	}
}

// DetectLocale chose locale in the following order:
// 1) ?locale=   (normalization "uk-UA" -> "uk")
// 2) Accept-Language header (normalization + q-factor sorting)
// 3) enum.DefaultLocale - default
func DetectLocale(w http.ResponseWriter, r *http.Request) string {
	if raw := r.URL.Query().Get("locale"); raw != "" {
		if loc, ok := normalizeToSupported(raw); ok {
			return loc
		}
	}

	if raw := r.Header.Get("Accept-Language"); raw != "" {
		if loc, ok := pickFromAcceptLanguage(raw); ok {
			return loc
		}
	}

	return enum.DefaultLocale
}

func normalizeToSupported(tag string) (string, bool) {
	if tag == "" {
		return "", false
	}
	primary := strings.ToLower(strings.SplitN(tag, "-", 2)[0]) // берем до '-'
	switch primary {
	case enum.LocaleEN, enum.LocaleRU, enum.LocaleUK:
		return primary, true
	default:
		return "", false
	}
}

func pickFromAcceptLanguage(header string) (string, bool) {
	type cand struct {
		tag string
		q   float64
		i   int
	}
	var items []cand

	parts := strings.Split(header, ",")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		tag := p
		q := 1.0

		if semi := strings.Index(p, ";"); semi >= 0 {
			tag = strings.TrimSpace(p[:semi])
			params := strings.Split(p[semi+1:], ";")
			for _, prm := range params {
				prm = strings.TrimSpace(prm)
				if strings.HasPrefix(prm, "q=") {
					if v, err := strconv.ParseFloat(strings.TrimPrefix(prm, "q="), 64); err == nil {
						q = v
					}
				}
			}
		}

		if tag == "" {
			continue
		}
		items = append(items, cand{tag: tag, q: q, i: i})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].q == items[j].q {
			return items[i].i < items[j].i
		}
		return items[i].q > items[j].q
	})

	for _, it := range items {
		if loc, ok := normalizeToSupported(it.tag); ok {
			return loc, true
		}
	}
	return "", false
}
