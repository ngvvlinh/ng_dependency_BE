package fabo

import (
	"context"

	cm "etop.vn/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=fabo

// +apix:path=/fabo.Page
type PageService interface {
	ConnectPages(context.Context, *ConnectPagesRequest) (*ConnectPagesResponse, error)
	RemovePages(context.Context, *RemovePagesRequest) (*cm.Empty, error)
	ListPages(context.Context, *ListPagesRequest) (*ListPagesResponse, error)
}
