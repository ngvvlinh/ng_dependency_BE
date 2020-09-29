package fbcustomerconversationsearch

import (
	"time"

	"o.o/capi/dot"
)

// +convert:create=FbExternalCommentSearch
type FbExternalCommentSearch struct {
	ExternalPageID      string
	ExternalMessageNorm string
	ExternalPostID      string
	ExternalUserID      string
	CreatedAt           time.Time
}

// +convert:create=FbExternalMessageSearch
type FbExternalMessageSearch struct {
	ExternalPageID         string
	ExternalMessageNorm    string
	ExternalConversationID string
	CreatedAt              time.Time
}

// +convert:create=FbCustomerConversationSearch
type FbCustomerConversationSearch struct {
	ID                   dot.ID
	ExternalUserNameNorm string
	ExternalPostID       string
	CreatedAt            time.Time
}
