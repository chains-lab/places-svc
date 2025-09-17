package handlers

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/config"
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
}

type Adapter struct {
	app *app.App
	log logium.Logger
	cfg config.Config
}

func NewAdapter(cfg config.Config, log logium.Logger, a *app.App) Adapter {
	return Adapter{
		app: a,
		log: log,
		cfg: cfg,
	}
}

func (a Adapter) Log(r *http.Request) logium.Logger {
	return a.log
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
