package partner

import (
	"context"

	"o.o/api/main/connectioning"
	extpartner "o.o/api/top/external/partner"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/connection_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/convertpb"
)

func (s *ShipmentConnectionService) GetConnections(ctx context.Context, r *GetConnectionsEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		PartnerID:          r.Context.Partner.ID,
		ConnectionType:     connection_type.Shipping,
		ConnectionMethod:   connection_type.ConnectionMethodDirect,
		ConnectionProvider: connection_type.ConnectionProviderPartner,
		Result:             nil,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	r.Result = &extpartner.GetConnectionsResponse{
		Connections: convertpb.PbShipmentConnections(query.Result),
	}

	return nil
}

func (s *ShipmentConnectionService) CreateConnection(ctx context.Context, r *CreateConnectionEndpoint) error {
	query := &connectioning.ListConnectionsQuery{
		PartnerID:        r.Context.Partner.ID,
		ConnectionMethod: connection_type.ConnectionMethodDirect,
	}
	if err := connectionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	if len(query.Result) > 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Bạn đã tạo kết nối vận chuyển. Không thể tạo thêm kết nối mới.")
	}

	cmd := &connectioning.CreateConnectionCommand{
		Name:      r.Name,
		PartnerID: r.Context.Partner.ID,
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
	}
	if err := connectionAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentConnection(cmd.Result)
	return nil
}

func (s *ShipmentConnectionService) UpdateConnection(ctx context.Context, r *UpdateConnectionEndpoint) error {
	if r.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "ID không được để trống")
	}
	cmd := &connectioning.UpdateConnectionCommand{
		ID:        r.ID,
		Name:      r.Name,
		ImageURL:  r.ImageURL,
		PartnerID: r.Context.Partner.ID,
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
	if err := connectionAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = convertpb.PbShipmentConnection(cmd.Result)
	return nil
}

func (s *ShipmentConnectionService) DeleteConnection(ctx context.Context, r *DeleteConnectionEndpoint) error {
	if r.Id == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "ID không được để trống")
	}
	cmd := &connectioning.DeleteConnectionCommand{
		ID:        r.Id,
		PartnerID: r.Context.Partner.ID,
	}
	if err := connectionAggregate.Dispatch(ctx, cmd); err != nil {
		return err
	}
	r.Result = &pbcm.DeletedResponse{Deleted: 1}
	return nil
}
