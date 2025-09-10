package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

func (p Place) DeletePlace(ctx context.Context, placeID uuid.UUID) error {
	place, err := p.Get(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return err
	}

	if place.Place.Status != constant.PlaceStatusInactive {
		return errx.ErrorPlaceForDeleteMustBeInactive.Raise(
			fmt.Errorf("place %s is not in inactive status", place.Place.ID.String()),
		)
	}

	err = p.locale.New().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete Location locale with id %s, cause: %w", placeID, err),
		)
	}

	return nil
}

type DeletePlacesFilter struct {
	Class         *string
	Status        *string
	CityID        *uuid.UUID
	DistributorID *uuid.UUID
	Verified      *bool
	Name          *string
	Address       *string
}

func (p Place) DeletePlaces(ctx context.Context, filter DeletePlacesFilter) error {
	query := p.query.New()

	if filter.Class != nil {
		query = query.FilterClass(*filter.Class)
	}
	if filter.Status != nil {
		query = query.FilterStatus(*filter.Status)
	}
	if filter.Verified != nil {
		query = query.FilterVerified(*filter.Verified)
	}
	if filter.CityID != nil {
		query = query.FilterCityID(*filter.CityID)
	}
	if filter.DistributorID != nil {
		query = query.FilterDistributorID(*filter.DistributorID)
	}
	if filter.Name != nil {
		query = query.FilterNameLike(*filter.Name)
	}
	if filter.Address != nil {
		query = query.FilterAddressLike(*filter.Address)
	}

	err := query.Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete locos, cause: %w", err),
		)
	}

	return nil
}
