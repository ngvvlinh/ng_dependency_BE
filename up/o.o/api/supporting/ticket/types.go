package ticket

import (
	"time"

	"o.o/api/top/types/etc/account_type"
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
	RefTicketID     dot.NullID // reference with another ticket

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

	CreatedBy     dot.ID
	CreatedSource account_type.AccountType
	CreatedName   string
	UpdatedBy     dot.ID
	ConfirmedBy   dot.ID
	ClosedBy      dot.ID

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ConfirmedAt time.Time
	ClosedAt    time.Time

	ConnectionID dot.ID
}

type TicketComment struct {
	ID                dot.ID
	TicketID          dot.ID
	CreatedBy         dot.ID
	CreatedName       string
	CreatedSource     account_type.AccountType
	AccountID         dot.ID
	ParentID          dot.ID
	ExternalCreatedAt string

	Message   string
	ImageUrls []string

	DeletedAt time.Time
	DeletedBy dot.ID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type TicketLabel struct {
	ID        dot.ID         `json:"id"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	ParentID  dot.ID         `json:"parent_id"`
	Color     string         `json:"color"`
	Children  []*TicketLabel `json:"children"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type TicketLabelExternal struct {
	ID           dot.ID
	ConnectionID dot.ID
	ExternalID   string
	ExternalName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TicketCreatingEvent struct {
	RefType   ticket_ref_type.TicketRefType
	RefID     dot.ID
	AccountID dot.ID
	RefCode   string
}
