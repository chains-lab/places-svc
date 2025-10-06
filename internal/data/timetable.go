package data

import (
	"context"

	"github.com/chains-lab/places-svc/internal/data/pgdb"
	"github.com/chains-lab/places-svc/internal/domain/models"
	"github.com/google/uuid"
)

func (d Database) SetTimetable(ctx context.Context, placeID uuid.UUID, intervals models.Timetable) error {
	stmt := make([]pgdb.PlaceTimetableRow, 0, len(intervals.Table))
	for _, interval := range intervals.Table {
		stmt = append(stmt, pgdb.PlaceTimetableRow{
			ID:       uuid.New(),
			PlaceID:  placeID,
			StartMin: interval.From.ToNumberMinutes(),
			EndMin:   interval.To.ToNumberMinutes(),
		})
	}

	return d.sql.timetables.New().Upsert(ctx, stmt...)
}

func (d Database) GetTimetableByPlaceID(ctx context.Context, placeID uuid.UUID) (models.Timetable, error) {
	rows, err := d.sql.timetables.New().FilterPlaceID(placeID).Select(ctx)
	if err != nil {
		return models.Timetable{}, err
	}

	intervals := make([]models.TimeInterval, 0, len(rows))
	for _, row := range rows {
		intervals = append(intervals, models.TimeInterval{
			From: models.NumberMinutesToMoment(row.StartMin),
			To:   models.NumberMinutesToMoment(row.EndMin),
		})
	}

	return models.Timetable{Table: intervals}, nil
}

func (d Database) DeleteTimetableByPlaceID(ctx context.Context, placeID uuid.UUID) error {
	return d.sql.timetables.New().FilterPlaceID(placeID).Delete(ctx)
}
