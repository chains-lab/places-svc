package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/chains-lab/places-svc/internal/data/pgdb"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/chains-lab/places-svc/internal/domain/services/place"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (d Database) CreatePlace(ctx context.Context, input models.PlaceDetails) error {
	return d.sql.places.Insert(ctx, placeModelToSchema(input))
}

func (d Database) GetPlaceByID(ctx context.Context, placeID uuid.UUID, locale string) (models.Place, error) {
	schema, err := d.sql.places.New().FilterID(placeID).GetWithDetails(ctx, locale)
	switch {
	case err == sql.ErrNoRows:
		return models.Place{}, nil
	case err != nil:
		return models.Place{}, err
	}

	return placeSchemaToModel(schema), nil
}

func (d Database) FilterPlaces(
	ctx context.Context,
	locale string,
	filter place.FilterParams,
	sort place.SortParams,
	page, size uint64,
) (models.PlacesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.places.New()

	if filter.Classes != nil && len(filter.Classes) > 0 {
		query = query.FilterClass(filter.Classes...)
	}
	if filter.Statuses != nil && len(filter.Statuses) > 0 {
		query = query.FilterStatus(filter.Statuses...)
	}
	if filter.CityID != nil {
		query = query.FilterCityID(*filter.CityID)
	}
	if filter.DistributorID != nil {
		query = query.FilterDistributorID(*filter.DistributorID)
	}
	if filter.Verified != nil {
		query = query.FilterVerified(*filter.Verified)
	}
	if filter.Name != nil {
		query = query.FilterNameLike(*filter.Name)
	}
	if filter.Address != nil {
		query = query.FilterAddressLike(*filter.Address)
	}
	if filter.Time != nil {
		query = query.FilterTimetableBetween(filter.Time.From.ToNumberMinutes(), filter.Time.To.ToNumberMinutes())
	}
	if filter.Location != nil {
		query = query.FilterWithinRadiusMeters(filter.Location.Point, filter.Location.RadiusM)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.PlacesCollection{}, err
	}

	query = query.Page(limit, offset)

	if sort.ByCreatedAt != nil {
		query = query.OrderByCreatedAt(*sort.ByCreatedAt)
	}
	if sort.ByDistance != nil && filter.Location != nil {
		query = query.OrderByDistance(filter.Location.Point, *sort.ByDistance)
	}

	rows, err := query.SelectWithDetails(ctx, locale)
	if err != nil {
		return models.PlacesCollection{}, err
	}

	collection := make([]models.Place, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, placeSchemaToModel(row))
	}

	return models.PlacesCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d Database) PlaceExists(ctx context.Context, placeID uuid.UUID) (bool, error) {
	count, err := d.sql.places.New().FilterID(placeID).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d Database) UpdatePlace(ctx context.Context, placeID uuid.UUID, params place.UpdateParams, updatedAt time.Time) error {
	query := d.sql.places.New()

	if params.Class != nil {
		query = query.UpdateClass(*params.Class)
	}
	if params.Point != nil {
		query = query.UpdatePoint(*params.Point)
	}
	if params.Address != nil {
		query = query.UpdateAddress(*params.Address)
	}
	if params.Phone != nil {
		if *params.Phone == "" {
			query = query.UpdatePhone(sql.NullString{Valid: false})
		} else {
			query = query.UpdatePhone(sql.NullString{String: *params.Phone, Valid: true})
		}
	}
	if params.Website != nil {
		if *params.Website == "" {
			query = query.UpdateWebsite(sql.NullString{Valid: false})
		} else {
			query = query.UpdateWebsite(sql.NullString{String: *params.Website, Valid: true})
		}
	}

	return query.FilterID(placeID).Update(ctx, updatedAt)
}

func (d Database) UpdateVerifiedPlace(ctx context.Context, placeID uuid.UUID, verified bool, updatedAt time.Time) error {
	return d.sql.places.New().FilterID(placeID).UpdateVerified(verified).Update(ctx, updatedAt)
}

func (d Database) UpdatePlaceStatus(ctx context.Context, placeID uuid.UUID, status string, updatedAt time.Time) error {
	return d.sql.places.New().FilterID(placeID).UpdateStatus(status).Update(ctx, updatedAt)
}

func (d Database) DeletePlace(ctx context.Context, placeID uuid.UUID) error {
	return d.sql.places.New().FilterID(placeID).Delete(ctx)
}

func placeModelToSchema(model models.PlaceDetails) pgdb.PlaceRow {
	res := pgdb.PlaceRow{
		ID:        model.ID,
		CityID:    model.CityID,
		Class:     model.Class,
		Status:    model.Status,
		Verified:  model.Verified,
		Point:     model.Point,
		Address:   model.Address,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
	if model.DistributorID != nil {
		res.DistributorID = uuid.NullUUID{UUID: *model.DistributorID, Valid: true}
	}
	if model.Website != nil {
		res.Website = sql.NullString{String: *model.Website, Valid: true}
	}
	if model.Phone != nil {
		res.Phone = sql.NullString{String: *model.Phone, Valid: true}
	}

	return res
}

func placeDetailsSchemaToModel(schema pgdb.PlaceRow) models.Place {
	res := models.Place{
		ID:        schema.ID,
		CityID:    schema.CityID,
		Class:     schema.Class,
		Status:    schema.Status,
		Verified:  schema.Verified,
		Point:     schema.Point,
		Address:   schema.Address,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}
	if schema.DistributorID.Valid {
		res.DistributorID = &schema.DistributorID.UUID
	}
	if schema.Website.Valid {
		res.Website = &schema.Website.String
	}
	if schema.Phone.Valid {
		res.Phone = &schema.Phone.String
	}

	return res
}

func placeSchemaToModel(schema pgdb.Place) models.Place {
	res := models.Place{
		ID:          schema.ID,
		CityID:      schema.CityID,
		Class:       schema.Class,
		Status:      schema.Status,
		Verified:    schema.Verified,
		Point:       schema.Point,
		Address:     schema.Address,
		Locale:      schema.Locale,
		Name:        schema.Name,
		Description: schema.Description,
		CreatedAt:   schema.CreatedAt,
		UpdatedAt:   schema.UpdatedAt,
	}
	if schema.DistributorID.Valid {
		res.DistributorID = &schema.DistributorID.UUID
	}
	if schema.Website.Valid {
		res.Website = &schema.Website.String
	}
	if schema.Phone.Valid {
		res.Phone = &schema.Phone.String
	}

	var timetable []models.TimeInterval
	for _, schemaInterval := range schema.Timetable {
		from := models.NumberMinutesToMoment(schemaInterval.StartMin)
		to := models.NumberMinutesToMoment(schemaInterval.EndMin)
		timetable = append(timetable, models.TimeInterval{
			From: from,
			To:   to,
		})
	}

	res.Timetable = models.Timetable{
		Table: timetable,
	}

	return res
}

func placeWithDetailsSchemaToModel(schema pgdb.Place) models.Place {
	res := models.Place{
		ID:        schema.ID,
		CityID:    schema.CityID,
		Class:     schema.Class,
		Status:    schema.Status,
		Verified:  schema.Verified,
		Point:     schema.Point,
		Address:   schema.Address,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}
	if schema.DistributorID.Valid {
		res.DistributorID = &schema.DistributorID.UUID
	}
	if schema.Website.Valid {
		res.Website = &schema.Website.String
	}
	if schema.Phone.Valid {
		res.Phone = &schema.Phone.String
	}

	return res
}
