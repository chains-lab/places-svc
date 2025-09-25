package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) ListLocalesForPlace(w http.ResponseWriter, r *http.Request) {
	pag, size := pagi.GetPagination(r)

	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		s.log.WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("failed to parse place_id: %w", err),
		})...)

		return
	}

	locales, err := s.domain.place.ListLocales(r.Context(), placeID, pag, size)
	if err != nil {
		s.log.WithError(err).Error("failed to get place locales")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.PlaceLocalesCollection(locales))
}
