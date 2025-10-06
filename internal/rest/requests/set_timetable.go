package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/places-svc/resources"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func SetTimetable(r *http.Request) (req resources.SetPlaceTimetable, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required, is.UUID),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.PlaceType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "place_id") != req.Data.Id.String() {
		errs["data/id"] = fmt.Errorf("query place_id param and body data/id do not match")
	}

	return req, errs.Filter()
}
