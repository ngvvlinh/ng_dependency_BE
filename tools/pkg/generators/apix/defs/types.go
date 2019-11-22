package defs

import (
	"etop.vn/backend/tools/pkg/generators/cq"
)

type Service struct {
	Name     string
	BasePath string
	APIPath  string
	Methods  []*Method
}

type Method = cq.HandlerDef

type Message = cq.Message
