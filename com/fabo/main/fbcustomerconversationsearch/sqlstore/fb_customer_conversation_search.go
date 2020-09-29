package sqlstore

import (
	"context"

	"o.o/api/fabo/fbcustomerconversationsearch"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbcustomerconversationsearch/convert"
	"o.o/backend/com/fabo/main/fbcustomerconversationsearch/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/validate"
)

type FbExternalCustomerConversationSearchStoreFactory func(ctx context.Context) *FbCustomerConversationSearchStore

func NewFbExternalCustomerConversationSearchStore(db *cmsql.Database) FbExternalCustomerConversationSearchStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbCustomerConversationSearchStore {
		return &FbCustomerConversationSearchStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbCustomerConversationSearchStore struct {
	ft      FbCustomerConversationSearchFilters
	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
}

func (s *FbCustomerConversationSearchStore) ByPageIDs(pageIds []string) *FbCustomerConversationSearchStore {
	s.preds = append(s.preds, sq.In("external_page_id", pageIds))
	return s
}

func (s *FbCustomerConversationSearchStore) ByExternalUserNameNorm(extUserName string) *FbCustomerConversationSearchStore {
	s.preds = append(s.preds, sq.NewExpr(
		"external_user_name_norm @@ ?::tsquery",
		validate.NormalizeSearchQueryAnd(extUserName)),
	)
	return s
}

func (s *FbCustomerConversationSearchStore) ListFbCustomerConversationSearchDB() (
	model.FbCustomerConversationSearchs, error,
) {
	var customerConvsDB model.FbCustomerConversationSearchs
	if err := s.query().Where(s.preds).OrderBy("created_at desc").Find(&customerConvsDB); err != nil {
		return nil, err
	}
	return customerConvsDB, nil
}

func (s *FbCustomerConversationSearchStore) ListFbCustomerConversationSearch() (
	[]*fbcustomerconversationsearch.FbCustomerConversationSearch, error,
) {
	customerConvsDB, err := s.ListFbCustomerConversationSearchDB()
	if err != nil {
		return nil, err
	}

	customerConvs := convert.Convert_fbcustomerconversationsearchmodel_FbCustomerConversationSearchs_fbcustomerconversationsearch_FbCustomerConversationSearchs(customerConvsDB)
	return customerConvs, nil
}
