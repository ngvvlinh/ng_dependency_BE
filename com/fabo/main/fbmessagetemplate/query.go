package fbmessagetemplate

import (
	"context"

	"o.o/api/fabo/fbmessagetemplate"
	"o.o/backend/com/fabo/main/fbmessagetemplate/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
)

var _ fbmessagetemplate.QueryService = &FbMessageTemplateQuery{}

type FbMessageTemplateQuery struct {
	fbmessagetemplateStore sqlstore.FbMessageTemplateStoreFactory
}

func NewFbMessagingQuery(database com.MainDB) *FbMessageTemplateQuery {
	return &FbMessageTemplateQuery{
		fbmessagetemplateStore: sqlstore.NewFbMessageTemplateStoreFactory(database),
	}
}

func FbMessagingQueryMessageBus(q *FbMessageTemplateQuery) fbmessagetemplate.QueryBus {
	b := bus.New()
	return fbmessagetemplate.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *FbMessageTemplateQuery) GetMessageTemplates(ctx context.Context, shopID dot.ID) ([]*fbmessagetemplate.FbMessageTemplate, error) {
	return q.fbmessagetemplateStore(ctx).ShopID(shopID).GetFbMessageTemplates()
}
