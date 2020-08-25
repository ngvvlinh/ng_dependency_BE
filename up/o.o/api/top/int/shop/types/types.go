package types

import (
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type TicketComment struct {
	ID        dot.ID `json:"id"`
	TicketID  dot.ID `json:"ticket_id"`
	CreatedBy dot.ID `json:"created_by"`
	AccountID dot.ID `json:"account_id"`
	ParentID  dot.ID `json:"parent_id"`

	Message  string `json:"message"`
	ImageUrl string `json:"image_url"`

	DeletedAt dot.Time `json:"deleted_at"`
	DeletedBy dot.ID   `json:"deleted_by"`

	CreatedAt dot.Time `json:"created_at"`
	UpdatedAt dot.Time `json:"updated_at"`
}

func (m *TicketComment) String() string { return jsonx.MustMarshalToString(m) }

type TicketLabel struct {
	ID       dot.ID         `json:"id"`
	Name     string         `json:"name"`
	Code     string         `json:"code"`
	Color    string         `json:"color"`
	ParentID dot.ID         `json:"parent_id"`
	Children []*TicketLabel `json:"children"`
}

func (m *TicketLabel) String() string { return jsonx.MustMarshalToString(m) }

type Ticket struct {
	ID              dot.ID   `json:"id"`
	Code            string   `json:"code"`
	AssignedUserIDs []dot.ID `json:"assigned_user_ids"`
	AccountID       dot.ID   `json:"account_id"`
	LabelIDs        []dot.ID `json:"label_ids"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
	AdminNote   string `json:"admin_note"`

	RefID   dot.ID                        `json:"ref_id"`
	RefType ticket_ref_type.TicketRefType `json:"ref_type"`
	RefCode string                        `json:"ref_code"`
	Source  ticket_source.TicketSource    `json:"source"`

	State  ticket_state.TicketState `json:"state"`
	Status status5.Status           `json:"status"`

	CreatedBy   dot.ID `json:"created_by"`
	UpdatedBy   dot.ID `json:"updated_by"`
	ConfirmedBy dot.ID `json:"confirmed_by"`
	ClosedBy    dot.ID `json:"closed_by"`

	CreatedAt   dot.Time `json:"created_at"`
	UpdatedAt   dot.Time `json:"updated_at"`
	ConfirmedAt dot.Time `json:"confirmed_at"`
	ClosedAt    dot.Time `json:"closed_at"`
}

func (m *Ticket) String() string { return jsonx.MustMarshalToString(m) }