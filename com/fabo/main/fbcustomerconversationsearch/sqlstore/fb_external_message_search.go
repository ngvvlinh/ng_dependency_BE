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

type FbExternalMessageSearchStoreFactory func(ctx context.Context) *FbExternalMessageSearchStore

func NewFbExternalMessageSearchStore(db *cmsql.Database) FbExternalMessageSearchStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalMessageSearchStore {
		return &FbExternalMessageSearchStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalMessageSearchStore struct {
	ft FbExternalMessageSearchFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
}

func (s *FbExternalMessageSearchStore) ByPageIDs(pageIDs []string) *FbExternalMessageSearchStore {
	s.preds = append(s.preds, sq.In("external_page_id", pageIDs))
	return s
}

func (s *FbExternalMessageSearchStore) ByExternalMessageNorm(msg string) *FbExternalMessageSearchStore {
	s.preds = append(s.preds, sq.NewExpr(
		"external_message_norm @@ ?::tsquery",
		validate.NormalizeSearchQueryAnd(msg)),
	)
	return s
}

func (s *FbExternalMessageSearchStore) ListExternalMessageSearchDB() (model.FbExternalMessageSearchs, error) {
	var mgss model.FbExternalMessageSearchs
	if err := s.query().Where(s.preds).OrderBy("created_at desc").Find(&mgss); err != nil {
		return nil, err
	}

	return mgss, nil
}

func (s *FbExternalMessageSearchStore) ListExternalMessageSearch() ([]*fbcustomerconversationsearch.FbExternalMessageSearch, error) {
	msgDBs, err := s.ListExternalMessageSearchDB()
	if err != nil {
		return nil, err
	}
	msgs := convert.Convert_fbcustomerconversationsearchmodel_FbExternalMessageSearchs_fbcustomerconversationsearch_FbExternalMessageSearchs(msgDBs)
	return msgs, nil
}
