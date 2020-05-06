package fabo

import (
	"context"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/top/int/fabo"
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
}

func NewCustomerConversationService(
	ss *session.Session,
	faboInfo *faboinfo.FaboInfo,
	fbMessagingQuery fbmessaging.QueryBus,
	fbMessagingAggr fbmessaging.CommandBus,
) *CustomerConversationService {
	s := &CustomerConversationService{
		ss:               ss,
		faboInfo:         faboInfo,
		fbMessagingQuery: fbMessagingQuery,
		fbMessagingAggr:  fbMessagingAggr,
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
		fbPageIDs := request.Filter.FbPageIDs
		if len(fbPageIDs) == 0 {
			faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID, s.ss.User().ID)
			if err != nil {
				return nil, err
			}
			fbPageIDs = faboInfo.FbPageIDs
		}
		listCustomerConversationsQuery.FbPageIDs = fbPageIDs
		listCustomerConversationsQuery.IsRead = request.Filter.IsRead
		listCustomerConversationsQuery.FbExternalUserID = request.Filter.FbExternalUserID
		listCustomerConversationsQuery.Type = request.Filter.Type
	}
	if len(listCustomerConversationsQuery.FbPageIDs) == 0 {
		faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID, s.ss.User().ID)
		if err != nil {
			return nil, err
		}
		listCustomerConversationsQuery.FbPageIDs = faboInfo.FbPageIDs
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
	faboInfo, err := s.faboInfo.GetFaboInfo(ctx, s.ss.Shop().ID, s.ss.User().ID)
	if err != nil {
		return nil, err
	}
	var fbExternalConversationIDs []string
	if request.Filter != nil {
		fbExternalConversationIDs = request.Filter.FbExternalConversationIDs
	}
	listFbExternalMessagesQuery := &fbmessaging.ListFbExternalMessagesQuery{
		FbPageIDs:         faboInfo.FbPageIDs,
		FbConversationIDs: fbExternalConversationIDs,
		Paging:            *paging,
	}
	if err := s.fbMessagingQuery.Dispatch(ctx, listFbExternalMessagesQuery); err != nil {
		return nil, err
	}

	return &fabo.FbMessagesResponse{
		FbMessages: convertpb.PbFbExternalMessages(listFbExternalMessagesQuery.Result.FbExternalMessages),
		Paging:     cmapi.PbCursorPageInfo(paging, &listFbExternalMessagesQuery.Result.Paging),
	}, nil
}
