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
}

// +apix:path=/fabo.CustomerConversation
type CustomerConversationService interface {
	ListCustomerConversations(context.Context, *ListCustomerConversationsRequest) (*FbCustomerConversationsResponse, error)
	ListMessages(context.Context, *ListMessagesRequest) (*FbMessagesResponse, error)
}
