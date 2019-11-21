package stocktaking

import (
	"context"

	"etop.vn/api/main/inventory"

	"etop.vn/api/meta"
)

// +gen:api

type Aggregate interface {
	CreateStocktake(context.Context, *CreateStocktakeRequest) (*ShopStocktake, error)

	UpdateStocktake(context.Context, *UpdateStocktakeRequest) (*ShopStocktake, error)

	ConfirmStocktake(context.Context, *ConfirmStocktakeRequest) (*ShopStocktake, error)

	CancelStocktake(context.Context, *CancelStocktakeRequest) (*ShopStocktake, error)
}

type QueryService interface {
	GetStocktakeByID(ctx context.Context, id int64, shopID int64) (*ShopStocktake, error)

	GetStocktakesByIDs(ctx context.Context, ids []int64, shopID int64) ([]*ShopStocktake, error)

	ListStocktake(context.Context, *ListStocktakeRequest) (*ListStocktakeResponse, error)
}

// +convert:create=ShopStocktake
type CreateStocktakeRequest struct {
	ShopID        int64
	TotalQuantity int32
	CreatedBy     int64
	Lines         []*StocktakeLine
	Note          string
}

// +convert:update=ShopStocktake
type UpdateStocktakeRequest struct {
	ShopID        int64
	ID            int64
	TotalQuantity int32
	UpdatedBy     int64
	Lines         []*StocktakeLine
	Note          string
}

type CancelStocktakeRequest struct {
	ShopID       int64
	ID           int64
	CancelReason string
}

type ConfirmStocktakeRequest struct {
	ID                   int64
	ShopID               int64
	ConfirmedBy          int64
	OverStock            bool
	AutoInventoryVoucher inventory.AutoInventoryVoucher
}

type ListStocktakeRequest struct {
	Page   meta.Paging
	ShopID int64
	Filter []meta.Filter
}

type ListStocktakeResponse struct {
	Stocktakes []*ShopStocktake
	PageInfo   meta.PageInfo
	Total      int32
}
