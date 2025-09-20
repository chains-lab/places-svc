package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/google/uuid"
)

func (h Handler) ListPlace(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var filters app.FilterListPlaces

	// -------- базовые фильтры --------
	if cities := q["city_id"]; len(cities) > 0 {
		cityIDs := make([]uuid.UUID, 0, len(cities))
		for _, c := range cities {
			id, err := uuid.Parse(c)
			if err != nil {
				ape.RenderErr(w, problems.InvalidParameter("city_id", fmt.Errorf("invalid uuid: %s", c)))
				return
			}
			cityIDs = append(cityIDs, id)
		}
		filters.CityIDs = cityIDs
	}

	if distributors := q["distributor_id"]; len(distributors) > 0 {
		distributorIDs := make([]uuid.UUID, 0, len(distributors))
		for _, d := range distributors {
			id, err := uuid.Parse(d)
			if err != nil {
				ape.RenderErr(w, problems.InvalidParameter("distributor_id", fmt.Errorf("invalid uuid: %s", d)))
				return
			}
			distributorIDs = append(distributorIDs, id)
		}
		filters.DistributorIDs = distributorIDs
	}

	if classes := q["class"]; len(classes) > 0 {
		filters.Classes = classes
	}
	if statuses := q["status"]; len(statuses) > 0 {
		filters.Statuses = statuses
	}
	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &[]string{name}[0]
	}
	if address := strings.TrimSpace(q.Get("address")); address != "" {
		filters.Address = &[]string{address}[0]
	}
	if verified := strings.TrimSpace(q.Get("verified")); verified != "" {
		switch verified {
		case "true":
			t := true
			filters.Verified = &t
		case "false":
			f := false
			filters.Verified = &f
		default:
			ape.RenderErr(w, problems.InvalidParameter("verified", fmt.Errorf("invalid boolean value: %s", verified)))
			return
		}
	}

	// -------- гео фильтр (point + radius) --------
	var geo *app.GeoFilterListPlaces
	if point := strings.TrimSpace(q.Get("point")); point != "" {
		parts := strings.Split(point, ",")
		if len(parts) != 2 {
			ape.RenderErr(w, problems.InvalidParameter("point", fmt.Errorf("invalid point format, expected 'lon,lat'")))
			return
		}
		var lon, lat float64
		if _, err := fmt.Sscanf(parts[0], "%f", &lon); err != nil {
			ape.RenderErr(w, problems.InvalidParameter("point", fmt.Errorf("invalid longitude value: %v", err)))
			return
		}
		if _, err := fmt.Sscanf(parts[1], "%f", &lat); err != nil {
			ape.RenderErr(w, problems.InvalidParameter("point", fmt.Errorf("invalid latitude value: %v", err)))
			return
		}
		geo = &app.GeoFilterListPlaces{Point: [2]float64{lon, lat}}
	}
	if radius := strings.TrimSpace(q.Get("radius")); radius != "" {
		var rM uint64
		if _, err := fmt.Sscanf(radius, "%d", &rM); err != nil || rM == 0 {
			ape.RenderErr(w, problems.InvalidParameter("radius", fmt.Errorf("invalid radius value: %v", err)))
			return
		}
		if geo == nil {
			ape.RenderErr(w, problems.InvalidParameter("radius", errors.New("radius requires 'point'")))
			return
		}
		geo.RadiusM = rM
	}
	if geo != nil && geo.RadiusM > 0 {
		filters.Location = geo
	}

	tf := strings.TrimSpace(q.Get("time_from"))
	tt := strings.TrimSpace(q.Get("time_to"))
	if tf != "" || tt != "" {
		if tf == "" || tt == "" {
			ape.RenderErr(w, problems.InvalidParameter("time_from,time_to", errors.New("both time_from and time_to must be provided")))
			return
		}
		from, err := parseMomentParam(tf)
		if err != nil {
			ape.RenderErr(w, problems.InvalidParameter("time_from", err))
			return
		}
		to, err := parseMomentParam(tt)
		if err != nil {
			ape.RenderErr(w, problems.InvalidParameter("time_to", err))
			return
		}
		filters.Time = &models.TimeInterval{From: from, To: to}
	}

	pag, sort := pagi.GetPagination(r)

	places, pagResp, err := h.app.ListPlaces(r.Context(), DetectLocale(w, r), filters, pag, sort)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}
	ape.Render(w, http.StatusOK, responses.PlacesCollection(places, pagResp))
}

func parseMomentParam(v string) (models.Moment, error) {
	v = strings.TrimSpace(v)
	parts := strings.Fields(v)
	if len(parts) != 2 {
		return models.Moment{}, fmt.Errorf("expected '<weekday> HH:MM', got %q", v)
	}
	wdToken := strings.ToLower(parts[0])
	wd, ok := wdMap[wdToken]
	if !ok {
		return models.Moment{}, fmt.Errorf("unknown weekday: %q", parts[0])
	}

	hh, mm := 0, 0
	if _, err := fmt.Sscanf(parts[1], "%d:%d", &hh, &mm); err != nil {
		return models.Moment{}, fmt.Errorf("invalid time %q: %v", parts[1], err)
	}
	if hh < 0 || hh > 23 || mm < 0 || mm > 59 {
		return models.Moment{}, fmt.Errorf("time out of range: %02d:%02d", hh, mm)
	}

	d := time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute
	return models.Moment{Weekday: wd, Time: d}, nil
}
