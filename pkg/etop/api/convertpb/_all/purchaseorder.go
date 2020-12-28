package _all

import (
	"o.o/api/main/purchaseorder"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func Convert_api_PurchaseOrderLine_To_core_PurchaseOrderLine(in *shop.PurchaseOrderLine) *purchaseorder.PurchaseOrderLine {
	if in == nil {
		return nil
	}
	return &purchaseorder.PurchaseOrderLine{
		ImageUrl:     in.ImageUrl,
		ProductName:  in.ProductName,
		VariantID:    in.VariantId,
		Quantity:     in.Quantity,
		PaymentPrice: in.PaymentPrice,
		Discount:     in.Discount,
	}
}

func Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(in []*shop.PurchaseOrderLine) []*purchaseorder.PurchaseOrderLine {
	out := make([]*purchaseorder.PurchaseOrderLine, len(in))
	for i := range in {
		out[i] = Convert_api_PurchaseOrderLine_To_core_PurchaseOrderLine(in[i])
	}

	return out
}

func PbPurchaseOrderLine(m *purchaseorder.PurchaseOrderLine) *shop.PurchaseOrderLine {
	if m == nil {
		return nil
	}
	return &shop.PurchaseOrderLine{
		ProductName:  m.ProductName,
		ImageUrl:     m.ImageUrl,
		ProductId:    m.ProductID,
		Code:         m.Code,
		Attributes:   m.Attributes,
		VariantId:    m.VariantID,
		Quantity:     m.Quantity,
		PaymentPrice: m.PaymentPrice,
	}
}

func PbPurchaseOrderLines(ms []*purchaseorder.PurchaseOrderLine) []*shop.PurchaseOrderLine {
	res := make([]*shop.PurchaseOrderLine, len(ms))
	for i, m := range ms {
		res[i] = PbPurchaseOrderLine(m)
	}
	return res
}

func PbPurchaseOrder(m *purchaseorder.PurchaseOrder) *shop.PurchaseOrder {
	if m == nil {
		return nil
	}
	return &shop.PurchaseOrder{
		Id:            m.ID,
		ShopId:        m.ShopID,
		SupplierId:    m.SupplierID,
		Supplier:      PbPurchaseOrderSupplier(m.Supplier),
		BasketValue:   m.BasketValue,
		DiscountLines: m.DiscountLines,
		TotalDiscount: m.TotalDiscount,
		FeeLines:      m.FeeLines,
		TotalFee:      m.TotalFee,
		TotalAmount:   m.TotalAmount,
		Code:          m.Code,
		Note:          m.Note,
		Status:        m.Status,
		Lines:         PbPurchaseOrderLines(m.Lines),
		PaidAmount:    m.PaidAmount,
		CreatedBy:     m.CreatedBy,
		CancelReason:  m.CancelReason,
		ConfirmedAt:   cmapi.PbTime(m.ConfirmedAt),
		CancelledAt:   cmapi.PbTime(m.CancelledAt),
		CreatedAt:     cmapi.PbTime(m.CreatedAt),
		UpdatedAt:     cmapi.PbTime(m.UpdatedAt),
	}
}

func PbPurchaseOrders(ms []*purchaseorder.PurchaseOrder) []*shop.PurchaseOrder {
	res := make([]*shop.PurchaseOrder, len(ms))
	for i, m := range ms {
		res[i] = PbPurchaseOrder(m)
	}
	return res
}

func PbPurchaseOrderSupplier(m *purchaseorder.PurchaseOrderSupplier) *shop.PurchaseOrderSupplier {
	if m == nil {
		return nil
	}
	return &shop.PurchaseOrderSupplier{
		FullName:           m.FullName,
		Phone:              m.Phone,
		Email:              m.Email,
		CompanyName:        m.CompanyName,
		TaxNumber:          m.TaxNumber,
		HeadquarterAddress: m.HeadquarterAddress,
		Deleted:            m.Deleted,
	}
}
