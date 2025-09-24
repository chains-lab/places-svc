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
	"github.com/chains-lab/places-svc/internal/domain/services/place"
)

func (h Service) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdatePlace(r)
	if err != nil {
		h.log.WithError(err).Error("error updating place")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	params := place.UpdateParams{}

	if req.Data.Attributes.Phone != nil {
		params.Phone = req.Data.Attributes.Phone
	}
	if req.Data.Attributes.Website != nil {
		params.Website = req.Data.Attributes.Website
	}
	if req.Data.Attributes.Class != nil {
		params.Class = req.Data.Attributes.Class
	}

	res, err := h.domain.Place.Update(
		r.Context(),
		req.Data.Id,
		DetectLocale(w, r),
		params,
	)
	if err != nil {
		h.log.WithError(err).Error("error updating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class %s not found", *params.Class)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
