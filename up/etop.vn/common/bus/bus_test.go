package bus

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type TestQuery struct {
	Value string
	Resp  string
}

type TestQueryA struct {
	Value string
}

type TestQueryB struct {
	Value string
}

func TestQueryHandlerReturnsError(t *testing.T) {
	bus := New()

	bus.AddHandler(func(ctx context.Context, query *TestQuery) error {
		return errors.New("handler error")
	})

	ctx := NewRootContext(context.Background())
	err := bus.Dispatch(ctx, &TestQuery{})

	if err == nil {
		t.Fatal("Send query failed " + err.Error())
	} else {
		t.Log("Handler error received ok")
	}
}

func TestQueryHandlerReturn(t *testing.T) {
	bus := New()

	bus.AddHandler(func(ctx context.Context, q *TestQuery) error {
		q.Resp = "hello from handler"
		return nil
	})

	ctx := NewRootContext(context.Background())
	query := &TestQuery{}
	err := bus.Dispatch(ctx, query)

	if err != nil {
		t.Fatal("Send query failed " + err.Error())
	} else if query.Resp != "hello from handler" {
		t.Fatal("Failed to get response from handler")
	}
}

func TestQueryMockHandlerReturnsError(t *testing.T) {
	bus := New()

	bus.MockHandler(func(q *TestQuery) error {
		return errors.New("test handler error")
	})

	ctx := context.Background()
	err := bus.Dispatch(ctx, &TestQuery{})

	if err == nil {
		t.Fatal("Send query failed " + err.Error())
	} else {
		t.Log("Handler error received ok")
	}
}

func TestQueryMockHandlerReturn(t *testing.T) {
	bus := New()

	bus.MockHandler(func(q *TestQuery) error {
		q.Resp = "hello from test handler"
		return nil
	})

	ctx := context.Background()
	query := &TestQuery{}
	err := bus.Dispatch(ctx, query)

	if err != nil {
		t.Fatal("Send query failed " + err.Error())
	} else if query.Resp != "hello from test handler" {
		t.Fatal("Failed to get response from handler")
	}
}

func TestEventListeners(t *testing.T) {
	bus := New()
	count := 0

	bus.AddEventListener(func(_ context.Context, query *TestQuery) error {
		count += 1
		return nil
	})

	bus.AddEventListener(func(_ context.Context, query *TestQuery) error {
		count += 10
		return nil
	})

	err := bus.Publish(Ctx(), &TestQuery{})

	if err != nil {
		t.Fatal("Publish event failed " + err.Error())
	} else if count != 11 {
		t.Fatal(fmt.Sprintf("Publish event failed, listeners called: %v, expected: %v", count, 11))
	}
}

func TestPrintStack(t *testing.T) {
	bus := New()

	bus.AddHandler(func(ctx context.Context, query *TestQuery) error {
		bus.Dispatch(ctx, &TestQueryA{"A1"})
		return bus.Dispatch(ctx, &TestQueryA{"A2"})
	})

	bus.AddHandler(func(ctx context.Context, query *TestQueryA) error {
		bus.Dispatch(ctx, &TestQueryB{query.Value + "-B1"})
		return bus.Dispatch(ctx, &TestQueryB{query.Value + "-B2"})
	})

	bus.AddHandler(func(ctx context.Context, query *TestQueryB) error {
		if query.Value == "A2-B2" {
			PrintStack(ctx)
			return errors.New("Error at A2-B2")
		}
		return nil
	})

	ctx := NewRootContext(context.Background())
	err := bus.Dispatch(ctx, &TestQuery{})
	if err != nil {
		PrintAllStack(ctx, false)
		PrintAllStack(ctx, true)
	}
}
