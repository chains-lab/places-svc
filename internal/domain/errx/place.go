package errx

import "github.com/chains-lab/ape"

// ErrorPlaceNotFound indicates that the specified place was not found in the system
var ErrorPlaceNotFound = ape.DeclareError("PLACE_NOT_FOUND")

// ErrorNeedAtLeastOneLocaleForPlace indicates that at least one locale must be provided for a place
var ErrorNeedAtLeastOneLocaleForPlace = ape.DeclareError("NEED_AT_LEAST_ONE_LOCALE_FOR_PLACE")

// ErrorPlaceForDeleteMustBeInactive indicates that a place must be inactive before it can be deleted
var ErrorPlaceForDeleteMustBeInactive = ape.DeclareError("PLACE_FOR_DELETE_MUST_BE_INACTIVE")

// ErrorInvalidLocale indicates that the provided locale is invalid or not supported
var ErrorCannotSetStatusBlocked = ape.DeclareError("CANNOT_SET_STATUS_BLOCKED")
