package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) SetLocalesForPlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	req, err := requests.SetLocalesForPlace(r)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	locales := make([]app.SetPlaceLocalParams, 0, len(req.Data.Attributes.Locales))
	for _, attr := range req.Data.Attributes.Locales {
		locales = append(locales, app.SetPlaceLocalParams{
			Locale:      attr.Locale,
			Name:        attr.Name,
			Description: attr.Description,
		})
	}

	err = a.app.SetPlaceLocales(r.Context(), placeID, locales...)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to set place locales")
		switch {

		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
