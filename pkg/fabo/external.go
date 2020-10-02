package fabo

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/shipping"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type ExtraShipmentService struct {
	session.Session

	ShippingQS   shipping.QueryBus
	ConnectionQS connectioning.QueryBus
}

func (s *ExtraShipmentService) Clone() fabo.ExtraShipmentService {
	res := *s
	return &res
}

func (s *ExtraShipmentService) CustomerReturnRate(
	ctx context.Context, req *fabo.CustomerReturnRateRequest,
) (*fabo.CustomerReturnRateResponse, error) {
	shopID := s.SS.Shop().ID

	listConnectionsQuery := &connectioning.ListConnectionsQuery{
		Status:             status3.WrapStatus(status3.P),
		ConnectionType:     connection_type.Shipping,
		ConnectionMethod:   connection_type.ConnectionMethodDirect,
		ConnectionProvider: connection_type.ConnectionProviderGHN,
	}
	if err := s.ConnectionQS.Dispatch(ctx, listConnectionsQuery); err != nil {
		return nil, err
	}
	connections := listConnectionsQuery.Result

	if len(connections) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "không tìm thấy connection hợp lệ")
	}
	var connectionIDs []dot.ID
	for _, conn := range connections {
		connectionIDs = append(connectionIDs, conn.ID)
	}

	listCustomerReturnRatesQuery := &shipping.ListCustomerReturnRatesQuery{
		ConnectionIDs: connectionIDs,
		ShopID:        shopID,
		Phone:         req.Phone,
	}
	if err := s.ShippingQS.Dispatch(ctx, listCustomerReturnRatesQuery); err != nil {
		return nil, err
	}
	customerReturnRateExtendeds := listCustomerReturnRatesQuery.Result

	var customerReturnRateExtendedsResp []*fabo.CustomerReturnRateExtended
	for _, customerReturnRateExtended := range customerReturnRateExtendeds {
		customerReturnRate := customerReturnRateExtended.CustomerReturnRate
		connection := customerReturnRateExtended.Connection

		customerReturnRateExtendedResp := &fabo.CustomerReturnRateExtended{
			ConnectionID:     connection.ID,
			ConnectionName:   connection.Name,
			ConnectionMethod: connection.ConnectionMethod,
			CustomerReturnRate: &fabo.CustomerReturnRate{
				Level:     customerReturnRate.Level,
				LevelCode: customerReturnRate.LevelCode,
				Rate:      customerReturnRate.Rate,
			},
		}
		customerReturnRateExtendedsResp = append(customerReturnRateExtendedsResp, customerReturnRateExtendedResp)
	}

	return &fabo.CustomerReturnRateResponse{
		CustomerReturnRateExtendeds: customerReturnRateExtendedsResp,
	}, nil
}
