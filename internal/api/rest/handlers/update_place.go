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
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdatePlace(r)
	if err != nil {
		a.log.WithError(err).Error("error updating place")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	placeId, err := uuid.Parse(req.Data.Id)
	if err != nil {
		a.log.WithError(err).Error("invalid place id")
		ape.RenderErr(w, problems.InvalidParameter("data/id", err))

		return
	}

	if placeId.String() != chi.URLParam(r, "place_id") {
		a.log.Error("place id in body does not match place id in URL")
		ape.RenderErr(w,
			problems.InvalidParameter("data/id", fmt.Errorf("query param and body do not match")),
			problems.InvalidParameter("place_id", fmt.Errorf("query param and body do not match")),
		)

		return
	}

	params := app.UpdatePlaceParams{}

	if req.Data.Attributes.Phone != nil {
		params.Phone = req.Data.Attributes.Phone
	}
	if req.Data.Attributes.Website != nil {
		params.Website = req.Data.Attributes.Website
	}
	if req.Data.Attributes.Class != nil {
		params.Class = req.Data.Attributes.Class
	}

	place, err := a.app.UpdatePlace(
		r.Context(),
		placeId,
		DetectLocale(w, r),
		params,
	)
	if err != nil {
		a.log.WithError(err).Error("error updating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("class %s not found", *params.Class)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(place))
}
