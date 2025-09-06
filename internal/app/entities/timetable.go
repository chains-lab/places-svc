package entities

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

type timetableQ interface {
	New() dbx.PlaceTimetablesQ
	Insert(ctx context.Context, in dbx.PlaceTimetable) error
	Get(ctx context.Context) (dbx.PlaceTimetable, error)
	Select(ctx context.Context) ([]dbx.PlaceTimetable, error)
	Delete(ctx context.Context) error

	FilterByID(id uuid.UUID) dbx.PlaceTimetablesQ
	FilterPlaceID(placeID uuid.UUID) dbx.PlaceTimetablesQ
	FilterBetween(startMin, endMin int) dbx.PlaceTimetablesQ

	Count(ctx context.Context) (uint64, error)
	Page(offset, limit uint64) dbx.PlaceTimetablesQ
}

type Timetable struct {
	queries timetableQ
}

func NewTimetable(db *sql.DB) Timetable {
	return Timetable{
		queries: dbx.NewPlaceTimetablesQ(db),
	}
}

func (t Timetable) CreateTimetable(ctx context.Context, placeID uuid.UUID, interval models.TimeInterval) error {
	count, err := t.queries.New().FilterPlaceID(placeID).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not create timetable, cause: %w", err),
		)
	}
	if count >= 70 {
		return errx.ErrorInternal.Raise( //TODO: custom error
			fmt.Errorf("could not create timetable, cause: max timetables reached"),
		)
	}

	start, end := interval.ToNumberMinutes()

	count, err = t.queries.New().FilterPlaceID(placeID).FilterBetween(start, end).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not create timetable, cause: %w", err),
		)
	}

	if count > 0 {
		err = t.queries.New().FilterPlaceID(placeID).FilterBetween(start, end).Delete(ctx)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not create timetable, cause: %w", err),
			)
		}
	}

	err = t.queries.New().Insert(ctx, dbx.PlaceTimetable{
		ID:       uuid.New(),
		PlaceID:  placeID,
		StartMin: start,
		EndMin:   end,
	})
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not create timetable, cause: %w", err),
		)
	}

	return nil
}

func (t Timetable) ListForPlaceID(ctx context.Context, placeID uuid.UUID) ([]models.TimeInterval, error) {
	rows, err := t.queries.New().FilterPlaceID(placeID).Select(ctx)
	if err != nil {
		return nil, errx.ErrorInternal.Raise(
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

	return intervals, nil
}

func (t Timetable) DeleteTimetable(ctx context.Context, ID uuid.UUID) error {
	err := t.queries.New().FilterPlaceID(ID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}

func (t Timetable) DeleteAllTimetableByPlaceID(ctx context.Context, placeID uuid.UUID) error {
	err := t.queries.New().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}
