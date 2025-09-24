package errx

import "github.com/chains-lab/ape"

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorInvalidLocale = ape.DeclareError("INVALID_LOCALE")
