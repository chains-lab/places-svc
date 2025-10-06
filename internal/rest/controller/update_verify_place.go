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

func (s Service) UpdateVerifiedPlace(w http.ResponseWriter, r *http.Request) {
	req, err := requests.UpdatePlaceVerified(r)
	if err != nil {
		s.log.WithError(err).Error("error parsing update place verified request")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	res, err := s.domain.place.Verify(r.Context(), req.Data.Id, DetectLocale(w, r), req.Data.Attributes.Verified)
	if err != nil {
		s.log.WithError(err).Error("failed to verify place")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
