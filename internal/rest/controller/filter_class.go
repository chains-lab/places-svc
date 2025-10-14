package controller

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/domain/services/class"
	"github.com/chains-lab/places-svc/internal/rest/responses"
	"github.com/chains-lab/restkit/pagi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (s Service) FilterClass(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var filters class.FilterParams

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
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"query": fmt.Errorf("invalid parent_cycle value: %s", parentCycle),
			})...)
			return
		}
	}

	pag, size := pagi.GetPagination(r)

	classes, err := s.domain.class.Filter(r.Context(), filters, pag, size)
	if err != nil {
		s.log.WithError(err).Error("failed to list classes")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.ClassesCollection(classes))
}
