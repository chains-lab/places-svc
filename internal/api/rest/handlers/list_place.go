package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/google/uuid"
)

func (a Adapter) ListPlace(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var filters app.FilterListPlaces

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
		if verified == "true" {
			t := true
			filters.Verified = &t
		} else if verified == "false" {
			f := false
			filters.Verified = &f
		} else {
			ape.RenderErr(w, problems.InvalidParameter(
				"verified",
				fmt.Errorf("invalid boolean value: %s", verified)),
			)
			return
		}
	}

	if point := strings.TrimSpace(q.Get("point")); point != "" {
		filters.Location = &app.GeoFilterListPlaces{}
		parts := strings.Split(point, ",")
		if len(parts) != 2 {
			ape.RenderErr(w, problems.InvalidParameter("point", fmt.Errorf("invalid point format, expected 'lon,lat'")))
			return
		}

		var lon, lat float64
		_, err := fmt.Sscanf(parts[0], "%f", &lon)
		if err != nil {
			ape.RenderErr(w, problems.InvalidParameter("point", fmt.Errorf("invalid longitude value: %v", err)))
			return
		}
		_, err = fmt.Sscanf(parts[1], "%f", &lat)
		if err != nil {
			ape.RenderErr(w, problems.InvalidParameter("point", fmt.Errorf("invalid latitude value: %v", err)))
			return
		}
		filters.Location = &app.GeoFilterListPlaces{
			Point: [2]float64{lon, lat},
		}
	}

	if radius := strings.TrimSpace(q.Get("radius")); radius != "" {
		var r uint64
		_, err := fmt.Sscanf(radius, "%d", &r)
		if err != nil || r == 0 {
			ape.RenderErr(w, problems.InvalidParameter("radius", fmt.Errorf("invalid radius value: %v", err)))
			return
		}
		filters.Location.RadiusM = r
	}

	pag, sort := pagi.GetPagination(r)

	places, pagResp, err := a.app.ListPlaces(r.Context(), DetectLocale(w, r), filters, pag, sort)
	if err != nil {
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.PlacesCollection(places, pagResp))
}
