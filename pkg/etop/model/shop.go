package model

import (
	cm "etop.vn/backend/pkg/common"
)

type CreateShopCommand struct {
	Name                        string
	OwnerID                     int64
	AddressID                   int64
	Address                     *Address
	Phone                       string
	BankAccount                 *BankAccount
	WebsiteURL                  string
	ImageURL                    string
	Email                       string
	AutoCreateFFM               bool
	IsTest                      bool
	URLSlug                     string
	CompanyInfo                 *CompanyInfo
	MoneyTransactionRRule       string
	SurveyInfo                  []*SurveyInfo
	ShippingServicePickStrategy []*ShippingServiceSelectStrategyItem

	Result *ShopExtended
}

type UpdateShopCommand struct {
	Shop   *Shop
	Result *ShopExtended
}

type DeleteShopCommand struct {
	ID      int64
	OwnerID int64
}

type SetDefaultAddressShopCommand struct {
	ShopID    int64
	Type      string
	AddressID int64

	Result struct {
		Updated int
	}
}

type GetShopQuery struct {
	ShopID int64

	Result *Shop
}

type GetShopsQuery struct {
	ShopIDs []int64

	Result struct {
		Shops []*Shop
	}
}

type GetShopExtendedQuery struct {
	ShopID int64

	IncludeDeleted bool

	Result *ShopExtended
}

type GetAllShopsQuery struct {
	Result []*Shop
}

type GetAllShopExtendedsQuery struct {
	Paging *cm.Paging

	Result struct {
		Shops []*ShopExtended
		Total int
	}
}

// GetShopWithPermissionQuery will set HasPermission to false if the user has no permission to access the shop
type GetShopWithPermissionQuery struct {
	ShopID int64
	UserID int64

	Result struct {
		Shop       *Shop
		Permission Permission
	}
}
