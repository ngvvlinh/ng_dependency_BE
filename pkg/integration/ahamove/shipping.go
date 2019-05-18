package ahamove

import (
	"context"

	"etop.vn/api/main/location"
	"etop.vn/api/main/shipnow/carrier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/logic/shipping_provider"
	"etop.vn/backend/pkg/etop/model"
	ahamoveclient "etop.vn/backend/pkg/integration/ahamove/client"
	shipnow_carrier "etop.vn/backend/pkg/services/shipnow-carrier"
)

var _ shipnow_carrier.ShipnowCarrier = &Carrier{}

const (
	AhamoveCodePublic = 'D'
)

type Carrier struct {
	clients  map[byte]*ahamoveclient.Client
	location location.Bus
}

func New(cfg Config, locationBus location.Bus) *Carrier {
	clientDefault := ahamoveclient.New(cfg.Env, cfg.AccountDefault)
	clients := map[byte]*ahamoveclient.Client{
		AhamoveCodePublic: clientDefault,
	}
	return &Carrier{
		clients:  clients,
		location: locationBus,
	}
}

func (c *Carrier) InitAllClients(ctx context.Context) error {
	for name, client := range c.clients {
		if err := client.Ping(); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "can not init client").
				WithMetap("client", name)
		}
	}
	return nil
}

func (p *Carrier) CreateExternalShipnow(ctx context.Context, cmd *carrier.CreateExternalShipnowCommand) (_err error) {
	return cm.ErrTODO
}

func (p *Carrier) CancelExternalShipnow(ctx context.Context, ffm *carrier.CancelExternalShipnowCommand) error {
	// code := ffm.ExternalShippingCode
	// cmd := &CancelOrderCommand{
	// 	ServiceID: ffm.ProviderServiceID,
	// 	LabelID:   code,
	// }
	// err := p.CancelOrder(ctx, cmd)
	// return err
	return nil
}

func (p *Carrier) GetShippingServices(ctx context.Context, args shipnow_carrier.GetShippingServiceArgs) ([]*carrier.AvailableShippingService, error) {

	// fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	// toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	// if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
	// 	return nil, err
	// }
	// fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	// toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	//
	// cmd := &CalcShippingFeeCommand{
	// 	ArbitraryID:      args.AccountID,
	// 	FromDistrictCode: args.FromDistrictCode,
	// 	ToDistrictCode:   args.ToDistrictCode,
	// 	Request: &ahamoveclient.CalcShippingFeeRequest{
	// 		Weight:          args.ChargeableWeight,
	// 		Value:           args.GetInsuranceAmount(),
	// 		PickingProvince: fromProvince.Name,
	// 		PickingDistrict: fromDistrict.Name,
	// 		Province:        toProvince.Name,
	// 		District:        toDistrict.Name,
	// 	},
	// }
	// err := p.CalcShippingFee(ctx, cmd)
	// services := cmd.Result
	// return services, err
	return nil, cm.ErrTODO
}

func (p *Carrier) GetAllShippingServices(ctx context.Context, args shipping_provider.GetShippingServicesArgs) ([]*model.AvailableShippingService, error) {
	// fromQuery := &location.GetLocationQuery{DistrictCode: args.FromDistrictCode}
	// toQuery := &location.GetLocationQuery{DistrictCode: args.ToDistrictCode}
	// if err := p.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
	// 	return nil, err
	// }
	// fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	// toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province
	//
	// cmd := &CalcShippingFeeCommand{
	// 	ArbitraryID:      args.AccountID,
	// 	FromDistrictCode: args.FromDistrictCode,
	// 	ToDistrictCode:   args.ToDistrictCode,
	// 	Request: &ahamoveclient.CalcShippingFeeRequest{
	// 		Weight:          args.ChargeableWeight,
	// 		Value:           args.GetInsuranceAmount(),
	// 		PickingProvince: fromProvince.Name,
	// 		PickingDistrict: fromDistrict.Name,
	// 		Province:        toProvince.Name,
	// 		District:        toDistrict.Name,
	// 	},
	// }
	// if err := p.CalcShippingFee(ctx, cmd); err != nil {
	// 	return nil, err
	// }
	// providerServices := cmd.Result
	//
	// // get ETOP services
	// etopServices := etop_shipping_price.GetEtopShippingServices(model.TypeGHTK, fromProvince, toProvince, toDistrict, args.ChargeableWeight)
	// etopServices, _ = etop_shipping_price.FillInfoEtopServices(providerServices, etopServices)
	//
	// allServices := append(providerServices, etopServices...)
	// return allServices, nil
	return nil, cm.ErrTODO
}
