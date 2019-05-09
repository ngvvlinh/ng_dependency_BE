package model

import (
	"time"

	cm "etop.vn/backend/pkg/common"
)

type GetSupplierQuery struct {
	SupplierID int64

	Result *Supplier
}

type GetSupplierExtendedQuery struct {
	SupplierID         int64
	KiotvietRetailerID string

	Result *SupplierExtended
}

type GetSuppliersQuery struct {
	Paging *cm.Paging
	IDs    []int64
	Status *Status3

	Result struct {
		Suppliers []*Supplier
		Total     int
	}
}

type GetSupplierExtendedsQuery struct {
	Paging *cm.Paging
	IDs    []int64
	Status *Status3

	Result struct {
		Suppliers []*SupplierExtended
		Total     int
	}
}

type GetSuppliersWithShipFromAddressQuery struct {
	IDs []int64

	Result struct {
		Suppliers []*SupplierShipFromAddress
	}
}

type GetAllSuppliersQuery struct {
	Result []*Supplier
}

// GetSupplierWithPermissionQuery will set HasPermission to false if the user has no permission to access the shop.
type GetSupplierWithPermissionQuery struct {
	SupplierID int64
	UserID     int64

	Result struct {
		Supplier   *Supplier
		Permission Permission
	}
}

type CreateSupplierKiotvietCommand struct {
	OwnerID int64
	SupplierInfo
	Kiotviet        SupplierKiotviet
	DefaultBranchID string
	IsTest          bool
	URLSlug         string

	Result struct {
		Supplier              *Supplier
		ProductSource         *ProductSource
		ProductSourceInternal *ProductSourceInternal
	}
}

type UpdateKiotvietAccessTokenCommand struct {
	ProductSourceID int64
	ClientToken     string
	ExpiresAt       time.Time
}

type VariantExternalWithQuantity struct {
	Variant *VariantExternal

	QuantityOnHand   int
	QuantityReserved int
}

type ExternalQuantity struct {
	ExternalProductID string
	BranchID          string
	QuantityOnHand    int
	QuantityReserved  int
}

type UpdatePriceRulesCommand struct {
	SupplierID int64
	PriceRules *SupplierPriceRules
}

type CompanyInfo struct {
	Name                string         `json:"name"`
	TaxCode             string         `json:"tax_code"`
	Address             string         `json:"address"`
	Website             string         `json:"website"`
	LegalRepresentative *ContactPerson `json:"legal_representative"`
}

type ContactPerson struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type UpdateSupplierCommand struct {
	Supplier *Supplier

	Result *SupplierExtended
}

type SetDefaultAddressSupplierCommand struct {
	SupplierID int64
	Type       string
	AddressID  int64

	Result struct {
		Updated int
	}
}
