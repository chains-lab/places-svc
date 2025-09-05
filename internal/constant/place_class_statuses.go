package constant

import "fmt"

const (
	PlaceClassStatusesActive   = "active"
	PlaceClassStatusesInactive = "inactive"
)

var placeClassStatuses = []string{
	PlaceClassStatusesActive,
	PlaceClassStatusesInactive,
}

var ErrorInvalidPlaceClassStatus = fmt.Errorf("invalid place class status, must be one of: %v", GetAllPlaceClassStatuses())

func IsValidPlaceClassStatus(status string) error {
	for _, s := range placeClassStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorInvalidPlaceClassStatus, status)
}

func GetAllPlaceClassStatuses() []string {
	return placeClassStatuses
}
