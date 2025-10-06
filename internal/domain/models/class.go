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

func (c Class) IsNil() bool {
	return c.Code == ""
}

type ClassesCollection struct {
	Data  []Class `json:"data"`
	Page  uint64  `json:"page"`
	Size  uint64  `json:"size"`
	Total uint64  `json:"total"`
}
