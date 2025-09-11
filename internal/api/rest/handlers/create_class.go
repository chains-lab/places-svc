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

func (a Adapter) CreateClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.CreateClass(r)
	if err != nil {
		a.log.WithError(err).Error("error creating class")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	params := app.CreateClassParams{
		Code: req.Data.Id,
		Icon: req.Data.Attributes.Icon,
		Name: req.Data.Attributes.Name,
	}
	if req.Data.Attributes.Parent != nil {
		params.Parent = req.Data.Attributes.Parent
	}

	class, err := a.app.CreateClass(r.Context(), params)
	if err != nil {
		a.log.WithError(err).Error("error creating class")
		switch {
		case errors.Is(err, errx.ErrorClassCodeAlreadyTaken):
			ape.RenderErr(w, problems.Conflict(
				fmt.Sprintf("class %s already exists", req.Data.Attributes.Name)),
			)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Class(class))
}
