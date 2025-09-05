package constant

import "fmt"

const (
	PlaceOwnershipPublic  = "public"
	PlaceOwnershipPrivate = "private"
)

var placeOwnerships = []string{
	PlaceOwnershipPublic,
	PlaceOwnershipPrivate,
}

var ErrorInvalidPlaceOwnership = fmt.Errorf("invalid place ownership, must be one of: %v", placeOwnerships)

func IsValidPlaceOwnership(ownership string) error {
	for _, o := range placeOwnerships {
		if o == ownership {
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrorInvalidPlaceOwnership, ownership)
}

func GetAllPlaceOwnerships() []string {
	return placeOwnerships
}
