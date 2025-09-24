package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chains-lab/places-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func CreateClass(r *http.Request) (req resources.CreateClass, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required, validation.RuneLength(0, 32)),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In(resources.ClassType)),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),

		"data/attributes/name": validation.Validate(
			req.Data.Attributes.Name, validation.Required, validation.RuneLength(1, 32)),
		"data/attributes/icon": validation.Validate(
			req.Data.Attributes.Icon, validation.RuneLength(1, 32)),
		"data/attributes/parent": validation.Validate(
			req.Data.Attributes.Parent, validation.RuneLength(1, 32)),
	}

	if req.Data.Attributes.Parent != nil && *req.Data.Attributes.Parent == req.Data.Id {
		errs["data/attributes/parent"] = fmt.Errorf("parent must not equal id")
	}

	return req, errs.Filter()
}
