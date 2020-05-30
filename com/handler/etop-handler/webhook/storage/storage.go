package storage

import (
	"context"

	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
)

type ChangesStore struct {
	db *cmsql.Database
}

func NewChangesStore(db com.WebhookDB) *ChangesStore {
	return &ChangesStore{
		db: db,
	}
}

func (cs *ChangesStore) Insert(ctx context.Context, data *model.Callback) error {
	return cs.db.ShouldInsert(data)
}
