package sqlstore

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
)

type FbExternalCommentStoreFactory func(ctx context.Context) *FbExternalCommentStore

func NewFbExternalCommentStore(db *cmsql.Database) FbExternalCommentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalCommentStore {
		return &FbExternalCommentStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalCommentStore struct {
	ft FbExternalCommentFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalCommentStore) CreateFbExternalComment(fbExternalComment *fbmessaging.FbExternalComment) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalCommentDB := new(model.FbExternalComment)
	if err := scheme.Convert(fbExternalComment, fbExternalCommentDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbExternalCommentDB)
	if err != nil {
		return err
	}

	var tempFbExternalComment model.FbExternalComment
	if err := s.query().Where(s.ft.ByID(fbExternalComment.ID)).ShouldGet(&tempFbExternalComment); err != nil {
		return err
	}
	fbExternalComment.CreatedAt = tempFbExternalComment.CreatedAt
	fbExternalComment.UpdatedAt = tempFbExternalComment.UpdatedAt

	return nil
}

func (s *FbExternalCommentStore) CreateFbExternalComments(fbExternalComments []*fbmessaging.FbExternalComment) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalCommentsDB := model.FbExternalComments(convert.Convert_fbmessaging_FbExternalComments_fbmessagingmodel_FbExternalComments(fbExternalComments))

	_, err := s.query().Upsert(&fbExternalCommentsDB)
	if err != nil {
		return err
	}
	return nil
}
