package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
)

func (h Service) DeactivateClass(w http.ResponseWriter, r *http.Request) {
	locale := DetectLocale(w, r)

	req, err := requests.DeactivateClass(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := h.domain.Class.Deactivate(r.Context(), req.Data.Id, locale, req.Data.Attributes.ReplacedClassCode)
	if err != nil {
		h.log.WithError(err).Errorf("error deactivating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class with code %s not found", req.Data.Id)))
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
