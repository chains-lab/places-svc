package meta

import (
	"context"
	"fmt"

	"github.com/chains-lab/restkit/token"
)

type ctxKey int

const (
	UserCtxKey ctxKey = iota
)

func User(ctx context.Context) (token.UserData, error) {
	if ctx == nil {
		return token.UserData{}, fmt.Errorf("mising context")
	}

	userData, ok := ctx.Value(UserCtxKey).(token.UserData)
	if !ok {
		return token.UserData{}, fmt.Errorf("mising context")
	}

	return userData, nil
}
