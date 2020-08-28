package clientv2

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/common/jsonx"
	"o.o/common/l"
)

const Token = "token"
const ShopID = "ShopId"

var ll = l.New()

type Client struct {
	affiliateID int
	clientID    int
	shopID      int

	baseUrl string
	token   string
	rclient *httpreq.Resty
}

func New(env string, cfg GHNAccountCfg) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}
	c := &Client{
		token:       cfg.Token,
		clientID:    cfg.ClientID,
		shopID:      cfg.ShopID,
		affiliateID: cfg.AffiliateID,
		rclient:     httpreq.NewResty(rcfg),
	}
	switch env {
	case cmenv.PartnerEnvTest, cmenv.PartnerEnvDev:
		c.baseUrl = "https://dev-online-gateway.ghn.vn/shiip/public-api"
	case cmenv.PartnerEnvProd:
		c.baseUrl = "https://online-gateway.ghn.vn/shiip/public-api"
	default:
		ll.Fatal("ghn: Invalid env")
	}
	return c
}

func (c *Client) ClientID() int {
	return c.clientID
}

func (c *Client) GetAffiliateID() string {
	if c.affiliateID > 0 {
		return strconv.Itoa(c.affiliateID)
	}
	return ""
}

func (c *Client) GetToken() string {
	return c.token
}

func (c *Client) Ping() error {
	req := &GetServicesRequest{
		FromDistrict: 1447,
		ToDistrict:   1443,
	}
	_, err := c.GetServices(context.Background(), req)
	return err
}

func (c *Client) SendOTPShopToAffiliateAccount(ctx context.Context, req *SendOTPShopAffiliateRequest) (*SendOTPShopAffiliateResponse, error) {
	var resp SendOTPShopAffiliateResponse
	err := c.sendRequest(ctx, "/v2/shop/affiliateOTP", c.shopID, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) CreateShopByAffiliateAccount(ctx context.Context, req *CreateShopAffiliateRequest) (*CreateShopAffiliateResponse, error) {
	var resp CreateShopAffiliateResponse
	err := c.sendRequest(ctx, "/v2/shop/affiliateCreate", c.shopID, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Use API add staff is mean manage your store.
// Doc: https://api.ghn.vn/home/docs/detail?id=54
func (c *Client) AffiliateCreateWithShop(ctx context.Context, req *AffiliateCreateWithShopRequest) error {
	return c.sendRequest(ctx, "/v2/shop/affiliateCreateWithShop", c.shopID, req, nil)
}

func (c *Client) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	var resp CreateOrderResponse
	err := c.sendRequest(ctx, "/v2/shipping-order/create", c.shopID, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) UpdateOrder(ctx context.Context, req *UpdateOrderRequest) error {
	err := c.sendRequest(ctx, "/v2/shipping-order/update", c.shopID, req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateOrderCOD(ctx context.Context, req *UpdateOrderCODRequest) error {
	err := c.sendRequest(ctx, "/v2/shipping-order/updateCOD", c.shopID, req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) FindAvailableServices(ctx context.Context, req *FindAvailableServicesRequest) (*FindAvailableServicesResponse, error) {
	var (
		// Key serviceID
		mapFee      = make(map[Int]*CalculateFeeResponse)
		mapLeadTime = make(map[Int]Time)
		mapService  = make(map[Int]*ServiceInfo)
		mu          sync.Mutex
		wg          sync.WaitGroup
	)

	// Get all services
	getServicesReq := &GetServicesRequest{
		FromDistrict: Int(req.FromDistrictID),
		ToDistrict:   Int(req.ToDistrictID),
	}
	getServicesResp, err := c.GetServices(ctx, getServicesReq)
	if err != nil {
		return nil, err
	}
	services := *getServicesResp

	// Get all fee for services
	for _, service := range services {
		// Ignore service have empty name
		if service.ShortName == "" {
			continue
		}
		mapService[service.ServiceID] = service
		calculateFeeReq := &CalculateFeeRequest{
			ShopID:         c.shopID,
			ServiceID:      service.ServiceID.Int(),
			ServiceTypeID:  service.ServiceTypeID.Int(),
			InsuranceValue: req.InsuranceFee,
			Coupon:         req.Coupon,
			FromDistrictID: req.FromDistrictID,
			FromWardCode:   req.FromWardCode,
			ToDistrictID:   req.ToDistrictID,
			ToWardCode:     req.ToWardCode,
			Height:         req.Height,
			Length:         req.Length,
			Weight:         req.Weight,
			Width:          req.Width,
		}

		getLeadTimeReq := &GetLeadTimeRequest{
			FromDistrictID: req.FromDistrictID,
			FromWardCode:   req.FromWardCode,
			ToDistrictID:   req.ToDistrictID,
			ToWardCode:     req.ToWardCode,
			ServiceID:      service.ServiceID.Int(),
		}
		wg.Add(1)
		go func(serviceID Int, _calculateFeeReq *CalculateFeeRequest, _getLeadTimeReq *GetLeadTimeRequest) {
			defer wg.Done()

			calculateFeeResp, err := c.CalculateFee(ctx, _calculateFeeReq)
			if err != nil {
				return
			}

			getLeadTimeResp, err := c.GetLeadTime(ctx, _getLeadTimeReq)
			if err != nil {
				return
			}

			// trường hợp gói cước tính lỗi thì sẽ ko được thêm vào mapFee
			mu.Lock()
			if calculateFeeResp != nil {
				mapFee[serviceID] = calculateFeeResp
				mapLeadTime[serviceID] = getLeadTimeResp.LeadTime
			}
			mu.Unlock()
		}(service.ServiceID, calculateFeeReq, getLeadTimeReq)
	}
	wg.Wait()

	// map fee for service
	var availableServices []*AvailableService
	for serviceID, fee := range mapFee {
		service := mapService[serviceID]
		expectedDeliveryTime := mapLeadTime[serviceID]
		availableService := &AvailableService{
			Name:                 service.ShortName,
			ServiceFee:           fee.Total,
			ServiceID:            serviceID,
			ExpectedDeliveryTime: expectedDeliveryTime,
		}
		availableServices = append(availableServices, availableService)
	}

	return &FindAvailableServicesResponse{
		AvailableServices: availableServices,
	}, nil
}

func (c *Client) GetWards(ctx context.Context, req *GetWardsRequest) (*GetWardsResponse, error) {
	var resp GetWardsResponse
	err := c.sendRequest(ctx, "/master-data/ward", c.shopID, req, &resp)
	return &resp, err
}

func (c *Client) GetServices(ctx context.Context, req *GetServicesRequest) (*GetServicesResponse, error) {
	var resp GetServicesResponse
	err := c.sendRequest(ctx, "/pack-service/all", c.shopID, req, &resp)
	return &resp, err
}

func (c *Client) CalculateFee(ctx context.Context, req *CalculateFeeRequest) (*CalculateFeeResponse, error) {
	var resp CalculateFeeResponse
	err := c.sendRequest(ctx, "/v2/shipping-order/fee", req.ShopID, req, &resp)
	return &resp, err
}

func (c *Client) GetLeadTime(ctx context.Context, req *GetLeadTimeRequest) (*GetLeadTimeResponse, error) {
	var resp GetLeadTimeResponse
	err := c.sendRequest(ctx, "/v2/shipping-order/leadtime", c.shopID, req, &resp)
	return &resp, err
}

func (c *Client) GetOrderInfo(ctx context.Context, req *GetOrderInfoRequest) (*GetOrderInfoResponse, error) {
	var resp GetOrderInfoResponse
	err := c.sendRequest(ctx, "/v2/shipping-order/detail", c.shopID, req, &resp)
	return &resp, err
}

func (c *Client) CancelOrder(ctx context.Context, req *CancelOrderRequest) error {
	return c.sendRequest(ctx, "/v2/switch-status/cancel", c.shopID, req, nil)
}

func (c *Client) GetShopByClientOwner(ctx context.Context, req *GetShopByClientOwnerRequest) (*GetShopByClientOwnerResponse, error) {
	var resp GetShopByClientOwnerResponse
	err := c.sendRequest(ctx, "/v2/client/shops-by-client-owner", c.shopID, req, &resp)
	return &resp, err
}

func (c *Client) AddClientContract(ctx context.Context, req *AddClientContractRequest) error {
	return c.sendRequest(ctx, "/v2/contract/add-client", c.shopID, req, nil)
}

func (c *Client) sendRequest(ctx context.Context, path string, shopID int, req, resp interface{}) error {
	var errResp ErrorResponse
	res, err := c.rclient.R().
		SetHeader(Token, c.token).
		SetHeader(ShopID, fmt.Sprintf("%d", shopID)).
		SetBody(req).
		SetResult(&errResp).
		SetError(&errResp).
		Post(URL(c.baseUrl, path))
	if err != nil {
		return cm.Error(cm.ExternalServiceError, "Lỗi kết nối với GHN", err)
	}

	status := res.StatusCode()
	switch {
	case status >= 200 && status < 300:
		if resp != nil {
			if httpreq.IsNullJsonRaw(errResp.Data) {
				return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giao Hang Nhanh: null response. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
			}
			if err = jsonx.Unmarshal(errResp.Data, resp); err != nil {
				return cm.Errorf(cm.ExternalServiceError, err, "Lỗi không xác định từ Giao Hang Nhanh: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", err, wl.X(ctx).CSEmail)
			}
		}
		return nil

	case status >= 400:
		var meta map[string]string
		if !httpreq.IsNullJsonRaw(errResp.Data) {
			if err = jsonx.Unmarshal(errResp.Data, &meta); err != nil {
				// The slow path
				var metaX map[string]interface{}
				_ = jsonx.Unmarshal(errResp.Data, &metaX)
				meta = make(map[string]string)
				for k, v := range metaX {
					meta[k] = fmt.Sprint(v)
				}
			}
			errResp.ErrorData = meta
		}

		// Handle "An error occur"
		if strings.Contains(errResp.Message.String(), "error") {
			return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Giao Hang Nhanh: %v. Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %.", errResp.Error(), wl.X(ctx).CSEmail).WithMetaM(meta)
		}

		return cm.Errorf(cm.ExternalServiceError, &errResp, "Lỗi từ Giao Hang Nhanh: %v. Nếu cần thêm thông tin vui lòng liên hệ %v.", errResp.Error(), wl.X(ctx).CSEmail).WithMetaM(meta)

	default:
		return cm.Errorf(cm.ExternalServiceError, nil, "Lỗi không xác định từ Giao Hang Nhanh: Invalid status (%v). Chúng tôi đang liên hệ với Giao Hang Nhanh để xử lý. Xin lỗi quý khách vì sự bất tiện này. Nếu cần thêm thông tin vui lòng liên hệ %v.", status, wl.X(ctx).CSEmail)
	}
}
