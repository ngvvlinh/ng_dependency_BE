package fbmessaging

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()
var scheme = conversion.Build(convert.RegisterConversions)

type FbExternalMessagingAggregate struct {
	db                          *cmsql.Database
	eventBus                    capi.EventBus
	fbExternalPostStore         sqlstore.FbExternalPostStoreFactory
	fbExternalCommentStore      sqlstore.FbExternalCommentStoreFactory
	fbExternalConversationStore sqlstore.FbExternalConversationStoreFactory
	fbExternalMessageStore      sqlstore.FbExternalMessageStoreFactory
	fbCustomerConversationStore sqlstore.FbCustomerConversationStoreFactory
}

func NewFbExternalMessagingAggregate(db *cmsql.Database, eventBus capi.EventBus) *FbExternalMessagingAggregate {
	return &FbExternalMessagingAggregate{
		db:                          db,
		eventBus:                    eventBus,
		fbExternalPostStore:         sqlstore.NewFbExternalPostStore(db),
		fbExternalCommentStore:      sqlstore.NewFbExternalCommentStore(db),
		fbExternalConversationStore: sqlstore.NewFbExternalConversationStore(db),
		fbExternalMessageStore:      sqlstore.NewFbExternalMessageStore(db),
		fbCustomerConversationStore: sqlstore.NewFbCustomerConversationStore(db),
	}
}

func (a *FbExternalMessagingAggregate) MessageBus() fbmessaging.CommandBus {
	b := bus.New()
	return fbmessaging.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbExternalMessagingAggregate) CreateFbExternalConversations(
	ctx context.Context, args fbmessaging.CreateFbExternalConversationsArgs,
) ([]*fbmessaging.FbExternalConversation, error) {
	newFbExternalConversations := make([]*fbmessaging.FbExternalConversation, 0, len(args.FbExternalConversations))
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalConversation := range args.FbExternalConversations {
			newFbExternalConversation := new(fbmessaging.FbExternalConversation)
			if err := scheme.Convert(fbExternalConversation, newFbExternalConversation); err != nil {
				return err
			}
			newFbExternalConversations = append(newFbExternalConversations, newFbExternalConversation)
		}
		if err := a.fbExternalConversationStore(ctx).CreateFbExternalConversations(newFbExternalConversations); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return newFbExternalConversations, nil
}

func (a *FbExternalMessagingAggregate) CreateFbExternalMessages(
	ctx context.Context, args fbmessaging.CreateFbExternalMessagesArgs,
) ([]*fbmessaging.FbExternalMessage, error) {
	newFbExternalMessages := make([]*fbmessaging.FbExternalMessage, 0, len(args.FbExternalMessages))
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalMessage := range args.FbExternalMessages {
			newFbExternalMessage := new(fbmessaging.FbExternalMessage)
			if err := scheme.Convert(fbExternalMessage, newFbExternalMessage); err != nil {
				return err
			}
			newFbExternalMessages = append(newFbExternalMessages, newFbExternalMessage)
		}
		if err := a.fbExternalMessageStore(ctx).CreateFbExternalMessages(newFbExternalMessages); err != nil {
			return err
		}
		event := &fbmessaging.FbExternalMessagesCreatedEvent{
			FbExternalMessages: newFbExternalMessages,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return newFbExternalMessages, nil
}

func (a *FbExternalMessagingAggregate) CreateFbCustomerConversations(
	ctx context.Context, args fbmessaging.CreateFbCustomerConversationsArgs,
) ([]*fbmessaging.FbCustomerConversation, error) {
	newFbCustomerConversations := make([]*fbmessaging.FbCustomerConversation, 0, len(args.FbCustomerConversations))
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbCustomerConversation := range args.FbCustomerConversations {
			newFbCustomerConversation := new(fbmessaging.FbCustomerConversation)
			if err := scheme.Convert(fbCustomerConversation, newFbCustomerConversation); err != nil {
				return err
			}
			newFbCustomerConversations = append(newFbCustomerConversations, newFbCustomerConversation)
		}
		if err := a.fbCustomerConversationStore(ctx).CreateFbCustomerConversations(newFbCustomerConversations); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return newFbCustomerConversations, nil
}
