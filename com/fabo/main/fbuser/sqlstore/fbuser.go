package sqlstore

import (
	"context"

	"etop.vn/api/fabo/fbusering"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/fabo/main/fbuser/convert"
	"etop.vn/backend/com/fabo/main/fbuser/model"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type FbUserStoreFactory func(ctx context.Context) *FbUserStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewFbUserStore(db *cmsql.Database) FbUserStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbUserStore {
		return &FbUserStore{
			db:    db,
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbUserStore struct {
	db *cmsql.Database
	ft FbUserFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbUserStore) ExternalID(externalID string) *FbUserStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbUserStore) UserID(userID dot.ID) *FbUserStore {
	s.preds = append(s.preds, s.ft.ByUserID(userID))
	return s
}

func (s *FbUserStore) Status(status status3.Status) *FbUserStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *FbUserStore) UpdateStatus(status int) (int, error) {
	query := s.query().Where(s.preds)
	updateStatus, err := query.Table("fb_user").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return updateStatus, err
}

func (s *FbUserStore) CreateFbUser(fbUser *fbusering.FbUser) error {
	sqlstore.MustNoPreds(s.preds)
	fbUserDB := new(model.FbUser)
	if err := scheme.Convert(fbUser, fbUserDB); err != nil {
		return err
	}
	_, err := s.query().Upsert(fbUserDB)
	if err != nil {
		return err
	}

	var tempFbUser model.FbUser
	if err := s.query().Where(s.ft.ByID(fbUser.ID)).ShouldGet(&tempFbUser); err != nil {
		return err
	}
	fbUser.CreatedAt = tempFbUser.CreatedAt
	fbUser.UpdatedAt = tempFbUser.UpdatedAt

	return nil
}

func (s *FbUserStore) GetFbUserDB() (*model.FbUser, error) {
	query := s.query().Where(s.preds)

	var fbUser model.FbUser
	err := query.ShouldGet(&fbUser)
	return &fbUser, err
}

func (s *FbUserStore) GetFbUser() (*fbusering.FbUser, error) {
	fbuser, err := s.GetFbUserDB()
	if err != nil {
		return nil, err
	}
	result := &fbusering.FbUser{}
	err = scheme.Convert(fbuser, result)
	if err != nil {
		return nil, err
	}
	return result, err
}
