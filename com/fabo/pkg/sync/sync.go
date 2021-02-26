package sync

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_live_video_status"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/fabo/main/fbpage/convert"
	fbpagemodel "o.o/backend/com/fabo/main/fbpage/model"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	faboRedis "o.o/backend/com/fabo/pkg/redis"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/scheduler"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var ll = l.New().WithChannel("sync_service")

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
	GetLiveVideos    TaskActionType = 897
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

type getAvatarArguments struct {
	externalUserID       string
	externalPageInternal *fbpaging.FbExternalPageInternal
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
	getAvatarArgs    *getAvatarArguments

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
	fbUseringAggr    fbusering.CommandBus
	fbUseringQuery   fbusering.QueryBus

	mu sync.Mutex
	rd *faboRedis.FaboRedis

	mapExternalPageAndTimeStart map[string]time.Time

	timeLimit   int
	timeToCrawl int
}

type Config struct {
	TimeLimit   int `yaml:"time_limit"`    // days
	TimeToCrawl int `yaml:"time_to_crawl"` // mins
}

func New(
	db com.MainDB,
	fbClient *fbclient.FbClient,
	fbMessagingAggr fbmessaging.CommandBus, fbMessagingQuery fbmessaging.QueryBus,
	fbUseringAggr fbusering.CommandBus, fbUseringQuery fbusering.QueryBus,
	fbRedis *faboRedis.FaboRedis, cfg Config,
) *Synchronizer {
	sched := scheduler.New(defaultNumWorkers)
	s := &Synchronizer{
		scheduler:                   sched,
		db:                          db,
		fbClient:                    fbClient,
		mapTaskArguments:            make(map[dot.ID]*TaskArguments),
		mapExternalPageAndTimeStart: make(map[string]time.Time),
		fbMessagingAggr:             fbMessagingAggr,
		fbMessagingQuery:            fbMessagingQuery,
		fbUseringAggr:               fbUseringAggr,
		fbUseringQuery:              fbUseringQuery,
		rd:                          fbRedis,
		timeLimit:                   cfg.TimeLimit,
		timeToCrawl:                 cfg.TimeToCrawl,
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
	ticker := time.NewTicker(time.Duration(s.timeToCrawl) * time.Minute)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		fbPageCombineds, err := listAllFbPagesActive(s.db)
		if err != nil {
			return err
		}

		//now := time.Now()
		for _, fbPageCombined := range fbPageCombineds {
			// ignore Page test

			isTestPage := strings.HasPrefix(fbPageCombined.FbExternalPage.ExternalName, fbclient.PrefixFanPageNameTest)
			if cmenv.IsProd() && isTestPage {
				continue
			}
			if !s.rd.IsLockCallAPIPage(fbPageCombined.FbExternalPage.ExternalID) {
				// Task get post
				s.addTaskGetPosts(fbPageCombined)
				s.addTaskGetLiveVideos(fbPageCombined)
			}

			if !s.rd.IsLockCallAPIMessenger(fbPageCombined.FbExternalPage.ExternalID) {
				// Task get conversation
				s.addTask(&TaskArguments{
					actionType:     GetConversations,
					accessToken:    fbPageCombined.FbExternalPageInternal.Token,
					shopID:         fbPageCombined.FbExternalPage.ShopID,
					pageID:         fbPageCombined.FbExternalPage.ID,
					externalPageID: fbPageCombined.FbExternalPage.ExternalID,
					fbPagingRequest: &model.FacebookPagingRequest{
						Limit: dot.Int(fbclient.DefaultLimitGetConversations),
						TimePagination: &model.TimePaginationRequest{
							Since: time.Now().AddDate(0, 0, -s.timeLimit),
						},
					},
				})
			}
		}
	}
	return nil
}

func (s *Synchronizer) addTaskGetPosts(fbPageCombined *fbpaging.FbExternalPageCombined) dot.ID {
	return s.addTask(&TaskArguments{
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

func (s *Synchronizer) addTaskGetLiveVideos(fbPageCombined *fbpaging.FbExternalPageCombined) dot.ID {
	return s.addTask(&TaskArguments{
		actionType:     GetLiveVideos,
		accessToken:    fbPageCombined.FbExternalPageInternal.Token,
		shopID:         fbPageCombined.FbExternalPage.ShopID,
		pageID:         fbPageCombined.FbExternalPage.ID,
		externalPageID: fbPageCombined.FbExternalPage.ExternalID,
		fbPagingRequest: &model.FacebookPagingRequest{
			Limit: dot.Int(fbclient.DefaultLimitGetLiveVideos),
			TimePagination: &model.TimePaginationRequest{
				Since: time.Now().AddDate(0, 0, -s.timeLimit),
			},
		},
	})
}

func (s *Synchronizer) syncCallbackLogs(id interface{}, p scheduler.Planner) (_err error) {
	taskID := id.(dot.ID)
	ctx := bus.Ctx()
	taskArgs := s.getTaskArguments(taskID)

	accessToken := taskArgs.accessToken
	shopID := taskArgs.shopID
	pageID := taskArgs.pageID
	externalPageID := taskArgs.externalPageID
	fbPagingReq := taskArgs.fbPagingRequest

	defer func() {
		if _err == nil {
			s.finishTask(taskID)
			return
		}

		facebookError, ok := _err.(*xerrors.APIError)
		if ok {
			code := facebookError.Meta["code"]
			switch code {
			case fbclient.AccessTokenHasExpired.String():
				// no-op
				return
			case fbclient.ApiTooManyCalls.String(), fbclient.ApplicationLimitReached.String():
				s.scheduler.AddAfter(taskID, defaultRecurrentFacebookThrottling, s.syncCallbackLogs)
				return
			case fbclient.RateLimitCallWithPage.String():
				xBusinessUseCaseUsage := facebookError.Meta[fbclient.XBusinessUseCaseUsage]
				var xBusinessUseCaseUsageHeader fbclient.XBusinessUseCaseUsageHeader
				if err := jsonx.Unmarshal([]byte(xBusinessUseCaseUsage), &xBusinessUseCaseUsageHeader); err != nil {
					return
				}

				estimatedTimeToRegainAccess := xBusinessUseCaseUsageHeader.GetEstimatedTimeToRegainAccessOfPage(externalPageID)
				_ = s.rd.LockCallAPIPage(externalPageID, estimatedTimeToRegainAccess)
				return
			case fbclient.RateLimitCallWithMessenger.String():
				xBusinessUseCaseUsage := facebookError.Meta[fbclient.XBusinessUseCaseUsage]
				var xBusinessUseCaseUsageHeader fbclient.XBusinessUseCaseUsageHeader
				if err := jsonx.Unmarshal([]byte(xBusinessUseCaseUsage), &xBusinessUseCaseUsageHeader); err != nil {
					return
				}

				estimatedTimeToRegainAccess := xBusinessUseCaseUsageHeader.GetEstimatedTimeToRegainAccessOfMessenger(externalPageID)
				_ = s.rd.LockCallAPIMessenger(externalPageID, estimatedTimeToRegainAccess)
				return
			default:
				if codeNum, err := strconv.ParseInt(code, 10, 64); err == nil {
					if 200 <= codeNum && codeNum <= 299 {
						// no-op
						return
					}
				}
			}
		}
		go ll.SendMessage(_err.Error())
		s.scheduler.AddAfter(taskID, defaultRecurrentFacebook, s.syncCallbackLogs)
	}()

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
	case GetLiveVideos:
		if err := s.HandleTaskGetLiveVideos(ctx, shopID, pageID, accessToken, externalPageID, fbPagingReq); err != nil {
			return err
		}
	}
	return nil
}

func (s *Synchronizer) handleTaskGetMessages(
	ctx context.Context, shopID dot.ID, pageID dot.ID, accessToken string,
	externalPageID string, taskArgs *TaskArguments, fbPagingReq *model.FacebookPagingRequest,
) error {
	fmt.Println("GetMessages")

	conversationID := taskArgs.getMessagesArgs.conversationID
	externalConversationID := taskArgs.getMessagesArgs.externalConversationID

	fbMessagesResp, err := s.fbClient.CallAPIListMessages(&fbclient.ListMessagesRequest{
		AccessToken:    accessToken,
		ConversationID: externalConversationID,
		PageID:         externalPageID,
		Pagination:     fbPagingReq,
	})
	if err != nil {
		return err
	}

	if fbMessagesResp.Messages == nil ||
		len(fbMessagesResp.Messages.MessagesData) == 0 ||
		fbMessagesResp.Messages.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	isFinished := false

	var messagesData []*model.MessageData
	mapPSIDAndProfile := make(map[string]*model.Profile)
	var newFbExternalMessages []*fbmessaging.FbExternalMessage
	for _, fbMessage := range fbMessagesResp.Messages.MessagesData {
		if time.Now().Sub(fbMessage.CreatedTime.ToTime()) > time.Duration(s.timeLimit)*24*time.Hour {
			isFinished = true
			continue
		}
		messagesData = append(messagesData, fbMessage)
		if fbMessage.From != nil && fbMessage.From.ID != "" && fbMessage.From.ID != externalPageID {
			mapPSIDAndProfile[fbMessage.From.ID] = &model.Profile{
				ID:   fbMessage.From.ID,
				Name: fbMessage.From.Name,
			}
		}

		if fbMessage.To != nil {
			for _, fbMessageTo := range fbMessage.To.Data {
				if fbMessageTo.ID != "" && fbMessageTo.ID != externalPageID {
					mapPSIDAndProfile[fbMessageTo.ID] = &model.Profile{
						ID:   fbMessageTo.ID,
						Name: fbMessageTo.Name,
					}
				}
			}
		}
	}

	var PSIDs []string
	for psid := range mapPSIDAndProfile {
		PSIDs = append(PSIDs, psid)
	}

	for psid, profile := range mapPSIDAndProfile {
		if profile.ProfilePic == "" {
			_profile, err := s.getProfile(accessToken, externalPageID, psid, profile)
			if err != nil {
				return err
			}
			mapPSIDAndProfile[psid] = _profile
		}
	}
	mapPSIDAndProfile[externalPageID] = &model.Profile{ProfilePic: fmt.Sprintf("https://graph.facebook.com/%s/picture?height=200&width=200&type=normal", externalPageID)}

	for _, messageData := range messagesData {
		var externalAttachments []*fbmessaging.FbMessageAttachment
		var externalShares []*fbmessaging.FbMessageShare
		if messageData.Attachments != nil {
			externalAttachments = fbclientconvert.ConvertMessageDataAttachments(messageData.Attachments.Data)
		}
		if messageData.Shares != nil {
			externalShares = fbclientconvert.ConvertMessageShares(messageData.Shares.Data)
		}

		if messageData.From != nil {
			id := messageData.From.ID
			messageData.From.Picture = &model.Picture{
				Data: model.PictureData{
					Url: mapPSIDAndProfile[id].ProfilePic,
				},
			}
		}
		if messageData.To != nil {
			for _, messageTo := range messageData.To.Data {
				id := messageTo.ID
				messageTo.Picture = &model.Picture{
					Data: model.PictureData{
						Url: mapPSIDAndProfile[id].ProfilePic,
					},
				}
			}
		}

		currentMessage := messageData.Message
		{
			var strs []string
			if currentMessage != "" {
				strs = append(strs, currentMessage)
			}
			// Get first share
			if messageData.Sticker == "" && len(externalShares) > 0 {
				if externalShares[0].Description != "" {
					strs = append(strs, externalShares[0].Description)
				}
				if externalShares[0].Link != "" {
					strs = append(strs, externalShares[0].Link)
				} else {
					strs = append(strs, externalShares[0].Name)
				}
			}
			currentMessage = strings.Join(strs, "\n")
		}

		newFbExternalMessages = append(newFbExternalMessages, &fbmessaging.FbExternalMessage{
			ID:                     cm.NewID(),
			ExternalConversationID: externalConversationID,
			ExternalPageID:         externalPageID,
			ExternalID:             messageData.ID,
			ExternalMessage:        currentMessage,
			ExternalSticker:        messageData.Sticker,
			ExternalTo:             fbclientconvert.ConvertObjectsTo(messageData.To),
			ExternalFrom:           fbclientconvert.ConvertObjectFrom(messageData.From),
			ExternalAttachments:    externalAttachments,
			ExternalMessageShares:  externalShares,
			ExternalCreatedTime:    messageData.CreatedTime.ToTime(),
			ExternalTimestamp:      int64(*messageData.CreatedTime) * 1000,
			InternalSource:         fb_internal_source.Facebook,
		})
	}

	if len(newFbExternalMessages) > 0 {
		if err := s.fbMessagingAggr.Dispatch(ctx, &fbmessaging.CreateFbExternalMessagesFromSyncCommand{
			FbExternalMessages: newFbExternalMessages,
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
	fbConversationsResp, err := s.fbClient.CallAPIListConversations(&fbclient.ListConversationsRequest{
		AccessToken: accessToken,
		PageID:      externalPageID,
		Pagination:  fbPagingReq,
	})
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
			PSID:                 externalUserID,
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
	fbExternalCommentsResp, err := s.fbClient.CallAPIListComments(&fbclient.ListCommentsRequest{
		AccessToken: accessToken,
		PostID:      externalPostID,
		PageID:      externalPageID,
		Pagination:  fbPagingReq,
	})
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
		if fbExternalComment.From == nil {
			continue
		}
		if fbExternalComment.From.ID == externalPageID {
			continue
		}

		// Map externalUserID (From)
		// Map externalParentID (Parent)

		var externalUserID, externalParentID, externalParentUserID string
		externalUserID = fbExternalComment.From.ID
		if fbExternalComment.Parent != nil {
			externalParentID = fbExternalComment.Parent.ID
			if fbExternalComment.Parent.From != nil {
				externalParentUserID = fbExternalComment.Parent.From.ID
			}
		}

		// Try get old message (if it create by our api or webhook),
		// if already exists do not change value of field `InternalSource`
		// otherwise set is default to `fb_internal_source.Facebook`
		commentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
			ExternalID: fbExternalComment.ID,
		}
		if err := s.fbMessagingQuery.Dispatch(ctx, commentQuery); err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}
		comment := commentQuery.Result
		internalSource := fb_internal_source.Facebook
		var createdBy dot.ID
		var isLiked, isPrivateReplied bool
		if comment != nil {
			internalSource = comment.InternalSource
			createdBy = comment.CreatedBy
			isLiked = comment.IsLiked
			isPrivateReplied = comment.IsPrivateReplied
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
			InternalSource:       internalSource,
			IsLiked:              isLiked,
			IsHidden:             fbExternalComment.IsHidden,
			IsPrivateReplied:     isPrivateReplied,
			CreatedBy:            createdBy,
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
	fbExternalPostResp, err := s.fbClient.CallAPIGetPost(&fbclient.GetPostRequest{
		AccessToken: accessToken,
		PostID:      externalChildPostID,
		PageID:      externalPageID,
	})
	if err != nil {
		return err
	}

	var totalComments, totalReactions int
	if fbExternalPostResp.CommentsSummary != nil && fbExternalPostResp.CommentsSummary.Summary != nil {
		totalComments = fbExternalPostResp.CommentsSummary.Summary.TotalCount
	}
	if fbExternalPostResp.ReactionsSummary != nil && fbExternalPostResp.ReactionsSummary.Summary != nil {
		totalReactions = fbExternalPostResp.ReactionsSummary.Summary.TotalCount
	}

	var createFbExternalPostsArgs []*fbmessaging.CreateFbExternalPostArgs
	createFbExternalPostsArgs = append(createFbExternalPostsArgs, &fbmessaging.CreateFbExternalPostArgs{
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
		TotalComments:       totalComments,
		TotalReactions:      totalReactions,
		StatusType:          fb_status_type.ParseFbStatusTypeWithDefault(fbExternalPostResp.StatusType, fb_status_type.Unknown),
	})

	updateOrCreateFbExternalPostsCmd := &fbmessaging.UpdateOrCreateFbExternalPostsFromSyncCommand{
		FbExternalPosts: createFbExternalPostsArgs,
	}
	if err := s.fbMessagingAggr.Dispatch(ctx, updateOrCreateFbExternalPostsCmd); err != nil {
		return err
	}
	if len(updateOrCreateFbExternalPostsCmd.Result) == 0 {
		return nil
	}

	childPostID := updateOrCreateFbExternalPostsCmd.Result[0].ID

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
	fbExternalPostsResp, err := s.fbClient.CallAPIListFeeds(&fbclient.ListFeedsRequest{
		AccessToken: accessToken,
		PageID:      externalPageID,
		Pagination:  fbPagingReq,
	})
	if err != nil {
		return err
	}

	// Finish task when data is empty
	if fbExternalPostsResp == nil || len(fbExternalPostsResp.Data) == 0 ||
		fbExternalPostsResp.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	// Get all child posts of each post
	// Create and Update posts
	mapExternalChildPostIDAndExternalPostID := make(map[string]string)
	var createFbExternalPostsArgs []*fbmessaging.CreateFbExternalPostArgs
	for _, fbPost := range fbExternalPostsResp.Data {
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

		var totalComments, totalReactions int
		if fbPost.CommentsSummary != nil && fbPost.CommentsSummary.Summary != nil {
			totalComments = fbPost.CommentsSummary.Summary.TotalCount
		}
		if fbPost.ReactionsSummary != nil && fbPost.ReactionsSummary.Summary != nil {
			totalReactions = fbPost.ReactionsSummary.Summary.TotalCount
		}

		createFbExternalPostsArgs = append(createFbExternalPostsArgs, &fbmessaging.CreateFbExternalPostArgs{
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
			TotalComments:       totalComments,
			TotalReactions:      totalReactions,
			StatusType:          fb_status_type.ParseFbStatusTypeWithDefault(fbPost.StatusType, fb_status_type.Unknown),
		})
	}

	updateOrCreateFbExternalPostsCmd := &fbmessaging.UpdateOrCreateFbExternalPostsFromSyncCommand{
		FbExternalPosts: createFbExternalPostsArgs,
	}
	if err := s.fbMessagingAggr.Dispatch(ctx, updateOrCreateFbExternalPostsCmd); err != nil {
		return err
	}

	mapExternalPostIDAndPostID := make(map[string]dot.ID)
	for _, fbExternalPost := range updateOrCreateFbExternalPostsCmd.Result {
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
	for _, fbExternalPost := range updateOrCreateFbExternalPostsCmd.Result {
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

func (s *Synchronizer) HandleTaskGetLiveVideos(
	ctx context.Context, shopID, pageID dot.ID,
	accessToken, externalPageID string, fbPagingReq *model.FacebookPagingRequest,
) error {
	fmt.Println("GetLiveVideos")

	// Call api (facebook) listPublishedPosts from facebook
	fbLiveVideosResp, err := s.fbClient.CallAPIListSimplifyLiveVideos(&fbclient.ListSimplifyLiveVideosRequest{
		AccessToken: accessToken,
		Pagination:  fbPagingReq,
	})
	if err != nil {
		return err
	}

	// Finish task when data is empty
	if fbLiveVideosResp == nil || len(fbLiveVideosResp.Data) == 0 ||
		fbLiveVideosResp.Paging.CompareFacebookPagingRequest(fbPagingReq) {
		return nil
	}

	for _, fbLiveVideo := range fbLiveVideosResp.Data {
		updateLiveVideoStatusCmd := &fbmessaging.UpdateLiveVideoStatusFromSyncCommand{
			ExternalID:              fbLiveVideo.GetExternalPostID(),
			ExternalLiveVideoStatus: fbLiveVideo.Status,
			LiveVideoStatus:         fb_live_video_status.ConvertToFbLiveVideoStatus(fbLiveVideo.Status),
		}
		if err := s.fbMessagingAggr.Dispatch(ctx, updateLiveVideoStatusCmd); err != nil {
			return err
		}
	}

	return nil
}

func (s *Synchronizer) getProfile(accessToken, externalPageID, PSID string, profileDefault *model.Profile) (*model.Profile, error) {
	profile, err := s.rd.LoadProfilePSID(externalPageID, PSID)
	switch err {
	// If profile not in redis then call api getProfileByPSID
	case redis.ErrNil:
		_profile, _err := s.fbClient.CallAPIGetProfileByPSID(&fbclient.GetProfileRequest{
			AccessToken:    accessToken,
			PSID:           PSID,
			PageID:         externalPageID,
			ProfileDefault: profileDefault,
		})
		if _err != nil {
			return nil, _err
		}
		if _err := s.rd.SaveProfilePSID(externalPageID, PSID, _profile); _err != nil {
			return nil, _err
		}
		return _profile, nil
	case nil:
		return profile, nil
	default:
		ll.SendMessagef("%v %v %v", externalPageID, PSID, err.Error())
		return nil, err
	}
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
			OrderBy("updated_at desc, id asc").
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

		// TODO(ngoc): refactor
		if err := db.
			In("id", listFbPageIDs).
			OrderBy("id").
			Limit(1000).
			Find(&fbExternalPageInternalModels); err != nil {
			return nil, err
		}

		fbExternalPageInternals := convert.Convert_fbpagemodel_FbExternalPageInternals_fbpaging_FbExternalPageInternals(fbExternalPageInternalModels)
		mapFbExternalPageInternal := make(map[string]*fbpaging.FbExternalPageInternal)
		for _, fbExternalPageInternal := range fbExternalPageInternals {
			mapFbExternalPageInternal[fbExternalPageInternal.ExternalID] = fbExternalPageInternal
		}

		for _, fbPage := range fbExternalPages {
			fromID = max(fromID, fbPage.ID)
			fbExternalPageCombineds = append(fbExternalPageCombineds, &fbpaging.FbExternalPageCombined{
				FbExternalPage:         fbPage,
				FbExternalPageInternal: mapFbExternalPageInternal[fbPage.ExternalID],
			})
		}
	}

	return fbExternalPageCombineds, nil
}

func ternaryString(statement bool, a, b string) string {
	if statement {
		return a
	}
	return b
}

func max(a, b dot.ID) dot.ID {
	if a > b {
		return a
	}
	return b
}
