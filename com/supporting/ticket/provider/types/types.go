package types

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/supporting/ticket"
)

type Driver interface {
	GetTicketDriver(
		env string,
		connection *connectioning.Connection,
		shopConnection *connectioning.ShopConnection,
	) (TicketProvider, error)
}

type TicketProvider interface {
	Ping(ctx context.Context) error
	CreateTicket(ctx context.Context, ticket *ticket.Ticket) (*ticket.Ticket, error)
	CreateTicketComment(ctx context.Context, ticket *ticket.TicketComment) (*ticket.TicketComment, error)
}
