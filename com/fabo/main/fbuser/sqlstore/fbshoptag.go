package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbShopTagStoreFactory func(ctx context.Context) *FbShopTagStore

type FbShopTagStore struct {
	db *cmsql.Database
	ft FbShopTagFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging
}

func NewFbShopTagStore(db *cmsql.Database) FbShopTagStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbShopTagStore {
		return &FbShopTagStore{
			db:    db,
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

func (s *FbShopTagStore) ByID(id dot.ID) *FbShopTagStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbShopTagStore) ByShopID(id dot.ID) *FbShopTagStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *FbShopTagStore) ByName(name string) *FbShopTagStore {
	s.preds = append(s.preds, s.ft.ByName(name))
	return s
}

func (s *FbShopTagStore) CreateShopTagDB(tag *model.FbShopTag) (*model.FbShopTag, error) {
	if err := s.query().ShouldInsert(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *FbShopTagStore) CreateShopTag(tag *fbusering.FbShopTag) error {
	sqlstore.MustNoPreds(s.preds)
	var tagDB *model.FbShopTag
	tagDB = convert.Convert_fbusering_FbShopTag_fbusermodel_FbShopTag(tag, tagDB)
	tagDB.ID = cm.NewID()
	tagDB, err := s.CreateShopTagDB(tagDB)
	if err != nil {
		return err
	}
	tag.ID = tagDB.ID
	return nil
}

func (s *FbShopTagStore) GetShopTagDB() (*model.FbShopTag, error) {
	query := s.query().Where(s.preds)

	tag := &model.FbShopTag{}
	err := query.ShouldGet(tag)
	return tag, err
}

func (s *FbShopTagStore) GetShopTag() (*fbusering.FbShopTag, error) {
	tagDB, err := s.GetShopTagDB()
	if err != nil {
		return nil, err
	}
	tag := &fbusering.FbShopTag{}
	tag = convert.Convert_fbusermodel_FbShopTag_fbusering_FbShopTag(tagDB, tag)
	return tag, err
}

func (s *FbShopTagStore) UpdateShopTagDB(tag *model.FbShopTag) (*model.FbShopTag, error) {
	query := s.query().Where(s.preds)

	id := tag.ID
	tag.ID = 0 // not update ID
	err := query.ShouldUpdate(tag)
	if err != nil {
		return nil, err
	}

	return s.ByID(id).GetShopTagDB()
}

func (s *FbShopTagStore) UpdateShopTag(tag *fbusering.FbShopTag) error {
	tagDB := &model.FbShopTag{}
	tagDB = convert.Convert_fbusering_FbShopTag_fbusermodel_FbShopTag(tag, tagDB)

	tagDB, err := s.UpdateShopTagDB(tagDB)
	if err != nil {
		return err
	}

	tag = convert.Convert_fbusermodel_FbShopTag_fbusering_FbShopTag(tagDB, tag)
	return nil
}

func (s *FbShopTagStore) ListShopTagDB() (model.FbShopTags, error) {
	query := s.query().Where(s.preds)
	var tags model.FbShopTags
	err := query.Find(&tags)
	return tags, err
}

func (s *FbShopTagStore) ListShopTag() ([]*fbusering.FbShopTag, error) {
	tagDBs, err := s.ListShopTagDB()
	if err != nil {
		return nil, err
	}

	var tags []*fbusering.FbShopTag
	tags = convert.Convert_fbusermodel_FbShopTags_fbusering_FbShopTags(tagDBs)
	return tags, nil
}

func (s *FbShopTagStore) DeleteShopTag() error {
	return s.query().Where(s.preds).ShouldDelete(&model.FbShopTag{})
}
