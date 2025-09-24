package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/go-chi/chi/v5"
)

func (h Service) DeactivateClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "class_code")
	locale := DetectLocale(w, r)

	req, err := requests.DeactivateClass(r)
	if err != nil {
		h.log.WithError(err).WithField("class_code", code).Error("error deleting place")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != code {
		h.log.Errorf("id mismatch with class_code parameter: %s != %s", req.Data.Id, code)
		ape.RenderErr(w,
			problems.InvalidPointer("data/id", fmt.Errorf("id mismatch with class_code parameter")),
			problems.InvalidParameter("class_code", fmt.Errorf("id mismatch with data/id")),
		)

		return
	}

	res, err := h.domain.Class.Deactivate(r.Context(), code, locale, req.Data.Attributes.ReplacedClassCode)
	if err != nil {
		h.log.WithError(err).Errorf("error deactivating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class with code %s not found", code)))
		case errors.Is(err, errx.ErrorClassDeactivateReplaceSame):
			ape.RenderErr(w, problems.InvalidParameter("data/attributes/replaced_class_code", err))
		case errors.Is(err, errx.ErrorClassDeactivateReplaceInactive):
			ape.RenderErr(w, problems.InvalidParameter("data/attributes/replaced_class_code", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(res))
}
