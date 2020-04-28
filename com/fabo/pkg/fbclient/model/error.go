package model

import (
	"o.o/capi/dot"
)

type FacebookError struct {
	Message        dot.NullString
	Type           dot.NullString
	Code           dot.NullInt
	ErrorSubcode   dot.NullInt
	ErrorUserTitle dot.NullString
	FbtraceId      dot.NullString
}
