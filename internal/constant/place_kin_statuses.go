package constant

import "fmt"

const (
	PlaceKindStatusActive   = "active"
	PlaceKindStatusInactive = "inactive"
)

var placeKindStatuses = []string{
	PlaceKindStatusActive,
	PlaceKindStatusInactive,
}

var ErrorInvalidPlaceKindStatus = fmt.Errorf("invalid place_kind status, must be one of: %v", placeKindStatuses)

func IsValidPlaceKindStatus(status string) error {
	for _, s := range placeKindStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorInvalidPlaceKindStatus, status)
}

func GetAllPlaceKindStatuses() []string {
	return placeKindStatuses
}
