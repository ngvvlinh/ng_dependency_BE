package bus

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/k0kubun/pp"

	"etop.vn/backend/pkg/common/logline"
)

type NodeContext struct {
	context.Context `json:"-"`

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

// Ctx is shorthand for NewRootContext(context.Background()) for quickly insert
// context to code.
func Ctx() *NodeContext {
	return NewRootContext(context.Background())
}

func NewRootContext(ctx context.Context) *NodeContext {
	nodes := make([]*NodeContext, 1, 16)
	node := &NodeContext{
		Context:     ctx,
		Nodes:       &nodes,
		ParentIndex: -1,
	}
	nodes[0] = node
	return node
}

func (n *NodeContext) WithMessage(msg interface{}) *NodeContext {
	node := &NodeContext{
		Context:     n.Context,
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
			pp.Println(node.Message)
		} else {
			pp.Println("error = "+node.Error.Error(), node.Message)
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
				pp.Println(n.Message)
			} else {
				pp.Println("error = "+n.Error.Error(), n.Message)
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
