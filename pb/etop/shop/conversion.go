package shop

import (
	haravanidentity "etop.vn/api/external/haravan/identity"
	"etop.vn/api/main/identity"
	pbcm "etop.vn/backend/pb/common"
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
		Id:        in.ID,
		ShopId:    in.ShopID,
		Subdomain: in.Subdomain,
		ExpiresAt: pbcm.PbTime(in.ExpiresAt),
		CreatedAt: pbcm.PbTime(in.CreatedAt),
		UpdatedAt: pbcm.PbTime(in.UpdatedAt),
	}
}
