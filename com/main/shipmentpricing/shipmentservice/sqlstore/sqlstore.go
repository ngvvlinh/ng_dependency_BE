package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/shipmentpricing/shipmentservice"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipmentpricing/shipmentservice/convert"
	"o.o/backend/com/main/shipmentpricing/shipmentservice/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sq/core"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShipmentServiceStore struct {
	ft    ShipmentServiceFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

	includeDeleted sqlstore.IncludeDeleted
}

type ShipmentServiceStoreFactory func(ctx context.Context) *ShipmentServiceStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewShipmentServiceStore(db *cmsql.Database) ShipmentServiceStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *ShipmentServiceStore {
		return &ShipmentServiceStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
			ctx: ctx,
		}
	}
}

func (ft *ShipmentServiceFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft *ShipmentServiceFilters) NotBelongWLPartner() sq.WriterTo {
	return ft.Filter("$.wl_partner_id IS NULL")
}

func (s *ShipmentServiceStore) ID(id dot.ID) *ShipmentServiceStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ShipmentServiceStore) IDs(ids ...dot.ID) *ShipmentServiceStore {
	s.preds = append(s.preds, sq.In("id", ids))
	return s
}

func (s *ShipmentServiceStore) Status(status status3.Status) *ShipmentServiceStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ShipmentServiceStore) ServiceID(serviceID string) *ShipmentServiceStore {
	s.preds = append(s.preds, sq.NewExpr("service_ids @> ?", core.Array{
		V: []string{serviceID},
	}))
	return s
}

func (s *ShipmentServiceStore) ConnectionID(connID dot.ID) *ShipmentServiceStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID))
	return s
}

func (s *ShipmentServiceStore) OptionalConnectionID(connID dot.ID) *ShipmentServiceStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connID).Optional())
	return s
}

func (s *ShipmentServiceStore) GetShipmentServiceDB() (*model.ShipmentService, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = query.OrderBy("created_at DESC")
	var service model.ShipmentService
	err := query.ShouldGet(&service)
	return &service, err
}

func (s *ShipmentServiceStore) GetShipmentService() (*shipmentservice.ShipmentService, error) {
	serviceDB, err := s.GetShipmentServiceDB()
	if err != nil {
		return nil, err
	}
	var res shipmentservice.ShipmentService
	if err := scheme.Convert(serviceDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ShipmentServiceStore) ListShipmentServiceDBs() (res []*model.ShipmentService, err error) {
	query := s.query().Where(s.preds).OrderBy("created_at DESC")
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	if wl.X(s.ctx).IsWLPartnerPOS() {
		query = query.Where(s.ft.NotBelongWLPartner())
	} else {
		query = s.ByWhiteLabelPartner(s.ctx, query)
	}
	err = query.Find((*model.ShipmentServices)(&res))
	return
}

func (s *ShipmentServiceStore) ListShipmentServices() (res []*shipmentservice.ShipmentService, _ error) {
	services, err := s.ListShipmentServiceDBs()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(services, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ShipmentServiceStore) CreateShipmentService(service *shipmentservice.ShipmentService) (*shipmentservice.ShipmentService, error) {
	sqlstore.MustNoPreds(s.preds)
	if service.ID == 0 {
		service.ID = cm.NewID()
	}
	service.WLPartnerID = wl.GetWLPartnerID(s.ctx)
	var serviceDB model.ShipmentService
	if err := scheme.Convert(service, &serviceDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&serviceDB); err != nil {
		return nil, err
	}
	return s.ID(service.ID).GetShipmentService()
}

func (s *ShipmentServiceStore) UpdateShipmentServiceDB(service *model.ShipmentService) (updated int, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	updated, err = query.Update(service)
	return
}

func (s *ShipmentServiceStore) UpdateShipmentService(service *shipmentservice.ShipmentService) error {
	sqlstore.MustNoPreds(s.preds)
	var serviceDB model.ShipmentService
	if err := scheme.Convert(service, &serviceDB); err != nil {
		return err
	}
	_, err := s.ID(service.ID).UpdateShipmentServiceDB(&serviceDB)
	return err
}

func (s *ShipmentServiceStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	return query.Table("shipment_service").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
}

func (s *ShipmentServiceStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
