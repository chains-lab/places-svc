package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/rest/requests"
	"github.com/chains-lab/places-svc/internal/rest/responses"
)

func (s Service) UpdatePlaceStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdatePlaceStatus(r)
	if err != nil {
		s.log.WithError(err).Error("error parsing update place status request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.place.UpdateStatus(r.Context(), req.Data.Id, DetectLocale(w, r), req.Data.Attributes.Status)
	if err != nil {
		s.log.WithError(err).WithField("place_id", req.Data.Id).Error("error activating place")
		switch {
		case errors.Is(err, errx.ErrorCannotSetStatusBlocked):
			ape.RenderErr(w, problems.Conflict("cannot set status to 'blocked'"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
