package requests

import (
	"encoding/json"
	"net/http"

	"github.com/chains-lab/places-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SetPlaceLocale(r *http.Request) (req resources.SetPlaceLocale, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(&req.Data, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.PlaceLocateType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}
}
