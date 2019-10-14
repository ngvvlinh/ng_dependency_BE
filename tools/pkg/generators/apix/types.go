package apix

import "go/types"

type Service struct {
	Name    string
	APIPath string
	Methods []*Method
}

type Method struct {
	Name     string
	Comment  string
	Request  Message
	Response Message
}

type Message struct {
	Type    types.Type
	PkgPath string
	Name    string
}
