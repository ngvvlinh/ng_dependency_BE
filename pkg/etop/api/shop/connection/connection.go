package connection

import (
	"context"

	"o.o/api/main/accountshipnow"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	api "o.o/api/top/int/shop"
	"o.o/api/top/int/types"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/shipping/carrier"
	shippingcarrier "o.o/backend/com/main/shipping/carrier"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ConnectionService struct {
	session.Session

	ShipmentManager    *shippingcarrier.ShipmentManager
	ConnectionQuery    connectioning.QueryBus
	ConnectionAggr     connectioning.CommandBus
	IdentityQuery      identity.QueryBus
	AccountshipnowAggr accountshipnow.CommandBus
}

func (s *ConnectionService) Clone() api.ConnectionService { res := *s; return &res }

func (s *ConnectionService) GetConnections(ctx context.Context, q *types.GetConnectionsRequest) (*types.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:    connection_type.Shipping,
		ConnectionSubtype: q.ConnectionSubtype,
		ConnectionMethod:  q.ConnectionMethod,
	}
	if q.ConnectionSubtype == connection_type.ConnectionSubtypeShipnow {
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

// GetAvailableConnections
//
// Lấy tất cả connection direct (shop có thể login)
func (s *ConnectionService) GetAvailableConnections(ctx context.Context, q *types.GetAvailableConnectionsRequest) (*types.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		ConnectionType:    connection_type.Shipping,
		ConnectionMethod:  connection_type.ConnectionMethodDirect,
		ConnectionSubtype: q.ConnectionSubtype,
		Status:            status3.WrapStatus(status3.P),
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
	query := &connectioning.ListShopConnectionsQuery{
		ShopID:  s.SS.Shop().ID,
		OwnerID: s.SS.Shop().OwnerID,
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

func (s *ConnectionService) LoginShopConnection(ctx context.Context, q *types.LoginShopConnectionRequest) (res *types.LoginShopConnectionResponse, err error) {
	identifier := cm.Coalesce(q.Identifier, q.Email)
	queryConn := &connectioning.GetConnectionByIDQuery{
		ID: q.ConnectionID,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}

	switch queryConn.Result.ConnectionSubtype {
	case connection_type.ConnectionSubtypeShipment:
		cmd := &carrier.ShopConnectionSignInArgs{
			ConnectionID: q.ConnectionID,
			ShopID:       s.SS.Shop().ID,
			Identifier:   identifier,
			Password:     q.Password,
		}
		return s.ShipmentManager.ShopConnectionSignIn(ctx, cmd)
	case connection_type.ConnectionSubtypeShipnow:
		phoneNorm, ok := validate.NormalizePhone(identifier)
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Số điện thoại không hợp lệ")
		}

		user := s.SS.User()
		// TODO: handle multiple shipnow connection
		// just handle ahamove now
		cmd := &accountshipnow.CreateExternalAccountAhamoveCommand{
			ShopID:       s.SS.Shop().ID,
			OwnerID:      user.ID,
			Phone:        phoneNorm.String(),
			Name:         user.FullName,
			ConnectionID: q.ConnectionID,
		}
		if err := s.AccountshipnowAggr.Dispatch(ctx, cmd); err != nil {
			return nil, err
		}
		return &types.LoginShopConnectionResponse{
			Code: "OK",
		}, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ConnectionID không hợp lệ")
	}
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
		OwnerID:      s.SS.Shop().OwnerID,
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
	cmd := &connectioning.CreateOrUpdateShopConnectionCommand{
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
