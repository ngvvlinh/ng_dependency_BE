package convert

import (
	"fmt"
	"time"

	"o.o/api/main/refund"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/refund/model  -> o.o/api/main/refund
// +gen:convert:  o.o/api/main/refund

const (
	MaxCodeNorm = 999999
	codePrefix  = "DTH"
)

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func convertCreate(in *refund.CreateRefundArgs, out *refund.Refund) *refund.Refund {
	if in == nil {
		return nil
	}
	if out == nil {
		out = &refund.Refund{}
	}
	apply_refund_CreateRefundArgs_refund_Refund(in, out)
	out.Status = status3.Z
	out.ID = cm.NewID()
	out.UpdatedBy = out.CreatedBy
	out.CreatedAt = time.Now()
	out.UpdatedAt = time.Now()
	return out
}
