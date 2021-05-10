package fbmessaging

import (
	context "context"
	"time"

	"o.o/api/fabo/fbmessaging/fb_comment_source"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbmessaging/fb_feed_type"
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_live_video_status"
	"o.o/api/fabo/fbmessaging/fb_post_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/api/meta"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

// +gen:api

type Aggregate interface {
	CreateFbExternalMessagesFromSync(context.Context, *CreateFbExternalMessagesFromSyncArgs) ([]*FbExternalMessage, error)
	CreateOrUpdateFbExternalMessages(context.Context, *CreateOrUpdateFbExternalMessagesArgs) ([]*FbExternalMessage, error)

	CreateOrUpdateFbExternalComments(context.Context, *CreateOrUpdateFbExternalCommentsArgs) ([]*FbExternalComment, error)

	CreateOrUpdateFbExternalConversations(context.Context, *CreateOrUpdateFbExternalConversationsArgs) ([]*FbExternalConversation, error)
	CreateOrUpdateFbExternalConversation(context.Context, *FbExternalConversation) (*FbExternalConversation, error)
	CreateFbCustomerConversations(context.Context, *CreateFbCustomerConversationsArgs) ([]*FbCustomerConversation, error)
	CreateOrUpdateFbCustomerConversations(context.Context, *CreateOrUpdateFbCustomerConversationsArgs) ([]*FbCustomerConversation, error)
	UpdateIsReadCustomerConversation(ctx context.Context, conversationCustomerID dot.ID, isRead bool) (int, error)

	UpdateFbCommentMessage(context.Context, *FbUpdateCommentMessageArgs) (int, error)
	RemoveComment(context.Context, *RemoveCommentArgs) error
	LikeOrUnLikeComment(context.Context, *LikeOrUnLikeCommentArgs) error
	HideOrUnHideComment(context.Context, *HideOrUnHideCommentArgs) error
	UpdateIsPrivateRepliedComment(context.Context, *UpdateIsPrivateRepliedCommentArgs) error

	UpdateOrCreateFbExternalPostsFromSync(context.Context, *UpdateOrCreateFbExternalPostsFromSyncArgs) ([]*FbExternalPost, error)
	UpdateLiveVideoStatusFromSync(context.Context, *UpdateLiveVideoStatusFromSyncArgs) (*FbExternalPost, error)
	SaveFbExternalPost(context.Context, *FbSavePostArgs) (*FbExternalPost, error)
	CreateFbExternalPost(context.Context, *FbCreatePostArgs) (*FbExternalPost, error)
	UpdateFbPostMessageAndPicture(context.Context, *FbUpdatePostMessageArgs) error
	CreateFbExternalPosts(context.Context, *CreateFbExternalPostsArgs) ([]*FbExternalPost, error)
	RemovePost(context.Context, *RemovePostArgs) error
	UpdateFbExternalPostTotalComments(ctx context.Context, externalID string) error
}

type QueryService interface {
	ListFbExternalMessagesByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalMessage, error)
	ListFbExternalMessages(context.Context, *ListFbExternalMessagesArgs) (*FbExternalMessagesResponse, error)
	ListLatestFbExternalMessages(_ context.Context, externalConversationIDs filter.Strings) ([]*FbExternalMessage, error)
	ListLatestCustomerFbExternalMessages(_ context.Context, externalConversationIDs filter.Strings) ([]*FbExternalMessage, error)

	GetLatestFbExternalComment(_ context.Context, externalPageID, externalPostID, externalUserID string) (*FbExternalComment, error)
	GetLatestFbExternalUserComment(_ context.Context, externalOwnerPostID, externalPostID, externalUserID string) (*FbExternalComment, error)
	GetLatestCustomerExternalComment(_ context.Context, externalPostID, externalUserID, externalPageID string) (*FbExternalComment, error)
	ListFbExternalComments(context.Context, *ListFbExternalCommentsArgs) (*FbExternalCommentsResponse, error)
	ListFbExternalCommentsByExternalIDs(context.Context, *ListFbExternalCommentsByIDsArgs) (*FbExternalCommentsResponse, error)

	GetExternalPostByExternalIDWithExternalCreatedTime(_ context.Context, externalID string, time time.Time) (*FbExternalPost, error)
	GetFbExternalPostByExternalID(_ context.Context, externalID string) (*FbExternalPost, error)
	GetFbExternalPostByExternalIDAndExternalUserID(_ context.Context, externalID, externalUserID string) (*FbExternalPost, error)
	GetFbExternalMessageByID(_ context.Context, ID dot.ID) (*FbExternalMessage, error)
	GetFbExternalMessageByExternalID(_ context.Context, externalID string) (*FbExternalMessage, error)
	GetFbExternalCommentByID(_ context.Context, ID dot.ID) (*FbExternalComment, error)
	GetFbExternalCommentByExternalID(_ context.Context, externalID string) (*FbExternalComment, error)
	GetLatestUpdateActiveComment(_ context.Context, extPostID string, extUserID string) (*FbExternalComment, error)
	GetFbExternalConversationByID(_ context.Context, ID dot.ID) (*FbExternalConversation, error)
	GetFbExternalPostByID(_ context.Context, ID dot.ID) (*FbExternalPost, error)
	ListFbExternalPostsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalPost, error)
	ListFbExternalPostsByIDs(_ context.Context, IDs filter.IDs) ([]*FbExternalPost, error)
	ListFbExternalPosts(context.Context, *LitFbExternalPostsArgs) (*FbExternalPostsResponse, error)

	GetFbExternalConversationByExternalIDAndExternalPageID(_ context.Context, externalID, externalPageID string) (*FbExternalConversation, error)
	GetFbExternalConversationByExternalPageIDAndExternalUserID(_ context.Context, externalPageID, externalUserID string) (*FbExternalConversation, error)
	ListFbExternalConversationsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalConversation, error)

	GetFbCustomerConversation(_ context.Context, customerConversationType fb_customer_conversation_type.FbCustomerConversationType, externalID, externalUserID string) (*FbCustomerConversation, error)
	GetFbCustomerConversationByID(_ context.Context, ID dot.ID) (*FbCustomerConversation, error)
	ListFbCustomerConversationsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbCustomerConversation, error)
	ListFbCustomerConversations(context.Context, *ListFbCustomerConversationsArgs) (*FbCustomerConversationsResponse, error)
	ListFbCustomerConversationsByExternalUserIDs(_ context.Context, extUserIDs []string, conversationType fb_customer_conversation_type.NullFbCustomerConversationType) ([]*FbCustomerConversation, error)
	ListFbCustomerConversationsByIDs(_ context.Context, IDs []dot.ID) ([]*FbCustomerConversation, error)
	ListFbCustomerConversationsByExtUserIDsAndExtIDs(_ context.Context, extUserIDs []string, extIDs []string) ([]*FbCustomerConversation, error)

	ListFbCustomerConversationStates(_ context.Context, IDs []dot.ID) ([]*FbCustomerConversationState, error)
}

// +convert:create=FbExternalMessage
type CreateFbExternalMessageArgs struct {
	ID                     dot.ID
	ExternalConversationID string
	ExternalPageID         string
	ExternalID             string
	ExternalMessage        string
	ExternalSticker        string
	ExternalTo             []*FbObjectTo
	ExternalFrom           *FbObjectFrom
	ExternalAttachments    []*FbMessageAttachment
	ExternalMessageShares  []*FbMessageShare
	ExternalCreatedTime    time.Time
	ExternalTimestamp      int64
	InternalSource         fb_internal_source.FbInternalSource
	CreatedBy              dot.ID
}

type CreateFbExternalMessagesArgs struct {
	FbExternalMessages []*CreateFbExternalMessageArgs
}

type CreateOrUpdateFbExternalMessagesArgs struct {
	FbExternalMessages []*CreateFbExternalMessageArgs
}

type CreateFbExternalMessagesFromSyncArgs struct {
	FbExternalMessages []*FbExternalMessage
}

// +convert:create=FbExternalPost
type CreateFbExternalPostArgs struct {
	ID                      dot.ID
	ExternalPageID          string
	ExternalUserID          string
	ExternalID              string
	ExternalParentID        string
	ExternalFrom            *FbObjectFrom
	ExternalPicture         string
	ExternalIcon            string
	ExternalMessage         string
	ExternalAttachments     []*PostAttachment
	ExternalCreatedTime     time.Time
	ExternalUpdatedTime     time.Time
	TotalComments           int
	TotalReactions          int
	Type                    fb_post_type.FbPostType
	FeedType                fb_feed_type.FbFeedType
	StatusType              fb_status_type.FbStatusType
	ExternalLiveVideoStatus string
	LiveVideoStatus         fb_live_video_status.FbLiveVideoStatus
}

type CreateFbExternalPostsArgs struct {
	FbExternalPosts []*CreateFbExternalPostArgs
}

type UpdateOrCreateFbExternalPostsFromSyncArgs struct {
	FbExternalPosts []*CreateFbExternalPostArgs
}

// +convert:update=FbExternalPost
type UpdateLiveVideoStatusFromSyncArgs struct {
	ExternalID              string
	ExternalLiveVideoStatus string
	LiveVideoStatus         fb_live_video_status.FbLiveVideoStatus
}

// +convert:create=FbExternalComment
type CreateFbExternalCommentArgs struct {
	ID                   dot.ID
	ExternalPostID       string
	ExternalPageID       string
	ExternalID           string
	ExternalUserID       string
	ExternalParentID     string
	ExternalParentUserID string
	ExternalMessage      string
	ExternalCommentCount int
	ExternalParent       *FbObjectParent
	ExternalFrom         *FbObjectFrom
	ExternalAttachment   *CommentAttachment
	ExternalCreatedTime  time.Time
	Source               fb_comment_source.FbCommentSource
	InternalSource       fb_internal_source.FbInternalSource

	IsLiked          bool
	IsHidden         bool
	IsPrivateReplied bool
	CreatedBy        dot.ID

	ExternalOwnerPostID string
	PostType            fb_post_type.FbPostType
}

type CreateOrUpdateFbExternalCommentsArgs struct {
	FbExternalComments []*CreateFbExternalCommentArgs
}

// +convert:create=FbExternalConversation
type CreateFbExternalConversationArgs struct {
	ID                   dot.ID
	ExternalPageID       string
	ExternalID           string
	PSID                 string
	ExternalUserID       string
	ExternalUserName     string
	ExternalLink         string
	ExternalUpdatedTime  time.Time
	ExternalMessageCount int
}

type CreateFbExternalConversationsArgs struct {
	FbExternalConversations []*CreateFbExternalConversationArgs
}

type CreateOrUpdateFbExternalConversationsArgs struct {
	FbExternalConversations []*CreateFbExternalConversationArgs
}

// +convert:create=FbCustomerConversation
type CreateFbCustomerConversationArgs struct {
	ID                         dot.ID
	ExternalPageID             string
	ExternalOwnerPostID        string
	ExternalID                 string
	ExternalUserID             string
	ExternalUserName           string
	ExternalFrom               *FbObjectFrom
	IsRead                     bool
	Type                       fb_customer_conversation_type.FbCustomerConversationType
	ExternalPostAttachments    []*PostAttachment
	ExternalCommentAttachment  *CommentAttachment
	ExternalMessageAttachments []*FbMessageAttachment
	LastMessage                string
	LastMessageAt              time.Time
	LastCustomerMessageAt      time.Time
	LastMessageExternalID      string
}

type CreateFbCustomerConversationsArgs struct {
	FbCustomerConversations []*CreateFbCustomerConversationArgs
}

type CreateOrUpdateFbCustomerConversationsArgs struct {
	FbCustomerConversations []*CreateFbCustomerConversationArgs
}

type ListFbExternalMessagesArgs struct {
	ExternalPageIDs         []string
	ExternalConversationIDs []string

	Paging meta.Paging
}

type FbExternalMessagesResponse struct {
	FbExternalMessages []*FbExternalMessage
	Paging             meta.PageInfo
}

type ListFbCustomerConversationsArgs struct {
	ExternalPageIDs []string
	ExternalUserID  dot.NullString
	IsRead          dot.NullBool
	Types           []fb_customer_conversation_type.FbCustomerConversationType

	Paging meta.Paging
}

type FbCustomerConversationsResponse struct {
	FbCustomerConversations []*FbCustomerConversation
	Paging                  meta.PageInfo
}

type ListFbExternalCommentsArgs struct {
	FbExternalPostID      string
	FbExternalUserID      string
	FbExternalPageID      string
	FbExternalOwnerPostID string

	Paging meta.Paging
}

type ListFbExternalCommentsByIDsArgs struct {
	FbExternalPostID string
	FbExternalPageID string
	ExternalIDs      []string
	Paging           meta.Paging
}

type FbExternalCommentsResponse struct {
	FbExternalComments []*FbExternalComment
	Paging             meta.PageInfo
}

type LitFbExternalPostsArgs struct {
	IsLiveVideo        dot.NullBool
	ExternalUserID     string
	ExternalPageIDs    []string
	ExternalStatusType fb_status_type.NullFbStatusType
	LiveVideoStatus    fb_live_video_status.NullFbLiveVideoStatus
	ExternalIDs        []string
	Paging             meta.Paging
}

type FbExternalPostsResponse struct {
	FbExternalPosts []*FbExternalPost
	Paging          meta.PageInfo
}

type FbCreatePostArgs struct {
	ExternalPageID string
	AccessToken    string
	Message        string
}

type FbSavePostArgs struct {
	ExternalPageID      string
	ExternalID          string
	ExternalFrom        *FbObjectFrom
	ExternalPicture     string
	ExternalIcon        string
	ExternalMessage     string
	ExternalAttachments []*PostAttachment `compare:"ignore"`
	ExternalCreatedTime time.Time
	ExternalParentID    string
	FeedType            fb_feed_type.FbFeedType
	StatusType          fb_status_type.FbStatusType
	Type                fb_post_type.FbPostType
}

type FbUpdatePostMessageArgs struct {
	ExternalPostID  string
	Message         string
	ExternalPicture string
}

type FbUpdateCommentMessageArgs struct {
	ExternalCommentID string
	Message           string
}

type RemovePostArgs struct {
	ExternalPostID string
	ExternalPageID string
}

type RemoveCommentArgs struct {
	ExternalCommentID string
}

type LikeOrUnLikeCommentArgs struct {
	ExternalCommentID string
	IsLiked           bool
}

type HideOrUnHideCommentArgs struct {
	ExternalCommentID string
	IsHidden          bool
}

type UpdateIsPrivateRepliedCommentArgs struct {
	ExternalCommentID string
	IsPrivateReplied  bool
}
