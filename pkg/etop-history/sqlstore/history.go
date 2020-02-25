package sqlstore

import (
	"context"

	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq/core"
	"etop.vn/capi/dot"
)

const (
	OpInsert = "INSERT"
	OpUpdate = "UPDATE"
	OpDelete = "DELETE"
)

type HistoryStoreFactory func(context.Context) *HistoryStore

type HistoryStore struct {
	query cmsql.QueryFactory
}

func NewHistoryStore(db *cmsql.Database) HistoryStoreFactory {
	return func(ctx context.Context) *HistoryStore {
		return &HistoryStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

func (s *HistoryStore) GetHistory(_model core.IGet, rid int64) (bool, error) {
	return s.query().Where("rid = ?", rid).Get(_model)
}

func (s *HistoryStore) ListCustomerRelationships(_models core.IFind, customerIDs, groupIDs []dot.ID, operation string) error {
	query := s.query()
	if len(customerIDs) != 0 {
		query = query.In("customer_id", customerIDs)
	}
	if len(groupIDs) != 0 {
		query = query.In("group_id", groupIDs)
	}
	query = query.Where("_op = ?", operation)
	return query.Find(_models)
}

func (s *HistoryStore) ListProductCollectionRelationships(_models core.IFind, shopID dot.ID, productIDs, collectionIDs []dot.ID, operation string) error {
	query := s.query().Where("shop_id = ?", shopID)
	if len(productIDs) != 0 {
		query = query.In("product_id", productIDs)
	}
	if len(collectionIDs) != 0 {
		query = query.In("collection_id", collectionIDs)
	}
	query = query.Where("_op = ?", operation)
	return query.Find(_models)
}
