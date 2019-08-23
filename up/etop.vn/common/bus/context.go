package bus

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/k0kubun/pp"

	"etop.vn/common/xerrors/logline"
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
func Ctx() *NodeContext {
	return NewRootContext(context.Background())
}

func NewRootContext(ctx context.Context) *NodeContext {
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

func (n *NodeContext) WithMessage(msg interface{}) *NodeContext {
	node := &NodeContext{
		Context:     n.Context,
		Values:      n.Values,
		Nodes:       n.Nodes,
		Parent:      n,
		Message:     msg,
		Index:       len(*n.Nodes),
		ParentIndex: n.Index,
		Start:       time.Now(),
	}
	*node.Nodes = append(*node.Nodes, node)
	return node
}

func (n *NodeContext) Value(key interface{}) interface{} {
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

func GetStack(ctx context.Context) []*NodeContext {
	node, ok := ctx.(*NodeContext)
	if !ok {
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
	node, ok := ctx.(*NodeContext)
	if !ok {
		return nil
	}
	return *node.Nodes
}

func PrintStack(ctx context.Context) {
	node, ok := ctx.(*NodeContext)
	if !ok {
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
	node, ok := ctx.(*NodeContext)
	if !ok {
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
	node, ok := ctx.(*NodeContext)
	if !ok {
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
			if n.Error == nil {
				fmt.Println(reflect.TypeOf(n.Message).Elem().Name())
			} else {
				fmt.Println(reflect.TypeOf(n.Message).Elem().Name(), "error="+n.Error.Error())
			}
		}
	}
}
