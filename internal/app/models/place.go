package models

import (
	"time"

	"github.com/google/uuid"
)

type Place struct {
	ID            uuid.UUID  `db:"id"`
	DistributorID *uuid.UUID `db:"distributor_id"`
	Type          string     `db:"type"`
	Status        string     `db:"status"`
	Ownership     string     `db:"ownership"`
	Name          string     `db:"name"`
	Description   *string    `db:"description"`
	Coords        Coords     `db:"coords"`
	Address       string     `db:"address"`
	Website       *string    `db:"website"`
	Phone         *string    `db:"phone"`
	UpdatedAt     time.Time  `db:"updated_at"`
	CreatedAt     time.Time  `db:"created_at"`
}

type Coords struct {
	Lon float64 `db:"lon"`
	Lat float64 `db:"lat"`
}
