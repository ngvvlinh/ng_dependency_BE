package portsip_pm

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/etelecom/provider"
	providertypes "o.o/backend/com/etelecom/provider/types"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	identitymodel "o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/sqlstore"
	etelecomxserviceclient "o.o/backend/pkg/integration/telecom/etelecom_xservice/client"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll = l.New()
)

type ProcessManager struct {
	connectionManager *connectionmanager.ConnectionManager
	connectionQS      connectioning.QueryBus
	connectionAggr    connectioning.CommandBus
	telecomManager    *provider.TelecomManager
	etelecomQS        etelecom.QueryBus
	etelecomAggr      etelecom.CommandBus
	accountAuthStore  sqlstore.AccountAuthStoreFactory
}

func New(
	evenBus bus.EventRegistry,
	connManager *connectionmanager.ConnectionManager,
	connQ connectioning.QueryBus,
	connA connectioning.CommandBus,
	telecomManager *provider.TelecomManager,
	etelecomQ etelecom.QueryBus,
	etelecomA etelecom.CommandBus,
	accountAuthStore sqlstore.AccountAuthStoreFactory,
) *ProcessManager {
	p := &ProcessManager{
		connectionManager: connManager,
		connectionQS:      connQ,
		connectionAggr:    connA,
		telecomManager:    telecomManager,
		etelecomQS:        etelecomQ,
		etelecomAggr:      etelecomA,
		accountAuthStore:  accountAuthStore,
	}
	p.registerEventHandlers(evenBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.TenantActivating)
}

func (m *ProcessManager) TenantActivating(ctx context.Context, event *etelecom.TenantActivingEvent) error {
	// Call api to connect the eB2B account to the tenant
	// So that can get call logs & update tenant settings
	// Implement:
	//      - Create shop_connection type portsip direct if not existed
	//      - Call portsip api: updateTrunkProvider
	//      - Create outbound rule

	if event.AccountID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing account ID")
	}
	if event.HotlineID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing hotline_id")
	}
	queryTenant := &etelecom.GetTenantQuery{
		ID: event.TenantID,
	}
	if err := m.etelecomQS.Dispatch(ctx, queryTenant); err != nil {
		return err
	}
	tenant := queryTenant.Result
	if tenant.ExternalData == nil || tenant.ExternalData.ID == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "Tenant does not exist in Portsip PBX")
	}

	queryHotline := &etelecom.GetHotlineQuery{
		ID: event.HotlineID,
	}
	if err := m.etelecomQS.Dispatch(ctx, queryHotline); err != nil {
		return err
	}
	hotline := queryHotline.Result
	if hotline.Status != status3.Z {
		return cm.Errorf(cm.FailedPrecondition, nil, "Hotline does not valid")
	}
	if hotline.Hotline == "" {
		return cm.Errorf(cm.FailedPrecondition, nil, "Missing hotline number")
	}

	// get auth_key
	authKey, err := m.getShopPartnerAPIKey(ctx, event.AccountID)
	if err != nil {
		return err
	}

	// prepair shop_connection for get portsip driver
	queryShopConn := &connectioning.GetShopConnectionQuery{
		OwnerID:      event.OwnerID,
		ConnectionID: tenant.ConnectionID,
		IsGlobal:     false,
	}
	err = m.connectionQS.Dispatch(ctx, queryShopConn)
	switch cm.ErrorCode(err) {
	case cm.NoError:
	case cm.NotFound:
		// create shop_connection
		cmdShopConn := &connectioning.CreateShopConnectionCommand{
			OwnerID:      event.OwnerID,
			ConnectionID: tenant.ConnectionID,
			Token:        "default_token",
			// expires_at will update when call api portsip
			TokenExpiresAt: time.Now(),
			TelecomData: &connectioning.ShopConnectionTelecomData{
				Username:     tenant.Name,
				Password:     tenant.Password,
				TenantToken:  authKey,
				TenantDomain: tenant.Domain,
			},
		}
		if _err := m.connectionAggr.Dispatch(ctx, cmdShopConn); _err != nil {
			return _err
		}
	default:
		return err
	}

	// step 2: update portsip trunk provider
	update := &providertypes.AddTenantToTrunkProviderRequest{
		TrunkProviderID: m.telecomManager.AdminPortsip.TrunkProviderDefaultID,
		TenantID:        tenant.ExternalData.ID,
		Hotline:         hotline.Hotline,
	}
	if err = m.updateTrunkProvider(ctx, update); err != nil {
		return cm.Errorf(cm.ErrorCode(err), nil, "Error when update trunk provider: %v", err.Error())
	}

	// step 3: create outbound rule
	if err = m.createOutboundRule(ctx, tenant.OwnerID, tenant.ConnectionID); err != nil {
		return cm.Errorf(cm.ErrorCode(err), nil, "Error when create outbound rule: %v", err.Error())
	}

	// step 4: setting tenant to get CDR
	xClient := etelecomxserviceclient.New(authKey)
	xReq := &etelecomxserviceclient.ConfigTenantCDRRequest{
		Name:     tenant.Name,
		Password: tenant.Password,
	}
	_, err = xClient.ConfigTenantCDR(ctx, xReq)
	if err != nil {
		ll.Error("xService config tenant error", l.Error(err))
	}

	// step 5: update hotline
	updateHotline := &etelecom.UpdateHotlineInfoCommand{
		ID:               hotline.ID,
		Status:           status3.P.Wrap(),
		TenantID:         tenant.ID,
		ConnectionID:     tenant.ConnectionID,
		ConnectionMethod: tenant.ConnectionMethod,
		OwnerID:          event.OwnerID,
	}
	return m.etelecomAggr.Dispatch(ctx, updateHotline)
}

func (m *ProcessManager) updateTrunkProvider(ctx context.Context, args *providertypes.AddTenantToTrunkProviderRequest) error {
	portsipAdminDriver, err := m.telecomManager.GetAdministratorPortsipDriver(ctx)
	if err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Please check Portsip admin account")
	}

	return portsipAdminDriver.AddTenantToTrunkProvider(ctx, args)
}

func (m *ProcessManager) createOutboundRule(ctx context.Context, ownerID, connID dot.ID) error {
	args := &provider.CreateOutboundRuleRequest{
		TrunkProviderID: m.telecomManager.AdminPortsip.TrunkProviderDefaultID,
		OwnerID:         ownerID,
		ConnectionID:    connID,
	}
	return m.telecomManager.CreateOutboundRule(ctx, args)
}

func (m *ProcessManager) getShopPartnerAPIKey(ctx context.Context, accountID dot.ID) (string, error) {
	auth, err := m.accountAuthStore(ctx).AccountID(accountID).Get()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		return auth.AuthKey, nil
	case cm.NotFound:
		aa := &identitymodel.AccountAuth{
			AccountID: accountID,
			Status:    status3.P,
		}
		err = m.accountAuthStore(ctx).Create(aa)
		return aa.AuthKey, err
	default:
		return "", err
	}
}
