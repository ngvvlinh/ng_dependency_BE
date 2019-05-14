package ahamoveclient

import (
	"encoding/json"

	"etop.vn/api/main/shipnow"
	"etop.vn/backend/pkg/etop/model"
)

type TransportType string

const (
	TransportRoad TransportType = "road"
	TransportFly  TransportType = "fly"
)

func (t TransportType) Name() string {
	switch t {
	case TransportRoad:
		return model.ShippingServiceNameStandard
	case TransportFly:
		return model.ShippingServiceNameFaster
	}
	return string(t)
}

type AhamoveAccount struct {
	AccountID string `yaml:"account_id"`
	Token     string `yaml:"token"`
}

type OrderState string
type OrderSubState string

// DeliveryStatus: path[x].status, x > 0
type DeliveryStatus string

const (
	StateConfirmed OrderState = "IDLE"       // Order confirmed
	StateAssigning OrderState = "ASSIGNING"  // Finding a supplier
	StateAccepted  OrderState = "ACCEPTED"   // Once supplier accepts the order, the order status will be changed to ACCEPTED
	StateInProcess OrderState = "IN PROCESS" // When the supplier collects the food at the restaurant, he will click “Pick up” and the order status will be changed to IN PROCESS
	StateCompleted OrderState = "COMPLETED"  // When the supplier completes the order, the order status will be changed to COMPLETED
	StateCancelled OrderState = "CANCELLED"

	StateBoarded    OrderSubState = "BOARDED"    // When the supplier has arrived at the pickup location (restaurant), our system will set sub_status of the order to BOARDED
	StateCompleting OrderSubState = "COMPLETING" // When the supplier location is nearby drop-off location, our system will set sub_status of the order to COMPLETING

	StateFfmCompleted DeliveryStatus = "COMPLETED" // When the supplier gives the order to recipient & collects cash, he will choose Complete and the according stop point status will be changed to COMPLETED
	StateFfmFailed    DeliveryStatus = "FAILED"    // When the supplier arrives at the recipient point, but the recipient does not show up, he will choose Fail and the according stop point status will be changed to FAILED
)

func (orderState OrderState) ToModel() shipnow.State {
	switch orderState {
	case StateConfirmed:
		return shipnow.StateCreated
	case StateAssigning:
		return shipnow.StateAssigning
	case StateAccepted:
		return shipnow.StatePicking
	case StateInProcess:
		return shipnow.StateDelivering
	case StateCompleted:
		return shipnow.StateDelivered
	case StateCancelled:
		return shipnow.StateCancelled
	default:
		return shipnow.StateDefault
	}
}

func (orderState OrderState) ToStatus5() model.Status5 {
	switch orderState {
	case StateCancelled:
		return model.S5Negative
	case StateCompleted:
		return model.S5Positive
	default:
		return model.S5SuperPos
	}
}

type CommonErrorResponse struct {
	Title string `url:"title"`
}

type CalcShippingFeeRequest struct {
	OrderTime int `url:"order_time"`
	// IdleUntil: always set IdleUntil = OrderTime
	IdleUntil      int                     `url:"idle_until"`
	Path           string                  `url:"path"`
	DeliveryPoints []*DeliveryPointRequest `url:"-"`
	ServiceID      string                  `url:"service_id"`
}

type DeliveryPointRequest struct {
	Address        string  `json:"address"`
	ProvinceCode   string  `json:"-"`
	DistrictCode   string  `json:"-"`
	WardCode       string  `json:"-"`
	Lat            float32 `json:"lat"`
	Lng            float32 `json:"lng"`
	Mobile         string  `json:"mobile"`
	Name           string  `json:"name"`
	COD            int     `json:"cod"`
	TrackingNumber string  `json:"tracking_number"`
	Remarks        string  `json:"remarks"`
	SenderName     string  `json:"sender_name"`
	SenderMobile   string  `json:"sender_mobile"`
}

func ConvertDeliveryPointsRequestToString(points []*DeliveryPointRequest) string {
	path, err := json.Marshal(points)
	if err != nil {
		return ""
	}
	return string(path)
}

type DeliveryPoint struct {
	// Supported address in template bellow, without lat, lng in parameters
	// {user_free_text_input},{ward},{province/district},{city}
	Address        string  `url:"address"`
	Mobile         string  `url:"mobile"`
	Name           string  `url:"name"`
	COD            int     `url:"cod"`
	TrackingNumber string  `url:"tracking_number"`
	Remarks        string  `url:"remarks"`
	SenderName     string  `url:"sender_name"`
	SenderMobile   string  `url:"sender_mobile"`
	Lat            float32 `url:"lat"`
	Lng            float32 `url:"lng"`

	RequirePod          bool `url:"require_pod"`          // True, # If true, supplier will be required proof of delivery before complete the transaction
	RequireVerification bool `url:"require_verification"` // True, # If true, a verification code will be sent to the receiver, then the supplier will need to enter this code to complete the transaction

	RatingByReceiver  int    `url:"rating_by_receiver"`  // 4,
	CommentByReceiver string `url:"comment_by_receiver"` // "Good",

	CompleteTime    float32 `url:"complete_time"`    // 1426671164,
	CompleteLat     float32 `url:"complete_lat"`     // 10.7890462,
	CompleteLng     float32 `url:"complete_lng"`     // 106.7763078,
	CompleteComment string  `url:"complete_comment"` // "Nice receiver",
	ImageUrl        string  `url:"image_url"`        // "//i.imgur.com/lchC2xz.jpg", # Image to show that the supplier has already completed the order, uploaded by the supplier
	PodInfo         string  `url:"pod_info"`         // "024792155", # Proof of delivery information that supplier has collected, can be recipient's ID number or image URL
	Status          string  `url:"status"`           // "COMPLETED"
}

type CalcShippingFeeResponse struct {
	Distance            float32 `json:"distance"`              // 1.144
	Duration            int     `json:"duration"`              // 288
	Currency            string  `json:"currency"`              // "VND"
	Discount            int     `json:"discount"`              // 30000
	DistanceFee         int     `json:"distance_fee"`          // 140000,
	StopFee             int     `json:"stop_fee"`              // 0,
	TotalFee            int     `json:"total_fee"`             // 275000,  # Total = Distance Fee + Request Fee + Stop Fee
	UserMainAccount     int     `json:"user_main_account"`     // 0,
	UserBonusAccount    int     `json:"user_bonus_account"`    // 0,
	TotalPay            int     `json:"total_pay"`             // 275000,
	Surcharge           string  `json:"surcharge"`             // 1.1 # 110% compared to normal fee, only returned if surcharge > 1,
	BookingTimeError    string  `json:"booking_time_error"`    // "Booking time not valid" #If booking time is not in opening hours period,
	PaymentErrorMessage string  `json:"payment_error_message"` // "Not enough credit" #If user does not have enough credit in balance
}

type CreateOrderRequest struct {
	Token     string `url:"token"`
	OrderTime int    `url:"order_time"`
	// IdleUntil: always set IdleUntil = OrderTime
	IdleUntil      int                     `url:"idle_until"`
	DeliveryPoints []*DeliveryPointRequest `url:"-"`
	Path           string                  `url:"path"`
	ServiceID      string                  `url:"service_id"`
	Remarks        string                  `url:"remarks"`
	// Payment method chose by user (BALANCE or CASH) - Optional
	PaymentMethod string `url:"payment_method"`
	Products      []Item `url:"-"`
	Items         string `url:"items"`
}

type Item struct {
	ID    string `url:"_id"`   // "TS",
	Num   int    `url:"num"`   // 2,
	Name  string `url:"name"`  // "Sua tuoi",
	Price int    `url:"price"` // 15000
}

func ConvertItemsToString(items []Item) string {
	res, err := json.Marshal(items)
	if err != nil {
		return ""
	}
	return string(res)
}

type CreateOrderResponse struct {
	OrderId    string `json:"order_id"`    // "QMB0HO",
	Status     string `json:"status"`      // "ASSIGNING",
	SharedLink string `json:"shared_link"` // "https://cloud.ahamove.com/share-order/51V42K/84909055578",
	Order      *Order `json:"order"`
}

type Order struct {
	Duration int              `json:"duration"` // 5015,
	Path     []*DeliveryPoint `json:"path"`

	// Status, Service, Request
	ID        string `json:"_id"`        // "1JFU54",
	Status    string `json:"status"`     // "CANCELLED",
	ServiceId string `json:"service_id"` // "SGN-TRICYCLE", # Service ID
	CityId    string `json:"city_id"`    // "SGN", # City ID

	// User ID & Supplier ID
	UserId       string `json:"user_id"`       // "84972709963", # User who created this order
	UserName     string `json:"user_name"`     // "Thao Vy", # User name
	SupplierId   string `json:"supplier_id"`   // "84909561477", # Supplier who accepted this order
	SupplierName string `json:"supplier_name"` // "Hieu Nguyen", # Supplier name

	// Partner
	Partner string `json:"partner"` // "topship", # If this order is created by 3rd-party partners

	// Time
	OrderTime      float32 `json:"order_time"`      // 1426297774, # Pick-up time
	CreateTime     float32 `json:"create_time"`     // 1426297174.649361, # Order create time
	AcceptTime     float32 `json:"accept_time"`     // 1426298634.189912, # Order accept time
	AcceptLat      float32 `json:"accept_lat"`      // 10.7890462,
	AcceptLng      float32 `json:"accept_lng"`      // 106.7763078,
	AcceptDistance float32 `json:"accept_distance"` // 1.2, # Distance from accepted location to pick up point
	AcceptDuration int     `json:"accept_duration"` // 150, # Time from accepted location to pick up point

	CancelTime    float32 `json:"cancel_time"`    // 1426671164.374609, # Order cancel time
	CancelComment string  `json:"cancel_comment"` // "Lich chuyen nha bi hoan", # Cancelled reason by user or supplier
	CancelByUser  bool    `json:"cancel_by_user"` // False, # True if this order is cancelled by user, otherwise False
	// Others: complete_time, complete_lat, complete_lng, accept_lat, accept_lng, fail_tome, fail_lat, fail_lng, fail_comment

	// Payment
	Currency    string  `json:"currency"`     // "VND",
	StopFee     int     `json:"stop_fee"`     // 0,
	RequestFee  int     `json:"request_fee"`  // 0,
	Distance    float32 `json:"distance"`     // 5.597,
	DistanceFee int     `json:"distance_fee"` // 142358,
	PromoCode   string  `json:"promo_code"`   // "AHAMOVE",
	Discount    int     `json:"discount"`     // 0,
	TotalFee    int     `json:"total_fee"`    // 142358, # Total fee = Distance Fee + Request Fee + Stop Fee - Discount
	VatFee      int     `json:"vat_fee"`      // 0,

	// Use credit from user account if available
	UserBonusAccount int `json:"user_bonus_account"` // 0,
	UserMainAccount  int `json:"user_main_account"`  // 0,
	TotalPay         int `json:"total_pay"`          // 142358, # Total pay = Total fee - User Main account - User Bonus account

	// Rating
	RatingByUser       int    `json:"rating_by_user"`        // 4,
	CommentByUser      string `json:"comment_by_user"`       // "Good",
	RatingBySupplier   int    `json:"rating_by_supplier"`    // 5,
	CommentBySupplier  string `json:"comment_by_supplier"`   // "Great customer",
	RatingByReceiver   int    `json:"rating_by_receiver"`    // 4,
	CommentByReceiver  string `json:"comment_by_receiver"`   // "Good driver",
	StoreRatingByUser  int    `json:"store_rating_by_user"`  // 4, # For food orders
	StoreCommentByUser string `json:"store_comment_by_user"` // "Delicious food", # For food orders

	// Others
	Remarks    string `json:"remarks"`     // "Den noi goi dien cho toi", # Note to supplier
	Remind     bool   `json:"remind"`      // True, # If this is advance booking and the system already reminds users
	AssignedBy string `json:"assigned_by"` // "84908842285", # If this order is assigned by a user or "auto" if by the system
	Index      int    `json:"index"`       // 1, # 0 if first order request from the user, 1 for second...
}

type GetOrderRequest struct {
	Token   string `url:"token"`
	OrderID string `url:"order_id"`
}

type CancelOrderRequest struct {
	Token    string `url:"token"`     // Token (User token or Supplier token, Admin if want to cancel a completed order)
	OrderId  string `url:"order_id"`  // Order ID
	Comment  string `url:"comment"`   // Optional
	ImageUrl string `url:"image_url"` // Optional, image URL to show the reason why supplier cancel an order
}
