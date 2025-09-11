package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) DeletePlace(w http.ResponseWriter, r *http.Request) {
	//initiator, err := meta.User(r.Context())
	//if err != nil {
	//	a.log.WithError(err).Error("failed to get user from context")
	//	ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))
	//
	//	return
	//}

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.log.WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	err = a.app.DeletePlace(r.Context(), placeID)
	if err != nil {
		a.log.WithError(err).Error("failed to delete place")
		switch {
		case errors.Is(err, errx.ErrorPlaceForDeleteMustBeInactive):
			ape.RenderErr(w, problems.PreconditionFailed("cannot delete place that is not inactive"))
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
