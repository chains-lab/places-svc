package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
)

func (s Service) ActivateClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	res, err := s.domain.class.Activate(r.Context(), code)
	if err != nil {
		s.log.WithError(err).Error("failed to activate class")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound("class not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(res))
}
