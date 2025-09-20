package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (h Handler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateClass(r)
	if err != nil {
		h.log.WithError(err).Error("error updating class")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	params := app.UpdateClassParams{}

	if req.Data.Attributes.Parent != nil {
		params.Parent = req.Data.Attributes.Parent
	}
	if req.Data.Attributes.Icon != nil {
		params.Icon = req.Data.Attributes.Icon
	}

	class, err := h.app.UpdateClass(
		r.Context(),
		req.Data.Id,
		DetectLocale(w, r),
		params,
	)
	if err != nil {
		h.log.WithError(err).Error("error updating class")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class %s not found", req.Data.Id)))
		case errors.Is(err, errx.ErrorClassParentCycle):
			ape.RenderErr(w, problems.PreconditionFailed(
				fmt.Sprintf("parent cycle detected for class with code %s", req.Data.Id)))
		case errors.Is(err, errx.ErrorClassParentEqualCode):
			ape.RenderErr(w, problems.PreconditionFailed(
				fmt.Sprintf("parent equal code for class with code %s", req.Data.Id)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(class))
}
