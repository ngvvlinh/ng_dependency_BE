package idutil

import (
	"o.o/api/top/types/etc/account_tag"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func NewID() dot.ID {
	return cm.NewID()
}

func NewShopID() dot.ID {
	return cm.NewIDWithTag(account_tag.TagShop)
}

func NewAffiliateID() dot.ID {
	return cm.NewIDWithTag(account_tag.TagAffiliate)
}

func IsPartnerID(id dot.ID) bool {
	return cm.GetTag(id) == account_tag.TagPartner
}

func IsShopID(id dot.ID) bool {
	return cm.GetTag(id) == account_tag.TagShop
}

func IsAffiliateID(id dot.ID) bool {
	return cm.GetTag(id) == account_tag.TagAffiliate
}

func IsEtopAccountID(id dot.ID) bool {
	return id == EtopAccountID
}

const EtopAccountID = account_tag.TagEtop
const EtopTradingAccountID = 1000015764575267699
const TopShipID = 1000030662086749358
