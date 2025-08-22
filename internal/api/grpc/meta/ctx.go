package meta

type ctxKey int

const (
	RequestIDCtxKey ctxKey = iota
	UserCtxKey
)
