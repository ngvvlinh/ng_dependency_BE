package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/schema"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
	"etop.vn/common/xerrors"
)

var (
	encoder = schema.NewEncoder()
	ll      = l.New()
)

const (
	PathShopInfo        = "shop"
	PathCarrierServices = "carrier_services"
)

func init() {
	encoder.SetAliasTag("url")
}

type Client struct {
	ApiKey  string
	Secret  string
	rclient *httpreq.Resty
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
		ApiKey:  cfg.APIKey,
		Secret:  cfg.Secret,
		rclient: httpreq.NewResty(rcfg),
	}
	return c
}

func (c *Client) GetShop(ctx context.Context, req *GetShopRequest) (*Shop, error) {
	var resp GetShopResponse
	if err := c.sendGetRequest(ctx, req.Connection, PathShopInfo, req, &resp, "Không thể lấy thông tin shop"); err != nil {
		return nil, err
	}
	return resp.Shop, nil
}

func (c *Client) ConnectCarrierService(ctx context.Context, req *ConnectCarrierServiceRequest) (*CarrierService, error) {
	req.CarrierService.CarrierServiceType = "api"
	var resp ConnectCarrierServiceResponse
	if err := c.sendPostRequest(ctx, req.Connection, PathCarrierServices, req, &resp, "Khổng thể tạo kết nối với nhà vận chuyển"); err != nil {
		return nil, err
	}
	return resp.CarrierService, nil
}

func (c *Client) GetCarrierServices(ctx context.Context, req *GetCarrierServicesRequest) ([]*CarrierServiceItem, error) {
	var resp GetCarrierServicesResponse
	if err := c.sendGetRequest(ctx, req.Connection, PathCarrierServices, nil, &resp, "Không thể lấy danh sách Carrier Services"); err != nil {
		return nil, err
	}
	return resp.CarrierServices, nil
}

func (c *Client) DeleteConnectedCarrierService(ctx context.Context, req *DeleteConnectedCarrierServiceRequest) error {
	path := fmt.Sprintf("%v/%v", PathCarrierServices, req.CarrierServiceID)
	if err := c.SendDeleteRequest(ctx, req.Connection, path, req, nil, "Không thể xóa kết nối nhà vận chuyển"); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccessToken(ctx context.Context, req *GetAccessTokenRequest) (*GetAccessTokenResponse, error) {
	var resp GetAccessTokenResponse

	formData := map[string]string{
		"subdomain":     req.Subdomain,
		"client_id":     c.ApiKey,
		"client_secret": c.Secret,
		"code":          req.Code,
		"grant_type":    "authorization_code",
		"redirect_uri":  req.RedirectURI,
	}

	urlStr := fmt.Sprintf("https://%v.myharavan.com/admin/oauth/access_token", req.Subdomain)
	res, err := c.rclient.R().
		SetFormData(formData).Post(urlStr)
	if err != nil {
		return nil, cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với Haravan")
	}

	if err = handleResponse(res, &resp, ""); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateShippingState(ctx context.Context, req *UpdateShippingStateRequest) error {
	if !req.State.CanUpdate() {
		return cm.Errorf(cm.InvalidArgument, nil, "API Haravan không hỗ trợ cập nhật trạng thái vận chuyển (%v).", req.State)
	}
	url := fmt.Sprintf("https://%v.myharavan.com/admin/orders/%v/fulfillments/%v/carrier_status/%v.json", req.Subdomain, req.OrderID, req.FulfillID, string(req.State))

	if err := c.SendPutRequest(ctx, req.Connection, url, req, nil, "Không thể cập nhật trạng thái giao hàng"); err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdatePaymentStatus(ctx context.Context, req *UpdatePaymentStatusRequest) error {
	if !req.Status.CanUpdate() {
		return cm.Errorf(cm.InvalidArgument, nil, "API Haravan không hỗ trợ cập nhật trạng thái thanh toán (%v).", req.Status)
	}
	url := fmt.Sprintf("https://%v.myharavan.com/admin/orders/%v/fulfillments/%v/cod_status/%v.json", req.Subdomain, req.OrderID, req.FulfillID, string(req.Status))

	if err := c.SendPutRequest(ctx, req.Connection, url, req, nil, "Không thể cập nhật trạng thái thanh toán"); err != nil {
		return err
	}
	return nil
}

func (c *Client) sendGetRequest(ctx context.Context, connection Connection, path string, req interface{}, resp interface{}, msg string) error {
	if connection.TokenStr == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu shop access token")
	}
	url := buildUrl(connection.Subdomain, path)
	res, err := c.rclient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", connection.TokenStr)).
		Get(url)
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với Haravan", err)
	}
	return handleResponse(res, resp, msg)
}

func (c *Client) sendPostRequest(ctx context.Context, connection Connection, path string, req interface{}, resp interface{}, msg string) error {
	if connection.TokenStr == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu shop access token")
	}
	url := buildUrl(connection.Subdomain, path)
	res, err := c.rclient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", connection.TokenStr)).
		SetBody(req).
		Post(url)
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với Haravan")
	}
	return handleResponse(res, resp, msg)
}

func (c *Client) SendDeleteRequest(ctx context.Context, connection Connection, path string, req interface{}, resp interface{}, msg string) error {
	if connection.TokenStr == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu shop access token")
	}
	url := buildUrl(connection.Subdomain, path)
	res, err := c.rclient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", connection.TokenStr)).
		SetBody(req).
		Delete(url)
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với Haravan")
	}
	return handleResponse(res, resp, msg)
}

func (c *Client) SendPutRequest(ctx context.Context, connection Connection, url string, req interface{}, resp interface{}, msg string) error {
	if connection.TokenStr == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Thiếu shop access token")
	}
	res, err := c.rclient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %v", connection.TokenStr)).
		Put(url)
	if err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Lỗi kết nối với Haravan")
	}
	return handleResponse(res, resp, msg)
}

func buildUrl(subdomain string, path string) string {
	return fmt.Sprintf("https://%v.myharavan.com/admin/%v.json", subdomain, path)
}

func handleResponse(res *httpreq.RestyResponse, result interface{}, msg string) error {
	status := res.StatusCode()
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Haravan: null response.")
			}
			if err := jsonx.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Haravan: %v", err)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		var errJSON xerrors.ErrorJSON
		if !httpreq.IsNullJsonRaw(body) {
			if err := jsonx.Unmarshal(body, &meta); err != nil {
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
			errJSON.Msg = msg
			errJSON.Meta = meta
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Haravan: %v. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn", errJSON.Error()).WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Haravan: %v. Invalid status (%v).", msg, status)
	}
}
