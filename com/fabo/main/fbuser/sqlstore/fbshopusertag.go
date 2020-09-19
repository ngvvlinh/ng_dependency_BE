package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbShopTagStoreFactory func(ctx context.Context) *FbShopTagStore

type FbShopTagStore struct {
	db *cmsql.Database
	ft FbShopUserTagFilters

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

func (s *FbShopTagStore) ByIDs(ids []dot.ID) *FbShopTagStore {
	s.preds = append(s.preds, sq.In("id", ids))
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

func (s *FbShopTagStore) CreateShopUserTagDB(tag *model.FbShopUserTag) (*model.FbShopUserTag, error) {
	if err := s.query().ShouldInsert(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *FbShopTagStore) CreateShopUserTag(tag *fbusering.FbShopUserTag) error {
	sqlstore.MustNoPreds(s.preds)
	var tagDB *model.FbShopUserTag
	tagDB = convert.Convert_fbusering_FbShopUserTag_fbusermodel_FbShopUserTag(tag, tagDB)
	tagDB.ID = cm.NewID()
	tagDB, err := s.CreateShopUserTagDB(tagDB)
	if err != nil {
		return err
	}
	tag.ID = tagDB.ID
	return nil
}

func (s *FbShopTagStore) GetShopUserTagDB() (*model.FbShopUserTag, error) {
	query := s.query().Where(s.preds)

	tag := &model.FbShopUserTag{}
	err := query.ShouldGet(tag)
	return tag, err
}

func (s *FbShopTagStore) GetShopUserTag() (*fbusering.FbShopUserTag, error) {
	tagDB, err := s.GetShopUserTagDB()
	if err != nil {
		return nil, err
	}
	tag := &fbusering.FbShopUserTag{}
	tag = convert.Convert_fbusermodel_FbShopUserTag_fbusering_FbShopUserTag(tagDB, tag)
	return tag, err
}

func (s *FbShopTagStore) GetShopTagDBs() (model.FbShopUserTags, error) {
	query := s.query().Where(s.preds)
	var tags model.FbShopUserTags
	err := query.Find(&tags)
	return tags, err
}

func (s *FbShopTagStore) GetShopUserTags() ([]*fbusering.FbShopUserTag, error) {
	tagDBs, err := s.GetShopTagDBs()
	if err != nil {
		return nil, err
	}
	tags := []*fbusering.FbShopUserTag{}
	tags = convert.Convert_fbusermodel_FbShopUserTags_fbusering_FbShopUserTags(tagDBs)
	return tags, err
}

func (s *FbShopTagStore) UpdateShopUserTagDB(tag *model.FbShopUserTag) (*model.FbShopUserTag, error) {
	query := s.query().Where(s.preds)

	id := tag.ID
	tag.ID = 0 // not update ID
	err := query.ShouldUpdate(tag)
	if err != nil {
		return nil, err
	}

	return s.ByID(id).GetShopUserTagDB()
}

func (s *FbShopTagStore) UpdateShopUserTag(tag *fbusering.FbShopUserTag) error {
	tagDB := &model.FbShopUserTag{}
	tagDB = convert.Convert_fbusering_FbShopUserTag_fbusermodel_FbShopUserTag(tag, tagDB)

	tagDB, err := s.UpdateShopUserTagDB(tagDB)
	if err != nil {
		return err
	}

	tag = convert.Convert_fbusermodel_FbShopUserTag_fbusering_FbShopUserTag(tagDB, tag)
	return nil
}

func (s *FbShopTagStore) ListShopUserTagsDB() (model.FbShopUserTags, error) {
	query := s.query().Where(s.preds)
	var tags model.FbShopUserTags
	err := query.Find(&tags)
	return tags, err
}

func (s *FbShopTagStore) ListShopUserTags() ([]*fbusering.FbShopUserTag, error) {
	tagDBs, err := s.ListShopUserTagsDB()
	if err != nil {
		return nil, err
	}

	var tags []*fbusering.FbShopUserTag
	tags = convert.Convert_fbusermodel_FbShopUserTags_fbusering_FbShopUserTags(tagDBs)
	return tags, nil
}

func (s *FbShopTagStore) DeleteShopUserTag() error {
	return s.query().Where(s.preds).ShouldDelete(&model.FbShopUserTag{})
}
