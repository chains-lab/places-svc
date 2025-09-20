package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

type timetableQ interface {
	New() dbx.PlaceTimetablesQ
	Insert(ctx context.Context, in ...dbx.PlaceTimetable) error
	Upsert(ctx context.Context, in ...dbx.PlaceTimetable) error
	Get(ctx context.Context) (dbx.PlaceTimetable, error)
	Select(ctx context.Context) ([]dbx.PlaceTimetable, error)
	Delete(ctx context.Context) error

	FilterByID(id uuid.UUID) dbx.PlaceTimetablesQ
	FilterPlaceID(placeID uuid.UUID) dbx.PlaceTimetablesQ
	FilterBetween(startMin, endMin int) dbx.PlaceTimetablesQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.PlaceTimetablesQ
}

func (p Place) SetTimetable(ctx context.Context, placeID uuid.UUID, intervals models.Timetable) (models.PlaceWithDetails, error) {
	place, err := p.Get(ctx, placeID, enum.DefaultLocale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	err = p.timetable.New().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not upsert timetable, cause: %w", err),
		)
	}

	stmt := make([]dbx.PlaceTimetable, 0, len(intervals.Table))
	for _, interval := range intervals.Table {
		start, end := interval.ToNumberMinutes()
		stmt = append(stmt, dbx.PlaceTimetable{
			ID:       uuid.New(),
			PlaceID:  placeID,
			StartMin: start,
			EndMin:   end,
		})
	}

	err = p.timetable.New().Insert(ctx, stmt...)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not upsert timetable, cause: %w", err),
		)
	}

	place.Timetable = intervals

	return place, nil
}

func (p Place) GetTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	rows, err := p.timetable.New().FilterPlaceID(placeID).Select(ctx)
	if err != nil {
		return models.Timetable{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not list timetable, cause: %w", err),
		)
	}

	intervals := make([]models.TimeInterval, 0, len(rows))
	for _, row := range rows {
		intervals = append(intervals, models.TimeInterval{
			From: models.NumberMinutesToMoment(row.StartMin),
			To:   models.NumberMinutesToMoment(row.EndMin),
		})
	}

	return models.Timetable{
		Table: intervals,
	}, nil
}

func (p Place) DeleteTimetable(ctx context.Context, placeID uuid.UUID) error {
	_, err := p.Get(ctx, placeID, enum.DefaultLocale)
	if err != nil {
		return err
	}

	err = p.timetable.New().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}
