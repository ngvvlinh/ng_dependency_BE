package mock

import (
	"fmt"

	sq2 "etop.vn/backend/pkg/common/sq"
)

// ErrorMock ...
type ErrorMock struct {
	Err    error
	Entry  *sq2.LogEntry
	Called int
}

// Error ...
type Error struct {
	Err   error
	Entry *sq2.LogEntry
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// Reset ...
func (m *ErrorMock) Reset() {
	fmt.Println()
	m.Err = nil
	m.Entry = nil
	m.Called = 0
}

// Mock ...
func (m *ErrorMock) Mock(err error, entry *sq2.LogEntry) error {
	m.Called++
	m.Err, m.Entry = err, entry

	if err == nil {
		return nil
	}
	return &Error{err, entry}
}
