package models

import "time"

type PlaceClass struct {
	Code       string    `json:"code"`
	FatherCode *string   `json:"father_code,omitempty"` // NULL для корней
	Status     string    `json:"status"`
	Icon       string    `json:"icon"`
	Path       string    `json:"path"` // ltree как text
	Name       string    `json:"name"`
	Locale     string    `json:"locale"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LocaleForKind struct {
	KindCode string `json:"kind_code"`
	Locale   string `json:"locale"`
	Name     string `json:"name"`
}
