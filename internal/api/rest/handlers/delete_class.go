package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) DeleteClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "class_code")

	err := a.app.DeleteClass(r.Context(), code)
	if err != nil {
		a.log.WithError(err).WithField("class_code", code).Error("error deleting place")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
