package fbmessaging

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

var _ fbmessaging.QueryService = &FbMessagingQuery{}

type FbMessagingQuery struct {
	db                          *cmsql.Database
	fbExternalPostStore         sqlstore.FbExternalPostStoreFactory
	fbExternalCommentStore      sqlstore.FbExternalCommentStoreFactory
	fbExternalConversationStore sqlstore.FbExternalConversationStoreFactory
	fbExternalMessagesStore     sqlstore.FbExternalMessageStoreFactory
	fbCustomerConversationStore sqlstore.FbCustomerConversationStoreFactory
}

func NewFbMessagingQuery(database com.MainDB) *FbMessagingQuery {
	return &FbMessagingQuery{
		db:                          database,
		fbExternalPostStore:         sqlstore.NewFbExternalPostStore(database),
		fbExternalCommentStore:      sqlstore.NewFbExternalCommentStore(database),
		fbExternalConversationStore: sqlstore.NewFbExternalConversationStore(database),
		fbExternalMessagesStore:     sqlstore.NewFbExternalMessageStore(database),
		fbCustomerConversationStore: sqlstore.NewFbCustomerConversationStore(database),
	}
}

func FbMessagingQueryMessageBus(q *FbMessagingQuery) fbmessaging.QueryBus {
	b := bus.New()
	return fbmessaging.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbMessagingQuery) ListFbCustomerConversations(
	ctx context.Context, args *fbmessaging.ListFbCustomerConversationsArgs,
) (*fbmessaging.FbCustomerConversationsResponse, error) {
	query := q.fbCustomerConversationStore(ctx)
	if len(args.ExternalPageIDs) != 0 {
		query = query.ExternalPageIDs(args.ExternalPageIDs)
	}
	if args.Type.Valid && args.Type.Enum != fb_customer_conversation_type.All {
		query = query.Type(args.Type.Enum)
	}
	if args.ExternalUserID.Valid {
		query = query.FbExternalID(args.ExternalUserID.String)
	}
	if args.IsRead.Valid {
		query = query.IsRead(args.IsRead.Bool)
	}
	fbCustomerConversations, err := query.WithPaging(args.Paging).ListFbCustomerConversations()
	if err != nil {
		return nil, err
	}
	return &fbmessaging.FbCustomerConversationsResponse{
		FbCustomerConversations: fbCustomerConversations,
		Paging:                  query.GetPaging(),
	}, nil
}

func (q *FbMessagingQuery) GetFbExternalConversationByExternalPageIDAndExternalUserID(
	ctx context.Context, externalPageID, externalUserID string,
) (*fbmessaging.FbExternalConversation, error) {
	return q.fbExternalConversationStore(ctx).ExternalPageID(externalPageID).ExternalUserID(externalUserID).GetFbExternalConversation()
}

func (q *FbMessagingQuery) GetFbExternalConversationByExternalIDAndExternalPageID(
	ctx context.Context, externalID, externalPageID string,
) (*fbmessaging.FbExternalConversation, error) {
	return q.fbExternalConversationStore(ctx).ExternalID(externalID).ExternalPageID(externalPageID).GetFbExternalConversation()
}

func (q *FbMessagingQuery) ListFbExternalConversationsByExternalIDs(
	ctx context.Context, externalIDs filter.Strings,
) ([]*fbmessaging.FbExternalConversation, error) {
	query := q.fbExternalConversationStore(ctx).ExternalIDs(externalIDs)
	fbExternalConversations, err := query.ListFbExternalConversations()
	if err != nil {
		return nil, err
	}
	return fbExternalConversations, nil
}

func (q *FbMessagingQuery) ListFbExternalMessagesByExternalIDs(
	ctx context.Context, externalIDs filter.Strings,
) ([]*fbmessaging.FbExternalMessage, error) {
	query := q.fbExternalMessagesStore(ctx).ExternalIDs(externalIDs)
	fbExternalMessages, err := query.ListFbExternalMessages()
	if err != nil {
		return nil, err
	}
	return fbExternalMessages, nil
}

func (q *FbMessagingQuery) ListFbCustomerConversationsByExternalIDs(
	ctx context.Context, externalIDs filter.Strings,
) ([]*fbmessaging.FbCustomerConversation, error) {
	query := q.fbCustomerConversationStore(ctx).ExternalIDs(externalIDs)
	fbExternalMessages, err := query.ListFbCustomerConversations()
	if err != nil {
		return nil, err
	}
	return fbExternalMessages, nil
}

func (q *FbMessagingQuery) ListFbExternalMessages(
	ctx context.Context, args *fbmessaging.ListFbExternalMessagesArgs,
) (*fbmessaging.FbExternalMessagesResponse, error) {
	query := q.fbExternalMessagesStore(ctx).ExternalPageIDs(args.ExternalPageIDs).WithPaging(args.Paging)
	if len(args.ExternalConversationIDs) > 0 {
		query = query.ExternalConversationIDs(args.ExternalConversationIDs)
	}
	fbExternalMessages, err := query.ListFbExternalMessages()
	if err != nil {
		return nil, err
	}
	return &fbmessaging.FbExternalMessagesResponse{
		FbExternalMessages: fbExternalMessages,
		Paging:             query.GetPaging(),
	}, nil
}

func (q *FbMessagingQuery) ListLatestFbExternalMessages(
	ctx context.Context, externalConversationIDs filter.Strings,
) ([]*fbmessaging.FbExternalMessage, error) {
	fbExternalMessages, err := q.fbExternalMessagesStore(ctx).ListLatestExternalMessages(externalConversationIDs)
	if err != nil {
		return nil, err
	}
	return fbExternalMessages, nil
}

func (q *FbMessagingQuery) ListFbExternalPostsByExternalIDs(
	ctx context.Context, externalIDs filter.Strings,
) ([]*fbmessaging.FbExternalPost, error) {
	fbExternalPosts, err := q.fbExternalPostStore(ctx).ExternalIDs(externalIDs).ListFbExternalPosts()
	if err != nil {
		return nil, err
	}
	return fbExternalPosts, nil
}

func (q *FbMessagingQuery) ListFbExternalPostsByIDs(
	ctx context.Context, IDs filter.IDs,
) ([]*fbmessaging.FbExternalPost, error) {
	fbExternalPosts, err := q.fbExternalPostStore(ctx).IDs(IDs).ListFbExternalPosts()
	if err != nil {
		return nil, err
	}
	return fbExternalPosts, nil
}

func (q *FbMessagingQuery) GetFbCustomerConversation(
	ctx context.Context, customerConversationType fb_customer_conversation_type.FbCustomerConversationType,
	externalID, externalUserID string,
) (*fbmessaging.FbCustomerConversation, error) {
	if customerConversationType == fb_customer_conversation_type.Comment {
		return q.fbCustomerConversationStore(ctx).
			FbExternalID(externalID).FbExternalUserID(externalUserID).GetFbCustomerConversation()
	}
	return q.fbCustomerConversationStore(ctx).
		FbExternalID(externalID).GetFbCustomerConversation()
}

func (q *FbMessagingQuery) GetLatestFbExternalComment(
	ctx context.Context, externalPageID, externalPostID, externalUserID string,
) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).GetLatestExternalComment(externalPageID, externalPostID, externalUserID)
}

func (q *FbMessagingQuery) GetLatestCustomerExternalComment(
	ctx context.Context, externalPostID, externalUserID string,
) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).GetLatestCustomerExternalComment(externalPostID, externalUserID)
}

func (q *FbMessagingQuery) ListFbExternalComments(
	ctx context.Context, args *fbmessaging.ListFbExternalCommentsArgs,
) (*fbmessaging.FbExternalCommentsResponse, error) {
	query := q.fbExternalCommentStore(ctx).WithPaging(args.Paging).
		ExternalPostID(args.FbExternalPostID).ExternalPageIDAndExternalUserID(args.FbExternalPageID, args.FbExternalUserID)
	fbExternalComments, err := query.ListFbExternalComments()
	if err != nil {
		return nil, err
	}
	return &fbmessaging.FbExternalCommentsResponse{
		FbExternalComments: fbExternalComments,
		Paging:             query.GetPaging(),
	}, nil
}

func (q *FbMessagingQuery) ListFbExternalCommentsByExternalIDs(
	ctx context.Context, args *fbmessaging.ListFbExternalCommentsByIDsArgs,
) (*fbmessaging.FbExternalCommentsResponse, error) {
	query := q.fbExternalCommentStore(ctx).WithPaging(args.Paging).ExternalIDs(args.ExternalIDs)
	fbExternalComments, err := query.ListFbExternalComments()
	if err != nil {
		return nil, err
	}
	return &fbmessaging.FbExternalCommentsResponse{
		FbExternalComments: fbExternalComments,
		Paging:             query.GetPaging(),
	}, nil
}

func (q *FbMessagingQuery) GetFbExternalPostByExternalID(
	ctx context.Context, externalID string,
) (*fbmessaging.FbExternalPost, error) {
	return q.fbExternalPostStore(ctx).ExternalID(externalID).GetFbExternalPost()
}

func (q *FbMessagingQuery) GetFbExternalMessageByID(
	ctx context.Context, ID dot.ID,
) (*fbmessaging.FbExternalMessage, error) {
	return q.fbExternalMessagesStore(ctx).ID(ID).GetFbExternalMessage()
}

func (q *FbMessagingQuery) GetFbExternalCommentByID(
	ctx context.Context, ID dot.ID,
) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).ID(ID).GetFbExternalComment()
}

func (q *FbMessagingQuery) GetFbExternalConversationByID(
	ctx context.Context, ID dot.ID,
) (*fbmessaging.FbExternalConversation, error) {
	return q.fbExternalConversationStore(ctx).ID(ID).GetFbExternalConversation()
}

func (q *FbMessagingQuery) GetFbCustomerConversationByID(ctx context.Context, ID dot.ID) (*fbmessaging.FbCustomerConversation, error) {
	return q.fbCustomerConversationStore(ctx).ID(ID).GetFbCustomerConversation()
}

func (q *FbMessagingQuery) GetFbExternalCommentByExternalID(ctx context.Context, externalID string) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).ExternalID(externalID).GetFbExternalComment()
}
