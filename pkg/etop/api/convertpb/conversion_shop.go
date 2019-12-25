package convertpb

import (
	haravanidentity "etop.vn/api/external/haravan/identity"
	"etop.vn/api/main/identity"
	"etop.vn/api/main/purchaseorder"
	"etop.vn/api/main/receipting"
	etop "etop.vn/api/top/int/etop"
	shop "etop.vn/api/top/int/shop"
	"etop.vn/backend/pkg/common/apifw/cmapi"
)

func Convert_core_XAccountAhamove_To_api_XAccountAhamove(in *identity.ExternalAccountAhamove, hideInfo bool) *shop.ExternalAccountAhamove {
	if in == nil {
		return nil
	}
	if hideInfo {
		return &shop.ExternalAccountAhamove{
			Id: in.ID,
		}
	}
	return &shop.ExternalAccountAhamove{
		Id:                  in.ID,
		Phone:               in.Phone,
		Name:                in.Name,
		ExternalVerified:    in.ExternalVerified,
		ExternalCreatedAt:   cmapi.PbTime(in.ExternalCreatedAt),
		CreatedAt:           cmapi.PbTime(in.CreatedAt),
		UpdatedAt:           cmapi.PbTime(in.UpdatedAt),
		LastSendVerifyAt:    cmapi.PbTime(in.LastSendVerifiedAt),
		ExternalTicketId:    in.ExternalTicketID,
		IdCardFrontImg:      in.IDCardFrontImg,
		IdCardBackImg:       in.IDCardBackImg,
		PortraitImg:         in.PortraitImg,
		UploadedAt:          cmapi.PbTime(in.UploadedAt),
		FanpageUrl:          in.FanpageURL,
		WebsiteUrl:          in.WebsiteURL,
		CompanyImgs:         in.CompanyImgs,
		BusinessLicenseImgs: in.BusinessLicenseImgs,
	}
}

func Convert_core_XAccountHaravan_To_api_XAccountHaravan(in *haravanidentity.ExternalAccountHaravan, hide bool) *shop.ExternalAccountHaravan {
	if in == nil {
		return nil
	}
	if hide {
		return &shop.ExternalAccountHaravan{
			Id: in.ID,
		}
	}
	return &shop.ExternalAccountHaravan{
		Id:                                in.ID,
		ShopId:                            in.ShopID,
		Subdomain:                         in.Subdomain,
		ExternalCarrierServiceId:          in.ExternalCarrierServiceID,
		ExternalConnectedCarrierServiceAt: cmapi.PbTime(in.ExternalConnectedCarrierServiceAt),
		ExpiresAt:                         cmapi.PbTime(in.ExpiresAt),
		CreatedAt:                         cmapi.PbTime(in.CreatedAt),
		UpdatedAt:                         cmapi.PbTime(in.UpdatedAt),
	}
}

func Convert_api_ReceiptLine_To_core_ReceiptLine(in *shop.ReceiptLine) *receipting.ReceiptLine {
	if in == nil {
		return nil
	}
	return &receipting.ReceiptLine{
		RefID:  in.RefId,
		Title:  in.Title,
		Amount: in.Amount,
	}
}

func Convert_api_ReceiptLines_To_core_ReceiptLines(in []*shop.ReceiptLine) []*receipting.ReceiptLine {
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
	}
}

func Convert_api_PurchaseOrderLines_To_core_PurchaseOrderLines(in []*shop.PurchaseOrderLine) []*purchaseorder.PurchaseOrderLine {
	out := make([]*purchaseorder.PurchaseOrderLine, len(in))
	for i := range in {
		out[i] = Convert_api_PurchaseOrderLine_To_core_PurchaseOrderLine(in[i])
	}

	return out
}
