package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/domain/modules/class"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (h Service) ListClass(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var filters class.FilterListParams

	if parent := q.Get("parent"); parent != "" {
		filters.Parent = &parent
	}

	if status := q.Get("status"); status != "" {
		filters.Status = &status
	}

	if parentCycle := q.Get("parent_cycle"); parentCycle != "" {
		if parentCycle == "true" {
			filters.ParentCycle = true
		} else if parentCycle == "false" {
			filters.ParentCycle = false
		} else {
			ape.RenderErr(w, problems.InvalidParameter(
				"verified",
				fmt.Errorf("invalid boolean value: %s", parentCycle)),
			)
			return
		}
	}

	pag, _ := pagi.GetPagination(r)
	locale := DetectLocale(w, r)

	classes, pagResp, err := h.domain.Class.List(r.Context(), locale, filters, pag)
	if err != nil {
		h.log.WithError(err).Error("failed to list classes")
		switch {
		case errors.Is(err, errx.ErrorClassStatusInvalid):
			ape.RenderErr(w, problems.InvalidParameter(
				"status",
				fmt.Errorf("invalid status value: %s", *filters.Status)),
			)
		case errors.Is(err, errx.ErrorParentClassNotFound):
			ape.RenderErr(w, problems.NotFound(
				fmt.Sprintf("parent class %s not found", *filters.Parent)),
			)
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.ClassesCollection(classes, pagResp))
}
