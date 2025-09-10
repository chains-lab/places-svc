package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetPlaceLocales(w http.ResponseWriter, r *http.Request) {
	pag, _ := pagi.GetPagination(r)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	locales, pagResp, err := a.app.ListPlaceLocales(r.Context(), placeID, pag)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get place locales")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.PlaceLocalesCollection(locales, pagResp))
}
