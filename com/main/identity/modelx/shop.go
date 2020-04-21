package modelx

import (
	adressmodel "o.o/backend/com/main/address/model"
	identitymodel "o.o/backend/com/main/identity/model"
	identitysharemodel "o.o/backend/com/main/identity/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

type CreateShopCommand struct {
	Name                        string
	OwnerID                     dot.ID
	AddressID                   dot.ID
	Address                     *adressmodel.Address
	Phone                       string
	BankAccount                 *identitysharemodel.BankAccount
	WebsiteURL                  dot.NullString
	ImageURL                    string
	Email                       string
	AutoCreateFFM               bool
	IsTest                      bool
	URLSlug                     string
	CompanyInfo                 *identitysharemodel.CompanyInfo
	MoneyTransactionRRule       string
	SurveyInfo                  []*identitymodel.SurveyInfo
	ShippingServicePickStrategy []*identitymodel.ShippingServiceSelectStrategyItem

	Result *identitymodel.ShopExtended
}

type UpdateShopCommand struct {
	Shop          *identitymodel.Shop
	AutoCreateFFM dot.NullBool
	Result        *identitymodel.ShopExtended
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

	Result *identitymodel.Shop
}

type GetShopsQuery struct {
	ShopIDs []dot.ID

	Result struct {
		Shops []*identitymodel.Shop
	}
}

type GetShopExtendedQuery struct {
	ShopID dot.ID

	IncludeDeleted bool

	Result *identitymodel.ShopExtended
}

type GetAllShopExtendedsQuery struct {
	Paging *cm.Paging

	Result struct {
		Shops []*identitymodel.ShopExtended
	}
}

// GetShopWithPermissionQuery will set HasPermission to false if the user has no permission to access the shop
type GetShopWithPermissionQuery struct {
	ShopID dot.ID
	UserID dot.ID

	Result struct {
		Shop       *identitymodel.Shop
		Permission identitymodel.Permission
	}
}
