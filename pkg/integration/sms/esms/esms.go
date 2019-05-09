package esms

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/validate"
)

var (
	ll = l.New()

	client = &http.Client{
		Timeout: 10 * time.Second,
	}
)

type Config struct {
	BaseURL   string `yaml:"base_url"   valid:"required"`
	APIKey    string `yaml:"api_key"    valid:"required"`
	SecretKey string `yaml:"secret_key" valid:"required"`
	BrandName string `yaml:"brand_name" valid:"required"`
}

const (
	APISendSMS              = "/SendMultipleMessage_V4_get"
	APIGetBalance           = "/GetBalance/%s/%s"
	APIGetSendStatus        = "/GetSendStatus"
	APIGetSMSReceiverStatus = "/GetSmsReceiverStatus_get"

	CodeSuccess = "100"

	SMSTypeBrandAds      = 1
	SMSTypeBrandCustomer = 2
	SMSTypeRandom        = 3
	SMSTypeStatic        = 4
	SMSTypeVerify        = 6
	SMSTypeOTP           = 7
	SMSTypeTenNumber     = 8
	SMSTypeTwoWay        = 13
)

// Response represents shared interface between response messages
type Response interface {
	GetError() error
}

type SendSMSRequest struct {
	Phone   string
	Content string
	SmsType int32
}

type SendSMSResponse struct {
	CodeResult      string
	CountRegenerate int32
	SMSID           string
	ErrorMessage    string
}

func (r *SendSMSResponse) GetError() error {
	if r.CodeResult != CodeSuccess {
		return cm.Error(cm.ExternalServiceError, "Không thể gửi tin nhắn", nil).WithMeta("msg", r.ErrorMessage)
	}
	return nil
}

type GetBalanceResponse struct {
	Balance      int32
	CodeResponse string
	ErrorMessage string
	UserID       int32
}

func (r *GetBalanceResponse) GetError() error {
	if r.CodeResponse != CodeSuccess {
		return cm.Error(cm.ExternalServiceError, "sms: "+r.ErrorMessage, nil).WithMeta("msg", r.ErrorMessage)
	}
	return nil
}

type GetSendStatusRequest struct {
	RefId string
}

type GetSendStatusResponse struct {
	CodeResponse  string
	SMSID         string
	SendFailed    int32
	SendStatus    int32
	SendSuccess   int32
	TotalReceiver int32
	TotalSend     int32
}

func (r GetSendStatusResponse) GetError() error {
	if r.CodeResponse != CodeSuccess {
		return cm.Error(cm.ExternalServiceError, "sms: Code "+r.CodeResponse, nil)
	}
	return nil
}

type GetSMSReceiverStatusRequest struct {
	RefId string
}

type Receiver struct {
	IsSent     bool
	Phone      string
	SentResult bool
}

type GetSMSReceiverStatusResponse struct {
	CodeResult   string
	ErrorMessage string
	ReceiverList []Receiver
}

func (r GetSMSReceiverStatusResponse) GetError() error {
	if r.CodeResult != CodeSuccess {
		return cm.Error(cm.ExternalServiceError, "sms: "+r.ErrorMessage, nil).WithMeta("msg", r.ErrorMessage)
	}
	return nil
}

// ESMS ...
type ESMS struct {
	cfg Config
}

// New ...
func New(cfg Config) *ESMS {
	_, err := validate.ValidateStruct(cfg)
	if err != nil {
		ll.Fatal("Invalid esms config", l.Error(err))
	}
	return &ESMS{cfg: cfg}
}

// Request ...
func (c *ESMS) Request(ctx context.Context, path string, resp Response) (err error) {
	t0 := time.Now()

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return cm.Error(cm.Internal, "", err)
	}
	req = req.WithContext(ctx)

	httpResp, err := client.Do(req)
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "", err)
	}

	// Decode response
	data, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "", err)
	}
	respData := data
	if len(respData) > 10000 {
		respData = respData[:5000]
	}
	t1 := time.Now()
	ll.Info("SMS→ "+c.cfg.BaseURL+path,
		l.Duration("t", t1.Sub(t0)),
		l.String("\n⇐", string(respData)))

	err = json.Unmarshal(data, resp)
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "", err)
	}
	return resp.GetError()
}

// GetBalance ...
func (c *ESMS) GetBalance(ctx context.Context) (*GetBalanceResponse, error) {
	var resp GetBalanceResponse
	path := fmt.Sprintf("%s/%s/%s", APIGetBalance, c.cfg.APIKey, c.cfg.SecretKey)
	err := c.Request(ctx, path, &resp)
	return &resp, err
}

// SendSMS ...
func (c *ESMS) SendSMS(ctx context.Context, smsType int32, phone string, content string) (*SendSMSResponse, error) {
	useBrand := smsType == SMSTypeBrandAds || smsType == SMSTypeBrandCustomer

	u, _ := url.ParseRequestURI(c.cfg.BaseURL + APISendSMS)
	q := u.Query()
	q.Add("Phone", phone)
	q.Add("Content", content)
	q.Add("ApiKey", c.cfg.APIKey)
	q.Add("SecretKey", c.cfg.SecretKey)
	q.Add("SmsType", fmt.Sprintf("%d", smsType))
	if useBrand {
		q.Add("Brandname", c.cfg.BrandName)
	}
	u.RawQuery = q.Encode()

	var resp SendSMSResponse
	err := c.Request(ctx, u.String(), &resp)
	return &resp, err
}

// GetSendStatus ...
func (c *ESMS) GetSendStatus(ctx context.Context, smsID string) (*GetSendStatusResponse, error) {
	u, _ := url.ParseRequestURI(c.cfg.BaseURL + APIGetSendStatus)
	q := u.Query()
	q.Add("ApiKey", c.cfg.APIKey)
	q.Add("SecretKey", c.cfg.SecretKey)
	q.Add("RefId", smsID)
	u.RawQuery = q.Encode()
	var resp GetSendStatusResponse
	err := c.Request(ctx, u.String(), &resp)
	return &resp, err
}

// GetSMSReceiverStatus ...
func (c *ESMS) GetSMSReceiverStatus(ctx context.Context, smsID string) (*GetSMSReceiverStatusResponse, error) {
	u, _ := url.ParseRequestURI(c.cfg.BaseURL + APIGetSMSReceiverStatus)
	q := u.Query()
	q.Add("ApiKey", c.cfg.APIKey)
	q.Add("SecretKey", c.cfg.SecretKey)
	q.Add("RefId", smsID)
	u.RawQuery = q.Encode()
	var resp GetSMSReceiverStatusResponse
	err := c.Request(ctx, u.String(), &resp)
	return &resp, err
}
