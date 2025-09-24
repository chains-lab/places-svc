package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

func (m Module) DeleteOne(ctx context.Context, placeID uuid.UUID) error {
	place, err := m.Get(ctx, placeID, enum.LocaleEN)
	if err != nil {
		return err
	}

	if place.Status != enum.PlaceStatusInactive {
		return errx.ErrorPlaceForDeleteMustBeInactive.Raise(
			fmt.Errorf("place %m is not in inactive status", place.ID.String()),
		)
	}

	err = m.db.PlaceLocales().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to delete Location locale with id %m, cause: %w", placeID, err),
		)
	}

	return nil
}

type DeleteFilter struct {
	Class         *string
	Status        *string
	CityID        *uuid.UUID
	DistributorID *uuid.UUID
	Verified      *bool
	Name          *string
	Address       *string
}

func (m Module) DeleteMany(ctx context.Context, filter DeleteFilter) error {
	query := m.db.Places()

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
