package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (s Service) GetClass(w http.ResponseWriter, r *http.Request) {
	res, err := s.domain.class.Get(r.Context(), chi.URLParam(r, "class"))
	if err != nil {
		s.log.WithError(err).WithField("class", chi.URLParam(r, "class")).Error("error getting class")
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
