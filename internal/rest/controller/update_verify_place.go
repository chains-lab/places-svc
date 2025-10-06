package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/rest/responses"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

func (s Service) UpdateVerifyPlace(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		s.log.WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("failed to parse place_id: %w", err),
		})...)

		return
	}

	var value bool
	q := r.URL.Query()
	if verified := strings.TrimSpace(q.Get("verified")); verified != "" {
		switch verified {
		case "true":
			value = true
		case "false":
			value = false
		default:
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("invalid verified value: %s", verified),
			})...)
			return
		}
	}

	res, err := s.domain.place.Verify(r.Context(), placeID, DetectLocale(w, r), value)
	if err != nil {
		s.log.WithError(err).Error("failed to verify place")
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
