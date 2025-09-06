package errx

import "github.com/chains-lab/ape"

var ErrorPlaceNotFound = ape.DeclareError("PLACE_NOT_FOUND")

var ErrorPlaceAlreadyExists = ape.DeclareError("PLACE_ALREADY_EXISTS")

// ErrorPlaceLocaleNotFound use only when we try to get locale for specific place, but it does not exist
var ErrorPlaceLocaleNotFound = ape.DeclareError("PLACE_LOCALE_NOT_FOUND")

var ErrorNeedAtLeastOneLocaleForPlace = ape.DeclareError("NEED_AT_LEAST_ONE_LOCALE_FOR_PLACE")

var ErrorCurrentLocaleNotFoundForPlace = ape.DeclareError("CURRENT_LOCALE_NOT_FOUND_FOR_PLACE")

var ErrorPlaceStatusInvalid = ape.DeclareError("PLACE_STATUS_INVALID")

var ErrorPlaceOwnershipInvalid = ape.DeclareError("PLACE_OWNERSHIP_INVALID")
