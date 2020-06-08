package fbmessaging

import (
	"context"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/backend/com/fabo/main/compare"
	"o.o/backend/com/fabo/main/fbmessaging/convert"
	"o.o/backend/com/fabo/main/fbmessaging/sqlstore"
	com "o.o/backend/com/main"
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
}

func NewFbExternalMessagingAggregate(db com.MainDB, eventBus capi.EventBus) *FbExternalMessagingAggregate {
	return &FbExternalMessagingAggregate{
		db:                               db,
		eventBus:                         eventBus,
		fbExternalPostStore:              sqlstore.NewFbExternalPostStore(db),
		fbExternalCommentStore:           sqlstore.NewFbExternalCommentStore(db),
		fbExternalConversationStore:      sqlstore.NewFbExternalConversationStore(db),
		fbExternalMessageStore:           sqlstore.NewFbExternalMessageStore(db),
		fbCustomerConversationStore:      sqlstore.NewFbCustomerConversationStore(db),
		fbCustomerConversationStateStore: sqlstore.NewFbCustomerConversationStateStore(db),
	}
}

func FbExternalMessagingAggregateMessageBus(a *FbExternalMessagingAggregate) fbmessaging.CommandBus {
	b := bus.New()
	return fbmessaging.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *FbExternalMessagingAggregate) CreateFbExternalMessages(
	ctx context.Context, args *fbmessaging.CreateFbExternalMessagesArgs,
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
				resultFbExternalMessages = append(resultFbExternalMessages, oldFbExternalMessage)

				if isEqual := compare.Compare(oldFbExternalMessage, newFbExternalMessage); isEqual {
					continue
				}
			} else {
				resultFbExternalMessages = append(resultFbExternalMessages, newFbExternalMessage)
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
	return resultFbExternalMessages, nil
}

func (a *FbExternalMessagingAggregate) CreateFbExternalConversations(
	ctx context.Context, args *fbmessaging.CreateFbExternalConversationsArgs,
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
	return newFbExternalConversations, nil
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

			newFbExternalPosts = append(newFbExternalPosts, newFbExternalPost)
		}

		if len(newFbExternalPosts) > 0 {
			if err := a.fbExternalPostStore(ctx).CreateFbExternalPosts(newFbExternalPosts); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return nil, err
	}
	return newFbExternalPosts, nil
}

func (a *FbExternalMessagingAggregate) CreateOrUpdateFbExternalPosts(
	ctx context.Context, args *fbmessaging.CreateOrUpdateFbExternalPostsArgs,
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

	resultFbExternalPosts := make([]*fbmessaging.FbExternalPost, 0, len(args.FbExternalPosts))
	newFbExternalPosts := make([]*fbmessaging.FbExternalPost, 0, len(args.FbExternalPosts))
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalPostArg := range args.FbExternalPosts {

			newFbExternalPost := new(fbmessaging.FbExternalPost)
			if err := scheme.Convert(fbExternalPostArg, newFbExternalPost); err != nil {
				return err
			}

			if oldFbExternalPost, ok := mapOldFbExternalPost[fbExternalPostArg.ExternalID]; ok {
				newFbExternalPost.ID = oldFbExternalPost.ID
				resultFbExternalPosts = append(resultFbExternalPosts, oldFbExternalPost)

				if isEqual := compare.Compare(oldFbExternalPost, newFbExternalPost); isEqual {
					continue
				}
			} else {
				resultFbExternalPosts = append(resultFbExternalPosts, newFbExternalPost)
			}

			newFbExternalPosts = append(newFbExternalPosts, newFbExternalPost)
		}

		if len(newFbExternalPosts) > 0 {
			if err := a.fbExternalPostStore(ctx).CreateFbExternalPosts(newFbExternalPosts); err != nil {
				return err
			}
		}

		return nil

	}); err != nil {
		return nil, err
	}
	return resultFbExternalPosts, nil
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
	var newFbExternalComments []*fbmessaging.FbExternalComment
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, fbExternalCommentArg := range args.FbExternalComments {
			newFbExternalComment := new(fbmessaging.FbExternalComment)
			if err := scheme.Convert(fbExternalCommentArg, newFbExternalComment); err != nil {
				return err
			}

			if oldFbExternalComment, ok := mapOldFbExternalComment[fbExternalCommentArg.ExternalID]; ok {
				newFbExternalComment.ID = oldFbExternalComment.ID
				resultFbExternalComments = append(resultFbExternalComments, oldFbExternalComment)

				if isEqual := compare.Compare(newFbExternalComment, oldFbExternalComment); isEqual {
					continue
				}
			} else {
				resultFbExternalComments = append(resultFbExternalComments, newFbExternalComment)
			}

			newFbExternalComments = append(newFbExternalComments, newFbExternalComment)
		}

		if len(newFbExternalComments) > 0 {
			if err := a.fbExternalCommentStore(ctx).CreateFbExternalComments(newFbExternalComments); err != nil {
				return err
			}

			event := &fbmessaging.FbExternalCommentsCreatedEvent{
				FbExternalComments: newFbExternalComments,
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
