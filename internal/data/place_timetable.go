package data

import (
	"context"

	"github.com/google/uuid"
)

type PlaceTimetablesQ interface {
	New() PlaceTimetablesQ

	Insert(ctx context.Context, in ...PlaceTimetable) error
	Upsert(ctx context.Context, in ...PlaceTimetable) error
	Get(ctx context.Context) (PlaceTimetable, error)
	Select(ctx context.Context) ([]PlaceTimetable, error)
	Delete(ctx context.Context) error

	FilterByID(id uuid.UUID) PlaceTimetablesQ
	FilterPlaceID(placeID uuid.UUID) PlaceTimetablesQ
	FilterBetween(start, end int) PlaceTimetablesQ

	Page(limit, offset uint64) PlaceTimetablesQ
	Count(ctx context.Context) (uint64, error)
}

type PlaceTimetable struct {
	ID       uuid.UUID `db:"id"        json:"id"`
	PlaceID  uuid.UUID `db:"place_id"  json:"place_id"`
	StartMin int       `db:"start_min" json:"start_min"`
	EndMin   int       `db:"end_min"   json:"end_min"`
}
