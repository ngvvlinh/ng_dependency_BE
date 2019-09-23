package test

import (
	"strconv"
)

// +gen:convert: etop.vn/backend/tools/pkg/generators/convert/test

type A struct {
	Value int
}

// +convert:type=A
type B struct {
	Value string
}

func ConvertAB(a *A, b *B) {
	b.Value = strconv.Itoa(a.Value)
}
