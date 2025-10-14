package meta

import (
	"context"
	"fmt"

	"github.com/chains-lab/restkit/auth"
)

type ctxKey int

const (
	UserCtxKey ctxKey = iota
)

func User(ctx context.Context) (auth.UserData, error) {
	if ctx == nil {
		return auth.UserData{}, fmt.Errorf("mising context")
	}

	userData, ok := ctx.Value(UserCtxKey).(auth.UserData)
	if !ok {
		return auth.UserData{}, fmt.Errorf("mising context")
	}

	return userData, nil
}
