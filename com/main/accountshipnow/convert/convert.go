package convert

import (
	"o.o/api/main/accountshipnow"
	accountshipnowmodel "o.o/backend/com/main/accountshipnow/model"
)

// +gen:convert: o.o/backend/com/main/accountshipnow/model -> o.o/api/main/accountshipnow
// +gen:convert: o.o/api/main/accountshipnow

func XAccountAhamove(in *accountshipnowmodel.ExternalAccountAhamove) *accountshipnow.ExternalAccountAhamove {
	res := &accountshipnow.ExternalAccountAhamove{}
	convert_accountshipnowmodel_ExternalAccountAhamove_accountshipnow_ExternalAccountAhamove(in, res)
	return res
}
