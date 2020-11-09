// +build wireinject

package ticket_all

import (
	"github.com/google/wire"

	"o.o/backend/com/etc/logging/ticket"
	ticketcarrier "o.o/backend/com/supporting/ticket/provider"
)

var WireSet = wire.NewSet(
	ticket.WireSet,
	ticketcarrier.WireSet,
	SupportedTicketDriver,
)
