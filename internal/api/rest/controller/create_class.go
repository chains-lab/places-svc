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
	"github.com/chains-lab/places-svc/internal/domain/services/class"
)

func (s Service) CreateClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.CreateClass(r)
	if err != nil {
		s.log.WithError(err).Error("error creating class")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	params := class.CreateParams{
		Code: req.Data.Id,
		Icon: req.Data.Attributes.Icon,
		Name: req.Data.Attributes.Name,
	}
	if req.Data.Attributes.Parent != nil {
		params.Parent = req.Data.Attributes.Parent
	}

	res, err := s.domain.class.Create(r.Context(), params)
	if err != nil {
		s.log.WithError(err).Error("error creating class")
		switch {
		case errors.Is(err, errx.ErrorClassCodeAlreadyTaken):
			ape.RenderErr(w, problems.Conflict(
				fmt.Sprintf("class %s already exists", req.Data.Id)),
			)
		case errors.Is(err, errx.ErrorParentClassNotFound):
			ape.RenderErr(w, problems.NotFound(
				fmt.Sprintf("parent class %s not found", *req.Data.Attributes.Parent)),
			)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusCreated, responses.Class(res))
}
