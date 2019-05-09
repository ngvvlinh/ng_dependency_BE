package ghn

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	ghnclient "etop.vn/backend/pkg/integration/ghn/client"
	"etop.vn/backend/pkg/integration/shipping"
)

type ClientType byte

const (
	GHNCodeDefault ClientType = 'D'
	GHNCodeExt     ClientType = 'E' // unused
)

var ll = l.New()

func init() {
	model.GetShippingServiceRegistry().RegisterNameFunc(model.TypeGHN, DecodeShippingServiceName)
}

func (c ClientType) String() string {
	return string(c)
}

func GetShippingSourceName(c ClientType, clientID int) string {
	return c.String() + strconv.Itoa(clientID)
}

func (c *Carrier) GetClients() []*ghnclient.Client {
	res := make([]*ghnclient.Client, 0, len(c.clients))
	for _, c := range c.clients {
		res = append(res, c)
	}
	return res
}

func (c *Carrier) GetShippingSourceNames() []string {
	res := make([]string, 0, len(c.clients))
	for code, client := range c.clients {
		res = append(res, GetShippingSourceName(code, client.ClientID()))
	}
	return res
}

func (c *Carrier) getClient(code ClientType) (*ghnclient.Client, error) {
	client := c.clients[code]
	if client != nil {
		return client, nil
	}
	if cm.IsDev() {
		return nil, cm.Error(cm.InvalidArgument, "DEVELOPMENT: No client for GHN", nil)
	}
	return nil, cm.Error(cm.InvalidArgument, "GHN: invalid client code", nil).
		WithMetap("Client Code", code)
}

func (c *Carrier) CreateOrder(ctx context.Context, cmd *RequestCreateOrderCommand) error {
	client, serviceID, err := c.ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}
	cmd.Request.ServiceID = serviceID
	cmd.Result, err = client.CreateOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) FindAvailableServices(ctx context.Context, cmd *RequestFindAvailableServicesCommand) error {
	type Result struct {
		Code   ClientType
		Result *ghnclient.FindAvailableServicesResponse
		Error  error
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex
	var defaultServiceID int
	var defaultProviderServiceID string

	wg.Add(len(c.clients))
	for code, c := range c.clients {
		go func(code ClientType, c *ghnclient.Client) {
			defer wg.Done()
			req := *cmd.Request // clone the request to prevent race condition
			resp, err := c.FindAvailableServices(ctx, &req)
			result := Result{code, resp, err}
			m.Lock()
			results = append(results, result)
			m.Unlock()
		}(code, c)
	}
	wg.Wait()

	var res []*model.AvailableShippingService
	for _, result := range results {
		if result.Error != nil {
			ll.Error("Error requesting GHN", l.Error(result.Error))
			continue
		}
		for _, service := range result.Result.AvailableServices {
			if service.Name == "" {
				return cm.Errorf(cm.ExternalServiceError, nil, "Error requesting GHN: Gói cước không hợp lệ")
			}

			// Generate ProviderServiceID
			providerServiceID := GenerateShippingServiceCode(string(result.Code), service.Name.String(), service.ServiceID.String())
			res = append(res, service.ToShippingService(providerServiceID))
			defaultServiceID = int(service.ServiceID)
			defaultProviderServiceID = providerServiceID
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].ProviderServiceID < res[j].ProviderServiceID
	})

	if len(res) > 0 && cmd.Request.InsuranceFee > 0 {
		// Include insurance:
		//
		// "FindAvailableServices" does not include insurance in calculation,
		// therefore we must call "CalculateFee" to get insurance fee, then add
		// it to all previous calculation.
		calcFeeCmd := &ghnclient.CalculateFeeRequest{
			Connection:     ghnclient.Connection{},
			Weight:         cmd.Request.Weight,
			Length:         cmd.Request.Length,
			Width:          cmd.Request.Width,
			Height:         cmd.Request.Height,
			FromDistrictID: cmd.Request.FromDistrictID,
			ToDistrictID:   cmd.Request.ToDistrictID,
			InsuranceFee:   cmd.Request.InsuranceFee,
			ServiceID:      defaultServiceID,
		}
		client, _, err := c.ParseServiceID(defaultProviderServiceID)
		if err != nil {
			return err
		}
		calcFeeResult, err := client.CalculateFee(ctx, calcFeeCmd)
		if err != nil {
			return cm.Errorf(cm.ExternalServiceError, err, "Lỗi từ GHN: Không thể tính được phí giao hàng (%v)", err)
		}
		insuranceFee := ghnclient.GetInsuranceFee(calcFeeResult.OrderCosts)
		for _, shippingservice := range res {
			shippingservice.ServiceFee += insuranceFee
		}
	}
	res = shipping.CalcServicesTime(model.TypeGHN, cmd.FromDistrict, cmd.ToDistrict, res)
	cmd.Result = res
	return nil
}

func GenerateShippingServiceCode(clientCode string, serviceName string, serviceID string) string {
	shortCode := strings.ToUpper(string(serviceName[0]))
	return clientCode + shortCode + serviceID
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) < 6 {
		return "", false
	}
	switch {
	case code[1] == 'C': // Chuẩn
		return model.ShippingServiceNameStandard, true
	case code[1] == 'N': // Nhanh
		return model.ShippingServiceNameFaster, true
	}
	return "", false
}

func (c *Carrier) ParseServiceCode(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func (c *Carrier) CalculateFee(ctx context.Context, cmd *RequestCalculateFeeCommand) error {
	var client *ghnclient.Client
	var err error
	if cmd.ServiceID == "" {
		// case ffm.ProviderServiceID is null
		// Get default client
		client = c.clients[GHNCodeDefault]
	} else {
		client, _, err = c.ParseServiceID(cmd.ServiceID)
		if err != nil {
			return err
		}
	}
	cmd.Result, err = client.CalculateFee(ctx, cmd.Request)
	return err
}

func (c *Carrier) GetOrder(ctx context.Context, cmd *RequestGetOrderCommand) error {
	var client *ghnclient.Client
	var err error
	if cmd.ServiceID == "" {
		// case ffm.ProviderServiceID is null
		// Get default client
		client = c.clients[GHNCodeDefault]
	} else {
		client, _, err = c.ParseServiceID(cmd.ServiceID)
		if err != nil {
			return err
		}
	}
	cmd.Result, err = client.GetOrderInfo(ctx, cmd.Request)
	return err
}

func (c *Carrier) CancelOrder(ctx context.Context, cmd *RequestCancelOrderCommand) error {
	client, _, err := c.ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}
	err = client.CancelOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) ReturnOrder(ctx context.Context, cmd *RequestReturnOrderCommand) error {
	client, _, err := c.ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}
	err = client.ReturnOrder(ctx, cmd.Request)
	return err
}

func (c *Carrier) GetOrderLogs(ctx context.Context, cmd *RequestGetOrderLogsCommand) error {
	client, _, err := c.ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	cmd.Result, err = client.GetOrderLogs(ctx, cmd.Request)
	return err
}

func (c *Carrier) ParseServiceID(code string) (client *ghnclient.Client, serviceID int, err error) {
	if code == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "Missing service id")
		return
	}
	if len(code) <= 3 {
		err = cm.Errorf(cm.InvalidArgument, nil, "Invalid service id")
		return
	}
	clientCode := ClientType(code[0])
	var serviceIDStr string
	switch clientCode {
	case GHNCodeDefault, GHNCodeExt:
		serviceIDStr = code[1:]
		// also trim the second character, see GenerateShippingServiceCode
		if serviceIDStr[0] > '9' {
			serviceIDStr = serviceIDStr[1:]
		}

	default:
		// Backward compatible: The old service ids are not prefixed by
		// ClientCode. They have the format “12345”. So we test whether the
		// first character is a number and associate those ids with the default
		// client (GHNCodeDefault).
		if clientCode < '0' || clientCode > '9' {
			err = cm.Errorf(cm.InvalidArgument, nil, "Provider service id has wrong format")
			return nil, 0, err
		}
		clientCode = GHNCodeDefault
		serviceIDStr = code
	}

	client, err = c.getClient(clientCode)
	if err != nil {
		return nil, 0, err
	}
	serviceID, err = strconv.Atoi(serviceIDStr)
	return
}
