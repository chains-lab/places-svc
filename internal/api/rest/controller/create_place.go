package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/meta"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/paulmach/orb"
)

func (h Service) CreatePlace(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		h.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreatePlace(r)
	if err != nil {
		h.log.WithError(err).Error("error creating place")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	params := place.CreateParams{
		CityID: req.Data.Attributes.CityId,
		Class:  req.Data.Attributes.Class,
		Point: orb.Point{
			req.Data.Attributes.Point.Lon,
			req.Data.Attributes.Point.Lat,
		},
		Locale:      req.Data.Attributes.Locale,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
	}
	if req.Data.Attributes.DistributorId != nil {
		params.DistributorID = req.Data.Attributes.DistributorId
	}
	if req.Data.Attributes.Phone != nil {
		params.Phone = req.Data.Attributes.Phone
	}
	if req.Data.Attributes.Website != nil {
		params.Website = req.Data.Attributes.Website
	}

	res, err := h.domain.Place.Create(r.Context(), params)
	if err != nil {
		h.log.WithError(err).Error("error creating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class with code %s not found", params.Class)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	h.log.Infof("created place with id %s by user %s", res.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Place(res))
}
