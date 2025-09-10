package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) GetTimetable(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	timetable, err := a.app.GetTimetable(r.Context(), placeID)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to get timetable")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}
	ape.Render(w, http.StatusOK, timetable)
}
