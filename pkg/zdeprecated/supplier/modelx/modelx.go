package modelx

import (
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

type VariantExternalWithQuantity struct {
	Variant *catalogmodel.VariantExternal

	QuantityOnHand   int
	QuantityReserved int
}

type GetSupplierQuery struct {
	SupplierID int64

	Result *model.Supplier
}

type GetSupplierExtendedQuery struct {
	SupplierID         int64
	KiotvietRetailerID string

	Result *model.SupplierExtended
}

type GetSuppliersQuery struct {
	Paging *cm.Paging
	IDs    []int64
	Status *model.Status3

	Result struct {
		Suppliers []*model.Supplier
		Total     int
	}
}

type GetSupplierExtendedsQuery struct {
	Paging *cm.Paging
	IDs    []int64
	Status *model.Status3

	Result struct {
		Suppliers []*model.SupplierExtended
		Total     int
	}
}

type GetSuppliersWithShipFromAddressQuery struct {
	IDs []int64

	Result struct {
		Suppliers []*model.SupplierShipFromAddress
	}
}

type GetAllSuppliersQuery struct {
	Result []*model.Supplier
}

// GetSupplierWithPermissionQuery will set HasPermission to false if the user has no permission to access the shop.
type GetSupplierWithPermissionQuery struct {
	SupplierID int64
	UserID     int64

	Result struct {
		Supplier   *model.Supplier
		Permission model.Permission
	}
}

type CreateSupplierKiotvietCommand struct {
	OwnerID int64
	model.SupplierInfo
	Kiotviet        model.SupplierKiotviet
	DefaultBranchID string
	IsTest          bool
	URLSlug         string

	Result struct {
		Supplier              *model.Supplier
		ProductSource         *model.ProductSource
		ProductSourceInternal *model.ProductSourceInternal
	}
}

type UpdateKiotvietAccessTokenCommand struct {
	ProductSourceID int64
	ClientToken     string
	ExpiresAt       time.Time
}

type ExternalQuantity struct {
	ExternalProductID string
	BranchID          string
	QuantityOnHand    int
	QuantityReserved  int
}

type UpdatePriceRulesCommand struct {
	SupplierID int64
	PriceRules *model.SupplierPriceRules
}

type UpdateSupplierCommand struct {
	Supplier *model.Supplier

	Result *model.SupplierExtended
}

type SetDefaultAddressSupplierCommand struct {
	SupplierID int64
	Type       string
	AddressID  int64

	Result struct {
		Updated int
	}
}
