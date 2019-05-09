package model

type GetPartner struct {
	PartnerID int64

	Result struct {
		Partner *Partner
	}
}

type GetPartnerRelationQuery struct {
	PartnerID         int64
	AccountID         int64
	ExternalAccountID string
	AuthKey           string

	Result struct {
		PartnerRelationFtShop
	}
}

type GetPartnerRelationsQuery struct {
	PartnerID int64
	OwnerID   int64

	Result struct {
		Relations []*PartnerRelationFtShop
	}
}

type GetPartnersFromRelationQuery struct {
	AccountIDs []int64

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
	AccountID  int64
	PartnerID  int64
	ExternalID string

	Result struct {
		*PartnerRelation
	}
}

type UpdatePartnerRelationCommand struct {
	AccountID  int64
	PartnerID  int64
	ExternalID string
}

type GetPartnersQuery struct {
	AvailableFromEtop bool

	Result struct {
		Partners []*Partner
	}
}
