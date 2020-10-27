package fbmessagetemplate

import (
	"context"

	"o.o/api/fabo/fbmessagetemplate"
	"o.o/api/top/types/common"
	"o.o/backend/com/fabo/main/fbmessagetemplate/sqlstore"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
)

var _ fbmessagetemplate.Aggregate = &FbMessageTemplateAggregate{}

type FbMessageTemplateAggregate struct {
	fbmessagetemplateStore sqlstore.FbMessageTemplateStoreFactory
}

func FbMessageTemplateAggregateMessageBus(a *FbMessageTemplateAggregate) fbmessagetemplate.CommandBus {
	b := bus.New()
	return fbmessagetemplate.NewAggregateHandler(a).RegisterHandlers(b)
}

func NewFbMessageTemplateAggregate(db com.MainDB) *FbMessageTemplateAggregate {
	return &FbMessageTemplateAggregate{
		fbmessagetemplateStore: sqlstore.NewFbMessageTemplateStoreFactory(db),
	}
}

func (s *FbMessageTemplateAggregate) CreateMessageTemplate(
	ctx context.Context, createTemplate *fbmessagetemplate.CreateFbMessageTemplate,
) (*fbmessagetemplate.FbMessageTemplate, error) {
	return s.fbmessagetemplateStore(ctx).CreateFbMessageTemplate(createTemplate)
}

func (s *FbMessageTemplateAggregate) UpdateMessageTemplate(ctx context.Context, template *fbmessagetemplate.UpdateFbMessageTemplate) (*common.Empty, error) {
	args := &fbmessagetemplate.UpdateFbMessageTemplate{
		ID:        template.ID,
		ShopID:    template.ShopID,
		Template:  template.Template,
		ShortCode: template.ShortCode,
	}
	err := s.fbmessagetemplateStore(ctx).UpdateMessageTemplate(args)
	if err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (s *FbMessageTemplateAggregate) DeleteMessageTemplate(ctx context.Context, args *fbmessagetemplate.DeleteFbMessageTemplate) (*common.Empty, error) {
	err := s.fbmessagetemplateStore(ctx).ID(args.ID).ShopID(args.ShopID).DeleteMessageTemplate()
	if err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}
