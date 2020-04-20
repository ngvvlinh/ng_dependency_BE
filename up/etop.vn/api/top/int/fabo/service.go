package fabo

import (
	"context"

	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=fabo

// +apix:path=/fabo.Session
type SessionService interface {
	InitSession(context.Context, *InitSessionRequest) (*InitSessionResponse, error)
}

// +apix:path=/fabo.Page
type PageService interface {
	RemoveFbPages(context.Context, *RemoveFbPagesRequest) (*cm.Empty, error)
	ListFbPages(context.Context, *ListFbPagesRequest) (*FbPagesResponse, error)
}
