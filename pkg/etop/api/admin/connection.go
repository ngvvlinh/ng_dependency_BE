package admin

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/top/int/admin"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ConnectionService struct {
	session.Session

	ConnectionAggr  connectioning.CommandBus
	ConnectionQuery connectioning.QueryBus
}

func (s *ConnectionService) Clone() admin.ConnectionService {
	res := *s
	return &res
}

func (s *ConnectionService) GetConnections(ctx context.Context, r *types.GetConnectionsRequest) (*types.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:   connection_type.Shipping,
		ConnectionMethod: r.ConnectionMethod,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
	}
	return result, nil
}

func (s *ConnectionService) ConfirmConnection(ctx context.Context, r *pbcm.IDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &connectioning.ConfirmConnectionCommand{
		ID: r.Id,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *ConnectionService) DisableConnection(ctx context.Context, r *pbcm.IDRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &connectioning.DisableConnectionCommand{
		ID: r.Id,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: cmd.Result}
	return result, nil
}

func (s *ConnectionService) CreateBuiltinConnection(ctx context.Context, r *types.CreateBuiltinConnectionRequest) (*types.Connection, error) {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
	}
	cmd := &connectioning.CreateBuiltinConnectionCommand{
		ID:    r.ConnectionID,
		Name:  r.Name,
		Token: r.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID: r.ExternalData.UserID,
			Email:  r.ExternalData.Email,
		},
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbConnection(cmd.Result)
	return result, nil
}

func (s *ConnectionService) GetBuiltinShopConnections(ctx context.Context, r *pbcm.Empty) (*types.GetShopConnectionsResponse, error) {
	query := &connectioning.ListGlobalShopConnectionsQuery{}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetShopConnectionsResponse{
		ShopConnections: convertpb.PbShopConnections(query.Result),
	}
	return result, nil
}

func (s *ConnectionService) UpdateBuiltinShopConnection(ctx context.Context, r *types.UpdateShopConnectionRequest) (*pbcm.UpdatedResponse, error) {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
	}
	query := &connectioning.GetConnectionByIDQuery{
		ID: r.ConnectionID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	conn := query.Result
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Connection không hợp lệ")
	}

	cmd := &connectioning.UpdateShopConnectionTokenCommand{
		ConnectionID: r.ConnectionID,
		Token:        r.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID: r.ExternalData.UserID,
			Email:  r.ExternalData.Email,
		},
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &pbcm.UpdatedResponse{Updated: 1}
	return result, nil
}

func (s *ConnectionService) GetConnectionServices(ctx context.Context, r *pbcm.IDRequest) (*types.GetConnectionServicesResponse, error) {
	query := &connectioning.ListConnectionServicesByIDQuery{
		ID: r.Id,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetConnectionServicesResponse{
		ConnectionService: convertpb.PbConnectionServices(query.Result),
	}
	return result, nil
}
