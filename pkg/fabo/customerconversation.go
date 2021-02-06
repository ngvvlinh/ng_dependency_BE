package fabo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"o.o/api/fabo/fbcustomerconversationsearch"
	"o.o/api/fabo/fbmessagetemplate"
	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_comment_action"
	"o.o/api/fabo/fbmessaging/fb_comment_source"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbmessaging/fb_internal_source"
	"o.o/api/fabo/fbmessaging/fb_post_type"
	"o.o/api/fabo/fbmessaging/fb_status_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/common"
	"o.o/backend/com/fabo/pkg/fbclient"
	fbclientconvert "o.o/backend/com/fabo/pkg/fbclient/convert"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/validate"
	convertxmin "o.o/backend/pkg/etop/apix/convertpb/_min"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/backend/pkg/fabo/faboinfo"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

var messageTemplateVariables = fabo.MessageTemplateVariableResponse{
	Variables: []*fabo.MessageTemplateVariable{
		{
			Code:  "$$fb_username",
			Label: "Tên khách hàng",
		},
		{
			Code:  "$$page_name",
			Label: "Tên page",
		},
		{
			Code:  "$$shop_name",
			Label: "Tên cửa hàng",
		},
	},
}

type CustomerConversationService struct {
	session.Session

	FaboPagesKit *faboinfo.FaboPagesKit
	FBClient     *fbclient.FbClient

	FBMessagingQuery       fbmessaging.QueryBus
	FBMessagingAggr        fbmessaging.CommandBus
	FBPagingQuery          fbpaging.QueryBus
	FBUserQuery            fbusering.QueryBus
	FbSearchQuery          fbcustomerconversationsearch.QueryBus
	FbMessageTemplateQuery fbmessagetemplate.QueryBus
	FbMessageTemplateAggr  fbmessagetemplate.CommandBus
}

type APIType string

const (
	APIComment           APIType = "comment"
	APIMessage           APIType = "message"
	CanNotCommentMsg             = "Không thể reply trên cuộc hội thoại này, có thể bài post hoặc comment này trên facebook đã bị xóa."
	SendMessageOutSide           = "Người nhận chưa reply trong vòng 24 giờ nên không thể gửi thêm tin nhắn"
	PersonNotAvailable           = "Có thể user này đã block page của bạn."
	ExpiredToken                 = "Token truy cập facebook trang của bạn đã hết hạn."
	UserCanNotReply              = "Không thể nhắn tin cho user này."
	AlreadyRepliedToUser         = "Bạn đã phản hồi nên không thể gửi tin nhắn"
	ReplyingTimeExpired          = "Đã quá thời hạn (7 ngày) kể từ khi KH bình luận, bạn không thể gửi tin nhắn cho KH này."
)

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
		externalPageIDsRequest := request.Filter.ExternalPageID

		if len(externalPageIDsRequest) == 0 {
			faboInfo, err := s.FaboPagesKit.GetPages(ctx, s.SS.Shop().ID)
			if err != nil {
				return nil, err
			}
			externalPageIDsRequest = faboInfo.ExternalPageIDs
		}
		listCustomerConversationsQuery.ExternalPageIDs = externalPageIDsRequest
		listCustomerConversationsQuery.IsRead = request.Filter.IsRead

		if request.Filter.ExternalUserID.Valid {
			listCustomerConversationsQuery.ExternalUserID = request.Filter.ExternalUserID
		}

		if request.Filter.Type.Valid {
			listCustomerConversationsQuery.Types = []fb_customer_conversation_type.FbCustomerConversationType{request.Filter.Type.Enum}
		} else {
			listCustomerConversationsQuery.Types = request.Filter.Types
		}
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
	customerConversations, err := s.buildFbCustomerConversations(ctx, listCustomerConversations)
	if err != nil {
		return nil, err
	}
	result := &fabo.FbCustomerConversationsResponse{
		CustomerConversations: customerConversations,
		Paging:                cmapi.PbCursorPageInfo(paging, &listCustomerConversationsQuery.Result.Paging),
	}
	return result, nil
}

// Api for full text search
func (s *CustomerConversationService) SearchCustomerConversations(
	ctx context.Context,
	request *fabo.SearchCustomerConversationRequest,
) (*fabo.SearchFbCustomerConversationsResponse, error) {
	textToSearch := request.Text
	shopID := s.SS.Shop().ID
	var customerConversations []*fbmessaging.FbCustomerConversation
	if textToSearch == "" {
		response := &fabo.SearchFbCustomerConversationsResponse{}
		return response, nil
	}

	pagesQuery := &fbpaging.ListActiveFbPagesByShopIDsQuery{
		ShopIDs: []dot.ID{shopID},
	}
	if err := s.FBPagingQuery.Dispatch(ctx, pagesQuery); err != nil {
		return nil, err
	}
	var extPageIDs []string
	for _, extPage := range pagesQuery.Result {
		extPageIDs = append(extPageIDs, extPage.ExternalID)
	}
	if len(extPageIDs) == 0 {
		return &fabo.SearchFbCustomerConversationsResponse{}, nil
	}

	// search by phone
	conversationsByPhone, err := s.searchCustomerConversationsByPhone(ctx, shopID, textToSearch)
	if err != nil {
		return nil, err
	}
	customerConversations = append(customerConversations, conversationsByPhone...)

	// search by external text message
	conversationsByExternalMessage, err := s.searchCustomerConversationByExternalMessage(ctx, extPageIDs, textToSearch)
	if err != nil {
		return nil, err
	}
	customerConversations = append(customerConversations, conversationsByExternalMessage...)

	// search by external user name
	conversationsByExternalUserName, err := s.searchByExternalUserName(ctx, extPageIDs, textToSearch)
	if err != nil {
		return nil, err
	}
	customerConversations = append(customerConversations, conversationsByExternalUserName...)

	// combine them together & return
	customerConversations = distinctCustomerConversations(customerConversations)
	faboCustomerConversations, err := s.buildFbCustomerConversations(ctx, customerConversations)
	if err != nil {
		return nil, err
	}
	response := &fabo.SearchFbCustomerConversationsResponse{
		CustomerConversations: faboCustomerConversations,
	}
	return response, nil
}

func (s *CustomerConversationService) searchCustomerConversationsByPhone(
	ctx context.Context, shopID dot.ID, phone string,
) ([]*fbmessaging.FbCustomerConversation, error) {
	normalizedPhone, isPhone := validate.NormalizePhone(phone)
	if !isPhone {
		return nil, nil
	}

	customerQuery := &fbusering.ListShopCustomerIDWithPhoneNormQuery{
		ShopID: shopID,
		Phone:  string(normalizedPhone),
	}
	if err := s.FBUserQuery.Dispatch(ctx, customerQuery); err != nil {
		return nil, err
	}

	customerIDs := customerQuery.Result
	queryFbExternalUserIDs := &fbusering.ListFbExternalUserIDsByShopCustomerIDsQuery{
		CustomerIDs: customerIDs,
	}
	if err := s.FBUserQuery.Dispatch(ctx, queryFbExternalUserIDs); err != nil {
		return nil, err
	}

	fbExtUserIDs := queryFbExternalUserIDs.Result
	customerByExternalUserIDsQuery := &fbmessaging.ListFbCustomerConversationsByExternalUserIDsQuery{
		ExtUserIDs: fbExtUserIDs,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, customerByExternalUserIDsQuery); err != nil {
		return nil, err
	}
	return customerByExternalUserIDsQuery.Result, nil
}

func (s *CustomerConversationService) searchCustomerConversationByExternalMessage(
	ctx context.Context, pageIDs []string, extMessage string,
) ([]*fbmessaging.FbCustomerConversation, error) {
	var result []*fbmessaging.FbCustomerConversation

	conversationsInComment, err := s.searchInComment(ctx, pageIDs, extMessage)
	if err != nil {
		return nil, err
	}
	result = append(result, conversationsInComment...)

	conversationsInMessage, err := s.searchInMessage(ctx, pageIDs, extMessage)
	if err != nil {
		return nil, err
	}
	result = append(result, conversationsInMessage...)

	return result, nil
}

func (s *CustomerConversationService) searchByExternalUserName(
	ctx context.Context, pageIDs []string, extUserName string,
) ([]*fbmessaging.FbCustomerConversation, error) {
	convSearchQuery := &fbcustomerconversationsearch.ListFbExternalConversationSearchQuery{
		PageIDs:     pageIDs,
		ExtUserName: extUserName,
	}
	if err := s.FbSearchQuery.Dispatch(ctx, convSearchQuery); err != nil {
		return nil, err
	}

	var conversationIDs []dot.ID
	for _, conv := range convSearchQuery.Result {
		conversationIDs = append(conversationIDs, conv.ID)
	}

	customerConversationsQuery := &fbmessaging.ListFbCustomerConversationsByIDsQuery{
		IDs: conversationIDs,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, customerConversationsQuery); err != nil {
		return nil, err
	}

	return customerConversationsQuery.Result, nil
}

func (s *CustomerConversationService) searchInComment(
	ctx context.Context, pageIDs []string, extMessage string,
) ([]*fbmessaging.FbCustomerConversation, error) {
	commentsQuery := &fbcustomerconversationsearch.ListFbExternalCommentSearchQuery{
		PageIDs:     pageIDs,
		ExternalMsg: extMessage,
	}
	if err := s.FbSearchQuery.Dispatch(ctx, commentsQuery); err != nil {
		return nil, err
	}
	var (
		postIDs    []string
		extUserIDs []string
	)
	for _, cmt := range commentsQuery.Result {
		postIDs = append(postIDs, cmt.ExternalPostID)
		extUserIDs = append(extUserIDs, cmt.ExternalUserID)
	}

	customerConversationQuery := &fbmessaging.ListFbCustomerConversationsByExtUserIDsAndExtIDsQuery{
		ExtUserIDs: extUserIDs,
		ExtIDs:     postIDs,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, customerConversationQuery); err != nil {
		return nil, err
	}

	return customerConversationQuery.Result, nil
}

func (s *CustomerConversationService) searchInMessage(
	ctx context.Context, pageIDs []string, extMessage string,
) ([]*fbmessaging.FbCustomerConversation, error) {
	messagesQuery := &fbcustomerconversationsearch.ListFbExternalMessageSearchQuery{
		PageIDs:     pageIDs,
		ExternalMsg: extMessage,
	}
	if err := s.FbSearchQuery.Dispatch(ctx, messagesQuery); err != nil {
		return nil, err
	}
	var externalConvIDs []string
	for _, msg := range messagesQuery.Result {
		externalConvIDs = append(externalConvIDs, msg.ExternalConversationID)
	}

	customerConversationQuery := &fbmessaging.ListFbCustomerConversationsByExternalIDsQuery{
		ExternalIDs: externalConvIDs,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, customerConversationQuery); err != nil {
		return nil, err
	}

	return customerConversationQuery.Result, nil
}

func (s *CustomerConversationService) GetCustomerConversationByID(
	ctx context.Context,
	request *fabo.GetCustomerConversationByIDRequest,
) (*fabo.GetCustomerConversationByIDResponse, error) {
	query := &fbmessaging.GetFbCustomerConversationByIDQuery{
		ID: request.ID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	// get and map avatar for user in conversation
	conversation := query.Result
	mapExternalUserIDAndImageURl, err := s.buildMapFbUserAvatar(ctx, []string{conversation.ExternalUserID})
	if err != nil {
		return nil, err
	}
	if conversation.ExternalFrom != nil {
		conversation.ExternalFrom.ImageURL = mapExternalUserIDAndImageURl[conversation.ExternalFrom.ID]
	}
	if conversation.ExternalUserID != "" {
		conversation.ExternalUserPictureURL = mapExternalUserIDAndImageURl[conversation.ExternalUserID]
	}
	return &fabo.GetCustomerConversationByIDResponse{
		Conversation: convertpb.PbFbCustomerConversation(conversation),
	}, nil
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

		mapExternalUserIDAndImageURl, err := s.buildMapFbUserAvatar(ctx, externalUserIDs)
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

	listFbExternalCommentsQuery := &fbmessaging.ListFbExternalCommentsQuery{
		FbExternalPostID: request.Filter.ExternalPostID,
		FbExternalUserID: request.Filter.ExternalUserID,
		Paging:           *paging,
	}

	if fbExternalPost.Type == fb_post_type.Page || fbExternalPost.Type == fb_post_type.Unknown { // backward compatible
		haveExternalPageID := false
		for _, externalPageID := range faboInfo.ExternalPageIDs {
			if externalPageID == fbExternalPost.ExternalPageID {
				haveExternalPageID = true
				break
			}
		}
		if !haveExternalPageID {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "external_page_ids not found")
		}
		listFbExternalCommentsQuery.FbExternalPageID = fbExternalPost.ExternalPageID
	} else {
		getFbExternalUserConnectedQuery := &fbusering.GetFbExternalUserConnectedByShopIDQuery{
			ShopID: s.SS.Shop().ID,
		}
		if err := s.FBUserQuery.Dispatch(ctx, getFbExternalUserConnectedQuery); err != nil {
			return nil, err
		}
		fbExternalUserConnected := getFbExternalUserConnectedQuery.Result
		listFbExternalCommentsQuery.FbExternalOwnerPostID = fbExternalUserConnected.ExternalID
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

		mapExternalUserIDAndImageURL, err := s.buildMapFbUserAvatar(ctx, externalUserIDs)
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

		mapExternalUserIDAndImageURL, err := s.buildMapFbUserAvatar(ctx, externalUserIDs)
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

	extCommentID := request.ExternalID
	requestExtUserID := request.ExternalUserID
	externalPostID := request.ExternalPostID
	externalPageID := request.ExternalPageID

	// Get post depends on externalPostID
	getFbExternalPostByExternalIDQuery := &fbmessaging.GetFbExternalPostByExternalIDQuery{
		ExternalID: externalPostID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalPostByExternalIDQuery); err != nil {
		return nil, err
	}
	if !getFbExternalPostByExternalIDQuery.Result.DeletedAt.IsZero() {
		return nil, cm.Error(cm.NotFound, CanNotCommentMsg, nil)
	}

	var comment *fbmessaging.FbExternalComment
	if requestExtUserID != "" {
		getFbExternalCommentQuery := &fbmessaging.GetLatestUpdateActiveCommentQuery{
			ExtPostID: externalPostID,
			ExtUserID: requestExtUserID,
		}
		if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalCommentQuery); err != nil {
			return nil, err
		}
		comment = getFbExternalCommentQuery.Result
	} else { // backward compatible
		getFbExternalCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
			ExternalID: extCommentID,
		}
		if err := s.FBMessagingQuery.Dispatch(ctx, getFbExternalCommentQuery); err != nil {
			return nil, err
		}
		comment = getFbExternalCommentQuery.Result
	}
	if !comment.DeletedAt.IsZero() {
		return nil, cm.Error(cm.NotFound, CanNotCommentMsg, nil)
	}

	// Get token depends on externalPageID
	getFbExternalPageInternalActiveQuery := &fbpaging.GetFbExternalPageInternalActiveByExternalIDQuery{
		ExternalID: externalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalActiveQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalActiveQuery.Result.Token

	sendCommentRequest := &fbclientmodel.SendCommentArgs{
		ID:            comment.ExternalID,
		Message:       request.Message,
		AttachmentURL: request.AttachmentURL,
	}
	sendCommentResponse, err := s.FBClient.CallAPISendComment(&fbclient.SendCommentRequest{
		AccessToken:     accessToken,
		SendCommentArgs: sendCommentRequest,
		PageID:          externalPageID,
	})
	if err != nil {
		return nil, s.handleAndConvertFacebookApiError(ctx, APIComment, comment.ExternalID, err)
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
		return nil, s.handleAndConvertFacebookApiError(ctx, APIComment, commentID, err)
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
				InternalSource:       fb_internal_source.Fabo, // message through our api, set it is `Fabo`
				CreatedBy:            s.SS.User().ID,
				PostType:             fb_post_type.Page,
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
	// TODO(khoa): check nil request.Message
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
		return nil, s.handleAndConvertFacebookApiError(ctx, APIMessage, PSID, err)
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
				InternalSource:         fb_internal_source.Fabo, // message through our api, set it is `Fabo`
				CreatedBy:              s.SS.User().ID,
			},
		},
	}
	if err := s.FBMessagingAggr.Dispatch(ctx, createFbExternalMessageCmd); err != nil {
		return nil, err
	}

	return convertpb.PbFbExternalMessage(createFbExternalMessageCmd.Result[0]), nil
}

func (s *CustomerConversationService) MessageTemplateVariables(ctx context.Context, req *common.Empty) (*fabo.MessageTemplateVariableResponse, error) {
	return &messageTemplateVariables, nil
}

func (s *CustomerConversationService) MessageTemplates(ctx context.Context, req *common.Empty) (*fabo.MessageTemplateResponse, error) {
	query := &fbmessagetemplate.GetMessageTemplatesQuery{
		ShopID: s.SS.Shop().ID,
	}
	if err := s.FbMessageTemplateQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	return &fabo.MessageTemplateResponse{Templates: convertpb.FbMessageTemplates(query.Result)}, nil
}

func (s *CustomerConversationService) CreateMessageTemplate(ctx context.Context, req *fabo.CreateMessageTemplateRequest) (*fabo.MessageTemplate, error) {
	if req.Template == "" {
		return nil, cm.Error(cm.InvalidArgument, "invalid template", nil)
	}

	cmd := &fbmessagetemplate.CreateMessageTemplateCommand{
		ShopID:    s.SS.Shop().ID,
		Template:  req.Template,
		ShortCode: req.ShortCode,
	}
	if err := s.FbMessageTemplateAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return convertpb.FbMessageTemplate(cmd.Result), nil
}

func (s *CustomerConversationService) UpdateMessageTemplate(ctx context.Context, req *fabo.UpdateMessageTemplateRequest) (*common.Empty, error) {
	cmd := &fbmessagetemplate.UpdateMessageTemplateCommand{
		ID:        req.ID,
		ShopID:    s.SS.Shop().ID,
		Template:  req.Template,
		ShortCode: req.ShortCode,
	}
	if err := s.FbMessageTemplateAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &common.Empty{}, nil
}

func (s *CustomerConversationService) DeleteMessageTemplate(ctx context.Context, req *fabo.DeleteMessageTemplateRequest) (*common.Empty, error) {
	cmd := &fbmessagetemplate.DeleteMessageTemplateCommand{
		ShopID: s.SS.Shop().ID,
		ID:     req.ID,
	}
	if err := s.FbMessageTemplateAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (s *CustomerConversationService) buildMapFbUserAvatar(ctx context.Context, externalUserIDs []string) (map[string]string, error) {
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

func (s *CustomerConversationService) handleAndConvertFacebookApiError(ctx context.Context, _type APIType, extID string, err error) error {
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
		removeCommentArgs := &fbmessaging.RemoveCommentCommand{
			ExternalCommentID: extID,
		}
		_ = s.FBMessagingAggr.Dispatch(ctx, removeCommentArgs)
		return cm.Errorf(cm.FacebookError, nil, CanNotCommentMsg).
			WithMetaM(metaError)
	case int(fbclient.MessageSentOutside):
		return cm.Errorf(cm.FacebookError, nil, SendMessageOutSide).
			WithMetaM(metaError)
	case int(fbclient.PersonNotAvailable):
		return cm.Errorf(cm.FacebookError, nil, PersonNotAvailable).
			WithMetaM(metaError)
	case int(fbclient.Expired):
		return cm.Errorf(cm.FacebookError, nil, ExpiredToken).
			WithMetaM(metaError)
	default:
		return err
	}
}

func (s *CustomerConversationService) LikeOrUnLikeComment(
	ctx context.Context, req *fabo.LikeOrUnLikeCommentRequest,
) (*common.Empty, error) {
	if req.Action != fb_comment_action.Like && req.Action != fb_comment_action.UnLike {
		return nil, cm.Error(cm.FailedPrecondition, "Action không hợp lệ", nil)
	}

	getFbCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
		ExternalID: req.ExternalCommentID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbCommentQuery); err != nil {
		return nil, err
	}
	fbComment := getFbCommentQuery.Result

	if fbComment.ExternalFrom != nil && fbComment.ExternalFrom.ID == fbComment.ExternalPageID {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể %v comment cuả page", req.Action.Label())
	}

	getFbExternalPageInternalByIDQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: req.ExternalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalByIDQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalByIDQuery.Result.Token

	if req.Action == fb_comment_action.Like {
		if _, err := s.FBClient.CallAPILikeComment(&fbclient.LikeCommentRequest{
			AccessToken: accessToken,
			PageID:      req.ExternalPageID,
			CommentID:   req.ExternalCommentID,
		}); err != nil {
			return nil, err
		}
	} else {
		if _, err := s.FBClient.CallAPIUnLikeComment(&fbclient.UnLikeCommentRequest{
			AccessToken: accessToken,
			PageID:      req.ExternalPageID,
			CommentID:   req.ExternalCommentID,
		}); err != nil {
			return nil, err
		}
	}

	likeOrUnLikeCommentCmd := &fbmessaging.LikeOrUnLikeCommentCommand{
		ExternalCommentID: req.ExternalCommentID,
	}
	if req.Action == fb_comment_action.Like {
		likeOrUnLikeCommentCmd.IsLiked = true
	}

	if err := s.FBMessagingAggr.Dispatch(ctx, likeOrUnLikeCommentCmd); err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (s *CustomerConversationService) HideOrUnHideComment(
	ctx context.Context, req *fabo.HideOrUnHideCommentRequest,
) (*common.Empty, error) {
	if req.Action != fb_comment_action.Hide && req.Action != fb_comment_action.UnHide {
		return nil, cm.Error(cm.FailedPrecondition, "Action không hợp lệ", nil)
	}

	getFbCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
		ExternalID: req.ExternalCommentID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbCommentQuery); err != nil {
		return nil, err
	}
	fbComment := getFbCommentQuery.Result

	if fbComment.ExternalFrom != nil && fbComment.ExternalFrom.ID == fbComment.ExternalPageID {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể %v comment của page", req.Action.Label())
	}

	getFbExternalPageInternalByIDQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: req.ExternalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalByIDQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalByIDQuery.Result.Token

	hideOrUnHideCommentReq := &fbclient.HideOrUnHideCommentRequest{
		AccessToken: accessToken,
		PageID:      req.ExternalPageID,
		CommentID:   req.ExternalCommentID,
	}
	if req.Action == fb_comment_action.Hide {
		hideOrUnHideCommentReq.IsHidden = true
	}
	if _, err := s.FBClient.CallAPIHideAndUnHideComment(hideOrUnHideCommentReq); err != nil {
		return nil, err
	}

	hideOrUnHideCommentCmd := &fbmessaging.HideOrUnHideCommentCommand{
		ExternalCommentID: req.ExternalCommentID,
	}
	if req.Action == fb_comment_action.Hide {
		hideOrUnHideCommentCmd.IsHidden = true
	}
	if err := s.FBMessagingAggr.Dispatch(ctx, hideOrUnHideCommentCmd); err != nil {
		return nil, err
	}

	return &common.Empty{}, nil
}

func (s *CustomerConversationService) SendPrivateReply(
	ctx context.Context, req *fabo.SendPrivateReplyRequest,
) (*common.Empty, error) {
	if req.Message == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Nội dung tin nhắn ko được để trống")
	}

	getFbCommentQuery := &fbmessaging.GetFbExternalCommentByExternalIDQuery{
		ExternalID: req.ExternalCommentID,
	}
	if err := s.FBMessagingQuery.Dispatch(ctx, getFbCommentQuery); err != nil {
		return nil, err
	}
	fbComment := getFbCommentQuery.Result

	if fbComment.ExternalFrom != nil && fbComment.ExternalFrom.ID == fbComment.ExternalPageID {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Không thể gửi tin nhắn cho page")
	}

	getFbExternalPageInternalByIDQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
		ExternalID: req.ExternalPageID,
	}
	if err := s.FBPagingQuery.Dispatch(ctx, getFbExternalPageInternalByIDQuery); err != nil {
		return nil, err
	}
	accessToken := getFbExternalPageInternalByIDQuery.Result.Token

	sendPrivateReplyReq := &fbclient.SendMessageRequest{
		AccessToken: accessToken,
		SendMessageArgs: &fbclientmodel.SendMessageArgs{
			Recipient: &fbclientmodel.RecipientSendMessageRequest{
				CommentID: req.ExternalCommentID,
			},
			Message: &fbclientmodel.MessageSendMessageRequest{
				Text: req.Message,
			},
		},
		PageID: req.ExternalPageID,
	}

	var changeIsPrivateReplies bool
	_, err := s.FBClient.CallAPISendMessage(sendPrivateReplyReq)
	changeIsPrivateReplies, err = s.handleErrorWhenSendPrivateReplies(err)
	if err != nil && !changeIsPrivateReplies {
		return nil, err
	}

	if changeIsPrivateReplies {
		updateIsPrivateRepliesCmd := &fbmessaging.UpdateIsPrivateRepliedCommentCommand{
			ExternalCommentID: req.ExternalCommentID,
			IsPrivateReplied:  true,
		}
		if err := s.FBMessagingAggr.Dispatch(ctx, updateIsPrivateRepliesCmd); err != nil {
			return nil, err
		}
	}

	return &common.Empty{}, nil
}

func (s *CustomerConversationService) ListLiveVideos(
	ctx context.Context, req *fabo.ListLiveVideosRequest,
) (*fabo.ListLiveVideosResponse, error) {
	var filterExternalPageIDs, externalPageIDsArgs, externalPageIDs []string

	paging, err := cmapi.CMCursorPaging(req.Paging)
	if err != nil {
		return nil, err
	}

	getFbExternalUserConnectedQuery := &fbusering.GetFbExternalUserConnectedByShopIDQuery{
		ShopID: s.SS.Shop().ID,
	}
	if err := s.FBUserQuery.Dispatch(ctx, getFbExternalUserConnectedQuery); err != nil {
		return nil, err
	}
	externalUserID := getFbExternalUserConnectedQuery.Result.ExternalID

	faboInfo, err := s.FaboPagesKit.GetPages(ctx, s.SS.Shop().ID)
	if err != nil {
		return nil, err
	}

	listFbExternalPostsQuery := &fbmessaging.ListFbExternalPostsQuery{
		ExternalStatusType: fb_status_type.AddedVideo.Wrap(),
		Paging:             *paging,
	}

	if req.Filter != nil {
		switch req.Filter.Type {
		case fb_post_type.Page:
			filterExternalPageIDs = req.Filter.ExternalPageIDs

			externalPageIDs = faboInfo.ExternalPageIDs
			if len(externalPageIDs) != 0 {
				mExternalPageIDs := map[string]bool{}

				for _, externalPageID := range externalPageIDs {
					mExternalPageIDs[externalPageID] = true
				}

				for _, filterExternalPageID := range filterExternalPageIDs {
					if _, ok := mExternalPageIDs[filterExternalPageID]; ok {
						externalPageIDsArgs = append(externalPageIDsArgs, filterExternalPageID)
					}
				}
			}

			if len(externalPageIDsArgs) == 0 {
				externalPageIDsArgs = externalPageIDs
			}

			listFbExternalPostsQuery.ExternalPageIDs = externalPageIDsArgs
		case fb_post_type.User:
			listFbExternalPostsQuery.ExternalUserID = externalUserID
		default:
			listFbExternalPostsQuery.ExternalPageIDs = faboInfo.ExternalPageIDs
			listFbExternalPostsQuery.ExternalUserID = externalUserID
		}

		listFbExternalPostsQuery.LiveVideoStatus = req.Filter.LiveVideoStatus
		listFbExternalPostsQuery.IsLiveVideo = req.Filter.IsLiveVideo
	}

	if err := s.FBMessagingQuery.Dispatch(ctx, listFbExternalPostsQuery); err != nil {
		return nil, err
	}

	return &fabo.ListLiveVideosResponse{
		FbExternalPosts: convertpb.PbFbExternalPosts(listFbExternalPostsQuery.Result.FbExternalPosts),
		Paging:          cmapi.PbCursorPageInfo(paging, &listFbExternalPostsQuery.Result.Paging),
	}, nil
}

func (s *CustomerConversationService) handleErrorWhenSendPrivateReplies(err error) (changeIsPrivateReplies bool, _ error) {
	if err == nil {
		return true, nil
	}

	apiErr, ok := err.(*xerrors.APIError)
	if !ok {
		return false, err
	}

	metaError := apiErr.Meta
	code, ok := metaError["code"]
	if !ok {
		return false, err
	}

	intCode, _err := strconv.Atoi(code)
	if _err != nil {
		return false, err
	}

	switch intCode {
	case int(fbclient.UserCanNotReply):
		return false, cm.Errorf(cm.FacebookError, nil, UserCanNotReply).
			WithMetaM(metaError)
	case int(fbclient.AlreadyRepliedTo):
		return true, cm.Errorf(cm.FacebookError, nil, AlreadyRepliedToUser).
			WithMetaM(metaError)
	case int(fbclient.ReplyingTimeExpired):
		return false, cm.Errorf(cm.FacebookError, nil, ReplyingTimeExpired).
			WithMetaM(metaError)
	default:
		return false, err
	}
}

func (s *CustomerConversationService) buildFbCustomerConversations(
	ctx context.Context, listCustomerConversations []*fbmessaging.FbCustomerConversation,
) ([]*fabo.FbCustomerConversation, error) {
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

		mapExternalUserIDAndImageURl, err := s.buildMapFbUserAvatar(ctx, externalUserIDs)
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

	fbUserExternalIDs := getFbExternalUserIDs(listCustomerConversations)
	listFbUserWithCustomersQuery := &fbusering.ListFbExternalUserWithCustomerByExternalIDsQuery{
		ShopID:      s.SS.Shop().ID,
		ExternalIDs: fbUserExternalIDs,
	}
	err := s.FBUserQuery.Dispatch(ctx, listFbUserWithCustomersQuery)
	if err != nil {
		return nil, err
	}
	var mapExternalIDFbUser = make(map[string]*fbusering.FbExternalUserWithCustomer)
	for _, v := range listFbUserWithCustomersQuery.Result {
		mapExternalIDFbUser[v.FbExternalUser.ExternalID] = v
	}

	listFbUserQuery := &fbusering.ListFbExternalUserByIDsQuery{
		ExtFbUserIDs: fbUserExternalIDs,
	}
	if err = s.FBUserQuery.Dispatch(ctx, listFbUserQuery); err != nil {
		return nil, err
	}
	mapFbExternalUserIDTagIds := buildMapFbExternalUserIDTagIds(listFbUserQuery.Result)

	customerConversations := convertpb.PbFbCustomerConversations(listCustomerConversations)
	for index, conversation := range customerConversations {
		if mapExternalIDFbUser[conversation.ExternalUserID] != nil {
			customerConversations[index].Customer = convertxmin.PbShopCustomer(mapExternalIDFbUser[conversation.ExternalUserID].ShopCustomer)
		}

		tagIDs, ok := mapFbExternalUserIDTagIds[conversation.ExternalUserID]
		if ok {
			customerConversations[index].ExternalUserTags = tagIDs
		}
	}
	return customerConversations, nil
}

func buildMapFbExternalUserIDTagIds(fbExternalUsers []*fbusering.FbExternalUser) map[string][]dot.ID {
	result := map[string][]dot.ID{}

	for _, user := range fbExternalUsers {
		result[user.ExternalID] = user.TagIDs
	}

	return result
}

func getFbExternalUserIDs(customerConversations []*fbmessaging.FbCustomerConversation) []string {
	visited := map[string]struct{}{}
	var result []string
	for _, conv := range customerConversations {
		if _, ok := visited[conv.ExternalUserID]; ok {
			continue
		}
		result = append(result, conv.ExternalUserID)
		visited[conv.ExternalUserID] = struct{}{}
	}
	return result
}

func distinctCustomerConversations(convs []*fbmessaging.FbCustomerConversation) []*fbmessaging.FbCustomerConversation {
	seen := map[dot.ID]struct{}{}
	var result []*fbmessaging.FbCustomerConversation
	for _, conv := range convs {
		if _, ok := seen[conv.ID]; !ok {
			result = append(result, conv)
			seen[conv.ID] = struct{}{}
		}
	}
	return result
}
