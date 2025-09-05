package constant

import "fmt"

const PlaceStatusActive = "active"
const PlaceStatusInactive = "inactive"
const PlaceStatusBlocked = "blocked"

var placeStatuses = []string{
	PlaceStatusActive,
	PlaceStatusInactive,
	PlaceStatusBlocked,
}

var ErrorInvalidPlaceStatus = fmt.Errorf("invalid place status, must be one of: %v", placeStatuses)

func IsValidPlaceStatus(status string) error {
	for _, s := range placeStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorInvalidPlaceStatus, status)
}

func GetAllPlaceStatuses() []string {
	return placeStatuses
}
