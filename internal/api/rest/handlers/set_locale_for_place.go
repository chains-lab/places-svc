package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/places-svc/internal/api/rest/requests"
	"github.com/chains-lab/places-svc/internal/app"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (a Adapter) SetPlaceLocale(w http.ResponseWriter, r *http.Request) {
	placeID, err := uuid.Parse(chi.URLParam(r, "place_id"))
	if err != nil {
		a.Log(r).WithError(err).Error("invalid place_id")
		ape.RenderErr(w, problems.InvalidParameter("place_id", err))

		return
	}

	req, err := requests.SetPlaceLocale(r)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid request body")
		ape.RenderErr(w, problems.BadRequest(err)...)

		return
	}

	s := strings.Split(req.Data.Id, ":")
	if len(s) != 2 {
		a.Log(r).WithError(err).Error("invalid id format")
		ape.RenderErr(w, problems.InvalidParameter("id", fmt.Errorf("invalid id: %s", req.Data.Id)))

		return
	}

	if placeID.String() != s[0] {
		a.Log(r).WithError(err).Error("place_id in URL does not match place_id in body")
		ape.RenderErr(w,
			problems.InvalidParameter("id", fmt.Errorf("place_id in URL does not match place_id in body")),
			problems.InvalidParameter("place_id", fmt.Errorf("place_id in URL does not match place_id in body")),
		)

		return
	}

	l := strings.TrimSpace(s[1])
	err = constant.IsValidLocaleSupported(l)
	if err != nil {
		a.Log(r).WithError(err).Error("invalid locale in id")
		ape.RenderErr(w, problems.InvalidPointer("data/id", fmt.Errorf("invalid locale in  id: %s", l)))

		return
	}

	locale := app.CreatePlaceLocalParams{
		Locale:      l,
		Name:        req.Data.Attributes.Name,
		Description: req.Data.Attributes.Description,
	}
	err = a.app.AddPlaceLocales(r.Context(), placeID, locale)
	if err != nil {
		a.Log(r).WithError(err).Error("failed to set place locale")
		ape.RenderErr(w, problems.InternalError())

		return
	}
}
