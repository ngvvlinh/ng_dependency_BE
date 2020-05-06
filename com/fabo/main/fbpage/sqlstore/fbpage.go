package sqlstore

import (
	"context"

	"o.o/api/fabo/fbpaging"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/convert"
	"o.o/backend/com/fabo/main/fbpage/model"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type FbExternalPageStoreFactory func(ctx context.Context) *FbExternalPageStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewFbExternalPageStore(db *cmsql.Database) FbExternalPageStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbExternalPageStore {
		return &FbExternalPageStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbExternalPageStore struct {
	ft FbExternalPageFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbExternalPageStore) WithPaging(paging meta.Paging) *FbExternalPageStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbExternalPageStore) Filters(filters meta.Filters) *FbExternalPageStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *FbExternalPageStore) ID(id dot.ID) *FbExternalPageStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbExternalPageStore) IDs(IDs []dot.ID) *FbExternalPageStore {
	s.preds = append(s.preds, sq.In("id", IDs))
	return s
}

func (s *FbExternalPageStore) OptionalShopID(shopID dot.ID) *FbExternalPageStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID).Optional())
	return s
}

func (s *FbExternalPageStore) ShopID(shopID dot.ID) *FbExternalPageStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *FbExternalPageStore) FbUserID(fbUserID dot.ID) *FbExternalPageStore {
	s.preds = append(s.preds, s.ft.ByFbUserID(fbUserID))
	return s
}

func (s *FbExternalPageStore) Status(status status3.Status) *FbExternalPageStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *FbExternalPageStore) UserID(userID dot.ID) *FbExternalPageStore {
	s.preds = append(s.preds, s.ft.ByUserID(userID))
	return s
}

func (s *FbExternalPageStore) ExternalIDs(externalIDs []string) *FbExternalPageStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbExternalPageStore) UpdateStatus(status int) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	updateStatus, err := query.Table("fb_external_page").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return updateStatus, err
}

func (s *FbExternalPageStore) UpdateConnectionStatus(connectionStatus int) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	updateConnectionStatus, err := query.Table("fb_external_page").UpdateMap(map[string]interface{}{
		"connection_status": connectionStatus,
	})
	return updateConnectionStatus, err
}

func (s *FbExternalPageStore) UpdateStatusAndConnectionStatus(status, connectionStatus int) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	updateStatusAndConnectionStatus, err := query.Table("fb_external_page").UpdateMap(map[string]interface{}{
		"status":            status,
		"connection_status": connectionStatus,
	})
	return updateStatusAndConnectionStatus, err
}

func (s *FbExternalPageStore) CreateFbExternalPage(fbExternalPage *fbpaging.FbExternalPage) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPageDB := new(model.FbExternalPage)
	if err := scheme.Convert(fbExternalPage, fbExternalPageDB); err != nil {
		return err
	}

	_, err := s.query().Upsert(fbExternalPageDB)
	if err != nil {
		return err
	}

	var tempFbExternalPage model.FbExternalPage
	if err := s.query().Where(s.ft.ByID(fbExternalPage.ID)).ShouldGet(&tempFbExternalPage); err != nil {
		return err
	}
	fbExternalPage.CreatedAt = tempFbExternalPage.CreatedAt
	fbExternalPage.UpdatedAt = tempFbExternalPage.UpdatedAt

	return nil
}

func (s *FbExternalPageStore) CreateFbExternalPages(fbExternalPages []*fbpaging.FbExternalPage) error {
	sqlstore.MustNoPreds(s.preds)
	fbExternalPagesDB := model.FbExternalPages(convert.Convert_fbpaging_FbExternalPages_fbpagemodel_FbExternalPages(fbExternalPages))

	_, err := s.query().Upsert(&fbExternalPagesDB)
	if err != nil {
		return err
	}

	return nil
}

func (s *FbExternalPageStore) ListFbExternalPagesDB() ([]*model.FbExternalPage, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbExternalPage, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbExternalPage)
	if err != nil {
		return nil, err
	}

	var fbExternalPages model.FbExternalPages
	err = query.Find(&fbExternalPages)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbExternalPages)
	return fbExternalPages, nil
}

func (s *FbExternalPageStore) ListFbPages() (result []*fbpaging.FbExternalPage, err error) {
	fbExternalPages, err := s.ListFbExternalPagesDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbExternalPages, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbExternalPageStore) IncludeDeleted() *FbExternalPageStore {
	s.includeDeleted = true
	return s
}
