package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/places-svc/internal/rest/responses"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) FilterPlace(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var filters place.FilterParams

	if cityID := strings.TrimSpace(q.Get("city_id")); cityID != "" {
		id, err := uuid.Parse(cityID)
		if err != nil {
			s.log.WithError(err).Error("invalid city_id")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("failed to parse city_id: %w", err),
			})...)
			return
		}
		filters.CityID = &id
	}

	if distributorID := strings.TrimSpace(q.Get("distributor_id")); distributorID != "" {
		id, err := uuid.Parse(distributorID)
		if err != nil {
			s.log.WithError(err).Error("invalid distributor_id")
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("failed to parse distributor_id: %w", err),
			})...)
			return
		}
		filters.DistributorID = &id
	}

	if name := strings.TrimSpace(q.Get("name")); name != "" {
		filters.Name = &[]string{name}[0]
	}

	if address := strings.TrimSpace(q.Get("address")); address != "" {
		filters.Address = &[]string{address}[0]
	}

	if classes := q["class"]; len(classes) > 0 {
		filters.Classes = classes
	}

	if statuses := q["status"]; len(statuses) > 0 {
		filters.Statuses = statuses
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
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("invalid verified value: %s", verified),
			})...)
			return
		}
	}

	var geo *place.FilterDistance
	if point := strings.TrimSpace(q.Get("point")); point != "" {
		parts := strings.Split(point, ",")
		if len(parts) != 2 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"point": fmt.Errorf("expected 'lon,lat', got %q", point),
			})...)

			return
		}
		var lon, lat float64
		if _, err := fmt.Sscanf(parts[0], "%f", &lon); err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"point": fmt.Errorf("invalid longitude value: %v", err),
			})...)

			return
		}
		if _, err := fmt.Sscanf(parts[1], "%f", &lat); err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"point": fmt.Errorf("invalid latitude value: %v", err),
			})...)
			return
		}
		geo = &place.FilterDistance{Point: [2]float64{lon, lat}}
	}

	if radius := strings.TrimSpace(q.Get("radius")); radius != "" {
		var rM uint64
		if _, err := fmt.Sscanf(radius, "%d", &rM); err != nil || rM == 0 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"radius": fmt.Errorf("invalid radius value: %v", err),
			})...)

			return
		}
		if geo == nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"radius": errors.New("the 'point' parameter is required when 'radius' is provided"),
			})...)

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
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"time": errors.New("both 'time_from' and 'time_to' parameters are required"),
			})...)

			return
		}
		from, err := parseMomentParam(tf)
		if err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"time_from": err,
			})...)

			return
		}
		to, err := parseMomentParam(tt)
		if err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"time_to": err,
			})...)

			return
		}
		filters.Time = &models.TimeInterval{From: from, To: to}
	}

	pag, size := pagi.GetPagination(r)

	sorts := pagi.SortFields(r)

	sort := place.SortParams{}
	for _, s := range sorts {
		switch s.Field {
		case "created_at":
			sort.ByCreatedAt = &s.Ascend
		case "distance":
			if geo == nil {
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"sort": errors.New("the 'point' parameter is required when sorting by distance"),
				})...)

				return
			}

			sort.ByDistance = &s.Ascend
		}

	}

	places, err := s.domain.place.Filter(r.Context(), DetectLocale(w, r), filters, sort, pag, size)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		return
	}
	ape.Render(w, http.StatusOK, responses.PlacesCollection(places))
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
