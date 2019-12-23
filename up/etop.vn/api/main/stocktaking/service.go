package stocktaking

import (
	"context"

	"etop.vn/api/main/inventory"
	"etop.vn/api/meta"
	"etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateStocktake(context.Context, *CreateStocktakeRequest) (*ShopStocktake, error)

	UpdateStocktake(context.Context, *UpdateStocktakeRequest) (*ShopStocktake, error)

	ConfirmStocktake(context.Context, *ConfirmStocktakeRequest) (*ShopStocktake, error)

	CancelStocktake(context.Context, *CancelStocktakeRequest) (*ShopStocktake, error)
}

type QueryService interface {
	GetStocktakeByID(ctx context.Context, id dot.ID, shopID dot.ID) (*ShopStocktake, error)

	GetStocktakesByIDs(ctx context.Context, ids []dot.ID, shopID dot.ID) ([]*ShopStocktake, error)

	ListStocktake(context.Context, *ListStocktakeRequest) (*ListStocktakeResponse, error)
}

// +convert:create=ShopStocktake
type CreateStocktakeRequest struct {
	ShopID        dot.ID
	TotalQuantity int
	CreatedBy     dot.ID
	Lines         []*StocktakeLine
	Note          string
}

// +convert:update=ShopStocktake
type UpdateStocktakeRequest struct {
	ShopID        dot.ID
	ID            dot.ID
	TotalQuantity int
	UpdatedBy     dot.ID
	Lines         []*StocktakeLine
	Note          string
}

type CancelStocktakeRequest struct {
	ShopID       dot.ID
	ID           dot.ID
	CancelReason string
}

type ConfirmStocktakeRequest struct {
	ID                   dot.ID
	ShopID               dot.ID
	ConfirmedBy          dot.ID
	OverStock            bool
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type ListStocktakeRequest struct {
	Page   meta.Paging
	ShopID dot.ID
	Filter []meta.Filter
}

type ListStocktakeResponse struct {
	Stocktakes []*ShopStocktake
	PageInfo   meta.PageInfo
}
