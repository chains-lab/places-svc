package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/services/plocale"
	"github.com/chains-lab/places-svc/internal/rest/requests"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) SetLocalesForPlace(w http.ResponseWriter, r *http.Request) {
	req, err := requests.SetLocalesForPlace(r)
	if err != nil {
		s.log.WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	locales := make([]plocale.SetParams, 0, len(req.Data.Attributes.Locales))
	for _, attr := range req.Data.Attributes.Locales {
		locales = append(locales, plocale.SetParams{
			Locale:      attr.Locale,
			Name:        attr.Name,
			Description: attr.Description,
		})
	}

	err = s.domain.plocale.SetForPlace(r.Context(), req.Data.Id, locales...)
	if err != nil {
		s.log.WithError(err).Error("failed to set place locales")
		switch {
		case errors.Is(err, errx.ErrorPlaceNotFound):
			ape.RenderErr(w, problems.NotFound("place not found"))
		case errors.Is(err, errx.ErrorNeedAtLeastOneLocaleForPlace):
			ape.RenderErr(w, problems.Conflict("place must have at least one locale"))
		case errors.Is(err, errx.ErrorInvalidLocale):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"locales": fmt.Errorf("invalid locale"),
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
