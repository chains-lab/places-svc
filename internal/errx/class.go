package errx

import "github.com/chains-lab/ape"

var ErrorClassNotFound = ape.DeclareError("CLASS_NOT_FOUND")

var ErrorClassCodeAlreadyTaken = ape.DeclareError("CLASS_CODE_ALREADY_TAKEN")

// ErrorClassLocaleNotFound use only when we try to get locale for specific class, but it does not exist
var ErrorClassLocaleNotFound = ape.DeclareError("CLASS_LOCALE_NOT_FOUND")

var ErrorCurrentLocaleNotFoundForClass = ape.DeclareError("CURRENT_LOCALE_NOT_FOUND_FOR_CLASS")

var ErrorNedAtLeastOneLocaleForClass = ape.DeclareError("NEED_AT_LEAST_ONE_LOCALE_FOR_CLASS")

var ErrorClassParentCycle = ape.DeclareError("CLASS_FATHER_CYCLE")

var ErrorClassParentEqualCode = ape.DeclareError("CLASS_PARENT_EQUAL_CODE")

var ErrorClassHasChildren = ape.DeclareError("CLASS_HAS_CHILDREN")

var ErrorClassDeactivateReplaceSame = ape.DeclareError("CLASS_DEACTIVATE_REPLACE_SAME")

var ErrorCannotDeleteActiveClass = ape.DeclareError("CANNOT_DELETE_ACTIVE_CLASS")

var ErrorCantDeleteClassWithPlaces = ape.DeclareError("CLASS_HAS_PLACES")

var ErrorClassDeactivateReplaceInactive = ape.DeclareError("CLASS_DEACTIVATE_REPLACE_INACTIVE")

var ErrorParentClassNotFound = ape.DeclareError("PARENT_CLASS_NOT_FOUND")

var ErrorClassStatusInvalid = ape.DeclareError("CLASS_STATUS_INVALID")

var ErrorClassStatusIsNotActive = ape.DeclareError("CLASS_STATUS_IS_NOT_ACTIVE")
