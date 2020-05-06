package sqlstore

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
)

type FbExternalPostStoreFactory func(ctx context.Context) *FbExternalPostStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewFbExternalPostStore(db *cmsql.Database) FbExternalPostStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalPostStore {
		return &FbExternalPostStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalPostStore struct {
	ft FbExternalPostFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalPostStore) CreateFbExternalPost(fbExternalPost *fbmessaging.FbExternalPost) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPostDB := new(model.FbExternalPost)
	if err := scheme.Convert(fbExternalPost, fbExternalPostDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbExternalPostDB)
	if err != nil {
		return err
	}

	var tempFbExternalPost model.FbExternalPost
	if err := s.query().Where(s.ft.ByID(fbExternalPost.ID)).ShouldGet(&tempFbExternalPost); err != nil {
		return err
	}
	fbExternalPost.CreatedAt = tempFbExternalPost.CreatedAt
	fbExternalPost.UpdatedAt = tempFbExternalPost.UpdatedAt

	return nil
}

func (s *FbExternalPostStore) CreateFbExternalPosts(fbExternalPosts []*fbmessaging.FbExternalPost) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPostsDB := model.FbExternalPosts(convert.Convert_fbmessaging_FbExternalPosts_fbmessagingmodel_FbExternalPosts(fbExternalPosts))

	_, err := s.query().Upsert(&fbExternalPostsDB)
	if err != nil {
		return err
	}
	return nil
}
