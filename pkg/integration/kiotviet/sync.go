package kiotviet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/scheduler"
	"etop.vn/backend/pkg/common/telebot"
	cmWrapper "etop.vn/backend/pkg/common/wrapper"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/kiotviet/ssm"
)

func init() {
	bus.AddHandler("kiotviet", SyncProductSource)
}

var (
	defaultZeroDate = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	gScheduler *scheduler.Scheduler
	gStates    MapStates
)

const (
	defaultNumWorkers = 32
	defaultRecurrent  = 5 * time.Minute
	defaultErrRecurr  = time.Hour
)

type SyncState = ssm.SyncState

type SavedSyncState struct {
	LastSyncAt time.Time  `json:"last_sync_at"`
	State      *SyncState `json:"state"`
	NextState  *SyncState `json:"next_state"`
	Done       bool       `json:"done"`
	Error      error      `json:"error"`
}

func (s *SavedSyncState) ToJSON() []byte {
	data, _ := json.Marshal(s)
	return data
}

type SupplierSyncState struct {
	ID      int64
	Webhook bool
}

type MapStates struct {
	m  map[int64]SupplierSyncState
	mu sync.RWMutex
}

func (m *MapStates) Get(id int64) (SupplierSyncState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.m[id]
	return s, ok
}

func (m *MapStates) Set(id int64, state SupplierSyncState) {
	m.mu.Lock()
	m.m[id] = state
	m.mu.Unlock()
}

func StartSync() {
	ctx := bus.NewRootContext(context.Background())
	gStates.m = make(map[int64]SupplierSyncState)

	query := &model.GetAllProductSourcesQuery{External: cm.PBool(true)}
	if err := bus.Dispatch(ctx, query); err != nil {
		cmWrapper.LogErrorAndTrace(ctx, err, "Can not sync ProductSource")
		return
	}

	sources := query.Result.Sources
	ll.Info("Sync all external ProductSources", l.Int("num", len(sources)))

	gScheduler = scheduler.New(defaultNumWorkers)
	for _, s := range sources {
		t := rand.Intn(int(time.Second))
		gScheduler.AddAfter(s.ID, time.Duration(t), syncProductSource).
			Recurrent(defaultRecurrent)
	}

	gScheduler.Start()
}

func StopSync() {
	gScheduler.Stop()
}

func SyncProductSource(ctx context.Context, cmd *SyncProductSourceCommand) error {
	if cmd.SourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SourceID", nil)
	}

	gScheduler.AddAfter(cmd.SourceID, time.Second, syncProductSource).
		Recurrent(defaultRecurrent)
	return nil
}

func syncProductSource(id interface{}, p scheduler.Planner) (_err error) {
	t0 := time.Now()
	sourceID := id.(int64)
	ctx := bus.NewRootContext(context.Background())
	defer func() {
		e := recover()
		if e != nil {
			_err = cm.Error(cm.RuntimePanic, cm.F("%v", e), nil)
			debug.PrintStack()
			bus.PrintAllStack(ctx, false)
			return
		}
		if _err != nil {
			fmt.Printf("%+v", _err)
			if cm.IsTrace(_err) || cm.ErrorCode(_err) == cm.Internal {
				bus.PrintAllStack(ctx, false)
			}
		}
	}()

	psQuery := &model.GetProductSourceExtendedQuery{
		GetProductSourceProps: model.GetProductSourceProps{
			ID: sourceID,
		},
	}
	if err := bus.Dispatch(ctx, psQuery); err != nil {
		return err
	}

	ps := psQuery.Result.ProductSource
	psi := psQuery.Result.ProductSourceInternal
	if psi == nil {
		return cm.Error(cm.Internal, cm.F("ProductSource is not external: %v", sourceID), nil)
	}

	extra := ps.ExtraInfo
	if extra == nil || extra.DefaultBranchID == "" {
		return cm.Error(cm.Internal, cm.F("Default branch id is empty (SourceID=%v)", sourceID), nil)
	}

	// We can `Add()` inside a defer because the top level functions already
	// call `Recurrent()`. So we can make sure that the job will always be added
	// back to the queue. When re-`Added()`, the scheduler will update priority
	// of the job inside its queue.
	defer func() {
		if _err != nil {
			p.AddAfter(id, defaultErrRecurr, syncProductSource).
				Recurrent(defaultErrRecurr)
		} else {
			p.AddAfter(id, defaultRecurrent, syncProductSource).
				Recurrent(defaultRecurrent)
		}

		switch cm.ErrorCode(_err) {
		case cm.NoError, cm.SkipSync:
		default:
			msg := fmt.Sprintf(`
ERROR: kiotviet-sync-service (%vms)
–– SUPPLIER: %v
%v
`, time.Now().Sub(t0)/time.Millisecond,
				ps.ID, _err)
			bus.Dispatch(ctx, &telebot.SendMessageCommand{Message: msg})
		}
	}()

	conn := &Connection{
		RetailerID: psi.Secret.RetailerID,
		TokenStr:   psi.AccessToken,
		ExpiresAt:  psi.ExpiresAt,
	}
	conn, err := getOrRenewToken(ctx, sourceID, conn)
	if err != nil {
		return err
	}

	{
		sstate, ok := gStates.Get(sourceID)
		if !ok {
			ensureCmd := &EnsureWebhooksCommand{conn, sourceID}
			if err := bus.Dispatch(ctx, ensureCmd); err != nil {
				ll.Error("Unable to register webhooks", l.Int64("SourceID", sourceID))
				return err
			}

			sstate = SupplierSyncState{
				ID:      sourceID,
				Webhook: true,
			}
			gStates.Set(sourceID, sstate)
		}
	}

	chErr := make(chan error, 2)
	go func() {
		state, err := recoverState(ps.SyncStateProducts)
		if err != nil {
			ll.Warn("Invalid sync state for products", l.Object("psi", psi))
		}
		syncCmd := &SyncProductsCommand{
			SourceID:   sourceID,
			BranchID:   extra.DefaultBranchID,
			Connection: conn,
			SyncState:  *state.NextState,
		}
		if err := syncProducts(syncCmd); err != nil {
			ll.Error("Sync products", l.Int64("supplier", sourceID), l.Error(err))
			fmt.Printf("%+v", err)
			if cm.IsTrace(err) || cm.ErrorCode(err) == cm.Internal {
				bus.PrintAllStack(ctx, false)
			}
		}
		chErr <- err
	}()
	go func() {
		state, err := recoverState(ps.SyncStateCategories)
		if err != nil {
			ll.Error("Invalid sync state for categories", l.Object("psi", psi))
		}
		syncCmd := &SyncCategoriesCommand{
			SourceID:   sourceID,
			Connection: conn,
			SyncState:  *state.NextState,
		}
		if err := syncCategories(syncCmd); err != nil {
			ll.Error("Sync categories", l.Int64("supplier", sourceID), l.Error(err))
			fmt.Printf("%+v", err)
			if cm.IsTrace(err) || cm.ErrorCode(err) == cm.Internal {
				bus.PrintAllStack(ctx, false)
			}
		}
		chErr <- err
	}()

	// TODO(qv): Save supplier sync state: error

	var errs cm.ErrorCollector
	errs.Collect(<-chErr, <-chErr)
	return errs.Any()
}

func recoverState(s []byte) (ss SavedSyncState, err error) {
	if len(s) != 0 {
		err = json.Unmarshal(json.RawMessage(s), &ss)
	}
	if ss.NextState == nil {
		ss.NextState = &SyncState{}
	}
	if ss.NextState.Page <= 0 {
		ss.NextState.Page = 1
	}
	if ss.NextState.Since.IsZero() {
		ss.NextState.Since = defaultZeroDate
	}
	return
}
