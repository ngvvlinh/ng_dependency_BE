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

type FbExternalCommentSearchStoreFactory func(ctx context.Context) *FbExternalCommentSearchStore

func NewFbExternalCommentSearchStore(db *cmsql.Database) FbExternalCommentSearchStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalCommentSearchStore {
		return &FbExternalCommentSearchStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalCommentSearchStore struct {
	ft FbExternalCommentSearchFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
}

func (s *FbExternalCommentSearchStore) ByPageIDs(pageIDs []string) *FbExternalCommentSearchStore {
	s.preds = append(s.preds, sq.In("external_page_id", pageIDs))
	return s
}

func (s *FbExternalCommentSearchStore) ByExternalMessageNorm(msg string) *FbExternalCommentSearchStore {
	s.preds = append(s.preds, sq.NewExpr("external_message_norm @@ ?::tsquery", validate.NormalizeSearchQueryAnd(msg)))
	return s
}

func (s *FbExternalCommentSearchStore) ListExternalCommentSearchDB() ([]*model.FbExternalCommentSearch, error) {
	var cmts model.FbExternalCommentSearchs
	if err := s.query().Where(s.preds).OrderBy("created_at desc").Find(&cmts); err != nil {
		return nil, err
	}

	return cmts, nil
}

func (s *FbExternalCommentSearchStore) ListExternalCommentSearch() ([]*fbcustomerconversationsearch.FbExternalCommentSearch, error) {
	cmtsDBs, err := s.ListExternalCommentSearchDB()
	if err != nil {
		return nil, err
	}

	return convert.Convert_fbcustomerconversationsearchmodel_FbExternalCommentSearchs_fbcustomerconversationsearch_FbExternalCommentSearchs(cmtsDBs), nil
}
