package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/resources"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func SetLocaleForClass(r *http.Request) (req resources.SetLocaleForClass, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required, is.UUID),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.PlaceLocaleType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),

		"data/attributes/locale": validation.Validate(
			req.Data.Attributes.Locale, validation.Required, validation.In(enum.GetAllLocales()),
		),
		"data/attributes/name": validation.Validate(
			req.Data.Attributes.Name, validation.Min(1), validation.Max(32),
		),
	}

	if chi.URLParam(r, "class_code") != req.Data.Id {
		errs["data/id"] = fmt.Errorf("query class_code param and body data/id do not match")
	}

	return req, errs.Filter()
}
