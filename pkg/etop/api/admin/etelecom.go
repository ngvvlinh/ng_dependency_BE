package admin

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/etelecom/usersetting"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/top/int/admin"
	etelecomtypes "o.o/api/top/int/etelecom/types"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	shopetelecom "o.o/backend/pkg/etop/api/shop/etelecom"
	"o.o/backend/pkg/etop/authorize/session"
)

type EtelecomService struct {
	session.Session

	EtelecomAggr     etelecom.CommandBus
	EtelecomQuery    etelecom.QueryBus
	UserSettingAggr  usersetting.CommandBus
	UserSettingQuery usersetting.QueryBus
	IdentityQuery    identity.QueryBus
}

func (s *EtelecomService) Clone() admin.EtelecomService {
	res := *s
	return &res
}

func (s *EtelecomService) CreateHotline(ctx context.Context, r *etelecomtypes.CreateHotlineRequest) (*etelecomtypes.Hotline, error) {
	cmd := &etelecom.CreateHotlineCommand{
		OwnerID:      r.OwnerID,
		Name:         r.Name,
		Hotline:      r.Hotline,
		Network:      r.Network,
		Description:  r.Description,
		IsFreeCharge: r.IsFreeCharge,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	res := shopetelecom.Convert_etelecom_Hotline_etelecomtypes_Hotline(cmd.Result, nil)
	return res, nil
}

func (s *EtelecomService) UpdateHotline(ctx context.Context, r *etelecomtypes.UpdateHotlineRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &etelecom.UpdateHotlineInfoCommand{
		ID:           r.ID,
		IsFreeCharge: r.IsFreeCharge,
		Name:         r.Name,
		Description:  r.Description,
		Status:       r.Status,
		Network:      r.Network,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *EtelecomService) DeleteHotline(ctx context.Context, r *pbcm.IDRequest) (*pbcm.DeletedResponse, error) {
	cmd := &etelecom.DeleteHotlineCommand{
		Id: r.Id,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.DeletedResponse{Deleted: 1}, nil
}

func (s *EtelecomService) GetHotlines(ctx context.Context, r *etelecomtypes.GetHotLinesRequest) (*etelecomtypes.GetHotLinesResponse, error) {
	query := &etelecom.ListHotlinesQuery{}
	if r.Filter != nil {
		query.OwnerID = r.Filter.OwnerID
		query.TenantID = r.Filter.TenantID
	}
	if err := s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := shopetelecom.Convert_etelecom_Hotlines_etelecomtypes_Hotlines(query.Result)
	return &etelecomtypes.GetHotLinesResponse{
		Hotlines: res,
	}, nil
}

func (s *EtelecomService) GetUserSettings(ctx context.Context, r *etelecomtypes.GetUserSettingsRequest) (*etelecomtypes.UserSettingsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	query := &usersetting.ListUserSettingsQuery{
		UserIDs: r.UserIDs,
		Paging:  *paging,
	}
	if err = s.UserSettingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := &etelecomtypes.UserSettingsResponse{
		UserSettings: convertpb.Convert_usersetting_UserSettings_api_UserSettings(query.Result.UserSettings),
		Paging:       cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}
	return res, nil
}

func (s *EtelecomService) UpdateUserSetting(ctx context.Context, r *etelecomtypes.UpdateUserSettingRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &usersetting.UpdateUserSettingCommand{
		UserID:              r.UserID,
		ExtensionChargeType: r.ExtensionChargeType,
	}
	if err := s.UserSettingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *EtelecomService) GetTenants(ctx context.Context, r *etelecomtypes.GetTenantsRequest) (*etelecomtypes.GetTenantsResponse, error) {
	paging, err := cmapi.CMCursorPaging(r.Paging)
	if err != nil {
		return nil, err
	}
	query := &etelecom.ListTenantsQuery{
		Paging: *paging,
	}
	if r.Filter != nil {
		query.OwnerID = r.Filter.OwnerID
	}
	if err = s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := shopetelecom.Convert_etelecom_Tenants_etelecomtypes_Tenants(query.Result.Tenants)
	return &etelecomtypes.GetTenantsResponse{
		Tenants: res,
		Paging:  cmapi.PbCursorPageInfo(paging, &query.Result.Paging),
	}, nil
}

func (s *EtelecomService) CreateTenant(ctx context.Context, r *etelecomtypes.AdminCreateTenantRequest) (*etelecomtypes.Tenant, error) {
	cmd := &etelecom.CreateTenantCommand{
		OwnerID:      r.OwnerID,
		AccountID:    r.AccountID,
		ConnectionID: r.ConnectionID,
	}
	if cmd.ConnectionID == 0 {
		cmd.ConnectionID = connectioning.DefaultDirectPortsipConnectionID
	}

	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	res := shopetelecom.Convert_etelecom_Tenant_etelecomtypes_Tenant(cmd.Result, nil)
	return res, nil
}

func (s *EtelecomService) ActivateTenant(ctx context.Context, r *etelecomtypes.ActivateTenantRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &etelecom.ActivateTenantCommand{
		OwnerID:      r.OwnerID,
		TenantID:     r.TenantID,
		HotlineID:    r.HotlineID,
		ConnectionID: r.ConnectionID,
	}
	if cmd.ConnectionID == 0 {
		cmd.ConnectionID = connectioning.DefaultDirectPortsipConnectionID
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *EtelecomService) RemoveHotlineOutOfTenant(ctx context.Context, r *etelecomtypes.RemoveHotlineOutOfTenantRequest) (*pbcm.RemovedResponse, error) {
	cmd := &etelecom.RemoveHotlineOutOfTenantCommand{
		HotlineID: r.HotlineID,
		OwnerID:   r.OwnerID,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.RemovedResponse{Removed: 1}, nil
}

func (s *EtelecomService) AddHotlineToTenant(ctx context.Context, r *etelecomtypes.AddHotlineToTenantRequest) (*pbcm.UpdatedResponse, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	cmd := &etelecom.ActiveHotlineForTenantCommand{
		HotlineID: r.HotlineID,
		TenantID:  r.TenantID,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}
