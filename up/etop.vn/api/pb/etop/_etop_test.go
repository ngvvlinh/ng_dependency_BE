package etop

import (
	"testing"

	"github.com/stretchr/testify/assert"

	status32 "etop.vn/api/pb/etop/etc/status3"
	status42 "etop.vn/api/pb/etop/etc/status4"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
)

func TestAccountType(t *testing.T) {
	tests := []struct {
		pb  AccountType
		tag int
		typ model.AccountType
	}{
		{AccountType_shop, model.TagShop, model.TypeShop},
		{AccountType_etop, model.EtopAccountID, model.TypeEtop},
	}
	for _, tt := range tests {
		assert.Equal(t, int(tt.pb), tt.tag, "Tag must be equal: %v", tt.pb)
		assert.Equal(t, tt.pb.String(), string(tt.typ), "Type string must be equal: %v", tt.pb)
	}
}

func TestStatus3Complement(t *testing.T) {
	var status status32.Status = 127
	assert.Equal(t, model.Status3(-1), *(status.ToModel()), "Convert 127 to -1")
	assert.Equal(t, 127, int(convertpb.Pb3(-1)), "Convert -1 back to 127")

	status = 0
	assert.Equal(t, model.Status3(0), *(status.ToModel()), "Leave 0 unchanged")
	assert.Equal(t, 0, int(convertpb.Pb3(0)), "Leave 0 unchanged")

	status = 1
	assert.Equal(t, model.Status3(1), *(status.ToModel()), "Leave 1 unchanged")
	assert.Equal(t, 1, int(convertpb.Pb3(1)), "Leave 1 unchanged")
}

func TestStatus4Complement(t *testing.T) {
	var status status42.Status = 127
	assert.Equal(t, model.Status4(-1), *(status.ToModel()), "Convert 127 to -1")
	assert.Equal(t, 127, int(convertpb.Pb4(-1)), "Convert -1 back to 127")

	status = 0
	assert.Equal(t, model.Status4(0), *(status.ToModel()), "Leave 0 unchanged")
	assert.Equal(t, 0, int(convertpb.Pb4(0)), "Leave 0 unchanged")

	status = 1
	assert.Equal(t, model.Status4(1), *(status.ToModel()), "Leave 1 unchanged")
	assert.Equal(t, 1, int(convertpb.Pb4(1)), "Leave 1 unchanged")
}
