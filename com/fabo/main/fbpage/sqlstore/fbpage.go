package sqlstore

import (
	"context"

	"etop.vn/api/fabo/fbpaging"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/fabo/main/fbpage/convert"
	"etop.vn/backend/com/fabo/main/fbpage/model"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type FbPageStoreFactory func(ctx context.Context) *FbPageStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewFbPageStore(db *cmsql.Database) FbPageStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FbPageStore {
		return &FbPageStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FbPageStore struct {
	ft FbPageFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *FbPageStore) WithPaging(paging meta.Paging) *FbPageStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FbPageStore) Filters(filters meta.Filters) *FbPageStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *FbPageStore) ID(id dot.ID) *FbPageStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FbPageStore) IDs(IDs []dot.ID) *FbPageStore {
	s.preds = append(s.preds, sq.In("id", IDs))
	return s
}

func (s *FbPageStore) ShopID(shopID dot.ID) *FbPageStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *FbPageStore) FbUserID(fbUserID dot.ID) *FbPageStore {
	s.preds = append(s.preds, s.ft.ByFbUserID(fbUserID))
	return s
}

func (s *FbPageStore) Status(status status3.Status) *FbPageStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *FbPageStore) UserID(userID dot.ID) *FbPageStore {
	s.preds = append(s.preds, s.ft.ByUserID(userID))
	return s
}

func (s *FbPageStore) ExternalIDs(externalIDs []string) *FbPageStore {
	s.preds = append(s.preds, sq.In("external_id", externalIDs))
	return s
}

func (s *FbPageStore) UpdateStatus(status int) (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	updateStatus, err := query.Table("fb_page").UpdateMap(map[string]interface{}{
		"status": status,
	})
	return updateStatus, err
}

func (s *FbPageStore) CreateFbPage(fbPage *fbpaging.FbPage) error {
	sqlstore.MustNoPreds(s.preds)
	fbPageDB := new(model.FbPage)
	if err := scheme.Convert(fbPage, fbPageDB); err != nil {
		return err
	}

	_, err := s.query().Insert(fbPageDB)
	if err != nil {
		return err
	}

	var tempFbPage model.FbPage
	if err := s.query().Where(s.ft.ByID(fbPage.ID)).ShouldGet(&tempFbPage); err != nil {
		return err
	}
	fbPage.CreatedAt = tempFbPage.CreatedAt
	fbPage.UpdatedAt = tempFbPage.UpdatedAt

	return nil
}

func (s *FbPageStore) CreateFbPages(fbPages []*fbpaging.FbPage) error {
	sqlstore.MustNoPreds(s.preds)
	fbPagesDB := model.FbPages(convert.Convert_fbpaging_FbPages_fbpagemodel_FbPages(fbPages))

	_, err := s.query().Upsert(&fbPagesDB)
	if err != nil {
		return err
	}

	return nil
}

func (s *FbPageStore) ListFbPagesDB() ([]*model.FbPage, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if !s.Paging.IsCursorPaging() && len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFbPage, s.ft.prefix)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterFbPage)
	if err != nil {
		return nil, err
	}

	var fbPages model.FbPages
	err = query.Find(&fbPages)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(fbPages)
	return fbPages, nil
}

func (s *FbPageStore) ListFbPages() (result []*fbpaging.FbPage, err error) {
	fbPages, err := s.ListFbPagesDB()
	if err != nil {
		return nil, err
	}
	if err = scheme.Convert(fbPages, &result); err != nil {
		return nil, err
	}
	return
}

func (s *FbPageStore) IncludeDeleted() *FbPageStore {
	s.includeDeleted = true
	return s
}
