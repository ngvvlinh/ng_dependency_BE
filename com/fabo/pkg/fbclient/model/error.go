package model

import (
	"etop.vn/capi/dot"
)

type FacebookError struct {
	Message        dot.NullString
	Type           dot.NullString
	Code           dot.NullInt
	ErrorSubcode   dot.NullInt
	ErrorUserTitle dot.NullString
	FbtraceId      dot.NullString
}
