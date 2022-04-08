package admin

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/shopping/setting"
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

	SettingAggr     setting.CommandBus
	ConnectionAggr  connectioning.CommandBus
	ConnectionQuery connectioning.QueryBus
	SettingQuery    setting.QueryBus
}

func (s *ConnectionService) Clone() admin.ConnectionService {
	res := *s
	return &res
}

func (s *ConnectionService) GetConnections(ctx context.Context, r *types.GetConnectionsRequest) (*types.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:    connection_type.Shipping,
		ConnectionMethod:  r.ConnectionMethod,
		ConnectionSubtype: r.ConnectionSubtype,
	}
	if r.ConnectionSubtype == connection_type.ConnectionSubtypeShipnow {
		// TopShip only support shipnow direct
		// shipnow builtin use for 3rd party
		query.ConnectionMethod = connection_type.ConnectionMethodDirect
	}

	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetConnectionsResponse{
		Connections: convertpb.PbConnections(query.Result),
	}
	return result, nil
}

func (s *ConnectionService) GetConnectDirectShipmentShopSetting(ctx context.Context, r *types.GetConnectDirectShipmentShopSettingRequest) (*types.GetConnectDirectShipmentSettingResponse, error) {
	query := &setting.GetShopSettingDirectShipmentQuery{
		ShopID: r.ShopID,
	}
	if err := s.SettingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	result := &types.GetConnectDirectShipmentSettingResponse{
		ShopID:                     query.ShopID,
		AllowConnectDirectShipment: query.Result.AllowConnectDirectShipment,
	}
	return result, nil
}

func (s *ConnectionService) UpdateConnectDirectShipmentShopSetting(ctx context.Context, r *types.UpdateDirectShipmentSettingRequest) (*types.UpdateDirectShipmentSettingResponse, error) {
	cmd := &setting.UpdateShopSettingDirectShipmentCommand{
		ShopID:                     r.ShopID,
		AllowConnectDirectShipment: r.AllowConnectDirectShipment,
	}

	if err := s.SettingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := &types.UpdateDirectShipmentSettingResponse{
		ShopID:                     cmd.ShopID,
		AllowConnectDirectShipment: cmd.AllowConnectDirectShipment,
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

	identifier := cm.Coalesce(r.ExternalData.Identifier, r.ExternalData.Email)
	cmd := &connectioning.CreateBuiltinConnectionCommand{
		ID:    r.ConnectionID,
		Name:  r.Name,
		Token: r.Token,
		ExternalData: &connectioning.ShopConnectionExternalData{
			UserID:     r.ExternalData.UserID,
			Identifier: identifier,
			ShopID:     r.ExternalData.ShopID,
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

func (s *ConnectionService) UpdateShopConnection(ctx context.Context, r *types.UpdateShopConnectionRequest) (*pbcm.UpdatedResponse, error) {
	query := &connectioning.GetConnectionByIDQuery{
		ID: r.ConnectionID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	conn := query.Result

	switch conn.ConnectionType {
	case connection_type.Shipping:
		if r.ExternalData == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "External data không được để trống")
		}
	case connection_type.Telecom:
		if r.TelecomData == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Telecom data không được để trống")
		}
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}

	var externalData *connectioning.ShopConnectionExternalData
	var telecomData *connectioning.ShopConnectionTelecomData
	if r.ExternalData != nil {
		identifier := cm.Coalesce(r.ExternalData.Email, r.ExternalData.Identifier)
		externalData = &connectioning.ShopConnectionExternalData{
			UserID:     r.ExternalData.UserID,
			Identifier: identifier,
			ShopID:     r.ExternalData.ShopID,
		}
	}
	if r.TelecomData != nil {
		telecomData = &connectioning.ShopConnectionTelecomData{
			Username:     r.TelecomData.Username,
			Password:     r.TelecomData.Password,
			TenantHost:   r.TelecomData.TenantHost,
			TenantToken:  r.TelecomData.TenantToken,
			TenantDomain: r.TelecomData.TenantDomain,
		}
	}

	cmd := &connectioning.UpdateShopConnectionCommand{
		ConnectionID:   r.ConnectionID,
		ShopID:         r.ShopID,
		OwnerID:        r.OwnerID,
		Token:          r.Token,
		TokenExpiresAt: r.TokenExpiresAt,
		ExternalData:   externalData,
		TelecomData:    telecomData,
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
