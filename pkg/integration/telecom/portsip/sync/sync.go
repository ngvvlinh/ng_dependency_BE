package sync

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"o.o/api/etelecom"
	"o.o/api/etelecom/call_direction"
	"o.o/api/etelecom/call_state"
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
	etelecomxservicedriver "o.o/backend/pkg/integration/telecom/etelecomxservice/driver"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var _ providertypes.TelecomSync = &PortsipSync{}

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

type PortsipSync struct {
	dbMain         *cmsql.Database
	telecomManager *telecomprovider.TelecomManager
	telecomQuery   etelecom.QueryBus
	telecomAggr    etelecom.CommandBus
	connectionAggr connectioning.CommandBus
	mutex          sync.Mutex

	// crawl call log
	schedulerCrawlCallLogs *scheduler.Scheduler
	crawlCallLogsTicker    *time.Ticker
	// key = {connection_id}-{owner_id}
	//		owner_id = 0 then key = {connection_id}
	// value = taskID
	mapTenantInProgress map[string]dot.ID
	mapTaskCrawlCallLog map[dot.ID]*TaskCrawlLogArguments
}

func New(
	db com.MainDB, telecomManager *telecomprovider.TelecomManager,
	telecomQS etelecom.QueryBus, telecomA etelecom.CommandBus,
	connectionA connectioning.CommandBus,
) *PortsipSync {
	return &PortsipSync{
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

func (s *PortsipSync) Init(ctx context.Context) error {
	return nil
}

func (s *PortsipSync) Start(ctx context.Context) error {
	s.schedulerCrawlCallLogs.Start()

	s.addCrawlCallLogsTasks(ctx)
	for {
		select {
		case <-s.crawlCallLogsTicker.C:
			s.addCrawlCallLogsTasks(ctx)
		}
	}
}

func (s *PortsipSync) Stop(ctx context.Context) error {
	s.schedulerCrawlCallLogs.Stop()
	return nil
}

// crawl for each shop_connection
// in this context we define shopConnection as tenant
func (s *PortsipSync) addCrawlCallLogsTasks(ctx context.Context) error {
	shopConnections, err := s.listShopConnections(ctx)
	if err != nil {
		return err
	}

	var taskIDs []dot.ID

	s.mutex.Lock()
	for _, shopConnection := range shopConnections {
		keyTenant := s.getKeyTenantInProgress(shopConnection.ConnectionID, shopConnection.OwnerID)

		// check tenant is in progress
		if _, ok := s.mapTenantInProgress[keyTenant]; ok {
			continue
		}

		// add task
		taskID := cm.NewID()
		// add tenant in to progress
		s.mapTenantInProgress[keyTenant] = taskID

		// add task with arguments
		s.mapTaskCrawlCallLog[taskID] = &TaskCrawlLogArguments{
			shopConnection: shopConnection,
		}
		taskIDs = append(taskIDs, taskID)
	}
	s.mutex.Unlock()

	// add tasks
	for _, taskID := range taskIDs {
		t := rand.Intn(int(time.Second))
		s.schedulerCrawlCallLogs.AddAfter(taskID, time.Duration(t), s.crawlCallLogs)
	}

	return nil
}

func (s *PortsipSync) crawlCallLogs(id interface{}, p scheduler.Planner) (err error) {
	ctx := bus.Ctx()

	taskArgumentID := id.(dot.ID)

	s.mutex.Lock()
	taskArguments := s.mapTaskCrawlCallLog[taskArgumentID]
	s.mutex.Unlock()

	ownerID := taskArguments.shopConnection.OwnerID
	connectionID := taskArguments.shopConnection.ConnectionID
	lastSyncAt := taskArguments.shopConnection.LastSyncAt
	// make sure we don't miss call logs (cdr) in 15 minutes
	lastSyncAt = lastSyncAt.Add(-15 * time.Minute)

	tenantKey := s.getKeyTenantInProgress(connectionID, ownerID)

	defer func() {
		// remove tenant in mapTenantInProgress
		s.mutex.Lock()
		delete(s.mapTenantInProgress, tenantKey)
		s.mutex.Unlock()

		if err != nil {
			sendError(ownerID, connectionID, err)
		}
	}()

	shopConn := taskArguments.shopConnection
	if shopConn.TelecomData == nil || shopConn.TelecomData.TenantToken == "" {
		return nil
	}
	etelecomXServiceDriver := etelecomxservicedriver.New(shopConn.TelecomData.TenantToken)

	var (
		scrollID        string
		getCallLogsResp *etelecomxservicedriver.GetCallLogsResponse
	)
	now := time.Now()
	tenant, err := s.getTenant(ctx, connectionID, ownerID)
	if err != nil {
		return err
	}

	var lastCallLogAt time.Time
	for true {
		getCallLogsReq := &etelecomxservicedriver.GetCallLogsRequest{
			ScrollID:  scrollID,
			StartedAt: lastSyncAt,
			EndedAt:   now,
		}
		getCallLogsResp, err = etelecomXServiceDriver.GetCallLogs(ctx, getCallLogsReq)
		if err != nil {
			return err
		}

		scrollID = getCallLogsResp.ScrollID
		if len(getCallLogsResp.CallLogs) == 0 {
			break
		}
		if lastCallLogAt.IsZero() {
			lastCallLogAt = getCallLogsResp.CallLogs[0].EndedAt
		}

		for _, callLogResp := range getCallLogsResp.CallLogs {
			if lastSyncAt.After(callLogResp.StartedAt) {
				break
			}
			_callsInfo := s.getCallInfo(ctx, tenant.ID, callLogResp)
			for _, info := range _callsInfo {
				cmdCreate := &etelecom.CreateOrUpdateCallLogFromCDRCommand{
					ExternalID:         callLogResp.CallID,
					StartedAt:          callLogResp.StartedAt,
					EndedAt:            callLogResp.EndedAt,
					Duration:           callLogResp.Duration,
					AudioURLs:          callLogResp.AudioURLs,
					ExternalDirection:  callLogResp.Direction,
					ExternalCallStatus: callLogResp.CallStatus,
					OwnerID:            ownerID,
					ConnectionID:       connectionID,
					ExternalSessionID:  info.SessionID,
					Callee:             info.Callee,
					Caller:             info.Caller,
					Direction:          info.Direction,
					CallState:          info.CallState,
					ExtensionID:        info.ExtensionID,
				}

				if _err := s.telecomAggr.Dispatch(ctx, cmdCreate); _err != nil {
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
		LastSyncAt:   lastCallLogAt,
	}
	if err = s.connectionAggr.Dispatch(ctx, updateShopConnectionCmd); err != nil {
		return err
	}

	return nil
}

func (s *PortsipSync) getKeyTenantInProgress(connectionID, ownerID dot.ID) string {
	if ownerID == 0 {
		return fmt.Sprintf("%d", connectionID.Int64())
	}
	return fmt.Sprintf("%d-%d", connectionID.Int64(), ownerID.Int64())
}

func (s *PortsipSync) listShopConnections(ctx context.Context) (shopConnections []*model.ShopConnection, err error) {
	// get connections of VHT
	var connectionIDs []dot.ID
	{
		var _connections model.Connections

		err = s.dbMain.
			Where("connection_type = ?", connection_type.Telecom).
			Where("connection_provider = ?", connection_type.ConnectionProviderPortsip).
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

		err = s.dbMain.
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

type callInfo struct {
	Callee      string
	Caller      string
	CallState   call_state.CallState
	ExtensionID dot.ID
	HotlineID   dot.ID
	Direction   call_direction.CallDirection
	SessionID   string
}

func (s *PortsipSync) getCallInfo(ctx context.Context, tenantID dot.ID, callLog *etelecomxservicedriver.CallLog) (res []*callInfo) {
	_callInfo := &callInfo{
		Callee:    callLog.Callee,
		Caller:    callLog.Caller,
		CallState: callLog.CallState,
		SessionID: callLog.SessionID,
	}

	res = append(res, _callInfo)

	switch callLog.Direction {
	case call_direction.In.String():
		// dựa vào call_targets để xác định ext nhận cuộc gọi
		// call target sẽ trả về 1 mảng các giá trị của extension
		// trường hợp porsip xài Queue hay Ring Group nó sẽ có tất cả ext của Queue hay Ring Group đó
		// cần tìm ra 1 extension tương ứng để tạo call log
		// Nếu ko tìm thấy ext => vẫn trả về kết quả
		_callInfo.Direction = call_direction.In
		extensionNumbersAnswered := []string{}
		extensionNumbersNotAnswered := []string{}
		targetsMap := make(map[string]*etelecomxservicedriver.CallTarget)
		for _, target := range callLog.CallTargets {
			targetsMap[target.TargetNumber] = target
			if target.CallState == call_state.Answered {
				extensionNumbersAnswered = append(extensionNumbersAnswered, target.TargetNumber)
			} else {
				extensionNumbersNotAnswered = append(extensionNumbersNotAnswered, target.TargetNumber)
			}
		}

		// find extension: priority extension number answered first
		ext, err := s.findExtension(ctx, extensionNumbersAnswered, tenantID)
		if err != nil {
			ext, err = s.findExtension(ctx, extensionNumbersNotAnswered, tenantID)
			if err != nil || ext == nil {
				return
			}
		}

		target, ok := targetsMap[ext.ExtensionNumber]
		if !ok {
			return
		}
		_callInfo.Callee = target.TargetNumber
		_callInfo.CallState = target.CallState
		_callInfo.ExtensionID = ext.ID
		return res

	case call_direction.Out.String():
		_callInfo.Direction = call_direction.Out
		extensionNumbers := []string{callLog.Caller}
		ext, err := s.findExtension(ctx, extensionNumbers, tenantID)
		if err != nil {
			return
		}
		_callInfo.ExtensionID = ext.ID
		return

	case call_direction.Ext.String():
		// tách làm 2 call log: gọi vào và gọi ra
		// gọi vào (in): dựa vào callee
		extensionCalleeNumbers := []string{callLog.Callee}
		extCallee, err := s.findExtension(ctx, extensionCalleeNumbers, tenantID)
		if err != nil {
			return
		}
		var result = []*callInfo{}
		callInfoIn := *_callInfo
		callInfoIn.Direction = call_direction.ExtIn
		callInfoIn.ExtensionID = extCallee.ID

		// gọi ra (out): dựa vào caller
		extensionCallerNumbers := []string{callLog.Caller}
		extCaller, err := s.findExtension(ctx, extensionCallerNumbers, tenantID)
		if err != nil {
			return
		}
		callInfoOut := *_callInfo
		callInfoOut.Direction = call_direction.ExtOut
		callInfoOut.ExtensionID = extCaller.ID
		callInfoOut.SessionID += "-" + call_direction.Out.String()

		result = append(result, &callInfoIn, &callInfoOut)
		return result

	default:
		return
	}
}

func (s *PortsipSync) findExtension(ctx context.Context, extNumbers []string, tenantID dot.ID) (ext *etelecom.Extension, err error) {
	query := &etelecom.ListExtensionsQuery{
		TenantID:         tenantID,
		ExtensionNumbers: extNumbers,
	}
	if err = s.telecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	extensions := query.Result
	if len(extensions) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "Extension not found: %v", extNumbers)
	}
	return extensions[0], nil
}

func (s *PortsipSync) getTenant(ctx context.Context, connID, ownerID dot.ID) (*etelecom.Tenant, error) {
	query := &etelecom.GetTenantByConnectionQuery{
		OwnerID:      ownerID,
		ConnectionID: connID,
	}
	if err := s.telecomQuery.Dispatch(ctx, query); err != nil {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Tenant not found").WithMetap("query", query)
	}
	return query.Result, nil
}
