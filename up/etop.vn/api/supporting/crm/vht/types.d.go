// Code generated by gen-cmd-query. DO NOT EDIT.

package vht

import (
	context "context"
	time "time"

	meta "etop.vn/api/meta"
)

type Command interface{ command() }
type Query interface{ query() }
type CommandBus struct{ bus meta.Bus }
type QueryBus struct{ bus meta.Bus }

func (c CommandBus) Dispatch(ctx context.Context, msg Command) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c QueryBus) Dispatch(ctx context.Context, msg Query) error {
	return c.bus.Dispatch(ctx, msg)
}
func (c CommandBus) DispatchAll(ctx context.Context, msgs ...Command) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (c QueryBus) DispatchAll(ctx context.Context, msgs ...Query) error {
	for _, msg := range msgs {
		if err := c.bus.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

type CreateOrUpdateCallHistoryByCallIDCommand struct {
	Direction       int32
	CdrID           string
	CallID          string
	SipCallID       string
	SdkCallID       string
	Cause           string
	Q850Cause       string
	FromExtension   string
	ToExtension     string
	FromNumber      string
	ToNumber        string
	Duration        int32
	TimeStarted     time.Time
	TimeConnected   time.Time
	TimeEnded       time.Time
	RecordingPath   string
	RecordingUrl    string
	RecordFileSize  int32
	EtopAccountID   int64
	VtigerAccountID string

	Result *VhtCallLog `json:"-"`
}

type CreateOrUpdateCallHistoryBySDKCallIDCommand struct {
	Direction       int32
	CdrID           string
	CallID          string
	SipCallID       string
	SdkCallID       string
	Cause           string
	Q850Cause       string
	FromExtension   string
	ToExtension     string
	FromNumber      string
	ToNumber        string
	Duration        int32
	TimeStarted     time.Time
	TimeConnected   time.Time
	TimeEnded       time.Time
	RecordingPath   string
	RecordingUrl    string
	RecordFileSize  int32
	EtopAccountID   int64
	VtigerAccountID string

	Result *VhtCallLog `json:"-"`
}

type PingServerVhtCommand struct {
	Result struct {
	} `json:"-"`
}

type SyncVhtCallHistoriesCommand struct {
	SyncTime time.Time

	Result struct {
	} `json:"-"`
}

type GetCallHistoriesQuery struct {
	Paging     *meta.Paging
	TextSearch string

	Result *GetCallHistoriesResponse `json:"-"`
}

type GetLastCallHistoryQuery struct {
	Offset int32
	Limit  int32
	Sort   []string

	Result *VhtCallLog `json:"-"`
}

// implement interfaces

func (q *CreateOrUpdateCallHistoryByCallIDCommand) command()    {}
func (q *CreateOrUpdateCallHistoryBySDKCallIDCommand) command() {}
func (q *PingServerVhtCommand) command()                        {}
func (q *SyncVhtCallHistoriesCommand) command()                 {}
func (q *GetCallHistoriesQuery) query()                         {}
func (q *GetLastCallHistoryQuery) query()                       {}

// implement conversion

func (q *CreateOrUpdateCallHistoryByCallIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *VhtCallLog) {
	return ctx,
		&VhtCallLog{
			Direction:       q.Direction,
			CdrID:           q.CdrID,
			CallID:          q.CallID,
			SipCallID:       q.SipCallID,
			SdkCallID:       q.SdkCallID,
			Cause:           q.Cause,
			Q850Cause:       q.Q850Cause,
			FromExtension:   q.FromExtension,
			ToExtension:     q.ToExtension,
			FromNumber:      q.FromNumber,
			ToNumber:        q.ToNumber,
			Duration:        q.Duration,
			TimeStarted:     q.TimeStarted,
			TimeConnected:   q.TimeConnected,
			TimeEnded:       q.TimeEnded,
			RecordingPath:   q.RecordingPath,
			RecordingUrl:    q.RecordingUrl,
			RecordFileSize:  q.RecordFileSize,
			EtopAccountID:   q.EtopAccountID,
			VtigerAccountID: q.VtigerAccountID,
		}
}

func (q *CreateOrUpdateCallHistoryBySDKCallIDCommand) GetArgs(ctx context.Context) (_ context.Context, _ *VhtCallLog) {
	return ctx,
		&VhtCallLog{
			Direction:       q.Direction,
			CdrID:           q.CdrID,
			CallID:          q.CallID,
			SipCallID:       q.SipCallID,
			SdkCallID:       q.SdkCallID,
			Cause:           q.Cause,
			Q850Cause:       q.Q850Cause,
			FromExtension:   q.FromExtension,
			ToExtension:     q.ToExtension,
			FromNumber:      q.FromNumber,
			ToNumber:        q.ToNumber,
			Duration:        q.Duration,
			TimeStarted:     q.TimeStarted,
			TimeConnected:   q.TimeConnected,
			TimeEnded:       q.TimeEnded,
			RecordingPath:   q.RecordingPath,
			RecordingUrl:    q.RecordingUrl,
			RecordFileSize:  q.RecordFileSize,
			EtopAccountID:   q.EtopAccountID,
			VtigerAccountID: q.VtigerAccountID,
		}
}

func (q *PingServerVhtCommand) GetArgs(ctx context.Context) (_ context.Context, _ *meta.Empty) {
	return ctx,
		&meta.Empty{}
}

func (q *SyncVhtCallHistoriesCommand) GetArgs(ctx context.Context) (_ context.Context, _ *SyncVhtCallHistoriesArgs) {
	return ctx,
		&SyncVhtCallHistoriesArgs{
			SyncTime: q.SyncTime,
		}
}

func (q *GetCallHistoriesQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetCallHistoriesArgs) {
	return ctx,
		&GetCallHistoriesArgs{
			Paging:     q.Paging,
			TextSearch: q.TextSearch,
		}
}

func (q *GetLastCallHistoryQuery) GetArgs(ctx context.Context) (_ context.Context, _ meta.Paging) {
	return ctx,
		meta.Paging{
			Offset: q.Offset,
			Limit:  q.Limit,
			Sort:   q.Sort,
		}
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCreateOrUpdateCallHistoryByCallID)
	b.AddHandler(h.HandleCreateOrUpdateCallHistoryBySDKCallID)
	b.AddHandler(h.HandlePingServerVht)
	b.AddHandler(h.HandleSyncVhtCallHistories)
	return CommandBus{b}
}

func (h AggregateHandler) HandleCreateOrUpdateCallHistoryByCallID(ctx context.Context, msg *CreateOrUpdateCallHistoryByCallIDCommand) error {
	result, err := h.inner.CreateOrUpdateCallHistoryByCallID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandleCreateOrUpdateCallHistoryBySDKCallID(ctx context.Context, msg *CreateOrUpdateCallHistoryBySDKCallIDCommand) error {
	result, err := h.inner.CreateOrUpdateCallHistoryBySDKCallID(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h AggregateHandler) HandlePingServerVht(ctx context.Context, msg *PingServerVhtCommand) error {
	return h.inner.PingServerVht(msg.GetArgs(ctx))
}

func (h AggregateHandler) HandleSyncVhtCallHistories(ctx context.Context, msg *SyncVhtCallHistoriesCommand) error {
	return h.inner.SyncVhtCallHistories(msg.GetArgs(ctx))
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	meta.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetCallHistories)
	b.AddHandler(h.HandleGetLastCallHistory)
	return QueryBus{b}
}

func (h QueryServiceHandler) HandleGetCallHistories(ctx context.Context, msg *GetCallHistoriesQuery) error {
	result, err := h.inner.GetCallHistories(msg.GetArgs(ctx))
	msg.Result = result
	return err
}

func (h QueryServiceHandler) HandleGetLastCallHistory(ctx context.Context, msg *GetLastCallHistoryQuery) error {
	result, err := h.inner.GetLastCallHistory(msg.GetArgs(ctx))
	msg.Result = result
	return err
}
