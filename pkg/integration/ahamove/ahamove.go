package ahamove

import (
	"context"
	"sync"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	ahamoveclient "etop.vn/backend/pkg/integration/ahamove/client"
)

var ll = l.New()

func (c *Carrier) getClient(code byte) (*ahamoveclient.Client, error) {
	client := c.clients[code]
	if client != nil {
		return client, nil
	}

	if cm.IsDev() {
		return nil, cm.Error(cm.InvalidArgument, "DEVELOPMENT: No client for ahamove", nil)
	}
	return nil, cm.Error(cm.InvalidArgument, "ahamove: invalid client code", nil)
}

func (c *Carrier) CalcShippingFee(ctx context.Context, cmd *CalcShippingFeeCommand) error {
	type Result struct {
		Code      byte
		ServiceID ServiceCode
		Result    *ahamoveclient.CalcShippingFeeResponse
		Error     error
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex

	services := GetAvailableServices(cmd.Request.DeliveryPoints)
	if len(services) == 0 {
		return cm.Error(cm.ExternalServiceError, "ahamove: Không có gói cước phù hợp", nil)
	}
	wg.Add(len(c.clients) * len(services))
	for code, c := range c.clients {
		for _, s := range services {
			go func(code byte, s *ShippingService, c *ahamoveclient.Client) {
				defer wg.Done()
				req := *cmd.Request
				req.ServiceID = string(s.ID)
				resp, err := c.CalcShippingFee(ctx, &req)
				m.Lock()
				result := Result{
					code, s.ID, resp, err,
				}
				results = append(results, result)
				m.Unlock()
			}(code, s, c)
		}
	}
	wg.Wait()
	if len(results) == 0 {
		return cm.Error(cm.ExternalServiceError, "Lỗi từ ahamove: Không thể lấy thông tin gói cước dịch vụ", nil)
	}

	generator := newServiceIDGenerator(cmd.ArbitraryID)
	var res []*model.AvailableShippingService
	for _, result := range results {
		providerServiceID, err := generator.GenerateServiceID(result.Code, result.ServiceID)
		if err != nil {
			return err
		}
		if result.Error != nil {
			continue
		}
		_r := ToShippingService(result.Result, result.ServiceID, providerServiceID)
		res = append(res, _r)
	}
	cmd.Result = res
	return nil
}

func (c *Carrier) CalcSingleShippingFee(ctx context.Context, cmd *CalcSingleShippingFeeCommand) error {
	// TODO
	return cm.ErrTODO
}

func (c *Carrier) CreateOrder(ctx context.Context, cmd *CreateOrderCommand) error {
	clientCode, serviceID, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(clientCode)
	if err != nil {
		return err
	}

	// detect transport from ServiceID
	cmd.Request.ServiceID = string(serviceID)
	cmd.Result, err = client.CreateOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) GetOrder(ctx context.Context, cmd *GetOrderCommand) error {
	clientCode, _, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(clientCode)
	if err != nil {
		return err
	}
	cmd.Result, err = client.GetOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) CancelOrder(ctx context.Context, cmd *CancelOrderCommand) error {
	clientCode, _, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(clientCode)
	if err != nil {
		return err
	}
	err = client.CancelOrder(ctx, cmd.Request)
	return err
}
