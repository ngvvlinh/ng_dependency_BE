package ahamove

import (
	shipnowtypes "etop.vn/api/main/shipnow/types"
	ahamoveclient "etop.vn/backend/pkg/integration/ahamove/client"
)

type CalcShippingFeeCommand struct {
	ArbitraryID int64 // This is provided as a seed, for stable randomization

	Request *ahamoveclient.CalcShippingFeeRequest
	Result  []*shipnowtypes.ShipnowService
}

type CalcSingleShippingFeeCommand struct {
	ServiceID string

	FromDistrictCode string
	ToDistrictCode   string

	Request *ahamoveclient.CalcShippingFeeRequest

	Result *shipnowtypes.ShipnowService
}

type CreateOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveclient.CreateOrderRequest
	Result  *ahamoveclient.CreateOrderResponse
}

type GetOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveclient.GetOrderRequest
	Result  *ahamoveclient.Order
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveclient.CancelOrderRequest
}

type GetServiceCommand struct {
	Request *ahamoveclient.GetServicesRequest
	Result  []*ahamoveclient.ServiceType
}

type RegisterAccountCommand struct {
	Request *ahamoveclient.RegisterAccountRequest
	Result  *ahamoveclient.RegisterAccountResponse
}

type GetAccountCommand struct {
	Request *ahamoveclient.GetAccountRequest
	Result  *ahamoveclient.Account
}

type VerifyAccountCommand struct {
	Request *ahamoveclient.VerifyAccountRequest
	Result  *ahamoveclient.VerifyAccountResponse
}
