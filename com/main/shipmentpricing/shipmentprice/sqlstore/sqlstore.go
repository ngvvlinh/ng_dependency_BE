package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/shipmentpricing/shipmentprice"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice/convert"
	"etop.vn/backend/com/main/shipmentpricing/shipmentprice/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/wl"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type ShipmentPriceStore struct {
	ft    ShipmentPriceFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

type ShipmentPriceStoreFactory func(ctx context.Context) *ShipmentPriceStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewShipmentPriceStore(db *cmsql.Database) ShipmentPriceStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShipmentPriceStore {
		return &ShipmentPriceStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (s *ShipmentPriceStore) ID(id dot.ID) *ShipmentPriceStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShipmentPriceStore) Status(status status3.Status) *ShipmentPriceStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ShipmentPriceStore) ShipmentServiceID(id dot.ID) *ShipmentPriceStore {
	s.preds = append(s.preds, s.ft.ByShipmentServiceID(id))
	return s
}

func (s *ShipmentPriceStore) OptionalShipmentServiceID(id dot.ID) *ShipmentPriceStore {
	s.preds = append(s.preds, s.ft.ByShipmentServiceID(id).Optional())
	return s
}

func (s *ShipmentPriceStore) ShipmentPriceListID(id dot.ID) *ShipmentPriceStore {
	s.preds = append(s.preds, s.ft.ByShipmentPriceListID(id))
	return s
}

func (s *ShipmentPriceStore) OptionalShipmentPriceListID(id dot.ID) *ShipmentPriceStore {
	s.preds = append(s.preds, s.ft.ByShipmentPriceListID(id).Optional())
	return s
}

func (s *ShipmentPriceStore) GetShipmentPriceDB() (*model.ShipmentPrice, error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	var price model.ShipmentPrice
	err := query.ShouldGet(&price)
	return &price, err
}

func (s *ShipmentPriceStore) GetShipmentPrice() (*shipmentprice.ShipmentPrice, error) {
	priceDB, err := s.GetShipmentPriceDB()
	if err != nil {
		return nil, err
	}
	var res shipmentprice.ShipmentPrice
	if err := scheme.Convert(priceDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ShipmentPriceStore) ListShipmentPriceDBs() (res []*model.ShipmentPrice, err error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	err = query.Find((*model.ShipmentPrices)(&res))
	return
}

func (s *ShipmentPriceStore) ListShipmentPrices() (res []*shipmentprice.ShipmentPrice, _ error) {
	priceDBs, err := s.ListShipmentPriceDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(priceDBs, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ShipmentPriceStore) CreateShipmentPrice(price *shipmentprice.ShipmentPrice) (*shipmentprice.ShipmentPrice, error) {
	sqlstore.MustNoPreds(s.preds)
	if price.ID == 0 {
		price.ID = cm.NewID()
	}
	var priceDB model.ShipmentPrice
	if err := scheme.Convert(price, &priceDB); err != nil {
		return nil, err
	}
	priceDB.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	if err := s.query().ShouldInsert(&priceDB); err != nil {
		return nil, err
	}
	return s.ID(price.ID).GetShipmentPrice()
}

func (s *ShipmentPriceStore) UpdateShipmentPriceDB(price *model.ShipmentPrice) error {
	query := s.query().Where(s.ft.ByID(price.ID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(price)
}

func (s *ShipmentPriceStore) UpdateShipmentPrice(price *shipmentprice.ShipmentPrice) (*shipmentprice.ShipmentPrice, error) {
	sqlstore.MustNoPreds(s.preds)
	var priceDB model.ShipmentPrice
	if err := scheme.Convert(price, &priceDB); err != nil {
		return nil, err
	}
	if err := s.UpdateShipmentPriceDB(&priceDB); err != nil {
		return nil, err
	}
	return s.ID(price.ID).GetShipmentPrice()
}

func (s *ShipmentPriceStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	_deleted, err := query.Table("shipment_price").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ShipmentPriceStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
