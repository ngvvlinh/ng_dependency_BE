package meta

import "context"

type Msg interface{}

type Bus interface {
	Dispatch(ctx context.Context, msg Msg) error
	DispatchAll(ctx context.Context, msgs ...Msg) error
}

type EventBus interface {
	Publish(ctx context.Context, msg Msg) error
}
