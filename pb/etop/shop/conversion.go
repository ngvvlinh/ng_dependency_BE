package shop

import (
	haravanidentity "etop.vn/api/external/haravan/identity"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	pbcm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/etop"
)

func Convert_core_XAccountAhamove_To_api_XAccountAhamove(in *identity.ExternalAccountAhamove) *ExternalAccountAhamove {
	if in == nil {
		return nil
	}
	return &ExternalAccountAhamove{
		Id:                  in.ID,
		Phone:               in.Phone,
		Name:                in.Name,
		ExternalVerified:    in.ExternalVerified,
		ExternalCreatedAt:   pbcm.PbTime(in.ExternalCreatedAt),
		CreatedAt:           pbcm.PbTime(in.CreatedAt),
		UpdatedAt:           pbcm.PbTime(in.UpdatedAt),
		LastSendVerifyAt:    pbcm.PbTime(in.LastSendVerifiedAt),
		ExternalTicketId:    in.ExternalTicketID,
		IdCardFrontImg:      in.IDCardFrontImg,
		IdCardBackImg:       in.IDCardBackImg,
		PortraitImg:         in.PortraitImg,
		UploadedAt:          pbcm.PbTime(in.UploadedAt),
		FanpageUrl:          in.FanpageURL,
		WebsiteUrl:          in.WebsiteURL,
		CompanyImgs:         in.CompanyImgs,
		BusinessLicenseImgs: in.BusinessLicenseImgs,
	}
}

func Convert_core_XAccountHaravan_To_api_XAccountHaravan(in *haravanidentity.ExternalAccountHaravan) *ExternalAccountHaravan {
	if in == nil {
		return nil
	}
	return &ExternalAccountHaravan{
		Id:                                in.ID,
		ShopId:                            in.ShopID,
		Subdomain:                         in.Subdomain,
		ExternalCarrierServiceId:          int32(in.ExternalCarrierServiceID),
		ExternalConnectedCarrierServiceAt: pbcm.PbTime(in.ExternalConnectedCarrierServiceAt),
		ExpiresAt:                         pbcm.PbTime(in.ExpiresAt),
		CreatedAt:                         pbcm.PbTime(in.CreatedAt),
		UpdatedAt:                         pbcm.PbTime(in.UpdatedAt),
	}
}

func Convert_api_ReceiptLine_To_core_ReceiptLine(in *ReceiptLine) *receipting.ReceiptLine {
	if in == nil {
		return nil
	}
	return &receipting.ReceiptLine{
		RefID:  in.RefId,
		Title:  in.Title,
		Amount: in.Amount,
	}
}

func Convert_api_ReceiptLines_To_core_ReceiptLines(in []*ReceiptLine) []*receipting.ReceiptLine {
	out := make([]*receipting.ReceiptLine, len(in))
	for i := range in {
		out[i] = Convert_api_ReceiptLine_To_core_ReceiptLine(in[i])
	}

	return out
}

func Convert_api_BankAccount_To_core_BankAccount(in *etop.BankAccount) *identity.BankAccount {
	if in == nil {
		return nil
	}
	return &identity.BankAccount{
		Name:          in.Name,
		Province:      in.Province,
		Branch:        in.Branch,
		AccountNumber: in.AccountNumber,
		AccountName:   in.AccountName,
	}
}

func Convert_api_PurchaseOrderLine_To_core_PurchaseOrderLine(in *PurchaseOrderLine) *purchaseorder.PurchaseOrderLine {
	if in == nil {
		return nil
	}
	return &purchaseorder.PurchaseOrderLine{
		VariantID: in.VariantId,
		Quantity:  in.Quantity,
		Price:     in.Price,
	}
}

func Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(in []*PurchaseOrderLine) []*purchaseorder.PurchaseOrderLine {
	out := make([]*purchaseorder.PurchaseOrderLine, len(in))
	for i := range in {
		out[i] = Convert_api_PurchaseOrderLine_To_core_PurchaseOrderLine(in[i])
	}

	return out
}
