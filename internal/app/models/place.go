package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type Place struct {
	ID            uuid.UUID  `json:"id"`
	CityID        uuid.UUID  `json:"city_id"`
	DistributorID *uuid.UUID `json:"distributor_id"`
	Class         string     `json:"class"`

	Status    string    `json:"status"`
	Verified  bool      `json:"verified"`
	Ownership string    `json:"ownership"`
	Point     orb.Point `json:"point"`
	Address   string    `json:"address"`

	Website *string `json:"website"`
	Phone   *string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlaceLocale struct {
	PlaceID     uuid.UUID `json:"place_id"`
	Locale      string    `json:"locale"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type PlaceWithDetails struct {
	Place     Place
	Locale    PlaceLocale
	Timetable Timetable
}
