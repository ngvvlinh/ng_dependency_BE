package fabo

import (
	"context"

	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=fabo

// +apix:path=/fabo.Page
type PageService interface {
	ConnectPages(context.Context, *ConnectPagesRequest) (*ConnectPagesResponse, error)
	RemovePages(context.Context, *RemovePagesRequest) (*cm.Empty, error)
	ListPages(context.Context, *ListPagesRequest) (*ListPagesResponse, error)
	CheckPermissions(context.Context, *CheckPagePermissionsRequest) (*CheckPagePermissionsResponse, error)
}

// +apix:path=/fabo.CustomerConversation
type CustomerConversationService interface {
	ListCustomerConversations(context.Context, *ListCustomerConversationsRequest) (*FbCustomerConversationsResponse, error)

	GetCustomerConversationByID(context.Context, *GetCustomerConversationByIDRequest) (*GetCustomerConversationByIDResponse, error)

	ListMessages(context.Context, *ListMessagesRequest) (*FbMessagesResponse, error)

	ListCommentsByExternalPostID(context.Context, *ListCommentsByExternalPostIDRequest) (*ListCommentsByExternalPostIDResponse, error)
	UpdateReadStatus(context.Context, *UpdateReadStatusRequest) (*cm.UpdatedResponse, error)

	SendMessage(context.Context, *SendMessageRequest) (*FbExternalMessage, error)

	SendComment(context.Context, *SendCommentRequest) (*FbExternalComment, error)

	CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error)
}

// +apix:path=/fabo.Customer
type CustomerService interface {
	CreateFbUserCustomer(ctx context.Context, request *CreateFbUserCustomerRequest) (*FbUserWithCustomer, error)

	ListFbUsers(ctx context.Context, request *ListFbUsersRequest) (*ListFbUsersResponse, error)
	GetFbUser(ctx context.Context, request *GetFbUserRequest) (*FbUserWithCustomer, error)
	ListCustomersWithFbUsers(ctx context.Context, request *ListCustomersWithFbUsersRequest) (*ListCustomersWithFbUsersResponse, error)
}
