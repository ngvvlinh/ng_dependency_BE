package list

import (
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/account_type"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etc/idutil"
	"o.o/capi/dot"
)

type PartnerInfo struct {
	OwnerID     dot.ID
	AccountID   dot.ID
	Fullname    string
	Email       string
	Phone       string
	Password    string
	AccountType account_type.AccountType
	AuthToken   string
}

func (p *PartnerInfo) Validate() error {
	if p.Email == "" || p.Phone == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Please provide phone or email")
	}
	if p.Fullname == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing fullname")
	}

	switch p.AccountType {
	case account_type.Shop:
		if p.AccountID == 0 {
			p.AccountID = cm.NewIDWithTag(account_tag.TagShop)
		} else {
			if !idutil.IsShopID(p.AccountID) {
				return cm.Errorf(cm.InvalidArgument, nil, "Account ID does not valid for type shop")
			}
		}
	case account_type.Partner:
		if p.AccountID == 0 {
			p.AccountID = cm.NewIDWithTag(account_tag.TagPartner)
		} else {
			if !idutil.IsPartnerID(p.AccountID) {
				return cm.Errorf(cm.InvalidArgument, nil, "Account ID does not valid for type partner")
			}
		}
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Account type invalid!")
	}

	if p.OwnerID == 0 {
		p.OwnerID = cm.NewID()
	}

	return nil
}
