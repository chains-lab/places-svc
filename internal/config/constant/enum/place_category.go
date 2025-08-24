package enum

import "fmt"

const (
	PlaceCategoryFodAndDrinks           = "food and drinks"
	PlaceCategoryShops                  = "shops"
	PlaceCategoryServices               = "services"
	PlaceCategoryHotelsAndAccommodation = "hotels and accommodation"
	PlaceCategoryActiveLeisure          = "active leisure"
	PlaceCategoryReligion               = "religion"
	PlaceCategoryOfficeAndFactories     = "office and factories"
	PlaceCategoryResidenceBuildings     = "residence buildings"
	PlaceCategoryEducation              = "education"
	PlaceCategoryHealthcare             = "healthcare"
	PlaceCategoryTransport              = "transport"
	PlaceCategoryCulture                = "culture"
)

var placeCategories = []string{
	PlaceCategoryFodAndDrinks,
	PlaceCategoryShops,
	PlaceCategoryServices,
	PlaceCategoryHotelsAndAccommodation,
	PlaceCategoryActiveLeisure,
	PlaceCategoryReligion,
	PlaceCategoryOfficeAndFactories,
	PlaceCategoryResidenceBuildings,
	PlaceCategoryEducation,
	PlaceCategoryHealthcare,
	PlaceCategoryTransport,
	PlaceCategoryCulture,
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
