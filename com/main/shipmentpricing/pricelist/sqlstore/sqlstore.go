package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/shipmentpricing/pricelist"
	"o.o/backend/com/main/shipmentpricing/pricelist/convert"
	"o.o/backend/com/main/shipmentpricing/pricelist/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShipmentPriceListStore struct {
	ft    ShipmentPriceListFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

type ShipmentPriceListStoreFactory func(ctx context.Context) *ShipmentPriceListStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewShipmentPriceListStore(db *cmsql.Database) ShipmentPriceListStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShipmentPriceListStore {
		return &ShipmentPriceListStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (ft *ShipmentPriceListFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShipmentPriceListFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *ShipmentPriceListStore) ID(id dot.ID) *ShipmentPriceListStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShipmentPriceListStore) IsActive(isActive bool) *ShipmentPriceListStore {
	s.preds = append(s.preds, s.ft.ByIsActive(isActive))
	return s
}

func (s *ShipmentPriceListStore) SubPriceListIDs(subPriceListIDs ...dot.ID) *ShipmentPriceListStore {
	s.preds = append(s.preds, sq.NewExpr("shipment_sub_price_list_ids && ?", core.Array{V: subPriceListIDs}))
	return s
}

func (s *ShipmentPriceListStore) GetShipmentPriceListDB() (*model.ShipmentPriceList, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	var priceList model.ShipmentPriceList
	err := query.ShouldGet(&priceList)
	return &priceList, err
}

func (s *ShipmentPriceListStore) GetShipmentPriceList() (*pricelist.ShipmentPriceList, error) {
	priceListDB, err := s.GetShipmentPriceListDB()
	if err != nil {
		return nil, err
	}
	var res pricelist.ShipmentPriceList
	if err := scheme.Convert(priceListDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ShipmentPriceListStore) ListShipmentPriceListDBs() (res []*model.ShipmentPriceList, err error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	err = query.Find((*model.ShipmentPriceLists)(&res))
	return
}

func (s *ShipmentPriceListStore) ListShipmentPriceLists() (res []*pricelist.ShipmentPriceList, _ error) {
	priceLists, err := s.ListShipmentPriceListDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(priceLists, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ShipmentPriceListStore) CreateShipmentPriceList(priceList *pricelist.ShipmentPriceList) (*pricelist.ShipmentPriceList, error) {
	sqlstore.MustNoPreds(s.preds)
	if priceList.ID == 0 {
		priceList.ID = cm.NewID()
	}
	priceList.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	var priceListDB model.ShipmentPriceList
	if err := scheme.Convert(priceList, &priceListDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&priceListDB); err != nil {
		return nil, err
	}
	return s.ID(priceList.ID).GetShipmentPriceList()
}

func (s *ShipmentPriceListStore) UpdateShipmentPriceList(priceList *pricelist.ShipmentPriceList) error {
	sqlstore.MustNoPreds(s.preds)
	var priceListDB model.ShipmentPriceList
	if err := scheme.Convert(priceList, &priceListDB); err != nil {
		return err
	}
	query := s.query().Where(s.ft.ByID(priceList.ID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(&priceListDB)
}

func (s *ShipmentPriceListStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_price_list").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *ShipmentPriceListStore) ActivePriceList() error {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_price_list").ShouldUpdateMap(map[string]interface{}{
		"is_active": true,
	})
}

func (s *ShipmentPriceListStore) DeactivePriceList() error {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_price_list").ShouldUpdateMap(map[string]interface{}{
		"is_active": false,
	})
}

func (s *ShipmentPriceListStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
