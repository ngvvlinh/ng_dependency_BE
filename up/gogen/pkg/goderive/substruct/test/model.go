package test

//go:generate ../../../../scripts/derive.sh
var _ = substructFoo1(&Foo1{}, &Foo{})

// Foo ...
type Foo struct {
	A    string
	I    int
	SS   []string
	PS   *string
	PSS  *[]string
	PSPS *[]*string

	XA string
}

// Foo1 ...
type Foo1 struct {
	A    string
	I    int
	SS   []string
	PS   *string
	PSS  *[]string
	PSPS *[]*string
}
