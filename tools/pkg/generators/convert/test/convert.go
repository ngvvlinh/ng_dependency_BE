package test

import (
	"strconv"
)

// +gen:convert: etop.vn/backend/tools/pkg/generators/convert/test

type S string

type A struct {
	Value   int
	Int     int64
	String  string
	Strings []string

	// C0 and C1 have custom conversions

	C  *C0
	Cs []*C0

	// D0 and D1 have no custom conversions

	D  *D0
	Ds []*D0

	// E is just a struct without any conversions

	E   E
	Ep  *E
	Es  []E
	Eps []*E
}

// +convert:type=A
type B struct {
	Value   string
	Int     int32
	String  S
	Strings []string

	C   *C1
	Cs  []*C1
	D   *D1
	Ds  []*D1
	E   E
	Ep  *E
	Es  []E
	Eps []*E
}

type C0 struct {
	Value int
}

// +convert:type=C0
type C1 struct {
	Value string
}

// +convert:type=C0
type C2 struct {
	C0
	X, Y, Z int
}

// +convert:type=C0
type C3 struct {
	*C0
	X, Y, Z int
}

type D0 struct {
	Value string
}

// +convert:type=D0
type D1 struct {
	Value string
}

type E struct {
	Value string
}

func ConvertAB(a *A, b *B) {
	convert_A_B(a, b)
	b.Value = strconv.Itoa(a.Value)
}

func ConvertC01(c0 *C0) *C1 {
	if c0 == nil {
		return nil
	}
	var c1 C1
	convert_C0_C1(c0, &c1)
	c1.Value = strconv.Itoa(c0.Value)
	return &c1
}

func ConvertC10(c1 *C1, c0 *C0) *C0 {
	if c1 == nil {
		return nil
	}
	convert_C1_C0(c1, c0)
	c0.Value, _ = strconv.Atoi(c1.Value)
	return c0
}
