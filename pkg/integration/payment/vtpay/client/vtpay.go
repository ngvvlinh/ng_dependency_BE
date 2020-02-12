package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpreq"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/l"
	"github.com/gorilla/schema"
)

var (
	ll      = l.New()
	encoder = schema.NewEncoder()
)

func init() {
	encoder.SetAliasTag("url")
}

type Config struct {
	Env          string `yaml:"env" valid:"required"`
	AccessCode   string `yaml:"access_code" valid:"required"`
	SecretKey    string `yaml:"secret_key" valid:"required"`
	MerchantCode string `yaml:"merchant_code" valid:"required"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_PAYMENT"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_VTPAY_ENV":           &c.Env,
		p + "_VTPAY_ACCESS_CODE":   &c.AccessCode,
		p + "_VTPAY_MERCHANT_CODE": &c.MerchantCode,
		p + "_VTPAY_SECRET_KEY":    &c.SecretKey,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env:          cmenv.PartnerEnvTest,
		AccessCode:   "d41d8cd98f00b204e9800998ecf8427eb36ee15cfa9bbaf984868d25cc21215a",
		SecretKey:    "d41d8cd98f00b204e9800998ecf8427e1c371a7425bf2b1a5fdf5a83c5371a40",
		MerchantCode: "ETOP",
	}
}

type Client struct {
	cfg          Config
	MerchantCode string
	rclient      *httpreq.Resty
	baseUrl      string
}

func New(cfg Config) *Client {
	if _, err := validate.ValidateStruct(cfg); err != nil {
		ll.Fatal("vtpay: invalid config", l.Error(err))
	}
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		cfg:          cfg,
		MerchantCode: cfg.MerchantCode,
		rclient:      httpreq.NewResty(rcfg),
	}
	switch cfg.Env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://sandbox.viettel.vn/"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://pay.bankplus.vn:8450/"
	default:
		ll.Fatal("vtpay: invalid env")
	}
	return c
}

func (c *Client) CheckSum(data string) string {
	hash := hmac.New(sha1.New, []byte(c.cfg.SecretKey))
	data = c.cfg.AccessCode + data
	_, err := io.WriteString(hash, data)
	if err != nil {
		panic(err)
	}
	macSum := hash.Sum(nil)
	dd := base64.StdEncoding.EncodeToString(macSum)
	return dd
}

func (c *Client) BuildUrlConnectPaymentGateway(ctx context.Context, req *ConnectPaymentGatewayRequest) (string, error) {
	req.MerchantCode = c.cfg.MerchantCode
	req.Command = "PAYMENT"
	req.Locale = "Vi"
	req.Version = "2.0"
	if err := req.Validate(); err != nil {
		return "", err
	}
	req.CheckSum = c.CheckSum(req.DataCheckSum())
	query, err := buildQueryString(req)
	if err != nil {
		return "", err
	}
	return c.baseUrl + "PaymentGateway/payment?" + query, nil
}

func buildQueryString(req interface{}) (string, error) {
	if req == nil {
		return "", nil
	}
	queryString := url.Values{}
	if err := encoder.Encode(req, queryString); err != nil {
		return "", err
	}
	return queryString.Encode(), nil
}

func (c *Client) GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error) {
	req.Cmd = "TRANS_INQUIRY"
	req.MerchantCode = c.cfg.MerchantCode
	req.Checksum = c.CheckSum(req.DataCheckSum())

	var resp GetTransactionResponse
	fullUrl := c.baseUrl + "PaymentAPI/webresources/postData"
	if err := c.sendPostFormRequest(ctx, fullUrl, req, &resp, "Không thể lấy thông tin giao dịch"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CancelTransaction(ctx context.Context, req *CancelTransactionRequest) (*CancelTransactionResponse, error) {
	req.Cmd = "REFUND_PAYMENT"
	req.RefundType = 1
	req.Version = "2.0"
	req.CheckSum = c.CheckSum(req.DataCheckSum())

	var resp CancelTransactionResponse
	fullUrl := c.baseUrl + "PaymentAPI/webresources/postData"
	if err := c.sendPostFormRequest(ctx, fullUrl, req, &resp, "Không thể hủy giao dịch"); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) sendPostFormRequest(ctx context.Context, fullURL string, req interface{}, resp interface{}, msg string) error {
	values := url.Values{}
	if req != nil {
		if err := encoder.Encode(req, values); err != nil {
			return err
		}
	}
	var formData = make(map[string]string)
	for key := range values {
		formData[key] = values.Get(key)
	}
	res, err := c.rclient.R().
		SetFormData(formData).Post(fullURL)
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối VTPay")
	}
	return httpreq.HandleResponse(ctx, res, resp, msg)
}
