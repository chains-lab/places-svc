package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/logium"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/config"
)

type Handler struct {
	app App
	log logium.Logger
	cfg config.Config
}

func NewHandler(cfg config.Config, log logium.Logger, a *app.App) Handler {
	return Handler{
		app: a,
		log: log,
		cfg: cfg,
	}
}

func (h Handler) Log(r *http.Request) logium.Logger {
	return h.log
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
