package ticket

import (
	"context"

	"o.o/api/meta"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_source"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

// +gen:api

type Aggregate interface {
	// ticket
	CreateTicket(context.Context, *CreateTicketArgs) (*Ticket, error)
	UpdateTicketInfo(context.Context, *UpdateTicketInfoArgs) (*Ticket, error)
	ConfirmTicket(context.Context, *ConfirmTicketArgs) (*Ticket, error)
	CloseTicket(context.Context, *CloseTicketArgs) (*Ticket, error)
	ReopenTicket(context.Context, *ReopenTicketArgs) (*Ticket, error)
	AssignTicket(context.Context, *AssignedTicketArgs) (*Ticket, error)
	UnassignTicket(context.Context, *UnssignTicketArgs) (*Ticket, error)

	//comment
	CreateTicketComment(context.Context, *CreateTicketCommentArgs) (*TicketComment, error)
	UpdateTicketComment(context.Context, *UpdateTicketCommentArgs) (*TicketComment, error)
	DeleteTicketComment(context.Context, *DeleteTicketCommentArgs) (int, error)

	//label
	CreateTicketLabel(context.Context, *CreateTicketLabelArgs) (*TicketLabel, error)
	UpdateTicketLabel(context.Context, *UpdateTicketLabelArgs) (*TicketLabel, error)
	DeleteTicketLabel(context.Context, *DeleteTicketLabelArgs) (int, error)
}

type QueryService interface {
	// ticket
	GetTicketByID(context.Context, *GetTicketByIDArgs) (*Ticket, error)
	ListTickets(context.Context, *GetTicketsArgs) (*GetTicketsResponse, error)

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
	Tree bool
}

// +convert:create=TicketLabel
type CreateTicketLabelArgs struct {
	Name     string
	Code     string
	Color    string
	ParentID dot.ID
}

// +convert:update=TicketLabel(ID)
type UpdateTicketLabelArgs struct {
	ID       dot.ID
	Color    string
	Name     dot.NullString
	Code     dot.NullString
	ParentID dot.NullID
}

type DeleteTicketLabelArgs struct {
	ID          dot.ID
	DeleteChild bool
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
	ID        dot.ID
	AccountID dot.ID
}

type GetTicketsArgs struct {
	Filter *FilterGetTicket
	Paging meta.Paging
}

type GetTicketsResponse struct {
	Tickets []*Ticket
	Paging  meta.PageInfo
}

type FilterGetTicket struct {
	IDs             []dot.ID
	CreatedBy       dot.ID
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
}

type DeleteTicketCommentArgs struct {
	TicketID  dot.ID
	AccountID dot.ID
	ID        dot.ID

	DeletedBy dot.ID
}

// +convert:create=TicketComment
type CreateTicketCommentArgs struct {
	CreatedBy dot.ID
	TicketID  dot.ID
	AccountID dot.ID
	ParentID  dot.ID

	Message  string
	ImageUrl string

	IsLeader bool
	IsAdmin  bool
}

// only creator can update
// +convert:update=TicketComment(ID,TicketID,AccountID)
type UpdateTicketCommentArgs struct {
	AccountID dot.ID
	ID        dot.ID
	UpdatedBy dot.ID

	Message string
}

// từ điển https://en.wiktionary.org/wiki/unassign
type UnssignTicketArgs struct {
	ID        dot.ID
	UpdatedBy dot.ID
}

type AssignedTicketArgs struct {
	ID              dot.ID
	UpdatedBy       dot.ID
	IsLeader        bool
	AssignedUserIDs []dot.ID
}

type ReopenTicketArgs struct {
	ID   dot.ID
	Note string
}

type CloseTicketArgs struct {
	IsLeader bool
	ID       dot.ID
	ClosedBy dot.ID
	Note     string
	State    ticket_state.TicketState
}

type ConfirmTicketArgs struct {
	IsLeader  bool
	ID        dot.ID
	ConfirmBy dot.ID
	Note      string
}

// +convert:update=Ticket
type UpdateTicketInfoArgs struct {
	ID        dot.ID
	Code      string
	AccountID dot.ID
	Labels    []int8

	Title       dot.NullString
	Description dot.NullString
	Note        dot.NullString

	RefID   dot.NullID
	RefType ticket_ref_type.NullTicketRefType
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

	CreatedBy dot.ID
}
