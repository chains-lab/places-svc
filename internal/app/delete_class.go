package app

import (
	"context"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (a App) DeleteClass(
	ctx context.Context,
	code string,
) error {
	c, err := a.classificator.Get(ctx, code, enum.DefaultLocale)
	if err != nil {
		return err
	}

	if c.Data.Status == enum.PlaceClassStatusesActive {
		return errx.ErrorCannotDeleteActiveClass.Raise(err)
	}

	places, _, err := a.place.List(ctx, enum.DefaultLocale,
		place.FilterListParams{Classes: []string{c.Data.Code}}, pagi.Request{}, nil)

	if len(places) != 0 {
		return errx.ErrorCantDeleteClassWithPlaces.Raise(err)
	}

	return a.classificator.DeleteClass(ctx, code)
}
