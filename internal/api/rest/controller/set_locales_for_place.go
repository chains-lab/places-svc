package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/domain/modules/place"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h Service) SetLocalesForPlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		h.log.WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	req, err := requests.SetLocalesForPlace(r)
	if err != nil {
		h.log.WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	locales := make([]place.SetLocaleParams, 0, len(req.Data.Attributes.Locales))
	for _, attr := range req.Data.Attributes.Locales {
		locales = append(locales, place.SetLocaleParams{
			Locale:      attr.Locale,
			Name:        attr.Name,
			Description: attr.Description,
		})
	}

	err = h.domain.Place.SetLocales(r.Context(), placeID, locales...)
	if err != nil {
		h.log.WithError(err).Error("failed to set place locales")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorNeedAtLeastOneLocaleForPlace):
			ape.RenderErr(w, problems.InvalidParameter("locales", err))
		case errors.Is(err, errx.ErrorInvalidLocale):
			ape.RenderErr(w, problems.InvalidParameter("locales", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
