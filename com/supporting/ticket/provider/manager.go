package provider

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/moneytx"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	"o.o/api/supporting/ticket"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/ticket/ticket_source"
	connectionmanager "o.o/backend/com/main/connectioning/manager"
	carriertypes "o.o/backend/com/supporting/ticket/provider/types"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/capi"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

type TicketManager struct {
	env               string
	connectionManager *connectionmanager.ConnectionManager
	eventBus          capi.EventBus
	ticketDriver      carriertypes.Driver
	connectionQS      connectioning.QueryBus
	shippingQS        shipping.QueryBus
	moneyTxQS         moneytx.QueryBus
	orderQS           ordering.QueryBus
}

func NewTicketManager(
	connectionManager *connectionmanager.ConnectionManager,
	eventBus capi.EventBus,
	ticketDriver carriertypes.Driver,
	connectionQuery connectioning.QueryBus,
) (*TicketManager, error) {
	return &TicketManager{
		connectionManager: connectionManager,
		eventBus:          eventBus,
		ticketDriver:      ticketDriver,
		connectionQS:      connectionQuery,
		env:               cmenv.PartnerEnv(),
	}, nil
}

func (m *TicketManager) GetTicketDriver(ctx context.Context, connectionID dot.ID) (carriertypes.TicketProvider, error) {
	connection, err := m.connectionManager.GetConnectionByID(ctx, connectionID)
	if err != nil {
		return nil, err
	}

	if connection.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "unsupported connection_method %v", connection.ConnectionMethod)
	}
	if connection.ConnectionType != connection_type.CRM {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "unsupported connection_type %v", connection.ConnectionType)
	}

	// connection
	getShopConnectionQuery := connectionmanager.GetShopConnectionArgs{
		ConnectionID: connectionID,
		IsGlobal:     true,
	}
	shopConnection, err := m.connectionManager.GetShopConnection(ctx, getShopConnectionQuery)
	if err != nil {
		return nil, err
	}

	if shopConnection.Status != status3.P || shopConnection.Token == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Connection does not valid (check status or token)")
	}

	return m.ticketDriver.GetTicketDriver(m.env, connection, shopConnection)
}

func (m *TicketManager) CreateTicket(ctx context.Context, ticketCore *ticket.Ticket) (*ticket.Ticket,  error) {
	if ticketCore.ConnectionID == 0 {
		return ticketCore, nil
	}
	switch ticketCore.Source {
	case ticket_source.WebPhone:
		driver, err := m.GetTicketDriver(ctx, ticketCore.ConnectionID)
		if err != nil {
			return nil, err
		}

		ticketCore, err = driver.CreateTicket(ctx, ticketCore)
		if err != nil {
			return nil, err
		}
	default:
		//no-op
	}

	return ticketCore, nil
}
