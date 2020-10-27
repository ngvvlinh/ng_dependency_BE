package fbmessagetemplate

import (
	"context"

	cm "o.o/api/top/types/common"
	"o.o/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateMessageTemplate(context.Context, *CreateFbMessageTemplate) (*FbMessageTemplate, error)
	UpdateMessageTemplate(context.Context, *UpdateFbMessageTemplate) (*cm.Empty, error)
	DeleteMessageTemplate(context.Context, *DeleteFbMessageTemplate) (*cm.Empty, error)
}

type QueryService interface {
	GetMessageTemplates(ctx context.Context, shopID dot.ID) ([]*FbMessageTemplate, error)
}
