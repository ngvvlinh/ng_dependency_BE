package testdst

import (
	"o.o/backend/zexp/etl/tests/testsrc"
)

// +gen:convert: o.o/backend/zexp/etl/tests/testdst  -> o.o/backend/zexp/etl/tests/testsrc

func ConvertAccount(in *testsrc.Account, out *Account) {
	convert_testsrc_Account_Account(in, out)
	out.FullName = in.FirstName + " " + in.LastName
}
