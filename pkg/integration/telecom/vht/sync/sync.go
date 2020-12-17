package sync

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/pp"

	"o.o/api/etelecom"
	"o.o/api/etelecom/call_log_direction"
	"o.o/api/main/connectioning"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	telecomprovider "o.o/backend/com/etelecom/provider"
	providertypes "o.o/backend/com/etelecom/provider/types"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/connectioning/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/scheduler"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	vhtclient "o.o/backend/pkg/integration/telecom/vht/client"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var _ providertypes.TelecomSync = &VHTSync{}

var ll = l.New().WithChannel(meta.ChannelTelecomProvider)

const (
	defaultNumWorkers                = 16
	defaultTimeCrawlCallLogs         = 5 * time.Minute
	uniqueExternalIDDirectionCallLog = "call_log_extension_id_external_id_direction_key"
	uniqueExternalIDCallLog          = "call_log_extension_id_external_id_key"
)

type TaskCrawlLogArguments struct {
	shopConnection *model.ShopConnection
}

type VHTSync struct {
	dbMain         *cmsql.Database
	telecomManager *telecomprovider.TelecomManager
	telecomQuery   etelecom.QueryBus
	telecomAggr    etelecom.CommandBus
	connectionAggr connectioning.CommandBus
	mutex          sync.Mutex

	// crawl call log
	schedulerCrawlCallLogs *scheduler.Scheduler
	crawlCallLogsTicker    *time.Ticker
	// key = {connection_id}-{shop_id}
	//		shop_id = 0 then key = {connection_id}
	// value = taskID
	mapTenantInProgress map[string]dot.ID
	mapTaskCrawlCallLog map[dot.ID]*TaskCrawlLogArguments
}

func New(
	db com.MainDB, telecomManager *telecomprovider.TelecomManager,
	telecomQS etelecom.QueryBus, telecomA etelecom.CommandBus,
	connectionA connectioning.CommandBus,
) *VHTSync {
	return &VHTSync{
		dbMain:                 db,
		telecomManager:         telecomManager,
		telecomQuery:           telecomQS,
		telecomAggr:            telecomA,
		connectionAggr:         connectionA,
		schedulerCrawlCallLogs: scheduler.New(defaultNumWorkers),
		crawlCallLogsTicker:    time.NewTicker(defaultTimeCrawlCallLogs),
		mapTenantInProgress:    make(map[string]dot.ID),
		mapTaskCrawlCallLog:    make(map[dot.ID]*TaskCrawlLogArguments),
	}
}

func (v *VHTSync) Init(ctx context.Context) error {
	return nil
}

func (v *VHTSync) Start(ctx context.Context) error {
	v.schedulerCrawlCallLogs.Start()

	v.addCrawlCallLogsTasks(ctx)
	for {
		select {
		case <-v.crawlCallLogsTicker.C:
			v.addCrawlCallLogsTasks(ctx)
		}
	}
}

func (v *VHTSync) Stop(ctx context.Context) error {
	v.schedulerCrawlCallLogs.Stop()
	return nil
}

// crawl for each shop_connection
// in this context we define shopConnection as tenant
func (v *VHTSync) addCrawlCallLogsTasks(ctx context.Context) error {
	shopConnections, err := v.listShopConnections(ctx)
	if err != nil {
		return err
	}

	var taskIDs []dot.ID

	v.mutex.Lock()
	for _, shopConnection := range shopConnections {
		pp.Println("shop connection :: ", shopConnection)
		keyTenant := v.getKeyTenantInProgress(shopConnection.ConnectionID, shopConnection.ShopID)

		// check tenant is in progress
		if _, ok := v.mapTenantInProgress[keyTenant]; ok {
			continue
		}

		// add task
		taskID := cm.NewID()
		// add tenant in to progress
		v.mapTenantInProgress[keyTenant] = taskID

		// add task with arguments
		v.mapTaskCrawlCallLog[taskID] = &TaskCrawlLogArguments{
			shopConnection: shopConnection,
		}
		taskIDs = append(taskIDs, taskID)
	}
	v.mutex.Unlock()

	// add tasks
	for _, taskID := range taskIDs {
		t := rand.Intn(int(time.Second))
		v.schedulerCrawlCallLogs.AddAfter(taskID, time.Duration(t), v.crawlCallLogs)
	}

	return nil
}

func (v *VHTSync) crawlCallLogs(id interface{}, p scheduler.Planner) (err error) {
	ctx := bus.Ctx()

	taskArgumentID := id.(dot.ID)

	v.mutex.Lock()
	taskArguments := v.mapTaskCrawlCallLog[taskArgumentID]
	v.mutex.Unlock()

	ownerID := taskArguments.shopConnection.OwnerID
	connectionID := taskArguments.shopConnection.ConnectionID
	lastSyncAt := taskArguments.shopConnection.LastSyncAt

	tenantKey := v.getKeyTenantInProgress(connectionID, ownerID)

	defer func() {
		// remove tenant in mapTenantInProgress
		v.mutex.Lock()
		delete(v.mapTenantInProgress, tenantKey)
		v.mutex.Unlock()

		if err != nil {
			sendError(ownerID, connectionID, err)
		}
	}()

	telecomDriver, err := v.telecomManager.GetTelecomDriver(ctx, connectionID, ownerID)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "Get driver error: %v", err.Error())
	}

	var (
		scrollID        string
		getCallLogsResp *providertypes.GetCallLogsResponse
	)
	now := time.Now()

	for true {
		getCallLogsReq := &providertypes.GetCallLogsRequest{
			StartedAt: lastSyncAt,
			EndedAt:   now,
			ScrollID:  scrollID,
		}
		getCallLogsResp, err = telecomDriver.GetCallLogs(ctx, getCallLogsReq)
		if err != nil {
			return err
		}

		scrollID = getCallLogsResp.ScrollID
		if len(getCallLogsResp.CallLogs) == 0 {
			break
		}

		for _, callLogResp := range getCallLogsResp.CallLogs {
			callState := vhtclient.VHTCallStatus(callLogResp.CallStatus).ToCallState()
			cmdCreate := &etelecom.CreateCallLogFromCDRCommand{
				ExternalID:         callLogResp.CallID,
				StartedAt:          callLogResp.StartedAt,
				EndedAt:            callLogResp.EndedAt,
				Duration:           callLogResp.Duration,
				Caller:             callLogResp.Caller,
				Callee:             callLogResp.Callee,
				AudioURLs:          callLogResp.AudioURLs,
				ExternalDirection:  callLogResp.Direction,
				ExternalCallStatus: callLogResp.CallStatus,
				CallState:          callState,
				CallStatus:         callState.ToStatus5(),
				OwnerID:            ownerID,
				ConnectionID:       connectionID,
			}

			if lastSyncAt.After(callLogResp.StartedAt) {
				return nil
			}
			var createCallLogCmds []*etelecom.CreateCallLogFromCDRCommand

			switch callLogResp.Direction {
			case call_log_direction.In.String(),
				call_log_direction.Out.String():
				direction, _ := call_log_direction.ParseCallLogDirection(callLogResp.Direction)
				cmd := *cmdCreate
				cmd.Direction = direction

				createCallLogCmds = append(createCallLogCmds, &cmd)

			case call_log_direction.Ext.String():
				// Api chỉ trả về 1 call log
				// để ghi nhận cuộc gọi từ extension -> extension thì cần tạo 2 call logs với direction (in, out)
				cmdIn := *cmdCreate
				cmdOut := *cmdCreate
				cmdIn.Direction = call_log_direction.In
				cmdOut.Direction = call_log_direction.Out
				createCallLogCmds = append(createCallLogCmds,
					&cmdIn, &cmdOut,
				)
			}

			for _, createCallLogCmd := range createCallLogCmds {
				if _err := v.telecomAggr.Dispatch(ctx, createCallLogCmd); _err != nil {
					if cm.ErrorCode(_err) == cm.NotFound {
						// ignore if extension or any things not found
						continue
					}
					if xerr, ok := _err.(*xerrors.APIError); ok && xerr.Err != nil {
						errMsg := xerr.Err.Error()
						if strings.Contains(errMsg, uniqueExternalIDCallLog) ||
							strings.Contains(errMsg, uniqueExternalIDDirectionCallLog) {
							// ignore if duplicate
							continue
						}
					}
					return _err
				}
			}
		}
	}

	updateShopConnectionCmd := &connectioning.UpdateShopConnectionLastSyncAtCommand{
		OwnerID:      ownerID,
		ConnectionID: connectionID,
		LastSyncAt:   now,
	}
	if err := v.connectionAggr.Dispatch(ctx, updateShopConnectionCmd); err != nil {
		return err
	}

	return nil
}

func (v *VHTSync) getKeyTenantInProgress(connectionID, shopID dot.ID) string {
	if shopID == 0 {
		return fmt.Sprintf("%d", connectionID.Int64())
	}
	return fmt.Sprintf("%d-%d", connectionID.Int64(), shopID.Int64())
}

func (v *VHTSync) listShopConnections(ctx context.Context) (shopConnections []*model.ShopConnection, err error) {
	// get connections of VHT
	var connectionIDs []dot.ID
	{
		var _connections model.Connections

		err = v.dbMain.
			Where("connection_type = ?", connection_type.Telecom).
			Where("connection_provider = ?", connection_type.ConnectionProviderVHT).
			Where("status = ?", status3.P).
			Limit(1000).
			Find(&_connections)
		if err != nil {
			return nil, cm.Errorf(cm.Internal, err, "Error when get telecom connections VHT: %v", err)
		}

		if len(_connections) == 0 {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Telecom connections VHT not found")
		}

		for _, connection := range _connections {
			connectionIDs = append(connectionIDs, connection.ID)
		}
	}

	createdAt := time.Time{}
	for {
		var _shopConnections model.ShopConnections

		err = v.dbMain.
			Where("created_at > ?", createdAt).
			Where("status = ?", status3.P).
			In("connection_id", connectionIDs).
			OrderBy("created_at asc").
			Limit(1000).
			Find(&_shopConnections)
		if err != nil {
			return nil, err
		}

		if len(_shopConnections) == 0 {
			break
		}

		shopConnections = append(shopConnections, _shopConnections...)
		createdAt = _shopConnections[len(_shopConnections)-1].CreatedAt
	}

	return
}

func sendError(ownerID, connectionID dot.ID, err error) {
	if ownerID != 0 {
		ll.SendMessagef("–––\n👹 Telecom-sync-service: VHT 👹\n- Method: Direct\n- ConnectionID: %v\n- OwnerID: %v \n- Error: %v\n---", connectionID, ownerID, err.Error())
		return
	}
	ll.SendMessagef("–––\n👹 Telecom-sync-service: VHT 👹\n- Method: Builtin\n- ConnectionID: %v\n- Error: %v\n---", connectionID, err.Error())
}
