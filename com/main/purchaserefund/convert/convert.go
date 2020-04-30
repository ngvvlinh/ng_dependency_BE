package convert

import (
	"fmt"
	"time"

	"o.o/api/main/purchaserefund"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/purchaserefund/model  -> o.o/api/main/purchaserefund
// +gen:convert: o.o/api/main/purchaserefund

const (
	MaxCodeNorm = 999999
	codePrefix  = "DTHN"
)

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func convertCreate(in *purchaserefund.CreatePurchaseRefundArgs, out *purchaserefund.PurchaseRefund) *purchaserefund.PurchaseRefund {
	if in == nil {
		return nil
	}
	if out == nil {
		out = &purchaserefund.PurchaseRefund{}
	}
	apply_purchaserefund_CreatePurchaseRefundArgs_purchaserefund_PurchaseRefund(in, out)
	out.Status = status3.Z
	out.ID = cm.NewID()
	out.UpdatedBy = out.CreatedBy
	out.CreatedAt = time.Now()
	out.UpdatedAt = time.Now()
	return out
}

func convertUpdate(in *purchaserefund.UpdatePurchaseRefundArgs, out *purchaserefund.PurchaseRefund) *purchaserefund.PurchaseRefund {
	if in == nil {
		return nil
	}
	if out == nil {
		out = &purchaserefund.PurchaseRefund{}
	}
	apply_purchaserefund_UpdatePurchaseRefundArgs_purchaserefund_PurchaseRefund(in, out)
	out.UpdatedBy = out.UpdatedBy
	out.UpdatedAt = time.Now()
	return out
}
