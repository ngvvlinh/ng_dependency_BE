package aggregate

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
)

func (a *EtelecomAggregate) CreateHotline(ctx context.Context, args *etelecom.CreateHotlineArgs) (*etelecom.Hotline, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	queryConn := &connectioning.GetConnectionByIDQuery{
		ID: args.ConnectionID,
	}
	if err := a.connectionQuery.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}
	conn := queryConn.Result
	if conn.ConnectionMethod == connection_type.ConnectionMethodDirect {
		if args.OwnerID == 0 {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing owner ID")
		}
	}

	var hotline etelecom.Hotline
	if err := scheme.Convert(args, &hotline); err != nil {
		return nil, err
	}
	hotline.ID = cm.NewID()
	hotline.ConnectionMethod = conn.ConnectionMethod
	return a.hotlineStore(ctx).CreateHotline(&hotline)
}

func (a *EtelecomAggregate) UpdateHotlineInfo(ctx context.Context, args *etelecom.UpdateHotlineInfoArgs) error {
	update := &etelecom.Hotline{
		IsFreeCharge: args.IsFreeCharge,
		Name:         args.Name,
		Description:  args.Description,
		Status:       args.Status.Apply(status3.Z),
	}
	return a.hotlineStore(ctx).ID(args.ID).UpdateHotline(update)
}
