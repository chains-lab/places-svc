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

func (a Adapter) GetClass(w http.ResponseWriter, r *http.Request) {
	locale := DetectLocale(w, r)

	class, err := a.app.GetClass(r.Context(), chi.URLParam(r, "class"), locale)
	if err != nil {
		a.log.WithError(err).WithField("class", chi.URLParam(r, "class")).Error("error getting class")
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
