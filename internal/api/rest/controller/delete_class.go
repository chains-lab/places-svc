package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
)

func (h Service) DeleteClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "class_code")

	err := h.domain.Class.Delete(r.Context(), code)
	if err != nil {
		h.log.WithError(err).WithField("class_code", code).Error("error deleting place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class with code %q not found", code)))
		case errors.Is(err, errx.ErrorCannotDeleteActiveClass):
			ape.RenderErr(w, problems.Forbidden("cannot delete active class"))
		case errors.Is(err, errx.ErrorCantDeleteClassWithPlaces):
			ape.RenderErr(w, problems.Forbidden("cannot delete class with places"))
		case errors.Is(err, errx.ErrorClassHasChildren):
			ape.RenderErr(w, problems.Forbidden("cannot delete class with children"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
