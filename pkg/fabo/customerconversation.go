package fabo

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/common"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/backend/pkg/fabo/faboinfo"
)

type CustomerConversationService struct {
	session.Session

	faboInfo         *faboinfo.FaboInfo
	fbMessagingQuery fbmessaging.QueryBus
	fbMessagingAggr  fbmessaging.CommandBus
	fbPagingQuery    fbpaging.QueryBus
	fbClient         *fbclient.FbClient
}

func NewCustomerConversationService(
	ss session.Session,
	faboInfo *faboinfo.FaboInfo,
	fbMessagingQuery fbmessaging.QueryBus,
	fbMessagingAggr fbmessaging.CommandBus,
	fbPagingQuery fbpaging.QueryBus,
	fbClient *fbclient.FbClient,
) *CustomerConversationService {
	s := &CustomerConversationService{
		Session:          ss,
		faboInfo:         faboInfo,
		fbMessagingQuery: fbMessagingQuery,
		fbMessagingAggr:  fbMessagingAggr,
		fbPagingQuery:    fbPagingQuery,
		fbClient:         fbClient,
	}
	return s
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
			if err := s.fbPagingQuery.Dispatch(ctx, listFbExternalPagesByIDsQuery); err != nil {
				return nil, err
			}
			for _, fbExternalPage := range listFbExternalPagesByIDsQuery.Result {
				externalPageIDsRequest = append(externalPageIDsRequest, fbExternalPage.ExternalID)
			}
		}

		if len(fbPageIDsRequest) == 0 && len(externalPageIDsRequest) == 0 {
			faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.SS.Shop().ID)
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
		faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.SS.Shop().ID)
		if err != nil {
			return nil, err
		}
		listCustomerConversationsQuery.ExternalPageIDs = faboInfo.ExternalPageIDs
	}
	if err := s.fbMessagingQuery.Dispatch(ctx, listCustomerConversationsQuery); err != nil {
		return nil, err
	}

	return &fabo.FbCustomerConversationsResponse{
		CustomerConversations: convertpb.PbFbCustomerConversations(listCustomerConversationsQuery.Result.FbCustomerConversations),
		Paging:                cmapi.PbCursorPageInfo(paging, &listCustomerConversationsQuery.Result.Paging),
	}, nil
}

func (s *CustomerConversationService) ListMessages(
	ctx context.Context, request *fabo.ListMessagesRequest,
) (*fabo.FbMessagesResponse, error) {
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.SS.Shop().ID)
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
	if err := s.fbMessagingQuery.Dispatch(ctx, listFbExternalMessagesQuery); err != nil {
		return nil, err
	}

	return &fabo.FbMessagesResponse{
		FbMessages: convertpb.PbFbExternalMessages(listFbExternalMessagesQuery.Result.FbExternalMessages),
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
	faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.SS.Shop().ID)
	if err != nil {
		return nil, err
	}
	getFbExternalPostQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: request.Filter.ExternalPostID,
	}
	if err := s.fbMessagingQuery.Dispatch(ctx, getFbExternalPostQuery); err != nil {
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
		FbExternalPageID: getFbExternalPostQuery.Result.ExternalPageID,
		Paging:           *paging,
	}
	if err := s.fbMessagingQuery.Dispatch(ctx, listFbExternalCommentsQuery); err != nil {
		return nil, err
	}

	var latestCustomerFbExternalComment *fbmessaging.FbExternalComment
	if request.Filter.ExternalUserID != "" {
		getLatestCustomerExternalCommentQuery := &fbmessaging.GetLatestCustomerExternalCommentQuery{
			ExternalPostID: request.Filter.ExternalPostID,
			ExternalUserID: request.Filter.ExternalUserID,
		}
		if err := s.fbMessagingQuery.Dispatch(ctx, getLatestCustomerExternalCommentQuery); err != nil {
			return nil, err
		}
		latestCustomerFbExternalComment = getLatestCustomerExternalCommentQuery.Result
	}

	parentID := getFbExternalPostQuery.Result.ExternalParentID
	queryParent := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: parentID,
	}
	if err = s.fbMessagingQuery.Dispatch(ctx, queryParent); err != nil {
		return nil, err
	}
	fbExternalParentPost := getFbExternalPostQuery.Result
	fbPost := convertpb.PbFbExternalPost(fbExternalPost)
	fbPost.ExternalParent = convertpb.PbFbExternalPost(fbExternalParentPost)

	return &fabo.ListCommentsByExternalPostIDResponse{
		FbComments: &fabo.FbCommentsResponse{
			FbComments: convertpb.PbFbExternalComments(listFbExternalCommentsQuery.Result.FbExternalComments),
			Paging:     cmapi.PbCursorPageInfo(paging, &listFbExternalCommentsQuery.Result.Paging),
		},
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
	if err := s.fbMessagingAggr.Dispatch(ctx, updateIsReadCustomerConversationCmd); err != nil {
		return nil, err
	}
	return &common.UpdatedResponse{
		Updated: updateIsReadCustomerConversationCmd.Result,
	}, nil
}

func (s *CustomerConversationService) SendComment(
	ctx context.Context, request *fabo.SendCommentRequest,
) (*fabo.FbExternalComment, error) {
	if request.ExternalPageID == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "missing external_page_id")
	}
	if request.ExternalID == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "missing external_id")
	}
	if request.ExternalPostID == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "missing external_post_id")
	}
	if request.Message == "" && request.AttachmentURL == "" {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "missing content")
	}

	getFbExternalPageInternalQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: request.ExternalPageID,
	}
	if err := s.fbPagingQuery.Dispatch(ctx, getFbExternalPageInternalQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalQuery.Result.Token

	getFbExternalPostByExternalIDQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: request.ExternalPostID,
	}
	if err := s.fbMessagingQuery.Dispatch(ctx, getFbExternalPostByExternalIDQuery); err != nil {
		return nil, err
	}

	sendCommentRequest := &fbclientmodel.SendCommentRequest{
		ID:            request.ExternalID,
		Message:       request.Message,
		AttachmentURL: request.AttachmentURL,
	}
	sendCommentResponse, err := s.fbClient.CallAPISendComment(accessToken, sendCommentRequest)
	if err != nil {
		return nil, err
	}

	newComment, err := s.fbClient.CallAPICommentByID(accessToken, sendCommentResponse.ID)
	if err != nil {
		return nil, err
	}

	var externalUserID, externalParentID, externalParentUserID string
	if newComment.From != nil {
		externalUserID = newComment.From.ID
	}
	if newComment.Parent != nil {
		externalParentID = newComment.Parent.ID
		if newComment.Parent.From != nil {
			externalParentUserID = newComment.Parent.From.ID
		}
	}
	createOrUpdateFbExternalCommentsCmd := &fbmessaging.CreateOrUpdateFbExternalCommentsCommand{
		FbExternalComments: []*fbmessaging.CreateFbExternalCommentArgs{
			{
				ID:                   cm.NewID(),
				ExternalPostID:       request.ExternalPostID,
				ExternalPageID:       request.ExternalPageID,
				ExternalID:           sendCommentResponse.ID,
				ExternalUserID:       externalUserID,
				ExternalParentID:     externalParentID,
				ExternalParentUserID: externalParentUserID,
				ExternalMessage:      newComment.Message,
				ExternalCommentCount: newComment.CommentCount,
				ExternalParent:       fbclientconvert.ConvertFbObjectParent(newComment.Parent),
				ExternalFrom:         fbclientconvert.ConvertObjectFrom(newComment.From),
				ExternalAttachment:   fbclientconvert.ConvertFbCommentAttachment(newComment.Attachment),
				ExternalCreatedTime:  newComment.CreatedTime.ToTime(),
			},
		},
	}

	if err := s.fbMessagingAggr.Dispatch(ctx, createOrUpdateFbExternalCommentsCmd); err != nil {
		return nil, err
	}

	return convertpb.PbFbExternalComment(createOrUpdateFbExternalCommentsCmd.Result[0]), nil
}

func (s *CustomerConversationService) SendMessage(
	ctx context.Context, request *fabo.SendMessageRequest,
) (*fabo.FbExternalMessage, error) {
	if request.Message == nil || request.Message.Type == "" ||
		(request.Message.URL == "" && request.Message.Text == "") {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "")
	}

	getFbExternalPageInternalByIDQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: request.ExternalPageID,
	}
	if err := s.fbPagingQuery.Dispatch(ctx, getFbExternalPageInternalByIDQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalByIDQuery.Result.Token

	getFbExternalConversationQuery := &fbmessaging.GetFbExternalConversationByExternalIDAndExternalPageIDQuery{
		ExternalID:     request.ExternalConversationID,
		ExternalPageID: request.ExternalPageID,
	}
	if err := s.fbMessagingQuery.Dispatch(ctx, getFbExternalConversationQuery); err != nil {
		return nil, err
	}
	PSID := getFbExternalConversationQuery.Result.PSID

	sendMessageRequest := &fbclientmodel.SendMessageRequest{
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

	sendMessageResp, err := s.fbClient.CallAPISendMessage(accessToken, sendMessageRequest)
	if err != nil {
		return nil, err
	}

	newMessage, err := s.fbClient.CallAPIGetMessage(accessToken, sendMessageResp.MessageID)
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
	if err := s.fbMessagingAggr.Dispatch(ctx, createFbExternalMessageCmd); err != nil {
		return nil, err
	}

	return convertpb.PbFbExternalMessage(createFbExternalMessageCmd.Result[0]), nil
}
