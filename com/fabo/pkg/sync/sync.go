package sync

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/convert"
	fbpagemodel "o.o/backend/com/fabo/main/fbpage/model"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/scheduler"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var ll = l.New()

const (
	defaultNumWorkers        = 32
	defaultRecurrentFacebook = 1 * time.Minute
	defaultRecurrent         = 5 * time.Minute
	defaultErrRecurr         = 10 * time.Minute
	maxLogsCount             = 500
)

type TaskActionType int

const (
	GetPosts             TaskActionType = 737
	GetConversations     TaskActionType = 363
	GetMessages          TaskActionType = 160
	GetCommentsByPostIDs TaskActionType = 655
)

type GetCommentsByPostIDsArguments struct {
	IDs []string
}

type GetMessagesArguments struct {
	ConversationID         dot.ID
	ExternalConversationID string
}

type TaskArguments struct {
	actionType     TaskActionType
	accessToken    string
	shopID         dot.ID
	externalPageID string
	pageID         dot.ID

	GetCommentsByPostIDsArgs *GetCommentsByPostIDsArguments
	GetMessages              *GetMessagesArguments

	fbPagingRequest *model.FacebookPagingRequest
}

func (t *TaskArguments) Clone() *TaskArguments {
	return &TaskArguments{
		accessToken:    t.accessToken,
		shopID:         t.shopID,
		externalPageID: t.externalPageID,
		pageID:         t.pageID,
	}
}

type Synchronizer struct {
	scheduler *scheduler.Scheduler

	db               *cmsql.Database
	fbClient         *fbclient.FbClient
	mapTaskArguments map[dot.ID]*TaskArguments

	fbMessagingAggr  fbmessaging.CommandBus
	fbMessagingQuery fbmessaging.QueryBus

	timeLimit int

	mu sync.Mutex
}

func New(
	db *cmsql.Database, fbClient *fbclient.FbClient,
	fbMessagingAggr fbmessaging.CommandBus, fbMessagingQuery fbmessaging.QueryBus,
	timeLimit int,
) *Synchronizer {
	sched := scheduler.New(defaultNumWorkers)
	s := &Synchronizer{
		scheduler:        sched,
		db:               db,
		fbClient:         fbClient,
		mapTaskArguments: make(map[dot.ID]*TaskArguments),
		fbMessagingAggr:  fbMessagingAggr,
		fbMessagingQuery: fbMessagingQuery,
		timeLimit:        timeLimit,
	}
	return s
}

func (s *Synchronizer) Init() error {
	s.scheduler.AddAfter(cm.NewID(), 1*time.Second, s.addJobs)
	return nil
}

func (s *Synchronizer) addJobs(id interface{}, p scheduler.Planner) (_err error) {
	fbPageCombineds, err := listAllFbPagesActive(s.db)
	if err != nil {
		return err
	}

	for _, fbPageCombined := range fbPageCombineds {
		// task get post of specific page
		{
			s.mu.Lock()
			taskID := cm.NewID()
			s.mapTaskArguments[taskID] = &TaskArguments{
				actionType:     GetPosts,
				accessToken:    fbPageCombined.FbExternalPageInternal.Token,
				shopID:         fbPageCombined.FbExternalPage.ShopID,
				pageID:         fbPageCombined.FbExternalPage.ID,
				externalPageID: fbPageCombined.FbExternalPage.ExternalID,
				fbPagingRequest: &model.FacebookPagingRequest{
					Limit: dot.Int(fbclient.DefaultLimitGetPosts),
					TimePagination: &model.TimePaginationRequest{
						Since: time.Now().AddDate(0, 0, -s.timeLimit),
					},
				},
			}
			t := rand.Intn(int(time.Second))
			s.scheduler.AddAfter(taskID, time.Duration(t), s.syncCallbackLogs)
			s.mu.Unlock()
		}

		// task get conversation
		{
			s.mu.Lock()
			taskID := cm.NewID()
			s.mapTaskArguments[taskID] = &TaskArguments{
				actionType:     GetConversations,
				accessToken:    fbPageCombined.FbExternalPageInternal.Token,
				shopID:         fbPageCombined.FbExternalPage.ShopID,
				pageID:         fbPageCombined.FbExternalPage.ID,
				externalPageID: fbPageCombined.FbExternalPage.ExternalID,
				fbPagingRequest: &model.FacebookPagingRequest{
					Limit: dot.Int(fbclient.DefaultLimitGetConversations),
				},
			}
			t := rand.Intn(int(time.Second))
			s.scheduler.AddAfter(taskID, time.Duration(t), s.syncCallbackLogs)
			s.mu.Unlock()
		}
	}

	s.scheduler.AddAfter(cm.NewID(), 30*time.Second, s.addJobs)

	return
}

func (s *Synchronizer) Start() {
	s.scheduler.Start()
}

func (s *Synchronizer) Stop() {
	s.scheduler.Stop()
}

func (s *Synchronizer) syncCallbackLogs(id interface{}, p scheduler.Planner) (_err error) {
	taskID := id.(dot.ID)
	ctx := bus.Ctx()
	s.mu.Lock()
	taskArgs := s.mapTaskArguments[taskID]
	s.mu.Unlock()

	accessToken := taskArgs.accessToken
	shopID := taskArgs.shopID
	pageID := taskArgs.pageID
	externalPageID := taskArgs.externalPageID
	fbPagingReq := taskArgs.fbPagingRequest

	switch taskArgs.actionType {
	case GetPosts:
		fmt.Println("GetPosts")
		//fbPostsResp, err := s.fbClient.CallAPIListPublishedPosts(accessToken, externalPageID, fbPagingReq)
		//if err != nil {
		//	return err
		//}
		//
		//if len(fbPostsResp.Data) == 0 ||
		//	fbPostsResp.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		//	return nil
		//}

	case GetConversations:
		fmt.Println("GetConversations")
		fbConversationsResp, err := s.fbClient.CallAPIListConversations(accessToken, externalPageID, fbPagingReq)
		if err != nil {
			// TODO: Ngoc classify error type
			return err
		}

		if len(fbConversationsResp.Conversations.ConversationsData) == 0 ||
			fbConversationsResp.Conversations.Paging.CompareFacebookPagingRequest(fbPagingReq) {
			return nil
		}

		isFinished := false
		mapFbExternalConversation := make(map[string]*model.Conversation)
		var fbConversationIDs []string
		for _, fbConversation := range fbConversationsResp.Conversations.ConversationsData {
			if time.Now().Sub(fbConversation.UpdatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
				isFinished = true
				continue
			}
			mapFbExternalConversation[fbConversation.ID] = fbConversation
			fbConversationIDs = append(fbConversationIDs, fbConversation.ID)
		}

		listFbConversationsQuery := &fbmessaging.ListFbExternalConversationsByExternalIDsQuery{
			ExternalIDs: fbConversationIDs,
		}
		if err := s.fbMessagingQuery.Dispatch(ctx, listFbConversationsQuery); err != nil {
			return err
		}

		mapOldFbExternalConversation := make(map[string]*fbmessaging.FbExternalConversation)
		for _, oldFbConversation := range listFbConversationsQuery.Result {
			mapOldFbExternalConversation[oldFbConversation.ExternalID] = oldFbConversation
		}

		mapExternalIDAndID := make(map[string]dot.ID)
		var fbExternalConversationsArgs []*fbmessaging.CreateFbExternalConversationArgs
		for externalID, fbExternalConversation := range mapFbExternalConversation {
			if oldFbExternalConversation, ok := mapOldFbExternalConversation[externalID]; ok {
				mapExternalIDAndID[oldFbExternalConversation.ExternalID] = oldFbExternalConversation.ID
				fbExternalConversationsArgs = append(fbExternalConversationsArgs, &fbmessaging.CreateFbExternalConversationArgs{
					ID:                   oldFbExternalConversation.ID,
					FbPageID:             oldFbExternalConversation.FbPageID,
					ExternalID:           oldFbExternalConversation.ExternalID,
					ExternalUserID:       oldFbExternalConversation.ExternalUserID,
					ExternalUserName:     oldFbExternalConversation.ExternalUserName,
					ExternalLink:         oldFbExternalConversation.ExternalLink,
					ExternalUpdatedTime:  fbExternalConversation.UpdatedTime.ToTime(),
					ExternalMessageCount: fbExternalConversation.MessageCount,
					LastMessage:          oldFbExternalConversation.LastMessage,
					LastMessageAt:        oldFbExternalConversation.LastMessageAt,
				})
			} else {
				var externalUserID, externalUserName string
				for _, sender := range fbExternalConversation.Senders.Data {
					if sender.ID != externalPageID {
						externalUserID = sender.ID
						externalUserName = sender.Name
						break
					}
				}
				ID := cm.NewID()
				mapExternalIDAndID[fbExternalConversation.ID] = ID
				fbExternalConversationsArgs = append(fbExternalConversationsArgs, &fbmessaging.CreateFbExternalConversationArgs{
					ID:                   ID,
					FbPageID:             pageID,
					ExternalID:           fbExternalConversation.ID,
					ExternalUserID:       externalUserID,
					ExternalUserName:     externalUserName,
					ExternalUpdatedTime:  fbExternalConversation.UpdatedTime.ToTime(),
					ExternalLink:         fbExternalConversation.Link,
					ExternalMessageCount: fbExternalConversation.MessageCount,
				})
			}
		}

		if len(fbExternalConversationsArgs) > 0 {
			if err := s.fbMessagingAggr.Dispatch(ctx, &fbmessaging.CreateFbExternalConversationsCommand{
				FbExternalConversations: fbExternalConversationsArgs,
			}); err != nil {
				return err
			}
		}

		listFbCustomerConversationsQuery := &fbmessaging.ListFbCustomerConversationsByExternalIDsQuery{
			ExternalIDs: fbConversationIDs,
		}
		if err := s.fbMessagingQuery.Dispatch(ctx, listFbCustomerConversationsQuery); err != nil {
			return err
		}

		mapOldFbCustomerConversation := make(map[string]*fbmessaging.FbCustomerConversation)
		for _, oldFbCustomerConversation := range listFbCustomerConversationsQuery.Result {
			mapOldFbCustomerConversation[oldFbCustomerConversation.ExternalID] = oldFbCustomerConversation
		}

		var fbCustomerConversationsArgs []*fbmessaging.CreateFbCustomerConversationArgs
		for externalID, fbExternalConversation := range mapFbExternalConversation {
			if _, ok := mapOldFbCustomerConversation[externalID]; !ok {
				var externalUserID, externalUserName string
				for _, sender := range fbExternalConversation.Senders.Data {
					if sender.ID != externalPageID {
						externalUserID = sender.ID
						externalUserName = sender.Name
						break
					}
				}
				ID := cm.NewID()
				mapExternalIDAndID[fbExternalConversation.ID] = ID
				// TODO: Ngoc fbclient get latest message (limit 1)
				fbCustomerConversationsArgs = append(fbCustomerConversationsArgs, &fbmessaging.CreateFbCustomerConversationArgs{
					ID:               ID,
					FbPageID:         pageID,
					ExternalID:       fbExternalConversation.ID,
					ExternalUserID:   externalUserID,
					ExternalUserName: externalUserName,
					IsRead:           false,
					Type:             fb_customer_conversation_type.Message,
				})
			}
		}

		if len(fbCustomerConversationsArgs) > 0 {
			if err := s.fbMessagingAggr.Dispatch(ctx, &fbmessaging.CreateFbCustomerConversationsCommand{
				FbCustomerConversations: fbCustomerConversationsArgs,
			}); err != nil {
				return err
			}
		}

		for _, fbConversation := range fbConversationsResp.Conversations.ConversationsData {
			if time.Now().Sub(fbConversation.UpdatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
				continue
			}
			s.mu.Lock()
			newTaskID := cm.NewID()
			s.mapTaskArguments[newTaskID] = &TaskArguments{
				actionType:  GetMessages,
				accessToken: accessToken,
				shopID:      shopID,
				pageID:      pageID,
				GetMessages: &GetMessagesArguments{
					ConversationID:         mapExternalIDAndID[fbConversation.ID],
					ExternalConversationID: fbConversation.ID,
				},
				externalPageID: externalPageID,
			}
			t := rand.Intn(int(time.Second))
			s.scheduler.AddAfter(newTaskID, time.Duration(t), s.syncCallbackLogs)
			s.mu.Unlock()
		}

		if !isFinished {
			s.mu.Lock()
			newTaskID := cm.NewID()
			t := rand.Intn(int(time.Second))
			s.mapTaskArguments[newTaskID] = &TaskArguments{
				actionType:      GetConversations,
				accessToken:     accessToken,
				shopID:          shopID,
				externalPageID:  externalPageID,
				pageID:          pageID,
				fbPagingRequest: fbConversationsResp.Conversations.Paging.ToPagingRequestAfter(fbclient.DefaultLimitGetConversations),
			}
			s.scheduler.AddAfter(newTaskID, time.Duration(t), s.syncCallbackLogs)
			s.mu.Unlock()
		}
	case GetMessages:
		fmt.Println("GetMessages")
		conversationID := taskArgs.GetMessages.ConversationID
		externalConversationID := taskArgs.GetMessages.ExternalConversationID
		fbMessagesResp, err := s.fbClient.CallAPIListMessages(accessToken, externalConversationID, fbPagingReq)
		if err != nil {
			// TODO: Ngoc classify error type
			return err
		}

		if len(fbMessagesResp.Messages.MessagesData) == 0 ||
			fbMessagesResp.Messages.Paging.CompareFacebookPagingRequest(fbPagingReq) {
			return nil
		}

		mapFbExternalMessage := make(map[string]*model.MessageData)
		var fbMessageIDs []string
		for _, fbMessage := range fbMessagesResp.Messages.MessagesData {
			if time.Now().Sub(fbMessage.CreatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
				continue
			}
			mapFbExternalMessage[fbMessage.ID] = fbMessage
			fbMessageIDs = append(fbMessageIDs, fbMessage.ID)
		}

		listFbMessagesQuery := &fbmessaging.ListFbExternalMessagesByExternalIDsQuery{
			ExternalIDs: fbMessageIDs,
		}
		if err := s.fbMessagingQuery.Dispatch(ctx, listFbMessagesQuery); err != nil {
			return err
		}

		mapOldFbExternalMessage := make(map[string]*fbmessaging.FbExternalMessage)
		for _, oldFbMessage := range listFbMessagesQuery.Result {
			mapOldFbExternalMessage[oldFbMessage.ExternalID] = oldFbMessage
		}

		var fbExternalMessagesArgs []*fbmessaging.CreateFbExternalMessageArgs
		for externalID, fbExternalMessage := range mapFbExternalMessage {
			var externalAttachments []*fbmessaging.FbMessageAttachment
			if fbExternalMessage.Attachments != nil {
				externalAttachments = fbclientconvert.ConvertMessageDataAttachments(fbExternalMessage.Attachments.Data)
			}
			if oldFbExternalMessage, ok := mapOldFbExternalMessage[externalID]; ok {
				fbExternalMessagesArgs = append(fbExternalMessagesArgs, &fbmessaging.CreateFbExternalMessageArgs{
					ID:                     oldFbExternalMessage.ID,
					FbConversationID:       oldFbExternalMessage.FbConversationID,
					ExternalConversationID: oldFbExternalMessage.ExternalConversationID,
					FbPageID:               oldFbExternalMessage.FbPageID,
					ExternalID:             oldFbExternalMessage.ExternalID,
					ExternalMessage:        fbExternalMessage.Message,
					ExternalTo:             fbclientconvert.ConvertObjectsTo(fbExternalMessage.To),
					ExternalFrom:           fbclientconvert.ConvertObjectFrom(fbExternalMessage.From),
					ExternalAttachments:    externalAttachments,
					ExternalCreatedTime:    oldFbExternalMessage.ExternalCreatedTime,
				})
			} else {
				fbExternalMessagesArgs = append(fbExternalMessagesArgs, &fbmessaging.CreateFbExternalMessageArgs{
					ID:                     cm.NewID(),
					FbConversationID:       conversationID,
					ExternalConversationID: externalConversationID,
					FbPageID:               pageID,
					ExternalID:             fbExternalMessage.ID,
					ExternalMessage:        fbExternalMessage.Message,
					ExternalTo:             fbclientconvert.ConvertObjectsTo(fbExternalMessage.To),
					ExternalFrom:           fbclientconvert.ConvertObjectFrom(fbExternalMessage.From),
					ExternalAttachments:    externalAttachments,
					ExternalCreatedTime:    fbExternalMessage.CreatedTime.ToTime(),
				})
			}
		}

		if len(fbExternalMessagesArgs) > 0 {
			if err := s.fbMessagingAggr.Dispatch(ctx, &fbmessaging.CreateFbExternalMessagesCommand{
				FbExternalMessages: fbExternalMessagesArgs,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

//func (s *Synchronizer) syncCallbackLogs(id interface{}, p scheduler.Planner) (_err error) {
//	taskID := id.(dot.ID)
//	s.mu.Lock()
//	taskArgs := s.mapTaskArguments[taskID]
//	s.mu.Unlock()
//
//	accessToken := taskArgs.accessToken
//	shopID := taskArgs.shopID
//	pageID := taskArgs.pageID
//	externalPageID := taskArgs.externalPageID
//	fbPagingReq := taskArgs.fbPagingRequest
//
//	switch taskArgs.actionType {
//	case GetPosts:
//		fmt.Println("GetPosts")
//		fbPostsResponse, err := s.fbClient.CallAPIListPublishedPosts(accessToken, externalPageID, fbPagingReq)
//		if err != nil {
//			return err
//		}
//
//		if len(fbPostsResponse.Data) == 0 {
//			return nil
//		}
//
//		var childPostIDs, postIDs []string
//		// TODO: Ngoc save posts to database
//		// Get comments and sub-posts
//		for _, fbPost := range fbPostsResponse.Data {
//			postIDs = append(postIDs, fbPost.ID)
//			if fbPost.Attachments != nil {
//				for _, attachment := range fbPost.Attachments.Data {
//					if attachment.Type == string(fbclient.Album) {
//						for _, subAttachment := range attachment.SubAttachments.Data {
//							childPostIDs = append(childPostIDs, subAttachment.Target.ID)
//						}
//						break
//					}
//				}
//			}
//		}
//
//		// Get next posts
//		{
//			if !fbPostsResponse.Paging.CompareFacebookPagingRequest(fbPagingReq) {
//				newTaskID := cm.NewID()
//				s.mu.Lock()
//				t := rand.Intn(int(time.Second))
//				s.mapTaskArguments[newTaskID] = &TaskArguments{
//					actionType:      GetPosts,
//					accessToken:     accessToken,
//					shopID:          shopID,
//					externalPageID:  externalPageID,
//					pageID:          pageID,
//					fbPagingRequest: fbPostsResponse.Paging.ToPagingRequestAfter(fbclient.DefaultLimitGetPosts),
//				}
//				s.scheduler.AddAfter(newTaskID, time.Duration(t), s.syncCallbackLogs)
//				s.mu.Unlock()
//			}
//		}
//
//		// Get child_posts
//		{
//			start, end := 0, 0
//			numOfChildPosts := len(childPostIDs)
//
//			for {
//				start = end
//				end = min(end+fbclient.MaximumIDs, numOfChildPosts)
//				fmt.Println("accessToken ", accessToken)
//				fbPostsResponse, err := s.fbClient.CallAPIListPostsByIDs(accessToken, childPostIDs[start:end])
//				if err != nil {
//					return err
//				}
//
//				if len(fbPostsResponse.Data) == 0 {
//					break
//				}
//
//				for _, fbPost := range fbPostsResponse.Data {
//					postIDs = append(postIDs, fbPost.ID)
//				}
//
//				if end == numOfChildPosts {
//					break
//				}
//			}
//		}
//
//		// Get comment summaries and create task getCommentsByPostIDs
//		{
//			startPostIDs, endPostIDs := 0, 0
//			for {
//				startPostIDs = endPostIDs
//				endPostIDs = min(endPostIDs+fbclient.MaximumIDs, len(postIDs))
//				fbCommentSummaries, err := s.fbClient.CallAPIListCommentSummaries(accessToken, postIDs[startPostIDs:endPostIDs])
//				if err != nil {
//					return err
//				}
//
//				{
//					var tempPostIDs []string
//					for postID := range fbCommentSummaries {
//						tempPostIDs = append(tempPostIDs, postID)
//					}
//
//					newTaskID := cm.NewID()
//					s.mu.Lock()
//					t := rand.Intn(int(time.Second))
//					s.mapTaskArguments[newTaskID] = &TaskArguments{
//						actionType:     GetCommentsByPostIDs,
//						accessToken:    accessToken,
//						shopID:         shopID,
//						externalPageID: externalPageID,
//						pageID:         pageID,
//						GetCommentsByPostIDsArgs: &GetCommentsByPostIDsArguments{
//							IDs: tempPostIDs,
//						},
//					}
//					s.scheduler.AddAfter(newTaskID, time.Duration(t), s.syncCallbackLogs)
//					s.mu.Unlock()
//				}
//				if endPostIDs >= len(postIDs) {
//					break
//				}
//			}
//		}
//	case GetConversations:
//	case GetCommentsByPostIDs:
//		fmt.Println("GetCommentsByPostIDs")
//
//		defer func() {
//			err := recover()
//			if err != nil {
//				if cm.ErrorCode(err.(error)) == cm.FacebookError { // TODO: Ngoc handle missing permissions
//					s.scheduler.AddAfter(taskID, defaultRecurrentFacebook, s.syncCallbackLogs)
//				} else {
//					s.scheduler.AddAfter(taskID, defaultRecurrent, s.syncCallbackLogs)
//				}
//			}
//		}()
//
//		getCommentsByPostIDs := taskArgs.GetCommentsByPostIDsArgs
//		fbCommentsResponse, err := s.fbClient.CallAPIListCommentsByPostIDs(accessToken, getCommentsByPostIDs.IDs)
//		if err != nil {
//			return err
//		}
//
//		for _, fbComment := range fbCommentsResponse.Data {
//			if fbComment.Parent == nil {
//				fmt.Printf("Comment: %v\n", fbComment.Message)
//			} else {
//				fmt.Printf("Reply comment: %v\n", fbComment.Message)
//			}
//		}
//	}
//	return nil
//}

func listAllFbPagesActive(db *cmsql.Database) (fbpaging.FbExternalPageCombineds, error) {
	fromID := dot.ID(0)

	// key is id
	var fbExternalPageCombineds []*fbpaging.FbExternalPageCombined

	for {
		var fbPageModels fbpagemodel.FbExternalPages

		if err := db.
			Where("id > ?", fromID.Int64()).
			Where("connection_status = ?", status3.P.Enum()).
			Where("status = ?", status3.P.Enum()).
			OrderBy("id").
			Limit(1000).
			Find(&fbPageModels); err != nil {
			return nil, err
		}

		if len(fbPageModels) == 0 {
			break
		}

		var listFbPageIDs []dot.ID
		fbExternalPages := convert.Convert_fbpagemodel_FbExternalPages_fbpaging_FbExternalPages(fbPageModels)
		for _, fbExternalPage := range fbExternalPages {
			listFbPageIDs = append(listFbPageIDs, fbExternalPage.ID)
		}

		var fbExternalPageInternalModels fbpagemodel.FbExternalPageInternals

		if err := db.
			In("id", listFbPageIDs).
			OrderBy("id").
			Limit(1000).
			Find(&fbExternalPageInternalModels); err != nil {
			return nil, err
		}

		fbExternalPageInternals := convert.Convert_fbpagemodel_FbExternalPageInternals_fbpaging_FbExternalPageInternals(fbExternalPageInternalModels)

		for i, fbPage := range fbExternalPages {
			fbExternalPageCombineds = append(fbExternalPageCombineds, &fbpaging.FbExternalPageCombined{
				FbExternalPage:         fbPage,
				FbExternalPageInternal: fbExternalPageInternals[i],
			})
		}

		fromID = fbExternalPages[len(fbExternalPages)-1].ID
	}

	return fbExternalPageCombineds, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
