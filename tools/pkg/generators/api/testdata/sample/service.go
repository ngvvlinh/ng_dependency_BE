package sample

import "context"

// +gen:apix
// +gen:apix:base-path=/test
// +gen:swagger:doc-path=zext/test/sample

// +apix:path=/sample-service
type SampleService interface {

	// +apix:path=/helloWorld
	HelloWorld(context.Context, *Empty) (*Response, error)

	// +apix:path=/legacy/something
	DoSomething(context.Context, *Request) (*Empty, error)
}

type Empty struct{}

func (_ *Empty) String() string { panic(".") }

type Request struct {
	Code int `json:"code"`

	Something string `json:"something"`
}

func (_ *Request) String() string { panic(".") }

type Response struct {
	Message string `json:"message"`
}

func (_ *Response) String() string { panic(".") }
