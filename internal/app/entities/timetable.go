package entities

//type Timetable struct {
//	queries dbx.PlaceTimetablesQ
//}
//
//func NewTimetable(db *sql.DB) Timetable {
//	return Timetable{
//		queries: dbx.NewPlaceTimetablesQ(db),
//	}
//}
//
//func (t Timetable) CreateTimetable(ctx context.Context, placeID uuid.UUID, interval models.TimeInterval) error {
//	count, err := t.queries.New().FilterByID(placeID).Count(ctx)
//	if err != nil {
//		return errx.ErrorInternal.Raise(
//			fmt.Errorf("could not create timetable, cause: %w", err),
//		)
//	}
//	if count >= 70 {
//		return errx.ErrorInternal.Raise( //TODO: custom error
//			fmt.Errorf("could not create timetable, cause: max timetables reached"),
//		)
//	}
//
//	start, end := interval.ToNumberMinutes()
//
//	count, err = t.queries.New().FilterByID(placeID).FilterBetween(start, end).Count(ctx)
//	if err != nil {
//		return errx.ErrorInternal.Raise(
//			fmt.Errorf("could not create timetable, cause: %w", err),
//		)
//	}
//
//	if count > 0 {
//		err = t.queries.New().FilterByID(placeID).FilterBetween(start, end).Delete(ctx)
//		if err != nil {
//			return errx.ErrorInternal.Raise(
//				fmt.Errorf("could not create timetable, cause: %w", err),
//			)
//		}
//	}
//
//	err = t.queries.New().Insert(ctx, dbx.PlaceTimetable{
//		ID:       uuid.New(),
//		PlaceID:  placeID,
//		StartMin: start,
//		EndMin:   end,
//	})
//	if err != nil {
//		return errx.ErrorInternal.Raise(
//			fmt.Errorf("could not create timetable, cause: %w", err),
//		)
//	}
//
//	return nil
//}
//
//func (t Timetable) ListForPlaceID(ctx context.Context, placeID uuid.UUID) ([]models.TimeInterval, error) {
//	rows, err := t.queries.New().FilterByID(placeID).Select(ctx)
//	if err != nil {
//		return nil, errx.ErrorInternal.Raise(
//			fmt.Errorf("could not list timetable, cause: %w", err),
//		)
//	}
//
//	intervals := make([]models.TimeInterval, 0, len(rows))
//	for _, row := range rows {
//		intervals = append(intervals, models.TimeInterval{
//			From: models.NumberMinutesToMoment(row.StartMin),
//			To:   models.NumberMinutesToMoment(row.EndMin),
//		})
//	}
//
//	return intervals, nil
//}
//
//func (t Timetable) DeleteTimetable(ctx context.Context, ID uuid.UUID) error {
//	err := t.queries.New().FilterByID(ID).Delete(ctx)
//	if err != nil {
//		return errx.ErrorInternal.Raise(
//			fmt.Errorf("could not delete timetable, cause: %w", err),
//		)
//	}
//
//	return nil
//}
//
//func (t Timetable) DeleteAllTimetableByPlaceID(ctx context.Context, placeID uuid.UUID) error {
//	err := t.queries.New().FilterByPlaceID(placeID).Delete(ctx)
//	if err != nil {
//		return errx.ErrorInternal.Raise(
//			fmt.Errorf("could not delete timetable, cause: %w", err),
//		)
//	}
//
//	return nil
//}
