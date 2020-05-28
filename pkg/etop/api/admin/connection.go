package admin

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
)

type ConnectionService struct {
	ConnectionAggr  connectioning.CommandBus
	ConnectionQuery connectioning.QueryBus
}

func (s *ConnectionService) Clone() *ConnectionService {
	res := *s
	return &res
}

func (s *ConnectionService) GetConnections(ctx context.Context, r *GetConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:   connection_type.Shipping,
		ConnectionMethod: r.ConnectionMethod,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
	}
	return nil
}

func (s *ConnectionService) ConfirmConnection(ctx context.Context, r *ConfirmConnectionEndpoint) error {
	cmd := &connectioning.ConfirmConnectionCommand{
		ID: r.Id,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *ConnectionService) DisableConnection(ctx context.Context, r *DisableConnectionEndpoint) error {
	cmd := &connectioning.DisableConnectionCommand{
		ID: r.Id,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *ConnectionService) CreateBuiltinConnection(ctx context.Context, r *CreateBuiltinConnectionEndpoint) error {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
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
		return err
	}
	r.Result = convertpb.PbConnection(cmd.Result)
	return nil
}

func (s *ConnectionService) GetBuiltinShopConnections(ctx context.Context, r *GetBuiltinShopConnectionsEndpoint) error {
	query := &connectioning.ListGlobalShopConnectionsQuery{}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetShopConnectionsResponse{
		ShopConnections: convertpb.PbShopConnections(query.Result),
	}
	return nil
}

func (s *ConnectionService) UpdateBuiltinShopConnection(ctx context.Context, r *UpdateBuiltinShopConnectionEndpoint) error {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
	}
	query := &connectioning.GetConnectionByIDQuery{
		ID: r.ConnectionID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	conn := query.Result
	if conn.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		return cm.Errorf(cm.FailedPrecondition, nil, "Connection không hợp lệ")
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
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: 1}
	return nil
}

func (s *ConnectionService) GetConnectionServices(ctx context.Context, r *GetConnectionServicesEndpoint) error {
	query := &connectioning.ListConnectionServicesByIDQuery{
		ID: r.Id,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetConnectionServicesResponse{
		ConnectionService: convertpb.PbConnectionServices(query.Result),
	}
	return nil
}
