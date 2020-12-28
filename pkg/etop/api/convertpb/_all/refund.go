package _all

import (
	"o.o/api/main/refund"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbRefunds(args []*refund.Refund) []*shop.Refund {
	var result []*shop.Refund
	for _, value := range args {
		result = append(result, PbRefund(value))
	}
	return result
}

func PbRefund(args *refund.Refund) *shop.Refund {
	var result = &shop.Refund{
		ID:              args.ID,
		ShopID:          args.ShopID,
		OrderID:         args.OrderID,
		Note:            args.Note,
		Code:            args.Code,
		AdjustmentLines: args.AdjustmentLines,
		TotalAdjustment: args.TotalAdjustment,
		Lines:           PbRefundLine(args.Lines),
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

func PbRefundLine(args []*refund.RefundLine) []*shop.RefundLine {
	var result []*shop.RefundLine
	for _, v := range args {
		result = append(result, &shop.RefundLine{
			VariantID:   v.VariantID,
			ProductID:   v.ProductID,
			Quantity:    v.Quantity,
			Code:        v.Code,
			ImageURL:    v.ImageURL,
			Name:        v.ProductName,
			RetailPrice: v.RetailPrice,
			Attributes:  v.Attributes,
			Adjustment:  v.Adjustment,
		})
	}
	return result
}
