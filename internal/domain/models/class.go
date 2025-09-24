package models

import (
	"time"
)

type ClassDetails struct {
	Code      string    `json:"code"`
	Parent    *string   `json:"parent,omitempty"`
	Status    string    `json:"status"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClassLocale struct {
	Class  string `json:"kind_code"`
	Locale string `json:"locale"`
	Name   string `json:"name"`
}

type Class struct {
	Code      string    `json:"code"`
	Parent    *string   `json:"parent,omitempty"`
	Status    string    `json:"status"`
	Icon      string    `json:"icon"`
	Locale    string    `json:"locale"`
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

type ClassLocaleCollection struct {
	Data  []ClassLocale `json:"data"`
	Page  uint          `json:"page"`
	Size  uint          `json:"size"`
	Total uint          `json:"total"`
}
