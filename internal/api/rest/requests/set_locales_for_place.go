package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/places-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SetLocalesForPlace(r *http.Request) (req resources.SetLocalesForPlace, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(&req.Data, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.PlaceLocaleType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	return req, errs.Filter()
}
