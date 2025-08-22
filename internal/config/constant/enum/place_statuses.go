package enum

import "fmt"

const (
	PlaceStatusActive   = "active"   // Place is active
	PlaceStatusInactive = "inactive" // Place is inactive
	PlaceStatusBlocked  = "blocked"  // Place is blocked
)

var placeStatuses = []string{
	PlaceStatusActive,
	PlaceStatusInactive,
	PlaceStatusBlocked,
}

// ErrorPlaceStatusNotSupported is returned when the provided place status is not supported.
var ErrorPlaceStatusNotSupported = fmt.Errorf("place status must be one of: %s", GetAllPlaceStatuses())

// ParsePlaceStatus checks if the provided status is valid and returns it.
func ParsePlaceStatus(status string) (string, error) {
	for _, s := range placeStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorPlaceStatusNotSupported)
}

func GetAllPlaceStatuses() []string {
	return placeStatuses
}
