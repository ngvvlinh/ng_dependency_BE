package ticket

import (
	"time"

	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/capi/dot"
)

// +gen:event:topic=event/Ticket

type Ticket struct {
	ExternalShippingCode string
	ExternalID           string

	ID              dot.ID
	Code            string
	AssignedUserIDs []dot.ID
	AccountID       dot.ID
	LabelIDs        []dot.ID

	Title       string
	Description string
	Note        string
	AdminNote   string

	RefID   dot.ID
	RefType ticket_ref_type.TicketRefType
	RefCode string
	Source  ticket_source.TicketSource

	State  ticket_state.TicketState
	Status status5.Status

	CreatedBy   dot.ID
	UpdatedBy   dot.ID
	ConfirmedBy dot.ID
	ClosedBy    dot.ID

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	ClosedAt    time.Time
}

type TicketComment struct {
	ID        dot.ID
	TicketID  dot.ID
	CreatedBy dot.ID
	AccountID dot.ID
	ParentID  dot.ID

	Message  string
	ImageUrl string

	DeletedAt time.Time
	DeletedBy dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type TicketLabel struct {
	ID       dot.ID         `json:"id"`
	Name     string         `json:"name"`
	Code     string         `json:"code"`
	ParentID dot.ID         `json:"parent_id"`
	Color    string         `json:"color"`
	Children []*TicketLabel `json:"children"`
}

type TicketCreatingEvent struct {
	RefType   ticket_ref_type.TicketRefType
	RefID     dot.ID
	AccountID dot.ID
	RefCode   string
}