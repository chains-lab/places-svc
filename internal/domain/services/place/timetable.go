package place

import (
	"context"
	"fmt"

	"github.com/chains-lab/enum"
	"github.com/chains-lab/places-svc/internal/data/schemas"
	"github.com/chains-lab/places-svc/internal/domain/errx"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (m Service) SetTimetable(
	ctx context.Context,
	placeID uuid.UUID,
	intervals models.Timetable,
) (models.PlaceWithDetails, error) {
	place, err := m.Get(ctx, placeID, enum.DefaultLocale)
	if err != nil {
		return models.PlaceWithDetails{}, err
	}

	err = m.db.PlaceTimetables().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not upsert timetable, cause: %w", err),
		)
	}

	stmt := make([]schemas.PlaceTimetable, 0, len(intervals.Table))
	for _, interval := range intervals.Table {
		start, end := interval.ToNumberMinutes()
		stmt = append(stmt, schemas.PlaceTimetable{
			ID:       uuid.New(),
			PlaceID:  placeID,
			StartMin: start,
			EndMin:   end,
		})
	}

	err = m.db.PlaceTimetables().Insert(ctx, stmt...)
	if err != nil {
		return models.PlaceWithDetails{}, errx.ErrorInternal.Raise(
			fmt.Errorf("could not upsert timetable, cause: %w", err),
		)
	}

	place.Timetable = intervals

	return place, nil
}

func (m Service) GetTimetable(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	rows, err := m.db.PlaceTimetables().FilterPlaceID(placeID).Select(ctx)
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

func (m Service) DeleteTimetable(ctx context.Context, placeID uuid.UUID) error {
	_, err := m.Get(ctx, placeID, enum.DefaultLocale)
	if err != nil {
		return err
	}

	err = m.db.PlaceTimetables().FilterPlaceID(placeID).Delete(ctx)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("could not delete timetable, cause: %w", err),
		)
	}

	return nil
}
