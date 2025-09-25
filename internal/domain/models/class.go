package models

import (
	"time"
)

type Class struct {
	Code      string    `json:"code"`
	Parent    *string   `json:"parent,omitempty"`
	Status    string    `json:"status"`
	Icon      string    `json:"icon"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClassesCollection struct {
	Data  []Class `json:"data"`
	Page  uint    `json:"page"`
	Size  uint    `json:"size"`
	Total uint    `json:"total"`
}
