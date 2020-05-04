package bus

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/k0kubun/pp"

	"o.o/common/xerrors/logline"
)

type WithValuer interface {
	WithValue(key, val interface{})
	ResetValue(key interface{})
}

type KV struct{ key, val interface{} }

type NodeContext struct {
	context.Context `json:"-"`

	Values      *[]KV
	Nodes       *[]*NodeContext `json:"-"`
	Index       int             `json:"-"`
	Parent      *NodeContext    `json:"-"`
	ParentIndex int             `json:"parent"`

	Message interface{}       `json:"msg"`
	Logs    []logline.LogLine `json:"logs"`
	Start   time.Time         `json:"-"`
	Time    time.Duration     `json:"t"`

	// It's ignored because the middleware will store the error
	Error error `json:"-"`
}

var _ WithValuer = &NodeContext{}

// Ctx is shorthand for NewRootContext(context.Background()) for quickly insert
// context to code.
func Ctx() context.Context {
	return NewRootContext(context.Background())
}

type ctxKey struct{}

func NewRootContext(ctx context.Context) context.Context {
	if node := GetContext(ctx); node != nil {
		return ctx
	}

	nodes := make([]*NodeContext, 1, 16)
	values := make([]KV, 0, 16)
	node := &NodeContext{
		Context:     ctx,
		Values:      &values,
		Nodes:       &nodes,
		ParentIndex: -1,
	}
	nodes[0] = node
	return node
}

func GetContext(ctx context.Context) *NodeContext {
	node, _ := ctx.Value(ctxKey{}).(*NodeContext)
	return node
}

func WithMessage(ctx context.Context, msg interface{}) *NodeContext {
	node := GetContext(ctx)
	if node == nil {
		panic("must be bus.NodeContext")
	}
	newNode := &NodeContext{
		Context:     ctx,
		Values:      node.Values,
		Nodes:       node.Nodes,
		Parent:      node,
		Message:     msg,
		Index:       len(*node.Nodes),
		ParentIndex: node.Index,
		Start:       time.Now(),
	}
	*node.Nodes = append(*node.Nodes, newNode)
	return newNode
}

func WithValue(ctx context.Context, key, val interface{}) {
	node := GetContext(ctx)
	if node == nil {
		panic("must be bus.NodeContext")
	}
	*node.Values = append(*node.Values, KV{key, val})
}

func ResetValue(ctx context.Context, key interface{}) {
	node := GetContext(ctx)
	if node == nil {
		panic("must be bus.NodeContext")
	}
	node.ResetValue(key)
}

func (n *NodeContext) Value(key interface{}) interface{} {
	if key == (ctxKey{}) {
		return n
	}
	values := *n.Values
	for i := len(values) - 1; i >= 0; i-- {
		kv := values[i]
		if kv.key == key {
			return kv.val
		}
	}
	return n.Context.Value(key)
}

func (n *NodeContext) WithValue(key, val interface{}) {
	*n.Values = append(*n.Values, KV{key, val})
}

func (n *NodeContext) ResetValue(key interface{}) {
	values := *n.Values
	for i := range values {
		if values[i].key == key {
			values[i].val = nil
		}
	}
}

func (n *NodeContext) AppendLog(line logline.LogLine) {
	n.Logs = append(n.Logs, line)
}

func GetStack(ctx context.Context) []*NodeContext {
	node := GetContext(ctx)
	if node == nil {
		return nil
	}

	s := make([]*NodeContext, 0, len(*node.Nodes))
	for node.Parent != nil {
		s = append(s, node)
		node = node.Parent
	}
	return s
}

func GetAllStack(ctx context.Context) []*NodeContext {
	node := GetContext(ctx)
	if node == nil {
		return nil
	}
	return *node.Nodes
}

func PrintStack(ctx context.Context) {
	node := GetContext(ctx)
	if node == nil {
		debug.PrintStack()
		log.Println("Must be bus.Context")
		return
	}

	for node.Parent != nil {
		if node.Error == nil {
			_, _ = pp.Println(node.Message)
		} else {
			_, _ = pp.Println("error = "+node.Error.Error(), node.Message)
		}
		node = node.Parent
	}
}

func PrintErrorStack(ctx context.Context) {
	node := GetContext(ctx)
	if node == nil {
		debug.PrintStack()
		log.Println("Must be bus.Context")
		return
	}

	nodes := *node.Nodes
	for i := len(nodes) - 1; i >= 0; i-- {
		n := nodes[i]
		if n.Error != nil {
			node = n
			break
		}
	}
	PrintStack(node)
}

func PrintAllStack(ctx context.Context, expanded bool) {
	node := GetContext(ctx)
	if node == nil {
		debug.PrintStack()
		log.Println("Must be bus.Context")
		return
	}

	nodes := *node.Nodes
	idents := make([]int, len(nodes))
	for i := 1; i < len(nodes); i++ {
		n := nodes[i]
		ident := 0
		for j := i - 1; j >= 0; j-- {
			if n.Parent == nodes[j] {
				ident = idents[j]
				idents[i] = ident + 1
				break
			}
		}
		for i := 0; i < ident; i++ {
			fmt.Print(".   ")
		}

		if expanded {
			if n.Error == nil {
				_, _ = pp.Println(n.Message)
			} else {
				_, _ = pp.Println("error = "+n.Error.Error(), n.Message)
			}
		} else {
			msgName := reflect.TypeOf(n.Message).Elem().Name()
			if n.Error == nil {
				fmt.Println(msgName)
			} else {
				fmt.Println(msgName, "error="+n.Error.Error())
			}
		}
	}
}
