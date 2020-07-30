package webhook

import (
	"o.o/api/meta"
	"o.o/api/supporting/ticket"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var ll = l.New().WithChannel(meta.ChannelShipmentCarrier)

type MainDB *cmsql.Database // TODO(vu): call the right service

type Webhook struct {
	db              *cmsql.Database
	TicketAggregate ticket.CommandBus
	TicketQuery     ticket.QueryBus
}

func New(
	db com.MainDB,
	ticketAggregate ticket.CommandBus,
	ticketQuery ticket.QueryBus,
) *Webhook {
	wh := &Webhook{
		db:              db,
		TicketAggregate: ticketAggregate,
		TicketQuery:     ticketQuery,
	}
	return wh
}

func (wh *Webhook) Register(rt *httpx.Router) {
	rt.POST("/webhook/ghn/ticket", wh.Callback)
}

func (wh *Webhook) Callback(c *httpx.Context) (_err error) {
	return nil
}
