package admin

import (
	"context"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/int/types"
	pbcm "etop.vn/api/top/types/common"
	"etop.vn/api/top/types/etc/connection_type"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

func (s *ConnectionService) GetConnections(ctx context.Context, r *GetConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType: connection_type.Shipping,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
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
	if err := connectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *ConnectionService) DisableConnection(ctx context.Context, r *DisableConnectionEndpoint) error {
	cmd := &connectioning.DisableConnectionCommand{
		ID: r.Id,
	}
	if err := connectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{Updated: cmd.Result}
	return nil
}

func (s *ConnectionService) CreateTopshipConnection(ctx context.Context, r *CreateTopshipConnectionEndpoint) error {
	if r.ExternalData == nil || r.ExternalData.UserID == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "UserID không được để trống")
	}
	cmd := &connectioning.CreateTopshipConnectionCommand{
		ID:    r.ConnectionID,
		Token: r.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID: r.ExternalData.UserID,
			Email:  r.ExternalData.Email,
		},
	}
	if err := connectionAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbConnection(cmd.Result)
	return nil
}

func (s *ConnectionService) GetConnectionServices(ctx context.Context, r *GetConnectionServicesEndpoint) error {
	query := &connectioning.ListConnectionServicesByIDQuery{
		ID: r.Id,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &types.GetConnectionServicesResponse{
		ConnectionService: convertpb.PbConnectionServices(query.Result),
	}
	return nil
}
