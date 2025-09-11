package errx

import "github.com/chains-lab/ape"

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorCategoryNotFound = ape.DeclareError("CATEGORY_NOT_FOUND")

var ErrorKindNotFound = ape.DeclareError("KIND_NOT_FOUND")

var ErrorInvalidLocale = ape.DeclareError("INVALID_LOCALE")
