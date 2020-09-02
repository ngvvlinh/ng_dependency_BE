package cmsql

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqltrace"
)

func MockPostgresDB() *Database {
	db := sq.Database{}
	sq.DollarMarker(&db)
	sq.DoubleQuoteEscape(&db)
	sq.UseArrayInsteadOfJSON(&db, true)
	return &Database{
		id: cm.NewID().Int64(),
		db: db,
	}
}

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

func monitorQuery(entry *sq.LogEntry) {
	sqltrace.Trace(entry)
}
