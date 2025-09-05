package models

import "time"

type Category struct {
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	Icon      string    `json:"icon"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryLocale struct {
	CategoryCode string `json:"category_code"`
	Locale       string `json:"locale"`
	Name         string `json:"name"`
}
