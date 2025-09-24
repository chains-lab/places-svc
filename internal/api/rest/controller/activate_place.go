package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h Service) ActivatePlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		h.log.WithError(err).Error("activate place")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	locale := DetectLocale(w, r)

	res, err := h.domain.Place.Activate(r.Context(), placeID, locale)
	if err != nil {
		h.log.WithError(err).WithField("place_id", placeID).Error("error activating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound("class not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
