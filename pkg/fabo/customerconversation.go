package fabo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_comment_source"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/common"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	convertpb2 "o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/backend/pkg/fabo/faboinfo"
	"o.o/common/xerrors"
)

type CustomerConversationService struct {
	session.Session

	FaboPagesKit     *faboinfo.FaboPagesKit
	FBMessagingQuery fbmessaging.QueryBus
	FBMessagingAggr  fbmessaging.CommandBus
	FBPagingQuery    fbpaging.QueryBus
	FBClient         *fbclient.FbClient
	FBUserQuery      fbusering.QueryBus
}

func (s *CustomerConversationService) Clone() fabo.CustomerConversationService {
	res := *s
	return &res
}

func (s *CustomerConversationService) ListCustomerConversations(
	ctx context.Context, request *fabo.ListCustomerConversationsRequest,
) (*fabo.FbCustomerConversationsResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	listCustomerConversationsQuery := &fbmessaging.ListFbCustomerConversationsQuery{
		Paging: *paging,
	}
	if request.Filter != nil {
		fbPageIDsRequest := request.Filter.FbPageIDs
		externalPageIDsRequest := request.Filter.ExternalPageID
		if len(fbPageIDsRequest) != 0 {
			listFbExternalPagesByIDsQuery := &fbpaging.ListFbExternalPagesByIDsQuery{
				IDs: fbPageIDsRequest,
			}
			if err := s.FBPagingQuery.Dispatch(ctx, listFbExternalPagesByIDsQuery); err != nil {
				return nil, err
			}
			for _, fbExternalPage := range listFbExternalPagesByIDsQuery.Result {
				externalPageIDsRequest = append(externalPageIDsRequest, fbExternalPage.ExternalID)
			}
		}

		if len(fbPageIDsRequest) == 0 && len(externalPageIDsRequest) == 0 {
			faboInfo, err := s.FaboPagesKit.GetPages(ctx, s.SS.Shop().ID)
			if err != nil {
				return nil, err
			}
			fbPageIDsRequest = faboInfo.FbPageIDs
			externalPageIDsRequest = faboInfo.ExternalPageIDs
		}
		listCustomerConversationsQuery.ExternalPageIDs = externalPageIDsRequest
		listCustomerConversationsQuery.IsRead = request.Filter.IsRead

		if request.Filter.FbExternalUserID.Valid {
			listCustomerConversationsQuery.ExternalUserID = request.Filter.FbExternalUserID
		} else {
			listCustomerConversationsQuery.ExternalUserID = request.Filter.ExternalUserID
		}
		listCustomerConversationsQuery.Type = request.Filter.Type
	}

	if len(listCustomerConversationsQuery.ExternalPageIDs) == 0 {
		faboInfo, err := s.FaboPagesKit.GetPages(ctx, s.SS.Shop().ID)
		if err != nil {
			return nil, err
		}
		listCustomerConversationsQuery.ExternalPageIDs = faboInfo.ExternalPageIDs
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, listCustomerConversationsQuery); err != nil {
		return nil, err
	}

	listCustomerConversations := listCustomerConversationsQuery.Result.FbCustomerConversations
	// get avatars
	{
		var externalUserIDs []string
		for _, customerConversation := range listCustomerConversations {
			if customerConversation.ExternalFrom != nil {
				externalUserIDs = append(externalUserIDs, customerConversation.ExternalFrom.ID)
			}
			if customerConversation.ExternalUserID != "" {
				externalUserIDs = append(externalUserIDs, customerConversation.ExternalUserID)
			}
		}

		mapExternalUserIDAndImageURl, err := s.getImageURLs(ctx, externalUserIDs)
		if err != nil {
			return nil, err
		}

		for _, customerConversation := range listCustomerConversations {
			if customerConversation.ExternalFrom != nil {
				customerConversation.ExternalFrom.ImageURL = mapExternalUserIDAndImageURl[customerConversation.ExternalFrom.ID]
			}
			if customerConversation.ExternalUserID != "" {
				customerConversation.ExternalUserPictureURL = mapExternalUserIDAndImageURl[customerConversation.ExternalUserID]
			}
		}
	}

	var fbUserExternalID []string
	for _, v := range listCustomerConversations {
		fbUserExternalID = append(fbUserExternalID, v.ExternalUserID)
	}
	listFbUserQuery := &fbusering.ListFbExternalUserWithCustomerByExternalIDsQuery{
		ShopID:      s.SS.Shop().ID,
		ExternalIDs: fbUserExternalID,
	}
	err = s.FBUserQuery.Dispatch(ctx, listFbUserQuery)
	if err != nil {
		return nil, err
	}
	var mapExternalIDFbUser = make(map[string]*fbusering.FbExternalUserWithCustomer)
	for _, v := range listFbUserQuery.Result {
		mapExternalIDFbUser[v.FbExternalUser.ExternalID] = v
	}
	result := &fabo.FbCustomerConversationsResponse{
		CustomerConversations: convertpb.PbFbCustomerConversations(listCustomerConversations),
		Paging:                cmapi.PbCursorPageInfo(paging, &listCustomerConversationsQuery.Result.Paging),
	}
	for k, v := range result.CustomerConversations {
		if mapExternalIDFbUser[v.ExternalUserID] != nil {
			result.CustomerConversations[k].Customer = convertpb2.PbShopCustomer(mapExternalIDFbUser[v.ExternalUserID].ShopCustomer)
		}
	}
	return result, nil
}

func (s *CustomerConversationService) ListMessages(
	ctx context.Context, request *fabo.ListMessagesRequest,
) (*fabo.FbMessagesResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	faboInfo, err := s.FaboPagesKit.GetPages(ctx, s.SS.Shop().ID)
	if err != nil {
		return nil, err
	}
	var externalConversationIDs []string
	if request.Filter != nil {
		if len(request.Filter.FbExternalConversationIDs) != 0 {
			externalConversationIDs = request.Filter.FbExternalConversationIDs
		} else {
			externalConversationIDs = request.Filter.ExternalConversationID
		}
	}
	listFbExternalMessagesQuery := &fbmessaging.ListFbExternalMessagesQuery{
		ExternalPageIDs:         faboInfo.ExternalPageIDs,
		ExternalConversationIDs: externalConversationIDs,
		Paging:                  *paging,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, listFbExternalMessagesQuery); err != nil {
		return nil, err
	}

	listMessages := listFbExternalMessagesQuery.Result.FbExternalMessages
	// get avatars
	{
		var externalUserIDs []string
		for _, message := range listMessages {
			if message.ExternalFrom != nil {
				externalUserIDs = append(externalUserIDs, message.ExternalFrom.ID)
			}
			for _, externalTo := range message.ExternalTo {
				externalUserIDs = append(externalUserIDs, externalTo.ID)
			}
		}

		mapExternalUserIDAndImageURl, err := s.getImageURLs(ctx, externalUserIDs)
		if err != nil {
			return nil, err
		}

		for _, message := range listMessages {
			if message.ExternalFrom != nil {
				message.ExternalFrom.ImageURL = mapExternalUserIDAndImageURl[message.ExternalFrom.ID]
			}
			for _, externalTo := range message.ExternalTo {
				externalTo.ImageURL = mapExternalUserIDAndImageURl[externalTo.ID]
			}
		}
	}

	return &fabo.FbMessagesResponse{
		FbMessages: convertpb.PbFbExternalMessages(listMessages),
		Paging:     cmapi.PbCursorPageInfo(paging, &listFbExternalMessagesQuery.Result.Paging),
	}, nil
}

func (s *CustomerConversationService) ListCommentsByExternalPostID(
	ctx context.Context, request *fabo.ListCommentsByExternalPostIDRequest,
) (*fabo.ListCommentsByExternalPostIDResponse, error) {
	// TODO: Ngoc add message
	if request.Filter == nil || request.Filter.ExternalPostID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing external_post_id")
	}
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	faboInfo, err := s.FaboPagesKit.GetPages(ctx, s.SS.Shop().ID)
	if err != nil {
		return nil, err
	}
	getFbExternalPostQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: request.Filter.ExternalPostID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalPostQuery); err != nil {
		return nil, err
	}
	fbExternalPost := getFbExternalPostQuery.Result

	haveExternalPageID := false
	for _, externalPageID := range faboInfo.ExternalPageIDs {
		if externalPageID == fbExternalPost.ExternalPageID {
			haveExternalPageID = true
			break
		}
	}
	//// TODO: Ngoc add message
	if !haveExternalPageID {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "")
	}

	listFbExternalCommentsQuery := &fbmessaging.ListFbExternalCommentsQuery{
		FbExternalPostID: request.Filter.ExternalPostID,
		FbExternalUserID: request.Filter.ExternalUserID,
		FbExternalPageID: fbExternalPost.ExternalPageID,
		Paging:           *paging,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, listFbExternalCommentsQuery); err != nil {
		return nil, err
	}

	listComments := listFbExternalCommentsQuery.Result.FbExternalComments
	// get avatars
	{
		var externalUserIDs []string
		for _, comment := range listComments {
			if comment.ExternalFrom != nil {
				externalUserIDs = append(externalUserIDs, comment.ExternalFrom.ID)
			}
		}

		mapExternalUserIDAndImageURL, err := s.getImageURLs(ctx, externalUserIDs)
		if err != nil {
			return nil, err
		}

		for _, comment := range listComments {
			if comment.ExternalFrom != nil {
				comment.ExternalFrom.ImageURL = mapExternalUserIDAndImageURL[comment.ExternalFrom.ID]
			}
		}
	}

	var latestCustomerFbExternalComment *fbmessaging.FbExternalComment
	if request.Filter.ExternalUserID != "" {
		getLatestCustomerExternalCommentQuery := &fbmessaging.GetLatestCustomerExternalCommentQuery{
			ExternalPostID: request.Filter.ExternalPostID,
			ExternalUserID: request.Filter.ExternalUserID,
			ExternalPageID: fbExternalPost.ExternalPageID,
		}
		if err := s.FBMessagingQuery.Dispatch(ctx, getLatestCustomerExternalCommentQuery); err != nil {
			return nil, err
		}
		latestCustomerFbExternalComment = getLatestCustomerExternalCommentQuery.Result
	}

	fbPost := convertpb.PbFbExternalPost(fbExternalPost)

	var commentParentExternalIDs []string
	for _, childrenComment := range listComments {
		if childrenComment.ExternalParentID != "" {
			commentParentExternalIDs = append(commentParentExternalIDs, childrenComment.ExternalParentID)
		}
	}
	listFbExternalCommentsParentQuery := &fbmessaging.ListFbExternalCommentsByExternalIDsQuery{
		ExternalIDs:      commentParentExternalIDs,
		FbExternalUserID: request.Filter.ExternalUserID,
		FbExternalPageID: getFbExternalPostQuery.Result.ExternalPageID,
	}
	if err = s.FBMessagingQuery.Dispatch(ctx, listFbExternalCommentsParentQuery); err != nil {
		return nil, err
	}
	var mapParentComment = make(map[string]*fbmessaging.FbExternalComment)
	for _, v := range listFbExternalCommentsParentQuery.Result.FbExternalComments {
		mapParentComment[v.ExternalID] = v
	}
	// TODO(ngoc): refactor
	// get avatars
	{
		var externalUserIDs []string
		for _, comment := range mapParentComment {
			if comment.ExternalFrom != nil {
				externalUserIDs = append(externalUserIDs, comment.ExternalFrom.ID)
			}
		}

		mapExternalUserIDAndImageURL, err := s.getImageURLs(ctx, externalUserIDs)
		if err != nil {
			return nil, err
		}

		for _, comment := range mapParentComment {
			if comment.ExternalFrom != nil {
				comment.ExternalFrom.ImageURL = mapExternalUserIDAndImageURL[comment.ExternalFrom.ID]
			}
		}
	}

	fbComments := &fabo.FbCommentsResponse{
		FbComments: convertpb.PbFbExternalComments(listComments),
		Paging:     cmapi.PbCursorPageInfo(paging, &listFbExternalCommentsQuery.Result.Paging),
	}
	for k, v := range fbComments.FbComments {
		if v.ExternalParentID != "" {
			fbComments.FbComments[k].ExternalParent = convertpb.PbFbExternalComment(mapParentComment[v.ExternalParentID])
		}
	}
	return &fabo.ListCommentsByExternalPostIDResponse{
		FbComments:                      fbComments,
		LatestCustomerFbExternalComment: convertpb.PbFbExternalComment(latestCustomerFbExternalComment),
		FbPost:                          fbPost,
	}, nil
}

func (s *CustomerConversationService) UpdateReadStatus(
	ctx context.Context, request *fabo.UpdateReadStatusRequest,
) (*common.UpdatedResponse, error) {
	updateIsReadCustomerConversationCmd := &fbmessaging.UpdateIsReadCustomerConversationCommand{
		ConversationCustomerID: request.CustomerConversationID,
		IsRead:                 request.Read,
	}
	if err := s.FBMessagingAggr.Dispatch(ctx, updateIsReadCustomerConversationCmd); err != nil {
		return nil, err
	}
	return &common.UpdatedResponse{
		Updated: updateIsReadCustomerConversationCmd.Result,
	}, nil
}

func (s *CustomerConversationService) SendComment(
	ctx context.Context, request *fabo.SendCommentRequest,
) (*fabo.FbExternalComment, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	externalID := request.ExternalID
	externalPostID := request.ExternalPostID
	externalPageID := request.ExternalPageID

	// Get comment depends on externalID
	getFbExternalCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
		ExternalID: externalID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalCommentQuery); err != nil {
		return nil, err
	}
	comment := getFbExternalCommentQuery.Result

	// Get token depends on externalPageID
	getFbExternalPageInternalActiveQuery := &fbpaging.GetFbExternalPageInternalActiveByExternalIDQuery{
		ExternalID: externalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalActiveQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalActiveQuery.Result.Token

	// Get post depends on externalPostID
	getFbExternalPostByExternalIDQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: externalPostID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalPostByExternalIDQuery); err != nil {
		return nil, err
	}

	// Send comment
	sendCommentRequest := &fbclientmodel.SendCommentArgs{
		ID:            request.ExternalID,
		Message:       request.Message,
		AttachmentURL: request.AttachmentURL,
	}
	sendCommentResponse, err := s.FBClient.CallAPISendComment(&fbclient.SendCommentRequest{
		AccessToken:     accessToken,
		SendCommentArgs: sendCommentRequest,
		PageID:          externalPageID,
	})
	if err != nil {
		return nil, convertApiError(err)
	}

	// Get comment
	var commentID string
	{
		// {page_id}_{post_id}
		postIDParts := strings.Split(request.ExternalPostID, "_")

		// {post_id}_{comment_id}
		commentID = fmt.Sprintf("%s_%s", postIDParts[1], sendCommentResponse.ID)
	}
	newComment, err := s.FBClient.CallAPICommentByID(&fbclient.GetCommentByIDRequest{
		AccessToken: accessToken,
		CommentID:   commentID,
		PageID:      externalPageID,
	})
	if err != nil {
		return nil, err
	}

	var externalParent *fbmessaging.FbObjectParent
	var externalUserID, externalParentID, externalParentUserID string
	if newComment.From != nil {
		externalUserID = newComment.From.ID
	}

	// Make user's comment to parent comment of new comment
	{
		externalParentID = comment.ExternalID
		externalParent = &fbmessaging.FbObjectParent{
			CreatedTime: comment.ExternalCreatedTime,
			Message:     comment.ExternalMessage,
			ID:          externalParentID,
		}
		if comment.ExternalFrom != nil {
			externalParentUserID = comment.ExternalFrom.ID
			externalParent.From = comment.ExternalFrom
		}

	}
	createOrUpdateFbExternalCommentsCmd := &fbmessaging.CreateOrUpdateFbExternalCommentsCommand{
		FbExternalComments: []*fbmessaging.CreateFbExternalCommentArgs{
			{
				ID:                   cm.NewID(),
				ExternalPostID:       externalPostID,
				ExternalPageID:       externalPageID,
				ExternalID:           sendCommentResponse.ID,
				ExternalUserID:       externalUserID,
				ExternalParentID:     externalParentID,
				ExternalParentUserID: externalParentUserID,
				ExternalMessage:      newComment.Message,
				ExternalCommentCount: newComment.CommentCount,
				ExternalParent:       externalParent,
				ExternalFrom:         fbclientconvert.ConvertObjectFrom(newComment.From),
				ExternalAttachment:   fbclientconvert.ConvertFbCommentAttachment(newComment.Attachment),
				ExternalCreatedTime:  newComment.CreatedTime.ToTime(),
				Source:               fb_comment_source.Web,
			},
		},
	}
	if err := s.FBMessagingAggr.Dispatch(ctx, createOrUpdateFbExternalCommentsCmd); err != nil {
		return nil, err
	}

	commentParent := convertpb.PbFbExternalComment(comment)
	result := convertpb.PbFbExternalComment(createOrUpdateFbExternalCommentsCmd.Result[0])
	result.ExternalParent = commentParent
	return result, nil
}

func (s *CustomerConversationService) CreatePost(
	ctx context.Context, request *fabo.CreatePostRequest,
) (*fabo.CreatePostResponse, error) {
	key := fmt.Sprintf("CreatePost %v-%v", request.ExternalPageID, request.Message)
	res, _, err := idempgroup.DoAndWrap(
		ctx, key, 7*60*time.Second, "tạo bài viết",
		func() (interface{}, error) { return s.createPost(ctx, request) })

	if err != nil {
		return nil, err
	}
	return res.(*fabo.CreatePostResponse), err

}

func (s *CustomerConversationService) createPost(
	ctx context.Context, request *fabo.CreatePostRequest,
) (*fabo.CreatePostResponse, error) {
	getFbExternalPageQuery := &fbpaging.GetFbExternalPageByExternalIDQuery{
		ExternalID: request.ExternalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageQuery); err != nil {
		return nil, err
	}
	externalPage := getFbExternalPageQuery.Result

	// Verify permissions
	if err := verifyScopes(appScopes, externalPage.ExternalPermissions); err != nil {
		return nil, err
	}

	getFbExternalPageInternalQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: request.ExternalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalQuery.Result.Token

	createPostCmd := &fbmessaging.CreateFbExternalPostCommand{
		ExternalPageID: request.ExternalPageID,
		Message:        request.Message,
		AccessToken:    accessToken,
	}

	if err := s.FBMessagingAggr.Dispatch(ctx, createPostCmd); err != nil {
		return nil, err
	}

	externalPostID := strings.Split(createPostCmd.Result.ExternalID, "_")[1]

	response := &fabo.CreatePostResponse{
		ExternalPostID: createPostCmd.Result.ExternalID,
		ExternalURL:    fmt.Sprintf("https://www.facebook.com/%s/posts/%s/", request.ExternalPageID, externalPostID),
	}
	return response, nil
}

func (s *CustomerConversationService) SendMessage(
	ctx context.Context, request *fabo.SendMessageRequest,
) (*fabo.FbExternalMessage, error) {
	request.Message.Text = strings.TrimSpace(request.Message.Text)
	request.Message.URL = strings.TrimSpace(request.Message.URL)
	if request.Message == nil || request.Message.Type == "" ||
		(request.Message.URL == "" && request.Message.Text == "") {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing message content")
	}

	getFbExternalPageInternalByIDQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: request.ExternalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalByIDQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalByIDQuery.Result.Token

	getFbExternalConversationQuery := &fbmessaging.GetFbExternalConversationByExternalIDAndExternalPageIDQuery{
		ExternalID:     request.ExternalConversationID,
		ExternalPageID: request.ExternalPageID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalConversationQuery); err != nil {
		return nil, err
	}
	PSID := getFbExternalConversationQuery.Result.PSID

	sendMessageRequest := &fbclientmodel.SendMessageArgs{
		Recipient: &fbclientmodel.RecipientSendMessageRequest{
			ID: PSID,
		},
	}

	if request.Message.Type == "text" {
		sendMessageRequest.Message = &fbclientmodel.MessageSendMessageRequest{
			Text: request.Message.Text,
		}
	} else {
		sendMessageRequest.Message = &fbclientmodel.MessageSendMessageRequest{
			Attachment: &fbclientmodel.AttachmentSendMessageRequest{
				Type: "image",
				Payload: fbclientmodel.PayloadAttachmentSendMessageRequest{
					Url: request.Message.URL,
				},
			},
		}
	}

	sendMessageResp, err := s.FBClient.CallAPISendMessage(&fbclient.SendMessageRequest{
		AccessToken:     accessToken,
		SendMessageArgs: sendMessageRequest,
		PageID:          request.ExternalPageID,
	})
	if err != nil {
		return nil, convertApiError(err)
	}

	newMessage, err := s.FBClient.CallAPIGetMessage(&fbclient.GetMessageRequest{
		AccessToken: accessToken,
		MessageID:   sendMessageResp.MessageID,
		PageID:      request.ExternalPageID,
	})
	if err != nil {
		return nil, err
	}

	var externalAttachments []*fbmessaging.FbMessageAttachment
	if newMessage.Attachments != nil {
		externalAttachments = fbclientconvert.ConvertMessageDataAttachments(newMessage.Attachments.Data)
	}
	createFbExternalMessageCmd := &fbmessaging.CreateOrUpdateFbExternalMessagesCommand{
		FbExternalMessages: []*fbmessaging.CreateFbExternalMessageArgs{
			{
				ID:                     cm.NewID(),
				ExternalConversationID: request.ExternalConversationID,
				ExternalPageID:         request.ExternalPageID,
				ExternalID:             newMessage.ID,
				ExternalMessage:        newMessage.Message,
				ExternalSticker:        newMessage.Sticker,
				ExternalTo:             fbclientconvert.ConvertObjectsTo(newMessage.To),
				ExternalFrom:           fbclientconvert.ConvertObjectFrom(newMessage.From),
				ExternalAttachments:    externalAttachments,
				ExternalCreatedTime:    newMessage.CreatedTime.ToTime(),
			},
		},
	}
	if err := s.FBMessagingAggr.Dispatch(ctx, createFbExternalMessageCmd); err != nil {
		return nil, err
	}

	return convertpb.PbFbExternalMessage(createFbExternalMessageCmd.Result[0]), nil
}

func (s *CustomerConversationService) getImageURLs(ctx context.Context, externalUserIDs []string) (map[string]string, error) {
	getFbUserQuery := &fbusering.ListFbExternalUsersByExternalIDsQuery{
		ExternalIDs: externalUserIDs,
	}
	if err := s.FBUserQuery.Dispatch(ctx, getFbUserQuery); err != nil {
		return nil, err
	}

	mapExternalUserIDAndImageURL := make(map[string]string)
	for _, externalUserID := range externalUserIDs {
		mapExternalUserIDAndImageURL[externalUserID] = ""
	}
	for _, fbExternalUser := range getFbUserQuery.Result {
		if fbExternalUser.ExternalInfo != nil && fbExternalUser.ExternalInfo.ImageURL != "" {
			mapExternalUserIDAndImageURL[fbExternalUser.ExternalID] = fbExternalUser.ExternalInfo.ImageURL
		}
	}

	return mapExternalUserIDAndImageURL, nil
}

// TODO(nakhoa17): make function more clear
func convertApiError(err error) error {
	apiErr, ok := err.(*xerrors.APIError)
	if !ok {
		return err
	}

	metaError := apiErr.Meta
	subCode, ok := metaError["sub_code"]
	if !ok {
		return err
	}

	intSubCode, _err := strconv.Atoi(subCode)
	if _err != nil {
		return err
	}

	switch intSubCode {
	case int(fbclient.ObjectNotExist):
		return cm.Errorf(cm.FacebookError, nil, "Cuộc hội thoại này không tồn tại, hoặc đã bị xóa.").
			WithMetaM(metaError)
	case int(fbclient.MessageSentOutside):
		return cm.Errorf(cm.FacebookError, nil, "Người nhận chưa reply trong vòng 24 giờ nên không thể gửi thêm tin nhắn").
			WithMetaM(metaError)
	case int(fbclient.Expired):
		return cm.Errorf(cm.FacebookError, nil, "Token truy cập facebook trang của bạn đã hết hạn.").
			WithMetaM(metaError)
	default:
		return err
	}
}
