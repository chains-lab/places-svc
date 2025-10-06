package place

import (
	"context"
	"fmt"
	"time"

	"github.com/chains-lab/places-svc/internal/domain/enum"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type CreateParams struct {
	CityID        uuid.UUID
	DistributorID *uuid.UUID
	Class         string
	Address       string
	Status        string
	Phone         *string
	Website       *string
	Point         orb.Point

	Locale      string
	Name        string
	Description string
}

func (s Service) Create(
	ctx context.Context,
	params CreateParams,
) (models.Place, error) {
	now := time.Now().UTC()

	placeID := uuid.New()

	place := models.Place{
		ID:        placeID,
		CityID:    params.CityID,
		Class:     params.Class,
		Status:    enum.PlaceStatusActive,
		Verified:  false,
		Point:     params.Point,
		Address:   params.Address,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if params.DistributorID != nil {
		place.DistributorID = params.DistributorID
	}
	if params.Website != nil {
		place.Website = params.Website
	}
	if params.Phone != nil {
		place.Phone = params.Phone
	}

	exist, err := s.db.ClassIsExistByCode(ctx, params.Class)
	if err != nil {
		return models.Place{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not verify class existence, cause %w", err),
		)
	}

	if !exist {
		return models.Place{}, errx.ErrorClassNotFound.Raise(
			fmt.Errorf("class %s not found", params.Class),
		)
	}

	var addr string
	if err = s.db.Transaction(ctx, func(ctx context.Context) error {
		err = s.db.CreatePlace(ctx, place.Details())
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not create place, cause %w", err),
			)
		}

		stmtLocale := models.PlaceLocale{
			PlaceID:     placeID,
			Locale:      params.Locale,
			Name:        params.Name,
			Description: params.Description,
		}
		err = s.db.CreatePlaceLocale(ctx, stmtLocale)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not create place locale, cause %w", err),
			)
		}

		addr, err = s.geo.Guess(ctx, params.Point)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not guess address for Location, cause %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Place{}, err
	}

	res := models.Place{
		ID:          placeID,
		CityID:      params.CityID,
		Class:       params.Class,
		Status:      params.Status,
		Verified:    false,
		Point:       params.Point,
		CreatedAt:   now,
		UpdatedAt:   now,
		Address:     addr,
		Locale:      params.Locale,
		Name:        params.Name,
		Description: params.Description,
		Timetable:   models.Timetable{},
	}
	if params.DistributorID != nil {
		res.DistributorID = params.DistributorID
	}
	if params.Website != nil {
		res.Website = params.Website
	}
	if params.Phone != nil {
		res.Phone = params.Phone
	}

	return res, nil
}
