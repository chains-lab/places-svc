package controller

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal"
	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
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

type Service struct {
	domain domain
	log    logium.Logger
	cfg    internal.Config
}

func New(cfg internal.Config, log logium.Logger, class Class, place Place, placesLocale PlaceLocales, timetable Timetable) Service {
	return Service{
		domain: domain{
			class:     class,
			place:     place,
			plocale:   placesLocale,
			timetable: timetable,
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
