package partnercarrier

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/top/external/partnercarrier"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type ShipmentConnectionService struct {
	session.Session

	ConnectionQuery connectioning.QueryBus
	ConnectionAggr  connectioning.CommandBus
}

func (s *ShipmentConnectionService) Clone() partnercarrier.ShipmentConnectionService {
	res := *s
	return &res
}

func (s *ShipmentConnectionService) GetConnections(ctx context.Context, r *pbcm.Empty) (*partnercarrier.GetConnectionsResponse, error) {
	query := &connectioning.ListConnectionsQuery{
		PartnerID:          s.SS.Partner().ID,
		ConnectionType:     connection_type.Shipping,
		ConnectionMethod:   connection_type.ConnectionMethodDirect,
		ConnectionProvider: connection_type.ConnectionProviderPartner,
		Result:             nil,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := &partnercarrier.GetConnectionsResponse{
		Connections: convertpb.PbShipmentConnections(query.Result),
	}
	return res, nil
}

func (s *ShipmentConnectionService) CreateConnection(ctx context.Context, r *partnercarrier.CreateConnectionRequest) (*partnercarrier.ShipmentConnection, error) {
	query := &connectioning.ListConnectionsQuery{
		PartnerID:        s.SS.Partner().ID,
		ConnectionType:   connection_type.Shipping,
		ConnectionMethod: connection_type.ConnectionMethodDirect,
	}
	if err := s.ConnectionQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	if len(query.Result) > 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Bạn đã tạo kết nối vận chuyển. Không thể tạo thêm kết nối mới.")
	}

	cmd := &connectioning.CreateConnectionCommand{
		Name:      r.Name,
		PartnerID: s.SS.Partner().ID,
		Driver:    "",
		DriverConfig: &connectioning.ConnectionDriverConfig{
			TrackingURL:            r.TrackingURL,
			CreateFulfillmentURL:   r.CreateFulfillmentURL,
			GetFulfillmentURL:      r.GetFulfillmentURL,
			GetShippingServicesURL: r.GetShippingServicesURL,
			CancelFulfillmentURL:   r.CancelFulfillmentURL,
			SignInURL:              r.SignInURL,
			SignUpURL:              r.SignUpURL,
		},
		ConnectionType:     connection_type.Shipping,
		ConnectionSubtype:  connection_type.ConnectionSubtypeShipment,
		ConnectionMethod:   connection_type.ConnectionMethodDirect,
		ConnectionProvider: connection_type.ConnectionProviderPartner,
		ImageURL:           r.ImageURL,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShipmentConnection(cmd.Result), nil
}

func (s *ShipmentConnectionService) UpdateConnection(ctx context.Context, r *partnercarrier.UpdateConnectionRequest) (*partnercarrier.ShipmentConnection, error) {
	if r.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ID không được để trống")
	}
	cmd := &connectioning.UpdateConnectionCommand{
		ID:        r.ID,
		Name:      r.Name,
		ImageURL:  r.ImageURL,
		PartnerID: s.SS.Partner().ID,
		DriverConfig: &connectioning.ConnectionDriverConfig{
			TrackingURL:            r.TrackingURL,
			CreateFulfillmentURL:   r.CreateFulfillmentURL,
			GetFulfillmentURL:      r.GetFulfillmentURL,
			GetShippingServicesURL: r.GetShippingServicesURL,
			CancelFulfillmentURL:   r.CancelFulfillmentURL,
			SignInURL:              r.SignInURL,
			SignUpURL:              r.SignUpURL,
		},
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return convertpb.PbShipmentConnection(cmd.Result), nil
}

func (s *ShipmentConnectionService) DeleteConnection(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	if r.Id == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "ID không được để trống")
	}
	cmd := &connectioning.DeleteConnectionCommand{
		ID:        r.Id,
		PartnerID: s.SS.Partner().ID,
	}
	if err := s.ConnectionAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.DeletedResponse{Deleted: 1}, nil
}
