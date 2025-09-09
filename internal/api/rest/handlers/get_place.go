package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetPlace(w http.ResponseWriter, r *http.Request) {
	locale := DetectLocale(w, r)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	place, err := a.app.GetPlace(r.Context(), placeID, locale)
	if err != nil {
		a.Log(r).WithError(err).WithField("place_id", placeID).Error("error getting place")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(place))
}
