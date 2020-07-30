package ghn

import (
	"context"

	"o.o/api/supporting/ticket"
	"o.o/backend/pkg/etop/logic/shipping_provider"
	ghnticket "o.o/backend/pkg/integration/shipping/ghn/ticket"
)

var _ shipping_provider.CarrierTicketDriver = &CarrierTicket{}

type CarrierTicket struct {
	clientTicket *ghnticket.Client
}

func NewCarrierTicket(cfg Config) *CarrierTicket {
	client := ghnticket.New("env")
	return &CarrierTicket{
		clientTicket: client,
	}
}

func (c CarrierTicket) CreateTicket(ctx context.Context, ticket *ticket.Ticket) (*ticket.Ticket, error) {
	_, err := c.clientTicket.CreateTicket(ctx, &ghnticket.CreateTicketRequest{
		OrderCode: ticket.ExternalShippingCode,
		//TODO(Nam) map
		Category:    "",
		Description: ticket.Description,
	})
	return ticket, err
}

func (c CarrierTicket) CreateComment(ctx context.Context, comment *ticket.TicketComment) (*ticket.TicketComment, error) {
	panic("implement me")
}
