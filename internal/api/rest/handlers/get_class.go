package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) GetClass(w http.ResponseWriter, r *http.Request) {
	locale := DetectLocale(w, r)

	class, err := a.app.GetClass(r.Context(), chi.URLParam(r, "class"), locale)
	if err != nil {
		a.log.WithError(err).WithField("class", chi.URLParam(r, "class")).Error("error getting class")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(class))
}
