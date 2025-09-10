package app

import (
	"context"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/entities/place"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
)

func (a App) DeleteClass(
	ctx context.Context,
	code string,
) error {
	c, err := a.classificator.Get(ctx, code, constant.DefaultLocale)
	if err != nil {
		return err
	}

	if c.Data.Status == constant.PlaceClassStatusesActive {
		return errx.ErrorCannotDeleteActiveClass.Raise(err)
	}

	places, _, err := a.place.List(ctx, constant.DefaultLocale,
		place.FilterListParams{Class: []string{c.Data.Code}}, pagi.Request{}, nil)

	if len(places) != 0 {
		return errx.ErrorClassHasPlaces.Raise(err)
	}

	return a.classificator.DeleteClass(ctx, code)
}
