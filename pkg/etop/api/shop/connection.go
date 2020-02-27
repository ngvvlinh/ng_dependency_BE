package shop

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/int/types"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/shipping/carrier"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

func (s *ConnectionService) GetConnections(ctx context.Context, q *GetConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType: connection_type.Shipping,
		Status:         status3.WrapStatus(status3.P),
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return nil
	}
	q.Result = &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
	}
	return nil
}

func (s *ConnectionService) GetAvailableConnections(ctx context.Context, q *GetAvailableConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:   connection_type.Shipping,
		ConnectionMethod: connection_type.ConnectionMethodDirect,
		Status:           status3.WrapStatus(status3.P),
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	q.Result = &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
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
	q.Result = &types.GetShopConnectionsResponse{
		ShopConnections: convertpb.PbShopConnections(query.Result),
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
	q.Result = convertpb.PbShopConnection(shopConnection)
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
	q.Result = convertpb.PbShopConnection(shopConnection)
	return nil
}

func (s *ConnectionService) DeleteShopConnection(ctx context.Context, q *DeleteShopConnectionEndpoint) error {
	cmd := &connectioning.DeleteShopConnectionCommand{
		ConnectionID: q.ConnectionID,
		ShopID:       q.Context.Shop.ID,
	}
	if err := connectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &pbcm.DeletedResponse{
		Deleted: cmd.Result,
	}
	return nil
}

func (s *ConnectionService) UpdateShopConnection(ctx context.Context, r *UpdateShopConnectionEndpoint) error {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
	}
	cmd := &connectioning.UpdateShopConnectionTokenCommand{
		ShopID:       r.Context.Shop.ID,
		ConnectionID: r.ConnectionID,
		Token:        r.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID: r.ExternalData.UserID,
			Email:  r.ExternalData.Email,
		},
	}
	if err := connectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}
