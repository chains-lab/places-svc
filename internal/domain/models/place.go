package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
)

type PlaceDetails struct {
	ID        uuid.UUID  `json:"id"`
	CityID    uuid.UUID  `json:"city_id"`
	CompanyID *uuid.UUID `json:"company_id"`
	Class     string     `json:"class"`

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

func (p PlaceDetails) IsNil() bool {
	return p.ID == uuid.Nil
}

type PlaceLocale struct {
	PlaceID     uuid.UUID `json:"place_id"`
	Locale      string    `json:"locale"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (p PlaceLocale) IsNil() bool {
	return p.PlaceID == uuid.Nil
}

type Place struct {
	ID        uuid.UUID  `json:"id"`
	CityID    uuid.UUID  `json:"city_id"`
	CompanyID *uuid.UUID `json:"company_id"`
	Class     string     `json:"class"`

	Status    string    `json:"status"`
	Verified  bool      `json:"verified"`
	Ownership string    `json:"ownership"`
	Point     orb.Point `json:"point"`
	Address   string    `json:"address"`

	Locale      string `json:"locale"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Website *string `json:"website"`
	Phone   *string `json:"phone"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Timetable Timetable
}

func (p Place) IsNil() bool {
	return p.ID == uuid.Nil
}

func (p Place) Details() PlaceDetails {
	return PlaceDetails{
		ID:        p.ID,
		CityID:    p.CityID,
		CompanyID: p.CompanyID,
		Class:     p.Class,
		Status:    p.Status,
		Verified:  p.Verified,
		Ownership: p.Ownership,
		Point:     p.Point,
		Address:   p.Address,
		Website:   p.Website,
		Phone:     p.Phone,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type PlacesCollection struct {
	Data  []Place `json:"data"`
	Page  uint64  `json:"page"`
	Size  uint64  `json:"size"`
	Total uint64  `json:"total"`
}

type PlaceLocaleCollection struct {
	Data  []PlaceLocale `json:"data"`
	Page  uint64        `json:"page"`
	Size  uint64        `json:"size"`
	Total uint64        `json:"total"`
}
