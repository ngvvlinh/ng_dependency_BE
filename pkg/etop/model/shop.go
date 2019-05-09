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

type GetShopCollectionQuery struct {
	ShopID       int64
	CollectionID int64

	Result *ShopCollection
}

type GetShopCollectionsQuery struct {
	ShopID        int64
	CollectionIDs []int64

	Result struct {
		Collections []*ShopCollection
	}
}

type CreateShopCollectionCommand struct {
	Collection *ShopCollection

	Result *ShopCollection
}

type UpdateShopCollectionCommand struct {
	Collection *ShopCollection

	Result *ShopCollection
}

type RemoveShopCollectionCommand struct {
	ShopID       int64
	CollectionID int64

	Result struct {
		Deleted int
	}
}

type CreateProductSourceCommand struct {
	ShopID int64
	Name   string
	Type   string

	Result *ProductSource
}

type GetShopProductSourcesCommand struct {
	UserID int64
	ShopID int64

	Result []*ProductSource
}

type CreateVariantCommand struct {
	ShopID          int64
	ProductSourceID int64
	ProductID       int64
	// In `Dép Adidas Adilette Slides - Full Đỏ`, product_name is "Dép Adidas Adilette Slides"
	ProductName string
	// In `Dép Adidas Adilette Slides - Full Đỏ`, name is "Full Đỏ"
	Name              string
	Description       string
	ShortDesc         string
	ImageURLs         []string
	Tags              []string
	Status            Status3
	ListPrice         int
	SKU               string
	Code              string
	QuantityAvailable int
	QuantityOnHand    int
	QuantityReserved  int
	CostPrice         int
	Attributes        []ProductAttribute
	DescHTML          string

	Result *ShopProductFtVariant
}

type ConnectProductSourceCommand struct {
	ProductSourceID int64
	ShopID          int64

	Result struct {
		Updated int
	}
}

type RemoveProductSourceCommand struct {
	ShopID int64

	Result struct {
		Updated int
	}
}

type CreateProductSourceCategoryCommand struct {
	ShopID            int64
	Name              string
	ProductSourceID   int64
	ProductSourceType string
	ParentID          int64

	Result *ProductSourceCategory
}

type UpdateProductsProductSourceCategoryCommand struct {
	CategoryID      int64
	ProductIDs      []int64
	ShopID          int64
	ProductSourceID int64

	Result struct {
		Updated int
	}
}
