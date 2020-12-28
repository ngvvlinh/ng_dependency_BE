package _all

import (
	"o.o/api/main/purchaserefund"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbPurchaseRefunds(args []*purchaserefund.PurchaseRefund) []*shop.PurchaseRefund {
	var result []*shop.PurchaseRefund
	for _, value := range args {
		result = append(result, PbPurchaseRefund(value))
	}
	return result
}

func PbPurchaseRefund(args *purchaserefund.PurchaseRefund) *shop.PurchaseRefund {
	var result = &shop.PurchaseRefund{
		ID:              args.ID,
		ShopID:          args.ShopID,
		PurchaseOrderID: args.PurchaseOrderID,
		Note:            args.Note,
		Code:            args.Code,
		TotalAdjustment: args.TotalAdjustment,
		AdjustmentLines: args.AdjustmentLines,
		Lines:           PbPurchaseRefundLine(args.Lines),
		CreatedAt:       cmapi.PbTime(args.CreatedAt),
		UpdatedAt:       cmapi.PbTime(args.UpdatedAt),
		CancelledAt:     cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:     cmapi.PbTime(args.ConfirmedAt),
		CreatedBy:       args.CreatedBy,
		UpdatedBy:       args.UpdatedBy,
		CancelReason:    args.CancelReason,
		Status:          args.Status,
		TotalAmount:     args.TotalAmount,
		BasketValue:     args.BasketValue,
	}
	return result
}

func PbPurchaseRefundLine(args []*purchaserefund.PurchaseRefundLine) []*shop.PurchaseRefundLine {
	var result []*shop.PurchaseRefundLine
	for _, v := range args {
		result = append(result, &shop.PurchaseRefundLine{
			VariantID:    v.VariantID,
			ProductID:    v.ProductID,
			Quantity:     v.Quantity,
			Code:         v.Code,
			ImageURL:     v.ImageURL,
			Name:         v.ProductName,
			PaymentPrice: v.PaymentPrice,
			Attributes:   v.Attributes,
			Adjustment:   v.Adjustment,
		})
	}
	return result
}
