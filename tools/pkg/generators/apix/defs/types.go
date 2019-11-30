package defs

import (
	"etop.vn/backend/tools/pkg/generators/api"
)

type Service struct {
	Name     string
	BasePath string
	APIPath  string
	Methods  []*Method
}

type Method = api.HandlerDef

type Message = api.Message
