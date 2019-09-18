package receipting

import (
	"context"

	"etop.vn/api/meta"
	"etop.vn/api/shopping"
	. "etop.vn/capi/dot"
)

// +gen:api

type Aggregate interface {
	CreateReceipt(ctx context.Context, _ *CreateReceiptArgs) (*Receipt, error)
	UpdateReceipt(ctx context.Context, _ *UpdateReceiptArgs) (*Receipt, error)
	DeleteReceipt(ctx context.Context, ID int64, shopID int64) (deleted int, _ error)
}

type QueryService interface {
	GetReceiptByID(context.Context, *shopping.IDQueryShopArg) (*Receipt, error)
	ListReceipts(context.Context, *shopping.ListQueryShopArgs) (*ReceiptsResponse, error)
	ListReceiptsByIDs(context.Context, *shopping.IDsQueryShopArgs) (*ReceiptsResponse, error)
	ListReceiptsByOrderIDs(context.Context, *shopping.IDsQueryShopArgs) (*ReceiptsResponse, error)
}

//-- queries --//

type ReceiptsResponse struct {
	Receipts []*Receipt
	Count    int32
	Paging   meta.PageInfo
}

type OrderIDsQueryArgs struct {
	OrderIDs []int64
	ShopID   int64
}

//-- commands --//

// +convert:create=Receipt
type CreateReceiptArgs struct {
	ShopID      int64
	TraderID    int64
	UserID      int64
	Code        string
	Title       string
	Description string
	Amount      int32
	OrderIDs    []int64
	Lines       []*ReceiptLine
}

// +convert:update=Receipt(ID,ShopID)
type UpdateReceiptArgs struct {
	ID          int64
	ShopID      int64
	TraderID    NullInt64
	UserID      NullInt64
	Title       NullString
	Description NullString
	Amount      NullInt32
	OrderIDs    []int64
	Lines       []*ReceiptLine
}
