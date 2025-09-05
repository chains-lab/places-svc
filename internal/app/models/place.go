package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type Place struct {
	ID            uuid.UUID     `json:"id"`
	CityID        uuid.UUID     `json:"city_id"`
	DistributorID uuid.NullUUID `json:"distributor_id"`
	TypeCode      string        `json:"type_code"`

	Status    string    `json:"status"`
	Verified  bool      `json:"verified"`
	Ownership string    `json:"ownership"`
	Point     orb.Point `json:"point"`

	Locale      string         `json:"locale"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	Description sql.NullString `json:"description"`

	Website *string `json:"website"`
	Phone   *string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PlaceLocale struct {
	PlaceID     uuid.UUID      `json:"place_id"`
	Locale      string         `json:"locale"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	Description sql.NullString `json:"description"`
}
