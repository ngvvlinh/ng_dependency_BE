// +build !generator

// Code generated by generator api. DO NOT EDIT.

package fbmessaging

import (
	context "context"

	fb_customer_conversation_type "o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	meta "o.o/api/meta"
	capi "o.o/capi"
	dot "o.o/capi/dot"
	filter "o.o/capi/filter"
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

type CreateFbCustomerConversationsCommand struct {
	FbCustomerConversations []*CreateFbCustomerConversationArgs

	Result []*FbCustomerConversation `json:"-"`
}

func (h AggregateHandler) HandleCreateFbCustomerConversations(ctx context.Context, msg *CreateFbCustomerConversationsCommand) (err error) {
	msg.Result, err = h.inner.CreateFbCustomerConversations(msg.GetArgs(ctx))
	return err
}

type CreateFbExternalConversationsCommand struct {
	FbExternalConversations []*CreateFbExternalConversationArgs

	Result []*FbExternalConversation `json:"-"`
}

func (h AggregateHandler) HandleCreateFbExternalConversations(ctx context.Context, msg *CreateFbExternalConversationsCommand) (err error) {
	msg.Result, err = h.inner.CreateFbExternalConversations(msg.GetArgs(ctx))
	return err
}

type CreateFbExternalMessagesCommand struct {
	FbExternalMessages []*CreateFbExternalMessageArgs

	Result []*FbExternalMessage `json:"-"`
}

func (h AggregateHandler) HandleCreateFbExternalMessages(ctx context.Context, msg *CreateFbExternalMessagesCommand) (err error) {
	msg.Result, err = h.inner.CreateFbExternalMessages(msg.GetArgs(ctx))
	return err
}

type CreateFbExternalPostsCommand struct {
	FbExternalPosts []*CreateFbExternalPostArgs

	Result []*FbExternalPost `json:"-"`
}

func (h AggregateHandler) HandleCreateFbExternalPosts(ctx context.Context, msg *CreateFbExternalPostsCommand) (err error) {
	msg.Result, err = h.inner.CreateFbExternalPosts(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateFbCustomerConversationsCommand struct {
	FbCustomerConversations []*CreateFbCustomerConversationArgs

	Result []*FbCustomerConversation `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateFbCustomerConversations(ctx context.Context, msg *CreateOrUpdateFbCustomerConversationsCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateFbCustomerConversations(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateFbExternalCommentsCommand struct {
	FbExternalComments []*CreateFbExternalCommentArgs

	Result []*FbExternalComment `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateFbExternalComments(ctx context.Context, msg *CreateOrUpdateFbExternalCommentsCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateFbExternalComments(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateFbExternalConversationsCommand struct {
	FbExternalConversations []*CreateFbExternalConversationArgs

	Result []*FbExternalConversation `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateFbExternalConversations(ctx context.Context, msg *CreateOrUpdateFbExternalConversationsCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateFbExternalConversations(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateFbExternalMessagesCommand struct {
	FbExternalMessages []*CreateFbExternalMessageArgs

	Result []*FbExternalMessage `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateFbExternalMessages(ctx context.Context, msg *CreateOrUpdateFbExternalMessagesCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateFbExternalMessages(msg.GetArgs(ctx))
	return err
}

type CreateOrUpdateFbExternalPostsCommand struct {
	FbExternalPosts []*CreateFbExternalPostArgs

	Result []*FbExternalPost `json:"-"`
}

func (h AggregateHandler) HandleCreateOrUpdateFbExternalPosts(ctx context.Context, msg *CreateOrUpdateFbExternalPostsCommand) (err error) {
	msg.Result, err = h.inner.CreateOrUpdateFbExternalPosts(msg.GetArgs(ctx))
	return err
}

type UpdateIsReadCustomerConversationCommand struct {
	ConversationCustomerID dot.ID
	IsRead                 bool

	Result int `json:"-"`
}

func (h AggregateHandler) HandleUpdateIsReadCustomerConversation(ctx context.Context, msg *UpdateIsReadCustomerConversationCommand) (err error) {
	msg.Result, err = h.inner.UpdateIsReadCustomerConversation(msg.GetArgs(ctx))
	return err
}

type GetFbCustomerConversationQuery struct {
	CustomerConversationType fb_customer_conversation_type.FbCustomerConversationType
	ExternalID               string
	ExternalUserID           string

	Result *FbCustomerConversation `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbCustomerConversation(ctx context.Context, msg *GetFbCustomerConversationQuery) (err error) {
	msg.Result, err = h.inner.GetFbCustomerConversation(msg.GetArgs(ctx))
	return err
}

type GetFbCustomerConversationByIDQuery struct {
	ID dot.ID

	Result *FbCustomerConversation `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbCustomerConversationByID(ctx context.Context, msg *GetFbCustomerConversationByIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbCustomerConversationByID(msg.GetArgs(ctx))
	return err
}

type GetFbExternalCommentByIDQuery struct {
	ID dot.ID

	Result *FbExternalComment `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbExternalCommentByID(ctx context.Context, msg *GetFbExternalCommentByIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbExternalCommentByID(msg.GetArgs(ctx))
	return err
}

type GetFbExternalConversationByExternalIDAndExternalPageIDQuery struct {
	ExternalID     string
	ExternalPageID string

	Result *FbExternalConversation `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbExternalConversationByExternalIDAndExternalPageID(ctx context.Context, msg *GetFbExternalConversationByExternalIDAndExternalPageIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbExternalConversationByExternalIDAndExternalPageID(msg.GetArgs(ctx))
	return err
}

type GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery struct {
	ExternalPageID string
	ExternalUserID string

	Result *FbExternalConversation `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbExternalConversationByExternalPageIDAndExternalUserID(ctx context.Context, msg *GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbExternalConversationByExternalPageIDAndExternalUserID(msg.GetArgs(ctx))
	return err
}

type GetFbExternalConversationByIDQuery struct {
	ID dot.ID

	Result *FbExternalConversation `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbExternalConversationByID(ctx context.Context, msg *GetFbExternalConversationByIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbExternalConversationByID(msg.GetArgs(ctx))
	return err
}

type GetFbExternalMessageByIDQuery struct {
	ID dot.ID

	Result *FbExternalMessage `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbExternalMessageByID(ctx context.Context, msg *GetFbExternalMessageByIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbExternalMessageByID(msg.GetArgs(ctx))
	return err
}

type GetFbExternalPostByExternalIDQuery struct {
	ExternalID string

	Result *FbExternalPost `json:"-"`
}

func (h QueryServiceHandler) HandleGetFbExternalPostByExternalID(ctx context.Context, msg *GetFbExternalPostByExternalIDQuery) (err error) {
	msg.Result, err = h.inner.GetFbExternalPostByExternalID(msg.GetArgs(ctx))
	return err
}

type GetLatestCustomerExternalCommentQuery struct {
	ExternalPostID string
	ExternalUserID string

	Result *FbExternalComment `json:"-"`
}

func (h QueryServiceHandler) HandleGetLatestCustomerExternalComment(ctx context.Context, msg *GetLatestCustomerExternalCommentQuery) (err error) {
	msg.Result, err = h.inner.GetLatestCustomerExternalComment(msg.GetArgs(ctx))
	return err
}

type GetLatestFbExternalCommentQuery struct {
	ExternalPageID string
	ExternalPostID string
	ExternalUserID string

	Result *FbExternalComment `json:"-"`
}

func (h QueryServiceHandler) HandleGetLatestFbExternalComment(ctx context.Context, msg *GetLatestFbExternalCommentQuery) (err error) {
	msg.Result, err = h.inner.GetLatestFbExternalComment(msg.GetArgs(ctx))
	return err
}

type ListFbCustomerConversationsQuery struct {
	ExternalPageIDs []string
	ExternalUserID  dot.NullString
	IsRead          dot.NullBool
	Type            fb_customer_conversation_type.NullFbCustomerConversationType
	Paging          meta.Paging

	Result *FbCustomerConversationsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListFbCustomerConversations(ctx context.Context, msg *ListFbCustomerConversationsQuery) (err error) {
	msg.Result, err = h.inner.ListFbCustomerConversations(msg.GetArgs(ctx))
	return err
}

type ListFbCustomerConversationsByExternalIDsQuery struct {
	ExternalIDs filter.Strings

	Result []*FbCustomerConversation `json:"-"`
}

func (h QueryServiceHandler) HandleListFbCustomerConversationsByExternalIDs(ctx context.Context, msg *ListFbCustomerConversationsByExternalIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFbCustomerConversationsByExternalIDs(msg.GetArgs(ctx))
	return err
}

type ListFbExternalCommentsQuery struct {
	FbExternalPostID string
	FbExternalUserID string
	FbExternalPageID string
	Paging           meta.Paging

	Result *FbExternalCommentsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalComments(ctx context.Context, msg *ListFbExternalCommentsQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalComments(msg.GetArgs(ctx))
	return err
}

type ListFbExternalCommentsByExternalIDsQuery struct {
	FbExternalPostID string
	FbExternalUserID string
	FbExternalPageID string
	ExternalIDs      []string
	Paging           meta.Paging

	Result *FbExternalCommentsResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalCommentsByExternalIDs(ctx context.Context, msg *ListFbExternalCommentsByExternalIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalCommentsByExternalIDs(msg.GetArgs(ctx))
	return err
}

type ListFbExternalConversationsByExternalIDsQuery struct {
	ExternalIDs filter.Strings

	Result []*FbExternalConversation `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalConversationsByExternalIDs(ctx context.Context, msg *ListFbExternalConversationsByExternalIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalConversationsByExternalIDs(msg.GetArgs(ctx))
	return err
}

type ListFbExternalMessagesQuery struct {
	ExternalPageIDs         []string
	ExternalConversationIDs []string
	Paging                  meta.Paging

	Result *FbExternalMessagesResponse `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalMessages(ctx context.Context, msg *ListFbExternalMessagesQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalMessages(msg.GetArgs(ctx))
	return err
}

type ListFbExternalMessagesByExternalIDsQuery struct {
	ExternalIDs filter.Strings

	Result []*FbExternalMessage `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalMessagesByExternalIDs(ctx context.Context, msg *ListFbExternalMessagesByExternalIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalMessagesByExternalIDs(msg.GetArgs(ctx))
	return err
}

type ListFbExternalPostsByExternalIDsQuery struct {
	ExternalIDs filter.Strings

	Result []*FbExternalPost `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalPostsByExternalIDs(ctx context.Context, msg *ListFbExternalPostsByExternalIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalPostsByExternalIDs(msg.GetArgs(ctx))
	return err
}

type ListFbExternalPostsByIDsQuery struct {
	IDs filter.IDs

	Result []*FbExternalPost `json:"-"`
}

func (h QueryServiceHandler) HandleListFbExternalPostsByIDs(ctx context.Context, msg *ListFbExternalPostsByIDsQuery) (err error) {
	msg.Result, err = h.inner.ListFbExternalPostsByIDs(msg.GetArgs(ctx))
	return err
}

type ListLatestFbExternalMessagesQuery struct {
	ExternalConversationIDs filter.Strings

	Result []*FbExternalMessage `json:"-"`
}

func (h QueryServiceHandler) HandleListLatestFbExternalMessages(ctx context.Context, msg *ListLatestFbExternalMessagesQuery) (err error) {
	msg.Result, err = h.inner.ListLatestFbExternalMessages(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateFbCustomerConversationsCommand) command()         {}
func (q *CreateFbExternalConversationsCommand) command()         {}
func (q *CreateFbExternalMessagesCommand) command()              {}
func (q *CreateFbExternalPostsCommand) command()                 {}
func (q *CreateOrUpdateFbCustomerConversationsCommand) command() {}
func (q *CreateOrUpdateFbExternalCommentsCommand) command()      {}
func (q *CreateOrUpdateFbExternalConversationsCommand) command() {}
func (q *CreateOrUpdateFbExternalMessagesCommand) command()      {}
func (q *CreateOrUpdateFbExternalPostsCommand) command()         {}
func (q *UpdateIsReadCustomerConversationCommand) command()      {}

func (q *GetFbCustomerConversationQuery) query()                                  {}
func (q *GetFbCustomerConversationByIDQuery) query()                              {}
func (q *GetFbExternalCommentByIDQuery) query()                                   {}
func (q *GetFbExternalConversationByExternalIDAndExternalPageIDQuery) query()     {}
func (q *GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery) query() {}
func (q *GetFbExternalConversationByIDQuery) query()                              {}
func (q *GetFbExternalMessageByIDQuery) query()                                   {}
func (q *GetFbExternalPostByExternalIDQuery) query()                              {}
func (q *GetLatestCustomerExternalCommentQuery) query()                           {}
func (q *GetLatestFbExternalCommentQuery) query()                                 {}
func (q *ListFbCustomerConversationsQuery) query()                                {}
func (q *ListFbCustomerConversationsByExternalIDsQuery) query()                   {}
func (q *ListFbExternalCommentsQuery) query()                                     {}
func (q *ListFbExternalCommentsByExternalIDsQuery) query()                        {}
func (q *ListFbExternalConversationsByExternalIDsQuery) query()                   {}
func (q *ListFbExternalMessagesQuery) query()                                     {}
func (q *ListFbExternalMessagesByExternalIDsQuery) query()                        {}
func (q *ListFbExternalPostsByExternalIDsQuery) query()                           {}
func (q *ListFbExternalPostsByIDsQuery) query()                                   {}
func (q *ListLatestFbExternalMessagesQuery) query()                               {}

// implement conversion

func (q *CreateFbCustomerConversationsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateFbCustomerConversationsArgs) {
	return ctx,
		&CreateFbCustomerConversationsArgs{
			FbCustomerConversations: q.FbCustomerConversations,
		}
}

func (q *CreateFbCustomerConversationsCommand) SetCreateFbCustomerConversationsArgs(args *CreateFbCustomerConversationsArgs) {
	q.FbCustomerConversations = args.FbCustomerConversations
}

func (q *CreateFbExternalConversationsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateFbExternalConversationsArgs) {
	return ctx,
		&CreateFbExternalConversationsArgs{
			FbExternalConversations: q.FbExternalConversations,
		}
}

func (q *CreateFbExternalConversationsCommand) SetCreateFbExternalConversationsArgs(args *CreateFbExternalConversationsArgs) {
	q.FbExternalConversations = args.FbExternalConversations
}

func (q *CreateFbExternalMessagesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateFbExternalMessagesArgs) {
	return ctx,
		&CreateFbExternalMessagesArgs{
			FbExternalMessages: q.FbExternalMessages,
		}
}

func (q *CreateFbExternalMessagesCommand) SetCreateFbExternalMessagesArgs(args *CreateFbExternalMessagesArgs) {
	q.FbExternalMessages = args.FbExternalMessages
}

func (q *CreateFbExternalPostsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateFbExternalPostsArgs) {
	return ctx,
		&CreateFbExternalPostsArgs{
			FbExternalPosts: q.FbExternalPosts,
		}
}

func (q *CreateFbExternalPostsCommand) SetCreateFbExternalPostsArgs(args *CreateFbExternalPostsArgs) {
	q.FbExternalPosts = args.FbExternalPosts
}

func (q *CreateOrUpdateFbCustomerConversationsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateFbCustomerConversationsArgs) {
	return ctx,
		&CreateOrUpdateFbCustomerConversationsArgs{
			FbCustomerConversations: q.FbCustomerConversations,
		}
}

func (q *CreateOrUpdateFbCustomerConversationsCommand) SetCreateOrUpdateFbCustomerConversationsArgs(args *CreateOrUpdateFbCustomerConversationsArgs) {
	q.FbCustomerConversations = args.FbCustomerConversations
}

func (q *CreateOrUpdateFbExternalCommentsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateFbExternalCommentsArgs) {
	return ctx,
		&CreateOrUpdateFbExternalCommentsArgs{
			FbExternalComments: q.FbExternalComments,
		}
}

func (q *CreateOrUpdateFbExternalCommentsCommand) SetCreateOrUpdateFbExternalCommentsArgs(args *CreateOrUpdateFbExternalCommentsArgs) {
	q.FbExternalComments = args.FbExternalComments
}

func (q *CreateOrUpdateFbExternalConversationsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateFbExternalConversationsArgs) {
	return ctx,
		&CreateOrUpdateFbExternalConversationsArgs{
			FbExternalConversations: q.FbExternalConversations,
		}
}

func (q *CreateOrUpdateFbExternalConversationsCommand) SetCreateOrUpdateFbExternalConversationsArgs(args *CreateOrUpdateFbExternalConversationsArgs) {
	q.FbExternalConversations = args.FbExternalConversations
}

func (q *CreateOrUpdateFbExternalMessagesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateFbExternalMessagesArgs) {
	return ctx,
		&CreateOrUpdateFbExternalMessagesArgs{
			FbExternalMessages: q.FbExternalMessages,
		}
}

func (q *CreateOrUpdateFbExternalMessagesCommand) SetCreateOrUpdateFbExternalMessagesArgs(args *CreateOrUpdateFbExternalMessagesArgs) {
	q.FbExternalMessages = args.FbExternalMessages
}

func (q *CreateOrUpdateFbExternalPostsCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateOrUpdateFbExternalPostsArgs) {
	return ctx,
		&CreateOrUpdateFbExternalPostsArgs{
			FbExternalPosts: q.FbExternalPosts,
		}
}

func (q *CreateOrUpdateFbExternalPostsCommand) SetCreateOrUpdateFbExternalPostsArgs(args *CreateOrUpdateFbExternalPostsArgs) {
	q.FbExternalPosts = args.FbExternalPosts
}

func (q *UpdateIsReadCustomerConversationCommand) GetArgs(ctx context.Context) (_ context.Context, conversationCustomerID dot.ID, isRead bool) {
	return ctx,
		q.ConversationCustomerID,
		q.IsRead
}

func (q *GetFbCustomerConversationQuery) GetArgs(ctx context.Context) (_ context.Context, customerConversationType fb_customer_conversation_type.FbCustomerConversationType, externalID string, externalUserID string) {
	return ctx,
		q.CustomerConversationType,
		q.ExternalID,
		q.ExternalUserID
}

func (q *GetFbCustomerConversationByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetFbExternalCommentByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetFbExternalConversationByExternalIDAndExternalPageIDQuery) GetArgs(ctx context.Context) (_ context.Context, externalID string, externalPageID string) {
	return ctx,
		q.ExternalID,
		q.ExternalPageID
}

func (q *GetFbExternalConversationByExternalPageIDAndExternalUserIDQuery) GetArgs(ctx context.Context) (_ context.Context, externalPageID string, externalUserID string) {
	return ctx,
		q.ExternalPageID,
		q.ExternalUserID
}

func (q *GetFbExternalConversationByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetFbExternalMessageByIDQuery) GetArgs(ctx context.Context) (_ context.Context, ID dot.ID) {
	return ctx,
		q.ID
}

func (q *GetFbExternalPostByExternalIDQuery) GetArgs(ctx context.Context) (_ context.Context, externalID string) {
	return ctx,
		q.ExternalID
}

func (q *GetLatestCustomerExternalCommentQuery) GetArgs(ctx context.Context) (_ context.Context, externalPostID string, externalUserID string) {
	return ctx,
		q.ExternalPostID,
		q.ExternalUserID
}

func (q *GetLatestFbExternalCommentQuery) GetArgs(ctx context.Context) (_ context.Context, externalPageID string, externalPostID string, externalUserID string) {
	return ctx,
		q.ExternalPageID,
		q.ExternalPostID,
		q.ExternalUserID
}

func (q *ListFbCustomerConversationsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListFbCustomerConversationsArgs) {
	return ctx,
		&ListFbCustomerConversationsArgs{
			ExternalPageIDs: q.ExternalPageIDs,
			ExternalUserID:  q.ExternalUserID,
			IsRead:          q.IsRead,
			Type:            q.Type,
			Paging:          q.Paging,
		}
}

func (q *ListFbCustomerConversationsQuery) SetListFbCustomerConversationsArgs(args *ListFbCustomerConversationsArgs) {
	q.ExternalPageIDs = args.ExternalPageIDs
	q.ExternalUserID = args.ExternalUserID
	q.IsRead = args.IsRead
	q.Type = args.Type
	q.Paging = args.Paging
}

func (q *ListFbCustomerConversationsByExternalIDsQuery) GetArgs(ctx context.Context) (_ context.Context, externalIDs filter.Strings) {
	return ctx,
		q.ExternalIDs
}

func (q *ListFbExternalCommentsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListFbExternalCommentsArgs) {
	return ctx,
		&ListFbExternalCommentsArgs{
			FbExternalPostID: q.FbExternalPostID,
			FbExternalUserID: q.FbExternalUserID,
			FbExternalPageID: q.FbExternalPageID,
			Paging:           q.Paging,
		}
}

func (q *ListFbExternalCommentsQuery) SetListFbExternalCommentsArgs(args *ListFbExternalCommentsArgs) {
	q.FbExternalPostID = args.FbExternalPostID
	q.FbExternalUserID = args.FbExternalUserID
	q.FbExternalPageID = args.FbExternalPageID
	q.Paging = args.Paging
}

func (q *ListFbExternalCommentsByExternalIDsQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListFbExternalCommentsByIDsArgs) {
	return ctx,
		&ListFbExternalCommentsByIDsArgs{
			FbExternalPostID: q.FbExternalPostID,
			FbExternalUserID: q.FbExternalUserID,
			FbExternalPageID: q.FbExternalPageID,
			ExternalIDs:      q.ExternalIDs,
			Paging:           q.Paging,
		}
}

func (q *ListFbExternalCommentsByExternalIDsQuery) SetListFbExternalCommentsByIDsArgs(args *ListFbExternalCommentsByIDsArgs) {
	q.FbExternalPostID = args.FbExternalPostID
	q.FbExternalUserID = args.FbExternalUserID
	q.FbExternalPageID = args.FbExternalPageID
	q.ExternalIDs = args.ExternalIDs
	q.Paging = args.Paging
}

func (q *ListFbExternalConversationsByExternalIDsQuery) GetArgs(ctx context.Context) (_ context.Context, externalIDs filter.Strings) {
	return ctx,
		q.ExternalIDs
}

func (q *ListFbExternalMessagesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *ListFbExternalMessagesArgs) {
	return ctx,
		&ListFbExternalMessagesArgs{
			ExternalPageIDs:         q.ExternalPageIDs,
			ExternalConversationIDs: q.ExternalConversationIDs,
			Paging:                  q.Paging,
		}
}

func (q *ListFbExternalMessagesQuery) SetListFbExternalMessagesArgs(args *ListFbExternalMessagesArgs) {
	q.ExternalPageIDs = args.ExternalPageIDs
	q.ExternalConversationIDs = args.ExternalConversationIDs
	q.Paging = args.Paging
}

func (q *ListFbExternalMessagesByExternalIDsQuery) GetArgs(ctx context.Context) (_ context.Context, externalIDs filter.Strings) {
	return ctx,
		q.ExternalIDs
}

func (q *ListFbExternalPostsByExternalIDsQuery) GetArgs(ctx context.Context) (_ context.Context, externalIDs filter.Strings) {
	return ctx,
		q.ExternalIDs
}

func (q *ListFbExternalPostsByIDsQuery) GetArgs(ctx context.Context) (_ context.Context, IDs filter.IDs) {
	return ctx,
		q.IDs
}

func (q *ListLatestFbExternalMessagesQuery) GetArgs(ctx context.Context) (_ context.Context, externalConversationIDs filter.Strings) {
	return ctx,
		q.ExternalConversationIDs
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
	b.AddHandler(h.HandleCreateFbCustomerConversations)
	b.AddHandler(h.HandleCreateFbExternalConversations)
	b.AddHandler(h.HandleCreateFbExternalMessages)
	b.AddHandler(h.HandleCreateFbExternalPosts)
	b.AddHandler(h.HandleCreateOrUpdateFbCustomerConversations)
	b.AddHandler(h.HandleCreateOrUpdateFbExternalComments)
	b.AddHandler(h.HandleCreateOrUpdateFbExternalConversations)
	b.AddHandler(h.HandleCreateOrUpdateFbExternalMessages)
	b.AddHandler(h.HandleCreateOrUpdateFbExternalPosts)
	b.AddHandler(h.HandleUpdateIsReadCustomerConversation)
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
	b.AddHandler(h.HandleGetFbCustomerConversation)
	b.AddHandler(h.HandleGetFbCustomerConversationByID)
	b.AddHandler(h.HandleGetFbExternalCommentByID)
	b.AddHandler(h.HandleGetFbExternalConversationByExternalIDAndExternalPageID)
	b.AddHandler(h.HandleGetFbExternalConversationByExternalPageIDAndExternalUserID)
	b.AddHandler(h.HandleGetFbExternalConversationByID)
	b.AddHandler(h.HandleGetFbExternalMessageByID)
	b.AddHandler(h.HandleGetFbExternalPostByExternalID)
	b.AddHandler(h.HandleGetLatestCustomerExternalComment)
	b.AddHandler(h.HandleGetLatestFbExternalComment)
	b.AddHandler(h.HandleListFbCustomerConversations)
	b.AddHandler(h.HandleListFbCustomerConversationsByExternalIDs)
	b.AddHandler(h.HandleListFbExternalComments)
	b.AddHandler(h.HandleListFbExternalCommentsByExternalIDs)
	b.AddHandler(h.HandleListFbExternalConversationsByExternalIDs)
	b.AddHandler(h.HandleListFbExternalMessages)
	b.AddHandler(h.HandleListFbExternalMessagesByExternalIDs)
	b.AddHandler(h.HandleListFbExternalPostsByExternalIDs)
	b.AddHandler(h.HandleListFbExternalPostsByIDs)
	b.AddHandler(h.HandleListLatestFbExternalMessages)
	return QueryBus{b}
}
