package model

import "etop.vn/capi/dot"

type GetPartner struct {
	PartnerID dot.ID

	Result struct {
		Partner *Partner
	}
}

type GetPartnerRelationQuery struct {
	PartnerID         dot.ID
	AccountID         dot.ID
	ExternalAccountID string
	AuthKey           string

	Result struct {
		PartnerRelationFtShop
	}
}

type GetPartnerRelationsQuery struct {
	PartnerID dot.ID
	OwnerID   dot.ID

	Result struct {
		Relations []*PartnerRelationFtShop
	}
}

type GetPartnersFromRelationQuery struct {
	AccountIDs []dot.ID

	Result struct {
		Partners []*Partner
	}
}

type CreatePartnerCommand struct {
	Partner *Partner

	Result struct {
		Partner *Partner
	}
}

type CreatePartnerRelationCommand struct {
	AccountID  dot.ID
	PartnerID  dot.ID
	ExternalID string

	Result struct {
		*PartnerRelation
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
		Partners []*Partner
	}
}
