package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/jsonx"
	"o.o/common/l"
	"time"
)

var ll = l.New()

type NTXOrderServiceCode string

const (
	UTMSource           = "api_etop"
	PathCalcShippingFee = "/v1/bill/calc-fee"
	PathCreateOrder     = "/v1/bill/create"
	PathCancelOrder     = "/v1/bill/destroy"

	OrderServiceCodeCH NTXOrderServiceCode = "STANDARD"
	OrderServiceCodeNH NTXOrderServiceCode = "FAST"
)

type Client struct {
	baseUrl       string
	PartnerID     int
	PaymentMethod int
	headers       map[string]string
	rclient       *httpreq.Resty
}

func New(env string, cfg Config) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		headers: map[string]string{
			"username": cfg.Username,
			"password": cfg.Password,
		},
		PartnerID: cfg.PartnerID,
		rclient:   httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://apisandbox.ntx.com.vn"
		c.PaymentMethod = 10
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://apiws.ntx.com.vn"
		c.PaymentMethod = 11
	default:
		ll.Fatal("NTX: Invalid env")
	}

	return c
}

func (c *Client) CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error) {
	var resp CalcShippingFeeResponse
	err := c.sendPostRequest(ctx, PathCalcShippingFee, req, &resp, "Không thể tính phí giao hàng")
	return &resp, err
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	var resp CreateOrderResponse
	err := c.sendPostRequest(ctx, PathCreateOrder, req, &resp, "Không thể tạo đơn hàng")
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CommonResponse, error) {
	var resp CommonResponse
	err := c.sendPostRequest(ctx, PathCancelOrder, req, &resp, "Không thể hủy đơn hàng")
	return &resp, err
}

func (c *Client) sendPostRequest(ctx context.Context, path string, body interface{}, resp ResponseInterface, msg string) error {
	res, err := c.rclient.R().
		SetBody(body).
		SetHeaders(c.headers).
		Post(c.baseUrl + path)
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với NTX: %v (%v)", msg, err)
	}
	err = handleResponse(ctx, res, resp, msg)
	return err
}

func handleResponse(ctx context.Context, res *httpreq.RestyResponse, result ResponseInterface, msg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ NTX: null response. Chúng tôi đang liên hệ với NTX để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
			}
			if err = jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ NTX: %v. Chúng tôi đang liên hệ với NTX để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
			}

			cr := result.GetCommonResponse()
			if !cr.Success {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ NTX: %v (%v)", msg, cr.Message).
					WithMeta("messsage", cr.Message)
			}

		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(body) {
			if err = jsonx.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ NTX. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ NTX: Invalid status (%v). Chúng tôi đang liên hệ với NTX để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v", status, wl.X(ctx).CSEmail)
	}
}
