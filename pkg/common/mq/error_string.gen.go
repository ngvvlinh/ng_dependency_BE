// Code generated by "stringer -type=Code -output=error_string.gen.go"; DO NOT EDIT.

package mq

import "fmt"

const _Code_name = "CodeOKCodeNotFoundCodeServiceUnavailableCodeMalformMessageCodeNoHandlerCodeNoTubeCodeTimeoutCodeInternalCodeUnknown"

var _Code_index = [...]uint8{0, 6, 18, 40, 58, 71, 81, 92, 104, 115}

func (i Code) String() string {
	if i < 0 || i >= Code(len(_Code_index)-1) {
		return fmt.Sprintf("Code(%d)", i)
	}
	return _Code_name[_Code_index[i]:_Code_index[i+1]]
}
