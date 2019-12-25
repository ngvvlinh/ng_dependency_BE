package aggregate

import (
	"context"

	"etop.vn/api/meta"
	crmvht "etop.vn/api/supporting/crm/vht"
	"etop.vn/backend/com/supporting/crm/vht/convert"
	"etop.vn/backend/com/supporting/crm/vht/sqlstore"
	syncvht "etop.vn/backend/com/supporting/crm/vht/sync"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
	vhtclient "etop.vn/backend/pkg/integration/vht/client"
)

var _ crmvht.Aggregate = &AggregateService{}

type AggregateService struct {
	VhtStore           sqlstore.VhtCallHistoriesFactory
	VhtClient          *vhtclient.Client
	SyncVhtCallHistory *syncvht.SyncVht
}

func New(db *cmsql.Database, vhtClient *vhtclient.Client) *AggregateService {
	return &AggregateService{
		VhtStore:  sqlstore.NewVhtCallHistoryStore(db),
		VhtClient: vhtClient,
	}
}

func (a *AggregateService) MessageBus() crmvht.CommandBus {
	b := bus.New()
	return crmvht.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a AggregateService) CreateOrUpdateCallHistoryBySDKCallID(ctx context.Context, req *crmvht.VhtCallLog) (*crmvht.VhtCallLog, error) {

	sdkCallID := req.SdkCallID
	if sdkCallID == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing sdkCallID in request", nil)
	}
	callHistoryModel := convert.ConvertToModel(req)
	callHistoryModel.SyncStatus = "Pending"

	_, err := a.VhtStore(ctx).BySdkCallID(sdkCallID).GetCallHistory()
	if err != nil {
		err = a.VhtStore(ctx).CreateVhtCallHistory(callHistoryModel)
	} else {
		err = a.VhtStore(ctx).BySdkCallID(sdkCallID).UpdateVhtCallHistory(callHistoryModel)
	}
	if err != nil {
		return nil, err
	}

	dbResult, err := a.VhtStore(ctx).BySdkCallID(sdkCallID).GetCallHistory()
	if err != nil {
		return nil, err
	}
	return convert.ConvertFromModel(dbResult), nil
}

func (a AggregateService) CreateOrUpdateCallHistoryByCallID(ctx context.Context, req *crmvht.VhtCallLog) (*crmvht.VhtCallLog, error) {
	callID := req.CallID
	if callID == "" {
		return nil, cm.Error(cm.InvalidArgument, "Missing sdkCallID in request", nil)
	}
	callHistoryModel := convert.ConvertToModel(req)

	_, err := a.VhtStore(ctx).ByCallID(callID).GetCallHistory()
	if err != nil {
		err = a.VhtStore(ctx).CreateVhtCallHistory(callHistoryModel)
	} else {
		err = a.VhtStore(ctx).ByCallID(callID).UpdateVhtCallHistory(callHistoryModel)
	}
	if err != nil {
		return nil, err
	}

	dbResult, err := a.VhtStore(ctx).ByCallID(callID).GetCallHistory()
	if err != nil {
		return nil, err
	}
	return convert.ConvertFromModel(dbResult), nil
}

func (a *AggregateService) PingServerVht(context.Context, *meta.Empty) error {
	if a.SyncVhtCallHistory == nil {
		a.SyncVhtCallHistory = syncvht.New(a.VhtStore, a.VhtClient)
	}
	return a.SyncVhtCallHistory.PingServerVht()
}

func (a *AggregateService) SyncVhtCallHistories(ctx context.Context, req *crmvht.SyncVhtCallHistoriesArgs) error {
	err := a.SyncVhtCallHistory.SyncVhtCallHistory(ctx, req.SyncTime)
	if err != nil {
		return err
	}
	err = a.SyncVhtCallHistory.SyncVhtCallHistoryPending(ctx)
	if err != nil {
		return err
	}
	return nil
}
