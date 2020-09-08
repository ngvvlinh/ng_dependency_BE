package bus

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestQuery struct {
	Value string
	Resp  string
}

func (q *TestQuery) GetTopic() string { return "topic" }

type TestQueryA struct {
	Value string
}

type TestQueryB struct {
	Value string
}

func TestQueryHandlerReturnsError(t *testing.T) {
	b := New()
	b.AddHandler(func(ctx context.Context, query *TestQuery) error {
		return errors.New("handler error")
	})

	ctx := NewRootContext(context.Background())
	err := b.Dispatch(ctx, &TestQuery{})
	assert.Error(t, err)
}

func TestQueryHandlerReturn(t *testing.T) {
	b := New()
	b.AddHandler(func(ctx context.Context, q *TestQuery) error {
		q.Resp = "hello from handler"
		return nil
	})

	ctx := NewRootContext(context.Background())
	query := &TestQuery{}
	err := b.Dispatch(ctx, query)
	assert.NoError(t, err)
}

func TestQueryMockHandlerReturnsError(t *testing.T) {
	b := New()
	b.MockHandler(func(q *TestQuery) error {
		return errors.New("test handler error")
	})

	ctx := Ctx()
	err := b.Dispatch(ctx, &TestQuery{})
	assert.Error(t, err)
}

func TestQueryMockHandlerReturn(t *testing.T) {
	b := New()
	b.MockHandler(func(q *TestQuery) error {
		q.Resp = "hello from test handler"
		return nil
	})

	ctx := Ctx()
	query := &TestQuery{}
	err := b.Dispatch(ctx, query)
	assert.NoError(t, err)
}

func TestQueryMockHandlerContext(t *testing.T) {
	b := New()
	b.MockHandler(func(ctx context.Context, q *TestQuery) error {
		q.Resp = "hello from test handler"
		return nil
	})

	ctx := Ctx()
	query := &TestQuery{}
	err := b.Dispatch(ctx, query)
	assert.NoError(t, err)
}

func TestEventListeners(t *testing.T) {
	b := New()
	count := 0

	b.AddEventListener(func(_ context.Context, query *TestQuery) error {
		count += 1
		return nil
	})
	b.AddEventListener(func(_ context.Context, query *TestQuery) error {
		count += 10
		return nil
	})

	err := b.Publish(Ctx(), &TestQuery{})
	if err != nil {
		t.Fatal("Publish event failed " + err.Error())
	} else if count != 11 {
		t.Fatal(fmt.Sprintf("Publish event failed, listeners called: %v, expected: %v", count, 11))
	}
}

func TestPrintStack(t *testing.T) {
	b := New()
	b.AddHandler(func(ctx context.Context, query *TestQuery) error {
		_ = b.Dispatch(ctx, &TestQueryA{"A1"})
		return b.Dispatch(ctx, &TestQueryA{"A2"})
	})
	b.AddHandler(func(ctx context.Context, query *TestQueryA) error {
		_ = b.Dispatch(ctx, &TestQueryB{query.Value + "-B1"})
		return b.Dispatch(ctx, &TestQueryB{query.Value + "-B2"})
	})
	b.AddHandler(func(ctx context.Context, query *TestQueryB) error {
		if query.Value == "A2-B2" {
			PrintStack(ctx)
			return errors.New("error at A2-B2")
		}
		return nil
	})

	ctx := NewRootContext(context.Background())
	err := b.Dispatch(ctx, &TestQuery{})
	if err != nil {
		PrintAllStack(ctx, false)
		PrintAllStack(ctx, true)
	}
}
