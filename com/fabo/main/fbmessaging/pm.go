package fbmessaging

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
)

type ProcessManager struct {
	eventBus     capi.EventBus
	fbmessagingQ fbmessaging.QueryBus
	fbmessagingA fbmessaging.CommandBus
}

func NewProcessManager(
	eventBus capi.EventBus,
	fbmessagingQuery fbmessaging.QueryBus,
	fbmessagingAggregate fbmessaging.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:     eventBus,
		fbmessagingQ: fbmessagingQuery,
		fbmessagingA: fbmessagingAggregate,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.CreateFbExternalMessages)
}

func (m *ProcessManager) CreateFbExternalMessages(ctx context.Context, event *fbmessaging.FbExternalMessagesCreatedEvent) error {
	var listExternalConversationIDs []string
	mapExternalConversationsIDs := make(map[string]bool)
	for _, fbExternalMessage := range event.FbExternalMessages {
		if ok := mapExternalConversationsIDs[fbExternalMessage.ExternalConversationID]; !ok {
			listExternalConversationIDs = append(listExternalConversationIDs, fbExternalMessage.ExternalConversationID)
		} else {
			mapExternalConversationsIDs[fbExternalMessage.ExternalConversationID] = true
		}
	}

	listLastExternalMessagesQuery := &fbmessaging.ListLatestFbExternalMessagesQuery{
		ExternalConversationIDs: listExternalConversationIDs,
	}
	if err := m.fbmessagingQ.Dispatch(ctx, listLastExternalMessagesQuery); err != nil {
		return err
	}
	mapLastExternalMessages := make(map[string]*fbmessaging.FbExternalMessage)
	for _, fbExternalMessage := range listLastExternalMessagesQuery.Result {
		mapLastExternalMessages[fbExternalMessage.ExternalConversationID] = fbExternalMessage
	}

	listFbCustomerConversationsQuery := &fbmessaging.ListFbCustomerConversationsByExternalIDsQuery{
		ExternalIDs: listExternalConversationIDs,
	}
	if err := m.fbmessagingQ.Dispatch(ctx, listFbCustomerConversationsQuery); err != nil {
		return err
	}

	var updateFbCustomerConversationsCmd []*fbmessaging.CreateFbCustomerConversationArgs
	for _, oldFbCustomerConversation := range listFbCustomerConversationsQuery.Result {
		updateFbCustomerConversationsCmd = append(updateFbCustomerConversationsCmd, &fbmessaging.CreateFbCustomerConversationArgs{
			ID:               oldFbCustomerConversation.ID,
			FbPageID:         oldFbCustomerConversation.FbPageID,
			ExternalID:       oldFbCustomerConversation.ExternalID,
			ExternalUserID:   oldFbCustomerConversation.ExternalUserID,
			ExternalUserName: oldFbCustomerConversation.ExternalUserName,
			IsRead:           false,
			Type:             oldFbCustomerConversation.Type,
			PostAttachments:  oldFbCustomerConversation.PostAttachments,
			LastMessage:      mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalMessage,
			LastMessageAt:    mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalCreatedTime,
		})
	}

	if err := m.fbmessagingA.Dispatch(ctx, &fbmessaging.CreateFbCustomerConversationsCommand{
		FbCustomerConversations: updateFbCustomerConversationsCmd,
	}); err != nil {
		return err
	}
	return nil
}
