package provider

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/main/connectioning"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	providertypes "o.o/backend/com/etelecom/provider/types"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

const (
	FiveMinutes             = 5 * time.Minute
	ExtensionPasswordLength = 10
)

type TelecomManager struct {
	eventBus          capi.EventBus
	connectionManager *connectionmanager.ConnectionManager
	telecomDriver     providertypes.Driver
	connectionQS      connectioning.QueryBus
	connectionAggr    connectioning.CommandBus
	identityQS        identity.QueryBus
	etelecomQS        etelecom.QueryBus
	AdminPortsip      providertypes.AdministratorTelecom
}

func NewTelecomManager(
	eventBus capi.EventBus,
	connectionManager *connectionmanager.ConnectionManager,
	telecomDriver providertypes.Driver,
	connectionQuery connectioning.QueryBus,
	connectionAggr connectioning.CommandBus,
	identityQuery identity.QueryBus,
	etelecomQuery etelecom.QueryBus,
	adminPortsip providertypes.AdministratorTelecom,
) (*TelecomManager, error) {
	return &TelecomManager{
		connectionManager: connectionManager,
		eventBus:          eventBus,
		telecomDriver:     telecomDriver,
		connectionQS:      connectionQuery,
		connectionAggr:    connectionAggr,
		identityQS:        identityQuery,
		etelecomQS:        etelecomQuery,
		AdminPortsip:      adminPortsip,
	}, nil
}

func (m *TelecomManager) GetTelecomDriver(ctx context.Context, connectionID, ownerID dot.ID) (providertypes.TelecomDriver, error) {
	connection, shopConnection, err := m.GetTelecomConnection(ctx, connectionID, ownerID)
	if err != nil {
		return nil, err
	}

	telecomDriver, err := m.telecomDriver.GetTelecomDriver(connection, shopConnection)
	if err != nil {
		return nil, err
	}

	// update token
	if err = m.generateToken(ctx, shopConnection, telecomDriver); err != nil {
		return nil, err
	}

	return telecomDriver, nil
}

func (m *TelecomManager) GetTelecomConnection(ctx context.Context, connectionID, ownerID dot.ID) (*connectioning.Connection, *connectioning.ShopConnection, error) {
	connection, err := m.connectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, nil, err
	}

	if connection.ConnectionType != connection_type.Telecom {
		return nil, nil, cm.Errorf(cm.FailedPrecondition, nil, "unsupported connection_type %v", connection.ConnectionType)
	}

	// connection
	getShopConnectionQuery := connectionmanager.GetShopConnectionArgs{
		ConnectionID: connectionID,
		OwnerID:      ownerID,
	}
	if connection.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		// ignore shopID
		getShopConnectionQuery.OwnerID = 0
		getShopConnectionQuery.IsGlobal = true
	}

	shopConnection, err := m.connectionManager.GetShopConnection(ctx, getShopConnectionQuery)
	if err != nil {
		return nil, nil, err
	}

	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid (check status or token)")
	}
	return connection, shopConnection, nil
}

func (m *TelecomManager) generateToken(ctx context.Context, shopConnection *connectioning.ShopConnection, telecomDriver providertypes.TelecomDriver) error {
	expiresAt := shopConnection.TokenExpiresAt
	if expiresAt.IsZero() {
		return nil
	}
	now := time.Now()
	// 5p trước khi hết hạn
	expiresAt.Add(-FiveMinutes)
	if expiresAt.After(now) {
		return nil
	}

	// re-generate token
	generateTokenResp, err := telecomDriver.GenerateToken(ctx)
	if err != nil {
		return err
	}

	// update shopConnection
	// shop connection telecom go with owner_id
	updateShopConnectionCmd := connectioning.CreateOrUpdateShopConnectionCommand{
		OwnerID:        shopConnection.OwnerID,
		ConnectionID:   shopConnection.ConnectionID,
		Token:          generateTokenResp.AccessToken,
		TokenExpiresAt: generateTokenResp.ExpiresAt,
		ExternalData:   shopConnection.ExternalData,
	}
	if err = m.connectionAggr.Dispatch(ctx, &updateShopConnectionCmd); err != nil {
		return err
	}
	*shopConnection = *updateShopConnectionCmd.Result
	return nil
}

func (m *TelecomManager) CreateExtension(ctx context.Context, ext *etelecom.Extension) (*CreateExtensionResponse, error) {
	hotlineQuery := &etelecom.GetHotlineQuery{
		ID: ext.HotlineID,
	}
	if err := m.etelecomQS.Dispatch(ctx, hotlineQuery); err != nil {
		return nil, err
	}
	hotline := hotlineQuery.Result

	driver, err := m.GetTelecomDriver(ctx, hotline.ConnectionID, hotline.OwnerID)
	if err != nil {
		return nil, err
	}

	var user *identity.User
	if ext.UserID != 0 {
		userQuery := &identity.GetUserByIDQuery{
			UserID: ext.UserID,
		}
		if err = m.identityQS.Dispatch(ctx, userQuery); err != nil {
			return nil, err
		}
		user = userQuery.Result
	}

	// get extension number
	extensionNumber := ext.ExtensionNumber
	if extensionNumber == "" {
		extQuery := &etelecom.GetPrivateExtensionNumberQuery{}
		if err = m.etelecomQS.Dispatch(ctx, extQuery); err != nil {
			return nil, err
		}
		extensionNumber = extQuery.Result
	}
	genPass := gencode.GenerateCode(gencode.Alphabet54, ExtensionPasswordLength)
	var profilePhone string
	profileName := hotline.Name
	if user != nil {
		profileName += " - " + user.FullName
		profilePhone = user.Phone
	}

	cmd := &providertypes.CreateExtensionRequest{
		ExtensionPassword: genPass,
		ExtensionNumber:   extensionNumber,
		Profile: &providertypes.ProfileExtension{
			FirstName: profileName,
			// Portsip định danh extension theo email
			// nên không thể tạo nhiều extension cho cùng 1 email được
			// Email:       user.Email,
			Phone:       profilePhone,
			Description: "",
		},
		Hotline: hotline.Hotline,
	}
	extResp, err := driver.CreateExtension(ctx, cmd)
	if err != nil {
		return nil, err
	}

	res := &CreateExtensionResponse{
		ExtensionID:       ext.ID,
		HotlineID:         hotline.ID,
		ExternalID:        extResp.ID,
		ExtensionNumber:   cmd.ExtensionNumber,
		ExtensionPassword: cmd.ExtensionPassword,
	}
	return res, nil
}

type getHotlineArgs struct {
	HotlineID    dot.ID
	OwnerID      dot.ID
	ConnectionID dot.ID
}

func (m *TelecomManager) getHotLine(ctx context.Context, args getHotlineArgs) (res *etelecom.Hotline, _ error) {
	conn, err := m.connectionManager.GetConnectionByID(ctx, args.ConnectionID)
	if err != nil {
		return nil, err
	}

	ownerID := args.OwnerID
	if conn.ConnectionMethod == connection_type.ConnectionMethodBuiltin {
		// builtin does not belong to owner_id
		ownerID = 0
	}
	if args.HotlineID != 0 {
		hotlineQuery := &etelecom.GetHotlineQuery{
			ID:      args.HotlineID,
			OwnerID: ownerID,
		}
		if err = m.etelecomQS.Dispatch(ctx, hotlineQuery); err != nil {
			return nil, err
		}
		return hotlineQuery.Result, nil
	}

	hotlineQuery := &etelecom.ListHotlinesQuery{
		OwnerID:      ownerID,
		ConnectionID: args.ConnectionID,
	}
	if err = m.etelecomQS.Dispatch(ctx, hotlineQuery); err != nil {
		return nil, err
	}
	hotlines := hotlineQuery.Result
	if len(hotlines) == 0 {
		return nil, cm.Errorf(cm.Internal, nil, "Vui lòng đăng ký hotline")
	}

	// get default hotline & make sure it's not empty
	for _, hotline := range hotlines {
		if hotline.Hotline != "" {
			res = hotline
			break
		}
	}
	if res == nil || res.Hotline == "" {
		return nil, cm.Errorf(cm.Internal, nil, "Vui lòng đăng ký số hotline")
	}
	return
}

func (m *TelecomManager) CreateOutboundRule(ctx context.Context, args *CreateOutboundRuleRequest) error {
	driver, err := m.GetTelecomDriver(ctx, args.ConnectionID, args.OwnerID)
	if err != nil {
		return err
	}
	req := &providertypes.CreateOutboundRuleRequest{
		TrunkProviderID: args.TrunkProviderID,
	}
	return driver.CreateOutboundRule(ctx, req)
}

func (m *TelecomManager) CreateTenantPortsip(ctx context.Context, tenant *etelecom.Tenant) (*CreateTenantResponse, error) {
	if tenant.Name == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing name")
	}
	if tenant.Password == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing password")
	}
	if tenant.Domain == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing domain")
	}

	userQuery := &identity.GetUserByIDQuery{
		UserID: tenant.OwnerID,
	}
	if err := m.identityQS.Dispatch(ctx, userQuery); err != nil {
		return nil, err
	}
	user := userQuery.Result

	req := &providertypes.CreateTenantRequest{
		Name:     tenant.Name,
		Domain:   tenant.Domain,
		Password: tenant.Password,
		Enable:   true,
		Profile: providertypes.TenantProfile{
			FirstName: "",
			LastName:  user.FullName,
			Email:     user.Email,
		},
	}

	driver, err := m.GetAdministratorPortsipDriver(ctx)
	if err != nil {
		return nil, err
	}
	tenantResp, err := driver.CreateTenant(ctx, req)
	if err != nil {
		return nil, err
	}
	return &CreateTenantResponse{
		TenantID:         tenant.ID,
		ExternalTenantID: tenantResp.ID,
	}, nil
}

func (m *TelecomManager) GetAdministratorPortsipDriver(ctx context.Context) (providertypes.TelecomAdminDriver, error) {
	expiresAt := m.AdminPortsip.TokenExpiresAt
	now := time.Now()
	// 5p trước khi hết hạn
	expiresAt.Add(-FiveMinutes)
	if expiresAt.After(now) {
		return m.AdminPortsip.Driver, nil
	}

	genTokenResp, err := m.AdminPortsip.Driver.GenerateToken(ctx)
	if err != nil {
		return nil, err
	}
	m.AdminPortsip.TokenExpiresAt = genTokenResp.ExpiresAt
	return m.AdminPortsip.Driver, nil
}
