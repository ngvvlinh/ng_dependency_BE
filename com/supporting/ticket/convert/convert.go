package convert

import (
	"time"

	"o.o/api/supporting/ticket"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/ticket/ticket_state"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/code/gencode"
)

// +gen:convert: o.o/backend/com/supporting/ticket/model -> o.o/api/supporting/ticket
// +gen:convert: o.o/api/supporting/ticket

func createLabel(in *ticket.CreateTicketLabelArgs, out *ticket.TicketLabel) {
	apply_ticket_CreateTicketLabelArgs_ticket_TicketLabel(in, out)
	out.ID = cm.NewID()
}

func createTicket(in *ticket.CreateTicketArgs, out *ticket.Ticket) {
	apply_ticket_CreateTicketArgs_ticket_Ticket(in, out)
	out.UpdatedBy = in.CreatedBy
	out.Status = status5.Z
	out.State = ticket_state.New
	out.Code = gencode.GenerateCodeWithChar("TK", time.Now())
	out.ID = cm.NewID()
}

func createComment(in *ticket.CreateTicketCommentArgs, out *ticket.TicketComment) {
	apply_ticket_CreateTicketCommentArgs_ticket_TicketComment(in, out)
	out.ID = cm.NewID()
}
