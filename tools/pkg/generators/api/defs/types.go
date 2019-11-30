package defs

import (
	"go/types"
)

type Kind string

const (
	KindQuery     = "QueryService"
	KindAggregate = "Aggregate"
	KindService   = "Service"
)

type Service struct {
	Kind    Kind
	Name    string
	APIPath string
	Methods []*Method
}

type Method struct {
	Name     string
	Comment  string
	Method   *types.Func
	Request  *Message
	Response *Message
}

type Message struct {
	Items ArgItems
}

type ArgItems []*ArgItem

type ArgItem struct {
	Inline bool
	Name   string
	Type   types.Type
	Var    *types.Var
	Ptr    bool
	Struct *types.Struct
}

type NodeType int

const (
	NodeNone = iota
	NodeField
	NodeStartInline
	NodeEndInline
)

type WalkFunc func(node NodeType, name string, field *types.Var, tag string) error

func (args ArgItems) Walk(fn WalkFunc) error {
	for _, arg := range args {
		if arg.Inline {
			s := arg.Struct
			if err := fn(NodeStartInline, arg.Name, arg.Var, ""); err != nil {
				return err
			}
			for i, n := 0, s.NumFields(); i < n; i++ {
				field := s.Field(i)
				if err := fn(NodeField, field.Name(), field, s.Tag(i)); err != nil {
					return err
				}
			}
			if err := fn(NodeEndInline, arg.Name, arg.Var, ""); err != nil {
				return err
			}
		} else {
			if err := fn(NodeField, arg.Name, arg.Var, ""); err != nil {
				return err
			}
		}
	}
	return nil
}
