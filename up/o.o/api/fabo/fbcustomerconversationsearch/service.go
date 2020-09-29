package fbcustomerconversationsearch

import "context"

// +gen:api

type QueryService interface {
	ListFbExternalCommentSearch(ctx context.Context, pageIDs []string, externalMsg string) ([]*FbExternalCommentSearch, error)
	ListFbExternalMessageSearch(ctx context.Context, pageIDs []string, externalMsg string) ([]*FbExternalMessageSearch, error)
	ListFbExternalConversationSearch(ctx context.Context, pageIDs []string, extUserName string) ([]*FbCustomerConversationSearch, error)
}
