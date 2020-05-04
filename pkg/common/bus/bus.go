package bus

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"o.o/capi"
)

var isTest = func() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}()

type HandlerFunc = interface{}
type CtxHandlerFunc func()
type Msg = interface{}

type Event = capi.Event

type Bus interface {
	Dispatch(ctx context.Context, msg Msg) error
	DispatchAll(ctx context.Context, msgs ...Msg) error
	Publish(ctx context.Context, msg Event) error

	AddHandler(handler HandlerFunc)
	AddHandlers(handlers ...HandlerFunc)
	MockHandler(handler HandlerFunc)
	AddEventListener(handler HandlerFunc)
	AddWildcardListener(handler HandlerFunc)
}

type EventRegistry interface {
	Publish(ctx context.Context, msg Event) error
	AddEventListener(handler interface{})
}

type InProcBus struct {
	handlers          map[reflect.Type]HandlerFunc
	listeners         map[reflect.Type][]HandlerFunc
	wildcardListeners []HandlerFunc

	expects map[reflect.Type]bool
}

var globalBus = New()

func New() Bus {
	bus := &InProcBus{}
	bus.handlers = make(map[reflect.Type]HandlerFunc)
	bus.listeners = make(map[reflect.Type][]HandlerFunc)
	bus.wildcardListeners = make([]HandlerFunc, 0)
	return bus
}

func (b *InProcBus) DispatchAll(ctx context.Context, msgs ...Msg) error {
	for _, msg := range msgs {
		if err := b.Dispatch(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (b *InProcBus) Dispatch(ctx context.Context, msg Msg) (_err error) {
	var msgType = reflect.TypeOf(msg).Elem()
	var handler = b.handlers[msgType]
	if handler == nil {
		return fmt.Errorf("bus: Handler not found for %s", msgType)
	}

	var params = make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(ctx)
	params[1] = reflect.ValueOf(msg)

	return call(ctx, msg, params, handler)
}

func (b *InProcBus) Publish(ctx context.Context, msg Event) (_err error) {
	var msgType = reflect.TypeOf(msg).Elem()
	var listeners = b.listeners[msgType]
	var params = make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(ctx)
	params[1] = reflect.ValueOf(msg)

	for _, listenerHandler := range listeners {
		if err := call(ctx, msg, params, listenerHandler); err != nil {
			return err
		}
	}
	for _, listenerHandler := range b.wildcardListeners {
		if err := call(ctx, msg, params, listenerHandler); err != nil {
			return err
		}
	}
	return nil
}

func call(ctx context.Context, msg Msg, params []reflect.Value, handler HandlerFunc) (_err error) {
	node := GetContext(ctx)
	if node != nil {
		newNode := WithMessage(ctx, msg)
		params[0] = reflect.ValueOf(newNode)
		defer func() {
			newNode.Error = _err
			newNode.Time = time.Since(newNode.Start)
		}()
	}
	ret := reflect.ValueOf(handler).Call(params)
	err, _ := ret[0].Interface().(error)
	return err
}

func (b *InProcBus) AddWildcardListener(handler HandlerFunc) {
	b.wildcardListeners = append(b.wildcardListeners, handler)
}

func (b *InProcBus) AddHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		panic(fmt.Sprintf("bus: Handler must be function (got %v)", handlerType.Name()))
	}
	if handlerType.NumIn() != 2 || handlerType.NumOut() != 1 {
		panic("bus: Handler must receive 2 params and return error")
	}
	if handlerType.In(0) != ctxType {
		panic(fmt.Sprintf("bus: Handler must receive the first param as context.Context (got %v)", handlerType.In(0)))
	}
	if handlerType.Out(0) != errType {
		panic(fmt.Sprintf("bus: Handler must return error (got %v)", handlerType.Out(0)))
	}

	queryType := handlerType.In(1).Elem()
	b.handlers[queryType] = handler
}

func (b *InProcBus) AddHandlers(handlers ...HandlerFunc) {
	for _, h := range handlers {
		b.AddHandler(h)
	}
}

var (
	ctxType = reflect.TypeOf((*context.Context)(nil)).Elem()
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

func (b *InProcBus) MockHandler(handler HandlerFunc) {
	if !isTest {
		panic("bus: Test handler must only be called during test")
	}

	handlerType := reflect.TypeOf(handler)
	if handlerType.NumOut() != 1 {
		panic("bus: Test handler must return 1 result")
	}
	if handlerType.Out(0) != errType {
		panic(fmt.Sprintf("bus: Test handler must return error (got %v)", handlerType.Out(0)))
	}
	switch {
	case handlerType.NumIn() == 1:
		queryType := handlerType.In(0).Elem()
		fnType := reflect.FuncOf(
			[]reflect.Type{ctxType, handlerType.In(0)},
			[]reflect.Type{errType},
			false,
		)
		wrapped := reflect.MakeFunc(fnType,
			func(args []reflect.Value) []reflect.Value {
				params := []reflect.Value{args[1]}
				return reflect.ValueOf(handler).Call(params)
			})
		b.handlers[queryType] = wrapped.Interface()

	case handlerType.NumIn() == 2:
		if handlerType.In(0) != ctxType {
			panic(fmt.Sprintf("bus: Test handler must receive the first param as context.Context (got %v)", handlerType.In(0)))
		}
		queryType := handlerType.In(1).Elem()
		b.handlers[queryType] = handler

	default:
		panic("bus: Test handler must receive 2 params and return error")
	}
}

func (b *InProcBus) AddEventListener(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	if handlerType.NumIn() != 2 || handlerType.NumOut() != 1 {
		panic("bus: Handler must receive 2 params and return error")
	}

	eventType := handlerType.In(1).Elem()
	_, exists := b.listeners[eventType]
	if !exists {
		b.listeners[eventType] = make([]HandlerFunc, 0)
	}
	b.listeners[eventType] = append(b.listeners[eventType], handler)
}

func (b *InProcBus) Expect(msg interface{}) {
	if b.expects == nil {
		b.expects = make(map[reflect.Type]bool)
	}
	typ := reflect.TypeOf(msg).Elem()
	b.expects[typ] = true
}

func typeName(t reflect.Type) string {
	return filepath.Base(t.PkgPath()) + "." + t.Name()
}

func (b *InProcBus) Validate() bool {
	exps := make([]reflect.Type, 0, len(b.expects))
	for name := range b.expects {
		exps = append(exps, name)
	}
	sort.Slice(exps, func(i, j int) bool {
		return typeName(exps[i]) < typeName(exps[j])
	})

	ok := true
	for _, typ := range exps {
		if b.handlers[typ] == nil {
			log.Println("bus: No handler for", typeName(typ))
			ok = false
		}
	}
	return ok
}

func (b *InProcBus) ValidateTypes(types []string) (res bool) {
	res = true
	for _, t := range types {
		ok := false
		for typ := range b.handlers {
			name := "*" + typ.PkgPath() + "." + typ.Name()
			if name == t {
				ok = true
				break
			}
		}
		if !ok {
			fmt.Println("Handler not found for:", t)
			res = false
		}
	}
	return
}

func Global() Bus {
	return globalBus
}

func AddHandler(implName string, handler HandlerFunc) {
	globalBus.AddHandler(handler)
}

func AddHandlers(implName string, handlers ...HandlerFunc) {
	for _, h := range handlers {
		AddHandler(implName, h)
	}
}

func MockHandler(handler HandlerFunc) {
	globalBus.MockHandler(handler)
}

func AddEventListener(handler HandlerFunc) {
	globalBus.AddEventListener(handler)
}

func AddWildcardListener(handler HandlerFunc) {
	globalBus.AddWildcardListener(handler)
}

func Dispatch(ctx context.Context, msg Msg) error {
	return globalBus.Dispatch(ctx, msg)
}

func ClearBusHandlers() {
	globalBus = New()
}

func Expect(msg interface{}) {
	globalBus.(*InProcBus).Expect(msg)
}

func Validate() bool {
	return globalBus.(*InProcBus).Validate()
}

func ValidateTypes(types []string) bool {
	return globalBus.(*InProcBus).ValidateTypes(types)
}
