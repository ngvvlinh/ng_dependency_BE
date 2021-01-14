package sqlstore

import (
	"context"

	"o.o/api/fabo/fbusering"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbuser/convert"
	"o.o/backend/com/fabo/main/fbuser/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalUserConnectedStoreFactory func(ctx context.Context) *FbExternalUserConnectedStore

func NewFbExternalUserConnectedStore(db *cmsql.Database) FbExternalUserConnectedStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalUserConnectedStore {
		return &FbExternalUserConnectedStore{
			db:    db,
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalUserConnectedStore struct {
	db *cmsql.Database
	ft FbExternalUserConnectedFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalUserConnectedStore) ExternalID(externalID string) *FbExternalUserConnectedStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *FbExternalUserConnectedStore) ExternalIDs(externalIDs []string) *FbExternalUserConnectedStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalUserConnectedStore) ShopID(shopID dot.ID) *FbExternalUserConnectedStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *FbExternalUserConnectedStore) WithPaging(paging meta.Paging) *FbExternalUserConnectedStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalUserConnectedStore) Filters(filters meta.Filters) *FbExternalUserConnectedStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *FbExternalUserConnectedStore) Status(status status3.Status) *FbExternalUserConnectedStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *FbExternalUserConnectedStore) UpdateStatus(status int) (int, error) {
	query := s.query().Where(s.preds)
	updateStatus, err := query.Table("fb_external_user_connected").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return updateStatus, err
}

func (s *FbExternalUserConnectedStore) CreateOrUpdateFbExternalUserConnected(fbExternalUserConnected *fbusering.FbExternalUserConnected) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalUserConnectedDB := new(model.FbExternalUserConnected)
	if err := scheme.Convert(fbExternalUserConnected, fbExternalUserConnectedDB); err != nil {
		return err
	}
	_, err := s.query().Upsert(fbExternalUserConnectedDB)
	if err != nil {
		return err
	}

	var tempFbExternalUserConnected model.FbExternalUserConnected
	if err := s.query().Where(s.ft.ByExternalID(fbExternalUserConnected.ExternalID)).ShouldGet(&tempFbExternalUserConnected); err != nil {
		return err
	}
	fbExternalUserConnected.CreatedAt = tempFbExternalUserConnected.CreatedAt
	fbExternalUserConnected.UpdatedAt = tempFbExternalUserConnected.UpdatedAt

	return nil
}

func (s *FbExternalUserConnectedStore) CreateFbExternalUserConnecteds(fbExternalUserConnecteds []*fbusering.FbExternalUserConnected) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalUserConnectedsDB := model.FbExternalUserConnecteds(convert.Convert_fbusering_FbExternalUserConnecteds_fbusermodel_FbExternalUserConnecteds(fbExternalUserConnecteds))

	if _, err := s.query().Upsert(&fbExternalUserConnectedsDB); err != nil {
		return err
	}
	return nil
}

func (s *FbExternalUserConnectedStore) GetFbExternalUserConnectedDB() (*model.FbExternalUserConnected, error) {
	query := s.query().Where(s.preds)

	var fbExternalUserConnected model.FbExternalUserConnected
	err := query.ShouldGet(&fbExternalUserConnected)
	return &fbExternalUserConnected, err
}

func (s *FbExternalUserConnectedStore) GetFbExternalUserConnected() (*fbusering.FbExternalUserConnected, error) {
	fbExternalUserConnected, err := s.GetFbExternalUserConnectedDB()
	if err != nil {
		return nil, err
	}

	result := &fbusering.FbExternalUserConnected{}
	if err = scheme.Convert(fbExternalUserConnected, result); err != nil {
		return nil, err
	}

	return result, err
}

func (s *FbExternalUserConnectedStore) ListFbExternalUserConnectedsDB() ([]*model.FbExternalUserConnected, error) {
	query := s.query().Where(s.preds)
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalUserConnected, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalUserConnected)
	if err != nil {
		return nil, err
	}

	var fbExternalUserConnecteds model.FbExternalUserConnecteds
	err = query.Find(&fbExternalUserConnecteds)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalUserConnecteds)
	return fbExternalUserConnecteds, nil
}

func (s *FbExternalUserConnectedStore) ListFbExternalUserConnecteds() (result []*fbusering.FbExternalUserConnected, err error) {
	fbExternalUserConnecteds, err := s.ListFbExternalUserConnectedsDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalUserConnecteds, &result); err != nil {
		return nil, err
	}
	return
}
