package client

import (
	"strings"

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

const MessageVersion = "1.4"
const MessageLanguage = "vi_VN"

type ReturnMode string

const (
	ReturnToRegisteredAddress ReturnMode = "01"
	ReturnToPickupAddress     ReturnMode = "02" // (only available with Ad-Hoc Pickup)
	ReturnToNewAddress        ReturnMode = "03"
	Abandon                   ReturnMode = "05"
)

type State string

const (
	// Submitted: 71005
	// Đã tạo đơn hàng thành công
	StateSubmitted State = "71005"

	// Shipment data received: 77123 - 99
	// ĐÃ NHẬN DỮ LIỆU ĐƠN HÀNG
	// DHL đạ ghi nhận dự liệu đơn hàng
	StateShipmentDataReceived State = "77123"

	// Shipment data edited: 77305 - 163
	// DỮ LIỆU ĐƠN HÀNG ĐÃ ĐƯỢC CHỈNH SỬA
	// Cập nhật sau khi đơn hàng đã được chỉnh sửa
	// -- ignore
	StateShipmentDataEdited State = "77305"

	// Shipment cancelled by customer: 77272 - 101
	// ĐƠN HÀNG HỦY THEO YÊU CẦU KHÁCH HÀNG
	StateShipmentCancelledByCustomer State = "77272"

	// Shipment picked up: 77206 - 130
	// ĐÃ LẤY HÀNG THÀNH CÔNG
	StateShipmentPickedUp State = "77206"

	// Shipment picked up failed: 77429 - 134
	// LẤY HÀNG KHÔNG THÀNH CÔNG
	StateShipmentPickedUpFailed State = "77429"

	// Shipment has arrived at service point: 77236 - 116
	// HÀNG ĐÃ ĐẾN ĐIỂM DỊCH VỤ
	StateShipmentHasArrivedAtServicePoint State = "77236"

	// En-route to facility: 77208 - 132
	// TRUNG CHUYỂN VỀ TRẠM
	// Hàng từ điểm dịch vụ được trung chuyển về Hub của DHL
	StateEnRouteToFacility State = "77208"

	// Processed at facility: 77015 - 220
	// ĐANG XỬ LÝ TẠI TRẠM
	// Hàng được qua xử lý tại Hub DHL và đơn hàng bắt đầu được tính cước
	StateProcessedAtFacility State = "77015"

	// Departed from Facility: 77169 - 399
	// ĐÃ KHỞI HÀNH TỪ TRẠM
	// Hàng xuất phát từ Hub DHL
	StateDepartedFromFacility77169 State = "77169"

	// Arrived at facility: 77178 - 470
	// ĐÃ ĐẾN TRẠM
	// Hàng đã đến trạm khu vực người nhận
	StateArrivedAtFacility State = "77178"

	// Processed at delivery facility: 77184 - 543/557
	// ĐANG XỬ LÝ TẠI TRẠM
	// Hàng đang xử lý tại trạm khu vực người nhận
	StateProcessedAtDeliveryFacility State = "77184"

	// Arrival at destination facility: 77066 - 470
	// ĐÃ ĐẾN TRẠM TẠI QUỐC GIA NGƯỜI NHẬN (DSP)
	// Hàng đến trạm phát của DSP
	StateArrivalAtDestinationFacility State = "77066"

	// Arrival at facility: 77177 - 468
	// ĐÃ ĐẾN TRẠM DSP
	// Hàng được chuyển dến DSP
	StateArrivalAtFacility State = "77177"

	// Departed from facility: 77180 - 481
	// ĐÃ KHỞI HÀNH TỪ TRẠM DSP
	StateDepartedFromFacility77180 State = "77180"

	// Out for Delivery: 77090 - 597
	// HÀNG ĐANG ĐƯỢC ĐI GIAO
	StateOutForDelivery State = "77090"

	// Delivery was refused: 77098 - 611
	// BỊ TỪ CHỐI GIAO HÀNG
	// Hàng sẽ được chuyển hoàn sau trạng thái này
	StateDeliveryWasRefused State = "77098"

	// Attempted delivery: closed on arrival: 77170 - 612
	// GIAO HÀNG KHÔNG THÀNH CÔNG: KHU VỰC ĐÓNG CỬA
	StateAttemptedDeliveryClosedOnArrival State = "77170"

	// Attempted delivery: no one at home: 77171 - 613
	// GIAO HÀNG KHÔNG THÀNH CÔNG: KHÔNG CÓ NGƯỜI Ở NHÀ
	StateAttemptedDeliveryNoOneAtHome State = "77171"

	// Shipment on hold for payment: 77172 - 614
	// HÀNG ĐANG CHỜ THANH TOÁN
	StateHoldForPayment State = "77172"

	// Delivery was attempted: 77173 - 615
	// GIAO HÀNG KHÔNG THÀNH CÔNG
	StateAttemptedDelivery State = "77173"

	// Delivery rescheduled: 77207 - 616
	// THAY ĐỔI THỜI GIAN GIAO HÀNG
	StateDeliveryRescheduled State = "77207"

	// Reject after opening: 77234 - 619
	// TỪ CHỐI NHẬN SAU KHI XEM HÀNG
	// Hàng sẽ được chuyển hoàn sau trạng thái này
	StateRejectAfterOpening State = "77234"

	// Attempted delivery: insufficient address: 77101 - 631
	// GIAO HÀNG KHÔNG THÀNH CÔNG: ĐỊA CHỈ KHÔNG CHÍNH XÁC
	// Hàng sẽ được chuyển hoàn sau trạng thái này
	StateAttemptedDeliveryInsufficientAddress State = "77101"

	// Shipment redirected to new address: 77102 - 632
	// HÀNG ĐƯỢC CHUYỂN ĐẾN ĐỊA CHỈ MỚI
	StateRedirectedToNewAddress State = "77102"

	// Delivery was attempted - restricted address: 621 - 621
	// GIAO HÀNG KHÔNG THÀNH CÔNG: KHU VỰC CẤM
	StateDeliveryWasAttemptedRestrictedAddress State = "621"

	// Delivery was attempted - unable to deliver as instructed
	StateDeliveryWasAttemptedUnableToDeliverAsInstructed State = "622"

	// Successfully delivered: 77093 - 600
	// HÀNG ĐÃ GIAO THÀNH CÔNG
	StateSuccessfullyDelivered State = "77093"

	// COD Amount has been deposited: 77223 - 618
	// TIỀN THU HỘ ĐÃ ĐƯỢC THU
	StateCODAmountHasBeenDeposited State = "77223"

	// COD Amount has been remitted: 77222 - 617
	// TIỀN THU HỘ ĐÃ ĐƯỢC CHUYỂN KHOẢN
	StateCODAmountHasBeenRemitted State = "77222"

	// Shipment on hod for return or disposal: 77117 - 800
	// HÀNG ĐANG CHỜ TRẢ VỀ HOẶC HỦY
	// Bắt đầu quá trình chuyển hoàn
	StateShipmentOnHodForReturnOrDisposal State = "77117"

	// Return shipment has departed facility: 77300 - 817
	// ĐƠN HÀNG HOÀN ĐÃ RỜI KHỎI TRẠM
	// Hàng hoàn chuyển từ Depot về Hub (Optional)
	StateReturnShipmentHasDepartedFacility State = "77300"

	// Return initiated: 77250 - 808
	// BẮT ĐẦU QUI TRÌNH TRẢ HÀNG
	// *Chỉ áp dụng returnMode 02 và 03
	StateReturnInitiated State = "77250"

	// Return shipment arrived at facility: 77298 - 810
	// ĐƠN HÀNG HOÀN ĐÃ ĐẾN TRẠM TRUNG CHUYỂN
	// Hàng hoàn về đến Hub (Optional)
	StateReturnShipmentArrivedAtFacility State = "77298"

	// Return shipment being processed: 77239 - 815
	// ĐANG XỬ LÝ TRẢ HÀNG
	// Hàng được xử lý tại Hub DHL
	StateReturnShipmentBeingProcessed State = "77239"

	// Return sorted to delivery facility: 77307 - 856
	// PHÂN LOẠI ĐƠN HÀNG HOÀN TẠI TRẠM GIAO
	// *Chỉ áp dụng returnMode 02 và 03
	StateReturnSortedToDeliveryFacility State = "77307"

	// Return processed at delivery facility: 77467 - 858
	// XỬ LÝ ĐƠN HÀNG HOÀN TẠI TRẠM GIAO
	// *Chỉ áp dụng returnMode 02 và 03
	StateReturnProcessedAtDeliveryFacility State = "77467"

	// Return shipment Out for Delivery 77232 - 820
	// HÀNG ĐANG ĐƯỢC HOÀN TRẢ
	StateReturnShipmentOutForDelivery State = "77232"

	// Returns shipment was successfully delivered: 77174 - 860
	// ĐÃ TRẢ HÀNG THÀNH CÔNG
	StateReturnsDelivered State = "77174"

	// Shipment abandoned by customer: 77273 - 780
	// ĐƠN HÀNG TIÊU HỦY
	// Hàng được hủy theo yêu cầu từ khách hàng
	StateShipmentAbandonedByCustomer State = "77273"

	// Item reported missing: 77243 - 690
	StateItemReportedMissing State = "77243"

	// Paper POD being processed: 77240 - 870
	// -- ignore
	StatePaperPODBeingProcessed State = "77240"

	// Paper POD has been successfully delivered: 77241 - 871
	// -- ignore
	StatePaperPODHasBeenSuccessfullyDelivered State = "77241"

	// Shipment was not delivered: 77270 - 625
	StateShipmentWasNotDelivered State = "77270"

	// Return delivery was attempted - recipient moved: 77421 - 885
	StateReturnDeliveryWasAtTemptedRecipientMoved State = "77421"

	// Return delivery was attempted - no one at home: 77423 - 886
	StateReturnDeliveryWasAttemptedNoOneAtHome State = "77423"

	// Return delivery was attempted - closed premise: 77425 - 887
	StateReturnDeliveryWasAttemptedClosedPremise State = "77425"

	// Return delivery rescheduled: 77427 - 880
	StateReturnDeliveryRescheduled State = "77427"

	// Return delivery refused: 77411 - 881
	StateReturnDeliveryRefused State = "77411"

	// Return shipment refused after, opening: 77419 - 884
	StateReturnShipmentRefusedAfterAndOpening State = "77419"

	// Return shipment pickup was unsuccessfull: 77224 - 627
	StateReturnShipmentPickupWasUnsuccessful State = "77224"

	// Returned to customer: 77286 - 863
	StateReturnedToCustomer State = "77286"
)

var mapState = map[string]State{
	"71005": StateSubmitted,
	"77123": StateShipmentDataReceived,
	"77305": StateShipmentDataEdited,
	"77272": StateShipmentCancelledByCustomer,
	"77206": StateShipmentPickedUp,
	"77429": StateShipmentPickedUpFailed,
	"77236": StateShipmentHasArrivedAtServicePoint,
	"77208": StateEnRouteToFacility,
	"77015": StateProcessedAtFacility,
	"77169": StateDepartedFromFacility77169,
	"77178": StateArrivedAtFacility,
	"77184": StateProcessedAtDeliveryFacility,
	"77066": StateArrivalAtDestinationFacility,
	"77177": StateArrivalAtFacility,
	"77180": StateDepartedFromFacility77180,
	"77090": StateOutForDelivery,
	"77098": StateDeliveryWasRefused,
	"77170": StateAttemptedDeliveryClosedOnArrival,
	"77171": StateAttemptedDeliveryNoOneAtHome,
	"77172": StateHoldForPayment,
	"77173": StateAttemptedDelivery,
	"77207": StateDeliveryRescheduled,
	"77234": StateRejectAfterOpening,
	"77101": StateAttemptedDeliveryInsufficientAddress,
	"77102": StateRedirectedToNewAddress,
	"621":   StateDeliveryWasAttemptedRestrictedAddress,
	"622":   StateDeliveryWasAttemptedUnableToDeliverAsInstructed,
	"77093": StateSuccessfullyDelivered,
	"77223": StateCODAmountHasBeenDeposited,
	"77222": StateCODAmountHasBeenRemitted,
	"77117": StateShipmentOnHodForReturnOrDisposal,
	"77300": StateReturnShipmentHasDepartedFacility,
	"77250": StateReturnInitiated,
	"77298": StateReturnShipmentArrivedAtFacility,
	"77239": StateReturnShipmentBeingProcessed,
	"77307": StateReturnSortedToDeliveryFacility,
	"77467": StateReturnProcessedAtDeliveryFacility,
	"77232": StateReturnShipmentOutForDelivery,
	"77174": StateReturnsDelivered,
	"77273": StateShipmentAbandonedByCustomer,
	"77243": StateItemReportedMissing,
	"77240": StatePaperPODBeingProcessed,
	"77241": StatePaperPODHasBeenSuccessfullyDelivered,
	"77270": StateShipmentWasNotDelivered,
	"77421": StateReturnDeliveryWasAtTemptedRecipientMoved,
	"77423": StateReturnDeliveryWasAttemptedNoOneAtHome,
	"77425": StateReturnDeliveryWasAttemptedClosedPremise,
	"77427": StateReturnDeliveryRescheduled,
	"77411": StateReturnDeliveryRefused,
	"77419": StateReturnShipmentRefusedAfterAndOpening,
	"77224": StateReturnShipmentPickupWasUnsuccessful,
	"77286": StateReturnedToCustomer,
}

var StatusMapping = map[string]string{
	"71005": "Submitted",
	"77123": "Shipment data received",
	"77305": "Shipment data edited",
	"77272": "Shipment cancelled by customer",
	"77206": "Shipment picked up",
	"77429": "Shipment picked up failed",
	"77236": "Shipment has arrived at service point",
	"77208": "En-route to facility",
	"77015": "Processed at facility",
	"77169": "Departed from facility",
	"77178": "Arrived at facility",
	"77184": "Processed at delivery facility",
	"77066": "Arrival a destination facility",
	"77177": "Arrival at facility",
	"77180": "Departed from facility",
	"77090": "Out for delivery",
	"77098": "Delivery was refused",
	"77170": "Delivery was attempted: closed premises",
	"77171": "Delivery was attempted: no one at home",
	"77172": "Shipment on hold for payment",
	"77173": "Delivery was attempted",
	"77207": "Delivery rescheduled",
	"77234": "Reject after opening",
	"77101": "Delivery was attempted: insufficient address",
	"77102": "Shipment redirected to new address",
	"621":   "Delivery was attempted: restricted access",
	"622":   "Delivery was attempted: unable to deliver as instructed",
	"77093": "Successfully delivered",
	"77223": "COD amount has been deposited",
	"77222": "COD amount has been remitted",
	"77117": "Shipment on hold for return or disposal",
	"77300": "Return shipment has departed facility",
	"77250": "Return initiated",
	"77298": "Return shipment arrived at facility",
	"77239": "Return shipment being processed",
	"77307": "Return sorted to delivery facility",
	"77467": "Return processed at delivery facility",
	"77232": "Return shipment out for delivery",
	"77174": "Return shipment was successfully delivered",
	"77273": "Shipment abandoned by customer",
	"77243": "Item reported missing",
	"77240": "Paper POD being processed",
	"77241": "Paper POD has been successfully delivered",
}

var SecondaryStatusMapping = map[string]string{
	"77327": "CUSTOMER REFUSED SHIPMENT DUE TO DELAY IN DELIVERY",
	"77329": "CUSTOMER REJECTED THE SHIPMENT",
	"77331": "CUSTOMER REFUSED SHIPMENT DUE TO DAMAGE IN TRANSIT",
	"77333": "THE COD AMOUNT MAY BE INCORRECT",
	"77335": "CUSTOMER HAS REQUESTED TO CANCEL THE ORDER",
	"77337": "CUSTOMER HAS ALREADY RECEIVED THE SHIPMENT",
	"77339": "UNABLE TO OPEN BOX",
	"77341": "REFUSED COD PAYMENT",
	"77343": "RECIPIENT NO-SHOW",
	"77345": "DUE TO COURIER WORKLOAD",
	"77347": "DUE TO HEAVY TRAFFIC",
	"77349": "DUE TO BAD WEATHER",
	"77351": "SHIPMENT WITH INCORRECT COURIER",
	"77353": "MISSING CUSTOMER ID WITH SHIPMENT",
	"77355": "DUE TO VEHICLE BREAKDOWN",
	"77357": "REQUEST FROM DHL TO ABORT DELIVERY ATTEMPT",
	"77359": "DUE TO STRIKE AND/OR ROADBLOCKS",
	"77361": "SHIPMENT DAMAGED IN TRANSIT",
	"77363": "COURIER CANNOT FIND THE LOCATION",
	"77365": "SHIPMENT LOST IN TRANSIT",
	"77367": "NO SAFE PLACE TO LEAVE THE SHIPMENT",
	"77369": "THE SHIPMENT IS TOO LARGE FOR THE MAILBOX",
	"77371": "UNABLE TO ENTER PREMISES",
	"77373": "NEIGHBOUR IS UNAVAILABLE",
	"77375": "WRONG CONTENT IN SHIPMENT",
	"77377": "SHIPMENT IS DAMAGED",
	"77379": "UNABLE TO TEST CONTENT",
	"77381": "CUSTOMER REFUSED SHIPMENT DUE TO MISSING CONTENT",
	"77383": "THE CUSTOMER HAS MOVED",
	"77385": "THE GUEST HAS CHECKED OUT FROM THE HOTEL",
	"77387": "CONSIGNEE REQUEST TO CHANGE THE NEW ADDRESS",
	"77389": "THE CUSTOMER IS NOT AT HOME",
	"77391": "THE GUEST HAS NOT ARRIVED AT THE HOTEL",
	"77393": "THE BUSINESS ADDRESS IS CLOSED",
	"77395": "THE CUSTOMER IS ON HOLIDAY",
	"77397": "TOO LATE - AFTER BUSINESS HOURS",
	"77399": "TOO EARLY - BEFORE BUSINESS HOURS",
	"77401": "THE CUSTOMER IS OUT FOR LUNCH",
	"77403": "THE CUSTOMER HAS RESCHEDULED DUE TO NOT AT HOME",
	"77405": "THE CUSTOMER HAS RESCHEDULED DUE TO COD AMOUNT NOT BEING AVAILABLE",
	"77407": "CUSTOMER HAS RESCHEDULED DUE TO NOT READY",
	"77309": "INCORRECT ADDRESS AND/OR PHONE NUMBER",
	"77311": "CUSTOMER NOT REACHABLE TO VERIFY ADDRESS",
	"77313": "PROVIDED CONTACT NUMBER IS INACTIVE",
	"77315": "PROVIDED CONTACT NUMBER DOES NOT EXIST",
	"77317": "PO BOX IS MISSING ADDRESS AND CONTACT DETAILS",
	"77319": "INCORRECT ZIPCODE",
	"77321": "NO SUCH ADDRESS",
	"77323": "NO SUCH CUSTOMER AT PROVIDED ADDRESS",
	"77325": "CUSTOMER REFUSED SHIPMENT",
	"77431": "Cancelled by Customer Service",
	"77433": "Parcels are oversize",
	"77435": "Pickup location is closed on arrival",
	"77437": "No packages to be picked up",
	"77441": "Capacity over limits",
	"77443": "Insufficient pickup address",
	"77445": "Closed on the public holiday",
	"77447": "Inadequate packaging",
	"77245": "Item Missing",
	"77246": "Rejected by customs",
	"77247": "Return to customer",
	"77248": "DG Reject",
	"77439": "Attemped after business hours",
	"77449": "Mechant moved to new address",
	"77451": "Parcel is not ready to be picked up",
	"77453": "Pickup rescheduled",
	"77455": "Duplicate booking",
	"77457": "Consumer moved to new address",
	"77459": "No one at home",
	"77461": "Pickup rescheduled",
	"77463": "Pickup was refused",
	"77288": "ITEM DAMAGED",
	"77290": "OTHER REASONS",
	"77292": "PROHIBITED ITEMS, CONFISCATED OR DELAYED",
	"77294": "SECURITY CHECK FAILURE",
}

func ToState(status string) State {
	state, ok := mapState[status]
	if !ok {
		return State(status)
	}
	return state
}

func (s State) ToModel() typesshipping.State {
	switch s {
	case StateSubmitted,
		StateShipmentDataReceived:
		return typesshipping.Created
	case StateShipmentPickedUp,
		StateShipmentPickedUpFailed:
		return typesshipping.Picking
	case
		StateEnRouteToFacility,
		StateProcessedAtFacility,
		StateDepartedFromFacility77169,
		StateArrivalAtDestinationFacility,
		StateProcessedAtDeliveryFacility,
		StateArrivalAtFacility,
		StateArrivedAtFacility,
		StateShipmentHasArrivedAtServicePoint,
		StateDepartedFromFacility77180:
		return typesshipping.Holding
	case StateOutForDelivery,
		StateRejectAfterOpening,
		StateDeliveryWasRefused,
		StateAttemptedDeliveryInsufficientAddress,
		StateAttemptedDeliveryClosedOnArrival,
		StateAttemptedDeliveryNoOneAtHome,
		StateHoldForPayment,
		StateAttemptedDelivery,
		StateRedirectedToNewAddress,
		StateDeliveryRescheduled,
		StateDeliveryWasAttemptedRestrictedAddress,
		StateDeliveryWasAttemptedUnableToDeliverAsInstructed,
		StateShipmentWasNotDelivered:
		return typesshipping.Delivering
	case StateSuccessfullyDelivered,
		StateCODAmountHasBeenDeposited,
		StateCODAmountHasBeenRemitted:
		return typesshipping.Delivered
	case StateShipmentOnHodForReturnOrDisposal,
		StateReturnShipmentOutForDelivery,
		StateReturnProcessedAtDeliveryFacility,
		StateReturnShipmentArrivedAtFacility,
		StateReturnShipmentBeingProcessed,
		StateReturnShipmentHasDepartedFacility,
		StateReturnSortedToDeliveryFacility,
		StateReturnInitiated,
		StateReturnDeliveryWasAttemptedNoOneAtHome,
		StateReturnDeliveryWasAtTemptedRecipientMoved,
		StateReturnDeliveryWasAttemptedClosedPremise,
		StateReturnDeliveryRescheduled,
		StateReturnDeliveryRefused,
		StateReturnShipmentRefusedAfterAndOpening,
		StateReturnShipmentPickupWasUnsuccessful:
		return typesshipping.Returning
	case StateReturnsDelivered,
		StateReturnedToCustomer:
		return typesshipping.Returned
	case StateShipmentAbandonedByCustomer,
		StateShipmentCancelledByCustomer:
		return typesshipping.Cancelled
	case StateItemReportedMissing:
		return typesshipping.Undeliverable
	default:
		return typesshipping.Unknown
	}
}

func (s State) ToSubstateModel() shippingsubstate.Substate {
	switch s {
	case StateShipmentPickedUpFailed:
		return shippingsubstate.PickFail
	case StateRejectAfterOpening,
		StateDeliveryWasRefused,
		StateAttemptedDeliveryInsufficientAddress,
		StateAttemptedDeliveryNoOneAtHome,
		StateAttemptedDelivery,
		StateDeliveryRescheduled,
		StateAttemptedDeliveryClosedOnArrival,
		StateHoldForPayment,
		StateDeliveryWasAttemptedRestrictedAddress,
		StateRedirectedToNewAddress,
		StateDeliveryWasAttemptedUnableToDeliverAsInstructed,
		StateShipmentWasNotDelivered:
		return shippingsubstate.DeliveryFail
	case StateReturnDeliveryWasAtTemptedRecipientMoved,
		StateReturnDeliveryWasAttemptedNoOneAtHome,
		StateReturnDeliveryWasAttemptedClosedPremise,
		StateReturnDeliveryRescheduled,
		StateReturnDeliveryRefused,
		StateReturnShipmentRefusedAfterAndOpening,
		StateReturnShipmentPickupWasUnsuccessful:
		return shippingsubstate.ReturnFail
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

type ErrorResponse struct {
	Code           String   `json:"code"`
	Message        String   `json:"message"`
	MessageDetails []String `json:"message_details"`
}

func (e ErrorResponse) Error() (s string) {
	defer func() {
		s = strs.TrimLastPunctuation(s)
	}()

	b := strings.Builder{}
	b.WriteString(e.Message.String())
	b.WriteString(" (")
	for i, v := range e.MessageDetails {
		b.WriteString(v.String())
		if i < len(e.MessageDetails)-1 {
			b.WriteString(". ")
		}
	}
	b.WriteString(")")
	return b.String()
}

// DHLAccountCfg được lấy ra từ shopconn.ExternalData
// ClientID ~ UserID
// AccountID ~ ShopID
type DHLAccountCfg struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	AccountID    string `yaml:"account_id"`
	Token        string `yaml:"token"`
}

type GenerateAccessTokenRequest struct {
	ClientID string
	Password string
}

type GenerateAccessTokenResponse struct {
	AccessTokenResponse *AccessTokenResponse `json:"accessTokenResponse"`
}

type AccessTokenResponse struct {
	Token            String                             `json:"token"`
	TokenType        String                             `json:"token_type"`
	ExpiresInSeconds Int                                `json:"expires_in_seconds"`
	ClientID         String                             `json:"client_id"`
	ResponseStatus   *GenerateAccessTokenResponseStatus `json:"responseStatus"`
}

type GenerateAccessTokenResponseStatus struct {
	Code           String `json:"code"`
	Message        String `json:"message"`
	MessageDetails String `json:"messageDetails"`
}

func (e *GenerateAccessTokenResponseStatus) ToError() *ErrorResponse {
	return &ErrorResponse{
		Code:           e.Code,
		Message:        e.Message,
		MessageDetails: []String{e.MessageDetails},
	}
}

type FindAvailableServicesResponse struct {
	AvailableServices []*AvailableService
}

type AvailableService struct {
	Name        String
	ServiceCode String
}

type CreateOrdersRequest struct {
	ManifestRequest *ManiFestReq `json:"manifestRequest"`
}

type ManiFestReq struct {
	Hdr *HdrReq `json:"hdr"`
	Bd  *BdReq  `json:"bd"`
}

type HdrReq struct {
	MessageType     string `json:"messageType"`
	MessageDateTime string `json:"messageDateTime"`
	MessageVersion  string `json:"messageVersion"`
	AccessToken     string `json:"accessToken"`
	MessageLanguage string `json:"messageLanguage"`
}

type BdReq struct {
	PickupAccountID string             `json:"pickupAccountId"`
	SoldToAccountID string             `json:"soldToAccountId"`
	PickupAddress   *AddressReq        `json:"pickupAddress"`
	ShipmentItems   []*ShipmentItemReq `json:"shipmentItems"`
}

type AddressReq struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	State    string `json:"state"`
	District string `json:"district"`
	Country  string `json:"country"`
}

func ToAddress(address *model.Address) *AddressReq {
	if address == nil {
		return nil
	}

	// Maximum province 20 characters
	// remove "Thành phố", "Tỉnh"
	province := address.Province
	if strings.HasPrefix(province, "Thành phố") {
		province = strings.TrimPrefix(province, "Thành phố")
	}
	if strings.HasPrefix(province, "Tỉnh") {
		province = strings.TrimPrefix(province, "Tỉnh")
	}

	return &AddressReq{
		Name:     address.FullName,
		Phone:    address.Phone,
		Address1: address.Address1,
		Address2: address.Address2,
		District: address.District,
		State:    province,
		Country:  "VN",
	}
}

type AddressResp struct {
	Name     String `json:"name"`
	Email    String `json:"email"`
	Phone    String `json:"phone"`
	Address1 String `json:"address1"`
	Address2 String `json:"address2"`
	City     String `json:"city"`
	District String `json:"district"`
	State    String `json:"state"`
	Country  String `json:"country"`
}

type ShipmentItemReq struct {
	ConsigneeAddress *AddressReq `json:"consigneeAddress"`
	ReturnAddress    *AddressReq `json:"returnAddress"`
	ShipmentID       string      `json:"shipmentID"`
	ReturnMode       string      `json:"returnMode"`
	PackageDesc      string      `json:"packageDesc"`
	TotalWeight      int         `json:"totalWeight"`
	TotalWeightUOM   string      `json:"totalWeightUOM"` // g
	Height           float64     `json:"height"`
	Length           float64     `json:"length"`
	Width            float64     `json:"width"`
	ProductCode      string      `json:"productCode"`
	CodValue         int         `json:"codValue"`
	InsuranceValue   float64     `json:"insuranceValue,omitempty"`

	// Doc:
	// Total declared value of the shipment (in 2 decimal points). Mandatory for Cross Border shipment, optional for Domestic shipment.
	// For Vietnam Domestic, totalValue must be a multiple of 500.
	TotalValue float64 `json:"totalValue"`
	Currency   string  `json:"currency"` // VND
	Remarks    string  `json:"remarks"`
}

type CreateOrdersResponse struct {
	ManifestResponse *ManifestResp `json:"manifestResponse"`
}

type ManifestResp struct {
	Hdr *HdrResp `json:"hdr"`
	Bd  *BdResp  `json:"bd"`
}

type HdrResp struct {
	MessageType     String `json:"messageType"`
	MessageDateTime Time   `json:"messageDateTime"`
	MessageVersion  String `json:"messageVersion"`
	MessageLanguage String `json:"messageLanguage"`
}

type BdResp struct {
	ShipmentItems  []*ShipmentItemResp `json:"shipmentItems"`
	ResponseStatus *ResponseStatus     `json:"responseStatus"`
}

type ShipmentItemResp struct {
	ShipmentID             String          `json:"shipmentID"`
	DeliveryConfirmationNo String          `json:"deliveryConfirmationNo"`
	ResponseStatus         *ResponseStatus `json:"responseStatus"`
}

type ResponseStatus struct {
	Code           String           `json:"code"`
	Message        String           `json:"message"`
	MessageDetails []*MessageDetail `json:"messageDetails"`
}

func (r *ResponseStatus) ToError() *ErrorResponse {
	var messageDetails []String
	for _, messageDetail := range r.MessageDetails {
		messageDetails = append(messageDetails, messageDetail.MessageDetail)
	}
	return &ErrorResponse{
		Code:           r.Code,
		Message:        r.Message,
		MessageDetails: messageDetails,
	}
}

type CancelOrderRequest struct {
	DeleteShipmentReq *DeleteShipmentReq `json:"deleteShipmentReq"`
}

type DeleteShipmentReq struct {
	Hdr *HdrReq              `json:"hdr"`
	Bd  *BdDeleteShipmentReq `json:"bd"`
}

type BdDeleteShipmentReq struct {
	PickupAccountId string                           `json:"pickupAccountId"`
	SoldToAccountId string                           `json:"soldToAccountId"`
	ShipmentItems   []*ShipmentItemDeleteShipmentReq `json:"shipmentItems"`
}

type ShipmentItemDeleteShipmentReq struct {
	ShipmentID string `json:"shipmentID"`
}

type MessageDetail struct {
	MessageDetail String `json:"messageDetail"`
}

type CancelOrderResponse struct {
	DeleteShipmentResp *DeleteShipmentResp `json:"deleteShipmentResp"`
}

type DeleteShipmentResp struct {
	Hdr *HdrResp              `json:"hdr"`
	Bd  *BdDeleteShipmentResp `json:"bd"`
}

type BdDeleteShipmentResp struct {
	ShipmentItems  []*ShipmentItemDeleteResp `json:"shipmentItems"`
	ResponseStatus *ResponseStatus           `json:"responseStatus"`
}

type ShipmentItemDeleteResp struct {
	ShipmentID     String          `json:"shipmentId"`
	ResponseStatus *ResponseStatus `json:"responseStatus"`
}

type TrackingOrdersRequest struct {
	TrackItemRequest *TrackItemReq `json:"trackItemRequest"`
}

type TrackItemReq struct {
	Hdr *HdrReq         `json:"hdr"`
	Bd  *BdTrackItemReq `json:"bd"`
}

type BdTrackItemReq struct {
	TrackingReferenceNumber []string `json:"trackingReferenceNumber"`
}

type TrackingOrdersResponse struct {
	TrackItemResponse *TrackItemResp `json:"trackItemResponse"`
}

type TrackItemResp struct {
	Hdr *HdrResp         `json:"hdr"`
	Bd  *BdTrackItemResp `json:"bd"`
}

type BdTrackItemResp struct {
	ShipmentItems  []*ShipmentItemTrackResp `json:"shipmentItems"`
	ResponseStatus *ResponseStatus          `json:"responseStatus"`
}

type ShipmentItemTrackResp struct {
	MasterShipmentID  String               `json:"masterShipmentId"`
	ShipmentID        String               `json:"shipmentID"`
	TrackingID        String               `json:"trackingID"`
	OrderNumber       String               `json:"orderNumber"`
	HandoverID        String               `json:"handoverID"`
	ShippingService   *ShippingServiceResp `json:"shippingService"`
	ConsigneeAddress  *AddressResp         `json:"consigneeAddress"`
	Weight            Int                  `json:"weight"`
	DimensionalWeight Int                  `json:"dimensionalWeight"`
	WeightUnit        String               `json:"weightUnit"`
	Events            []*EventResp         `json:"events"`
}

func (s ShipmentItemTrackResp) GetWeight() int {
	if s.Weight > s.DimensionalWeight {
		return s.Weight.Int()
	}
	return s.DimensionalWeight.Int()
}

func (s ShipmentItemTrackResp) GetLatestEvent() *EventResp {
	if len(s.Events) == 0 {
		return nil
	}

	latestTime := s.Events[0].DateTime.ToTime()
	latestID := 0

	for i := 1; i < len(s.Events); i++ {
		currEvent := s.Events[i]
		currEventTime := currEvent.DateTime.ToTime()
		if currEventTime.After(latestTime) {
			latestTime = currEventTime
			latestID = i
		}
	}

	return s.Events[latestID]

}

type EventResp struct {
	Status          String       `json:"status"`
	SecondaryStatus String       `json:"secondaryStatus"`
	SecondaryEvent  String       `json:"secondaryEvent"`
	Description     String       `json:"description"`
	DateTime        Time         `json:"dateTime"`
	Timezone        String       `json:"timezone"`
	Address         *AddressResp `json:"address"`
}

type ShippingServiceResp struct {
	ProductCode String `json:"productCode"`
	ProductName String `json:"productName"`
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}
