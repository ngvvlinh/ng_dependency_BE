package fabo

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbpaging"
	"o.o/api/top/int/fabo"
	"o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/fabo/convertpb"
	"o.o/backend/pkg/fabo/faboinfo"
)

type CustomerConversationService struct {
	session.Sessioner
	ss *session.Session

	faboInfo         *faboinfo.FaboInfo
	fbMessagingQuery fbmessaging.QueryBus
	fbMessagingAggr  fbmessaging.CommandBus
	fbPagingQuery    fbpaging.QueryBus
}

func NewCustomerConversationService(
	ss *session.Session,
	faboInfo *faboinfo.FaboInfo,
	fbMessagingQuery fbmessaging.QueryBus,
	fbMessagingAggr fbmessaging.CommandBus,
	fbPagingQuery fbpaging.QueryBus,
) *CustomerConversationService {
	s := &CustomerConversationService{
		ss:               ss,
		faboInfo:         faboInfo,
		fbMessagingQuery: fbMessagingQuery,
		fbMessagingAggr:  fbMessagingAggr,
		fbPagingQuery:    fbPagingQuery,
	}
	return s
}

func (s *CustomerConversationService) Clone() fabo.CustomerConversationService {
	res := *s
	res.Sessioner, res.ss = s.ss.Split()
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
			faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID)
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
		faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID)
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
	faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID)
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
	if request.Filter.ExternalPostID == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "")
	}
	paging, err := cmapi.CMCursorPaging(request.Paging)
	if err != nil {
		return nil, err
	}
	faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID)
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

	return &fabo.ListCommentsByExternalPostIDResponse{
		FbComments: &fabo.FbCommentsResponse{
			FbComments: convertpb.PbFbExternalComments(listFbExternalCommentsQuery.Result.FbExternalComments),
			Paging:     cmapi.PbCursorPageInfo(paging, &listFbExternalCommentsQuery.Result.Paging),
		},
		FbPost: convertpb.PbFbExternalPost(fbExternalPost),
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
