package fbmessaging

import (
	"context"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	faboRedis "o.o/backend/com/fabo/pkg/redis"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

var _ fbmessaging.QueryService = &FbMessagingQuery{}

type FbMessagingQuery struct {
	db                               *cmsql.Database
	fbExternalPostStore              sqlstore.FbExternalPostStoreFactory
	fbExternalCommentStore           sqlstore.FbExternalCommentStoreFactory
	fbExternalConversationStore      sqlstore.FbExternalConversationStoreFactory
	fbExternalConversationStateStore sqlstore.FbCustomerConversationStateStoreFactory
	fbExternalMessagesStore          sqlstore.FbExternalMessageStoreFactory
	fbCustomerConversationStore      sqlstore.FbCustomerConversationStoreFactory
	rd                               *faboRedis.FaboRedis
}

func NewFbMessagingQuery(database com.MainDB, faboRedis *faboRedis.FaboRedis) *FbMessagingQuery {
	return &FbMessagingQuery{
		db:                               database,
		fbExternalPostStore:              sqlstore.NewFbExternalPostStore(database),
		fbExternalCommentStore:           sqlstore.NewFbExternalCommentStore(database),
		fbExternalConversationStore:      sqlstore.NewFbExternalConversationStore(database),
		fbExternalConversationStateStore: sqlstore.NewFbCustomerConversationStateStore(database),
		fbExternalMessagesStore:          sqlstore.NewFbExternalMessageStore(database),
		fbCustomerConversationStore:      sqlstore.NewFbCustomerConversationStore(database),
		rd:                               faboRedis,
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
	fbCustomerConversations, err := query.WithPaging(args.Paging).ListFbCustomerConversations()

	fbCustomerConversations, err = q.mapStateFbCustomerConversations(ctx, fbCustomerConversations, args.IsRead)
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
	return q.mapStateFbCustomerConversations(ctx, fbExternalMessages, dot.NullBool{Valid: false})
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

func (q *FbMessagingQuery) ListLatestCustomerFbExternalMessages(
	ctx context.Context, externalConversationIDs filter.Strings,
) ([]*fbmessaging.FbExternalMessage, error) {
	fbExternalMessages, err := q.fbExternalMessagesStore(ctx).ListLatestCustomerExternalMessages(externalConversationIDs)
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
	err = q.mapPostParent(ctx, fbExternalPosts)
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
	err = q.mapPostParent(ctx, fbExternalPosts)
	if err != nil {
		return nil, err
	}
	return fbExternalPosts, nil
}

func (q *FbMessagingQuery) ListFbExternalPosts(
	ctx context.Context, args *fbmessaging.LitFbExternalPostsArgs,
) (*fbmessaging.FbExternalPostsResponse, error) {
	query := q.fbExternalPostStore(ctx).WithPaging(args.Paging)

	if args.ExternalStatusType.Valid {
		query = query.ExternalStatusType(args.ExternalStatusType.Enum)
	}
	if len(args.ExternalIDs) > 0 {
		query = query.ExternalIDs(args.ExternalIDs)
	}
	if len(args.ExternalPageIDs) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "external_page_ids must not be null")
	} else {
		query = query.ExternalPageIDs(args.ExternalPageIDs)
	}

	fbExternalPosts, err := query.ListFbExternalPosts()
	if err != nil {
		return nil, err
	}

	return &fbmessaging.FbExternalPostsResponse{
		FbExternalPosts: fbExternalPosts,
		Paging:          query.GetPaging(),
	}, nil
}

func (q *FbMessagingQuery) mapPostParent(ctx context.Context, posts []*fbmessaging.FbExternalPost) error {
	var postParentIDs []string
	mapPost := make(map[string]*fbmessaging.FbExternalPost)
	for _, post := range posts {
		postParentIDs = append(postParentIDs, post.ExternalParentID)
	}
	postParents, err := q.fbExternalPostStore(ctx).ExternalIDs(postParentIDs).ListFbExternalPosts()
	if err != nil {
		return err
	}
	for _, postParent := range postParents {
		mapPost[postParent.ExternalID] = postParent
	}
	for _, post := range posts {
		if mapPost[post.ExternalParentID] != nil {
			post.ExternalParent = mapPost[post.ExternalParentID]
		}
	}
	return nil
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
	ctx context.Context, externalPostID, externalUserID, externalPageID string,
) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).GetLatestCustomerExternalComment(externalPostID, externalUserID, externalPageID)
}

func (q *FbMessagingQuery) ListFbExternalComments(
	ctx context.Context, args *fbmessaging.ListFbExternalCommentsArgs,
) (_ *fbmessaging.FbExternalCommentsResponse, err error) {
	query := q.fbExternalCommentStore(ctx).WithPaging(args.Paging).
		ExternalPostID(args.FbExternalPostID)

	var fbExternalComments []*fbmessaging.FbExternalComment
	if args.FbExternalUserID != "" {
		if args.FbExternalPageID != args.FbExternalUserID {
			query = query.ExternalPageIDAndExternalUserID(args.FbExternalPageID, args.FbExternalUserID)
			fbExternalComments, err = query.ListFbExternalComments()
			if err != nil {
				return nil, err
			}
		} else {
			fbExternalComments, err = query.ListFbExternalCommentsOfPage(args.FbExternalPageID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		fbExternalComments, err = query.ListFbExternalComments()
		if err != nil {
			return nil, err
		}
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

	fbExternalPost, err := q.fbExternalPostStore(ctx).ExternalID(externalID).GetFbExternalPost()
	if err != nil {
		return nil, err
	}
	err = q.mapPostParent(ctx, []*fbmessaging.FbExternalPost{fbExternalPost})
	if err != nil {
		return nil, err
	}
	return fbExternalPost, nil
}

func (q *FbMessagingQuery) GetExternalPostByExternalIDWithExternalCreatedTime(
	ctx context.Context, externalID string, created time.Time,
) (*fbmessaging.FbExternalPost, error) {
	fbExternalPost, err := q.
		fbExternalPostStore(ctx).
		ExternalID(externalID).
		ExternalCreatedTime(created).
		GetFbExternalPost()
	if err != nil {
		return nil, err
	}
	return fbExternalPost, nil
}

func (q *FbMessagingQuery) GetFbExternalMessageByID(
	ctx context.Context, ID dot.ID,
) (*fbmessaging.FbExternalMessage, error) {
	return q.fbExternalMessagesStore(ctx).ID(ID).GetFbExternalMessage()
}

func (q *FbMessagingQuery) GetFbExternalMessageByExternalID(
	ctx context.Context, externalID string,
) (*fbmessaging.FbExternalMessage, error) {
	return q.fbExternalMessagesStore(ctx).ExternalID(externalID).GetFbExternalMessage()
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

func (q *FbMessagingQuery) GetFbCustomerConversationByID(
	ctx context.Context, ID dot.ID,
) (*fbmessaging.FbCustomerConversation, error) {
	fbCustomerConversation, err := q.fbCustomerConversationStore(ctx).ID(ID).GetFbCustomerConversation()
	if err != nil {
		return nil, err
	}
	fbCustomerConversationState, err := q.fbExternalConversationStateStore(ctx).ID(ID).GetFbCustomerConversationState()
	if err != nil {
		return nil, err
	}

	fbCustomerConversation.IsRead = fbCustomerConversationState.IsRead
	return fbCustomerConversation, nil
}

func (q *FbMessagingQuery) GetFbExternalCommentByExternalID(
	ctx context.Context, externalID string,
) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).ExternalID(externalID).GetFbExternalComment()
}

func (q *FbMessagingQuery) ListFbCustomerConversationStates(
	ctx context.Context, IDs []dot.ID,
) ([]*fbmessaging.FbCustomerConversationState, error) {
	return q.fbExternalConversationStateStore(ctx).IDs(IDs...).ListFbCustomerConversationStates()
}

func (q *FbMessagingQuery) GetLatestUpdateActiveComment(
	ctx context.Context, extPostID, extUserID string,
) (*fbmessaging.FbExternalComment, error) {
	return q.fbExternalCommentStore(ctx).
		ExternalUserID(extUserID).
		ExternalPostID(extPostID).
		GetLatestUpdatedActiveComment()
}

func (q *FbMessagingQuery) ListFbCustomerConversationsByExternalUserIDs(
	ctx context.Context, extUserIDs []string, conversationType fb_customer_conversation_type.NullFbCustomerConversationType,
) ([]*fbmessaging.FbCustomerConversation, error) {
	query := q.fbCustomerConversationStore(ctx).ExternalUserIDs(extUserIDs)
	if conversationType.Valid {
		query = query.Type(conversationType.Enum)
	}
	conversations, err := query.ListFbCustomerConversations()
	if err != nil {
		return nil, err
	}
	return q.mapStateFbCustomerConversations(ctx, conversations, dot.NullBool{Valid: false})
}

func (q *FbMessagingQuery) ListFbCustomerConversationsByIDs(
	ctx context.Context, ids []dot.ID,
) ([]*fbmessaging.FbCustomerConversation, error) {
	conversations, err := q.fbCustomerConversationStore(ctx).IDs(ids).ListFbCustomerConversations()
	if err != nil {
		return nil, err
	}
	return q.mapStateFbCustomerConversations(ctx, conversations, dot.NullBool{Valid: false})
}

func (q *FbMessagingQuery) ListFbCustomerConversationsByExtUserIDsAndExtIDs(
	ctx context.Context, extUserIDs, extIDs []string,
) ([]*fbmessaging.FbCustomerConversation, error) {
	conversations, err := q.
		fbCustomerConversationStore(ctx).
		ExternalIDs(extIDs).
		ExternalUserIDs(extUserIDs).
		ListFbCustomerConversations()
	if err != nil {
		return nil, err
	}
	return q.mapStateFbCustomerConversations(ctx, conversations, dot.NullBool{Valid: false})
}

func (q *FbMessagingQuery) mapStateFbCustomerConversations(
	ctx context.Context,
	fbCustomerConversations []*fbmessaging.FbCustomerConversation,
	isRead dot.NullBool,
) ([]*fbmessaging.FbCustomerConversation, error) {
	var fbCustomerConversationsResult []*fbmessaging.FbCustomerConversation
	var fbCustomerConversationIDs []dot.ID
	for _, fbCustomerConversation := range fbCustomerConversations {
		fbCustomerConversationIDs = append(fbCustomerConversationIDs, fbCustomerConversation.ID)
	}
	getFbCustomerConversationStateQuery := q.fbExternalConversationStateStore(ctx).IDs(fbCustomerConversationIDs...)
	if isRead.Valid {
		getFbCustomerConversationStateQuery = getFbCustomerConversationStateQuery.IsRead(isRead.Bool)
	}

	fbCustomerConversationStates, err := getFbCustomerConversationStateQuery.ListFbCustomerConversationStates()
	if err != nil {
		return nil, err
	}

	mapFbCustomerConversationStates := make(map[dot.ID]*fbmessaging.FbCustomerConversationState)
	for _, fbCustomerConversationState := range fbCustomerConversationStates {
		mapFbCustomerConversationStates[fbCustomerConversationState.ID] = fbCustomerConversationState
	}

	for _, fbCustomerConversation := range fbCustomerConversations {
		if fbCustomerConversationState, ok := mapFbCustomerConversationStates[fbCustomerConversation.ID]; ok {
			fbCustomerConversation.IsRead = fbCustomerConversationState.IsRead
			fbCustomerConversationsResult = append(fbCustomerConversationsResult, fbCustomerConversation)
		}
	}

	return fbCustomerConversationsResult, nil
}
