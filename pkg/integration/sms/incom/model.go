package incom

import "encoding/json"

const (
	ByBrandNameID     FeeTypeID        = "0"
	ByPrefixID        FeeTypeID        = "1"
	UnAccentContentID MsgContentTypeID = "0"
	AccentContentID   MsgContentTypeID = "12"
)

type SendSMSRequest struct {
	Submission Submission `json:"submission"`
}

type Submission struct {
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Sms       []SMS  `json:"sms"`
}

type SMS struct {
	ID        string `json:"id"`
	BrandName string `json:"brandname"`
	Text      string `json:"text"`
	To        string `json:"to"`
	// 0: Gửi bằng brandname, 1: Gửi bằng đầu số
	FeeTypeID FeeTypeID `json:"feetypeid"`
	// 0: Gửi tin không dấu, 12: Gửi tin có dấu
	MSGContentTypeID MsgContentTypeID `json:"msgcontenttypeid"`
}

type FeeTypeID string

type MsgContentTypeID string

type CommonResponse struct {
	Response *SendSMSResponse `json:"response"`
}
type SendSMSResponse struct {
	Submission *SubmissionRes `json:"submission"`
}

type SubmissionRes struct {
	Sms []json.RawMessage `json:"sms"`
}

type SMSRes struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

var SMSResultCodeMap = map[string]string{
	"0":  "Gửi thành công",
	"2":  "Lỗi hệ thống",
	"3":  "Sai user hoặc mật khẩu",
	"4":  "IP không được phép",
	"5":  "Chưa khai báo brandname/dịch vụ",
	"6":  "Lặp nội dung",
	"7":  "Thuê bao từ chối nhận tin",
	"8":  "Không được phép gửi tin",
	"9":  "Chưa khai báo template",
	"10": "Định dạng thuê bao không đúng",
	"11": "Có tham số không hợp lệ",
	"12": "Tài khoản không đúng",
	"13": "Gửi tin: lỗi kết nối",
	"14": "Tài khoản không đủ",
	"15": "Tài khoản hết hạn",
	"16": "Hết hạn dịch vụ",
	"17": "Hết hạn mức gửi test",
	"18": "Hủy gửi tin (CSKH)",
	"19": "Hủy gửi tin (KD)",
	"20": "Gateway chưa hỗ trợ Unicode",
	"21": "Chưa set giá trả trước",
	"22": "Tài khoản chưa kích hoạt",
	"25": "Chưa khai báo partner cho user",
	"26": "Chưa khai báo GateOwner cho user",
	"27": "Gửi tin: gate trả mã lỗi",
	"31": "Số điện thoại 11 số đã chuyển sang 10 số",
}
