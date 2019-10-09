package capi

import (
	"context"
)

type Bus interface {
	Dispatch(ctx context.Context, msg interface{}) error
	DispatchAll(ctx context.Context, msgs ...interface{}) error
}

type Event interface {
	GetTopic() string
}

type EventBus interface {
	Publish(ctx context.Context, msg Event) error
}
