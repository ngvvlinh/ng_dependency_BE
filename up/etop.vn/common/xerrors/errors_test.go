package xerrors

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"etop.vn/common/l"
)

func TestCustomErrorCode(t *testing.T) {
	for k, v := range mapCustomCodes {
		if v == nil {
			t.Errorf("Missing custom error code for %v", k)
		}
	}
}

func TestErrorInterface(t *testing.T) {
	err := errors.New("ABC")
	if _, ok := err.(IError); !ok {
		t.Error("Must satisfy errors interface")
	}
}

func TestError(t *testing.T) {
	cause := fmt.Errorf("Foo")

	t.Run("Simple", func(t *testing.T) {
		err := newError(true, true, NotFound, "Simple error", nil)
		fmt.Println("Println", err)
		fmt.Println("Message", err.Message)
		fmt.Println("Cause  ", err.Err)
		fmt.Printf("Format: %+v\n", err)
	})

	t.Run("Wrap", func(t *testing.T) {
		err := newError(true, true, NotFound, "Wrap error", cause)
		fmt.Println("Println", err)
		fmt.Println("Message", err.Message)
		fmt.Println("Cause  ", err.Err)
		fmt.Printf("Format: %+v\n", err)
	})

	t.Run("Multiple levels", func(t *testing.T) {
		err := newError(true, true, NotFound, "Error A", cause)
		err = newError(true, true, AlreadyExists, "Error B", err)
		err = newError(true, true, InvalidArgument, "Error C", err)
		fmt.Println("Println", err)
		fmt.Println("Message", err.Message)
		fmt.Println("Cause  ", err.Err)
		fmt.Printf("Format: %+v\n", err)
	})
}

func TestMapError(t *testing.T) {
	err1 := Error(NotFound, "Not Found", nil)
	err2 := MapError(err1).
		Map(AlreadyExists, InvalidArgument, "Invalid").
		Map(NotFound, FailedPrecondition, "Failed").
		Default(Internal, "Unexpected")

	xerr := interface{}(err2).(*APIError)
	assert.Equal(t, xerr.Code, FailedPrecondition)
	assert.Equal(t, xerr.Message, "Failed")
}

func TestErrorJSON(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		err := Error(NotFound, "Foo", nil)

		v := assertJSONError(t, err)
		assert.Equal(t, &jsonError{
			Code: "not_found",
			Msg:  "Foo",
			Logs: json.RawMessage([]byte(`[]`)),
		}, v)
	})

	t.Run("Simple with trace", func(t *testing.T) {
		err := ErrorTrace(NotFound, "Foo", nil)

		v := assertJSONError(t, err)
		assert.Equal(t, &jsonError{
			Code:  "not_found",
			Msg:   "Foo",
			Logs:  json.RawMessage([]byte(`[]`)),
			Stack: v.Stack,
		}, v)
		assert.Equal(t, 0, strings.Index(v.Stack,
			"\netop.vn/common/xerrors"))
	})

	t.Run("Advanced", func(t *testing.T) {
		cause := errors.New("os: file not found")
		err := ErrorTrace(NotFound, "Foo", cause).
			Log("Oops!", l.Int("int", 1), l.String("str", "._."), l.Any("iface", nil)).
			Log(">.<")
		{
			v := assertJSONError(t, err)
			assert.Equal(t, &jsonError{
				Code:  "not_found",
				Msg:   "Foo",
				Err:   "os: file not found",
				Logs:  v.Logs,
				Stack: v.Stack,
			}, v)
		}

		err = Error(FailedPrecondition, "Bar", err).
			Log("Baz!!")
		{
			v := assertJSONError(t, err)
			assert.Equal(t, &jsonError{
				Code:  "failed_precondition",
				Msg:   "Bar",
				Err:   "os: file not found",
				Orig:  "Foo",
				Logs:  v.Logs,
				Stack: v.Stack,
			}, v)
		}

		err = Error(WrongPassword, "Quix", err)
		{
			v := assertJSONError(t, err)
			assert.Equal(t, &jsonError{
				Code:  "unauthenticated",
				XCode: "wrong_password",
				Err:   "os: file not found",
				Msg:   "Quix",
				Orig:  "Foo",
				Logs:  v.Logs,
				Stack: v.Stack,
			}, v)
			logs := string(v.Logs)
			assert.True(t, strings.Contains(logs, `"str":"._."`))
			assert.True(t, strings.Contains(logs, `"int":1`))
			assert.True(t, strings.Contains(logs, `"iface":null`))
			assert.True(t, strings.Contains(logs, `"@msg":"\u003e.\u003c"`))
			assert.True(t, strings.Contains(logs, `"@msg":"Oops!"`))
			assert.True(t, strings.Contains(logs, `"@msg":"Baz!!"`))
		}
	})
}

type jsonError struct {
	Code  string          `json:"code"`
	XCode string          `json:"xcode"`
	Err   string          `json:"err"`
	Msg   string          `json:"msg"`
	Orig  string          `json:"orig"`
	Logs  json.RawMessage `json:"logs"`
	Stack string          `json:"stack"`
}

func assertJSONError(t *testing.T, err *APIError) *jsonError {
	data, jsonErr := err.MarshalJSON()
	assert.NoError(t, jsonErr)

	data, jsonErr = json.Marshal(err)
	assert.NoError(t, jsonErr)

	var v jsonError
	jsonErr = json.Unmarshal(data, &v)
	if jsonErr != nil {
		t.Errorf("Got error while decoding JSON: %v", jsonErr)
	}
	return &v
}
