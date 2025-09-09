package handlers

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/go-chi/chi/v5"
)

func (a Adapter) DeactivateClass(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "class_code")

	locale := DetectLocale(w, r)

	req, err := requests.DeactivateClass(r)
	if err != nil {
		a.log.WithError(err).WithField("class_code", code).Error("error deleting place")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	if req.Data.Id != code {
		ape.RenderErr(w,
			problems.InvalidPointer("data/id", fmt.Errorf("id mismatch with class_code parameter")),
			problems.InvalidParameter("class_code", fmt.Errorf("id mismatch with data/id")),
		)

		return
	}

	class, err := a.app.DeactivateClass(r.Context(), code, locale, req.Data.Attributes.ReplacedClassCode)
	if err != nil {
		a.log.WithError(err).WithField("class_code", code).Error("error deleting place")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.Class(class))
}
