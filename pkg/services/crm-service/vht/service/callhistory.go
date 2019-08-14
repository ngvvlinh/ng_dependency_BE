package service

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/backend/pb/services/crmservice"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/services/crm-service/model"
)

func (s *VhtService) GetCallHistories(ctx context.Context, req *crmservice.GetCallHistoriesRequest) (*crmservice.GetCallHistoriesResponse, error) {
	page := req.Page
	perpage := req.Perpage
	if perpage == 0 {
		perpage = 100
	}
	if perpage > 1000 {
		perpage = 1000
	}
	var paging meta.Paging
	paging.Offset = page
	paging.Limit = perpage

	//search in db
	textSearch := req.TextSearch
	var dbResult []*model.VhtCallHistory
	var err error
	if textSearch != "" {
		dbResult, err = s.vhtCallHistories(ctx).Paging(paging).SearchVhtCallHistories(textSearch)
	} else {
		dbResult, err = s.vhtCallHistories(ctx).Paging(paging).GetCallHistories()
	}
	if err != nil {
		return nil, err
	}

	var vhtCallHistoryResult []*crmservice.VHTCallLog
	for i := 0; i < len(dbResult); i++ {
		dataRow := dbResult[i]
		callHistoryRow := s.ConvertModel2Proto(dataRow)
		vhtCallHistoryResult = append(vhtCallHistoryResult, callHistoryRow)
	}

	return &crmservice.GetCallHistoriesResponse{
		VhtCallLog: vhtCallHistoryResult,
	}, nil
}

func (s *VhtService) CreateOrUpdateBySDKCallID(ctx context.Context, req *crmservice.VHTCallLog) (*crmservice.VHTCallLog, error) {

	sdkCallID := req.SdkCallId
	if sdkCallID == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing sdkCallID in request", nil)
	}
	callHistoryModel := s.ConvertProto2Model(req)
	callHistoryModel.SyncStatus = "Pending"

	_, err := s.vhtCallHistories(ctx).BySdkCallID(sdkCallID).GetCallHistory()
	if err != nil {
		err = s.vhtCallHistories(ctx).CreateVhtCallHistory(callHistoryModel)
	} else {
		err = s.vhtCallHistories(ctx).BySdkCallID(sdkCallID).UpdateVhtCallHistory(callHistoryModel)
	}
	if err != nil {
		return nil, err
	}

	dbResult, err := s.vhtCallHistories(ctx).BySdkCallID(sdkCallID).GetCallHistory()
	if err != nil {
		return nil, err
	}
	return s.ConvertModel2Proto(dbResult), nil
}

func (s *VhtService) CreateOrUpdateCallHistoryByCallID(ctx context.Context, req *crmservice.VHTCallLog) (*crmservice.VHTCallLog, error) {

	callID := req.CallId
	if callID == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing sdkCallID in request", nil)
	}
	callHistoryModel := s.ConvertProto2Model(req)

	_, err := s.vhtCallHistories(ctx).ByCallID(callID).GetCallHistory()
	if err != nil {
		err = s.vhtCallHistories(ctx).CreateVhtCallHistory(callHistoryModel)
	} else {
		err = s.vhtCallHistories(ctx).ByCallID(callID).UpdateVhtCallHistory(callHistoryModel)
	}
	if err != nil {
		return nil, err
	}

	dbResult, err := s.vhtCallHistories(ctx).ByCallID(callID).GetCallHistory()
	if err != nil {
		return nil, err
	}
	return s.ConvertModel2Proto(dbResult), nil
}
