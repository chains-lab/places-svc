package handlers

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) SetLocalesForClass(w http.ResponseWriter, r *http.Request) {
	classCode, err := uuid.Parse(chi.URLParam(r, "class_code"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid class_code")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	req, err := requests.SetLocalesForClass(r)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	locales := make([]app.SetClassLocaleParams, 0, len(req.Data.Attributes.Locales))
	for _, attr := range req.Data.Attributes.Locales {
		locales = append(locales, app.SetClassLocaleParams{
			Locale: attr.Locale,
			Name:   attr.Name,
		})
	}

	err = a.app.SetLocalesForClass(r.Context(), classCode.String(), locales...)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to set class locales")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
