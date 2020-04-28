package storage

import (
	"context"

	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
)

type ChangesStore struct {
	db *cmsql.Database
}

func NewChangesStore(db *cmsql.Database) *ChangesStore {
	return &ChangesStore{
		db: db,
	}
}

func (cs *ChangesStore) Insert(ctx context.Context, data *model.Callback) error {
	return cs.db.ShouldInsert(data)
}
