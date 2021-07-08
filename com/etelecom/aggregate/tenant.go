package aggregate

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

const MaxTenantRetry = 10

func (a *EtelecomAggregate) CreateTenant(ctx context.Context, args *etelecom.CreateTenantArgs) (*etelecom.Tenant, error) {
	if args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing owner ID")
	}
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID")
	}

	queryUser := &identity.GetUserByIDQuery{
		UserID: args.OwnerID,
	}
	if err := a.identityQuery.Dispatch(ctx, queryUser); err != nil {
		return nil, err
	}
	user := queryUser.Result

	queryConn := &connectioning.GetConnectionByIDQuery{
		ID: args.ConnectionID,
	}
	if err := a.connectionQuery.Dispatch(ctx, queryConn); err != nil {
		return nil, err
	}
	conn := queryConn.Result
	if conn.ConnectionProvider != connection_type.ConnectionProviderPortsip && conn.ConnectionMethod != connection_type.ConnectionMethodDirect {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Does not support this connection")
	}

	tenantID := dot.ID(0)
	tenant, err := a.tenantStore(ctx).OwnerID(args.OwnerID).ConnectionID(args.ConnectionID).GetTenant()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		tenantID = tenant.ID
		if tenant.ExternalData != nil && tenant.ExternalData.ID != "" {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Tenant đã được tạo cho người dùng này")
		}
	case cm.NotFound:
		// create new one
		tenantID = cm.NewID()
		tenant = &etelecom.Tenant{
			ID:               tenantID,
			ConnectionID:     conn.ID,
			ConnectionMethod: conn.ConnectionMethod,
			OwnerID:          args.OwnerID,
			Status:           status3.Z.Wrap(),
		}
	default:
		return nil, err
	}

	var shop *identity.Shop
	if args.AccountID != 0 {
		queryShop := &identity.GetShopByIDQuery{
			ID: args.AccountID,
		}
		if err = a.identityQuery.Dispatch(ctx, queryShop); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Account ID không hợp lệ")
		}
		shop = queryShop.Result
		if shop.OwnerID != args.OwnerID {
			return nil, cm.Errorf(cm.InvalidArgument, err, "account_id không thuộc owner_id")
		}
	}

	for index := 0; index < MaxTenantRetry; index++ {
		info := getTenantInfo(user, shop, index)
		tenant.Name = info.Name
		tenant.Domain = info.Domain
		tenant.Password = info.Password
		externalTenantResp, _err := a.telecomManager.CreateTenantPortsip(ctx, tenant)
		if _err != nil {
			if cm.ErrorXCode(_err) == cm.PortsipNameOrDomainIncorrect && index < MaxTenantRetry-1 {
				// Tenant bị trùng tên hoặc domain
				// Tạo tenant với tên mới
				continue
			}
			return nil, _err
		}

		// tạo tenent thành công
		tenant.ExternalData = &etelecom.TenantExternalData{ID: externalTenantResp.ExternalTenantID}
		tenant, err = a.tenantStore(ctx).CreateTenant(tenant)
		if err != nil {
			return nil, err
		}
		break
	}

	return tenant, nil
}

type TenantInfo struct {
	Name     string
	Domain   string
	Password string
}

func getTenantInfo(user *identity.User, shop *identity.Shop, index int) *TenantInfo {
	fullname := user.FullName
	if shop != nil {
		fullname = shop.Name
	}
	fullname = normalizeTenantName(fullname)
	userID := user.ID.String()
	idx := userID[len(userID)-4:]
	name := fullname + "-" + idx
	if index != 0 {
		name += "-" + strconv.Itoa(index)
	}
	name += "-" + cmenv.Env().String()

	domain := name + ".eb2b.vn"
	password := "eB2B@" + gencode.GenerateCode(gencode.Alphabet54, 8)
	return &TenantInfo{
		Name:     name,
		Domain:   domain,
		Password: password,
	}
}

func normalizeTenantName(name string) string {
	nameRegexp := regexp.MustCompile(`[!@#$%^&*()_={}?<>/|~-]`)
	name = nameRegexp.ReplaceAllString(name, "")
	name = validate.NormalizeUnaccent(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}

type UpdateExtTenantInfoArgs struct {
	TenantID     dot.ID
	ExternalData *etelecom.TenantExternalData
}

func (a *EtelecomAggregate) updateExternalTenantInfo(ctx context.Context, args *UpdateExtTenantInfoArgs) error {
	update := &etelecom.Tenant{
		ExternalData: &etelecom.TenantExternalData{
			ID: args.ExternalData.ID,
		},
	}
	return a.tenantStore(ctx).ID(args.TenantID).UpdateTenant(update)
}

func (a *EtelecomAggregate) DeleteTenant(ctx context.Context, id dot.ID) error {
	_, err := a.tenantStore(ctx).ID(id).SoftDelete()
	return err
}

func (a *EtelecomAggregate) ActivateTenant(ctx context.Context, args *etelecom.ActivateTenantArgs) (*etelecom.Tenant, error) {
	// 1. Call api to connect the eB2B account to the tenant
	//      So that can get call logs & update tenant settings
	//
	// 2. Call api Portsip PBX to setting tenant
	//    - Update Trunk provider
	//    - Create outbound rules

	if args.TenantID == 0 && args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Please provide at least owner_id or tenant_id")
	}

	tenant, err := a.tenantStore(ctx).OptionalID(args.TenantID).OptionalOwnerID(args.OwnerID).ConnectionID(args.ConnectionID).GetTenant()
	if err != nil {
		return nil, err
	}
	tenantID := tenant.ID
	if tenant.Status.Enum == status3.P {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Tenant was activated")
	}

	event := &etelecom.TenantActivingEvent{
		TenantID:  tenantID,
		OwnerID:   tenant.OwnerID,
		HotlineID: args.HotlineID,
	}
	if err = a.eventBus.Publish(ctx, event); err != nil {
		return nil, err
	}

	// TODO: update tenant active
	update := &etelecom.Tenant{
		Status: status3.P.Wrap(),
	}
	if err = a.tenantStore(ctx).ID(tenantID).UpdateTenant(update); err != nil {
		return nil, err
	}
	return a.tenantStore(ctx).ID(tenantID).GetTenant()
}
