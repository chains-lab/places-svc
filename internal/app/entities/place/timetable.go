package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/places-svc/internal/app/models"
	"github.com/chains-lab/places-svc/internal/constant"
	"github.com/chains-lab/places-svc/internal/dbx"
	"github.com/chains-lab/places-svc/internal/errx"
	"github.com/google/uuid"
)

type timetableQ interface {
	New() dbx.PlaceTimetablesQ
	Insert(ctx context.Context, in ...dbx.PlaceTimetable) error
	Get(ctx context.Context) (dbx.PlaceTimetable, error)
	Select(ctx context.Context) ([]dbx.PlaceTimetable, error)
	Delete(ctx context.Context) error

	FilterByID(id uuid.UUID) dbx.PlaceTimetablesQ
	FilterPlaceID(placeID uuid.UUID) dbx.PlaceTimetablesQ
	FilterBetween(startMin, endMin int) dbx.PlaceTimetablesQ

	Count(ctx context.Context) (uint64, error)
	Page(offset, limit uint64) dbx.PlaceTimetablesQ
}

// AddTimetable deprecated: use SetTimetable instead
func (p Place) AddTimetable(ctx context.Context, placeID uuid.UUID, interval models.TimeInterval) error {
	_, err := p.Get(ctx, placeID, constant.DefaultLocale)
	if err != nil {
		return err
	}

	count, err := p.timetable.New().FilterPlaceID(placeID).Count(ctx)
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

	count, err = p.timetable.New().FilterPlaceID(placeID).FilterBetween(start, end).Count(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not create timetable, cause: %w", err),
		)
	}

	if count > 0 {
		err = p.timetable.New().FilterPlaceID(placeID).FilterBetween(start, end).Delete(ctx)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("could not create timetable, cause: %w", err),
			)
		}
	}

	err = p.timetable.New().Insert(ctx, dbx.PlaceTimetable{
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

func (p Place) SetTimetable(ctx context.Context, placeID uuid.UUID, intervals ...models.TimeInterval) error {
	_, err := p.Get(ctx, placeID, constant.DefaultLocale)
	if err != nil {
		return err
	}

	err = p.timetable.New().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not upsert timetable, cause: %w", err),
		)
	}

	stmt := make([]dbx.PlaceTimetable, 0, len(intervals))
	for _, interval := range intervals {
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
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not upsert timetable, cause: %w", err),
		)
	}

	return nil
}

// GetTimetable deprecated: use ListTimetable instead
func (p Place) GetTimetable(ctx context.Context, ID uuid.UUID) (models.TimeInterval, error) {
	row, err := p.timetable.New().FilterByID(ID).Get(ctx)
	if err != nil {
		return models.TimeInterval{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not get timetable, cause: %w", err),
		)
	}

	return models.TimeInterval{
		From: models.NumberMinutesToMoment(row.StartMin),
		To:   models.NumberMinutesToMoment(row.EndMin),
	}, nil
}

func (p Place) ListTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
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

// DeleteTimetable deprecated: use DeleteAllTimetable instead
func (p Place) DeleteTimetable(ctx context.Context, ID uuid.UUID) error {
	err := p.timetable.New().FilterPlaceID(ID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}

func (p Place) DeleteAllTimetable(ctx context.Context, placeID uuid.UUID) error {
	err := p.timetable.New().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}
