package fbmessaging

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
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

func NewFbMessagingQuery(database *cmsql.Database) *FbMessagingQuery {
	return &FbMessagingQuery{
		db:                          database,
		fbExternalPostStore:         sqlstore.NewFbExternalPostStore(database),
		fbExternalCommentStore:      sqlstore.NewFbExternalCommentStore(database),
		fbExternalConversationStore: sqlstore.NewFbExternalConversationStore(database),
		fbExternalMessagesStore:     sqlstore.NewFbExternalMessageStore(database),
		fbCustomerConversationStore: sqlstore.NewFbCustomerConversationStore(database),
	}
}

func (q *FbMessagingQuery) MessageBus() fbmessaging.QueryBus {
	b := bus.New()
	return fbmessaging.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbMessagingQuery) ListFbCustomerConversations(
	ctx context.Context, args *fbmessaging.ListFbCustomerConversationsArgs,
) (*fbmessaging.FbCustomerConversationsResponse, error) {
	query := q.fbCustomerConversationStore(ctx).FbPageIDs(args.FbPageIDs)
	if args.Type.Valid && args.Type.Enum != fb_customer_conversation_type.All {
		query = query.Type(args.Type.Enum)
	}
	if args.FbExternalUserID.Valid {
		query = query.FbExternalID(args.FbExternalUserID.String)
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
	query := q.fbExternalMessagesStore(ctx).FbPageIDs(args.FbPageIDs).WithPaging(args.Paging)
	if len(args.FbConversationIDs) > 0 {
		query = query.ExternalConversationIDs(args.FbConversationIDs)
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

func (q *FbMessagingQuery) GetFbExternalPostByExternalID(
	ctx context.Context, externalID string,
) (*fbmessaging.FbExternalPost, error) {
	return q.fbExternalPostStore(ctx).ExternalID(externalID).GetFbExternalPost()
}
