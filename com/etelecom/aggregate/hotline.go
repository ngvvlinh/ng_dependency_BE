package aggregate

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
)

func (a *EtelecomAggregate) CreateHotline(ctx context.Context, args *etelecom.CreateHotlineArgs) (*etelecom.Hotline, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	var hotline etelecom.Hotline
	if err := scheme.Convert(args, &hotline); err != nil {
		return nil, err
	}
	hotline.ID = cm.NewID()
	return a.hotlineStore(ctx).CreateHotline(&hotline)
}

func (a *EtelecomAggregate) UpdateHotlineInfo(ctx context.Context, args *etelecom.UpdateHotlineInfoArgs) error {
	update := &etelecom.Hotline{
		IsFreeCharge:     args.IsFreeCharge,
		Name:             args.Name,
		Description:      args.Description,
		Status:           args.Status.Apply(status3.Z),
		TenantID:         args.TenantID,
		ConnectionID:     args.ConnectionID,
		ConnectionMethod: args.ConnectionMethod,
		OwnerID:          args.OwnerID,
		Network:          args.Network,
	}
	return a.hotlineStore(ctx).ID(args.ID).UpdateHotline(update)
}
