package shop

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/int/shop"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/backend/com/main/shipping/carrier"
)

func (s *ConnectionService) GetConnections(ctx context.Context, q *GetConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType: connection_type.Shipping,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return nil
	}
	q.Result = &shop.GetConnectionsResponse{Connections: PbConnections(query.Result)}
	return nil
}

func (s *ConnectionService) GetAvailableConnections(ctx context.Context, q *GetAvailableConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:   connection_type.Shipping,
		ConnectionMethod: connection_type.ConnectionMethodDirect,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetConnectionsResponse{
		Connections: PbConnections(query.Result),
	}
	return nil
}

func (s *ConnectionService) GetShopConnections(ctx context.Context, q *GetShopConnectionsEndpoint) error {
	query := &connectioning.ListShopConnectionsByShopIDQuery{
		ShopID: q.Context.Shop.ID,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &shop.GetShopConnectionsResponse{
		ShopConnections: PbShopConnections(query.Result),
	}
	return nil
}

func (s *ConnectionService) RegisterShopConnection(ctx context.Context, q *RegisterShopConnectionEndpoint) error {
	cmd := &carrier.ShopConnectionSignUpArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       q.Context.Shop.ID,
		Name:         q.Name,
		Email:        q.Email,
		Password:     q.Password,
		Phone:        q.Phone,
		Province:     q.Province,
		District:     q.District,
		Address:      q.Address,
	}
	shopConnection, err := shipmentManager.ShopConnectionSignUp(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbShopConnection(shopConnection)
	return nil
}

func (s *ConnectionService) LoginShopConnection(ctx context.Context, q *LoginShopConnectionEndpoint) error {
	cmd := &carrier.ShopConnectionSignInArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       q.Context.Shop.ID,
		Email:        q.Email,
		Password:     q.Password,
	}
	shopConnection, err := shipmentManager.ShopConnectionSignIn(ctx, cmd)
	if err != nil {
		return err
	}
	q.Result = PbShopConnection(shopConnection)
	return nil
}
