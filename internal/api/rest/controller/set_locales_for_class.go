package controller

import (
	"errors"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/domain/modules/class"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h Service) SetLocalesForClass(w http.ResponseWriter, r *http.Request) {
	classCode, err := uuid.Parse(chi.URLParam(r, "class_code"))
	if err != nil {
		h.log.WithError(err).Error("invalid class_code")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.SetLocalesForClass(r)
	if err != nil {
		h.log.WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	locales := make([]class.SetLocaleParams, 0, len(req.Data.Attributes.Locales))
	for _, attr := range req.Data.Attributes.Locales {
		locales = append(locales, class.SetLocaleParams{
			Locale: attr.Locale,
			Name:   attr.Name,
		})
	}

	err = h.domain.Class.SetLocales(r.Context(), classCode.String(), locales...)
	if err != nil {
		h.log.WithError(err).Error("failed to set class locales")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound("class not found"))
		case errors.Is(err, errx.ErrorNedAtLeastOneLocaleForClass):
			ape.RenderErr(w, problems.InvalidParameter("locales", err))
		case errors.Is(err, errx.ErrorInvalidLocale):
			ape.RenderErr(w, problems.InvalidParameter("locales", err))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
