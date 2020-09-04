package shipping_provider

import (
	"context"

	"o.o/api/supporting/ticket"
)

type CarrierTicketDriver interface {
	CreateTicket(context.Context, *ticket.Ticket) (*ticket.Ticket, error)
	CreateComment(context.Context, *ticket.TicketComment) (*ticket.TicketComment, error)
}
