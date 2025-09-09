package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
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
		ape.RenderErr(w, problems.InternalError())

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
