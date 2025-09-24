package schemas

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
	ID       uuid.UUID `storage:"id"        json:"id"`
	PlaceID  uuid.UUID `storage:"place_id"  json:"place_id"`
	StartMin int       `storage:"start_min" json:"start_min"`
	EndMin   int       `storage:"end_min"   json:"end_min"`
}
