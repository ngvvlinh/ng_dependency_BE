// +build !generator

// Code generated by generator api. DO NOT EDIT.

package notify

import (
	context "context"

	capi "o.o/capi"
	dot "o.o/capi/dot"
)

type CommandBus struct{ bus capi.Bus }
type QueryBus struct{ bus capi.Bus }

func NewCommandBus(bus capi.Bus) CommandBus { return CommandBus{bus} }
func NewQueryBus(bus capi.Bus) QueryBus     { return QueryBus{bus} }

func (b CommandBus) Dispatch(ctx context.Context, msg interface{ command() }) error {
	return b.bus.Dispatch(ctx, msg)
}
func (b QueryBus) Dispatch(ctx context.Context, msg interface{ query() }) error {
	return b.bus.Dispatch(ctx, msg)
}

type CreateUserNotifySettingCommand struct {
	UserID        dot.ID
	DisableTopics []string

	Result *UserNotiSetting `json:"-"`
}

func (h AggregateHandler) HandleCreateUserNotifySetting(ctx context.Context, msg *CreateUserNotifySettingCommand) (err error) {
	msg.Result, err = h.inner.CreateUserNotifySetting(msg.GetArgs(ctx))
	return err
}

type DisableTopicCommand struct {
	UserID dot.ID
	Topic  string

	Result *UserNotiSetting `json:"-"`
}

func (h AggregateHandler) HandleDisableTopic(ctx context.Context, msg *DisableTopicCommand) (err error) {
	msg.Result, err = h.inner.DisableTopic(msg.GetArgs(ctx))
	return err
}

type EnableTopicCommand struct {
	UserID dot.ID
	Topic  string

	Result *UserNotiSetting `json:"-"`
}

func (h AggregateHandler) HandleEnableTopic(ctx context.Context, msg *EnableTopicCommand) (err error) {
	msg.Result, err = h.inner.EnableTopic(msg.GetArgs(ctx))
	return err
}

type GetOrCreateUserNotifySettingCommand struct {
	UserID        dot.ID
	DisableTopics []string

	Result *UserNotiSetting `json:"-"`
}

func (h AggregateHandler) HandleGetOrCreateUserNotifySetting(ctx context.Context, msg *GetOrCreateUserNotifySettingCommand) (err error) {
	msg.Result, err = h.inner.GetOrCreateUserNotifySetting(msg.GetArgs(ctx))
	return err
}

type GetUserNotifySettingQuery struct {
	UserID dot.ID

	Result *UserNotiSetting `json:"-"`
}

func (h QueryServiceHandler) HandleGetUserNotifySetting(ctx context.Context, msg *GetUserNotifySettingQuery) (err error) {
	msg.Result, err = h.inner.GetUserNotifySetting(msg.GetArgs(ctx))
	return err
}

// implement interfaces

func (q *CreateUserNotifySettingCommand) command()      {}
func (q *DisableTopicCommand) command()                 {}
func (q *EnableTopicCommand) command()                  {}
func (q *GetOrCreateUserNotifySettingCommand) command() {}

func (q *GetUserNotifySettingQuery) query() {}

// implement conversion

func (q *CreateUserNotifySettingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *CreateUserNotifySettingArgs) {
	return ctx,
		&CreateUserNotifySettingArgs{
			UserID:        q.UserID,
			DisableTopics: q.DisableTopics,
		}
}

func (q *CreateUserNotifySettingCommand) SetCreateUserNotifySettingArgs(args *CreateUserNotifySettingArgs) {
	q.UserID = args.UserID
	q.DisableTopics = args.DisableTopics
}

func (q *DisableTopicCommand) GetArgs(ctx context.Context) (_ context.Context, _ *DisableTopicArgs) {
	return ctx,
		&DisableTopicArgs{
			UserID: q.UserID,
			Topic:  q.Topic,
		}
}

func (q *DisableTopicCommand) SetDisableTopicArgs(args *DisableTopicArgs) {
	q.UserID = args.UserID
	q.Topic = args.Topic
}

func (q *EnableTopicCommand) GetArgs(ctx context.Context) (_ context.Context, _ *EnableTopicArgs) {
	return ctx,
		&EnableTopicArgs{
			UserID: q.UserID,
			Topic:  q.Topic,
		}
}

func (q *EnableTopicCommand) SetEnableTopicArgs(args *EnableTopicArgs) {
	q.UserID = args.UserID
	q.Topic = args.Topic
}

func (q *GetOrCreateUserNotifySettingCommand) GetArgs(ctx context.Context) (_ context.Context, _ *GetOrCreateUserNotifySettingArgs) {
	return ctx,
		&GetOrCreateUserNotifySettingArgs{
			UserID:        q.UserID,
			DisableTopics: q.DisableTopics,
		}
}

func (q *GetOrCreateUserNotifySettingCommand) SetGetOrCreateUserNotifySettingArgs(args *GetOrCreateUserNotifySettingArgs) {
	q.UserID = args.UserID
	q.DisableTopics = args.DisableTopics
}

func (q *GetUserNotifySettingQuery) GetArgs(ctx context.Context) (_ context.Context, _ *GetUserNotiSettingArgs) {
	return ctx,
		&GetUserNotiSettingArgs{
			UserID: q.UserID,
		}
}

func (q *GetUserNotifySettingQuery) SetGetUserNotiSettingArgs(args *GetUserNotiSettingArgs) {
	q.UserID = args.UserID
}

// implement dispatching

type AggregateHandler struct {
	inner Aggregate
}

func NewAggregateHandler(service Aggregate) AggregateHandler { return AggregateHandler{service} }

func (h AggregateHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) CommandBus {
	b.AddHandler(h.HandleCreateUserNotifySetting)
	b.AddHandler(h.HandleDisableTopic)
	b.AddHandler(h.HandleEnableTopic)
	b.AddHandler(h.HandleGetOrCreateUserNotifySetting)
	return CommandBus{b}
}

type QueryServiceHandler struct {
	inner QueryService
}

func NewQueryServiceHandler(service QueryService) QueryServiceHandler {
	return QueryServiceHandler{service}
}

func (h QueryServiceHandler) RegisterHandlers(b interface {
	capi.Bus
	AddHandler(handler interface{})
}) QueryBus {
	b.AddHandler(h.HandleGetUserNotifySetting)
	return QueryBus{b}
}
