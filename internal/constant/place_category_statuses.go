package constant

import "fmt"

const (
	PlaceCategoryStatusActive   = "active"
	PlaceCategoryStatusInactive = "inactive"
)

var placeCategoryStatuses = []string{
	PlaceCategoryStatusActive,
	PlaceCategoryStatusInactive,
}

var ErrorInvalidPlaceCategoryStatus = fmt.Errorf("invalid place_category status, must be one of: %v", placeCategoryStatuses)

func IsValidPlaceCategoryStatus(status string) error {
	for _, s := range placeCategoryStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorInvalidPlaceCategoryStatus, status)
}

func GetAllPlaceCategoryStatuses() []string {
	return placeCategoryStatuses
}
