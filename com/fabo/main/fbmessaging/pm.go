package fbmessaging

import (
	"context"
	"fmt"
	"sync"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

type ProcessManager struct {
	eventBus     capi.EventBus
	fbmessagingQ fbmessaging.QueryBus
	fbmessagingA fbmessaging.CommandBus
	fbpagingQ    fbpaging.QueryBus
	fbuseringQ   fbusering.QueryBus
	fbuseringA   fbusering.CommandBus
}

func NewProcessManager(
	eventBus capi.EventBus,
	fbmessagingQuery fbmessaging.QueryBus,
	fbmessagingAggregate fbmessaging.CommandBus,
	fbpagingQuery fbpaging.QueryBus, fbuseringQ fbusering.QueryBus,
	fbuseringA fbusering.CommandBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:     eventBus,
		fbmessagingQ: fbmessagingQuery,
		fbmessagingA: fbmessagingAggregate,
		fbpagingQ:    fbpagingQuery,
		fbuseringQ:   fbuseringQ,
		fbuseringA:   fbuseringA,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.HandleFbExternalMessagesCreatedEvent)
	eventBus.AddEventListener(m.HandleFbExternalCommentsCreatedEvent)
	eventBus.AddEventListener(m.HandleFbExternalConversationsCreatedEvent)
}

func (m *ProcessManager) HandleFbExternalMessagesCreatedEvent(
	ctx context.Context, event *fbmessaging.FbExternalMessagesCreatedEvent,
) error {
	if len(event.FbExternalMessages) == 0 {
		return nil
	}

	var externalConversationIDs []string
	mapLastExternalMessages := make(map[string]*fbmessaging.FbExternalMessage)
	{
		setExternalConversationIDs := NewSet()
		for _, fbExternalMessage := range event.FbExternalMessages {
			setExternalConversationIDs.Push(fbExternalMessage.ExternalConversationID)
		}
		externalConversationIDs = setExternalConversationIDs.Strings()

		listLastExternalMessagesQuery := &fbmessaging.ListLatestFbExternalMessagesQuery{
			ExternalConversationIDs: externalConversationIDs,
		}
		if err := m.fbmessagingQ.Dispatch(ctx, listLastExternalMessagesQuery); err != nil {
			return err
		}

		for _, fbExternalMessage := range listLastExternalMessagesQuery.Result {
			mapLastExternalMessages[fbExternalMessage.ExternalConversationID] = fbExternalMessage
		}
	}

	listFbCustomerConversationsQuery := &fbmessaging.ListFbCustomerConversationsByExternalIDsQuery{
		ExternalIDs: externalConversationIDs,
	}
	if err := m.fbmessagingQ.Dispatch(ctx, listFbCustomerConversationsQuery); err != nil {
		return err
	}

	var updateFbCustomerConversationsCmd []*fbmessaging.CreateFbCustomerConversationArgs
	for _, oldFbCustomerConversation := range listFbCustomerConversationsQuery.Result {
		externalMessageAttachments := mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalAttachments
		if mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalSticker != "" {
			externalMessageAttachments = []*fbmessaging.FbMessageAttachment{
				convertExternalStickerToMessageAttachment(mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalSticker),
			}
		}
		updateFbCustomerConversationsCmd = append(updateFbCustomerConversationsCmd, &fbmessaging.CreateFbCustomerConversationArgs{
			ID:                         oldFbCustomerConversation.ID,
			FbPageID:                   oldFbCustomerConversation.FbPageID,
			ExternalPageID:             oldFbCustomerConversation.ExternalPageID,
			ExternalID:                 oldFbCustomerConversation.ExternalID,
			ExternalUserID:             oldFbCustomerConversation.ExternalUserID,
			ExternalUserName:           oldFbCustomerConversation.ExternalUserName,
			ExternalFrom:               mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalFrom,
			IsRead:                     false,
			Type:                       oldFbCustomerConversation.Type,
			ExternalPostAttachments:    oldFbCustomerConversation.ExternalPostAttachments,
			ExternalCommentAttachment:  oldFbCustomerConversation.ExternalCommentAttachment,
			ExternalMessageAttachments: externalMessageAttachments,
			LastMessage:                mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalMessage,
			LastMessageAt:              mapLastExternalMessages[oldFbCustomerConversation.ExternalID].ExternalCreatedTime,
		})
	}

	if len(updateFbCustomerConversationsCmd) > 0 {
		if err := m.fbmessagingA.Dispatch(ctx, &fbmessaging.CreateFbCustomerConversationsCommand{
			FbCustomerConversations: updateFbCustomerConversationsCmd,
		}); err != nil {
			return err
		}
	}

	// Handle create new FbExternalUser (customer)
	{
		mapExternalUserID := make(map[string]bool)
		var fbObjectFroms []*fbmessaging.FbObjectFrom
		for _, externalMessage := range event.FbExternalMessages {
			if externalMessage.ExternalFrom == nil {
				continue
			}
			if _, ok := mapExternalUserID[externalMessage.ExternalFrom.ID]; !ok {
				mapExternalUserID[externalMessage.ExternalFrom.ID] = true
				fbObjectFroms = append(fbObjectFroms, externalMessage.ExternalFrom)
			}
		}

		// TODO: Ngoc
		//if err := m.handleCreateExternalCustomerUser(ctx, fbObjectFroms); err != nil {
		//	return err
		//}
	}

	return nil
}

func (m *ProcessManager) HandleFbExternalCommentsCreatedEvent(
	ctx context.Context, event *fbmessaging.FbExternalCommentsCreatedEvent,
) error {
	if len(event.FbExternalComments) == 0 {
		return nil
	}

	externalPostIDsSet := NewSet()
	mapExternalPostIDAndMapExternalUserIDAndExternalPageID := make(map[string]map[string]string)
	for _, fbExternalComment := range event.FbExternalComments {
		// Ignore ExternalUserID = "" and ExternalFrom = nil
		if fbExternalComment.ExternalUserID == "" || fbExternalComment.ExternalFrom == nil {
			continue
		}

		// Ignore UserID == PageID
		if fbExternalComment.ExternalUserID == fbExternalComment.ExternalPageID &&
			(fbExternalComment.ExternalParentUserID == "" || fbExternalComment.ExternalParent == nil || fbExternalComment.ExternalParent.From == nil) {
			continue
		}

		if _, ok := mapExternalPostIDAndMapExternalUserIDAndExternalPageID[fbExternalComment.ExternalPostID]; !ok {
			mapExternalPostIDAndMapExternalUserIDAndExternalPageID[fbExternalComment.ExternalPostID] = make(map[string]string)
			externalPostIDsSet.Push(fbExternalComment.ExternalPostID)
		}

		tempExternalUserID := fbExternalComment.ExternalUserID
		if fbExternalComment.ExternalUserID == fbExternalComment.ExternalPageID {
			tempExternalUserID = fbExternalComment.ExternalParentUserID
		}

		mapExternalPostIDAndMapExternalUserIDAndExternalPageID[fbExternalComment.ExternalPostID][tempExternalUserID] = fbExternalComment.ExternalPageID
	}

	listExternalPostsByExternalIDsQuery := &fbmessaging.ListFbExternalPostsByExternalIDsQuery{
		ExternalIDs: externalPostIDsSet.Strings(),
	}
	if err := m.fbmessagingQ.Dispatch(ctx, listExternalPostsByExternalIDsQuery); err != nil {
		return err
	}
	mapExternalPost := make(map[string]*fbmessaging.FbExternalPost)
	for _, externalPost := range listExternalPostsByExternalIDsQuery.Result {
		mapExternalPost[externalPost.ExternalID] = externalPost
	}

	for externalPostID, mapExternalUserIDAndExternalPageID := range mapExternalPostIDAndMapExternalUserIDAndExternalPageID {
		for externalUserID, externalPageID := range mapExternalUserIDAndExternalPageID {
			getLatestFbExternalCommentQuery := &fbmessaging.GetLatestFbExternalCommentQuery{
				ExternalPostID: externalPostID,
				ExternalUserID: externalUserID,
				ExternalPageID: externalPageID,
			}
			if err := m.fbmessagingQ.Dispatch(ctx, getLatestFbExternalCommentQuery); err != nil {
				return err
			}
			lastFbExternalComment := getLatestFbExternalCommentQuery.Result

			externalUserIDCustomerConversation := lastFbExternalComment.ExternalUserID
			externalUserNameCustomerConversation := lastFbExternalComment.ExternalFrom.Name
			externalFromCustomerConversation := lastFbExternalComment.ExternalFrom
			if lastFbExternalComment.ExternalUserID == lastFbExternalComment.ExternalPageID {
				externalUserIDCustomerConversation = lastFbExternalComment.ExternalParent.From.ID
				externalUserNameCustomerConversation = lastFbExternalComment.ExternalParent.From.Name
			}

			getFbCustomerConversationQuery := &fbmessaging.GetFbCustomerConversationQuery{
				ExternalID:               externalPostID,
				ExternalUserID:           externalUserIDCustomerConversation,
				CustomerConversationType: fb_customer_conversation_type.Comment,
			}
			if err := m.fbmessagingQ.Dispatch(ctx, getFbCustomerConversationQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
				return err
			}
			oldFbCustomerConversation := getFbCustomerConversationQuery.Result

			ID := cm.NewID()
			if oldFbCustomerConversation != nil {
				ID = oldFbCustomerConversation.ID
			}

			if err := m.fbmessagingA.Dispatch(ctx, &fbmessaging.CreateFbCustomerConversationsCommand{
				FbCustomerConversations: []*fbmessaging.CreateFbCustomerConversationArgs{
					{
						ID:                        ID,
						FbPageID:                  lastFbExternalComment.FbPageID,
						ExternalPageID:            lastFbExternalComment.ExternalPageID,
						ExternalID:                lastFbExternalComment.ExternalPostID,
						ExternalUserID:            externalUserIDCustomerConversation,
						ExternalUserName:          externalUserNameCustomerConversation,
						ExternalFrom:              externalFromCustomerConversation,
						IsRead:                    false,
						Type:                      fb_customer_conversation_type.Comment,
						ExternalPostAttachments:   mapExternalPost[externalPostID].ExternalAttachments,
						ExternalCommentAttachment: lastFbExternalComment.ExternalAttachment,
						LastMessage:               lastFbExternalComment.ExternalMessage,
						LastMessageAt:             lastFbExternalComment.ExternalCreatedTime,
					},
				},
			}); err != nil {
				return err
			}
		}
	}

	// Handle create new FbExternalUser (customer)
	{
		mapExternalUserID := make(map[string]bool)
		var fbObjectFroms []*fbmessaging.FbObjectFrom
		for _, externalComment := range event.FbExternalComments {
			if externalComment.ExternalFrom == nil {
				continue
			}
			if _, ok := mapExternalUserID[externalComment.ExternalFrom.ID]; !ok {
				mapExternalUserID[externalComment.ExternalFrom.ID] = true
				fbObjectFroms = append(fbObjectFroms, externalComment.ExternalFrom)
			}
		}

		// TODO: Ngoc
		//if err := m.handleCreateExternalCustomerUser(ctx, fbObjectFroms); err != nil {
		//	return err
		//}
	}

	return nil
}

func (m *ProcessManager) HandleFbExternalConversationsCreatedEvent(
	ctx context.Context, event *fbmessaging.FbExternalConversationsCreatedEvent,
) error {
	if len(event.FbExternalConversations) == 0 {
		return nil
	}

	mapOldExternalConversationID := make(map[string]bool)
	{
		externalConversationIDsSet := NewSet()
		for _, externalConversation := range event.FbExternalConversations {
			externalConversationIDsSet.Push(externalConversation.ExternalID)
		}

		listOldCustomerConversationsQuery := &fbmessaging.ListFbCustomerConversationsByExternalIDsQuery{
			ExternalIDs: externalConversationIDsSet.Strings(),
		}
		if err := m.fbmessagingQ.Dispatch(ctx, listOldCustomerConversationsQuery); err != nil {
			return err
		}

		for _, oldExternalConversation := range listOldCustomerConversationsQuery.Result {
			mapOldExternalConversationID[oldExternalConversation.ExternalID] = true
		}
	}

	externalConversations := event.FbExternalConversations

	var createFbCustomerConversationsArgs []*fbmessaging.CreateFbCustomerConversationArgs
	for _, externalConversation := range externalConversations {
		if _, ok := mapOldExternalConversationID[externalConversation.ExternalID]; ok {
			continue
		}
		createFbCustomerConversationsArgs = append(createFbCustomerConversationsArgs, &fbmessaging.CreateFbCustomerConversationArgs{
			ID:               cm.NewID(),
			FbPageID:         externalConversation.FbPageID,
			ExternalPageID:   externalConversation.ExternalPageID,
			ExternalID:       externalConversation.ExternalID,
			ExternalUserID:   externalConversation.ExternalUserID,
			ExternalUserName: externalConversation.ExternalUserName,
			Type:             fb_customer_conversation_type.Message,
		})
	}

	if len(createFbCustomerConversationsArgs) > 0 {
		if err := m.fbmessagingA.Dispatch(ctx, &fbmessaging.CreateFbCustomerConversationsCommand{
			FbCustomerConversations: createFbCustomerConversationsArgs,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProcessManager) handleCreateExternalCustomerUser(
	ctx context.Context, fbObjectFroms []*fbmessaging.FbObjectFrom,
) error {
	if len(fbObjectFroms) == 0 {
		return nil
	}

	var externalUserIDs []string
	for _, fbObjectFrom := range fbObjectFroms {
		externalUserIDs = append(externalUserIDs, fbObjectFrom.ID)
	}

	listExternalUsersByExternalIDsQuery := &fbusering.ListFbExternalUsersByExternalIDQuery{
		ExternalIDs: externalUserIDs,
	}
	if err := m.fbuseringQ.Dispatch(ctx, listExternalUsersByExternalIDsQuery); err != nil {
		return err
	}

	mapOldExternalUser := make(map[string]*fbusering.FbExternalUser)
	for _, oldExternalUser := range listExternalUsersByExternalIDsQuery.Result {
		mapOldExternalUser[oldExternalUser.ExternalID] = oldExternalUser
	}

	var createFbExternalUsersArgs []*fbusering.CreateFbExternalUserArgs
	for _, fbObjectFrom := range fbObjectFroms {
		if oldExternalUser, ok := mapOldExternalUser[fbObjectFrom.ID]; ok && oldExternalUser.ExternalInfo == nil {
			createFbExternalUsersArgs = append(createFbExternalUsersArgs, &fbusering.CreateFbExternalUserArgs{
				ID:         oldExternalUser.ID,
				UserID:     oldExternalUser.UserID,
				ExternalID: oldExternalUser.ExternalID,
				ExternalInfo: &fbusering.FbExternalUserInfo{
					Name: fbObjectFrom.Name,
				},
				Status: oldExternalUser.Status,
			})
		} else {
			createFbExternalUsersArgs = append(createFbExternalUsersArgs, &fbusering.CreateFbExternalUserArgs{
				ID:         cm.NewID(),
				ExternalID: fbObjectFrom.ID,
				ExternalInfo: &fbusering.FbExternalUserInfo{
					Name: fbObjectFrom.Name,
				},
				Status: status3.P,
			})
		}
	}

	return nil
}

func convertExternalStickerToMessageAttachment(externalSticker string) *fbmessaging.FbMessageAttachment {
	if externalSticker == "" {
		return nil
	}
	return &fbmessaging.FbMessageAttachment{
		ImageData: &fbmessaging.FbMessageAttachmentImageData{
			Width:           240,
			Height:          240,
			URL:             externalSticker,
			PreviewURL:      externalSticker,
			RenderAsSticker: true,
		},
		MimeType: "sticker",
		FileURL:  externalSticker,
	}
}

type Set struct {
	mapValues map[interface{}]bool
	mu        sync.Mutex
}

func NewSet() *Set {
	return &Set{
		mapValues: make(map[interface{}]bool),
		mu:        sync.Mutex{},
	}
}

func (h *Set) Push(arg interface{}) *Set {
	h.mu.Lock()
	h.mapValues[arg] = true
	h.mu.Unlock()
	return h
}

func (h *Set) Delete(arg interface{}) *Set {
	h.mu.Lock()
	delete(h.mapValues, arg)
	h.mu.Unlock()
	return h
}

func (h *Set) Values() []interface{} {
	var result []interface{}
	h.mu.Lock()
	for key, _ := range h.mapValues {
		result = append(result, key)
	}
	h.mu.Unlock()
	return result
}

func (h *Set) Strings() []string {
	var result []string
	h.mu.Lock()
	for key, _ := range h.mapValues {
		switch v := key.(type) {
		case string:
			result = append(result, v)
		default:
			panic(fmt.Sprintf("%v is not of type string", key))
		}
	}
	h.mu.Unlock()
	return result
}

func (h *Set) IDs() []dot.ID {
	var result []dot.ID
	h.mu.Lock()
	for key, _ := range h.mapValues {
		switch v := key.(type) {
		case dot.ID:
			result = append(result, v)
		default:
			panic(fmt.Sprintf("%v %T is not of type string", key, key))
		}
	}
	h.mu.Unlock()
	return result
}
