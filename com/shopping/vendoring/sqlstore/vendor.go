package sqlstore

import (
	"context"
	"time"

	"etop.vn/backend/pkg/common/scheme"

	"etop.vn/api/shopping/customering"

	"etop.vn/backend/com/shopping/vendoring/convert"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/vendoring"
	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/com/shopping/vendoring/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
)

type VendorStoreFactory func(ctx context.Context) *VendorStore

func NewVendorStore(db cmsql.Database) VendorStoreFactory {
	return func(ctx context.Context) *VendorStore {
		return &VendorStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type VendorStore struct {
	ft ShopVendorFilters

	query   func() cmsql.QueryInterface
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *VendorStore) Paging(paging meta.Paging) *VendorStore {
	s.paging = paging
	return s
}

func (s *VendorStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *VendorStore) Filters(filters meta.Filters) *VendorStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *VendorStore) ID(id int64) *VendorStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *VendorStore) IDs(ids ...int64) *VendorStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *VendorStore) ShopID(id int64) *VendorStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *VendorStore) OptionalShopID(id int64) *VendorStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *VendorStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopVendor)(nil))
}

func (s *VendorStore) CreateVendor(vendor *vendoring.ShopVendor) error {
	sqlstore.MustNoPreds(s.preds)
	trader := &customeringmodel.ShopTrader{
		ID:     vendor.ID,
		ShopID: vendor.ShopID,
		Type:   customering.VendorType,
	}
	vendorDB := new(model.ShopVendor)
	if err := scheme.Convert(vendor, vendorDB); err != nil {
		return err
	}
	if _, err := s.query().Insert(trader, vendorDB); err != nil {
		return err
	}

	var tempVendor model.ShopVendor
	if err := s.query().Where(s.ft.ByID(vendor.ID), s.ft.ByShopID(vendor.ShopID)).ShouldGet(&tempVendor); err != nil {
		return err
	}

	vendor.CreatedAt = tempVendor.CreatedAt
	vendor.UpdatedAt = tempVendor.UpdatedAt

	return nil
}

func (s *VendorStore) UpdateVendorDB(vendor *model.ShopVendor) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ft.ByID(vendor.ID)).UpdateAll().ShouldUpdate(vendor)
	return err
}

func (s *VendorStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_vendor").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *VendorStore) DeleteCustomer() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopVendor)(nil))
	return int(n), err
}

func (s *VendorStore) GetVendorDB() (*model.ShopVendor, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var vendor model.ShopVendor
	err := query.ShouldGet(&vendor)
	return &vendor, err
}

func (s *VendorStore) GetVendor() (vendorResult *vendoring.ShopVendor, _ error) {
	vendor, err := s.GetVendorDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_vendoringmodel_ShopVendor_vendoring_ShopVendor(vendor, vendorResult), nil
}

func (s *VendorStore) ListVendorsDB() ([]*model.ShopVendor, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortVendor)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterVendor)
	if err != nil {
		return nil, err
	}

	var vendors model.ShopVendors
	err = query.Find(&vendors)
	return vendors, err
}

func (s *VendorStore) ListVendors() ([]*vendoring.ShopVendor, error) {
	vendors, err := s.ListVendorsDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_vendoringmodel_ShopVendors_vendoring_ShopVendors(vendors), nil
}
