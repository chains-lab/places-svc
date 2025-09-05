package errx

import "github.com/chains-lab/ape"

var ErrorClassNotFound = ape.DeclareError("CLASS_NOT_FOUND")

var ErrorClassAlreadyExists = ape.DeclareError("CLASS_ALREADY_EXISTS")

// ErrorClassLocaleNotFound use only when we try to get locale for specific class, but it does not exist
var ErrorClassLocaleNotFound = ape.DeclareError("CLASS_LOCALE_NOT_FOUND")
