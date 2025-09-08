package models

import "time"

type PlaceClass struct {
	Code      string    `json:"code"`
	Parent    *string   `json:"parent,omitempty"`
	Status    string    `json:"status"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LocaleForClass struct {
	Class  string `json:"kind_code"`
	Locale string `json:"locale"`
	Name   string `json:"name"`
}

type PlaceClassWithLocale struct {
	Data   PlaceClass
	Locale LocaleForClass
}
