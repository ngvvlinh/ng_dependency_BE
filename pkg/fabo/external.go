package fabo

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/etc/connection_type"
	"o.o/backend/com/main/shipping/carrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/authorize/session"
	ghnclientv2 "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	ghndriverv2 "o.o/backend/pkg/integration/shipping/ghn/driverv2"
)

type ExtraShipmentService struct {
	session.Session

	ShipmentManager *carrier.ShipmentManager
	ConnectionQS    connectioning.QueryBus
}

func (s *ExtraShipmentService) Clone() fabo.ExtraShipmentService {
	res := *s
	return &res
}

func (s *ExtraShipmentService) CustomerReturnRate(
	ctx context.Context, req *fabo.CustomerReturnRateRequest,
) (*fabo.CustomerReturnRateResponse, error) {
	getConnectionQuery := &connectioning.GetConnectionByIDQuery{
		ID: req.ConnectionID,
	}
	if err := s.ConnectionQS.Dispatch(ctx, getConnectionQuery); err != nil {
		return nil, err
	}

	connection := getConnectionQuery.Result
	if connection.ConnectionMethod != connection_type.ConnectionMethodDirect ||
		connection.ConnectionProvider != connection_type.ConnectionProviderGHN {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "connection không hợp lệ")
	}

	shopID := s.SS.Shop().ID
	shipmentDriver, err := s.ShipmentManager.GetShipmentDriver(ctx, connection.ID, shopID)
	if err != nil {
		return nil, err
	}
	ghnDriver := shipmentDriver.(*ghndriverv2.GHNDriver)
	ghnClient := ghnDriver.GetClient()

	etlCustomerRateReq := &ghnclientv2.CustomerReturnRateRequest{
		Phone: req.Phone,
	}
	etlCustomerRateResp, err := ghnClient.CustomerReturnRate(ctx, etlCustomerRateReq)
	if err != nil {
		return nil, err
	}

	return &fabo.CustomerReturnRateResponse{
		Level:     etlCustomerRateResp.Level.String(),
		LevelCode: etlCustomerRateResp.LevelCode.String(),
		Rate:      float64(etlCustomerRateResp.Rate),
	}, nil
}
