package convert

import (
	"o.o/api/main/bankstatement"
	bankstatementmodel "o.o/backend/com/main/bankstatement/model"
)

// +gen:convert: o.o/backend/com/main/bankstatement/model -> o.o/api/main/bankstatement
// +gen:convert: o.o/api/main/bankstatement

func BankStatementToModel(in *bankstatement.BankStatement) (out *bankstatementmodel.BankStatement) {
	if out == nil {
		out = &bankstatementmodel.BankStatement{}
	}
	convert_bankstatement_BankStatement_bankstatementmodel_BankStatement(in, out)
	out.OtherInfo = in.OtherInfo

	return out
}
