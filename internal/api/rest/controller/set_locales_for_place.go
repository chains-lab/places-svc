package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
)

func (h Service) SetLocalesForPlace(w http.ResponseWriter, r *http.Request) {
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

	err = h.domain.Place.SetLocales(r.Context(), req.Data.Id, locales...)
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
