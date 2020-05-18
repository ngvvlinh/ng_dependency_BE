package fbmessaging

import (
	context "context"
	"time"

	"o.o/capi/filter"

	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/meta"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateFbExternalMessages(context.Context, *CreateFbExternalMessagesArgs) ([]*FbExternalMessage, error)
	CreateOrUpdateFbExternalMessages(context.Context, *CreateOrUpdateFbExternalMessagesArgs) ([]*FbExternalMessage, error)

	CreateFbExternalPosts(context.Context, *CreateFbExternalPostsArgs) ([]*FbExternalPost, error)
	CreateOrUpdateFbExternalPosts(context.Context, *CreateOrUpdateFbExternalPostsArgs) ([]*FbExternalPost, error)

	CreateOrUpdateFbExternalComments(context.Context, *CreateOrUpdateFbExternalCommentsArgs) ([]*FbExternalComment, error)

	CreateFbExternalConversations(context.Context, *CreateFbExternalConversationsArgs) ([]*FbExternalConversation, error)
	CreateOrUpdateFbExternalConversations(context.Context, *CreateOrUpdateFbExternalConversationsArgs) ([]*FbExternalConversation, error)

	CreateFbCustomerConversations(context.Context, *CreateFbCustomerConversationsArgs) ([]*FbCustomerConversation, error)
	CreateOrUpdateFbCustomerConversations(context.Context, *CreateOrUpdateFbCustomerConversationsArgs) ([]*FbCustomerConversation, error)
	UpdateIsReadCustomerConversation(ctx context.Context, conversationCustomerID dot.ID, isRead bool) (int, error)
}

type QueryService interface {
	ListFbExternalMessagesByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalMessage, error)
	ListFbExternalMessages(context.Context, *ListFbExternalMessagesArgs) (*FbExternalMessagesResponse, error)
	ListLatestFbExternalMessages(_ context.Context, externalConversationIDs filter.Strings) ([]*FbExternalMessage, error)

	GetLatestFbExternalComment(_ context.Context, externalPageID, externalPostID, externalUserID string) (*FbExternalComment, error)
	GetLatestCustomerExternalComment(_ context.Context, externalPostID, externalUserID string) (*FbExternalComment, error)
	ListFbExternalComments(context.Context, *ListFbExternalCommentsArgs) (*FbExternalCommentsResponse, error)

	GetFbExternalPostByExternalID(_ context.Context, externalID string) (*FbExternalPost, error)
	ListFbExternalPostsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalPost, error)
	ListFbExternalPostsByIDs(_ context.Context, IDs filter.IDs) ([]*FbExternalPost, error)

	GetFbExternalConversationByExternalIDAndExternalPageID(_ context.Context, externalID, externalPageID string) (*FbExternalConversation, error)
	GetFbExternalConversationByExternalPageIDAndExternalUserID(_ context.Context, externalPageID, externalUserID string) (*FbExternalConversation, error)
	ListFbExternalConversationsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalConversation, error)

	GetFbCustomerConversation(_ context.Context, customerConversationType fb_customer_conversation_type.FbCustomerConversationType, externalID, externalUserID string) (*FbCustomerConversation, error)
	ListFbCustomerConversationsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbCustomerConversation, error)
	ListFbCustomerConversations(context.Context, *ListFbCustomerConversationsArgs) (*FbCustomerConversationsResponse, error)
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
	ExternalCreatedTime    time.Time
}

type CreateFbExternalMessagesArgs struct {
	FbExternalMessages []*CreateFbExternalMessageArgs
}

type CreateOrUpdateFbExternalMessagesArgs struct {
	FbExternalMessages []*CreateFbExternalMessageArgs
}

// +convert:create=FbExternalPost
type CreateFbExternalPostArgs struct {
	ID                  dot.ID
	ExternalPageID      string
	ExternalID          string
	ExternalParentID    string
	ExternalFrom        *FbObjectFrom
	ExternalPicture     string
	ExternalIcon        string
	ExternalMessage     string
	ExternalAttachments []*PostAttachment
	ExternalCreatedTime time.Time
	ExternalUpdatedTime time.Time
}

type CreateFbExternalPostsArgs struct {
	FbExternalPosts []*CreateFbExternalPostArgs
}

type CreateOrUpdateFbExternalPostsArgs struct {
	FbExternalPosts []*CreateFbExternalPostArgs
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
	Type            fb_customer_conversation_type.NullFbCustomerConversationType

	Paging meta.Paging
}

type FbCustomerConversationsResponse struct {
	FbCustomerConversations []*FbCustomerConversation
	Paging                  meta.PageInfo
}

type ListFbExternalCommentsArgs struct {
	FbExternalPostID string
	FbExternalUserID string
	FbExternalPageID string

	Paging meta.Paging
}

type FbExternalCommentsResponse struct {
	FbExternalComments []*FbExternalComment
	Paging             meta.PageInfo
}
