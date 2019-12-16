package testdst

import (
	"etop.vn/backend/zexp/etl/tests/testsrc"
)

// +gen:convert: etop.vn/backend/zexp/etl/tests/testdst -> etop.vn/backend/zexp/etl/tests/testsrc

func ConvertAccount(in *testsrc.Account, out *Account) {
	convert_testsrc_Account_Account(in, out)
	out.FullName = in.FirstName + " " + in.LastName
}
