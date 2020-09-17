package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/resty.v1"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()

type Client struct {
	baseUrl      string
	token        string
	clientID     string
	clientSecret string
	accountID    string

	rclient *httpreq.Resty
}

type sendRequestArgs struct {
	path string
	req  interface{}
	resp interface{}
}

type Method string
type DHLOrderServiceCode string
type PickupAccountID string

// PickupAccountID là ID 3 kho ở Bắc, Trung, Nam của tài khoản TopShip trên DHL
// Địa chỉ lấy hàng được truyền vào trong từng request tạo đơn
// nhưng cần truyền PickupAccountID tương ứng với địa chỉ lấy hàng trên
var PickupAccountIDNorth, PickupAccountIDMiddle, PickupAccountIDSouth PickupAccountID

const (
	GetMethod  Method = "GET"
	PostMethod Method = "POST"

	// Giao tối ưu
	OrderServiceCodeSPD DHLOrderServiceCode = "SPD"
	// Giao nhanh
	OrderServiceCodePDE DHLOrderServiceCode = "PDE"
	// Giao tiêu chuẩn
	OrderServiceCodePDO DHLOrderServiceCode = "PDO"
)

func New(env string, cfg DHLAccountCfg) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:        cfg.Token,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		accountID:    cfg.AccountID,
		rclient:      httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		PickupAccountIDNorth = "5351118"
		PickupAccountIDMiddle = "5351118"
		PickupAccountIDSouth = "5351118"
		c.baseUrl = "https://apitest.dhlecommerce.asia/rest"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://api.dhlecommerce.asia/rest"
	default:
		ll.Fatal("DHL: Invalid env")
	}
	return c
}

func (c *Client) UpdateToken(newToken string) {
	c.token = newToken
}

func (c *Client) GenerateAccessToken(ctx context.Context) (*GenerateAccessTokenResponse, error) {
	var resp GenerateAccessTokenResponse

	if err := c.sendGetRequest(ctx, sendRequestArgs{
		path: fmt.Sprintf("/v1/OAuth/AccessToken?clientId=%s&password=%s", c.clientID, c.clientSecret),
		resp: &resp,
	}); err != nil {
		return nil, err
	}

	accessTokenResp := resp.AccessTokenResponse
	if accessTokenResp != nil && accessTokenResp.ResponseStatus != nil && accessTokenResp.ResponseStatus.Code != "100000" {
		return nil, ErrorResp(ctx, accessTokenResp.ResponseStatus.ToError())
	}
	return &resp, nil
}

func (c *Client) CreateOrders(ctx context.Context, req *CreateOrdersRequest) (*CreateOrdersResponse, error) {
	var resp CreateOrdersResponse

	req.ManifestRequest.Hdr = c.generateHdrRequest("SHIPMENT")
	req.ManifestRequest.Bd.SoldToAccountID = c.accountID

	if err := c.sendPostRequest(ctx, sendRequestArgs{
		path: "/v3/Shipment",
		req:  req,
		resp: &resp,
	}); err != nil {
		return nil, err
	}

	if resp.ManifestResponse != nil && resp.ManifestResponse.Bd != nil && len(resp.ManifestResponse.Bd.ShipmentItems) > 0 {
		shipmentItem := resp.ManifestResponse.Bd.ShipmentItems[0]
		responseStatus := shipmentItem.ResponseStatus
		if responseStatus != nil && responseStatus.Code != "200" {
			return nil, ErrorResp(ctx, responseStatus.ToError())
		}
	} else {
		responseStatus := resp.ManifestResponse.Bd.ResponseStatus
		return nil, ErrorResp(ctx, responseStatus.ToError())
	}

	return &resp, nil
}

func (c *Client) FindAvailableServices(ctx context.Context) *FindAvailableServicesResponse {
	return &FindAvailableServicesResponse{
		AvailableServices: []*AvailableService{
			{
				Name:        "Giao tối ưu",
				ServiceCode: String(OrderServiceCodeSPD),
			},
			{
				Name:        "Giao nhanh",
				ServiceCode: String(OrderServiceCodePDE),
			},
			{
				Name:        "Giao tiêu chuẩn",
				ServiceCode: String(OrderServiceCodePDO),
			},
		},
	}
}

func (c *Client) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CancelOrderResponse, error) {
	var resp CancelOrderResponse

	req.DeleteShipmentReq.Hdr = c.generateHdrRequest("DELETESHIPMENT")
	req.DeleteShipmentReq.Bd.SoldToAccountId = c.accountID

	if err := c.sendPostRequest(ctx, sendRequestArgs{
		path: "/v2/Label/Delete",
		req:  req,
		resp: &resp,
	}); err != nil {
		return nil, err
	}

	if resp.DeleteShipmentResp != nil && resp.DeleteShipmentResp.Bd != nil && len(resp.DeleteShipmentResp.Bd.ShipmentItems) > 0 {
		shipmentItem := resp.DeleteShipmentResp.Bd.ShipmentItems[0]
		responseStatus := shipmentItem.ResponseStatus
		if responseStatus != nil && responseStatus.Code != "200" {
			return nil, ErrorResp(ctx, responseStatus.ToError())
		}
	} else {
		responseStatus := resp.DeleteShipmentResp.Bd.ResponseStatus
		return nil, ErrorResp(ctx, responseStatus.ToError())
	}

	return &resp, nil
}

func (c *Client) TrackingOrder(ctx context.Context, req *TrackingOrdersRequest) (*TrackingOrdersResponse, error) {
	var resp TrackingOrdersResponse

	req.TrackItemRequest.Hdr = c.generateHdrRequest("TRACKITEM")

	if err := c.sendPostRequest(ctx, sendRequestArgs{
		path: "/v3/Tracking",
		req:  req,
		resp: &resp,
	}); err != nil {
		return nil, err
	}

	if resp.TrackItemResponse != nil && resp.TrackItemResponse.Bd != nil && len(resp.TrackItemResponse.Bd.ShipmentItems) == 0 {
	} else {
		responseStatus := resp.TrackItemResponse.Bd.ResponseStatus
		if responseStatus.Code != "200" {
			return nil, ErrorResp(ctx, responseStatus.ToError())
		}
	}
	return &resp, nil
}

func (c *Client) generateHdrRequest(messageType string) *HdrReq {
	if messageType == "" {
		return nil
	}
	nowStr := time.Now().Format("2006-01-02T15:04:05-07:00")
	return &HdrReq{
		MessageType:     messageType,
		MessageDateTime: nowStr,
		MessageVersion:  MessageVersion,
		AccessToken:     c.token,
		MessageLanguage: MessageLanguage,
	}
}

func (c *Client) sendPostRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, PostMethod, args)
}

func (c *Client) sendGetRequest(ctx context.Context, args sendRequestArgs) error {
	return c.sendRequest(ctx, GetMethod, args)
}

func ErrorResp(ctx context.Context, err *ErrorResponse) error {
	return cm.Errorf(cm.ExternalServiceError, err, "Lỗi từ DHL: %v. Nếu cần thêm thông tin vui lòng liên hệ %v", err.Error(), wl.X(ctx).CSEmail)
}

func (c *Client) sendRequest(ctx context.Context, method Method, args sendRequestArgs) error {
	var res *resty.Response
	var err error

	request := c.rclient.R().
		SetBody(args.req)

	switch method {
	case PostMethod:
		res, err = request.Post(URL(c.baseUrl, args.path))
	case GetMethod:
		res, err = request.Get(URL(c.baseUrl, args.path))
	default:
		panic(fmt.Sprintf("unsupported method %v", method))
	}
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với DHL", err)
	}

	status := res.StatusCode()
	body := res.Body()
	fmt.Println(string(body))
	switch {
	case status == 200:
		if err = jsonx.Unmarshal(body, args.resp); err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ DHL: %v. Chúng tôi đang liên hệ với DHL để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
		}
		return nil
	case status >= 400:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ DHL. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ DHL: Invalid status (%v). Chúng tôi đang liên hệ với DHL để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
