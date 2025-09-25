package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) DeactivatePlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		s.log.WithError(err).Error("activate place")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("failed to parse place_id: %w", err),
		})...)

		return
	}

	locale := DetectLocale(w, r)

	res, err := s.domain.place.Deactivate(r.Context(), placeID, locale)
	if err != nil {
		s.log.WithError(err).Errorf("error deactivating place with id %s", placeID)
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound(fmt.Sprintf("place with id %s not found", placeID)))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Place(res))
}
