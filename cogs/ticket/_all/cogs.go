package ticket_all

import (
	"o.o/api/main/connectioning"
	"o.o/api/main/contact"
	"o.o/api/main/shipping"
	"o.o/api/top/types/etc/connection_type"
	tickettypes "o.o/backend/com/supporting/ticket/provider/types"
	cm "o.o/backend/pkg/common"
	suitecrmclient "o.o/backend/pkg/integration/crm/suitecrm/client"
	suitecrmdriver "o.o/backend/pkg/integration/crm/suitecrm/driver"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()

type TicketDriver struct {
	shippingQS shipping.QueryBus
	contactQS  contact.QueryBus
	eventBus   capi.EventBus
}

func SupportedTicketDriver(
	eventBus capi.EventBus,
	shippingQuery shipping.QueryBus,
	contactQuery contact.QueryBus,
) tickettypes.Driver {
	return TicketDriver{
		eventBus:   eventBus,
		shippingQS: shippingQuery,
		contactQS:  contactQuery,
	}
}

func (d TicketDriver) GetTicketDriver(
	env string,
	connection *connectioning.Connection,
	shopConnection *connectioning.ShopConnection,
) (tickettypes.TicketProvider, error) {
	switch connection.ConnectionProvider {
	case connection_type.ConnectionProviderSuiteCRM:
		if shopConnection.Token == "" {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "token must not be null")
		}
		cfg := &suitecrmclient.SuiteCRMCfg{
			Token: shopConnection.Token,
		}
		driver := suitecrmdriver.New(env, cfg, d.shippingQS, d.contactQS)
		return driver, nil
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection không hợp lệ")
	}
}
