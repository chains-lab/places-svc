package handlers

import (
	"fmt"
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/api/rest/responses"
	"github.com/chains-lab/places-svc/internal/app"
)

func (a Adapter) ListClass(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var filters app.FilterClassesParams

	if parent := q.Get("parent"); parent != "" {
		filters.Parent = &parent
	}

	if status := q.Get("status"); status != "" {
		filters.Status = &status
	}

	if name := q.Get("name"); name != "" {
		filters.Name = &name
	}

	if parentCycle := q.Get("parent_cycle"); parentCycle != "" {
		if parentCycle == "true" {
			t := true
			filters.ParentCycle = &t
		} else if parentCycle == "false" {
			f := false
			filters.ParentCycle = &f
		} else {
			ape.RenderErr(w, problems.InvalidParameter(
				"verified",
				fmt.Errorf("invalid boolean value: %s", parentCycle)),
			)
			return
		}
	}

	pag, _ := pagi.GetPagination(r)

	classes, pagResp, err := a.app.ListClasses(r.Context(), filters, pag)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to list classes")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}

		return
	}

	ape.Render(w, http.StatusOK, responses.ClassesCollection(classes, pagResp))
}
