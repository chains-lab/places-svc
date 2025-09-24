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

func (h Service) GetPlace(w http.ResponseWriter, r *http.Request) {
	locale := DetectLocale(w, r)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		h.log.WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	res, err := h.domain.Place.Get(r.Context(), placeID, locale)
	if err != nil {
		h.log.WithError(err).WithField("place_id", placeID).Error("error getting place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
