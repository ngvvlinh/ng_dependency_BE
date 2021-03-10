package client

import (
	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_fee_type"
	"o.o/api/top/types/etc/try_on"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
)

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type PartnerAccountCfg struct {
	Connection  *connectioning.Connection
	Token       string
	AffiliateID string
}

func validateConnection(conn *connectioning.Connection) error {
	if conn.ConnectionProvider != connection_type.ConnectionProviderPartner || conn.DriverConfig == nil {
		return cm.Errorf(cm.FailedPrecondition, nil, "Không thể khởi tạo direct shippment. Connection không hợp lệ")
	}
	driverCfg := conn.DriverConfig
	errMsg := ""
	switch {
	case driverCfg.CreateFulfillmentURL == "":
		errMsg = "Thiếu create_fulfillment_url"
	case driverCfg.TrackingURL == "":
		errMsg = "Thiếu tracking_url"
	case driverCfg.GetFulfillmentURL == "":
		errMsg = "Thiếu get_fulfillment_url"
	case driverCfg.CancelFulfillmentURL == "":
		errMsg = "Thiếu cancel_fulfillment_url"
	case driverCfg.GetShippingServicesURL == "":
		errMsg = "Thiếu get_shipping_services_url"
	}
	if errMsg != "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Không thể khởi tạo direct shipment. %v", errMsg)
	}
	return nil
}

type ResponseInterface interface {
	GetCommonResponse() CommonResponse
}

type CommonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c CommonResponse) GetCommonResponse() CommonResponse {
	return c
}

type GetShippingServicesRequest struct {
	BasketValue      int           `json:"basket_value"`
	TotalWeight      int           `json:"total_weight"`
	PickupAddress    SimpleAddress `json:"pickup_address"`
	ShippingAddress  SimpleAddress `json:"shipping_address"`
	IncludeInsurance bool          `json:"include_insurance"`
	// @Deprecated use CODAmount instead
	TotalCODAmount int `json:"total_cod_amount"`
	CODAmount      int `json:"cod_amount"`
}

type SimpleAddress struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	Province string `json:"province"`
	District string `json:"district"`
	Ward     string `json:"ward"`
}

type GetShippingServiceResponse struct {
	CommonResponse
	Data []*ShippingService `json:"data"`
}

type ShippingService struct {
	ServiceCode        String `json:"service_code"`
	Name               String `json:"name"`
	ServiceFee         Int    `json:"service_fee"`
	ServiceFeeMain     Int    `json:"service_fee_main"`
	ExpectedPickAt     Time   `json:"expected_pick_at"`
	ExpectedDeliveryAt Time   `json:"expected_delivery_at"`
}

type CreateFulfillmentRequest struct {
	PickupAddress   SimpleAddress `json:"pickup_address"`
	ShippingAddress SimpleAddress `json:"shipping_address"`
	Lines           []*ItemLine   `json:"lines"`
	TotalWeight     int           `json:"total_weight"`
	BasketValue     int           `json:"basket_value"`
	// @Deprecated use CODAmount instead
	TotalCODAmount      int              `json:"total_cod_amount"`
	CODAmount           int              `json:"cod_amount"`
	ShippingNote        string           `json:"shipping_note"`
	IncludeInsurance    bool             `json:"include_insurance"`
	ShippingServiceCode string           `json:"shipping_service_code"`
	ShippingFee         int              `json:"shipping_fee"`
	AffiliateID         string           `json:"affiliate_id"`
	TryOn               try_on.TryOnCode `json:"try_on"`
}

type ItemLine struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type CreateFulfillmentResponse struct {
	CommonResponse
	Data *Fulfillment `json:"data"`
}

type Fulfillment struct {
	FulfillmentID      String              `json:"fulfillment_id"`
	ShippingCode       String              `json:"shipping_code"`
	SortCode           String              `json:"sort_code"`
	ShippingFee        Int                 `json:"shipping_fee"`
	TrackingURL        String              `json:"tracking_url"`
	ExpectedPickAt     Time                `json:"expected_pick_at"`
	ExpectedDeliveryAt Time                `json:"expected_delivery_at"`
	ShippingFeeLines   []*ShippingFeeLine  `json:"shipping_fee_lines"`
	ShippingState      shippingstate.State `json:"shipping_state"`
	TryOn              try_on.TryOnCode    `json:"try_on"`
}

type ShippingFeeLine struct {
	Cost            Int                               `json:"cost"`
	ShippingFeeType shipping_fee_type.ShippingFeeType `json:"shipping_fee_type"`
}

type GetFulfillmentRequest struct {
	FulfillmentID String `json:"fulfillment_id"`
	ShippingCode  String `json:"shipping_code"`
}

type GetFulfillmentResponse struct {
	CommonResponse
	Data *Fulfillment `json:"data"`
}

type CancelFulfillmentRequest struct {
	FulfillmentID string `json:"fulfillment_id"`
	ShippingCode  string `json:"shipping_code"`
}

func (s *ShippingService) validate() error {
	errMsg := ""
	switch {
	case s.ServiceCode == "":
		errMsg = "Missing service_code"
	case s.Name == "":
		errMsg = "Missing service name"
	// case s.ExpectedPickAt.IsZero():
	// 	errMsg = "Missing expected_pick_at"
	// case s.ExpectedDeliveryAt.IsZero():
	// 	errMsg = "Missing expected_delivery_at"
	case s.ServiceFee == 0:
		errMsg = "Missing service_fee"
		// case s.ServiceFeeMain == 0:
		// 	errMsg = "Missing service_fee_main"
	}
	if errMsg != "" {
		return cm.Errorf(cm.InvalidArgument, nil, errMsg)
	}
	return nil
}

func (ffm *Fulfillment) validate() error {
	errMsg := ""
	switch {
	case ffm.FulfillmentID == "":
		errMsg = "Missing fulfillment_id"
	case ffm.ShippingCode == "":
		errMsg = "Missing shipping_code"
	case ffm.TrackingURL == "":
		errMsg = "Missing tracking_url"
	case ffm.ShippingState == 0:
		errMsg = "Missing shipping_state"
	// case ffm.ExpectedPickAt.IsZero():
	// 	errMsg = "Missing expected_pick_at"
	// case ffm.ExpectedDeliveryAt.IsZero():
	// 	errMsg = "Missing expected_delivered_at"
	case ffm.ShippingFeeLines == nil || len(ffm.ShippingFeeLines) == 0:
		errMsg = "Missing shipping_fee_lines"
	}
	for _, line := range ffm.ShippingFeeLines {
		if line.ShippingFeeType == 0 {
			errMsg = "shipping_fee_type does not valid"
			break
		}
	}
	if errMsg != "" {
		return cm.Errorf(cm.InvalidArgument, nil, errMsg)
	}
	return nil
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	CommonResponse
	Data *AccountData `json:"data"`
}

type AccountData struct {
	UserID String `json:"user_id"`
	Token  String `json:"token"`
}

func (a *AccountData) validate() error {
	if a == nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Account can not be empty")
	}
	errMsg := ""
	switch {
	case a.UserID == "":
		errMsg = "Missing user_id"
	case a.Token == "":
		errMsg = "Missing token"
	}
	if errMsg != "" {
		return cm.Errorf(cm.InvalidArgument, nil, errMsg)
	}
	return nil
}

type SignUpRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type SignUpResponse struct {
	CommonResponse
	Data *AccountData `json:"data"`
}
