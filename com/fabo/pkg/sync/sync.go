package sync

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"o.o/api/fabo/fbmessaging"
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
	"o.o/backend/pkg/common/extservice/telebot"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var ll = l.New()

const (
	defaultNumWorkers                  = 32
	defaultRecurrentFacebook           = 1 * time.Minute
	defaultRecurrent                   = 5 * time.Minute
	defaultErrRecurr                   = 10 * time.Minute
	defaultRecurrentFacebookThrottling = 60 * time.Minute
	maxLogsCount                       = 500
)

type TaskActionType int

const (
	GetPosts         TaskActionType = 737
	GetChildPost     TaskActionType = 230
	GetComments      TaskActionType = 655
	GetConversations TaskActionType = 363
	GetMessages      TaskActionType = 160
)

type getCommentsArguments struct {
	externalPostID string
	postID         dot.ID
}

type getChildPostArguments struct {
	externalChildPostID string
	externalPostID      string
}

type getMessagesArguments struct {
	conversationID         dot.ID
	externalConversationID string
}

type TaskArguments struct {
	actionType  TaskActionType
	accessToken string
	shopID      dot.ID

	pageID         dot.ID
	externalPageID string

	getCommentsArgs  *getCommentsArguments
	getChildPostArgs *getChildPostArguments
	getMessagesArgs  *getMessagesArguments

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

	bot *telebot.Channel
	mu  sync.Mutex

	timeLimit int
}

func New(
	db *cmsql.Database, fbClient *fbclient.FbClient,
	fbMessagingAggr fbmessaging.CommandBus, fbMessagingQuery fbmessaging.QueryBus,
	bot *telebot.Channel, timeLimit int,
) *Synchronizer {
	sched := scheduler.New(defaultNumWorkers)
	s := &Synchronizer{
		scheduler:        sched,
		db:               db,
		fbClient:         fbClient,
		mapTaskArguments: make(map[dot.ID]*TaskArguments),
		fbMessagingAggr:  fbMessagingAggr,
		fbMessagingQuery: fbMessagingQuery,
		bot:              bot,
		timeLimit:        timeLimit,
	}
	return s
}

func (s *Synchronizer) Init() error {
	s.scheduler.AddAfter(cm.NewID(), 1*time.Second, s.addJobs)
	return nil
}

func (s *Synchronizer) Start() {
	s.scheduler.Start()
}

func (s *Synchronizer) Stop() {
	s.scheduler.Stop()
}

func (s *Synchronizer) addTask(taskArguments *TaskArguments) (taskID dot.ID) {
	if taskArguments == nil {
		return
	}

	s.mu.Lock()
	taskID = cm.NewID()
	s.mapTaskArguments[taskID] = taskArguments
	t := rand.Intn(int(time.Second))
	s.scheduler.AddAfter(taskID, time.Duration(t), s.syncCallbackLogs)
	s.mu.Unlock()
	return
}

func (s *Synchronizer) finishTask(taskID dot.ID) {
	s.mu.Lock()
	delete(s.mapTaskArguments, taskID)
	s.mu.Unlock()
}

func (s *Synchronizer) getTaskArguments(taskID dot.ID) (taskArguments *TaskArguments) {
	s.mu.Lock()
	taskArguments = s.mapTaskArguments[taskID]
	s.mu.Unlock()
	return
}

func (s *Synchronizer) addJobs(id interface{}, p scheduler.Planner) (_err error) {
	fbPageCombineds, err := listAllFbPagesActive(s.db)
	if err != nil {
		return err
	}

	for _, fbPageCombined := range fbPageCombineds {
		// Task get post of specific page
		{
			s.addTask(&TaskArguments{
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
			})
		}

		// Task get conversation
		//{
		//	s.addTask(&TaskArguments{
		//		actionType:     GetConversations,
		//		accessToken:    fbPageCombined.FbExternalPageInternal.Token,
		//		shopID:         fbPageCombined.FbExternalPage.ShopID,
		//		pageID:         fbPageCombined.FbExternalPage.ID,
		//		externalPageID: fbPageCombined.FbExternalPage.ExternalID,
		//		fbPagingRequest: &model.FacebookPagingRequest{
		//			Limit: dot.Int(fbclient.DefaultLimitGetConversations),
		//			TimePagination: &model.TimePaginationRequest{
		//				Since: time.Now().AddDate(0, 0, -s.timeLimit),
		//			},
		//		},
		//	})
		//}
	}

	s.scheduler.AddAfter(cm.NewID(), 5*time.Minute, s.addJobs)

	return
}

func (s *Synchronizer) syncCallbackLogs(id interface{}, p scheduler.Planner) (_err error) {
	taskID := id.(dot.ID)
	ctx := bus.Ctx()

	defer func() {
		if _err == nil {
			s.finishTask(taskID)
			return
		}

		facebookError := _err.(*xerrors.APIError)
		code := facebookError.Meta["code"]
		switch code {
		case fbclient.AccessTokenHasExpired.String():
			// no-op
			return
		case fbclient.ApiTooManyCalls.String(), fbclient.ApplicationLimitReached.String():
			s.scheduler.AddAfter(taskID, defaultRecurrentFacebookThrottling, s.syncCallbackLogs)
			return
		default:
			if codeNum, err := strconv.ParseInt(code, 10, 64); err == nil {
				if 200 <= codeNum && codeNum <= 299 {
					// no-op
					return
				}
			}
		}
		go s.bot.SendMessage(_err.Error())
		s.scheduler.AddAfter(taskID, defaultRecurrentFacebook, s.syncCallbackLogs)
	}()

	taskArgs := s.getTaskArguments(taskID)

	accessToken := taskArgs.accessToken
	shopID := taskArgs.shopID
	pageID := taskArgs.pageID
	externalPageID := taskArgs.externalPageID
	fbPagingReq := taskArgs.fbPagingRequest

	switch taskArgs.actionType {
	case GetPosts:
		if err := s.HandleTaskGetPosts(ctx, shopID, pageID, accessToken, externalPageID, fbPagingReq); err != nil {
			return err
		}
	case GetChildPost:
		if err := s.handleTaskGetChildPost(ctx, shopID, pageID, accessToken, externalPageID, taskArgs); err != nil {
			return err
		}
	case GetComments:
		if err := s.handleTaskGetComments(ctx, pageID, accessToken, externalPageID, fbPagingReq, taskArgs); err != nil {
			return err
		}
	case GetConversations:
		if err := s.handleTaskGetConversations(ctx, shopID, pageID, accessToken, externalPageID, fbPagingReq); err != nil {
			return err
		}
	case GetMessages:
		if err := s.handleTaskGetMessages(ctx, shopID, pageID, accessToken, externalPageID, taskArgs, fbPagingReq); err != nil {
			return err
		}
	}
	return nil
}

func (s *Synchronizer) handleTaskGetMessages(
	ctx context.Context, shopID dot.ID, pageID dot.ID, accessToken string,
	externalPageID string, taskArgs *TaskArguments, fbPagingReq *model.FacebookPagingRequest,
) error {
	fmt.Println("getMessagesArgs")

	conversationID := taskArgs.getMessagesArgs.conversationID
	externalConversationID := taskArgs.getMessagesArgs.externalConversationID

	fbMessagesResp, err := s.fbClient.CallAPIListMessages(accessToken, externalConversationID, fbPagingReq)
	if err != nil {
		return err
	}

	if fbMessagesResp.Messages == nil ||
		len(fbMessagesResp.Messages.MessagesData) == 0 ||
		fbMessagesResp.Messages.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	isFinished := false
	var fbExternalMessagesArgs []*fbmessaging.CreateFbExternalMessageArgs
	for _, fbMessage := range fbMessagesResp.Messages.MessagesData {
		if time.Now().Sub(fbMessage.CreatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
			isFinished = true
			continue
		}
		var externalAttachments []*fbmessaging.FbMessageAttachment
		if fbMessage.Attachments != nil {
			externalAttachments = fbclientconvert.ConvertMessageDataAttachments(fbMessage.Attachments.Data)
		}
		fbExternalMessagesArgs = append(fbExternalMessagesArgs, &fbmessaging.CreateFbExternalMessageArgs{
			ID:                     cm.NewID(),
			ExternalConversationID: externalConversationID,
			ExternalPageID:         externalPageID,
			ExternalID:             fbMessage.ID,
			ExternalMessage:        fbMessage.Message,
			ExternalSticker:        fbMessage.Sticker,
			ExternalTo:             fbclientconvert.ConvertObjectsTo(fbMessage.To),
			ExternalFrom:           fbclientconvert.ConvertObjectFrom(fbMessage.From),
			ExternalAttachments:    externalAttachments,
			ExternalCreatedTime:    fbMessage.CreatedTime.ToTime(),
		})
	}

	if len(fbExternalMessagesArgs) > 0 {
		if err := s.fbMessagingAggr.Dispatch(ctx, &fbmessaging.CreateOrUpdateFbExternalMessagesCommand{
			FbExternalMessages: fbExternalMessagesArgs,
		}); err != nil {
			return err
		}
	}

	if !isFinished {
		s.addTask(&TaskArguments{
			actionType:     GetMessages,
			accessToken:    accessToken,
			shopID:         shopID,
			externalPageID: externalPageID,
			pageID:         pageID,
			getMessagesArgs: &getMessagesArguments{
				conversationID:         conversationID,
				externalConversationID: externalConversationID,
			},
			fbPagingRequest: fbMessagesResp.Messages.Paging.ToPagingRequestAfter(fbclient.DefaultLimitGetMessages),
		})
	}
	return nil
}

func (s *Synchronizer) handleTaskGetConversations(
	ctx context.Context, shopID dot.ID, pageID dot.ID,
	accessToken string, externalPageID string, fbPagingReq *model.FacebookPagingRequest,
) error {
	fmt.Println("GetConversations")

	// Call api list conversations that depends on externalPageID
	fbConversationsResp, err := s.fbClient.CallAPIListConversations(accessToken, externalPageID, "", fbPagingReq)
	if err != nil {
		// TODO: Ngoc classify error type
		return err
	}

	// Finish task when data response is empty
	if fbConversationsResp.Conversations == nil ||
		len(fbConversationsResp.Conversations.ConversationsData) == 0 ||
		fbConversationsResp.Conversations.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	isFinished := false
	var fbExternalConversationsArgs []*fbmessaging.CreateFbExternalConversationArgs
	for _, fbConversation := range fbConversationsResp.Conversations.ConversationsData {
		if time.Now().Sub(fbConversation.UpdatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
			isFinished = true
			continue
		}

		var externalUserID, externalUserName string
		for _, sender := range fbConversation.Senders.Data {
			if sender.ID != externalPageID {
				externalUserID = sender.ID
				externalUserName = sender.Name
				break
			}
		}
		fbExternalConversationsArgs = append(fbExternalConversationsArgs, &fbmessaging.CreateFbExternalConversationArgs{
			ID:                   cm.NewID(),
			ExternalID:           fbConversation.ID,
			ExternalPageID:       externalPageID,
			ExternalUserID:       externalUserID,
			ExternalUserName:     externalUserName,
			ExternalLink:         fbConversation.Link,
			ExternalUpdatedTime:  fbConversation.UpdatedTime.ToTime(),
			ExternalMessageCount: fbConversation.MessageCount,
		})
	}

	mapExternalIDAndID := make(map[string]dot.ID)
	if len(fbExternalConversationsArgs) > 0 {
		createOrUpdateFbExternalConversationsCmd := &fbmessaging.CreateOrUpdateFbExternalConversationsCommand{
			FbExternalConversations: fbExternalConversationsArgs,
		}
		if err := s.fbMessagingAggr.Dispatch(ctx, createOrUpdateFbExternalConversationsCmd); err != nil {
			return err
		}

		for _, fbConversation := range createOrUpdateFbExternalConversationsCmd.Result {
			mapExternalIDAndID[fbConversation.ExternalID] = fbConversation.ID
		}
	}

	for _, fbConversation := range fbConversationsResp.Conversations.ConversationsData {
		if time.Now().Sub(fbConversation.UpdatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
			continue
		}
		s.addTask(&TaskArguments{
			actionType:     GetMessages,
			accessToken:    accessToken,
			shopID:         shopID,
			externalPageID: externalPageID,
			pageID:         pageID,
			getMessagesArgs: &getMessagesArguments{
				conversationID:         mapExternalIDAndID[fbConversation.ID],
				externalConversationID: fbConversation.ID,
			},
			fbPagingRequest: &model.FacebookPagingRequest{
				Limit: dot.Int(fbclient.DefaultLimitGetMessages),
			},
		})
	}

	if !isFinished {
		s.addTask(&TaskArguments{
			actionType:      GetConversations,
			accessToken:     accessToken,
			shopID:          shopID,
			externalPageID:  externalPageID,
			pageID:          pageID,
			fbPagingRequest: fbConversationsResp.Conversations.Paging.ToPagingRequestAfter(fbclient.DefaultLimitGetConversations),
		})
	}
	return nil
}

func (s *Synchronizer) handleTaskGetComments(
	ctx context.Context, pageID dot.ID, accessToken string,
	externalPageID string, fbPagingReq *model.FacebookPagingRequest, taskArgs *TaskArguments,
) error {
	fmt.Println("GetComments")

	// Get values from arguments
	externalPostID := taskArgs.getCommentsArgs.externalPostID
	//postID := taskArgs.getCommentsArgs.postID

	// Call api list comments that depends on (externalPostID)
	fbExternalCommentsResp, err := s.fbClient.CallAPIListComments(accessToken, externalPostID, fbPagingReq)
	if err != nil {
		return err
	}

	// Finish task when data response is empty
	if fbExternalCommentsResp.Comments == nil ||
		len(fbExternalCommentsResp.Comments.CommentData) == 0 ||
		fbExternalCommentsResp.Comments.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	var createOrUpdateFbExternalCommentsArgs []*fbmessaging.CreateFbExternalCommentArgs
	for _, fbExternalComment := range fbExternalCommentsResp.Comments.CommentData {
		if fbExternalComment.IsHidden {
			continue
		}
		var externalUserID, externalParentID, externalParentUserID string
		if fbExternalComment.From != nil {
			externalUserID = fbExternalComment.From.ID
		}
		if fbExternalComment.Parent != nil {
			externalParentID = fbExternalComment.Parent.ID
			if fbExternalComment.Parent.From != nil {
				externalParentUserID = fbExternalComment.Parent.From.ID
			}
		}
		createOrUpdateFbExternalCommentsArgs = append(createOrUpdateFbExternalCommentsArgs, &fbmessaging.CreateFbExternalCommentArgs{
			ID:                   cm.NewID(),
			ExternalPostID:       externalPostID,
			ExternalPageID:       externalPageID,
			ExternalID:           fbExternalComment.ID,
			ExternalUserID:       externalUserID,
			ExternalParentID:     externalParentID,
			ExternalParentUserID: externalParentUserID,
			ExternalMessage:      fbExternalComment.Message,
			ExternalCommentCount: fbExternalComment.CommentCount,
			ExternalParent:       fbclientconvert.ConvertFbObjectParent(fbExternalComment.Parent),
			ExternalFrom:         fbclientconvert.ConvertObjectFrom(fbExternalComment.From),
			ExternalAttachment:   fbclientconvert.ConvertFbCommentAttachment(fbExternalComment.Attachment),
			ExternalCreatedTime:  fbExternalComment.CreatedTime.ToTime(),
		})
	}

	if len(createOrUpdateFbExternalCommentsArgs) == 0 {
		return nil
	}

	createOrUpdateFbExternalCommentsCmd := &fbmessaging.CreateOrUpdateFbExternalCommentsCommand{
		FbExternalComments: createOrUpdateFbExternalCommentsArgs,
	}
	if err := s.fbMessagingAggr.Dispatch(ctx, createOrUpdateFbExternalCommentsCmd); err != nil {
		return err
	}
	return nil
}

func (s *Synchronizer) handleTaskGetChildPost(
	ctx context.Context, shopID dot.ID, pageID dot.ID,
	accessToken string, externalPageID string, taskArgs *TaskArguments,
) error {
	fmt.Println("GetChildPost")

	// Get values from arguments
	externalChildPostID := taskArgs.getChildPostArgs.externalChildPostID
	externalPostID := taskArgs.getChildPostArgs.externalPostID

	// Call api get (child) post that depends on postID
	fbExternalPostResp, err := s.fbClient.CallAPIGetPost(externalChildPostID, accessToken)
	if err != nil {
		return err
	}

	var createOrUpdateFbExternalPostsArgs []*fbmessaging.CreateFbExternalPostArgs

	createOrUpdateFbExternalPostsArgs = append(createOrUpdateFbExternalPostsArgs, &fbmessaging.CreateFbExternalPostArgs{
		ID:                  cm.NewID(),
		ExternalPageID:      externalPageID,
		ExternalID:          externalChildPostID,
		ExternalParentID:    externalPostID,
		ExternalFrom:        fbclientconvert.ConvertObjectFrom(fbExternalPostResp.From),
		ExternalPicture:     fbExternalPostResp.FullPicture,
		ExternalIcon:        fbExternalPostResp.Icon,
		ExternalMessage:     ternaryString(fbExternalPostResp.Message != "", fbExternalPostResp.Message, fbExternalPostResp.Story),
		ExternalAttachments: fbclientconvert.ConvertAttachments(fbExternalPostResp.Attachments),
		ExternalCreatedTime: fbExternalPostResp.CreatedTime.ToTime(),
		ExternalUpdatedTime: fbExternalPostResp.UpdatedTime.ToTime(),
	})

	createOrUpdateFbExternalPostsCmd := &fbmessaging.CreateOrUpdateFbExternalPostsCommand{
		FbExternalPosts: createOrUpdateFbExternalPostsArgs,
	}
	if err := s.fbMessagingAggr.Dispatch(ctx, createOrUpdateFbExternalPostsCmd); err != nil {
		return err
	}
	childPostID := createOrUpdateFbExternalPostsCmd.Result[0].ID

	s.addTask(&TaskArguments{
		actionType:     GetComments,
		accessToken:    accessToken,
		shopID:         shopID,
		pageID:         pageID,
		externalPageID: externalPageID,
		getCommentsArgs: &getCommentsArguments{
			externalPostID: externalChildPostID,
			postID:         childPostID,
		},
		fbPagingRequest: &model.FacebookPagingRequest{
			Limit: dot.Int(fbclient.DefaultLimitGetComments),
			TimePagination: &model.TimePaginationRequest{
				Since: time.Now().AddDate(0, 0, -s.timeLimit),
			},
		},
	})
	return nil
}

func (s *Synchronizer) HandleTaskGetPosts(
	ctx context.Context, shopID dot.ID, pageID dot.ID,
	accessToken string, externalPageID string, fbPagingReq *model.FacebookPagingRequest,
) error {
	fmt.Println("GetPosts")

	// Call api (facebook) listPublishedPosts from facebook
	fbExternalPostsResp, err := s.fbClient.CallAPIListPublishedPosts(accessToken, externalPageID, fbPagingReq)
	if err != nil {
		return err
	}

	// Finish task when data is empty
	if len(fbExternalPostsResp.Data) == 0 ||
		fbExternalPostsResp.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	// Get all child posts of each post
	// Create and Update posts
	mapExternalChildPostIDAndExternalPostID := make(map[string]string)
	var createOrUpdateFbExternalPostsArgs []*fbmessaging.CreateFbExternalPostArgs
	for _, fbPost := range fbExternalPostsResp.Data {
		if fbPost.ID == "472044163146832_1113902165627692" {
			fmt.Println("a")
		}
		if fbPost.Attachments != nil {
			for _, attachment := range fbPost.Attachments.Data {
				// TODO: Ngoc add enum
				if attachment.Type == "album" {
					for _, subAttachment := range attachment.SubAttachments.Data {
						if subAttachment.Type == "photo" {
							mapExternalChildPostIDAndExternalPostID[fmt.Sprintf("%s_%s", externalPageID, subAttachment.Target.ID)] = fbPost.ID
						}
					}
				}
			}
		}

		createOrUpdateFbExternalPostsArgs = append(createOrUpdateFbExternalPostsArgs, &fbmessaging.CreateFbExternalPostArgs{
			ID:                  cm.NewID(),
			ExternalPageID:      externalPageID,
			ExternalID:          fbPost.ID,
			ExternalFrom:        fbclientconvert.ConvertObjectFrom(fbPost.From),
			ExternalPicture:     fbPost.FullPicture,
			ExternalIcon:        fbPost.Icon,
			ExternalMessage:     ternaryString(fbPost.Message != "", fbPost.Message, fbPost.Story),
			ExternalAttachments: fbclientconvert.ConvertAttachments(fbPost.Attachments),
			ExternalCreatedTime: fbPost.CreatedTime.ToTime(),
			ExternalUpdatedTime: fbPost.UpdatedTime.ToTime(),
		})
	}

	createOrUpdateFbExternalPostsCmd := &fbmessaging.CreateOrUpdateFbExternalPostsCommand{
		FbExternalPosts: createOrUpdateFbExternalPostsArgs,
	}
	if err := s.fbMessagingAggr.Dispatch(ctx, createOrUpdateFbExternalPostsCmd); err != nil {
		return err
	}

	mapExternalPostIDAndPostID := make(map[string]dot.ID)
	for _, fbExternalPost := range createOrUpdateFbExternalPostsCmd.Result {
		mapExternalPostIDAndPostID[fbExternalPost.ExternalID] = fbExternalPost.ID
	}

	// Create tasks getChildPost
	for externalChildPostID, externalPostID := range mapExternalChildPostIDAndExternalPostID {
		s.addTask(&TaskArguments{
			actionType:     GetChildPost,
			accessToken:    accessToken,
			shopID:         shopID,
			pageID:         pageID,
			externalPageID: externalPageID,
			getChildPostArgs: &getChildPostArguments{
				externalChildPostID: externalChildPostID,
				externalPostID:      externalPostID,
			},
		})
	}

	// Create tasks getComments
	for _, fbExternalPost := range createOrUpdateFbExternalPostsCmd.Result {
		s.addTask(&TaskArguments{
			actionType:     GetComments,
			accessToken:    accessToken,
			shopID:         shopID,
			pageID:         pageID,
			externalPageID: externalPageID,
			getCommentsArgs: &getCommentsArguments{
				externalPostID: fbExternalPost.ExternalID,
				postID:         mapExternalPostIDAndPostID[fbExternalPost.ExternalID],
			},
			fbPagingRequest: &model.FacebookPagingRequest{
				Limit: dot.Int(fbclient.DefaultLimitGetComments),
			},
		})
	}

	return nil
}

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

func ternaryString(statement bool, a, b string) string {
	if statement {
		return a
	}
	return b
}
