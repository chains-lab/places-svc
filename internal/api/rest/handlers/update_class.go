package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/app"
)

func (a Adapter) UpdateClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdateClass(r)
	if err != nil {
		a.log.WithError(err).Error("error updating class")
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

	class, err := a.app.UpdateClass(
		r.Context(),
		req.Data.Id,
		DetectLocale(w, r),
		params,
	)
	if err != nil {
		a.log.WithError(err).Error("error updating class")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(class))
}
