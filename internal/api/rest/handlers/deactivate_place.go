package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h Handler) DeactivatePlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		h.log.WithError(err).Error("deactivate place")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	locale := DetectLocale(w, r)

	place, err := h.app.DeactivatePlace(r.Context(), placeID, locale)
	if err != nil {
		h.log.WithError(err).Errorf("error deactivating place with id %s", placeID)
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(place))
}
