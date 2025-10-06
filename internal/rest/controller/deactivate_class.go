package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/rest/requests"
	"github.com/chains-lab/places-svc/internal/rest/responses"
)

func (s Service) DeactivateClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.DeactivateClass(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.class.Deactivate(r.Context(), req.Data.Id, req.Data.Attributes.ReplacedClassCode)
	if err != nil {
		s.log.WithError(err).Errorf("error deactivating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class with code %s not found", req.Data.Id)))
		case errors.Is(err, errx.ErrorClassDeactivateReplaceSame):
			ape.RenderErr(w, problems.Conflict(fmt.Sprintf("class cannot replace itself")))
		case errors.Is(err, errx.ErrorClassDeactivateReplaceInactive):
			ape.RenderErr(w, problems.Conflict(fmt.Sprintf("class cannot replace an inactive class")))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(res))
}
