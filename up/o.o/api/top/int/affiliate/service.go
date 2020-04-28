package affiliate

import (
	"context"

	etop "o.o/api/top/int/etop"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:swagger:doc-path=etop/affiliate

// +apix:path=/affiliate.Misc
type MiscService interface {
	VersionInfo(context.Context, *cm.Empty) (*cm.VersionInfoResponse, error)
}

// +apix:path=/affiliate.Account
type AccountService interface {
	RegisterAffiliate(context.Context, *RegisterAffiliateRequest) (*etop.Affiliate, error)
	UpdateAffiliate(context.Context, *UpdateAffiliateRequest) (*etop.Affiliate, error)
	UpdateAffiliateBankAccount(context.Context, *UpdateAffiliateBankAccountRequest) (*etop.Affiliate, error)
	DeleteAffiliate(context.Context, *cm.IDRequest) (*cm.Empty, error)
}
