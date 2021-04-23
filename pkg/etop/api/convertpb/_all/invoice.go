package _all

import (
	"o.o/api/main/invoicing"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbInvoice(in *invoicing.InvoiceFtLine) *types.Invoice {
	if in == nil {
		return nil
	}
	return &types.Invoice{
		ID:            in.ID,
		AccountID:     in.AccountID,
		TotalAmount:   in.TotalAmount,
		Description:   in.Description,
		PaymentID:     in.PaymentID,
		Status:        in.Status,
		PaymentStatus: in.PaymentStatus,
		Customer:      PbSubscriptionCustomer(in.Customer),
		Lines:         PbInvoiceLines(in.Lines),
		CreatedAt:     cmapi.PbTime(in.CreatedAt),
		UpdatedAt:     cmapi.PbTime(in.UpdatedAt),
		ReferralType:  in.ReferralType,
		ReferralIDs:   in.ReferralIDs,
		Classify:      in.Classify,
		Type:          in.Type,
	}
}

func PbInvoices(items []*invoicing.InvoiceFtLine) []*types.Invoice {
	result := make([]*types.Invoice, len(items))
	for i, item := range items {
		result[i] = PbInvoice(item)
	}
	return result
}

func PbInvoiceLine(in *invoicing.InvoiceLine) *types.InvoiceLine {
	if in == nil {
		return nil
	}
	return &types.InvoiceLine{
		ID:           in.ID,
		LineAmount:   in.LineAmount,
		Price:        in.Price,
		Quantity:     in.Quantity,
		Description:  in.Description,
		InvoiceID:    in.InvoiceID,
		ReferralType: in.ReferralType,
		ReferralID:   in.ReferralID,
		CreatedAt:    cmapi.PbTime(in.CreatedAt),
		UpdatedAt:    cmapi.PbTime(in.UpdatedAt),
	}
}

func PbInvoiceLines(items []*invoicing.InvoiceLine) []*types.InvoiceLine {
	result := make([]*types.InvoiceLine, len(items))
	for i, item := range items {
		result[i] = PbInvoiceLine(item)
	}
	return result
}

func Convert_api_InvoiceLine_To_core_InvoiceLine(in *types.InvoiceLine) *invoicing.InvoiceLine {
	if in == nil {
		return nil
	}
	return &invoicing.InvoiceLine{
		ID:           in.ID,
		LineAmount:   in.LineAmount,
		Price:        in.Price,
		Quantity:     in.Quantity,
		Description:  in.Description,
		InvoiceID:    in.InvoiceID,
		ReferralType: in.ReferralType,
		ReferralID:   in.ReferralID,
		CreatedAt:    in.CreatedAt.ToTime(),
		UpdatedAt:    in.UpdatedAt.ToTime(),
	}
}

func Convert_api_InvoiceLines_To_core_InvoiceLines(items []*types.InvoiceLine) []*invoicing.InvoiceLine {
	result := make([]*invoicing.InvoiceLine, len(items))
	for i, item := range items {
		result[i] = Convert_api_InvoiceLine_To_core_InvoiceLine(item)
	}
	return result
}
