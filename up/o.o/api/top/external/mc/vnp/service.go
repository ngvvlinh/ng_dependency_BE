package vnp

import (
	"context"

	"o.o/api/top/external/shop"
	"o.o/api/top/external/types"
	cm "o.o/api/top/types/common"
)

// +gen:apix
// +gen:apix:base-path=/v1
// +gen:swagger:doc-path=external/mc/vnp
// +gen:swagger:title: MoveCrop - VNPost
// +gen:swagger:version=v1
// +gen:swagger:description=description.md

var _ shop.ShipnowService = ShipnowService(nil)

// +apix:path=/vnposts
// +swagger:tag: Giao Siêu Tốc
type ShipnowService interface {

	// Ping
	//
	// Kiểm tra token là hợp lệ và server đang hoạt động.
	//
	// +apix:path=ping
	Ping(context.Context, *cm.Empty) (*cm.Empty, error)

	// Get Services
	//
	// Lấy danh sách các dịch vụ Siêu Tốc.
	//
	// +apix:path=getservicesvnpost
	GetShipnowServices(context.Context, *types.GetShipnowServicesRequest) (*types.GetShipnowServicesResponse, error)

	// Create Order
	//
	// Tạo đơn giao Siêu Tốc.
	//
	// +apix:path=createordervnpost
	CreateShipnowFulfillment(context.Context, *types.CreateShipnowFulfillmentRequest) (*types.ShipnowFulfillment, error)

	// Cancel Order
	//
	// Huỷ đơn giao Siêu Tốc.
	//
	// +apix:path=cancelordervnpost
	CancelShipnowFulfillment(context.Context, *types.CancelShipnowFulfillmentRequest) (*cm.UpdatedResponse, error)

	// Get Order
	//
	// Lấy thông tin đơn giao Siêu Tốc.
	//
	// +apix:path=getordervnpost
	GetShipnowFulfillment(context.Context, *types.FulfillmentIDRequest) (*types.ShipnowFulfillment, error)
}
