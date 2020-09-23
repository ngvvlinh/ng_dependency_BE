package client

import (
	"math"
	"strings"
	"time"

	typesshipping "o.o/api/top/types/etc/shipping"
	shippingsubstate "o.o/api/top/types/etc/shipping/substate"
	"o.o/api/top/types/etc/status5"
	"o.o/backend/com/main/address/model"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/common/strs"
)

type (
	Bool   = httpreq.Bool
	Float  = httpreq.Float
	Int    = httpreq.Int
	String = httpreq.String
	Time   = httpreq.Time
)

type NinjaVanCfg struct {
	Token     string `json:"token"`
	ClientID  string `json:"client_id"`
	SecretKey string `json:"secret_key"`
}

type State string

const (
	// A Order has been created and is in Staging phase
	StateStaging State = "Staging"

	// Đang lấy
	// Order has been confirmed and is pending pickup
	StatePendingPickup State = "Pending Pickup"

	// Chờ pick up lại
	// A van has been dispatched to pick up Order
	StateVanEnrouteToPickup State = "Van en-route to pickup"

	// Đã lấy
	// Order as been picked up and is en-route to Sorting Hub
	StateEnrouteToSortingHub State = "En-route to Sorting Hub"

	// Đã lấy
	// Order has arrived at Sorting Hub and has been processed successfully
	StateArrivedAtSortingHub State = "Arrived at Sorting Hub"

	// Đã lấy
	// Order has arrived at Origin Hub and has been processed successfully
	StateArrivedAtOriginHub State = "Arrived at Origin Hub"

	// Đang giao
	// Order is on van, en-route to delivery
	StateOnVehicleForDelivery State = "On Vehicle for Delivery"

	// Đã giao
	// Delivery has been successfully completed
	StateCompleted State = "Completed"

	// Delivery has failed and the Order is pending re-schedule
	StatePendingReschedule State = "Pending Reschedule"

	// Lấy không thành công
	// Pickup has failed and the Order is awaiting re-schedule
	StatePickupFail State = "Pickup fail"

	// Đã huỷ
	// Order has been cancelled
	StateCancelled State = "Cancelled"

	// Đã trả
	// Delivery of Order has failed repeatedly, sending back to Sender
	StateReturnedToSender State = "Returned to Sender"

	// The parcel size of an Order has been changed
	StateParcelSize State = "Parcel Size"

	// Chờ giao
	// The parcel has been placed at the Distribution Point for customer collection
	StateArrivedAtDistributionPoint State = "Arrived at Distribution Point"

	// Đã giao
	// Proof of delivery is now ready
	StateSuccessfulDelivery State = "Successful Delivery"

	// Đã lấy
	// Proof of pickup is now ready
	StateSuccessFulPickup State = "Successful Pickup"

	// The parcel weight of an Order has been changed
	StateParcelWeight State = "Parcel Weight"

	// Order is in cross border leg or is pending tax payment from consignee if required
	StateCrossBorderTransit State = "Cross Border Transit"

	// Order is ready for pickup from customs warehouse
	StateCustomsCleared State = "Customs Cleared"

	// Order is in customs clearance exception
	StateCustomsHeld State = "Customs Held"

	// Đang trả
	// Return to Sender request has been triggered
	StateReturnToSenderTriggered State = "Return to Sender Triggered"

	// Order has been received at Distribution Point and is pending pickup
	StatePendingPickupAtDistributionPoint State = "Pending Pickup at Distribution Point"

	// The parcel size, or parcel weight, or parcel dimensions of an Order has been changed
	StateParcelMeasurementsUpdate State = "Parcel Measurements Update"
)

func (s State) ToModel() typesshipping.State {
	switch s {
	case StateStaging:
		return typesshipping.Created
	case StateCancelled:
		return typesshipping.Cancelled
	case StatePendingPickup, StatePickupFail:
		return typesshipping.Picking
	case StateArrivedAtSortingHub,
		StateSuccessFulPickup,
		StateArrivedAtDistributionPoint,
		StateArrivedAtOriginHub,
		StatePendingPickupAtDistributionPoint,
		StateEnrouteToSortingHub:
		return typesshipping.Holding
	case StateOnVehicleForDelivery,
		StatePendingReschedule:
		return typesshipping.Delivering
	case StateVanEnrouteToPickup,
		StateReturnToSenderTriggered:
		return typesshipping.Returning
	case StateReturnedToSender:
		return typesshipping.Returned
	case StateSuccessfulDelivery, StateCompleted:
		return typesshipping.Delivered
	default:
		return typesshipping.Unknown
	}
}

func (s State) ToSubstateModel() shippingsubstate.Substate {
	switch s {
	case StatePendingReschedule:
		return shippingsubstate.DeliveryFail
	case StatePickupFail:
		return shippingsubstate.PickFail
	default:
		return 0
	}
}

func (s State) ToStatus5() status5.Status {
	switch s.ToModel() {
	case typesshipping.Cancelled:
		return status5.N
	case typesshipping.Returned, typesshipping.Returning, typesshipping.Undeliverable:
		return status5.NS
	case typesshipping.Delivered:
		return status5.P
	}

	return status5.S
}

type AddressType string

type TimeZone string

type CountryCode string

const (
	TimeZoneHCM = "Asia/Ho_Chi_Minh"

	LayoutISO = "2006-01-02"

	ClientCredentials string = "client_credentials"

	VN CountryCode = "VN"

	ThreeDays = 3 * 24 * time.Hour
)

var MapReasons = map[State]map[string]string{
	StatePickupFail: {
		"Nobody At Location":     "No one present at the pickup location",
		"Inaccurate Address":     "Could not locate pickup location",
		"Parcel Not Available":   "Item was not available at the pickup location, usually because shipper was unexpectedly out of stock, or packing was delayed.",
		"Parcel Too Bulky":       "Parcel was too bulky to be picked up, usually because size of parcel might have been mis-reported.",
		"Cancellation Requested": "Pickup was cancelled on request of the shipper, i.e. fr some reason the shipper did not want to fulfill the order at all.",
	},
	StatePendingReschedule: {
		"Address is correct but customer is not available":                 "Address is correct but customer is not available",
		"Customer requested change of delivery date / time":                "Customer requested change of delivery date / time",
		"Customer requested change of delivery location":                   "Customer requested change of delivery location",
		"Address on AWB is correct":                                        "Address on AWB is correct",
		"Address on AWB is incorrect":                                      "Address on AWB is incorrect",
		"Address was outside of driver coverage area":                      "Address was outside of driver coverage area",
		"I had insufficient time to complete all my deliveries":            "I had insufficient time to complete all my deliveries",
		"Vehicle breakdown":                                                "Vehicle breakdown",
		"Delay due to unexpected traffic conditions":                       "Delay due to unexpected traffic conditions",
		"Delay due to natural disasters or nationwide emergencies":         "Delay due to natural disasters or nationwide emergencies",
		"Unattempted - Parcel issues":                                      "Unattempted - Parcel issues",
		"Office address, closed":                                           "Office address, closed",
		"Office address, open but no one to receive":                       "Office address, open but no one to receive",
		"Residential address, but no one to receive":                       "Residential address, but no one to receive",
		"Package is fine - unable to collect COD":                          "Package is fine - unable to collect COD",
		"Package is fine - Customer wishes to cancel order":                "Package is fine - Customer wishes to cancel order",
		"Package is defective - Damaged":                                   "Package is defective - Damaged",
		"Package is defective - Wrong item inside":                         "Package is defective - Wrong item inside",
		"Location is inaccessible (restricted area)":                       "Location is inaccessible (restricted area)",
		"Refused entry by security personnel":                              "Refused entry by security personnel",
		"Driver cannot find location (lat-long issues)":                    "Driver cannot find location (lat-long issues)",
		"Incomplete Address provided (no unit number, block number, etc.)": "Incomplete Address provided (no unit number, block number, etc.)",
		"Specified Address was incorrect or Recipient has moved":           "Specified Address was incorrect or Recipient has moved",
	},
	StateReturnedToSender: {
		"Return to sender: Nobody at address":                       "No one present at the delivery location",
		"Return to sender: Unable to find Address":                  "Could not locate delivery location",
		"Return to sender: Item refused at Doorstep":                "Item was refused at the delivery location",
		"Return to sender: Refused to pay COD":                      "Customer refused to pay Cash-on-delivery amount",
		"Return to sender: Customer delayed beyond delivery period": "Customer requested for deliver re-schedule at a date far in the future (eg. 1 month later) or the delivery failed 3 times. A failure of 3 times would include 2 attempts that were re-scheduled by the customer.",
		"Return to sender: Cancelled by Shipper":                    "Order was cancelled by Shipper",
	},
}

type ErrorResponse struct {
	ErrorResponseDetail *ErrorResponseDetail `json:"error"`
}

type ErrorResponseDetail struct {
	Code      String         `json:"code"`
	RequestID String         `json:"request_id"`
	Title     String         `json:"title"`
	Message   String         `json:"message"`
	Details   []*ErrorDetail `json:"details"`
}

func (e *ErrorResponse) Error() (s string) {
	defer func() {
		s = strs.TrimLastPunctuation(s)
	}()

	errorResponse := e.ErrorResponseDetail
	b := strings.Builder{}
	b.WriteString(errorResponse.Message.String())
	b.WriteString("(")
	for _, detail := range errorResponse.Details {
		b.WriteString("Field: ")
		b.WriteString(detail.Field.String())
		b.WriteRune(',')
		b.WriteString("Message: ")
		b.WriteString(detail.Message.String())
		b.WriteRune(',')
		b.WriteString("Reason: ")
		b.WriteString(detail.Reason.String())
		break
	}
	b.WriteString(")")
	return b.String()
}

type ErrorDetail struct {
	Reason  String `json:"reason"`
	Field   String `json:"field"`
	Message String `json:"message"`
}

type GenerateAccessTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type GenerateAccessTokenResponse struct {
	AccessToken String `json:"access_token"`
	Expires     Int    `json:"expires"`
	TokenType   String `json:"token_type"`
	ExpiresIn   Int    `json:"expires_in"`
}

type CreateOrderRequest struct {
	ServiceType             string       `json:"service_type"`
	ServiceLevel            string       `json:"service_level"`
	RequestedTrackingNumber string       `json:"requested_tracking_number"`
	Reference               *Reference   `json:"reference"`
	From                    *AddressInfo `json:"from"`
	To                      *AddressInfo `json:"to"`
	ParcelJob               *ParcelJob   `json:"parcel_job"`
	Marketplace             *Marketplace `json:"marketplace"`
}

type Marketplace struct {
	SellerID          string `json:"seller_id"`
	SellerCompanyName string `json:"seller_company_name"`
}

type CreateOrderResponse struct {
	TrackingNumber String       `json:"tracking_number"`
	ServiceType    String       `json:"service_type"`
	ServiceLevel   String       `json:"service_level"`
	OrderReference *Reference   `json:"order_reference"`
	From           *AddressInfo `json:"from"`
	To             *AddressInfo `json:"to"`
	ParcelJob      *ParcelJob   `json:"parcel_job"`
}

type Reference struct {
	MerchantOrderNumber string `json:"merchant_order_number"`
}

func ToAddress(address *model.Address) *AddressInfo {
	if address == nil {
		return nil
	}
	return &AddressInfo{
		Name:        address.FullName,
		PhoneNumber: address.Phone,
		Email:       address.Email,
		Address: &Address{
			Address1: address.Address1,
			Address2: address.Address2,
			Country:  string(VN),
			Province: address.Province,
			District: address.District,
			Ward:     address.Ward,
		},
	}
}

type AddressInfo struct {
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	Address     *Address `json:"address"`
}

type Address struct {
	Address1    string `json:"address1"` // required
	Address2    string `json:"address2"`
	Country     string `json:"country"`      // required
	Province    string `json:"province"`     // required
	District    string `json:"district"`     // required
	Ward        string `json:"ward"`         // required
	AddressType string `json:"address_type"` // required
}

type ParcelJob struct {
	IsPickupRequired     bool              `json:"is_pickup_required"` // required
	PickupInstructions   string            `json:"pickup_instructions"`
	PickupAddress        *AddressInfo      `json:"pickup_address"`
	DeliveryStartDate    string            `json:"delivery_start_date"` // yyyy-MM-dd
	DeliveryTimeslot     *DeliveryTimeslot `json:"delivery_timeslot"`
	DeliveryInstructions string            `json:"delivery_instructions"`
	CashOnDelivery       float64           `json:"cash_on_delivery"`
	InsuredValue         float64           `json:"insured_value"`
	Dimensions           *Dimensions       `json:"dimensions"`
}

type DeliveryTimeslot struct {
	StartTime string `json:"start_time"` // yyyy-MM-dd
	EndTime   string `json:"end_time"`   // yyyy-MM-dd
	TimeZone  string `json:"timezone"`
}

type Dimensions struct {
	Weight float64 `json:"weight"` // kg
	Length float64 `json:"length"` // cm
	Width  float64 `json:"width"`  // cm
	Height float64 `json:"height"` // cm
}

type CancelOrderResponse struct {
	TrackingID String `json:"trackingId"`
	Status     String `json:"status"`
	// TODO(ngoc): parse time
	UpdatedAt String `json:"updatedAt"`
}

type FindAvailableServicesResponse struct {
	AvailableServices []*AvailableService
}

type AvailableService struct {
	Name String `json:"Name"`
}

type CallbackOrder struct {
	ID                   String       `json:"id"`
	ShipperID            Int          `json:"shipper_id"`
	ShipperRefNo         String       `json:"shipper_ref_no"`
	ShipperOrderRefNo    String       `json:"shipper_order_ref_no"`
	TrackingRefNo        String       `json:"tracking_ref_no"`
	Timestamp            Time         `json:"timestamp"`
	PreviousStatus       String       `json:"previous_status"`
	Status               String       `json:"status"`
	TrackingID           String       `json:"tracking_id"`
	OrderID              String       `json:"order_id"`
	PreviousSize         String       `json:"previous_size"`
	NewSize              String       `json:"new_size"`
	PreviousWeight       Float        `json:"previous_weight"`
	NewWeight            Float        `json:"new_weight"`
	Comments             String       `json:"comments"`
	PreviousMeasurements *Measurement `json:"previous_measurements"`
	NewMeasurements      *Measurement `json:"new_measurements"`
}

func (c *CallbackOrder) GetWeight(oldWeight int) int {
	// NinjaVan trả về khối lượng đơn vị kg
	newWeight, weightMeasurement := float64(0), float64(0)
	if c.PreviousWeight != c.NewWeight {
		newWeight = float64(newWeight)
	}
	if c.NewMeasurements != nil {
		weightMeasurement = math.Max(float64(c.NewMeasurements.MeasuredWeight), float64(c.NewMeasurements.VolumetricWeight))
	}
	res := math.Max(newWeight, weightMeasurement)
	if res != 0 {
		return int(res * 1000)
	}
	return oldWeight
}

type Measurement struct {
	Width            Float  `json:"width"`
	Height           Float  `json:"height"`
	Length           Float  `json:"length"`
	Size             String `json:"size"`
	MeasuredWeight   Float  `json:"measured_weight"`
	VolumetricWeight Float  `json:"Volumetric_weight"`
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}
