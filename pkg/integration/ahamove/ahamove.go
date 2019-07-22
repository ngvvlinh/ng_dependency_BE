package ahamove

import (
	"context"
	"strings"
	"sync"

	"etop.vn/api/main/shipnow/carrier"
	shipnowtypes "etop.vn/api/main/shipnow/types"
	cm "etop.vn/backend/pkg/common"
	ahamoveclient "etop.vn/backend/pkg/integration/ahamove/client"
	"etop.vn/common/l"
)

var ll = l.New()

func (c *Carrier) CalcShippingFee(ctx context.Context, cmd *CalcShippingFeeCommand) error {
	type Result struct {
		Service *ShippingService
		Result  *ahamoveclient.CalcShippingFeeResponse
		Error   error
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex

	services := c.GetAvailableServices(ctx, cmd.Request.DeliveryPoints)
	if len(services) == 0 {
		return cm.Error(cm.ExternalServiceError, "ahamove: Không có gói cước phù hợp", nil)
	}
	wg.Add(len(services))
	for _, s := range services {
		go func(s *ShippingService) {
			defer wg.Done()
			req := *cmd.Request
			req.ServiceID = string(s.Code)
			resp, err := c.client.CalcShippingFee(ctx, &req)
			m.Lock()
			result := Result{
				s, resp, err,
			}
			results = append(results, result)
			m.Unlock()
		}(s)
	}
	wg.Wait()
	if len(results) == 0 {
		return cm.Error(cm.ExternalServiceError, "Lỗi từ ahamove: Không thể lấy thông tin gói cước dịch vụ", nil)
	}

	generator := newServiceIDGenerator(cmd.ArbitraryID)
	var res []*shipnowtypes.ShipnowService
	for _, result := range results {
		providerServiceID, err := generator.GenerateServiceID(result.Service.Code)
		if err != nil {
			return err
		}
		if result.Error != nil {
			continue
		}
		_r := ToShipnowService(result.Result, result.Service, providerServiceID)
		res = append(res, _r)
	}
	cmd.Result = res
	return nil
}

func (c *Carrier) CreateOrder(ctx context.Context, cmd *CreateOrderCommand) error {
	serviceID, err := parseServiceCode(cmd.ServiceID)
	if err != nil {
		return err
	}

	cmd.Request.ServiceID = string(serviceID)
	cmd.Result, err = c.client.CreateOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) GetOrder(ctx context.Context, cmd *GetOrderCommand) (err error) {
	cmd.Result, err = c.client.GetOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) CancelOrder(ctx context.Context, cmd *CancelOrderCommand) (err error) {
	err = c.client.CancelOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) GetServices(ctx context.Context, cmd *GetServiceCommand) error {
	services, err := c.client.GetServices(ctx, cmd.Request)
	cmd.Result = services
	return err
}

func (c *CarrierAccount) RegisterAccount(ctx context.Context, cmd *RegisterAccountCommand) (err error) {
	cmd.Result, err = c.client.RegisterAccount(ctx, cmd.Request)
	return err
}

func (c *CarrierAccount) GetAccount(ctx context.Context, cmd *GetAccountCommand) (err error) {
	cmd.Result, err = c.client.GetAccount(ctx, cmd.Request)
	return err
}

func (c *CarrierAccount) VerifyAccount(ctx context.Context, cmd *VerifyAccountCommand) (err error) {
	cmd.Result, err = c.client.VerifyAccount(ctx, cmd.Request)
	return err
}

func ToShipnowService(sfResp *ahamoveclient.CalcShippingFeeResponse, service *ShippingService, providerServiceID string) *shipnowtypes.ShipnowService {
	if sfResp == nil {
		return nil
	}

	res := &shipnowtypes.ShipnowService{
		Carrier:     carrier.Ahamove,
		Name:        service.Name,
		Code:        providerServiceID,
		Fee:         int32(sfResp.TotalFee),
		Description: service.Description,
	}
	// BIKE/POOL: discount, total_fee, total_pay
	// SAMEDAY: partner_discount, partner_fee, partner_pay
	// Ahamove đang fix, sau này sẽ dùng total_fee hết
	if strings.Contains(service.Code, string(SAMEDAY)) {
		res.Fee = int32(sfResp.PartnerFee)
	}

	// Avoid fee == 0
	if res.Fee == 0 {
		return nil
	}
	return res
}
