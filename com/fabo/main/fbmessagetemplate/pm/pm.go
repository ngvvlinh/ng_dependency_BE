package pm

import (
	"context"

	"o.o/api/fabo/fbmessagetemplate"
	"o.o/api/main/identity"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus capi.EventBus

	messageTemplateAggregate fbmessagetemplate.CommandBus
}

var messageTemplateSamples = []fbmessagetemplate.FbMessageTemplate{
	{
		ShortCode: "Xin chào",
		Template:  "Chào bạn [Tên khách hàng], shop có thể giúp gì cho bạn?",
	},
	{
		ShortCode: "Hết hàng",
		Template:  "Sản phẩm đã hết hàng, shop sẽ liên hệ bạn khi hàng về ạ!",
	},
	{
		ShortCode: "Kiểm tra hàng",
		Template:  "Dạ chào anh/chị! Shop sẽ kiểm tra lại hàng và báo anh/chị ngay!",
	},
}

func NewProcessManager(
	eventBus bus.EventRegistry,
	messageTemplateAggregate fbmessagetemplate.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		eventBus:                 eventBus,
		messageTemplateAggregate: messageTemplateAggregate,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.AccountCreated)
}

func (m *ProcessManager) AccountCreated(ctx context.Context, event *identity.AccountCreatedEvent) error {
	for _, messageTemplate := range messageTemplateSamples {
		if err := m.messageTemplateAggregate.Dispatch(ctx, &fbmessagetemplate.CreateMessageTemplateCommand{
			ShopID:    event.ShopID,
			Template:  messageTemplate.Template,
			ShortCode: messageTemplate.ShortCode,
		}); err != nil {
			return err
		}
	}
	return nil
}
