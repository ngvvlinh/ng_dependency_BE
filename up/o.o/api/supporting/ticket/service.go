package ticket

import (
	"context"
	"time"

	"o.o/api/meta"
	cm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/api/top/types/etc/ticket/ticket_type"
	"o.o/capi/dot"
	"o.o/capi/filter"
	"o.o/common/xerrors"
)

// +gen:api

type Aggregate interface {
	// ticket
	CreateTicket(context.Context, *CreateTicketArgs) (*Ticket, error)
	UpdateTicketInfo(context.Context, *UpdateTicketInfoArgs) (*cm.UpdatedResponse, error)
	ConfirmTicket(context.Context, *ConfirmTicketArgs) (*Ticket, error)
	CloseTicket(context.Context, *CloseTicketArgs) (*Ticket, error)
	ReopenTicket(context.Context, *ReopenTicketArgs) (*Ticket, error)
	AssignTicket(context.Context, *AssignedTicketArgs) (*Ticket, error)
	UnassignTicket(context.Context, *UnassignTicketArgs) (*Ticket, error)
	UpdateTicketRefTicketID(context.Context, *UpdateTicketRefTicketIDArgs) (*cm.UpdatedResponse, error)

	//comment
	CreateTicketComment(context.Context, *CreateTicketCommentArgs) (*TicketComment, error)
	UpdateTicketComment(context.Context, *UpdateTicketCommentArgs) (*TicketComment, error)
	DeleteTicketComment(context.Context, *DeleteTicketCommentArgs) (int, error)

	//label
	CreateTicketLabel(context.Context, *CreateTicketLabelArgs) (*TicketLabel, error)
	UpdateTicketLabel(context.Context, *UpdateTicketLabelArgs) (*TicketLabel, error)
	DeleteTicketLabel(context.Context, *DeleteTicketLabelArgs) (int, error)

	CreateTicketLabelExternal(context.Context, *CreateTicketLabelExternalArgs) (*TicketLabelExternal, error)
	UpdateTicketLabelExternal(context.Context, *UpdateTicketLabelExternalArgs) (*TicketLabelExternal, error)
	DeleteTicketLabelExternal(context.Context, *DeleteTicketLabelExternalArgs) (int, error)
}

type QueryService interface {
	// ticket
	GetTicketByID(context.Context, *GetTicketByIDArgs) (*Ticket, error)
	GetTicketByExternalID(context.Context, *GetTicketByExternalIDArgs) (*Ticket, error)
	ListTickets(context.Context, *GetTicketsArgs) (*ListTicketsResponse, error)
	ListTicketsByRefTicketID(context.Context, *ListTicketsByRefTicketIDArgs) ([]*Ticket, error)

	//comment
	GetTicketCommentByID(context.Context, *GetTicketCommentByIDArgs) (*TicketComment, error)
	ListTicketComments(context.Context, *GetTicketCommentsArgs) (*GetTicketCommentsResponse, error)

	// label
	GetTicketLabelByID(context.Context, *GetTicketLabelByIDArgs) (*TicketLabel, error)
	ListTicketLabels(context.Context, *GetTicketLabelsArgs) (*GetTicketLabelsResponse, error)
}

type GetTicketLabelsResponse struct {
	TicketLabels []*TicketLabel
}

type GetTicketLabelByIDArgs struct {
	ID dot.ID
}

type GetTicketLabelsArgs struct {
	Type   ticket_type.NullTicketType
	ShopID dot.ID
	Tree   bool
}

// +convert:create=TicketLabel
type CreateTicketLabelArgs struct {
	ShopID   dot.ID
	Type     ticket_type.TicketType
	Name     string
	Code     string
	Color    string
	ParentID dot.ID
}

// +convert:update=TicketLabel(ID)
type UpdateTicketLabelArgs struct {
	ID       dot.ID
	ShopID   dot.ID
	Type     ticket_type.TicketType
	Color    string
	Name     dot.NullString
	Code     dot.NullString
	ParentID dot.NullID
}

type DeleteTicketLabelArgs struct {
	ID          dot.ID
	ShopID      dot.ID
	Type        ticket_type.TicketType
	DeleteChild bool
}

// +convert:create=TicketLabelExternal
type CreateTicketLabelExternalArgs struct {
	ConnectionID dot.ID
	ExternalID   string
	ExternalName string
}

// +convert:update=TicketLabelExternal
type UpdateTicketLabelExternalArgs struct {
	ID           dot.ID
	ExternalName string
}

type DeleteTicketLabelExternalArgs struct {
	ID dot.ID
}

type GetTicketCommentByIDArgs struct {
	ID        dot.ID
	AccountID dot.ID
}

type GetTicketCommentsArgs struct {
	Filter *FilterGetTicketComment
	Paging meta.Paging
}

type FilterGetTicketComment struct {
	IDs       []dot.ID
	Title     string
	CreatedBy dot.ID
	ParentID  dot.ID
	AccountID dot.ID
	TicketID  dot.ID
}

type GetTicketCommentsResponse struct {
	TicketComments []*TicketComment
	Paging         meta.PageInfo
}

type GetTicketByIDArgs struct {
	ID              dot.ID
	AccountID       dot.ID
	AssignedUserIDs []dot.ID
	CreatedBy       dot.ID
}

type GetTicketByExternalIDArgs struct {
	ExternalID string
}

type GetTicketsArgs struct {
	Filter *FilterGetTicket
	Paging meta.Paging
}

type ListTicketsByRefTicketIDArgs struct {
	AccountID       dot.ID
	RefTicketID     dot.ID
	AssignedUserIDs []dot.ID
	CreatedBy       dot.ID
}

type ListTicketsResponse struct {
	Tickets []*Ticket
	Paging  meta.PageInfo
}

type FilterGetTicket struct {
	IDs             []dot.ID
	CreatedBy       dot.ID
	ConfirmedBy     dot.ID
	ClosedBy        dot.ID
	AccountID       dot.ID
	LabelIDs        []dot.ID
	Title           filter.FullTextSearch
	AssignedUserIDs []dot.ID
	RefID           dot.ID
	RefType         ticket_ref_type.TicketRefType
	Code            string
	State           ticket_state.TicketState
	RefCode         string
	Types           []ticket_type.TicketType
}

type DeleteTicketCommentArgs struct {
	AccountID dot.ID
	ID        dot.ID
	IsAdmin   bool

	DeletedBy dot.ID
}

// +convert:create=TicketComment
type CreateTicketCommentArgs struct {
	CreatedBy     dot.ID
	CreatedName   string
	CreatedSource account_type.AccountType
	TicketID      dot.ID
	AccountID     dot.ID
	ParentID      dot.ID

	Message   string
	ImageUrls []string

	IsLeader bool
	IsAdmin  bool
}

// +convert:create=TicketComment
type CreateTicketCommentWebhookArgs struct {
	TicketID          dot.ID
	AccountID         dot.ID
	ParentID          dot.ID
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExternalCreatedAt string
	Message           string
	ImageUrl          string
}

// only creator can update
// +convert:update=TicketComment(ID,TicketID,AccountID)
type UpdateTicketCommentArgs struct {
	AccountID dot.ID
	ID        dot.ID
	UpdatedBy dot.ID

	Message   string
	ImageUrls []string
}

// từ điển https://en.wiktionary.org/wiki/unassign
type UnassignTicketArgs struct {
	ID        dot.ID
	UpdatedBy dot.ID
}

type AssignedTicketArgs struct {
	ID              dot.ID
	UpdatedBy       dot.ID
	IsLeader        bool // for ticket: system
	AssignedUserIDs []dot.ID
}

type ReopenTicketArgs struct {
	ID   dot.ID
	Note string
}

type CloseTicketArgs struct {
	IsLeader bool // for ticket: system
	ID       dot.ID
	ClosedBy dot.ID
	Note     string
	State    ticket_state.TicketState
}

type ConfirmTicketArgs struct {
	IsLeader  bool // for ticket: system
	ID        dot.ID
	ConfirmBy dot.ID
	Note      string
}

// +convert:update=Ticket
type UpdateTicketInfoArgs struct {
	ID          dot.ID
	Code        string
	AccountID   dot.ID
	Labels      []dot.ID
	RefTicketID dot.ID

	Title       string
	Description string
	Note        dot.NullString

	RefID   dot.ID
	RefType ticket_ref_type.TicketRefType
	Source  ticket_source.NullTicketSource

	// new, accepted, processing, closed
	Status status4.Status
	// States closed: success, fail, ignore
	State string
}

// +convert:create=Ticket
type CreateTicketArgs struct {
	AssignedUserIDs []dot.ID
	AccountID       dot.ID
	LabelIDs        []dot.ID
	RefTicketID     dot.NullID

	Title       string
	Description string

	// user note
	Note string
	// admin note
	AdminNote string

	RefID   dot.ID
	RefType ticket_ref_type.TicketRefType
	RefCode string
	Source  ticket_source.TicketSource

	CreatedBy     dot.ID
	CreatedSource account_type.AccountType
	CreatedName   string

	Type ticket_type.TicketType
}

type UpdateTicketRefTicketIDArgs struct {
	ID          dot.ID
	RefTicketID dot.NullID
}

func (m *CreateTicketLabelExternalArgs) Validate() error {
	if m.ConnectionID == 0 {
		return EditErrorMsg("ConnectionID")
	}

	if m.ExternalID == "" {
		return EditErrorMsg("ExternalID")
	}

	if m.ExternalName == "" {
		return EditErrorMsg("ExternalName")
	}

	return nil
}

func (m *UpdateTicketLabelExternalArgs) Validate() error {
	if m.ExternalName == "" {
		return EditErrorMsg("ExternalName")
	}

	if m.ID == 0 {
		return EditErrorMsg("ID")
	}

	return nil
}

func EditErrorMsg(str string) error {
	return xerrors.Errorf(xerrors.InvalidArgument, nil, "Vui lòng nhập thông tin bắt buộc, thiếu %v", str)
}
