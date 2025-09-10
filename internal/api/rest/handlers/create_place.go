package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/meta"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

func (a Adapter) CreatePlace(w http.ResponseWriter, r *http.Request) {
	initiator, err := meta.User(r.Context())
	if err != nil {
		a.log.WithError(err).Error("failed to get user from context")
		ape.RenderErr(w, problems.Unauthorized("failed to get user from context"))

		return
	}

	req, err := requests.CreatePlace(r)
	if err != nil {
		a.log.WithError(err).Error("error creating place")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	cityID, err := uuid.Parse(req.Data.Attributes.CityId)
	if err != nil {
		a.log.WithError(err).Error("invalid city_id")
		ape.RenderErr(w, problems.InvalidPointer("data/attributes/city_id", err))

		return
	}

	params := app.CreatePlaceParams{
		CityID: cityID,
		Class:  req.Data.Attributes.Class,
		Point: orb.Point{
			req.Data.Attributes.Point.Lon,
			req.Data.Attributes.Point.Lat,
		},
	}
	if req.Data.Attributes.DistributorId != nil {
		distributorID, err := uuid.Parse(*req.Data.Attributes.DistributorId)
		if err != nil {
			a.log.WithError(err).Error("invalid distributor_id")
			ape.RenderErr(w, problems.InvalidPointer("data/attributes/distributor_id", err))

			return
		}

		//TODO make a post for check user is distributor or not

		params.DistributorID = &distributorID
	}
	if req.Data.Attributes.Phone != nil {
		params.Phone = req.Data.Attributes.Phone
	}
	if req.Data.Attributes.Website != nil {
		params.Website = req.Data.Attributes.Website
	}

	locate := app.CreatePlaceLocalParams{
		Locale:      req.Data.Attributes.Locale,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
	}

	place, err := a.app.CreatePlace(r.Context(), params, locate)
	if err != nil {
		a.log.WithError(err).Error("error creating place")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	a.log.Infof("created place with id %s by user %s", place.Place.ID, initiator.ID)

	ape.Render(w, http.StatusCreated, responses.Place(place))
}
