package types

import (
	"database/sql"

	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq/core"
)

type Model interface {
	SQLTableName() string
	SQLSelect(core.SQLWriter) error
	SQLScan(core.Opts, *sql.Rows) error
	SQLInsert(core.SQLWriter) error
	SQLUpsert(core.SQLWriter) error
}

type DataSource struct {
	DB    *cmsql.Database
	Model Model
}

type ModelPair struct {
	Source *DataSource
	Target *DataSource
}

func NewModelPair(fromDB *cmsql.Database, fromModel Model, toDB *cmsql.Database, toModel Model) *ModelPair {
	return &ModelPair{
		Source: NewDataSource(fromDB, fromModel),
		Target: NewDataSource(toDB, toModel),
	}
}

func NewDataSource(db *cmsql.Database, model Model) *DataSource {
	return &DataSource{
		DB:    db,
		Model: model,
	}
}
