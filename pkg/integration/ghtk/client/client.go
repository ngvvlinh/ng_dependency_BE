package ghtkclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	resty "gopkg.in/resty.v1"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpreq"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/common/l"
)

var ll = l.New()
var encoder = schema.NewEncoder()

func init() {
	encoder.SetAliasTag("url")
}

type Client struct {
	baseUrl string
	token   string
	rclient *resty.Client
}

const (
	PathCalcShippingFee = "/services/shipment/fee"
	PathCreateOrder     = "/services/shipment/order"
	PathGetOrder        = "/services/shipment/v2"
	PathCancelOrder     = "/services/shipment/cancel"
)

func New(env string, cfg GhtkAccount) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	c := &Client{
		token:   cfg.Token,
		rclient: resty.NewWithClient(client).SetDebug(true),
	}
	switch env {
	case cm.PartnerEnvTest:
		c.baseUrl = "https://dev.ghtk.vn/"
	case cm.PartnerEnvProd:
		c.baseUrl = "https://services.giaohangtietkiem.vn/"
	default:
		ll.Fatal("ghtk: Invalid ENV")
	}
	return c
}

func (c *Client) Ping() error {
	req := &CalcShippingFeeRequest{
		Weight:          1000,
		Value:           3000000,
		PickingProvince: "Hà Nội",
		PickingDistrict: "Quận Hai Bà Trưng",
		Address:         "163 Tô Hiến Thành",
		Province:        "Hồ Chí Minh",
		District:        "Quận 10",
		Transport:       "road",
	}
	_, err := c.CalcShippingFee(context.Background(), req)
	return err
}

func (c *Client) TestCreateOrder() error {
	req := &CreateOrderRequest{
		Products: []*ProductRequest{
			&ProductRequest{
				Name:     "bút",
				Weight:   0.1,
				Quantity: 1,
			},
			&ProductRequest{
				Name:     "tẩy",
				Weight:   0.2,
				Quantity: 1,
			},
		},
		Order: &OrderRequest{
			ID:           "1234.5678",
			PickName:     "HCM-nội thành",
			PickAddress:  "590 CMT8 P.11",
			PickProvince: "TP. Hồ Chí Minh",
			PickDistrict: "Quận 3",
			PickTel:      "0911222333",
			Tel:          "0911222333",
			Name:         "GHTK - HCM - Noi Thanh",
			Address:      "123 nguyễn chí thanh",
			Province:     "TP. Hồ Chí Minh",
			District:     "Quận 1",
			IsFreeship:   1,
			PickMoney:    47000,
			Note:         "Khối lượng tính cước tối đa: 1.00 kg",
			Value:        3000000,
		},
	}
	_, err := c.CreateOrder(context.Background(), req)
	return err
}

func (c *Client) TestGetOrder() error {
	labelID := "S1858017.SG11.10I.236273443"
	_, err := c.GetOrder(context.Background(), labelID, "")
	return err
}

func (c *Client) TestCancelOrder() error {
	labelID := "S1858017.298467358"
	_, err := c.CancelOrder(context.Background(), labelID, "")
	return err
}

func (c *Client) CalcShippingFee(ctx context.Context, req *CalcShippingFeeRequest) (*CalcShippingFeeResponse, error) {
	req.PickingProvince = validate.NormalizeSearchSimple(req.PickingProvince)
	req.PickingDistrict = validate.NormalizeSearchSimple(req.PickingDistrict)
	req.PickingWard = validate.NormalizeSearchSimple(req.PickingWard)
	req.PickingAddress = validate.NormalizeSearchSimple(req.PickingAddress)
	req.Province = validate.NormalizeSearchSimple(req.Province)
	req.District = validate.NormalizeSearchSimple(req.District)
	req.Ward = validate.NormalizeSearchSimple(req.Ward)
	req.Address = validate.NormalizeSearchSimple(req.Address)
	var resp CalcShippingFeeResponse
	err := c.sendGetRequest(ctx, PathCalcShippingFee, req, &resp,
		"Không thể tính phí giao hàng")
	return &resp, err
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// Setup freeship: COD sẽ chỉ thu người nhận hàng số tiền bằng pick_money
	req.Order.IsFreeship = 1
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp CreateOrderResponse
	err := c.sendPostRequest(ctx, PathCreateOrder, req, &resp, "Không thể tạo đơn hàng")
	return &resp, err
}

func (c *Client) GetOrder(ctx context.Context, labelID, orderPartnerID string) (*GetOrderResponse, error) {
	endPoint := PathGetOrder
	if labelID != "" {
		endPoint += "/" + labelID
	} else if orderPartnerID != "" {
		endPoint += "/" + orderPartnerID
	}
	var resp GetOrderResponse
	err := c.sendGetRequest(ctx, endPoint, nil, &resp, "Không thể lấy thông tin đơn hàng")
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, labelID, orderPartnerID string) (*CommonResponse, error) {
	endPoint := PathCancelOrder
	if labelID != "" {
		endPoint += "/" + labelID
	} else if orderPartnerID != "" {
		endPoint += "/" + orderPartnerID
	}
	var resp CommonResponse
	err := c.sendPostRequest(ctx, endPoint, nil, &resp, "Không thể hủy đơn hàng")
	return &resp, err
}

func (c *Client) sendGetRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	queryString := url.Values{}
	if req != nil {
		err := encoder.Encode(req, queryString)
		if err != nil {
			return cm.Error(cm.Internal, "", err)
		}
	}

	res, err := c.rclient.R().
		SetQueryString(queryString.Encode()).
		SetHeader("token", c.token).
		Get(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với GHTK", err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

func (c *Client) sendPostRequest(ctx context.Context, path string, req interface{}, resp ResponseInterface, msg string) error {
	res, err := c.rclient.R().
		SetBody(req).
		SetHeader("token", c.token).
		Post(model.URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với GHTK", err)
	}
	err = handleResponse(res, resp, msg)
	return err
}

func handleResponse(res *resty.Response, result ResponseInterface, msg string) error {
	status := res.StatusCode()
	var err error
	body := res.Body()
	switch {
	case status >= 200 && status < 300:
		if result != nil {
			if httpreq.IsNullJsonRaw(body) {
				return cm.Error(cm.ExternalServiceError, "Lỗi không xác định từ Giaohangtietkiem: null response. Chúng tôi đang liên hệ với Giaohangtietkiem để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", nil)
			}
			if err = json.Unmarshal(body, result); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Giaohangtietkiem: %v. Chúng tôi đang liên hệ với Giaohangtietkiem để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", err)
			}
			cr := result.GetCommonResponse()
			if !cr.Success {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Giaohangtietkiem: %v (%v)", msg, cr.Message).
					WithMeta("messsage", cr.Message)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(body) {
			if err = json.Unmarshal(body, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = json.Unmarshal(body, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
		}

		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi từ Giaohangtietkiem. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.").WithMetaM(meta)
	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giaohangtietkiem: Invalid status (%v). Chúng tôi đang liên hệ với Giaohangtietkiem để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ hotro@etop.vn.", status)
	}
}
