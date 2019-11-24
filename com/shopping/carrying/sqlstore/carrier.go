package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/meta"
	"etop.vn/api/shopping/carrying"
	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/com/shopping/carrying/convert"
	"etop.vn/backend/com/shopping/carrying/model"
	customeringmodel "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sqlstore"
	"etop.vn/capi/dot"
)

type CarrierStoreFactory func(context.Context) *CarrierStore

func NewCarrierStore(db *cmsql.Database) CarrierStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *CarrierStore {
		return &CarrierStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type CarrierStore struct {
	ft ShopCarrierFilters

	query   cmsql.QueryFactory
	preds   []interface{}
	filters meta.Filters
	paging  meta.Paging

	includeDeleted sqlstore.IncludeDeleted
}

func (s *CarrierStore) Paging(paging meta.Paging) *CarrierStore {
	s.paging = paging
	return s
}

func (s *CarrierStore) GetPaging() meta.PageInfo {
	return meta.FromPaging(s.paging)
}

func (s *CarrierStore) Filters(filters meta.Filters) *CarrierStore {
	if s.filters == nil {
		s.filters = filters
	} else {
		s.filters = append(s.filters, filters...)
	}
	return s
}

func (s *CarrierStore) ID(id dot.ID) *CarrierStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *CarrierStore) IDs(ids ...dot.ID) *CarrierStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *CarrierStore) ShopID(id dot.ID) *CarrierStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *CarrierStore) OptionalShopID(id dot.ID) *CarrierStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *CarrierStore) Count() (uint64, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	return query.Count((*model.ShopCarrier)(nil))
}

func (s *CarrierStore) CreateCarrier(carrier *carrying.ShopCarrier) error {
	sqlstore.MustNoPreds(s.preds)
	trader := &customeringmodel.ShopTrader{
		ID:     carrier.ID,
		ShopID: carrier.ShopID,
		Type:   tradering.CarrierType,
	}
	carrierDB := new(model.ShopCarrier)
	if err := scheme.Convert(carrier, carrierDB); err != nil {
		return err
	}
	if _, err := s.query().Insert(trader, carrierDB); err != nil {
		return err
	}

	var tempCarrier model.ShopCarrier
	if err := s.query().Where(s.ft.ByID(carrier.ID), s.ft.ByShopID(carrier.ShopID)).ShouldGet(&tempCarrier); err != nil {
		return err
	}

	carrier.CreatedAt = tempCarrier.CreatedAt
	carrier.UpdatedAt = tempCarrier.UpdatedAt
	return nil
}

func (s *CarrierStore) UpdateCarrierDB(carrier *model.ShopCarrier) error {
	sqlstore.MustNoPreds(s.preds)
	err := s.query().Where(s.ft.ByID(carrier.ID)).UpdateAll().ShouldUpdate(carrier)
	return err
}

func (s *CarrierStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_carrier").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return int(_deleted), err
}

func (s *CarrierStore) DeleteCarrier() (int, error) {
	n, err := s.query().Where(s.preds).Delete((*model.ShopCarrier)(nil))
	return int(n), err
}

func (s *CarrierStore) GetCarrierDB() (*model.ShopCarrier, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var carrier model.ShopCarrier
	err := query.ShouldGet(&carrier)
	return &carrier, err
}

func (s *CarrierStore) GetCarrier() (carrierResult *carrying.ShopCarrier, _ error) {
	carrier, err := s.GetCarrierDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_carryingmodel_ShopCarrier_carrying_ShopCarrier(carrier, carrierResult), nil
}

func (s *CarrierStore) ListCarriersDB() ([]*model.ShopCarrier, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query, err := sqlstore.LimitSort(query, &s.paging, SortCarrier)
	if err != nil {
		return nil, err
	}
	query, _, err = sqlstore.Filters(query, s.filters, FilterCarrier)
	if err != nil {
		return nil, err
	}

	var carriers model.ShopCarriers
	err = query.Find(&carriers)
	return carriers, err
}

func (s *CarrierStore) ListCarriers() ([]*carrying.ShopCarrier, error) {
	carriers, err := s.ListCarriersDB()
	if err != nil {
		return nil, err
	}
	return convert.Convert_carryingmodel_ShopCarriers_carrying_ShopCarriers(carriers), nil
}
