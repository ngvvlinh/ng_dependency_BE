package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion/convert"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type PriceListPromotionStore struct {
	ft    ShipmentPriceListPromotionFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	sqlstore.Paging

	ctx            context.Context
	includeDeleted sqlstore.IncludeDeleted
}

var SortShipmentPricePromotion = map[string]string{
	"created_at":     "created_at",
	"priority_point": "priority_point",
}

type PriceListStorePromotionFactory func(ctx context.Context) *PriceListPromotionStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewPriceListStorePromotion(db *cmsql.Database) PriceListStorePromotionFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *PriceListPromotionStore {
		return &PriceListPromotionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (ft *ShipmentPriceListPromotionFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShipmentPriceListPromotionFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *PriceListPromotionStore) WithPaging(paging meta.Paging) *PriceListPromotionStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *PriceListPromotionStore) ID(id dot.ID) *PriceListPromotionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PriceListPromotionStore) OptionalConnectionID(connID dot.ID) *PriceListPromotionStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID).Optional())
	return s
}

func (s *PriceListPromotionStore) OptionalPriceListID(priceListID dot.ID) *PriceListPromotionStore {
	s.preds = append(s.preds, s.ft.ByPriceListID(priceListID).Optional())
	return s
}

func (s *PriceListPromotionStore) ConnectionID(connID dot.ID) *PriceListPromotionStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID))
	return s
}

func (s *PriceListPromotionStore) Status(status status3.Status) *PriceListPromotionStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *PriceListPromotionStore) ActiveDate(now time.Time) *PriceListPromotionStore {
	s.preds = append(s.preds, sq.NewExpr("date_from < ? AND date_to >= ?", now, now))
	return s
}

func (s *PriceListPromotionStore) GetShipmentPriceListPromotionDB() (*model.ShipmentPriceListPromotion, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	var priceListPromotion model.ShipmentPriceListPromotion
	err := query.ShouldGet(&priceListPromotion)
	return &priceListPromotion, err
}

func (s *PriceListPromotionStore) GetShipmentPriceListPromotion() (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	priceListPromotionDB, err := s.GetShipmentPriceListPromotionDB()
	if err != nil {
		return nil, err
	}
	var res pricelistpromotion.ShipmentPriceListPromotion
	if err := scheme.Convert(priceListPromotionDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *PriceListPromotionStore) ListShipmentPriceListPromotionDBs() (res []*model.ShipmentPriceListPromotion, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
	query, err = sqlstore.LimitSort(query, &s.Paging, SortShipmentPricePromotion)
	if err != nil {
		return nil, err
	}

	if wl.X(s.ctx).IsWLPartnerPOS() {
		query = query.Where(s.ft.NotBelongWLPartner())
	} else {
		query = s.ByWhiteLabelPartner(s.ctx, query)
	}
	err = query.Find((*model.ShipmentPriceListPromotions)(&res))
	return
}

func (s *PriceListPromotionStore) ListShipmentPriceListPromotions() (res []*pricelistpromotion.ShipmentPriceListPromotion, _ error) {
	priceListPromotions, err := s.ListShipmentPriceListPromotionDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(priceListPromotions, &res); err != nil {
		return nil, err
	}
	return
}

func (s *PriceListPromotionStore) CreateShipmentPriceListPromotion(priceListPromotion *pricelistpromotion.ShipmentPriceListPromotion) (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	sqlstore.MustNoPreds(s.preds)
	if priceListPromotion.ID == 0 {
		priceListPromotion.ID = cm.NewID()
	}
	priceListPromotion.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	var priceListPromotionDB model.ShipmentPriceListPromotion
	if err := scheme.Convert(priceListPromotion, &priceListPromotionDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&priceListPromotionDB); err != nil {
		return nil, err
	}
	return s.ID(priceListPromotion.ID).GetShipmentPriceListPromotion()
}

func (s *PriceListPromotionStore) UpdateShipmentPriceListPromotion(priceListPromotion *pricelistpromotion.ShipmentPriceListPromotion) error {
	sqlstore.MustNoPreds(s.preds)
	var priceListPromotionDB model.ShipmentPriceListPromotion
	if err := scheme.Convert(priceListPromotion, &priceListPromotionDB); err != nil {
		return err
	}
	query := s.query().Where(s.ft.ByID(priceListPromotion.ID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.ShouldUpdate(&priceListPromotionDB)
}

func (s *PriceListPromotionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_price_list_promotion").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *PriceListPromotionStore) SetUndefaultPriceList() error {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_price_list_promotion").ShouldUpdateMap(map[string]interface{}{
		"is_default": false,
	})
}

func (s *PriceListPromotionStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
