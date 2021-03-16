package fbmessaging

import (
	"context"
	"strings"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/backend/com/fabo/main/compare"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

type FbExternalMessagingAggregate struct {
	db                               *cmsql.Database
	eventBus                         capi.EventBus
	fbExternalPostStore              sqlstore.FbExternalPostStoreFactory
	fbExternalCommentStore           sqlstore.FbExternalCommentStoreFactory
	fbExternalConversationStore      sqlstore.FbExternalConversationStoreFactory
	fbExternalMessageStore           sqlstore.FbExternalMessageStoreFactory
	fbCustomerConversationStore      sqlstore.FbCustomerConversationStoreFactory
	fbCustomerConversationStateStore sqlstore.FbCustomerConversationStateStoreFactory
	fbClient                         *fbclient.FbClient
}

func NewFbExternalMessagingAggregate(db com.MainDB, eventBus capi.EventBus, client *fbclient.FbClient) *FbExternalMessagingAggregate {
	return &FbExternalMessagingAggregate{
		db:                               db,
		eventBus:                         eventBus,
		fbExternalPostStore:              sqlstore.NewFbExternalPostStore(db),
		fbExternalCommentStore:           sqlstore.NewFbExternalCommentStore(db),
		fbExternalConversationStore:      sqlstore.NewFbExternalConversationStore(db),
		fbExternalMessageStore:           sqlstore.NewFbExternalMessageStore(db),
		fbCustomerConversationStore:      sqlstore.NewFbCustomerConversationStore(db),
		fbCustomerConversationStateStore: sqlstore.NewFbCustomerConversationStateStore(db),
		fbClient:                         client,
	}
}

func FbExternalMessagingAggregateMessageBus(a *FbExternalMessagingAggregate) fbmessaging.CommandBus {
	b := bus.New()
	return fbmessaging.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbExternalMessagingAggregate) CreateFbExternalMessagesFromSync(
	ctx context.Context, args *fbmessaging.CreateFbExternalMessagesFromSyncArgs,
) ([]*fbmessaging.FbExternalMessage, error) {
	var fbExternalMessageIDs []string
	for _, fbExternalMessage := range args.FbExternalMessages {
		fbExternalMessageIDs = append(fbExternalMessageIDs, fbExternalMessage.ExternalID)
	}

	oldFbExternalMessages, err := a.fbExternalMessageStore(ctx).ExternalIDs(fbExternalMessageIDs).ListFbExternalMessages()
	if err != nil {
		return nil, err
	}
	mapOldFbExternalMessage := make(map[string]*fbmessaging.FbExternalMessage)
	for _, oldFbExternalMessage := range oldFbExternalMessages {
		mapOldFbExternalMessage[oldFbExternalMessage.ExternalID] = oldFbExternalMessage
	}

	var newFbExternalMessages []*fbmessaging.FbExternalMessage
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalMessageArg := range args.FbExternalMessages {
			newFbExternalMessage := new(fbmessaging.FbExternalMessage)
			*newFbExternalMessage = *fbExternalMessageArg

			// ignore old message
			if _, ok := mapOldFbExternalMessage[fbExternalMessageArg.ExternalID]; ok {
				continue
			}

			newFbExternalMessages = append(newFbExternalMessages, newFbExternalMessage)
		}

		if len(newFbExternalMessages) > 0 {
			if err := a.fbExternalMessageStore(ctx).CreateFbExternalMessages(newFbExternalMessages); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalMessagesCreatedEvent{
				FbExternalMessages: newFbExternalMessages,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return newFbExternalMessages, nil
}

func (a *FbExternalMessagingAggregate) CreateOrUpdateFbExternalMessages(
	ctx context.Context, args *fbmessaging.CreateOrUpdateFbExternalMessagesArgs,
) ([]*fbmessaging.FbExternalMessage, error) {
	var fbExternalMessageIDs []string
	for _, fbExternalMessage := range args.FbExternalMessages {
		fbExternalMessageIDs = append(fbExternalMessageIDs, fbExternalMessage.ExternalID)
	}

	oldFbExternalMessages, err := a.fbExternalMessageStore(ctx).ExternalIDs(fbExternalMessageIDs).ListFbExternalMessages()
	if err != nil {
		return nil, err
	}
	mapOldFbExternalMessage := make(map[string]*fbmessaging.FbExternalMessage)
	for _, oldFbExternalMessage := range oldFbExternalMessages {
		mapOldFbExternalMessage[oldFbExternalMessage.ExternalID] = oldFbExternalMessage
	}

	var resultFbExternalMessages []*fbmessaging.FbExternalMessage
	var newFbExternalMessages []*fbmessaging.FbExternalMessage
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalMessageArg := range args.FbExternalMessages {
			newFbExternalMessage := new(fbmessaging.FbExternalMessage)
			if err := scheme.Convert(fbExternalMessageArg, newFbExternalMessage); err != nil {
				return err
			}

			if oldFbExternalMessage, ok := mapOldFbExternalMessage[fbExternalMessageArg.ExternalID]; ok {
				newFbExternalMessage.ID = oldFbExternalMessage.ID
				if err := a.fbExternalMessageStore(ctx).ExternalID(newFbExternalMessage.ExternalID).UpdateFbExternalMessage(newFbExternalMessage); err != nil {
					return err
				}

				resultFbExternalMessages = append(resultFbExternalMessages, newFbExternalMessage)
			} else {
				resultFbExternalMessages = append(resultFbExternalMessages, newFbExternalMessage)
				newFbExternalMessages = append(newFbExternalMessages, newFbExternalMessage)
			}
		}

		if len(newFbExternalMessages) > 0 {
			if err := a.fbExternalMessageStore(ctx).CreateFbExternalMessages(newFbExternalMessages); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalMessagesCreatedEvent{
				FbExternalMessages: newFbExternalMessages,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		return nil

	}); err != nil {
		return nil, err
	}
	return resultFbExternalMessages, nil
}

func (a *FbExternalMessagingAggregate) CreateOrUpdateFbExternalConversations(
	ctx context.Context, args *fbmessaging.CreateOrUpdateFbExternalConversationsArgs,
) ([]*fbmessaging.FbExternalConversation, error) {
	var fbExternalConversationIDs []string
	for _, fbExternalConversation := range args.FbExternalConversations {
		fbExternalConversationIDs = append(fbExternalConversationIDs, fbExternalConversation.ExternalID)
	}

	oldFbExternalConversations, err := a.fbExternalConversationStore(ctx).ExternalIDs(fbExternalConversationIDs).ListFbExternalConversations()
	if err != nil {
		return nil, err
	}
	mapOldFbExternalConversation := make(map[string]*fbmessaging.FbExternalConversation)
	for _, oldFbExternalConversation := range oldFbExternalConversations {
		mapOldFbExternalConversation[oldFbExternalConversation.ExternalID] = oldFbExternalConversation
	}

	var resultFbExternalConversations []*fbmessaging.FbExternalConversation
	var newFbExternalConversations []*fbmessaging.FbExternalConversation
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalConversationArg := range args.FbExternalConversations {
			newFbExternalConversation := new(fbmessaging.FbExternalConversation)
			if err := scheme.Convert(fbExternalConversationArg, newFbExternalConversation); err != nil {
				return err
			}

			if oldFbExternalConversation, ok := mapOldFbExternalConversation[fbExternalConversationArg.ExternalID]; ok {
				newFbExternalConversation.ID = oldFbExternalConversation.ID
				resultFbExternalConversations = append(resultFbExternalConversations, oldFbExternalConversation)

				if isEqual := compare.Compare(oldFbExternalConversation, newFbExternalConversation); isEqual {
					continue
				}
			} else {
				resultFbExternalConversations = append(resultFbExternalConversations, newFbExternalConversation)
			}

			newFbExternalConversations = append(newFbExternalConversations, newFbExternalConversation)
		}

		if len(newFbExternalConversations) > 0 {
			if err := a.fbExternalConversationStore(ctx).CreateFbExternalConversations(newFbExternalConversations); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalConversationsCreatedEvent{
				FbExternalConversations: newFbExternalConversations,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return nil, err
	}
	return resultFbExternalConversations, nil
}

func (a *FbExternalMessagingAggregate) CreateOrUpdateFbExternalConversation(
	ctx context.Context, fbExternalConversationArgs *fbmessaging.FbExternalConversation,
) (*fbmessaging.FbExternalConversation, error) {
	fbExternalConversationID := fbExternalConversationArgs.ExternalID

	oldFbExternalConversation, err := a.fbExternalConversationStore(ctx).ExternalID(fbExternalConversationID).GetFbExternalConversation()
	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}

	fbExternalConversation := new(fbmessaging.FbExternalConversation)
	*fbExternalConversation = *fbExternalConversationArgs

	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if oldFbExternalConversation != nil {
			// update fbExternalConversation
			fbExternalConversation.ID = oldFbExternalConversation.ID
			if err := a.fbExternalConversationStore(ctx).ExternalID(fbExternalConversationID).UpdateFbExternalConversation(fbExternalConversation); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalConversationsUpdatedEvent{
				FbExternalConversations: []*fbmessaging.FbExternalConversation{fbExternalConversation},
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		} else {
			// create fbExternalConversation
			if err := a.fbExternalConversationStore(ctx).CreateFbExternalConversation(fbExternalConversation); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalConversationsCreatedEvent{
				FbExternalConversations: []*fbmessaging.FbExternalConversation{fbExternalConversation},
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return fbExternalConversation, nil
}

func (a *FbExternalMessagingAggregate) CreateFbCustomerConversations(
	ctx context.Context, args *fbmessaging.CreateFbCustomerConversationsArgs,
) ([]*fbmessaging.FbCustomerConversation, error) {
	newFbCustomerConversations := make([]*fbmessaging.FbCustomerConversation, 0, len(args.FbCustomerConversations))
	newFbCustomerConversationStates := make([]*fbmessaging.FbCustomerConversationState, 0, len(args.FbCustomerConversations))

	for _, fbCustomerConversation := range args.FbCustomerConversations {
		newFbCustomerConversation := new(fbmessaging.FbCustomerConversation)
		if err := scheme.Convert(fbCustomerConversation, newFbCustomerConversation); err != nil {
			return nil, err
		}
		newFbCustomerConversations = append(newFbCustomerConversations, newFbCustomerConversation)
		newFbCustomerConversationStates = append(newFbCustomerConversationStates, &fbmessaging.FbCustomerConversationState{
			ID:             newFbCustomerConversation.ID,
			IsRead:         newFbCustomerConversation.IsRead,
			ExternalPageID: newFbCustomerConversation.ExternalPageID,
			UpdatedAt:      time.Now(),
		})
	}

	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if len(newFbCustomerConversations) > 0 {
			if err := a.fbCustomerConversationStore(ctx).CreateFbCustomerConversations(newFbCustomerConversations); err != nil {
				return err
			}

			if err := a.fbCustomerConversationStateStore(ctx).CreateFbCustomerConversationStates(newFbCustomerConversationStates); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return newFbCustomerConversations, nil
}

func (a *FbExternalMessagingAggregate) CreateOrUpdateFbCustomerConversations(
	ctx context.Context, args *fbmessaging.CreateOrUpdateFbCustomerConversationsArgs,
) ([]*fbmessaging.FbCustomerConversation, error) {
	var customerConversationArgsTypeMessage []*fbmessaging.CreateFbCustomerConversationArgs
	var customerConversationArgsTypeComment []*fbmessaging.CreateFbCustomerConversationArgs
	var newCustomerConversations []*fbmessaging.FbCustomerConversation

	for _, customerConversationArg := range args.FbCustomerConversations {
		switch customerConversationArg.Type {
		case fb_customer_conversation_type.Comment:
			customerConversationArgsTypeComment = append(customerConversationArgsTypeComment, customerConversationArg)
		case fb_customer_conversation_type.Message:
			customerConversationArgsTypeMessage = append(customerConversationArgsTypeMessage, customerConversationArg)
		}
	}

	// Handle type message
	if len(customerConversationArgsTypeMessage) > 0 {
		mapOldCustomerConversation := make(map[string]*fbmessaging.FbCustomerConversation)
		{
			var externalIDs []string
			mapExternalID := make(map[string]bool)
			for _, customerConversationArg := range customerConversationArgsTypeMessage {
				mapExternalID[customerConversationArg.ExternalID] = true
			}
			for externalID := range mapExternalID {
				externalIDs = append(externalIDs, externalID)
			}

			customerConversations, err := a.fbCustomerConversationStore(ctx).ExternalIDs(externalIDs).ListFbCustomerConversations()
			if err != nil {
				return nil, err
			}
			for _, customerConversation := range customerConversations {
				mapOldCustomerConversation[customerConversation.ExternalID] = customerConversation
			}
		}

		for _, customerConversation := range customerConversationArgsTypeMessage {
			ID := customerConversation.ID
			var isRead bool
			createdAt, updatedAt := time.Now(), time.Now()
			if oldCustomerConversation, ok := mapOldCustomerConversation[customerConversation.ExternalID]; ok {
				ID = oldCustomerConversation.ID
				isRead = oldCustomerConversation.IsRead
				createdAt = oldCustomerConversation.CreatedAt
				updatedAt = oldCustomerConversation.UpdatedAt
			}
			newCustomerConversations = append(newCustomerConversations, &fbmessaging.FbCustomerConversation{
				ID:               ID,
				ExternalPageID:   customerConversation.ExternalPageID,
				ExternalID:       customerConversation.ExternalUserID,
				ExternalUserID:   customerConversation.ExternalUserID,
				ExternalUserName: customerConversation.ExternalUserName,
				IsRead:           isRead,
				Type:             fb_customer_conversation_type.Message,
				LastMessage:      customerConversation.LastMessage,
				LastMessageAt:    customerConversation.LastMessageAt,
				CreatedAt:        createdAt,
				UpdatedAt:        updatedAt,
			})
		}
	}

	// Handle type comment
	if len(customerConversationArgsTypeComment) > 0 {
		mapExternalIDandMapExternalUserIDAndOldCustomerConversation := make(map[string]map[string]*fbmessaging.FbCustomerConversation)
		{
			mapExternalIDandExternalUserID := make(map[string]map[string]bool)
			for _, customerConversationArg := range customerConversationArgsTypeComment {
				if _, ok := mapExternalIDandExternalUserID[customerConversationArg.ExternalID]; !ok {
					mapExternalIDandExternalUserID[customerConversationArg.ExternalID] = make(map[string]bool)
				}
				mapExternalIDandExternalUserID[customerConversationArg.ExternalID][customerConversationArg.ExternalUserID] = true
			}

			query := a.fbCustomerConversationStore(ctx)
			for externalID, mapExternalUserID := range mapExternalIDandExternalUserID {
				for externalUserID := range mapExternalUserID {
					query = query.ExternalIDAndExternalUserID(externalID, externalUserID)
				}
			}
			customerConversations, err := query.ListFbCustomerConversations()
			if err != nil {
				return nil, err
			}
			for _, customerConversation := range customerConversations {
				if _, ok := mapExternalIDandMapExternalUserIDAndOldCustomerConversation[customerConversation.ExternalID]; !ok {
					mapExternalIDandMapExternalUserIDAndOldCustomerConversation[customerConversation.ExternalID] = make(map[string]*fbmessaging.FbCustomerConversation)
				}
				mapExternalIDandMapExternalUserIDAndOldCustomerConversation[customerConversation.ExternalID][customerConversation.ExternalUserID] = customerConversation
			}
		}

		for _, customerConversation := range customerConversationArgsTypeComment {
			ID := customerConversation.ID
			var isRead bool
			createdAt, updatedAt := time.Now(), time.Now()
			if oldCustomerConversation, ok := mapExternalIDandMapExternalUserIDAndOldCustomerConversation[customerConversation.ExternalID][customerConversation.ExternalUserID]; ok {
				ID = oldCustomerConversation.ID
				isRead = oldCustomerConversation.IsRead
				createdAt = oldCustomerConversation.CreatedAt
				updatedAt = oldCustomerConversation.UpdatedAt
			}
			newCustomerConversations = append(newCustomerConversations, &fbmessaging.FbCustomerConversation{
				ID:                        ID,
				ExternalPageID:            customerConversation.ExternalUserID,
				ExternalID:                customerConversation.ExternalID,
				ExternalUserID:            customerConversation.ExternalUserID,
				ExternalUserName:          customerConversation.ExternalUserName,
				Type:                      fb_customer_conversation_type.Comment,
				IsRead:                    isRead,
				ExternalPostAttachments:   customerConversation.ExternalPostAttachments,
				ExternalCommentAttachment: customerConversation.ExternalCommentAttachment,
				LastMessage:               customerConversation.LastMessage,
				LastMessageAt:             customerConversation.LastMessageAt,
				CreatedAt:                 createdAt,
				UpdatedAt:                 updatedAt,
			})
		}
	}

	if len(newCustomerConversations) > 0 {
		if err := a.fbCustomerConversationStore(ctx).CreateFbCustomerConversations(newCustomerConversations); err != nil {
			return nil, err
		}
	}

	return newCustomerConversations, nil

}

func (a *FbExternalMessagingAggregate) CreateFbExternalPosts(
	ctx context.Context, args *fbmessaging.CreateFbExternalPostsArgs,
) ([]*fbmessaging.FbExternalPost, error) {
	newFbExternalPosts := make([]*fbmessaging.FbExternalPost, 0, len(args.FbExternalPosts))
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalPost := range args.FbExternalPosts {
			newFbExternalPost := new(fbmessaging.FbExternalPost)
			if err := scheme.Convert(fbExternalPost, newFbExternalPost); err != nil {
				return err
			}

			if fbExternalPost.StatusType == fb_status_type.AddedVideo {
				newFbExternalPost.IsLiveVideo = true
			}

			newFbExternalPosts = append(newFbExternalPosts, newFbExternalPost)
		}

		if len(newFbExternalPosts) > 0 {
			if err := a.fbExternalPostStore(ctx).UpsertFbExternalPosts(newFbExternalPosts); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return nil, err
	}
	return newFbExternalPosts, nil
}

func (a *FbExternalMessagingAggregate) UpdateOrCreateFbExternalPostsFromSync(
	ctx context.Context, args *fbmessaging.UpdateOrCreateFbExternalPostsFromSyncArgs,
) ([]*fbmessaging.FbExternalPost, error) {
	fbExternalPostIDs := make([]string, 0, len(args.FbExternalPosts))
	for _, fbExternalPost := range args.FbExternalPosts {
		fbExternalPostIDs = append(fbExternalPostIDs, fbExternalPost.ExternalID)
	}

	oldFbExternalPosts, err := a.fbExternalPostStore(ctx).ExternalIDs(fbExternalPostIDs).ListFbExternalPosts()
	if err != nil {
		return nil, err
	}
	mapOldFbExternalPost := make(map[string]*fbmessaging.FbExternalPost)
	for _, oldFbExternalPost := range oldFbExternalPosts {
		mapOldFbExternalPost[oldFbExternalPost.ExternalID] = oldFbExternalPost
	}

	var resultFbExternalPosts []*fbmessaging.FbExternalPost
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalPostArg := range args.FbExternalPosts {
			newFbExternalPost := new(fbmessaging.FbExternalPost)
			if err := scheme.Convert(fbExternalPostArg, newFbExternalPost); err != nil {
				return err
			}
			if fbExternalPostArg.StatusType == fb_status_type.AddedVideo {
				newFbExternalPost.IsLiveVideo = true
			}

			if oldFbExternalPost, ok := mapOldFbExternalPost[fbExternalPostArg.ExternalID]; !ok {
				if err := a.fbExternalPostStore(ctx).CreateFbExternalPost(newFbExternalPost); err != nil {
					return err
				}
			} else {
				newFbExternalPost.ID = oldFbExternalPost.ID
				if err := a.fbExternalPostStore(ctx).ExternalID(fbExternalPostArg.ExternalID).UpdateFbExternalPost(newFbExternalPost); err != nil {
					return err
				}
			}

			resultFbExternalPosts = append(resultFbExternalPosts, newFbExternalPost)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	return resultFbExternalPosts, nil
}

func (a *FbExternalMessagingAggregate) UpdateLiveVideoStatusFromSync(
	ctx context.Context, args *fbmessaging.UpdateLiveVideoStatusFromSyncArgs,
) (*fbmessaging.FbExternalPost, error) {
	updatedFbExternalPost := new(fbmessaging.FbExternalPost)
	if err := scheme.Convert(args, updatedFbExternalPost); err != nil {
		return nil, err
	}
	updatedFbExternalPost.IsLiveVideo = true

	if err := a.fbExternalPostStore(ctx).ExternalID(args.ExternalID).UpdateFbExternalPost(updatedFbExternalPost); err != nil {
		return nil, err
	}

	// update type of fbCustomerConversation to live_video
	if _, err := a.fbCustomerConversationStore(ctx).ExternalID(args.ExternalID).UpdateType(fb_customer_conversation_type.LiveVideo); err != nil {
		return nil, err
	}

	return a.fbExternalPostStore(ctx).ExternalID(args.ExternalID).GetFbExternalPost()
}

func (a *FbExternalMessagingAggregate) CreateOrUpdateFbExternalComments(
	ctx context.Context, args *fbmessaging.CreateOrUpdateFbExternalCommentsArgs,
) ([]*fbmessaging.FbExternalComment, error) {
	fbExternalCommentIDs := make([]string, 0, len(args.FbExternalComments))
	for _, fbExternalComment := range args.FbExternalComments {
		fbExternalCommentIDs = append(fbExternalCommentIDs, fbExternalComment.ExternalID)
	}

	oldFbExternalComments, err := a.fbExternalCommentStore(ctx).ExternalIDs(fbExternalCommentIDs).ListFbExternalComments()
	if err != nil {
		return nil, err
	}
	mapOldFbExternalComment := make(map[string]*fbmessaging.FbExternalComment)
	for _, oldFbExternalComment := range oldFbExternalComments {
		mapOldFbExternalComment[oldFbExternalComment.ExternalID] = oldFbExternalComment
	}

	var resultFbExternalComments []*fbmessaging.FbExternalComment
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalCommentArg := range args.FbExternalComments {
			newFbExternalComment := new(fbmessaging.FbExternalComment)
			if err := scheme.Convert(fbExternalCommentArg, newFbExternalComment); err != nil {
				return err
			}

			if oldFbExternalComment, ok := mapOldFbExternalComment[fbExternalCommentArg.ExternalID]; ok {
				newFbExternalComment.ID = oldFbExternalComment.ID
				newFbExternalComment.CreatedAt = oldFbExternalComment.CreatedAt
				resultFbExternalComments = append(resultFbExternalComments, oldFbExternalComment)

				if isEqual := compare.CompareFbExternalComments(oldFbExternalComment, newFbExternalComment); isEqual {
					continue
				}
			} else {
				resultFbExternalComments = append(resultFbExternalComments, newFbExternalComment)
			}

			if err := a.fbExternalCommentStore(ctx).CreateOrUpdateFbExternalComment(newFbExternalComment); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalCommentCreatedOrUpdatedEvent{
				FbExternalComment: newFbExternalComment,
			}
			if err := a.eventBus.Publish(ctx, event); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return resultFbExternalComments, nil
}

func (a *FbExternalMessagingAggregate) UpdateIsReadCustomerConversation(
	ctx context.Context, conversationCustomerID dot.ID, isRead bool,
) (int, error) {
	return a.fbCustomerConversationStateStore(ctx).ID(conversationCustomerID).UpdateIsRead(isRead)
}

func (a *FbExternalMessagingAggregate) CreateFbExternalPost(ctx context.Context, args *fbmessaging.FbCreatePostArgs) (*fbmessaging.FbExternalPost, error) {
	args.Message = strings.TrimSpace(args.Message)
	if args.Message == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing post message content")
	}

	createPostRequest := &fbclientmodel.CreatePostRequest{
		Message: args.Message,
	}
	post, err := a.fbClient.CallAPICreatePost(&fbclient.CreatePostRequest{
		AccessToken: args.AccessToken,
		PageID:      args.ExternalPageID,
		Content:     createPostRequest,
	})
	if err != nil {
		return nil, err
	}

	return &fbmessaging.FbExternalPost{
		ExternalID: post.ID,
	}, nil
}

// Temp use
func (a *FbExternalMessagingAggregate) SaveFbExternalPost(
	ctx context.Context, post *fbmessaging.FbSavePostArgs,
) (*fbmessaging.FbExternalPost, error) {
	extPost := &fbmessaging.FbExternalPost{
		ID:                  cm.NewID(),
		ExternalPageID:      post.ExternalPageID,
		ExternalID:          post.ExternalID,
		ExternalFrom:        post.ExternalFrom,
		ExternalPicture:     post.ExternalPicture,
		ExternalIcon:        post.ExternalIcon,
		ExternalMessage:     post.ExternalMessage,
		ExternalAttachments: post.ExternalAttachments,
		ExternalCreatedTime: post.ExternalCreatedTime,
		ExternalParentID:    post.ExternalParentID,
		FeedType:            post.FeedType,
	}
	if err := a.fbExternalPostStore(ctx).UpsertFbExternalPost(extPost); err != nil {
		return nil, err
	}
	return extPost, nil
}

func (a *FbExternalMessagingAggregate) UpdateFbPostMessageAndPicture(
	ctx context.Context, feedMessage *fbmessaging.FbUpdatePostMessageArgs,
) error {
	return a.fbExternalPostStore(ctx).ExternalID(feedMessage.ExternalPostID).UpdatePostMessageAndPicture(feedMessage.Message, feedMessage.ExternalPicture)
}

func (a *FbExternalMessagingAggregate) UpdateFbCommentMessage(
	ctx context.Context, updateArgs *fbmessaging.FbUpdateCommentMessageArgs,
) (int, error) {
	return a.fbExternalCommentStore(ctx).ExternalID(updateArgs.ExternalCommentID).UpdateMessage(updateArgs.Message)
}

func (a *FbExternalMessagingAggregate) RemovePost(
	ctx context.Context, removeArgs *fbmessaging.RemovePostArgs,
) error {
	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		_, err := a.fbExternalPostStore(ctx).ExternalID(removeArgs.ExternalPostID).SoftDelete()
		if err != nil {
			// If post not found, create new post with deleted_at
			if cm.ErrorCode(err) == cm.NotFound {
				extPost := &fbmessaging.FbExternalPost{
					ID:             cm.NewID(),
					ExternalPageID: removeArgs.ExternalPageID,
					ExternalID:     removeArgs.ExternalPostID,
					DeletedAt:      time.Now(),
				}
				return a.fbExternalPostStore(ctx).UpsertFbExternalPost(extPost)
			}
			return err
		}

		// Remove all child posts that belong to.
		if _, err := a.fbExternalPostStore(ctx).
			ExternalParentID(removeArgs.ExternalPostID).
			SoftDelete(); err != nil {
			return err
		}

		// Remove all comments that belong to.
		if _, err := a.fbExternalCommentStore(ctx).
			ExternalPostID(removeArgs.ExternalPostID).
			SoftDelete(); err != nil {
			return err
		}

		// Remove all customer conversations that belong to.
		if _, err := a.fbCustomerConversationStore(ctx).
			ExternalID(removeArgs.ExternalPostID).
			SoftDelete(); err != nil {
			return err
		}

		return nil
	})
}

func (a *FbExternalMessagingAggregate) RemoveComment(ctx context.Context, removeArgs *fbmessaging.RemoveCommentArgs) error {
	return a.db.InTransaction(ctx, func(queryInterface cmsql.QueryInterface) error {
		commentID := removeArgs.ExternalCommentID
		if _, err := a.fbExternalCommentStore(ctx).
			ExternalIDOrExternalParentID(commentID, commentID).
			SoftDelete(); err != nil {
			return err
		}
		return nil
	})
}

func (a *FbExternalMessagingAggregate) LikeOrUnLikeComment(
	ctx context.Context, args *fbmessaging.LikeOrUnLikeCommentArgs,
) error {
	if _, err := a.fbExternalCommentStore(ctx).ExternalID(args.ExternalCommentID).GetFbExternalComment(); err != nil {
		return err
	}

	_, err := a.fbExternalCommentStore(ctx).ExternalID(args.ExternalCommentID).UpdateIsLiked(args.IsLiked)
	return err
}

func (a *FbExternalMessagingAggregate) HideOrUnHideComment(
	ctx context.Context, args *fbmessaging.HideOrUnHideCommentArgs,
) error {
	if _, err := a.fbExternalCommentStore(ctx).ExternalID(args.ExternalCommentID).GetFbExternalComment(); err != nil {
		return err
	}

	_, err := a.fbExternalCommentStore(ctx).ExternalID(args.ExternalCommentID).UpdateIsHidden(args.IsHidden)
	return err
}

func (a *FbExternalMessagingAggregate) UpdateIsPrivateRepliedComment(
	ctx context.Context, args *fbmessaging.UpdateIsPrivateRepliedCommentArgs,
) error {
	if _, err := a.fbExternalCommentStore(ctx).ExternalID(args.ExternalCommentID).GetFbExternalComment(); err != nil {
		return err
	}

	if _, err := a.fbExternalCommentStore(ctx).ExternalID(args.ExternalCommentID).UpdateIsPrivateReplied(args.IsPrivateReplied); err != nil {
		return err
	}
	return nil
}
