package model

import (
	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

type CreateShopCommand struct {
	Name                        string
	OwnerID                     dot.ID
	AddressID                   dot.ID
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
	Shop          *Shop
	AutoCreateFFM *bool
	Result        *ShopExtended
}

type DeleteShopCommand struct {
	ID      dot.ID
	OwnerID dot.ID
}

type SetDefaultAddressShopCommand struct {
	ShopID    dot.ID
	Type      string
	AddressID dot.ID

	Result struct {
		Updated int
	}
}

type GetShopQuery struct {
	ShopID dot.ID

	Result *Shop
}

type GetShopsQuery struct {
	ShopIDs []dot.ID

	Result struct {
		Shops []*Shop
	}
}

type GetShopExtendedQuery struct {
	ShopID dot.ID

	IncludeDeleted bool

	Result *ShopExtended
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
	ShopID dot.ID
	UserID dot.ID

	Result struct {
		Shop       *Shop
		Permission Permission
	}
}
