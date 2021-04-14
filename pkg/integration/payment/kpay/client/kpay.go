package client

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"github.com/Luzifer/go-openssl/v4"
	"golang.org/x/net/context"
	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/common/l"
)

var (
	ll = l.New()
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	DELETE Method = "DELETE"
)

type Config struct {
	Env         string `yaml:"env"`
	PublicKey   string `yaml:"public_key"`
	PrivateKey  string `yaml:"private_key"`
	XApiClient  string `yaml:"x_api_client"`
	AccessToken string `yaml:"access_token"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_PAYMENT"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_KPAY_ENV":          &c.Env,
		p + "_KPAY_PUBLIC_KEY":   &c.PublicKey,
		p + "_KPAY_PRIVATE_KEY":  &c.PrivateKey,
		p + "_KPAY_X_API_CLIENT": &c.XApiClient,
		p + "_KPAY_ACCESS_TOKEN": &c.AccessToken,
	}.MustLoad()
}

func DefaultConfig() Config {
	return Config{
		Env:         cmenv.PartnerEnvTest,
		PublicKey:   "-----BEGIN PUBLIC KEY-----\nMFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAMFNA4t9npeTL02Y+YJFGRFLJGyZRkAT\n0ypuLnZary3tHq4h8M8uyGR0p5KpJrw4S+zQnsbV71Pq8Lyl2btEQ/MCAwEAAQ==\n-----END PUBLIC KEY-----",
		PrivateKey:  "-----BEGIN RSA PRIVATE KEY-----\nMIIBOwIBAAJBAKVMPhKo4nwXVjargWGRx1uqHmNr8RnvN0ti+ErcYmJ62g/xl1E3\nHPrEM4UQab7lWOl3P0mwKo/RplTlvTx7EEUCAwEAAQJAUQiPZZZlcW//U83AH5TX\nppg/TX4dNBmRgeOC1TA1CxFjLhCZ6sZl9fjiKYx3QMLhiLjlcbpvCkOKgu/tZFdQ\nAQIhAOx55TO/3sCDFJbH1FYwl4d2jAJ+fx4b/uzUK0GVpC4dAiEAsvHwRwv88eou\nZmi4PswOT1ENCXVDp/PuGNBTmI7VckkCIQCWV8ESr4+tESlL22vyCB6ubv4Ar++d\ncusWIqYFol+pOQIhAJkSRqSvx4808L4fpEPrj+4fehR2IArCqhDF3EyrGkEpAiBY\nkDYAYHuGPNax+W30BIgf7qybIijEPawEAiK/M+p/qQ==\n-----END RSA PRIVATE KEY-----",
		XApiClient:  "201605456902",
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NDA5OSwiYXBwSWQiOiIyMDE2MDU0NTY5MDIiLCJtZXJjaGFudElkIjoxNTYyMDYsImFjY291bnRJZCI6NTA0LCJ0eXBlIjoiQVBQIiwiaWF0IjoxNjE2MTQ0OTY2fQ.T-9Kkvq1LxZl4WJkjl7hfMGrb8MYH4zIPNQK0QNauRI",
	}
}

type Client struct {
	cfg        Config
	rClient    *httpreq.Resty
	baseUrl    string
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New(cfg Config) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		cfg:        cfg,
		rClient:    httpreq.NewResty(rcfg),
		privateKey: bytesToPrivateKey([]byte(cfg.PrivateKey)),
		publicKey:  bytesToPublicKey([]byte(cfg.PublicKey)),
	}
	switch cfg.Env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://sbx-pg.payme.vn"
	case cmenv.PartnerEnvProd:
		c.baseUrl = ""
	default:
		ll.Fatal("KPay: invalid env")
	}
	return c
}

func (c Client) CreateTransaction(ctx context.Context, args *CreateTransactionRequest) (*CreateTransactionResponse, error) {
	var resp CreateTransactionResponse
	if err := c.sendPostRequest(ctx, "/v1/Payment/Generate", args, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c Client) GetTransaction(ctx context.Context, args *GetTransactionRequest) (*GetTransactionResponse, error) {
	var resp GetTransactionResponse
	path := fmt.Sprintf("/v1/Payment/Information/%s", args.Transaction)
	if err := c.sendGetRequest(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c Client) CancelTransaction(ctx context.Context, args *CancelTransactionRequest) (*CancelTransactionResponse, error) {
	var resp CancelTransactionResponse
	if err := c.sendPostRequest(ctx, "/v1/Payment/Cancel/Pending", args, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c Client) sendGetRequest(ctx context.Context, path string, resp interface{}) error {
	return c.sendRequest(ctx, path, GET, nil, resp)
}

func (c Client) sendPostRequest(ctx context.Context, path string, payload, resp interface{}) error {
	return c.sendRequest(ctx, path, POST, payload, resp)
}

func (c Client) sendRequest(
	ctx context.Context, path string,
	method Method, payload, resp interface{},
) (err error) {
	var (
		bodyResp BodyResponse
		res      *httpreq.RestyResponse
		req      *resty.Request
	)

	key := cm.NewID().String()
	xApiKey := base64.StdEncoding.EncodeToString(encryptWithPublicKey([]byte(key), c.publicKey))
	xApiAction := string(encryptAES(key, path))

	var xApiMessage string
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}
		xApiMessage = string(encryptAES(key, string(payloadBytes)))
	}

	xApiValidate := createHash(fmt.Sprintf("%s%s%s%s%s", xApiAction, method, c.cfg.AccessToken, xApiMessage, key))
	headers := map[string]string{
		"x-api-client":   c.cfg.XApiClient,
		"x-api-key":      xApiKey,
		"x-api-action":   xApiAction,
		"x-api-validate": xApiValidate,
		"authorization":  c.cfg.AccessToken,
	}
	requestBody := &BodyRequest{
		XAPIMessage: xApiMessage,
	}

	req = c.rClient.R().
		SetHeaders(headers).
		SetBody(requestBody).
		SetResult(&bodyResp).
		SetError(&bodyResp)

	switch method {
	case GET:
		res, err = req.Get(url(c.baseUrl, path))
	case POST:
		res, err = req.Post(url(c.baseUrl, path))
	case DELETE:
		res, err = req.Delete(url(c.baseUrl, path))
	default:
		return cm.Errorf(cm.Internal, nil, "KPay: unsupported method %v", method)
	}

	status := res.StatusCode()
	switch {
	case status >= 200 && status < 300:
		if resp != nil {
			return c.handleResponseBody(ctx, bodyResp, res.Header().Get("x-api-key"), resp)
		}
		return nil
	case status >= 400:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ KPay: nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ KPay: Invalid status (%v). Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}

func (c Client) handleResponseBody(ctx context.Context, bodyResp BodyResponse, xApiKey string, resp interface{}) error {
	var decryptedBody DecryptedBodyResponse
	xApiKeyDecode, err := base64.StdEncoding.DecodeString(xApiKey)
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ KPay: %v. Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
	}
	decryptedKey := decryptWithPrivateKey(xApiKeyDecode, c.privateKey)

	decryptedBodyBytes := DecryptAES(string(decryptedKey), bodyResp.XApiMessage)
	fmt.Println(string(decryptedBodyBytes))
	if err := json.Unmarshal(decryptedBodyBytes, &decryptedBody); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ KPay: %v. Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
	}

	if _, ok := mapErrorCode[decryptedBody.Code]; !ok {
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ KPay: Invalid errorCode (%v). Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", decryptedBody.Code, wl.X(ctx).CSEmail)
	}

	if decryptedBody.Code != RequestSuccess {
		var errResponse ErrorResponse
		if err := json.Unmarshal(decryptedBody.Data, &errResponse); err != nil {
			return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ KPay: %v. Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ KPay: %v. Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", errResponse.Message, wl.X(ctx).CSEmail)
	}

	if err := json.Unmarshal(decryptedBody.Data, &resp); err != nil {
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ KPay: %v. Chúng tôi đang liên hệ với KPay để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
	}

	return nil
}

// bytesToPrivateKey bytes to private key
func bytesToPrivateKey(pri []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(pri)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		ll.Info("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			ll.Fatal("bytesToPrivateKey: %v", l.Object("err", err))
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		ll.Fatal("bytesToPrivateKey: %v", l.Object("err", err))
	}
	return key
}

// decryptWithPrivateKey encrypts data with private key
func decryptWithPrivateKey(msg []byte, pri *rsa.PrivateKey) []byte {
	ciphertext, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, pri, msg, nil)
	if err != nil {
		ll.Panic("decryptWithPrivateKey: %v", l.Object("err", err))
	}
	return ciphertext
}

// bytesToPublicKey bytes to public key
func bytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		ll.Info("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			ll.Panic("bytesToPublicKey: %v", l.Object("err", err))
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		ll.Panic("bytesToPublicKey: %v", l.Object("err", err))
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		ll.Panic("bytesToPublicKey: %v", l.Object("err", err))
	}
	return key
}

// encryptWithPublicKey encrypts data with public key
func encryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	ciphertext, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		ll.Panic("encryptWithPublicKey: %v", l.Object("err", err))
	}
	return ciphertext
}

func encryptAES(key, plainString string) []byte {
	o := openssl.New()

	encryptKey, err := o.EncryptBytes(key, []byte(plainString), openssl.BytesToKeyMD5)
	if err != nil {
		ll.Panic("encryptAES: %v", l.Object("err", err))
	}
	return encryptKey
}

func DecryptAES(key, plainString string) []byte {
	o := openssl.New()

	decryptKey, err := o.DecryptBytes(key, []byte(plainString), openssl.BytesToKeyMD5)
	if err != nil {
		ll.Panic("decryptAES: %v", l.Object("err", err))
	}
	return decryptKey
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
