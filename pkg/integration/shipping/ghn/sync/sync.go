package sync

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	shipmodel "etop.vn/backend/com/main/shipping/model"
	shippingmodelx "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/scheduler"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
	"etop.vn/backend/pkg/integration/shipping/ghn/update"
	"etop.vn/capi/dot"
	"etop.vn/capi/util"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var ll = l.New()

const (
	defaultNumWorkers = 32
	defaultRecurrent  = 5 * time.Minute
	defaultErrRecurr  = 10 * time.Minute
	maxLogsCount      = 500
)

type Synchronizer struct {
	carrier   *ghn.Carrier
	scheduler *scheduler.Scheduler
}

func New(carrier *ghn.Carrier) *Synchronizer {
	sched := scheduler.New(defaultNumWorkers)
	s := &Synchronizer{
		carrier:   carrier,
		scheduler: sched,
	}
	return s
}

func (s *Synchronizer) InitStatesFromDatabase(ctx context.Context) error {
	shippingSourceNames := s.carrier.GetShippingSourceNames()
	query := &model.GetShippingSources{
		Type:  model.TypeGHN,
		Names: shippingSourceNames,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}

	shippingSources := query.Result
	for _, ss := range shippingSources {
		t := rand.Intn(int(time.Second))
		s.scheduler.AddAfter(ss.ID, time.Duration(t), s.syncCallbackLogs)
	}
	return nil
}

func (s *Synchronizer) Start() {
	s.scheduler.Start()
}

func (s *Synchronizer) Stop() {
	s.scheduler.Stop()
}

func (s *Synchronizer) syncCallbackLogs(id interface{}, p scheduler.Planner) (_err error) {
	t0 := time.Now()
	ctx := bus.Ctx()
	shippingSourceID := id.(dot.ID)

	shippingSourceQuery := &model.GetShippingSource{ID: shippingSourceID}
	if err := bus.Dispatch(ctx, shippingSourceQuery); err != nil {
		return err
	}
	serviceID := shippingSourceQuery.Result.ShippingSource.Name

	shippingSourceInternal := shippingSourceQuery.Result.ShippingSourceInternal
	fromTime := int64(1)
	if shippingSourceInternal != nil && !shippingSourceInternal.LastSyncAt.IsZero() {
		fromTime = shippingSourceInternal.LastSyncAt.UnixNano() / int64(time.Millisecond)
	}
	var ghnOrderLogs []*ghnclient.OrderLog
	skip := 0
	for {
		client, _, err := s.carrier.ParseServiceID(serviceID)
		if err != nil {
			return err
		}
		req := &ghnclient.OrderLogsRequest{
			FromTime: fromTime,
			Skip:     skip,
		}
		resp, err := client.GetOrderLogs(ctx, req)
		if err != nil {
			return err
		}
		if len(resp.Logs) == 0 {
			break
		}
		ghnOrderLogs = append(ghnOrderLogs, resp.Logs...)
		skip = len(ghnOrderLogs)
		if len(ghnOrderLogs) > maxLogsCount {
			break
		}
	}

	defer func() {
		err := recover()
		if err != nil {
			_err = cm.ErrorTracef(cm.Internal, nil, "panic: %v", err)
		}
		if xerrors.IsTrace(_err) {
			bus.PrintAllStack(ctx, true)
		}
		if _err != nil {
			// If there was error, we won't retry so frequently.
			if _, ok := s.scheduler.Peek(id); ok {
				p.AddAfter(id, defaultErrRecurr, s.syncCallbackLogs).Recurrent(defaultErrRecurr)
			}
		}

		switch cm.ErrorCode(_err) {
		case cm.NoError, cm.SkipSync:
		default:
			msg := fmt.Sprintf(`
		ERROR: shipping-sync-service (%vms)
		%v
		`, time.Since(t0)/time.Millisecond, _err)
			bus.Dispatch(ctx, &telebot.SendMessageCommand{Message: msg})
		}
	}()

	p.AddAfter(id, defaultRecurrent, s.syncCallbackLogs).Recurrent(defaultRecurrent)
	if len(ghnOrderLogs) == 0 {
		ll.Info("ghnOrderLogs is empty")
		return nil
	}

	// ghnOrderLogs: order logs already sort by time asc
	// externalCode = ffm ID (in Etop)
	ffmIDs := make([]dot.ID, 0, len(ghnOrderLogs))
	ffmIDMap := make(map[dot.ID]dot.ID)
	ffmsMap := make(map[dot.ID]*shipmodel.Fulfillment)
	for _, log := range ghnOrderLogs {
		if log.OrderInfo.ExternalCode != "" {
			externalCode, err := util.ParseID(log.OrderInfo.ExternalCode.String())
			if err == nil {
				ffmIDMap[externalCode] = externalCode
			}
		}
	}
	for _, id := range ffmIDMap {
		ffmIDs = append(ffmIDs, id)
	}

	ll.Info("Callback Logs ", l.Int("len ghnOrderLogs", len(ghnOrderLogs)))
	ffmsQuery := &shippingmodelx.GetFulfillmentsQuery{
		IDs: ffmIDs,
	}
	if err := bus.Dispatch(ctx, ffmsQuery); err != nil {
		return err
	}
	for _, ffm := range ffmsQuery.Result.Fulfillments {
		ffmsMap[ffm.ID] = ffm
	}

	updateFfmMap := make(map[dot.ID]*shipmodel.Fulfillment)
	for _, oLog := range ghnOrderLogs {
		externalCode, err := util.ParseID(oLog.OrderInfo.ExternalCode.String())
		if err != nil {
			continue
		}
		ffm := ffmsMap[externalCode]
		if ffm == nil {
			continue
		}
		logUpdateAt := oLog.UpdateTime.ToTime()
		ffmLastSyncAt := ffm.LastSyncAt
		if logUpdateAt.After(ffmLastSyncAt) {
			msgFakeCallback := oLog.OrderInfo.ToFakeCallbackOrder()
			ghnOrder := oLog.OrderInfo.ToGHNOrder()
			ffm = update.CalcUpdateFulfillment(ffm, msgFakeCallback, ghnOrder)
			ffm.LastSyncAt = logUpdateAt

			// update ffms Map
			ffmsMap[externalCode] = ffm
			updateFfmMap[ffm.ID] = ffm
		}
	}
	updateFfms := make([]*shipmodel.Fulfillment, 0, len(updateFfmMap))
	for _, ffm := range updateFfmMap {
		updateFfms = append(updateFfms, ffm)
	}
	ll.Info("Callback Logs ", l.Int("len updateFfms", len(updateFfms)))
	cmd := &shippingmodelx.SyncUpdateFulfillmentsCommand{
		ShippingSourceID: shippingSourceID,
		LastSyncAt:       ghnOrderLogs[len(ghnOrderLogs)-1].UpdateTime.ToTime(),
		Fulfillments:     updateFfms,
	}

	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	return nil
}
