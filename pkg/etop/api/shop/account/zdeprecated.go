package account

import (
	"context"

	"o.o/api/main/identity"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/bus"
)

// deprecated
func (s *AccountService) UpdateExternalAccountAhamoveVerificationImages(ctx context.Context, r *api.UpdateXAccountAhamoveVerificationRequest) (*pbcm.UpdatedResponse, error) {
	if err := validateUrl(r.IdCardFrontImg, r.IdCardBackImg, r.PortraitImg, r.WebsiteUrl, r.FanpageUrl); err != nil {
		return nil, err
	}
	if err := validateUrl(r.BusinessLicenseImgs...); err != nil {
		return nil, err
	}
	if err := validateUrl(r.CompanyImgs...); err != nil {
		return nil, err
	}

	query := &identitymodelx.GetUserByIDQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	user := query.Result
	phone := user.Phone

	cmd := &identity.UpdateExternalAccountAhamoveVerificationCommand{
		OwnerID:             user.ID,
		Phone:               phone,
		IDCardFrontImg:      r.IdCardFrontImg,
		IDCardBackImg:       r.IdCardBackImg,
		PortraitImg:         r.PortraitImg,
		WebsiteURL:          r.WebsiteUrl,
		FanpageURL:          r.FanpageUrl,
		CompanyImgs:         r.CompanyImgs,
		BusinessLicenseImgs: r.BusinessLicenseImgs,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}

	result := &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return result, nil
}