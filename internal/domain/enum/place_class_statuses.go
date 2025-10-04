package enum

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

func CheckPlaceClassStatus(status string) error {
	for _, s := range placeClassStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s': %w", status, ErrorInvalidPlaceClassStatus)
}

func GetAllPlaceClassStatuses() []string {
	return placeClassStatuses
}
