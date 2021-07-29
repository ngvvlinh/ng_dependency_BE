package aggregate

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/top/types/etc/status3"
	providertypes "o.o/backend/com/etelecom/provider/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

func (a *EtelecomAggregate) CreateHotline(ctx context.Context, args *etelecom.CreateHotlineArgs) (*etelecom.Hotline, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}
	_, err := a.hotlineStore(ctx).Hotline(args.Hotline).GetHotline()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}
	if err == nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Hotline was existed")
	}

	var hotline etelecom.Hotline
	if err = scheme.Convert(args, &hotline); err != nil {
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

func (a *EtelecomAggregate) DeleteHotline(ctx context.Context, id dot.ID) error {
	hotline, err := a.hotlineStore(ctx).ID(id).GetHotline()
	if err != nil {
		return err
	}
	if hotline.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Hotline was activated. Can not delete")
	}
	_, err = a.hotlineStore(ctx).ID(id).SoftDelete()
	return err
}

func (a *EtelecomAggregate) RemoveHotlineOutOfTenant(ctx context.Context, args *etelecom.RemoveHotlineOutOfTenantArgs) error {
	if err := args.Validate(); err != nil {
		return err
	}
	hotline, err := a.hotlineStore(ctx).ID(args.HotlineID).OptionalOwnerID(args.OwnerID).GetHotline()
	if err != nil {
		return err
	}
	if hotline.Status != status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Hotline was not activated")
	}

	tenant, err := a.tenantStore(ctx).OwnerID(args.OwnerID).ID(hotline.TenantID).GetTenant()
	if err != nil && cm.ErrorCode(err) == cm.NotFound {
		return cm.Errorf(cm.NotFound, nil, "Hotline does not belongs to any Tenant")
	}
	if err != nil {
		return err
	}
	if !tenant.Status.Valid || tenant.Status.Enum != status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Tenant %v was not activated", tenant.Name)
	}

	return a.txDBEtelecom.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// unactive tenant + hotline
		if err = a.tenantStore(ctx).ID(tenant.ID).UpdateTenantStatus(status3.Z); err != nil {
			return err
		}
		if err = a.hotlineStore(ctx).ID(hotline.ID).UpdateHotlineStatus(status3.Z); err != nil {
			return err
		}
		event := &etelecom.RemovedHotlineOutOfTenantEvent{
			OwnerID:       args.OwnerID,
			HotlineNumber: hotline.Hotline,
			TenantID:      hotline.TenantID,
		}
		return a.eventBus.Publish(ctx, event)
	})
}

func (a *EtelecomAggregate) ActiveHotlineForTenant(ctx context.Context, args *etelecom.ActiveHotlineForTenantArgs) error {
	tenant, err := a.tenantStore(ctx).OptionalOwnerID(args.OwnerID).ID(args.TenantID).GetTenant()
	if err != nil {
		return err
	}
	if !tenant.Status.Valid || tenant.Status.Enum != status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Tenant %v was not activated", tenant.Name)
	}

	hotline, err := a.hotlineStore(ctx).OptionalOwnerID(args.OwnerID).ID(args.HotlineID).GetHotline()
	if err != nil {
		return err
	}
	if hotline.Status == status3.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Hotline was activated.")
	}
	if hotline.Hotline == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "Missing hotline number")
	}

	update := &providertypes.AddHotlineToTenantInTrunkProviderRequest{
		TrunkProviderID: a.telecomManager.AdminPortsip.TrunkProviderDefaultID,
		TenantID:        tenant.ExternalData.ID,
		Hotline:         hotline.Hotline,
	}
	if err := a.updateTrunkProvider(ctx, update); err != nil {
		return cm.Errorf(cm.ErrorCode(err), nil, "Error when update trunk provider: %v", err.Error())
	}

	updateHotlineArgs := &etelecom.UpdateHotlineInfoArgs{
		ID:               hotline.ID,
		Status:           status3.P.Wrap(),
		TenantID:         tenant.ID,
		ConnectionID:     tenant.ConnectionID,
		ConnectionMethod: tenant.ConnectionMethod,
		OwnerID:          tenant.OwnerID,
	}
	return a.UpdateHotlineInfo(ctx, updateHotlineArgs)
}

func (a *EtelecomAggregate) updateTrunkProvider(ctx context.Context, args *providertypes.AddHotlineToTenantInTrunkProviderRequest) error {
	portsipAdminDriver, err := a.telecomManager.GetAdministratorPortsipDriver(ctx)
	if err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Please check Portsip admin account")
	}

	return portsipAdminDriver.AddHotlineToTenantInTrunkProvider(ctx, args)
}
