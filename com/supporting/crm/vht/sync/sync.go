package sync

import (
	"context"
	"time"

	"etop.vn/backend/com/supporting/crm/vht/model"
	"etop.vn/backend/com/supporting/crm/vht/sqlstore"
	cm "etop.vn/backend/pkg/common"
	vhtclient "etop.vn/backend/pkg/integration/vht/client"
	"etop.vn/common/jsonx"
)

type SyncVht struct {
	VhtCallHistoryStore sqlstore.VhtCallHistoriesFactory
	VhtClient           *vhtclient.Client
}

func New(vhtstore sqlstore.VhtCallHistoriesFactory, client *vhtclient.Client) *SyncVht {
	return &SyncVht{
		VhtCallHistoryStore: vhtstore,
		VhtClient:           client,
	}
}

func (s *SyncVht) SyncVhtCallHistory(ctx context.Context, lasTimeSync time.Time) error {
	fromDate := lasTimeSync.Unix()
	toDate := time.Now().Unix()
	lasTimeSync = time.Now()

	queryDTO := &vhtclient.VHTHistoryQueryDTO{
		Page:        1,
		Limit:       50,
		DateStarted: fromDate,
		DateEnded:   toDate,
		SortBy:      "time_started",
		SortType:    "ASC",
	}
	for true {
		result, err := s.VhtClient.GetHistories(queryDTO)
		if err != nil {
			return err
		}
		for i := 0; i < len(result.Items); i++ {
			data := vhtclient.ConvertToModel(result.Items[i])
			err = data.BeforeInsertOrUpdate()
			if err != nil {
				return err
			}
			query := s.VhtCallHistoryStore(ctx).ByCallID(data.CallID)
			var oldData *model.VhtCallHistory
			oldData, err = query.GetCallHistory()
			data.OData = ""
			if err != nil && cm.ErrorCode(err) == cm.NotFound {
				err = s.VhtCallHistoryStore(ctx).CreateVhtCallHistory(data)
			} else if err == nil {
				if oldData.SyncStatus == "Done" {
					continue
				}

				oldDataMarshal := jsonx.MustMarshalToString(oldData)
				data.OData = oldDataMarshal
				err = s.VhtCallHistoryStore(ctx).ByCallID(data.CallID).UpdateVhtCallHistory(data)
			} else {
				return err
			}
		}
		if len(result.Items) < 50 {
			break
		}
		queryDTO.Page = queryDTO.Page + 1
	}
	return nil
}

func (s *SyncVht) SyncVhtCallHistoryPending(ctx context.Context) error {
	historiesPending, err := s.VhtCallHistoryStore(ctx).ByStatus("Pending").GetCallHistories()
	if err != nil {
		return nil
	}
	for i := 0; i < len(historiesPending); i++ {
		var result *vhtclient.VhtCallHistory
		result, err = s.VhtClient.GetHistoryBySDKCallID(historiesPending[i].SdkCallID)
		if err != nil {
			return err
		}
		data := vhtclient.ConvertToModel(result)
		err = s.VhtCallHistoryStore(ctx).BySdkCallID(data.SdkCallID).UpdateVhtCallHistory(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SyncVht) PingServerVht() error {
	return s.VhtClient.PingServerVht()
}
