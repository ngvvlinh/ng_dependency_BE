package kiotviet

import (
	"sync"
	"time"

	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/etop/model"
)

type Int = httpreq.Int
type String = httpreq.String
type Time = httpreq.Time

type Connection struct {
	RetailerID string
	TokenStr   string
	ExpiresAt  time.Time

	m sync.Mutex
}

func (c *Connection) URL(path string) string {
	return BaseURLPublic + path
}

type AccessTokenResponse struct {
	TokenStr     string `json:"access_token"`
	ExpiresIn    Int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	ErrorStr       string `json:"error"`
	ResponseStatus struct {
		ErrorCode string `json:"errorCode"`
		Message   string `json:"message"`
	} `json:"responseStatus"`
}

func (r *ErrorResponse) Error() string {
	if r.ErrorStr != "" {
		return r.ErrorStr
	}
	if r.ResponseStatus.ErrorCode != "" {
		if r.ResponseStatus.Message == "" {
			return "Mã " + r.ResponseStatus.ErrorCode
		}
		return r.ResponseStatus.Message +
			" (" + r.ResponseStatus.ErrorCode + ")"
	}
	return "No error"
}

type TotalResponse struct {
	Total    Int `json:"total"`
	PageSize Int `json:"pageSize"`
}

type Branch struct {
	ID            String `json:"id"`
	Name          string `json:"branchName"`
	Code          string `json:"branchCode"`
	ContactNumber string `json:"contactNumber"`
	RetailerID    String `json:"retailerId"`
	Address       string `json:"address"`
	UpdatedAt     Time   `json:"modifiedDate"`
	CreatedAt     Time   `json:"createdDate"`
}

type Category struct {
	ID         String `json:"categoryId"`
	Name       string `json:"categoryName"`
	ParentID   String `json:"parentId"`
	RetailerID String `json:"retailerID"`
	UpdatedAt  Time   `json:"modifiedDate"`
	CreatedAt  Time   `json:"createdDate"`
}

type Product struct {
	ID              String       `json:"id"`
	RetailerID      String       `json:"retailerId"`
	Code            string       `json:"code"`
	Name            string       `json:"name"`
	FullName        string       `json:"fullName"`
	CategoryID      String       `json:"categoryId"`
	AllowsSale      bool         `json:"allowsSale"`
	HasVariants     bool         `json:"hasVariants"`
	BasePrice       Int          `json:"basePrice"`
	Unit            string       `json:"unit"`
	Units           []*Unit      `json:"units"`
	MasterProductID String       `json:"masterProductId"`
	MasterUnitID    String       `json:"masterUnitId"`
	Conversionvalue float64      `json:"conversionvalue"`
	Description     string       `json:"description"`
	UpdatedAt       Time         `json:"modifiedDate"`
	IsActive        bool         `json:"isActive"`
	Attributes      []*Attribute `json:"attributes"`
	Inventories     []*Inventory `json:"inventories"`
	Images          []string     `json:"images"`
}

type Unit struct {
	ID        String `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	FullName  string `json:"fullName"`
	Unit      string `json:"unit"`
	BasePrice Int    `json:"basePrice”"`
}

type Attribute struct {
	ID    String `json:"productId"`
	Name  string `json:"attributeName"`
	Value string `json:"attributeValue"`
}

type Inventory struct {
	ProductID   String `json:"productId"`
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
	BranchID    String `json:"branchId"`
	BranchName  string `json:"branchName"`
	Cost        Int    `json:"cost"`
	OnHand      Int    `json:"onHand"`
	Reserved    Int    `json:"reserved"`
}

// -- Commands

type RequestAccessTokenCommand struct {
	ClientID     string
	ClientSecret string
	Result       struct {
		TokenStr  string
		ExpiresAt time.Time
	}
}

type GetOrRenewTokenCommand struct {
	SupplierID int64
	Result     struct {
		// Supplier *model.Supplier
		// Kiotviet *model.SupplierKiotviet
		*Connection
	}
}

type SyncProductSourceCommand struct {
	SourceID      int64
	FromBeginning bool
}

type SyncCategoriesCommand struct {
	SourceID   int64
	Connection *Connection
	SyncState  SyncState

	Result struct {
		Response  *UpdatedCategoriesResponse
		NextState SyncState
		Done      bool
	}
}

type SyncProductsCommand struct {
	SourceID   int64
	BranchID   string
	Connection *Connection
	SyncState  SyncState

	Result struct {
		Response  *UpdatedProductsResponse
		NextState SyncState
		Done      bool
	}
}

type BranchResult struct {
	TotalResponse
	Data       []*Branch `json:"data"`
	DeletedIDs []String  `json:"removedIds"`
}

type UpdatedCategoriesResponse struct {
	TotalResponse
	Data       []*Category `json:"data"`
	DeletedIDs []String    `json:"removedIds"`
}

type UpdatedProductsResponse struct {
	TotalResponse
	Data       []*Product `json:"data"`
	DeletedIDs []String   `json:"removedIds"`
}

type RequestUpdatedCategoriesCommand struct {
	Connection *Connection
	UpdatedQuery

	Result *UpdatedCategoriesResponse
}

type RequestUpdatedProductsCommand struct {
	Connection *Connection
	UpdatedQuery

	Result *UpdatedProductsResponse
}

func ToModelBranches(bs []*Branch) []*model.KiotvietBranch {
	res := make([]*model.KiotvietBranch, len(bs))
	for i, b := range bs {
		res[i] = ToModelBranch(b)
	}
	return res
}

func ToModelBranch(b *Branch) *model.KiotvietBranch {
	return &model.KiotvietBranch{
		ID:            string(b.ID),
		Name:          b.Name,
		Code:          b.Code,
		ContactNumber: b.ContactNumber,
		RetailerID:    string(b.RetailerID),
		Address:       b.Address,
		UpdatedAt:     time.Time(b.UpdatedAt),
		CreatedAt:     time.Time(b.CreatedAt),
	}
}

type EnsureWebhooksCommand struct {
	*Connection
	SuppierID int64
}

type RequestWebhooksCommand struct {
	*Connection

	Result *WebhooksResponse
}

type WebhookRequestItem struct {
	Type        string `json:"Type"`
	URL         string `json:"Url"`
	IsActive    bool   `json:"IsActive"`
	Description string `json:"Description"`
}

type WebhooksRequest struct {
	Webhook WebhookRequestItem `json:"Webhook"`
}

type WebhooksResponse struct {
	TotalResponse
	Data []*Webhook
}

type Webhook struct {
	ID          String "id"
	Type        String "type"
	URL         String "url"
	IsActive    bool   "isActive"
	Description String "description"
	RetailerID  String "retailerId"
}

type RequestCreateWebhookCommand struct {
	*Connection
	Webhook WebhookRequestItem
}

type RequestDeleteWebhookCommand struct {
	*Connection
	ID string
}
