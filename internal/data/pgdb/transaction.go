package pgdb

import (
	"context"
	"database/sql"
)

type txKeyType struct{}

var TxKey = txKeyType{}

func TxFromCtx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	return tx, ok
}
