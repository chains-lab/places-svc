package models

import "time"

type PlaceClass struct {
	Code      string    `json:"code"`
	Father    *string   `json:"father,omitempty"` // NULL для корней
	Status    string    `json:"status"`
	Icon      string    `json:"icon"`
	Path      string    `json:"path"` // ltree как text
	Name      string    `json:"name"`
	Locale    string    `json:"locale"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LocaleForClass struct {
	Class  string `json:"kind_code"`
	Locale string `json:"locale"`
	Name   string `json:"name"`
}
