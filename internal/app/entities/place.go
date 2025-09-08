package entities

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/pagi"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type placeQ interface {
	New() dbx.PlacesQ

	Insert(ctx context.Context, input dbx.Place) error
	Get(ctx context.Context) (dbx.Place, error)
	Select(ctx context.Context) ([]dbx.Place, error)
	Update(ctx context.Context, input dbx.UpdatePlaceParams) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.PlacesQ
	FilterClass(class ...string) dbx.PlacesQ
	FilterStatus(status ...string) dbx.PlacesQ
	FilterOwnership(ownership ...string) dbx.PlacesQ
	FilterCityID(cityID ...uuid.UUID) dbx.PlacesQ
	FilterDistributorID(distributorID ...uuid.UUID) dbx.PlacesQ
	FilterVerified(verified bool) dbx.PlacesQ
	FilterNameLike(name string) dbx.PlacesQ
	FilterAddressLike(address string) dbx.PlacesQ
	FilterWithinRadiusMeters(point orb.Point, radiusM uint64) dbx.PlacesQ
	FilterWithinBBox(minLon, minLat, maxLon, maxLat float64) dbx.PlacesQ
	FilterWithinPolygonWKT(polyWKT string) dbx.PlacesQ
	FilterTimetableBetween(start, end int) dbx.PlacesQ

	WithLocale(locale string) dbx.PlacesQ

	GetWithLocale(ctx context.Context, locale string) (dbx.PlaceWithLocale, error)
	SelectWithLocale(ctx context.Context, locale string) ([]dbx.PlaceWithLocale, error)

	OrderByCreatedAt(ascend bool) dbx.PlacesQ
	OrderByDistance(point orb.Point, ascend bool) dbx.PlacesQ

	Page(limit, offset uint64) dbx.PlacesQ
	Count(ctx context.Context) (uint64, error)
}

type Place struct {
	query   placeQ
	localeQ placeLocaleQ
}

func NewPlace(db *sql.DB) Place {
	return Place{
		query:   dbx.NewPlacesQ(db),
		localeQ: dbx.NewPlaceLocalesQ(db),
	}
}

type CreatePlaceParams struct {
	ID            uuid.UUID
	CityID        uuid.UUID
	DistributorID *uuid.UUID
	Class         string
	Status        string
	Website       *string
	Phone         *string
	Ownership     string
	Point         orb.Point
}

type CreatePlaceLocalParams struct {
	Locale      string
	Name        string
	Address     string
	Description *string
}

func (p Place) Create(
	ctx context.Context,
	params CreatePlaceParams,
	locale CreatePlaceLocalParams,
) (models.PlaceWithLocale, error) {
	now := time.Now().UTC()

	stmt := dbx.Place{
		ID:        params.ID,
		CityID:    params.CityID,
		Class:     params.Class,
		Status:    params.Status,
		Verified:  false,
		Ownership: params.Ownership,
		Point:     params.Point,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if params.DistributorID != nil {
		stmt.DistributorID = uuid.NullUUID{UUID: *params.DistributorID, Valid: true}
	}
	if params.Website != nil {
		stmt.Website = sql.NullString{String: *params.Website, Valid: true}
	}
	if params.Phone != nil {
		stmt.Phone = sql.NullString{String: *params.Phone, Valid: true}
	}

	err := p.query.New().Insert(ctx, stmt)
	if err != nil {
		return models.PlaceWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not create Location, cause %w", err),
		)
	}

	stmtLocale := dbx.PlaceLocale{
		PlaceID: params.ID,
		Locale:  locale.Locale,
		Name:    locale.Name,
		Address: locale.Address,
	}
	err = p.localeQ.Insert(ctx, stmtLocale)
	if err != nil {
		return models.PlaceWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not create Location locale, cause %w", err),
		)
	}

	place := models.Place{
		ID:        params.ID,
		CityID:    params.CityID,
		Class:     params.Class,
		Status:    params.Status,
		Verified:  false,
		Ownership: params.Ownership,
		Point:     params.Point,
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

	paramsLocale := models.LocaleForPlace{
		PlaceID: params.ID,
		Locale:  locale.Locale,
		Name:    locale.Name,
		Address: locale.Address,
	}
	if locale.Description != nil {
		paramsLocale.Description = locale.Description
	}

	return models.PlaceWithLocale{
		Data:   place,
		Locale: paramsLocale,
	}, nil
}

func (p Place) GetPlaceByID(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
) (models.PlaceWithLocale, error) {
	err := constant.IsValidLocaleSupported(locale)
	if err != nil {
		locale = constant.LocaleEN
	}

	place, err := p.query.New().FilterID(placeID).GetWithLocale(ctx, locale)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.PlaceWithLocale{}, errx.ErrorPlaceNotFound.Raise(
				fmt.Errorf("Location with id %s not found, cause %w", placeID, err),
			)
		default:
			return models.PlaceWithLocale{}, errx.ErrorInternal.Raise(
				fmt.Errorf("failed to get Location with id %s: %w", placeID, err),
			)
		}
	}

	return placeWithLocaleModelFromDB(place), nil
}

type SearchPlacesFilter struct {
	Class         []string
	Status        []string
	Ownership     []string
	CityID        []uuid.UUID
	DistributorID []uuid.UUID
	Verified      *bool
	Name          *string
	Address       *string

	Location *SearchPlaceDistanceFilter

	Locale *string
}

type SearchPlaceDistanceFilter struct {
	Point   orb.Point
	RadiusM uint64
}

func (p Place) SearchPlaces(
	ctx context.Context,
	filter SearchPlacesFilter,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.PlaceWithLocale, pagi.Response, error) {
	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1 // +1 чтобы определить наличие next
	offset := (pag.Page - 1) * pag.Size

	query := p.query.New()

	if len(filter.Class) > 0 && filter.Class != nil {
		query = query.FilterClass(filter.Class...)
	}
	if len(filter.Status) > 0 && filter.Status != nil {
		query = query.FilterStatus(filter.Status...)
	}
	if filter.Verified != nil {
		query = query.FilterVerified(*filter.Verified)
	}
	if len(filter.Ownership) > 0 && filter.Ownership != nil {
		query = query.FilterOwnership(filter.Ownership...)
	}
	if len(filter.CityID) > 0 && filter.CityID != nil {
		query = query.FilterCityID(filter.CityID...)
	}
	if len(filter.DistributorID) > 0 && filter.DistributorID != nil {
		query = query.FilterDistributorID(filter.DistributorID...)
	}
	if filter.Name != nil {
		query = query.FilterNameLike(*filter.Name)
	}
	if filter.Address != nil {
		query = query.FilterAddressLike(*filter.Address)
	}
	if filter.Location.RadiusM > 0 && filter.Location != nil {
		query = query.FilterWithinRadiusMeters(filter.Location.Point, filter.Location.RadiusM)
	}

	count, err := p.query.Count(ctx)

	for _, s := range sort {
		switch s.Field {
		case "created_at":
			query = query.OrderByCreatedAt(s.Ascend)
		case "distance":
			if filter.Location != nil {
				query = query.OrderByDistance(filter.Location.Point, s.Ascend)
			}
		}
	}

	loc := constant.LocaleEN
	if filter.Locale != nil {
		err = constant.IsValidLocaleSupported(loc)
		if err == nil {
			loc = *filter.Locale
		}
	}

	rows, err := query.Page(limit, offset).SelectWithLocale(ctx, loc)
	if err != nil {
		return nil, pagi.Response{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to search locos, cause: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	places := make([]models.PlaceWithLocale, 0, len(rows))
	for _, row := range rows {
		places = append(places, placeWithLocaleModelFromDB(row))
	}

	return places, pagi.Response{
		Page:  pag.Page,
		Size:  pag.Size,
		Total: count,
	}, nil
}

type UpdatePlaceParams struct {
	Class     *string
	Status    *string
	Verified  *bool
	Ownership *string
	Point     *orb.Point
	Website   *string
	Phone     *string
}

func (p Place) UpdatePlace(
	ctx context.Context,
	placeID uuid.UUID,
	locale string,
	params UpdatePlaceParams,
) (models.PlaceWithLocale, error) {
	place, err := p.GetPlaceByID(ctx, placeID, locale) //TODO locale
	if err != nil {
		return models.PlaceWithLocale{}, err
	}

	stmt := dbx.UpdatePlaceParams{
		UpdatedAt: time.Now().UTC(),
	}
	if params.Class != nil {
		stmt.Class = params.Class
		place.Data.Class = *params.Class
	}
	if params.Status != nil {
		stmt.Status = params.Status
		place.Data.Status = *params.Status
	}
	if params.Verified != nil {
		stmt.Verified = params.Verified
		place.Data.Verified = *params.Verified
	}
	if params.Ownership != nil {
		err := constant.IsValidPlaceOwnership(*params.Ownership)
		if err != nil {
			return models.PlaceWithLocale{}, err
		}
		stmt.Ownership = params.Ownership
		place.Data.Ownership = *params.Ownership
	}
	if params.Point != nil {
		stmt.Point = params.Point
		place.Data.Point = *params.Point
	}
	if params.Website != nil {
		switch *params.Website {
		case "":
			stmt.Website = &sql.NullString{Valid: false}
			place.Data.Website = nil
		default:
			stmt.Website = &sql.NullString{String: *params.Website, Valid: true}
			place.Data.Website = params.Website
		}
	}
	if params.Phone != nil {
		switch *params.Phone {
		case "":
			stmt.Phone = &sql.NullString{Valid: false}
			place.Data.Phone = nil
		default:
			stmt.Phone = &sql.NullString{String: *params.Phone, Valid: true}
			place.Data.Phone = params.Phone
		}
	}
	err = p.query.New().Update(ctx, stmt)
	if err != nil {
		return models.PlaceWithLocale{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update Location with id %s, cause: %w", placeID, err),
		)
	}

	return place, nil
}

func (p Place) DeletePlaceByID(ctx context.Context, placeID uuid.UUID) error {
	_, err := p.GetPlaceByID(ctx, placeID, constant.LocaleEN)
	if err != nil {
		return err
	}

	err = p.localeQ.New().FilterPlaceID(placeID).Delete(ctx)
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
	Ownership     *string
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
	if filter.Ownership != nil {
		query = query.FilterOwnership(*filter.Ownership)
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

type UpdatePlacesFilter struct {
	Class         []string
	CityID        []uuid.UUID
	DistributorID []uuid.UUID
}

type UpdatePlacesParams struct {
	Class     *string
	Status    *string
	Ownership *string
	Verified  *bool
}

func (p Place) UpdatePlaces(
	ctx context.Context,
	filter UpdatePlacesFilter,
	params UpdatePlaceParams,
) error {
	query := p.query.New()

	if len(filter.Class) > 0 && filter.Class != nil {
		query = query.FilterClass(filter.Class...)
	}
	if len(filter.CityID) > 0 && filter.CityID != nil {
		query = query.FilterCityID(filter.CityID...)
	}
	if len(filter.DistributorID) > 0 && filter.DistributorID != nil {
		query = query.FilterDistributorID(filter.DistributorID...)
	}

	stmt := dbx.UpdatePlaceParams{}
	if params.Class != nil {
		stmt.Class = params.Class
	}
	if params.Status != nil {
		err := constant.IsValidPlaceStatus(*params.Status)
		if err != nil {
			return errx.ErrorPlaceStatusInvalid.Raise(
				fmt.Errorf("invalid place status, cause: %w", err),
			)
		}

		stmt.Status = params.Status
	}
	if params.Ownership != nil {
		err := constant.IsValidPlaceOwnership(*params.Ownership)
		if err != nil {
			return errx.ErrorPlaceOwnershipInvalid.Raise(
				fmt.Errorf("invalid place ownership, cause: %w", err),
			)
		}

		stmt.Ownership = params.Ownership
	}
	if params.Verified != nil {
		stmt.Verified = params.Verified
	}

	err := query.Update(ctx, stmt)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update locos, cause: %w", err),
		)
	}

	return nil
}

func placeWithLocaleModelFromDB(in dbx.PlaceWithLocale) models.PlaceWithLocale {
	out := models.PlaceWithLocale{
		Data: models.Place{
			ID:        in.ID,
			CityID:    in.CityID,
			Class:     in.Class,
			Status:    in.Status,
			Verified:  in.Verified,
			Ownership: in.Ownership,
			Point:     in.Point,
			CreatedAt: in.CreatedAt,
			UpdatedAt: in.UpdatedAt,
		},
		Locale: models.LocaleForPlace{
			PlaceID: in.ID,
		},
	}
	if in.DistributorID.Valid {
		out.Data.DistributorID = &in.DistributorID.UUID
	}
	if in.Website.Valid {
		out.Data.Website = &in.Website.String
	}
	if in.Phone.Valid {
		out.Data.Phone = &in.Phone.String
	}
	if in.Description.Valid {
		out.Locale.Description = &in.Description.String
	}
	if in.Locale != nil {
		out.Locale.Locale = *in.Locale
	}
	if in.Name != nil {
		out.Locale.Name = *in.Name
	}
	if in.Address != nil {
		out.Locale.Address = *in.Address
	}

	return out
}
