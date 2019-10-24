package cmsql

import "context"

type QueryFactory func() QueryInterface

func NewQueryFactory(ctx context.Context, db *Database) QueryFactory {
	return func() QueryInterface {
		return GetTxOrNewQuery(ctx, db)
	}
}

func GetTxOrNewQuery(ctx context.Context, db *Database) QueryInterface {
	tx := ctx.Value(db.TxKey())
	if tx == nil {
		return db.WithContext(ctx)
	}
	return tx.(Tx)
}
