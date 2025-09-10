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

func (a Adapter) ListPlaceLocales(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.log.WithError(err).Error("error parsing place_id")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	pag, _ := pagi.GetPagination(r)

	locales, pagResp, err := a.app.ListPlaceLocales(r.Context(), placeID, pag)
	if err != nil {
		a.log.WithError(err).Error("error listing place locales")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.PlaceLocalesCollection(locales, pagResp))
}
