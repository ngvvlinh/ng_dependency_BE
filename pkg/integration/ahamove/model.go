package ahamove

import (
	"etop.vn/api/main/shipnow/carrier"

	shipnowtypes "etop.vn/api/main/shipnow/types"
	ahamoveClient "etop.vn/backend/pkg/integration/ahamove/client"
)

type CalcShippingFeeCommand struct {
	ArbitraryID int64 // This is provided as a seed, for stable randomization

	Request *ahamoveClient.CalcShippingFeeRequest
	Result  []*shipnowtypes.ShipnowService
}

type CalcSingleShippingFeeCommand struct {
	ServiceID string

	FromDistrictCode string
	ToDistrictCode   string

	Request *ahamoveClient.CalcShippingFeeRequest

	Result *shipnowtypes.ShipnowService
}

type CreateOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveClient.CreateOrderRequest
	Result  *ahamoveClient.CreateOrderResponse
}

type GetOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveClient.GetOrderRequest
	Result  *ahamoveClient.Order
}

type CancelOrderCommand struct {
	ServiceID string // Required for detecting which client

	Request *ahamoveClient.CancelOrderRequest
}

func ToShippingService(sfResp *ahamoveClient.CalcShippingFeeResponse, serviceID ServiceCode, providerServiceID string) *shipnowtypes.ShipnowService {
	if sfResp == nil {
		return nil
	}
	service := ServicesIndexID[serviceID]
	return &shipnowtypes.ShipnowService{
		Carrier:            carrier.Ahamove,
		Name:               service.Name,
		Code:               providerServiceID,
		Fee:                int32(sfResp.TotalFee),
		ExpectedPickupAt:   nil,
		ExpectedDeliveryAt: nil,
	}
}

type RegisterAccountCommand struct {
	Request *ahamoveClient.RegisterAccountRequest
	Result  *ahamoveClient.RegisterAccountResponse
}

type GetAccountCommand struct {
	Request *ahamoveClient.GetAccountRequest
	Result  *ahamoveClient.Account
}
