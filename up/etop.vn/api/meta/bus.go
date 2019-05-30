package meta

import "context"

type Bus interface {
	Dispatch(ctx context.Context, msg interface{}) error
	DispatchAll(ctx context.Context, msgs ...interface{}) error
}

type EventBus interface {
	Publish(ctx context.Context, msg interface{}) error
}
