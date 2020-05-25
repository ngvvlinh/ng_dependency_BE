package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/shipmentpricing/shopshipmentpricelist"
	"o.o/api/meta"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist/convert"
	"o.o/backend/com/main/shipmentpricing/shopshipmentpricelist/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShopPriceListStore struct {
	ft    ShopShipmentPriceListFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context
	sqlstore.Paging

	includeDeleted sqlstore.IncludeDeleted
}

type ShopPriceListStoreFactory func(ctx context.Context) *ShopPriceListStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewShopPriceListStore(db *cmsql.Database) ShopPriceListStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShopPriceListStore {
		return &ShopPriceListStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (ft *ShopShipmentPriceListFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShopShipmentPriceListFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *ShopPriceListStore) WithPaging(paging meta.Paging) *ShopPriceListStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *ShopPriceListStore) ShopID(id dot.ID) *ShopPriceListStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *ShopPriceListStore) ShipmentPriceListID(id dot.ID) *ShopPriceListStore {
	s.preds = append(s.preds, s.ft.ByShipmentPriceListID(id))
	return s
}

func (s *ShopPriceListStore) ShipmentPriceListIDs(ids []dot.ID) *ShopPriceListStore {
	s.preds = append(s.preds, sq.In("shipment_price_list_id", ids))
	return s
}

func (s *ShopPriceListStore) OptionalShipmentPriceListID(id dot.ID) *ShopPriceListStore {
	s.preds = append(s.preds, s.ft.ByShipmentPriceListID(id).Optional())
	return s
}

func (s *ShopPriceListStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}

func (s *ShopPriceListStore) GetShopPriceListDB() (*model.ShopShipmentPriceList, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	var priceList model.ShopShipmentPriceList
	err := query.ShouldGet(&priceList)
	return &priceList, err
}

func (s *ShopPriceListStore) GetShopPriceList() (*shopshipmentpricelist.ShopShipmentPriceList, error) {
	priceListDB, err := s.GetShopPriceListDB()
	if err != nil {
		return nil, err
	}
	var res shopshipmentpricelist.ShopShipmentPriceList
	if err := scheme.Convert(priceListDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ShopPriceListStore) ListShopPriceListDBs() (res []*model.ShopShipmentPriceList, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = []string{"-created_at"}
	}
	query, err = sqlstore.LimitSort(query, &s.Paging, SortShopShipmentPriceList)
	if err != nil {
		return nil, err
	}
	err = query.Find((*model.ShopShipmentPriceLists)(&res))
	return
}

func (s *ShopPriceListStore) ListShopPriceLists() (res []*shopshipmentpricelist.ShopShipmentPriceList, err error) {
	priceLists, err := s.ListShopPriceListDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(priceLists, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ShopPriceListStore) CreateShopPriceList(priceList *shopshipmentpricelist.ShopShipmentPriceList) (*shopshipmentpricelist.ShopShipmentPriceList, error) {
	sqlstore.MustNoPreds(s.preds)
	var priceListDB model.ShopShipmentPriceList
	if err := scheme.Convert(priceList, &priceListDB); err != nil {
		return nil, err
	}
	if priceList.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop ID").WithMetap("sqlstore", "CreateShopPriceList")
	}
	priceListDB.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	if err := s.query().ShouldInsert(&priceListDB); err != nil {
		return nil, err
	}
	return s.ShopID(priceList.ShopID).GetShopPriceList()
}

func (s *ShopPriceListStore) UpdateShopPriceList(priceList *shopshipmentpricelist.ShopShipmentPriceList) error {
	sqlstore.MustNoPreds(s.preds)
	var priceListDB model.ShopShipmentPriceList
	if priceList.ShopID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing shop ID").WithMetap("sqlstore", "UpdateShopPriceList")
	}
	if err := scheme.Convert(priceList, &priceListDB); err != nil {
		return err
	}
	query := s.query().Where(s.ft.ByShopID(priceList.ShopID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(&priceListDB)
}

func (s *ShopPriceListStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shop_shipment_price_list").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}
