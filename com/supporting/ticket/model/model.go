package model

import (
	"time"

	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/capi/dot"
)

// +sqlgen
type Ticket struct {
	ID              dot.ID
	Code            string
	AssignedUserIDs []dot.ID
	AccountID       dot.ID
	LabelIDs        []dot.ID
	ExternalID      string
	RefTicketID     dot.NullID // reference with another ticket

	Title       string
	Description string
	Note        string

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

	CreatedAt   time.Time `sq:"create"`
	UpdatedAt   time.Time `sq:"update"`
	ConfirmedAt time.Time
	ClosedAt    time.Time

	WLPartnerID  dot.ID
	ConnectionID dot.ID
}

// +sqlgen:           Ticket    as t
// +sqlgen:left-join: TicketSearch as ts on t.id = ts.id
type TicketExtended struct {
	*Ticket
	TicketSearch *TicketSearch
}

// +sqlgen
type TicketSearch struct {
	ID        dot.ID
	TitleNorm string
}

// +sqlgen
type TicketComment struct {
	ID            dot.ID
	TicketID      dot.ID
	CreatedBy     dot.ID
	CreatedSource account_type.AccountType
	CreatedName   string
	AccountID     dot.ID
	ParentID      dot.ID

	ExternalID        string
	ExternalCreatedAt time.Time

	Message   string
	ImageUrls []string

	DeletedAt time.Time
	DeletedBy dot.ID

	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
}

// +sqlgen
type TicketLabel struct {
	ID        dot.ID
	Name      string
	Code      string
	Color     string
	ParentID  dot.ID
	CreatedAt time.Time `sq:"create"`
	UpdatedAt time.Time `sq:"update"`
	DeletedAt time.Time

	WLPartnerID dot.ID
}

// +sqlgen
type TicketLabelTicketLabelExternal struct {
	TicketLabelID         dot.ID
	TicketLabelExternalID dot.ID
	DeletedAt             time.Time

	TicketLabel         *TicketLabel         `sq:"-"`
	TicketLabelExternal *TicketLabelExternal `sq:"-"`
}

// +sqlgen
type TicketLabelExternal struct {
	ID           dot.ID
	ConnectionID dot.ID
	ExternalID   string
	ExternalName string
	CreatedAt    time.Time `sq:"create"`
	UpdatedAt    time.Time `sq:"update"`
	DeletedAt    time.Time
}
