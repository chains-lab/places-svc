package enum

import "fmt"

const PlaceStatusActive = "active"
const PlaceStatusInactive = "inactive"
const PlaceStatusBlocked = "blocked"
const PlaceStatusDeprecated = "deprecated"

var placeStatuses = []string{
	PlaceStatusActive,
	PlaceStatusInactive,
	PlaceStatusBlocked,
	PlaceStatusDeprecated,
}

var ErrorInvalidPlaceStatus = fmt.Errorf("invalid place status, must be one of: %v", placeStatuses)

func CheckPlaceStatus(status string) error {
	for _, s := range placeStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s': %w", status, ErrorInvalidPlaceStatus)
}

func GetAllPlaceStatuses() []string {
	return placeStatuses
}
