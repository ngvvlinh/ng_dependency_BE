package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/shipmentpricing/subpricelist"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipmentpricing/subpricelist/convert"
	"o.o/backend/com/main/shipmentpricing/subpricelist/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShipmentSubPriceListStore struct {
	ft    ShipmentSubPriceListFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

type ShipmentSubPriceListStoreFactory func(ctx context.Context) *ShipmentSubPriceListStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewShipmentSubPriceListStore(db *cmsql.Database) ShipmentSubPriceListStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShipmentSubPriceListStore {
		return &ShipmentSubPriceListStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (ft *ShipmentSubPriceListFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShipmentSubPriceListFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *ShipmentSubPriceListStore) ID(id dot.ID) *ShipmentSubPriceListStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShipmentSubPriceListStore) IDs(ids ...dot.ID) *ShipmentSubPriceListStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShipmentSubPriceListStore) OptionalConnectionID(connID dot.ID) *ShipmentSubPriceListStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID).Optional())
	return s
}

func (s *ShipmentSubPriceListStore) Status(status status3.Status) *ShipmentSubPriceListStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ShipmentSubPriceListStore) GetShipmentSubPriceListDB() (*model.ShipmentSubPriceList, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	var priceList model.ShipmentSubPriceList
	err := query.ShouldGet(&priceList)
	return &priceList, err
}

func (s *ShipmentSubPriceListStore) GetShipmentSubPriceList() (*subpricelist.ShipmentSubPriceList, error) {
	serviceDB, err := s.GetShipmentSubPriceListDB()
	if err != nil {
		return nil, err
	}
	var res subpricelist.ShipmentSubPriceList
	if err := scheme.Convert(serviceDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ShipmentSubPriceListStore) ListShipmentSubPriceListDBs() (res []*model.ShipmentSubPriceList, err error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	err = query.Find((*model.ShipmentSubPriceLists)(&res))
	return
}

func (s *ShipmentSubPriceListStore) ListShipmentSubPriceLists() (res []*subpricelist.ShipmentSubPriceList, _ error) {
	subPriceLists, err := s.ListShipmentSubPriceListDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(subPriceLists, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ShipmentSubPriceListStore) CreateShipmentSubPriceList(subPriceList *subpricelist.ShipmentSubPriceList) (*subpricelist.ShipmentSubPriceList, error) {
	sqlstore.MustNoPreds(s.preds)
	if subPriceList.ID == 0 {
		subPriceList.ID = cm.NewID()
	}
	subPriceList.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	var subPriceListDB model.ShipmentSubPriceList
	if err := scheme.Convert(subPriceList, &subPriceListDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&subPriceListDB); err != nil {
		return nil, err
	}
	return s.ID(subPriceList.ID).GetShipmentSubPriceList()
}

func (s *ShipmentSubPriceListStore) UpdateShipmentSubPriceList(subPriceList *subpricelist.ShipmentSubPriceList) error {
	sqlstore.MustNoPreds(s.preds)
	var subPriceListDB model.ShipmentSubPriceList
	if err := scheme.Convert(subPriceList, &subPriceListDB); err != nil {
		return err
	}
	query := s.query().Where(s.ft.ByID(subPriceList.ID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(&subPriceListDB)
}

func (s *ShipmentSubPriceListStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_sub_price_list").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *ShipmentSubPriceListStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
