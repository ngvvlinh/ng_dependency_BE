package connection

import (
	"context"

	"o.o/api/main/connectioning"
	api "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipping/carrier"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ConnectionService struct {
	session.Session

	ShipmentManager *shippingcarrier.ShipmentManager
	ConnectionQuery connectioning.QueryBus
	ConnectionAggr  connectioning.CommandBus
}

func (s *ConnectionService) Clone() api.ConnectionService { res := *s; return &res }

func (s *ConnectionService) GetConnections(ctx context.Context, q *pbcm.Empty) (*types.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType: connection_type.Shipping,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
	}
	return result, nil
}

func (s *ConnectionService) GetAvailableConnections(ctx context.Context, q *pbcm.Empty) (*types.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:   connection_type.Shipping,
		ConnectionMethod: connection_type.ConnectionMethodDirect,
		Status:           status3.WrapStatus(status3.P),
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
	}
	return result, nil
}

func (s *ConnectionService) GetShopConnections(ctx context.Context, q *pbcm.Empty) (*types.GetShopConnectionsResponse, error) {
	query := &connectioning.ListShopConnectionsByShopIDQuery{
		ShopID: s.SS.Shop().ID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetShopConnectionsResponse{
		ShopConnections: convertpb.PbShopConnections(query.Result),
	}
	return result, nil
}

func (s *ConnectionService) RegisterShopConnection(ctx context.Context, q *types.RegisterShopConnectionRequest) (*types.ShopConnection, error) {
	cmd := &carrier.ShopConnectionSignUpArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       s.SS.Shop().ID,
		Name:         q.Name,
		Identifier:   q.Email,
		Password:     q.Password,
		Phone:        q.Phone,
		Province:     q.Province,
		District:     q.District,
		Address:      q.Address,
	}
	shopConnection, err := s.ShipmentManager.ShopConnectionSignUp(ctx, cmd)
	if err != nil {
		return nil, err
	}
	result := convertpb.PbShopConnection(shopConnection)
	return result, nil
}

func (s *ConnectionService) LoginShopConnection(ctx context.Context, q *types.LoginShopConnectionRequest) (*types.LoginShopConnectionResponse, error) {
	identifier := cm.Coalesce(q.Email, q.Identifier)
	cmd := &carrier.ShopConnectionSignInArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       s.SS.Shop().ID,
		Identifier:   identifier,
		Password:     q.Password,
	}
	result, err := s.ShipmentManager.ShopConnectionSignIn(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ConnectionService) LoginShopConnectionWithOTP(ctx context.Context, q *types.LoginShopConnectionWithOTPRequest) (*types.LoginShopConnectionWithOTPResponse, error) {
	cmd := &carrier.ShopConnectionSignInWithOTPArgs{
		ConnectionID: q.ConnectionID,
		ShopID:       s.SS.Shop().ID,
		Identifier:   q.Identifier,
		OTP:          q.OTP,
	}
	_, err := s.ShipmentManager.ShopConnectionSignInWithOTP(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &types.LoginShopConnectionWithOTPResponse{
		Code: "OK",
	}, nil
}

func (s *ConnectionService) DeleteShopConnection(ctx context.Context, q *types.DeleteShopConnectionRequest) (*pbcm.DeletedResponse, error) {
	cmd := &connectioning.DeleteShopConnectionCommand{
		ConnectionID: q.ConnectionID,
		ShopID:       s.SS.Shop().ID,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.DeletedResponse{
		Deleted: cmd.Result,
	}
	return result, nil
}

func (s *ConnectionService) UpdateShopConnection(ctx context.Context, r *types.UpdateShopConnectionRequest) (*pbcm.UpdatedResponse, error) {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
	}

	identifier := cm.Coalesce(r.ExternalData.Email, r.ExternalData.Identifier)
	cmd := &connectioning.UpdateShopConnectionTokenCommand{
		ShopID:       s.SS.Shop().ID,
		ConnectionID: r.ConnectionID,
		Token:        r.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID:     r.ExternalData.UserID,
			Identifier: identifier,
			ShopID:     r.ExternalData.ShopID,
		},
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}
