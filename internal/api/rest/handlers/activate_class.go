package handlers

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) ActivateClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	locale := DetectLocale(w, r)

	class, err := a.app.ActivateClass(r.Context(), code, locale)
	if err != nil {
		a.log.WithError(err).Error("failed to activate class")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound("class not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(class))
}
