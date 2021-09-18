// +build !generator

// Code generated by generator api. DO NOT EDIT.

package ticket

import (
	context "context"

	meta "o.o/api/meta"
	common "o.o/api/top/types/common"
	account_type "o.o/api/top/types/etc/account_type"
	status4 "o.o/api/top/types/etc/status4"
	ticket_ref_type "o.o/api/top/types/etc/ticket/ticket_ref_type"
	ticket_source "o.o/api/top/types/etc/ticket/ticket_source"
	ticket_state "o.o/api/top/types/etc/ticket/ticket_state"
	ticket_type "o.o/api/top/types/etc/ticket/ticket_type"
	capi "o.o/capi"
	dot "o.o/capi/dot"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type AssignTicketCommand struct {
	ID              dot.ID
	UpdatedBy       dot.ID
	IsLeader        bool
	AssignedUserIDs []dot.ID

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleAssignTicket(ctx context.Context, msg *AssignTicketCommand) (err error) {
	msg.Result, err = h.inner.AssignTicket(msg.GetArgs(ctx))
	return err
}

type CloseTicketCommand struct {
	IsLeader bool
	ID       dot.ID
	ClosedBy dot.ID
	Note     string
	State    ticket_state.TicketState

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleCloseTicket(ctx context.Context, msg *CloseTicketCommand) (err error) {
	msg.Result, err = h.inner.CloseTicket(msg.GetArgs(ctx))
	return err
}

type ConfirmTicketCommand struct {
	IsLeader  bool
	ID        dot.ID
	ConfirmBy dot.ID
	Note      string

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleConfirmTicket(ctx context.Context, msg *ConfirmTicketCommand) (err error) {
	msg.Result, err = h.inner.ConfirmTicket(msg.GetArgs(ctx))
	return err
}

type CreateTicketCommand struct {
	AssignedUserIDs []dot.ID
	AccountID       dot.ID
	LabelIDs        []dot.ID
	RefTicketID     dot.NullID
	Title           string
	Description     string
	Note            string
	AdminNote       string
	RefID           dot.ID
	RefType         ticket_ref_type.TicketRefType
	RefCode         string
	Source          ticket_source.TicketSource
	CreatedBy       dot.ID
	CreatedSource   account_type.AccountType
	CreatedName     string
	Type            ticket_type.TicketType

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleCreateTicket(ctx context.Context, msg *CreateTicketCommand) (err error) {
	msg.Result, err = h.inner.CreateTicket(msg.GetArgs(ctx))
	return err
}

type CreateTicketCommentCommand struct {
	CreatedBy     dot.ID
	CreatedName   string
	CreatedSource account_type.AccountType
	TicketID      dot.ID
	AccountID     dot.ID
	ParentID      dot.ID
	Message       string
	ImageUrls     []string
	IsLeader      bool
	IsAdmin       bool

	Result *TicketComment `json:"-"`
}

func (h AggregateHandler) HandleCreateTicketComment(ctx context.Context, msg *CreateTicketCommentCommand) (err error) {
	msg.Result, err = h.inner.CreateTicketComment(msg.GetArgs(ctx))
	return err
}

type CreateTicketLabelCommand struct {
	ShopID   dot.ID
	Type     ticket_type.TicketType
	Name     string
	Code     string
	Color    string
	ParentID dot.ID

	Result *TicketLabel `json:"-"`
}

func (h AggregateHandler) HandleCreateTicketLabel(ctx context.Context, msg *CreateTicketLabelCommand) (err error) {
	msg.Result, err = h.inner.CreateTicketLabel(msg.GetArgs(ctx))
	return err
}

type CreateTicketLabelExternalCommand struct {
	ConnectionID dot.ID
	ExternalID   string
	ExternalName string

	Result *TicketLabelExternal `json:"-"`
}

func (h AggregateHandler) HandleCreateTicketLabelExternal(ctx context.Context, msg *CreateTicketLabelExternalCommand) (err error) {
	msg.Result, err = h.inner.CreateTicketLabelExternal(msg.GetArgs(ctx))
	return err
}

type DeleteTicketCommentCommand struct {
	AccountID dot.ID
	ID        dot.ID
	IsAdmin   bool
	DeletedBy dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteTicketComment(ctx context.Context, msg *DeleteTicketCommentCommand) (err error) {
	msg.Result, err = h.inner.DeleteTicketComment(msg.GetArgs(ctx))
	return err
}

type DeleteTicketLabelCommand struct {
	ID          dot.ID
	ShopID      dot.ID
	Type        ticket_type.TicketType
	DeleteChild bool

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteTicketLabel(ctx context.Context, msg *DeleteTicketLabelCommand) (err error) {
	msg.Result, err = h.inner.DeleteTicketLabel(msg.GetArgs(ctx))
	return err
}

type DeleteTicketLabelExternalCommand struct {
	ID dot.ID

	Result int `json:"-"`
}

func (h AggregateHandler) HandleDeleteTicketLabelExternal(ctx context.Context, msg *DeleteTicketLabelExternalCommand) (err error) {
	msg.Result, err = h.inner.DeleteTicketLabelExternal(msg.GetArgs(ctx))
	return err
}

type ReopenTicketCommand struct {
	ID   dot.ID
	Note string

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleReopenTicket(ctx context.Context, msg *ReopenTicketCommand) (err error) {
	msg.Result, err = h.inner.ReopenTicket(msg.GetArgs(ctx))
	return err
}

type UnassignTicketCommand struct {
	ID        dot.ID
	UpdatedBy dot.ID

	Result *Ticket `json:"-"`
}

func (h AggregateHandler) HandleUnassignTicket(ctx context.Context, msg *UnassignTicketCommand) (err error) {
	msg.Result, err = h.inner.UnassignTicket(msg.GetArgs(ctx))
	return err
}

type UpdateTicketCommentCommand struct {
	AccountID dot.ID
	ID        dot.ID
	UpdatedBy dot.ID
	Message   string
	ImageUrls []string

	Result *TicketComment `json:"-"`
}

func (h AggregateHandler) HandleUpdateTicketComment(ctx context.Context, msg *UpdateTicketCommentCommand) (err error) {
	msg.Result, err = h.inner.UpdateTicketComment(msg.GetArgs(ctx))
	return err
}

type UpdateTicketInfoCommand struct {
	ID          dot.ID
	Code        string
	AccountID   dot.ID
	Labels      []dot.ID
	RefTicketID dot.ID
	Title       string
	Description string
	Note        dot.NullString
	RefID       dot.ID
	RefType     ticket_ref_type.TicketRefType
	Source      ticket_source.NullTicketSource
	Status      status4.Status
	State       string

	Result *common.UpdatedResponse `json:"-"`
}

func (h AggregateHandler) HandleUpdateTicketInfo(ctx context.Context, msg *UpdateTicketInfoCommand) (err error) {
	msg.Result, err = h.inner.UpdateTicketInfo(msg.GetArgs(ctx))
	return err
}

type UpdateTicketLabelCommand struct {
	ID       dot.ID
	ShopID   dot.ID
	Type     ticket_type.TicketType
	Color    string
	Name     dot.NullString
	Code     dot.NullString
	ParentID dot.NullID

	Result *TicketLabel `json:"-"`
}

func (h AggregateHandler) HandleUpdateTicketLabel(ctx context.Context, msg *UpdateTicketLabelCommand) (err error) {
	msg.Result, err = h.inner.UpdateTicketLabel(msg.GetArgs(ctx))
	return err
}

type UpdateTicketLabelExternalCommand struct {
	ID           dot.ID
	ExternalName string

	Result *TicketLabelExternal `json:"-"`
}

func (h AggregateHandler) HandleUpdateTicketLabelExternal(ctx context.Context, msg *UpdateTicketLabelExternalCommand) (err error) {
	msg.Result, err = h.inner.UpdateTicketLabelExternal(msg.GetArgs(ctx))
	return err
}

type UpdateTicketRefTicketIDCommand struct {
	ID          dot.ID
	RefTicketID dot.NullID

	Result *common.UpdatedResponse `json:"-"`
}

func (h AggregateHandler) HandleUpdateTicketRefTicketID(ctx context.Context, msg *UpdateTicketRefTicketIDCommand) (err error) {
	msg.Result, err = h.inner.UpdateTicketRefTicketID(msg.GetArgs(ctx))
	return err
}

type GetTicketByExternalIDQuery struct {
	ExternalID string

	Result *Ticket `json:"-"`
}

func (h QueryServiceHandler) HandleGetTicketByExternalID(ctx context.Context, msg *GetTicketByExternalIDQuery) (err error) {
	msg.Result, err = h.inner.GetTicketByExternalID(msg.GetArgs(ctx))
	return err
}

type GetTicketByIDQuery struct {
	ID              dot.ID
	AccountID       dot.ID
	AssignedUserIDs []dot.ID
	CreatedBy       dot.ID

	Result *Ticket `json:"-"`
}

func (h QueryServiceHandler) HandleGetTicketByID(ctx context.Context, msg *GetTicketByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTicketByID(msg.GetArgs(ctx))
	return err
}

type GetTicketCommentByIDQuery struct {
	ID        dot.ID
	AccountID dot.ID

	Result *TicketComment `json:"-"`
}

func (h QueryServiceHandler) HandleGetTicketCommentByID(ctx context.Context, msg *GetTicketCommentByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTicketCommentByID(msg.GetArgs(ctx))
	return err
}

type GetTicketLabelByIDQuery struct {
	ID dot.ID

	Result *TicketLabel `json:"-"`
}

func (h QueryServiceHandler) HandleGetTicketLabelByID(ctx context.Context, msg *GetTicketLabelByIDQuery) (err error) {
	msg.Result, err = h.inner.GetTicketLabelByID(msg.GetArgs(ctx))
	return err
}

type ListTicketCommentsQuery struct {
	Filter *FilterGetTicketComment
	Paging meta.Paging

	Result *GetTicketCommentsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListTicketComments(ctx context.Context, msg *ListTicketCommentsQuery) (err error) {
	msg.Result, err = h.inner.ListTicketComments(msg.GetArgs(ctx))
	return err
}

type ListTicketLabelsQuery struct {
	Type   ticket_type.NullTicketType
	ShopID dot.ID
	Tree   bool

	Result *GetTicketLabelsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListTicketLabels(ctx context.Context, msg *ListTicketLabelsQuery) (err error) {
	msg.Result, err = h.inner.ListTicketLabels(msg.GetArgs(ctx))
	return err
}

type ListTicketsQuery struct {
	Filter    *FilterGetTicket
	Paging    meta.Paging
	IsLeader  bool
	HasFilter bool

	Result *ListTicketsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListTickets(ctx context.Context, msg *ListTicketsQuery) (err error) {
	msg.Result, err = h.inner.ListTickets(msg.GetArgs(ctx))
	return err
}

type ListTicketsByRefTicketIDQuery struct {
	AccountID       dot.ID
	RefTicketID     dot.ID
	AssignedUserIDs []dot.ID
	CreatedBy       dot.ID

	Result []*Ticket `json:"-"`
}

func (h QueryServiceHandler) HandleListTicketsByRefTicketID(ctx context.Context, msg *ListTicketsByRefTicketIDQuery) (err error) {
	msg.Result, err = h.inner.ListTicketsByRefTicketID(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *AssignTicketCommand) command()              {}
func (q *CloseTicketCommand) command()               {}
func (q *ConfirmTicketCommand) command()             {}
func (q *CreateTicketCommand) command()              {}
func (q *CreateTicketCommentCommand) command()       {}
func (q *CreateTicketLabelCommand) command()         {}
func (q *CreateTicketLabelExternalCommand) command() {}
func (q *DeleteTicketCommentCommand) command()       {}
func (q *DeleteTicketLabelCommand) command()         {}
func (q *DeleteTicketLabelExternalCommand) command() {}
func (q *ReopenTicketCommand) command()              {}
func (q *UnassignTicketCommand) command()            {}
func (q *UpdateTicketCommentCommand) command()       {}
func (q *UpdateTicketInfoCommand) command()          {}
func (q *UpdateTicketLabelCommand) command()         {}
func (q *UpdateTicketLabelExternalCommand) command() {}
func (q *UpdateTicketRefTicketIDCommand) command()   {}

func (q *GetTicketByExternalIDQuery) query()    {}
func (q *GetTicketByIDQuery) query()            {}
func (q *GetTicketCommentByIDQuery) query()     {}
func (q *GetTicketLabelByIDQuery) query()       {}
func (q *ListTicketCommentsQuery) query()       {}
func (q *ListTicketLabelsQuery) query()         {}
func (q *ListTicketsQuery) query()              {}
func (q *ListTicketsByRefTicketIDQuery) query() {}

// implement conversion

func (q *AssignTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *AssignedTicketArgs) {
	return ctx,
		&AssignedTicketArgs{
			ID:              q.ID,
			UpdatedBy:       q.UpdatedBy,
			IsLeader:        q.IsLeader,
			AssignedUserIDs: q.AssignedUserIDs,
		}
}

func (q *AssignTicketCommand) SetAssignedTicketArgs(args *AssignedTicketArgs) {
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
	q.IsLeader = args.IsLeader
	q.AssignedUserIDs = args.AssignedUserIDs
}

func (q *CloseTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CloseTicketArgs) {
	return ctx,
		&CloseTicketArgs{
			IsLeader: q.IsLeader,
			ID:       q.ID,
			ClosedBy: q.ClosedBy,
			Note:     q.Note,
			State:    q.State,
		}
}

func (q *CloseTicketCommand) SetCloseTicketArgs(args *CloseTicketArgs) {
	q.IsLeader = args.IsLeader
	q.ID = args.ID
	q.ClosedBy = args.ClosedBy
	q.Note = args.Note
	q.State = args.State
}

func (q *ConfirmTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ConfirmTicketArgs) {
	return ctx,
		&ConfirmTicketArgs{
			IsLeader:  q.IsLeader,
			ID:        q.ID,
			ConfirmBy: q.ConfirmBy,
			Note:      q.Note,
		}
}

func (q *ConfirmTicketCommand) SetConfirmTicketArgs(args *ConfirmTicketArgs) {
	q.IsLeader = args.IsLeader
	q.ID = args.ID
	q.ConfirmBy = args.ConfirmBy
	q.Note = args.Note
}

func (q *CreateTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateTicketArgs) {
	return ctx,
		&CreateTicketArgs{
			AssignedUserIDs: q.AssignedUserIDs,
			AccountID:       q.AccountID,
			LabelIDs:        q.LabelIDs,
			RefTicketID:     q.RefTicketID,
			Title:           q.Title,
			Description:     q.Description,
			Note:            q.Note,
			AdminNote:       q.AdminNote,
			RefID:           q.RefID,
			RefType:         q.RefType,
			RefCode:         q.RefCode,
			Source:          q.Source,
			CreatedBy:       q.CreatedBy,
			CreatedSource:   q.CreatedSource,
			CreatedName:     q.CreatedName,
			Type:            q.Type,
		}
}

func (q *CreateTicketCommand) SetCreateTicketArgs(args *CreateTicketArgs) {
	q.AssignedUserIDs = args.AssignedUserIDs
	q.AccountID = args.AccountID
	q.LabelIDs = args.LabelIDs
	q.RefTicketID = args.RefTicketID
	q.Title = args.Title
	q.Description = args.Description
	q.Note = args.Note
	q.AdminNote = args.AdminNote
	q.RefID = args.RefID
	q.RefType = args.RefType
	q.RefCode = args.RefCode
	q.Source = args.Source
	q.CreatedBy = args.CreatedBy
	q.CreatedSource = args.CreatedSource
	q.CreatedName = args.CreatedName
	q.Type = args.Type
}

func (q *CreateTicketCommentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateTicketCommentArgs) {
	return ctx,
		&CreateTicketCommentArgs{
			CreatedBy:     q.CreatedBy,
			CreatedName:   q.CreatedName,
			CreatedSource: q.CreatedSource,
			TicketID:      q.TicketID,
			AccountID:     q.AccountID,
			ParentID:      q.ParentID,
			Message:       q.Message,
			ImageUrls:     q.ImageUrls,
			IsLeader:      q.IsLeader,
			IsAdmin:       q.IsAdmin,
		}
}

func (q *CreateTicketCommentCommand) SetCreateTicketCommentArgs(args *CreateTicketCommentArgs) {
	q.CreatedBy = args.CreatedBy
	q.CreatedName = args.CreatedName
	q.CreatedSource = args.CreatedSource
	q.TicketID = args.TicketID
	q.AccountID = args.AccountID
	q.ParentID = args.ParentID
	q.Message = args.Message
	q.ImageUrls = args.ImageUrls
	q.IsLeader = args.IsLeader
	q.IsAdmin = args.IsAdmin
}

func (q *CreateTicketLabelCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateTicketLabelArgs) {
	return ctx,
		&CreateTicketLabelArgs{
			ShopID:   q.ShopID,
			Type:     q.Type,
			Name:     q.Name,
			Code:     q.Code,
			Color:    q.Color,
			ParentID: q.ParentID,
		}
}

func (q *CreateTicketLabelCommand) SetCreateTicketLabelArgs(args *CreateTicketLabelArgs) {
	q.ShopID = args.ShopID
	q.Type = args.Type
	q.Name = args.Name
	q.Code = args.Code
	q.Color = args.Color
	q.ParentID = args.ParentID
}

func (q *CreateTicketLabelExternalCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateTicketLabelExternalArgs) {
	return ctx,
		&CreateTicketLabelExternalArgs{
			ConnectionID: q.ConnectionID,
			ExternalID:   q.ExternalID,
			ExternalName: q.ExternalName,
		}
}

func (q *CreateTicketLabelExternalCommand) SetCreateTicketLabelExternalArgs(args *CreateTicketLabelExternalArgs) {
	q.ConnectionID = args.ConnectionID
	q.ExternalID = args.ExternalID
	q.ExternalName = args.ExternalName
}

func (q *DeleteTicketCommentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteTicketCommentArgs) {
	return ctx,
		&DeleteTicketCommentArgs{
			AccountID: q.AccountID,
			ID:        q.ID,
			IsAdmin:   q.IsAdmin,
			DeletedBy: q.DeletedBy,
		}
}

func (q *DeleteTicketCommentCommand) SetDeleteTicketCommentArgs(args *DeleteTicketCommentArgs) {
	q.AccountID = args.AccountID
	q.ID = args.ID
	q.IsAdmin = args.IsAdmin
	q.DeletedBy = args.DeletedBy
}

func (q *DeleteTicketLabelCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteTicketLabelArgs) {
	return ctx,
		&DeleteTicketLabelArgs{
			ID:          q.ID,
			ShopID:      q.ShopID,
			Type:        q.Type,
			DeleteChild: q.DeleteChild,
		}
}

func (q *DeleteTicketLabelCommand) SetDeleteTicketLabelArgs(args *DeleteTicketLabelArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.Type = args.Type
	q.DeleteChild = args.DeleteChild
}

func (q *DeleteTicketLabelExternalCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DeleteTicketLabelExternalArgs) {
	return ctx,
		&DeleteTicketLabelExternalArgs{
			ID: q.ID,
		}
}

func (q *DeleteTicketLabelExternalCommand) SetDeleteTicketLabelExternalArgs(args *DeleteTicketLabelExternalArgs) {
	q.ID = args.ID
}

func (q *ReopenTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *ReopenTicketArgs) {
	return ctx,
		&ReopenTicketArgs{
			ID:   q.ID,
			Note: q.Note,
		}
}

func (q *ReopenTicketCommand) SetReopenTicketArgs(args *ReopenTicketArgs) {
	q.ID = args.ID
	q.Note = args.Note
}

func (q *UnassignTicketCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UnassignTicketArgs) {
	return ctx,
		&UnassignTicketArgs{
			ID:        q.ID,
			UpdatedBy: q.UpdatedBy,
		}
}

func (q *UnassignTicketCommand) SetUnassignTicketArgs(args *UnassignTicketArgs) {
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
}

func (q *UpdateTicketCommentCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateTicketCommentArgs) {
	return ctx,
		&UpdateTicketCommentArgs{
			AccountID: q.AccountID,
			ID:        q.ID,
			UpdatedBy: q.UpdatedBy,
			Message:   q.Message,
			ImageUrls: q.ImageUrls,
		}
}

func (q *UpdateTicketCommentCommand) SetUpdateTicketCommentArgs(args *UpdateTicketCommentArgs) {
	q.AccountID = args.AccountID
	q.ID = args.ID
	q.UpdatedBy = args.UpdatedBy
	q.Message = args.Message
	q.ImageUrls = args.ImageUrls
}

func (q *UpdateTicketInfoCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateTicketInfoArgs) {
	return ctx,
		&UpdateTicketInfoArgs{
			ID:          q.ID,
			Code:        q.Code,
			AccountID:   q.AccountID,
			Labels:      q.Labels,
			RefTicketID: q.RefTicketID,
			Title:       q.Title,
			Description: q.Description,
			Note:        q.Note,
			RefID:       q.RefID,
			RefType:     q.RefType,
			Source:      q.Source,
			Status:      q.Status,
			State:       q.State,
		}
}

func (q *UpdateTicketInfoCommand) SetUpdateTicketInfoArgs(args *UpdateTicketInfoArgs) {
	q.ID = args.ID
	q.Code = args.Code
	q.AccountID = args.AccountID
	q.Labels = args.Labels
	q.RefTicketID = args.RefTicketID
	q.Title = args.Title
	q.Description = args.Description
	q.Note = args.Note
	q.RefID = args.RefID
	q.RefType = args.RefType
	q.Source = args.Source
	q.Status = args.Status
	q.State = args.State
}

func (q *UpdateTicketLabelCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateTicketLabelArgs) {
	return ctx,
		&UpdateTicketLabelArgs{
			ID:       q.ID,
			ShopID:   q.ShopID,
			Type:     q.Type,
			Color:    q.Color,
			Name:     q.Name,
			Code:     q.Code,
			ParentID: q.ParentID,
		}
}

func (q *UpdateTicketLabelCommand) SetUpdateTicketLabelArgs(args *UpdateTicketLabelArgs) {
	q.ID = args.ID
	q.ShopID = args.ShopID
	q.Type = args.Type
	q.Color = args.Color
	q.Name = args.Name
	q.Code = args.Code
	q.ParentID = args.ParentID
}

func (q *UpdateTicketLabelExternalCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateTicketLabelExternalArgs) {
	return ctx,
		&UpdateTicketLabelExternalArgs{
			ID:           q.ID,
			ExternalName: q.ExternalName,
		}
}

func (q *UpdateTicketLabelExternalCommand) SetUpdateTicketLabelExternalArgs(args *UpdateTicketLabelExternalArgs) {
	q.ID = args.ID
	q.ExternalName = args.ExternalName
}

func (q *UpdateTicketRefTicketIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *UpdateTicketRefTicketIDArgs) {
	return ctx,
		&UpdateTicketRefTicketIDArgs{
			ID:          q.ID,
			RefTicketID: q.RefTicketID,
		}
}

func (q *UpdateTicketRefTicketIDCommand) SetUpdateTicketRefTicketIDArgs(args *UpdateTicketRefTicketIDArgs) {
	q.ID = args.ID
	q.RefTicketID = args.RefTicketID
}

func (q *GetTicketByExternalIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketByExternalIDArgs) {
	return ctx,
		&GetTicketByExternalIDArgs{
			ExternalID: q.ExternalID,
		}
}

func (q *GetTicketByExternalIDQuery) SetGetTicketByExternalIDArgs(args *GetTicketByExternalIDArgs) {
	q.ExternalID = args.ExternalID
}

func (q *GetTicketByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketByIDArgs) {
	return ctx,
		&GetTicketByIDArgs{
			ID:              q.ID,
			AccountID:       q.AccountID,
			AssignedUserIDs: q.AssignedUserIDs,
			CreatedBy:       q.CreatedBy,
		}
}

func (q *GetTicketByIDQuery) SetGetTicketByIDArgs(args *GetTicketByIDArgs) {
	q.ID = args.ID
	q.AccountID = args.AccountID
	q.AssignedUserIDs = args.AssignedUserIDs
	q.CreatedBy = args.CreatedBy
}

func (q *GetTicketCommentByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketCommentByIDArgs) {
	return ctx,
		&GetTicketCommentByIDArgs{
			ID:        q.ID,
			AccountID: q.AccountID,
		}
}

func (q *GetTicketCommentByIDQuery) SetGetTicketCommentByIDArgs(args *GetTicketCommentByIDArgs) {
	q.ID = args.ID
	q.AccountID = args.AccountID
}

func (q *GetTicketLabelByIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketLabelByIDArgs) {
	return ctx,
		&GetTicketLabelByIDArgs{
			ID: q.ID,
		}
}

func (q *GetTicketLabelByIDQuery) SetGetTicketLabelByIDArgs(args *GetTicketLabelByIDArgs) {
	q.ID = args.ID
}

func (q *ListTicketCommentsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketCommentsArgs) {
	return ctx,
		&GetTicketCommentsArgs{
			Filter: q.Filter,
			Paging: q.Paging,
		}
}

func (q *ListTicketCommentsQuery) SetGetTicketCommentsArgs(args *GetTicketCommentsArgs) {
	q.Filter = args.Filter
	q.Paging = args.Paging
}

func (q *ListTicketLabelsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketLabelsArgs) {
	return ctx,
		&GetTicketLabelsArgs{
			Type:   q.Type,
			ShopID: q.ShopID,
			Tree:   q.Tree,
		}
}

func (q *ListTicketLabelsQuery) SetGetTicketLabelsArgs(args *GetTicketLabelsArgs) {
	q.Type = args.Type
	q.ShopID = args.ShopID
	q.Tree = args.Tree
}

func (q *ListTicketsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetTicketsArgs) {
	return ctx,
		&GetTicketsArgs{
			Filter:    q.Filter,
			Paging:    q.Paging,
			IsLeader:  q.IsLeader,
			HasFilter: q.HasFilter,
		}
}

func (q *ListTicketsQuery) SetGetTicketsArgs(args *GetTicketsArgs) {
	q.Filter = args.Filter
	q.Paging = args.Paging
	q.IsLeader = args.IsLeader
	q.HasFilter = args.HasFilter
}

func (q *ListTicketsByRefTicketIDQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListTicketsByRefTicketIDArgs) {
	return ctx,
		&ListTicketsByRefTicketIDArgs{
			AccountID:       q.AccountID,
			RefTicketID:     q.RefTicketID,
			AssignedUserIDs: q.AssignedUserIDs,
			CreatedBy:       q.CreatedBy,
		}
}

func (q *ListTicketsByRefTicketIDQuery) SetListTicketsByRefTicketIDArgs(args *ListTicketsByRefTicketIDArgs) {
	q.AccountID = args.AccountID
	q.RefTicketID = args.RefTicketID
	q.AssignedUserIDs = args.AssignedUserIDs
	q.CreatedBy = args.CreatedBy
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleAssignTicket)
	b.AddHandler(h.HandleCloseTicket)
	b.AddHandler(h.HandleConfirmTicket)
	b.AddHandler(h.HandleCreateTicket)
	b.AddHandler(h.HandleCreateTicketComment)
	b.AddHandler(h.HandleCreateTicketLabel)
	b.AddHandler(h.HandleCreateTicketLabelExternal)
	b.AddHandler(h.HandleDeleteTicketComment)
	b.AddHandler(h.HandleDeleteTicketLabel)
	b.AddHandler(h.HandleDeleteTicketLabelExternal)
	b.AddHandler(h.HandleReopenTicket)
	b.AddHandler(h.HandleUnassignTicket)
	b.AddHandler(h.HandleUpdateTicketComment)
	b.AddHandler(h.HandleUpdateTicketInfo)
	b.AddHandler(h.HandleUpdateTicketLabel)
	b.AddHandler(h.HandleUpdateTicketLabelExternal)
	b.AddHandler(h.HandleUpdateTicketRefTicketID)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetTicketByExternalID)
	b.AddHandler(h.HandleGetTicketByID)
	b.AddHandler(h.HandleGetTicketCommentByID)
	b.AddHandler(h.HandleGetTicketLabelByID)
	b.AddHandler(h.HandleListTicketComments)
	b.AddHandler(h.HandleListTicketLabels)
	b.AddHandler(h.HandleListTickets)
	b.AddHandler(h.HandleListTicketsByRefTicketID)
	return QueryBus{b}
}
