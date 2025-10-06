package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) UpdatePlaceStatus(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		s.log.WithError(err).Error("activate place")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("failed to parse place_id: %w", err),
		})...)

		return
	}

	res, err := s.domain.place.UpdateStatus(r.Context(), placeID, DetectLocale(w, r), chi.URLParam(r, "status"))
	if err != nil {
		s.log.WithError(err).WithField("place_id", placeID).Error("error activating place")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound("class not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
