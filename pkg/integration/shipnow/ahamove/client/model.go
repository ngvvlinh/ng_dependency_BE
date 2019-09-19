package client

import (
	"encoding/json"

	"etop.vn/backend/pkg/common/httpreq"

	"etop.vn/api/main/etop"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
)

type Config struct {
	Env    string `yaml:"env"`
	Name   string `yaml:"name"`
	ApiKey string `yaml:"api_key"`
}

type (
	Bool  = httpreq.Bool
	Float = httpreq.Float
	Int   = httpreq.Int
	Time  = httpreq.Time
)

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_AHAMOVE"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENV":     &c.Env,
		p + "_NAME":    &c.Name,
		p + "_API_KEY": &c.ApiKey,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env:    cm.PartnerEnvTest,
		Name:   "ahamove_test",
		ApiKey: "860160f707a7a8afc7bffa3e54630f40",
	}
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

	StateDeliveryCompleted DeliveryStatus = "COMPLETED" // When the supplier gives the order to recipient & collects cash, he will choose Complete and the according stop point status will be changed to COMPLETED
	StateDeliveryFailed    DeliveryStatus = "FAILED"    // When the supplier arrives at the recipient point, but the recipient does not show up, he will choose Fail and the according stop point status will be changed to FAILED

	PaymentMethodCash = "CASH"
)

func (orderState OrderState) ToCoreState() shipnowtypes.State {
	switch orderState {
	case StateConfirmed:
		return shipnowtypes.StateCreated
	case StateAssigning:
		return shipnowtypes.StateAssigning
	case StateAccepted:
		return shipnowtypes.StatePicking
	case StateInProcess:
		return shipnowtypes.StateDelivering
	case StateCompleted:
		return shipnowtypes.StateDelivered
	case StateCancelled:
		return shipnowtypes.StateCancelled
	default:
		return shipnowtypes.StateDefault
	}
}

func (orderState OrderState) ToStatus5() etop.Status5 {
	switch orderState {
	case StateCancelled:
		return etop.S5Negative
	case StateCompleted:
		return etop.S5Positive
	default:
		return etop.S5SuperPos
	}
}

func (s DeliveryStatus) ToCoreState(currentState shipnowtypes.State) shipnowtypes.State {
	switch s {
	case StateDeliveryCompleted:
		return shipnowtypes.StateDelivered
	case StateDeliveryFailed:
		return shipnowtypes.StateReturned
	default:
		return currentState
	}
}

func (s DeliveryStatus) ToStatus5() etop.Status5 {
	switch s {
	case StateDeliveryCompleted:
		return etop.S5Positive
	case StateDeliveryFailed:
		return etop.S5Negative
	default:
		return etop.S5SuperPos
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
	// Payment method chose by user (BALANCE or CASH or MOMO)
	PaymentMethod string `url:"payment_method"`
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
	COD            int32   `json:"cod"`
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
	Address        string  `json:"address"`
	Mobile         string  `json:"mobile"`
	Name           string  `json:"name"`
	COD            int     `json:"cod"`
	TrackingNumber string  `json:"tracking_number"`
	Remarks        string  `json:"remarks"`
	SenderName     string  `json:"sender_name"`
	SenderMobile   string  `json:"sender_mobile"`
	Lat            float32 `json:"lat"`
	Lng            float32 `json:"lng"`

	RequirePod          bool `json:"require_pod"`          // True, # If true, supplier will be required proof of delivery before complete the transaction
	RequireVerification bool `json:"require_verification"` // True, # If true, a verification code will be sent to the receiver, then the supplier will need to enter this code to complete the transaction

	RatingByReceiver  int    `json:"rating_by_receiver"`  // 4,
	CommentByReceiver string `json:"comment_by_receiver"` // "Good",

	CompleteTime    float32 `json:"complete_time"`    // 1426671164,
	CompleteLat     float32 `json:"complete_lat"`     // 10.7890462,
	CompleteLng     float32 `json:"complete_lng"`     // 106.7763078,
	CompleteComment string  `json:"complete_comment"` // "Nice receiver",
	ImageUrl        string  `json:"image_url"`        // "//i.imgur.com/lchC2xz.jpg", # Image to show that the supplier has already completed the order, uploaded by the supplier
	PodInfo         string  `json:"pod_info"`         // "024792155", # Proof of delivery information that supplier has collected, can be recipient's ID number or image URL
	Status          string  `json:"status"`           // "COMPLETED"
}

type CalcShippingFeeResponse struct {
	Distance         Float  `json:"distance"`           // 1.144
	Duration         Int    `json:"duration"`           // 288
	Currency         string `json:"currency"`           // "VND"
	Discount         Int    `json:"discount"`           // 30000
	DistanceFee      Int    `json:"distance_fee"`       // 140000,
	StopFee          Int    `json:"stop_fee"`           // 0,
	TotalFee         Int    `json:"total_fee"`          // 275000,  # Total = Distance Fee + Request Fee + Stop Fee
	UserMainAccount  Int    `json:"user_main_account"`  // 0,
	UserBonusAccount Int    `json:"user_bonus_account"` // 0,
	TotalPay         Int    `json:"total_pay"`          // 275000,
	// Surcharge           string  `json:"surcharge"`             // 1.1 # 110% compared to normal fee, only returned if surcharge > 1,
	BookingTimeError    string `json:"booking_time_error"`    // "Booking time not valid" #If booking time is not in opening hours period,
	PaymentErrorMessage string `json:"payment_error_message"` // "Not enough credit" #If user does not have enough credit in balance
	PartnerFee          Int    `json:"partner_fee"`
	PartnerDistanceFee  Int    `json:"partner_distance_fee"`
	PartnerDiscount     Int    `json:"partner_discount"`
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
	// Method which user chooses to pay for this order (Available methods: CASH)
	PaymentMethod string `url:"payment_method"`
	// Products      []Item `url:"-"`
	// Items         string `url:"items"`
}

type Item struct {
	ID    string `json:"_id"`   // "TS",
	Num   int    `json:"num"`   // 2,
	Name  string `json:"name"`  // "Sua tuoi",
	Price int    `json:"price"` // 15000
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
	Duration Int              `json:"duration"` // 5015,
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
	OrderTime      Float `json:"order_time"`      // 1426297774, # Pick-up time
	CreateTime     Float `json:"create_time"`     // 1426297174.649361, # Order create time
	AcceptTime     Float `json:"accept_time"`     // 1426298634.189912, # Order accept time
	AcceptLat      Float `json:"accept_lat"`      // 10.7890462,
	AcceptLng      Float `json:"accept_lng"`      // 106.7763078,
	AcceptDistance Float `json:"accept_distance"` // 1.2, # Distance from accepted location to pick up point
	AcceptDuration Int   `json:"accept_duration"` // 150, # Time from accepted location to pick up point

	CancelTime    Float  `json:"cancel_time"`    // 1426671164.374609, # Order cancel time
	CancelComment string `json:"cancel_comment"` // "Lich chuyen nha bi hoan", # Cancelled reason by user or supplier
	CancelByUser  bool   `json:"cancel_by_user"` // False, # True if this order is cancelled by user, otherwise False
	// Others: complete_time, complete_lat, complete_lng, accept_lat, accept_lng, fail_tome, fail_lat, fail_lng, fail_comment

	// Payment
	Currency    string `json:"currency"`     // "VND",
	StopFee     Int    `json:"stop_fee"`     // 0,
	RequestFee  Int    `json:"request_fee"`  // 0,
	Distance    Float  `json:"distance"`     // 5.597,
	DistanceFee Int    `json:"distance_fee"` // 142358,
	PromoCode   string `json:"promo_code"`   // "AHAMOVE",
	Discount    Int    `json:"discount"`     // 0,
	TotalFee    Int    `json:"total_fee"`    // 142358, # Total fee = Distance Fee + Request Fee + Stop Fee - Discount
	VatFee      Int    `json:"vat_fee"`      // 0,

	// Use credit from user account if available
	UserBonusAccount Int `json:"user_bonus_account"` // 0,
	UserMainAccount  Int `json:"user_main_account"`  // 0,
	TotalPay         Int `json:"total_pay"`          // 142358, # Total pay = Total fee - User Main account - User Bonus account

	// Rating
	RatingByUser       Int    `json:"rating_by_user"`        // 4,
	CommentByUser      string `json:"comment_by_user"`       // "Good",
	RatingBySupplier   Int    `json:"rating_by_supplier"`    // 5,
	CommentBySupplier  string `json:"comment_by_supplier"`   // "Great customer",
	RatingByReceiver   Int    `json:"rating_by_receiver"`    // 4,
	CommentByReceiver  string `json:"comment_by_receiver"`   // "Good driver",
	StoreRatingByUser  Int    `json:"store_rating_by_user"`  // 4, # For food orders
	StoreCommentByUser string `json:"store_comment_by_user"` // "Delicious food", # For food orders

	// Others
	Remarks    string `json:"remarks"`     // "Den noi goi dien cho toi", # Note to supplier
	Remind     bool   `json:"remind"`      // True, # If this is advance booking and the system already reminds users
	AssignedBy string `json:"assigned_by"` // "84908842285", # If this order is assigned by a user or "auto" if by the system
	Index      Int    `json:"index"`       // 1, # 0 if first order request from the user, 1 for second...
	// Items: chỉ sử dụng cho food, tài xế dựa vào đây để mua hàng
	// Items      []Item `json:"items"`
}

type GetOrderRequest struct {
	Token   string `url:"token"`
	OrderID string `url:"order_id"`
}

type CancelOrderRequest struct {
	Token    string `url:"token"`     // ApiKey (User token or Supplier token, Admin if want to cancel a completed order)
	OrderId  string `url:"order_id"`  // Order ID
	Comment  string `url:"comment"`   // Optional
	ImageUrl string `url:"image_url"` // Optional, image URL to show the reason why supplier cancel an order
}

type RegisterAccountRequest struct {
	ApiKey  string `url:"api_key"`
	Mobile  string `url:"mobile"`
	Name    string `url:"name"`
	Address string `url:"address"`
}

type RegisterAccountResponse struct {
	Token string `json:"token"`
}

type GetAccountRequest struct {
	Token string `url:"token"`
}

type VerifyAccountRequest struct {
	Token       string `url:"token"`
	Description string `url:"description"`
	Subject     string `url:"subject"`
	Type        string `url:"type"`
}

type VerifyAccountResponse struct {
	Ticket Ticket `json:"ticket"`
}

type Ticket struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	Subject     string `json:"subject"`
}

type Account struct {
	ID            string   `json:"_id"`
	Name          string   `json:"name"`
	Partner       string   `json:"partner"`
	Email         string   `json:"email"`
	ReferralCode  string   `json:"referral_code"`
	CountryCode   string   `json:"country_code"`
	Currency      string   `json:"currency"`
	AccountStatus string   `json:"account_status"`
	CreateTime    Float    `json:"create_time"`
	LastActivity  Float    `json:"last_activity"`
	Points        Float    `json:"points"`
	BannedList    []string `json:"banned_list"`
	NumRating     Float    `json:"num_rating"`
	Rating        Float    `json:"rating"`
	Verified      bool     `json:"verified"`
}

type GetServicesRequest struct {
	CityID string `url:"city_id"`
}

type ServiceType struct {
	ID              string `json:"_id"`          // "SGN-TRICYCLE"
	Name            string `json:"name"`         // "Xe ba gac (1.8 x 1.3m)"
	NameViVn        string `json:"name_vi_vn"`   // "Xe ba gac (1.8 x 1.3m)"
	Currency        string `json:"currency"`     // "VND"
	CityID          string `json:"city_id"`      // "SGN"
	IconUrl         string `json:"icon_url"`     // "http//apistg.ahamove.com/images/tricycle.png"
	DistanceFee     string `json:"distance_fee"` // "120000 if x <= 4 else 120000 + (x - 4)  14000",
	StopFee         Int    `json:"stop_fee"`     // 10000
	MinStopPoints   Int    `json:"min_stop_points"`
	MaxStopPoints   Int    `json:"max_stop_points"`   // 5
	DescriptionViVn string `json:"description_vi_vn"` // "Giao hàng trong 1h"
	MaxCOD          Int    `json:"max_cod"`
	COD             Int    `json:"cod"`
}
