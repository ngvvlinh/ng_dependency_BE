package one

// +gen
// +gen:sample=10

type A struct {
}

// +gen:b: this directive should be ignored
type B int

// +gen:last: 20: number:int * x
