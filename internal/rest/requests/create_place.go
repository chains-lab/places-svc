package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func newDecodeError(what string, err error) error {
	return validation.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}

func CreatePlace(r *http.Request) (req resources.CreatePlace, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.PlaceType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),

		"data/attributes/locale": validation.Validate(
			req.Data.Attributes.Locale, validation.Required, validation.In(enum.GetAllLocales())),
		"data/attributes/name": validation.Validate(
			req.Data.Attributes.Name, validation.Required, validation.Length(1, 255)),
		"data/attributes/description": validation.Validate(
			req.Data.Attributes.Description, validation.Required, validation.Length(0, 1024)),
		"data/attributes/website": validation.Validate(
			req.Data.Attributes.Website, validation.Length(0, 255)),
		"data/attributes/phone": validation.Validate(
			req.Data.Attributes.Phone, validation.Length(0, 32)),
	}

	return req, errs.Filter()
}
