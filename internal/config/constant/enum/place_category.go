package enum

import "fmt"

const (
	PlaceCategoryOther      = "other"      // Other category
	PlaceCategoryRestaurant = "restaurant" // Restaurant category
	PlaceCategoryStore      = "store"      // Store category
)

var placeCategories = []string{
	PlaceCategoryOther,
	PlaceCategoryRestaurant,
	PlaceCategoryStore,
}

// ErrorPlaceCategoryNotSupported is returned when the provided place category is not supported.
var ErrorPlaceCategoryNotSupported = fmt.Errorf("place category must be one of: %s", GetAllPlaceCategories())

func ParsePlaceCategory(category string) (string, error) {
	for _, c := range placeCategories {
		if c == category {
			return c, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", category, ErrorPlaceCategoryNotSupported)
}

func GetAllPlaceCategories() []string {
	return placeCategories
}
