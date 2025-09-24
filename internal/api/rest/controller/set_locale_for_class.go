package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) SetLocaleForClass(w http.ResponseWriter, r *http.Request) {
	req, err := requests.SetLocaleForClass(r)
	if err != nil {
		s.log.WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	loc := class.SetLocaleParams{
		Locale: req.Data.Attributes.Locale,
		Name:   req.Data.Attributes.Name,
	}

	err = s.domain.Class.SetLocale(r.Context(), req.Data.Id, loc)
	if err != nil {
		s.log.WithError(err).Error("failed to set class locales")
		switch {
		case errors.Is(err, errx.ErrorClassNotFound):
			ape.RenderErr(w, problems.NotFound("class not found"))
		case errors.Is(err, errx.ErrorCannotDeleteDefaultLocaleForClass):
			ape.RenderErr(w, problems.Conflict("default locale can not be changed"))
		case errors.Is(err, errx.ErrorInvalidLocale):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"locales": fmt.Errorf("invalid locale: %s", req.Data.Attributes.Locale),
			})...)
		case errors.Is(err, errx.ErrorClassNameAlreadyTaken):
			ape.RenderErr(w, problems.Conflict("class with the same name already exists in this locale"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
