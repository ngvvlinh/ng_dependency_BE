package affiliate

import (
	"context"

	cm "etop.vn/api/pb/common"
	"etop.vn/api/pb/etop"
	aff "etop.vn/api/pb/etop/affiliate"
)

// +gen:apix
// +gen:apix:doc-path=etop/affiliate

// +apix:path=/affiliate.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/affiliate.Account
type AccountService interface {
	RegisterAffiliate(context.Context, *aff.RegisterAffiliateRequest) (*etop.Affiliate, error)
	UpdateAffiliate(context.Context, *aff.UpdateAffiliateRequest) (*etop.Affiliate, error)
	UpdateAffiliateBankAccount(context.Context, *aff.UpdateAffiliateBankAccountRequest) (*etop.Affiliate, error)
	DeleteAffiliate(context.Context, *cm.IDRequest) (*cm.Empty, error)
}
