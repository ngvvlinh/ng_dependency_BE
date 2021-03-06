package _all

import (
	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	telecomtypes "o.o/backend/com/etelecom/provider/types"
	cm "o.o/backend/pkg/common"
	portsipclient "o.o/backend/pkg/integration/telecom/portsip/client"
	portsipdriver "o.o/backend/pkg/integration/telecom/portsip/driver"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()

type TelecomDriver struct {
	eventBus capi.EventBus
}

func SupportedTelecomDriver(
	eventBus capi.EventBus,
) telecomtypes.Driver {
	return TelecomDriver{
		eventBus: eventBus,
	}
}

func (t TelecomDriver) GetTelecomDriver(
	connection *connectioning.Connection,
	shopConnection *connectioning.ShopConnection,
) (telecomtypes.TelecomDriver, error) {
	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderPortsip:
		if shopConnection.Token == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "token must not be null")
		}
		exData := shopConnection.TelecomData
		if exData == nil {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "telecom_data must not be null")
		}
		if exData.TenantToken == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "ShopConnection missing tenant token. connection_id = %v, owner_id = %v", shopConnection.ConnectionID, shopConnection.OwnerID)
		}
		if exData.Username == "" || exData.Password == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "ShopConnection Telecom missing username or password. connection_id = %v, owner_id = %v", shopConnection.ConnectionID, shopConnection.OwnerID)
		}
		// need tenant domain to telecom client sign in sip
		if exData.TenantDomain == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "ShopConnection Telecom missing tenant domain. connection_id = %v, owner_id = %v", shopConnection.ConnectionID, shopConnection.OwnerID)
		}
		cfg := portsipclient.PortsipAccountCfg{
			Password:    shopConnection.TelecomData.Password,
			Token:       shopConnection.Token,
			Username:    shopConnection.TelecomData.Username,
			TenantHost:  shopConnection.TelecomData.TenantHost,
			TenantToken: shopConnection.TelecomData.TenantToken,
		}
		driver := portsipdriver.New(cfg)
		return driver, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection kh??ng h???p l???")
	}
}
