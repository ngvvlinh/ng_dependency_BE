package capi

import (
	"context"
	"fmt"
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

// TODO: common interface
type Message interface {
	fmt.Stringer
}
