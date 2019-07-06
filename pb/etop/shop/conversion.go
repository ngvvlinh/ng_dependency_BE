package shop

import (
	"etop.vn/api/main/identity"
	pbcm "etop.vn/backend/pb/common"
)

func Convert_core_XAhamoveAccount_To_api_XAhamoveAccount(in *identity.ExternalAccountAhamove) *ExternalAccountAhamove {
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
