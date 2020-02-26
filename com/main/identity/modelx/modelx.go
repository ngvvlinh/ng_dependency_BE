package modelx

import (
	"etop.vn/api/top/types/etc/account_type"
	identitymodel "etop.vn/backend/com/main/identity/model"
	"etop.vn/capi/dot"
)

type GetByIDArgs struct {
	ID dot.ID
}

type UpdateAccountURLSlugCommand struct {
	AccountID dot.ID
	URLSlug   string
}

type GetAccountAuthQuery struct {
	AuthKey     string
	AccountType account_type.AccountType
	AccountID   dot.ID

	Result struct {
		AccountAuth *identitymodel.AccountAuth
		Account     identitymodel.AccountInterface
	}
}

type GetPartner struct {
	PartnerID dot.ID

	Result struct {
		Partner *identitymodel.Partner
	}
}

type GetPartnerRelationQuery struct {
	PartnerID         dot.ID
	AccountID         dot.ID
	ExternalAccountID string
	ExternalUserID    string
	AuthKey           string

	Result struct {
		identitymodel.PartnerRelationFtShop
	}
}

type GetPartnerRelationsQuery struct {
	PartnerID dot.ID
	OwnerID   dot.ID

	Result struct {
		Relations []*identitymodel.PartnerRelationFtShop
	}
}

type GetPartnersFromRelationQuery struct {
	AccountIDs []dot.ID

	Result struct {
		Partners []*identitymodel.Partner
	}
}

type CreatePartnerCommand struct {
	Partner *identitymodel.Partner

	Result struct {
		Partner *identitymodel.Partner
	}
}

type CreatePartnerRelationCommand struct {
	AccountID  dot.ID
	UserID     dot.ID
	PartnerID  dot.ID
	ExternalID string

	Result struct {
		*identitymodel.PartnerRelation
	}
}

type UpdatePartnerRelationCommand struct {
	AccountID  dot.ID
	PartnerID  dot.ID
	ExternalID string
}

type GetPartnersQuery struct {
	AvailableFromEtop bool

	Result struct {
		Partners []*identitymodel.Partner
	}
}
