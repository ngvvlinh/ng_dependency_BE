package _all

import (
	"o.o/api/subscripting/invoice"
	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/subscripting/subscriptionproduct"
	subscriptingtypes "o.o/api/subscripting/types"
	"o.o/api/top/int/types"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbSubrProduct(in *subscriptionproduct.SubscriptionProduct) *types.SubscriptionProduct {
	if in == nil {
		return nil
	}
	return &types.SubscriptionProduct{
		ID:          in.ID,
		Name:        in.Name,
		Type:        in.Type,
		Description: in.Description,
		ImageURL:    in.ImageURL,
		Status:      in.Status,
		CreatedAt:   cmapi.PbTime(in.CreatedAt),
		UpdatedAt:   cmapi.PbTime(in.UpdatedAt),
	}
}

func PbSubrProducts(items []*subscriptionproduct.SubscriptionProduct) []*types.SubscriptionProduct {
	result := make([]*types.SubscriptionProduct, len(items))
	for i, item := range items {
		result[i] = PbSubrProduct(item)
	}
	return result
}

func PbSubrPlan(in *subscriptionplan.SubscriptionPlan) *types.SubscriptionPlan {
	if in == nil {
		return nil
	}
	return &types.SubscriptionPlan{
		ID:            in.ID,
		Name:          in.Name,
		Price:         in.Price,
		Status:        in.Status,
		Description:   in.Description,
		ProductID:     in.ProductID,
		Interval:      in.Interval,
		IntervalCount: in.IntervalCount,
		CreatedAt:     cmapi.PbTime(in.CreatedAt),
		UpdatedAt:     cmapi.PbTime(in.UpdatedAt),
	}
}

func PbSubrPlans(items []*subscriptionplan.SubscriptionPlan) []*types.SubscriptionPlan {
	result := make([]*types.SubscriptionPlan, len(items))
	for i, item := range items {
		result[i] = PbSubrPlan(item)
	}
	return result
}

func PbSubscription(in *subscription.SubscriptionFtLine) *types.Subscription {
	if in == nil {
		return nil
	}
	return &types.Subscription{
		ID:                   in.ID,
		AccountID:            in.AccountID,
		CancelAtPeriodEnd:    in.CancelAtPeriodEnd,
		CurrentPeriodStartAt: cmapi.PbTime(in.CurrentPeriodStartAt),
		CurrentPeriodEndAt:   cmapi.PbTime(in.CurrentPeriodEndAt),
		Status:               in.Status,
		BillingCycleAnchorAt: cmapi.PbTime(in.BillingCycleAnchorAt),
		StartedAt:            cmapi.PbTime(in.StartedAt),
		Lines:                PbSubsriptionLines(in.Lines),
		Customer:             PbSubscriptionCustomer(in.Customer),
		CreatedAt:            cmapi.PbTime(in.CreatedAt),
		UpdatedAt:            cmapi.PbTime(in.UpdatedAt),
	}
}

func PbSubscriptionCustomer(in *subscriptingtypes.CustomerInfo) *types.SubrCustomer {
	if in == nil {
		return nil
	}
	return &types.SubrCustomer{
		FullName: in.FullName,
		Email:    in.Email,
		Phone:    in.Phone,
	}
}

func Convert_api_SubrCustomer_To_core_SubrCustomer(in *types.SubrCustomer) *subscriptingtypes.CustomerInfo {
	if in == nil {
		return nil
	}
	return &subscriptingtypes.CustomerInfo{
		FullName: in.FullName,
		Email:    in.Email,
		Phone:    in.Phone,
	}
}

func PbSubscriptionLine(in *subscription.SubscriptionLine) *types.SubscriptionLine {
	if in == nil {
		return nil
	}
	return &types.SubscriptionLine{
		ID:             in.ID,
		PlanID:         in.PlanID,
		SubscriptionID: in.SubscriptionID,
		Quantity:       in.Quantity,
		CreatedAt:      cmapi.PbTime(in.CreatedAt),
		UpdatedAt:      cmapi.PbTime(in.UpdatedAt),
	}
}

func PbSubsriptionLines(items []*subscription.SubscriptionLine) []*types.SubscriptionLine {
	result := make([]*types.SubscriptionLine, len(items))
	for i, item := range items {
		result[i] = PbSubscriptionLine(item)
	}
	return result
}

func PbSubscriptions(items []*subscription.SubscriptionFtLine) []*types.Subscription {
	result := make([]*types.Subscription, len(items))
	for i, item := range items {
		result[i] = PbSubscription(item)
	}
	return result
}

func Convert_api_SubscriptionLine_To_core_SubscriptionLine(in *types.SubscriptionLine) *subscription.SubscriptionLine {
	if in == nil {
		return nil
	}
	return &subscription.SubscriptionLine{
		ID:             in.ID,
		PlanID:         in.PlanID,
		SubscriptionID: in.SubscriptionID,
		Quantity:       in.Quantity,
		CreatedAt:      in.CreatedAt.ToTime(),
		UpdatedAt:      in.UpdatedAt.ToTime(),
	}
}

func Convert_api_SubscriptionLines_To_core_SubscriptionLines(items []*types.SubscriptionLine) []*subscription.SubscriptionLine {
	result := make([]*subscription.SubscriptionLine, len(items))
	for i, item := range items {
		result[i] = Convert_api_SubscriptionLine_To_core_SubscriptionLine(item)
	}
	return result
}

func PbInvoice(in *invoice.InvoiceFtLine) *types.Invoice {
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
	}
}

func PbInvoices(items []*invoice.InvoiceFtLine) []*types.Invoice {
	result := make([]*types.Invoice, len(items))
	for i, item := range items {
		result[i] = PbInvoice(item)
	}
	return result
}

func PbInvoiceLine(in *invoice.InvoiceLine) *types.InvoiceLine {
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

func PbInvoiceLines(items []*invoice.InvoiceLine) []*types.InvoiceLine {
	result := make([]*types.InvoiceLine, len(items))
	for i, item := range items {
		result[i] = PbInvoiceLine(item)
	}
	return result
}

func Convert_api_InvoiceLine_To_core_InvoiceLine(in *types.InvoiceLine) *invoice.InvoiceLine {
	if in == nil {
		return nil
	}
	return &invoice.InvoiceLine{
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

func Convert_api_InvoiceLines_To_core_InvoiceLines(items []*types.InvoiceLine) []*invoice.InvoiceLine {
	result := make([]*invoice.InvoiceLine, len(items))
	for i, item := range items {
		result[i] = Convert_api_InvoiceLine_To_core_InvoiceLine(item)
	}
	return result
}
