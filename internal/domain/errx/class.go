package errx

import "github.com/chains-lab/ape"

// ErrorClassNotFound is used when we try to get/update/delete class by code, but it does not exist
// Its 404 - Not Found
var ErrorClassNotFound = ape.DeclareError("CLASS_NOT_FOUND")

// ErrorClassStatusInvalid is used when we try to set/get invalid status to class
// Its 400 - Bad Request
var ErrorClassStatusInvalid = ape.DeclareError("CLASS_STATUS_INVALID")

// ErrorClassCodeAlreadyTaken is used when we try to create class with code that already exists
// Its 409 - Conflict
var ErrorClassCodeAlreadyTaken = ape.DeclareError("CLASS_CODE_ALREADY_TAKEN")

// ErrorClassNameAlreadyTaken is used when we try to create/update class with name that already exists for specific locale
// Its 409 - Conflict
var ErrorClassNameAlreadyTaken = ape.DeclareError("CLASS_NAME_ALREADY_TAKEN")

// ErrorCannotDeleteDefaultLocaleForClass use only when we try to delete default locale for specific class, but it is default
var ErrorCannotDeleteDefaultLocaleForClass = ape.DeclareError("CANNOT_DELETE_DEFAULT_LOCALE_FOR_CLASS")

// ErrorParentClassNotFound is used when we try to create/update class with parent that does not exist
// Its 404 - Not Found
var ErrorParentClassNotFound = ape.DeclareError("PARENT_CLASS_NOT_FOUND")

// ErrorClassParentCycle is used when we try to create/update class with parent that creates cycle
// Its 409 - Conflict
var ErrorClassParentCycle = ape.DeclareError("CLASS_FATHER_CYCLE")

// ErrorClassNameExists is used when we try to create/update class with name that already exists
// Its 409 - Conflict
var ErrorClassNameExists = ape.DeclareError("CLASS_NAME_EXISTS")

// ErrorClassHasChildren is used when we try to delete class that has children
// Its 409 - Conflict
var ErrorClassHasChildren = ape.DeclareError("CLASS_HAS_CHILDREN")

// ErrorClassDeactivateReplaceSame is used when we try to deactivate class and replace with itself
// Its 409 - Conflict
var ErrorClassDeactivateReplaceSame = ape.DeclareError("CLASS_DEACTIVATE_REPLACE_SAME")

// ErrorCannotDeleteActiveClass is used when we try to delete class that is active
// Its 409 - Conflict
var ErrorCannotDeleteActiveClass = ape.DeclareError("CANNOT_DELETE_ACTIVE_CLASS")

// ErrorCantDeleteClassWithPlaces is used when we try to delete class that has places
// Its 409 - Conflict
var ErrorCantDeleteClassWithPlaces = ape.DeclareError("CLASS_HAS_PLACES")

// ErrorClassDeactivateReplaceInactive is used when we try to deactivate class and replace with inactive class
// Its 409 - Conflict
var ErrorClassDeactivateReplaceInactive = ape.DeclareError("CLASS_DEACTIVATE_REPLACE_INACTIVE")

var ErrorClassStatusIsNotActive = ape.DeclareError("CLASS_STATUS_IS_NOT_ACTIVE")
