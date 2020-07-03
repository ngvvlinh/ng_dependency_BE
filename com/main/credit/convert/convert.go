package convert

import (
	"time"

	"o.o/api/main/credit"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/credit/model -> o.o/api/main/credit
// +gen:convert: o.o/api/main/credit

func createCreditVoucher(args *credit.CreateCreditArgs, out *credit.Credit) {
	apply_credit_CreateCreditArgs_credit_Credit(args, out)
	out.ID = cm.NewID()
	out.UpdatedAt = time.Now()
	out.CreatedAt = time.Now()
}
