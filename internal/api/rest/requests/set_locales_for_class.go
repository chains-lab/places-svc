package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/resources"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func SetLocalesForClass(r *http.Request) (req resources.SetLocalesForClass, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required, is.UUID),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.PlaceLocaleType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	for i, loc := range req.Data.Attributes.Locales {
		key := fmt.Sprintf("data/attributes/locales/%d/locale", i)
		errs[key] = validation.Validate(
			loc.Locale,
			validation.Required,
			validation.In(enum.GetAllLocales()),
		)

		nameKey := fmt.Sprintf("data/attributes/locales/%d/name", i)
		errs[nameKey] = validation.Validate(loc.Name, validation.RuneLength(0, 32))
	}

	if chi.URLParam(r, "class_code") != req.Data.Id {
		errs["data/id"] = fmt.Errorf("query class_code param and body data/id do not match")
	}

	return req, errs.Filter()
}
