package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) ActivateClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	locale := DetectLocale(w, r)

	class, err := a.app.ActivateClass(r.Context(), code, locale)
	if err != nil {
		a.log.WithError(err).Error("failed to activate class")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(class))
}
