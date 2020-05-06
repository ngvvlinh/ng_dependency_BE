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
	CreateFbExternalMessages(context.Context, CreateFbExternalMessagesArgs) ([]*FbExternalMessage, error)

	CreateFbExternalConversations(context.Context, CreateFbExternalConversationsArgs) ([]*FbExternalConversation, error)

	CreateFbCustomerConversations(context.Context, CreateFbCustomerConversationsArgs) ([]*FbCustomerConversation, error)
}

type QueryService interface {
	ListFbExternalMessagesByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalMessage, error)
	ListFbExternalMessages(context.Context, *ListFbExternalMessagesArgs) (*FbExternalMessagesResponse, error)
	ListLatestFbExternalMessages(_ context.Context, externalConversationIDs filter.Strings) ([]*FbExternalMessage, error)

	ListFbExternalConversationsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbExternalConversation, error)

	ListFbCustomerConversationsByExternalIDs(_ context.Context, externalIDs filter.Strings) ([]*FbCustomerConversation, error)
	ListFbCustomerConversations(context.Context, *ListFbCustomerConversationsArgs) (*FbCustomerConversationsResponse, error)
}

// +convert:create=FbExternalMessage
type CreateFbExternalMessageArgs struct {
	ID                     dot.ID
	FbConversationID       dot.ID
	ExternalConversationID string
	FbPageID               dot.ID
	ExternalID             string
	ExternalMessage        string
	ExternalTo             []*FbObjectTo
	ExternalFrom           *FbObjectFrom
	ExternalAttachments    []*FbMessageAttachment
	ExternalCreatedTime    time.Time
}

type CreateFbExternalMessagesArgs struct {
	FbExternalMessages []*CreateFbExternalMessageArgs
}

// +convert:create=FbExternalConversation
type CreateFbExternalConversationArgs struct {
	ID                   dot.ID
	FbPageID             dot.ID
	ExternalID           string
	ExternalUserID       string
	ExternalUserName     string
	ExternalLink         string
	ExternalUpdatedTime  time.Time
	ExternalMessageCount int
	LastMessage          string
	LastMessageAt        time.Time
}

type CreateFbExternalConversationsArgs struct {
	FbExternalConversations []*CreateFbExternalConversationArgs
}

// +convert:create=FbCustomerConversation
type CreateFbCustomerConversationArgs struct {
	ID               dot.ID
	FbPageID         dot.ID
	ExternalID       string
	ExternalUserID   string
	ExternalUserName string
	IsRead           bool
	Type             fb_customer_conversation_type.FbCustomerConversationType
	PostAttachments  []*PostAttachment
	LastMessage      string
	LastMessageAt    time.Time
}

type CreateFbCustomerConversationsArgs struct {
	FbCustomerConversations []*CreateFbCustomerConversationArgs
}

type ListFbExternalMessagesArgs struct {
	FbPageIDs         []dot.ID
	FbConversationIDs []string

	Paging meta.Paging
}

type FbExternalMessagesResponse struct {
	FbExternalMessages []*FbExternalMessage
	Paging             meta.PageInfo
}

type ListFbCustomerConversationsArgs struct {
	FbPageIDs        []dot.ID
	FbExternalUserID dot.NullString
	IsRead           dot.NullBool
	Type             fb_customer_conversation_type.NullFbCustomerConversationType

	Paging meta.Paging
}

type FbCustomerConversationsResponse struct {
	FbCustomerConversations []*FbCustomerConversation
	Paging                  meta.PageInfo
}
