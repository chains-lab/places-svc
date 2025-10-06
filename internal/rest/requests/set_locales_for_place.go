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

func SetLocalesForPlace(r *http.Request) (req resources.SetLocalesForPlace, err error) {
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
		errs[nameKey] = validation.Validate(loc.Name, validation.RuneLength(0, 128))

		descriptionKey := fmt.Sprintf("data/attributes/locales/%d/description", i)
		errs[descriptionKey] = validation.Validate(loc.Description, validation.RuneLength(0, 1024))
	}

	if chi.URLParam(r, "place_id") != req.Data.Id.String() {
		errs["data/id"] = fmt.Errorf("query place_id param and body data/id do not match")
	}

	return req, errs.Filter()
}
