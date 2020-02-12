package vtpost

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"etop.vn/api/main/location"
	"etop.vn/api/top/types/etc/shipping_fee_type"
	"etop.vn/api/top/types/etc/shipping_provider"
	"etop.vn/api/top/types/etc/status5"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	shippingsharemodel "etop.vn/backend/com/main/shipping/sharemodel"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/shipping"
	vtpostclient "etop.vn/backend/pkg/integration/shipping/vtpost/client"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var ll = l.New()

const (
	SecretCode       = int64(1080)
	VTPostCodePublic = 'D'
)

func init() {
	model.GetShippingServiceRegistry().RegisterNameFunc(shipping_provider.VTPost, DecodeShippingServiceName)
}

func (c *Carrier) getClient(ctx context.Context, code byte) (vtpostclient.Client, error) {
	client := c.clients[code]
	if client != nil {
		// TODO: move to underlying goroutine
		changed, err := client.AutoLoginAndRefreshToken(ctx)
		if err != nil {
			return nil, err
		}
		if changed {
			if err = CreateShippingSource(VTPostCodePublic, client); err != nil {
				return nil, err
			}
		}
		return client, nil
	}

	if cmenv.IsDev() {
		return nil, cm.Error(cm.InvalidArgument, "DEVELOPMENT: No client for Vtpost", nil)
	}
	return nil, cm.Error(cm.InvalidArgument, "vtpost: invalid client code", nil)
}

func (c *Carrier) CalcShippingFee(ctx context.Context, cmd *CalcShippingFeeAllServicesArgs) error {
	type Result struct {
		Code   byte
		Result *vtpostclient.ShippingFeeService
		Error  error
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(len(c.clients))
	for code, client := range c.clients {
		go func(code byte, c vtpostclient.Client) {
			defer wg.Done()
			req := *cmd.Request // clone the request to prevent race condition
			resp, err := c.CalcShippingFeeAllServices(ctx, &req)
			m.Lock()
			for _, service := range resp {
				result := Result{code, service, err}
				results = append(results, result)
			}
			m.Unlock()
		}(code, client)
	}

	wg.Wait()
	if len(results) == 0 {
		return cm.Error(cm.ExternalServiceError, "Lỗi từ vtPost: không thể lấy thông tin gói cước dịch vụ", nil).
			WithMeta("reason", "timeout")
	}
	generator := newServiceIDGenerator(cmd.ArbitraryID.Int64())
	var res []*model.AvailableShippingService
	client, err := c.getClient(ctx, VTPostCodePublic)
	if err != nil {
		return err
	}
	for _, result := range results {
		// always generate service id, even if the result is error
		serviceCode := vtpostclient.VTPostOrderServiceCode(result.Result.MaDVChinh)
		providerServiceID, err := generator.GenerateServiceID(result.Code, serviceCode)
		if err != nil {
			continue
		}
		if result.Error != nil {
			continue
		}
		// ignore this service
		ignoreServices := []string{
			vtpostclient.OrderServiceCodeV60.String(),
		}
		if cm.StringsContain(ignoreServices, serviceCode.String()) {
			continue
		}

		// recall get price to get exactly shipping fee for each service
		query := &vtpostclient.CalcShippingFeeRequest{
			SenderProvince:   cmd.Request.SenderProvince,
			SenderDistrict:   cmd.Request.SenderDistrict,
			ReceiverProvince: cmd.Request.ReceiverProvince,
			ReceiverDistrict: cmd.Request.ReceiverDistrict,
			OrderService:     serviceCode,
			ProductWeight:    cmd.Request.ProductWeight,
			ProductPrice:     cmd.Request.ProductPrice,
			MoneyCollection:  cmd.Request.MoneyCollection,
		}
		resp, err := client.CalcShippingFee(ctx, query)
		if err != nil {
			continue
		}
		result.Result.GiaCuoc = resp.Data.MoneyTotal

		now := time.Now()
		expectedPickTime := shipping.CalcPickTime(shipping_provider.VTPost, now)
		thoigian := result.Result.ThoiGian // has format: "12 giờ"
		thoigian = strings.Replace(thoigian, " giờ", "", -1)
		hours, err := strconv.Atoi(thoigian)
		var expectedDeliveryDuration time.Duration
		if err != nil {
			expectedDeliveryDuration = CalcDeliveryDuration(serviceCode, cmd.FromProvince, cmd.ToProvince, cmd.FromDistrict, cmd.ToDistrict)
		} else {
			expectedDeliveryDuration = time.Duration(hours) * time.Hour
		}
		expectedDeliveryTime := expectedPickTime.Add(expectedDeliveryDuration)

		resItem := result.Result.ToAvailableShippingService(providerServiceID, expectedPickTime, expectedDeliveryTime)
		res = append(res, resItem)
	}
	res = shipping.CalcServicesTime(shipping_provider.VTPost, cmd.FromDistrict, cmd.ToDistrict, res)
	cmd.Result = res
	return nil
}

func (c *Carrier) GetShippingFeeLines(ctx context.Context, cmd *GetShippingFeeLinesCommand) error {
	clientCode, orderService, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}
	client, err := c.getClient(ctx, clientCode)
	if err != nil {
		return err
	}
	req := *cmd.Request
	req.OrderService = orderService
	res, err := client.CalcShippingFee(ctx, &req)
	if err != nil {
		return err
	}

	now := time.Now()
	expectedPickTime := shipping.CalcPickTime(shipping_provider.VTPost, now)
	expectedDeliveryDuration := CalcDeliveryDuration(orderService, cmd.FromProvince, cmd.ToProvince, cmd.FromDistrict, cmd.ToDistrict)
	expectedDeliveryTime := expectedPickTime.Add(expectedDeliveryDuration)
	lines, err := res.Data.CalcAndConvertShippingFeeLines()
	if err != nil {
		return err
	}

	cmd.Result = &GetShippingFeeLineResponse{
		ShippingFeeLines:   lines,
		ExpectedPickAt:     expectedPickTime,
		ExpectedDeliveryAt: expectedDeliveryTime,
	}
	return nil
}

func (c *Carrier) createOrder(ctx context.Context, cmd *CreateOrderArgs) error {
	clientCode, orderService, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(ctx, clientCode)
	if err != nil {
		return err
	}

	// detect transport from ServiceID
	req := *cmd.Request
	req.OrderService = orderService
	cmd.Result, err = client.CreateOrder(ctx, &req)
	return err
}

func (c *Carrier) cancelOrder(ctx context.Context, cmd *CancelOrderCommand) error {
	clientCode, _, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(ctx, clientCode)
	if err != nil {
		return err
	}
	cmd.Result, err = client.CancelOrder(ctx, cmd.Request)
	return err
}

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, orderMsg vtpostclient.CallbackOrderData) *shipmodel.Fulfillment {
	if !shipping.CanUpdateFulfillment(ffm) {
		return ffm
	}

	now := time.Now()
	data, _ := jsonx.Marshal(orderMsg)
	statusCode := orderMsg.OrderStatus
	vtpostStatus := vtpostclient.ToVTPostShippingState(statusCode)
	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: now,
		ExternalShippingData:      data,
		ExternalShippingState:     orderMsg.StatusName,
		ExternalShippingStateCode: strconv.Itoa(statusCode),
		ExternalShippingStatus:    vtpostStatus.ToStatus5(),
		ShippingState:             vtpostStatus.ToModel(ffm.ShippingState),
		EtopDiscount:              ffm.EtopDiscount,
		ShippingStatus:            vtpostStatus.ToShippingStatus5(ffm.ShippingState),
	}

	// Update price + weight
	if ffm.TotalWeight != orderMsg.ProductWeight {
		changeWeightNote := shipping.ChangeWeightNote(ffm.TotalWeight, orderMsg.ProductWeight)
		update.TotalWeight = orderMsg.ProductWeight
		update.AdminNote = ffm.AdminNote + "\n" + changeWeightNote
	}
	if ffm.ShippingFeeShop != orderMsg.MoneyTotal && shipping.CanUpdateFulfillmentFeelines(ffm) {
		// keep all shipping fee lines except shippingFeeMain
		mainFee := orderMsg.MoneyTotal
		for _, line := range ffm.ProviderShippingFeeLines {
			if line.ShippingFeeType == shipping_fee_type.Main {
				continue
			}
			mainFee = mainFee - line.Cost
		}
		if mainFee >= 0 {
			for _, line := range ffm.ProviderShippingFeeLines {
				if line.ShippingFeeType == shipping_fee_type.Main {
					line.Cost = mainFee
				}
			}
		}
		update.ProviderShippingFeeLines = ffm.ProviderShippingFeeLines
		update.ShippingFeeShopLines = shippingsharemodel.GetShippingFeeShopLines(update.ProviderShippingFeeLines, false, dot.NullInt{})
	}

	// Only update status5 if the current status is not ending status
	newStatus := vtpostStatus.ToStatus5()
	// UpdateInfo ClosedAt
	if newStatus == status5.N || newStatus == status5.NS || newStatus == status5.P {
		if ffm.ExternalShippingClosedAt.IsZero() {
			update.ClosedAt = now
		}
		if ffm.ClosedAt.IsZero() {
			update.ClosedAt = now
		}
	}
	return update
}

func CalcDeliveryDuration(orderService vtpostclient.VTPostOrderServiceCode, fromProvince, toProvince *location.Province, fromDistrict, toDistrict *location.District) (duration time.Duration) {
	serviceName := orderService.Name()
	switch serviceName {
	case model.ShippingServiceNameFaster:
		duration = CalcDeliveryDurationFastService(fromProvince, toProvince, fromDistrict, toDistrict)
	case model.ShippingServiceNameStandard:
		duration = CalcDeliveryDurationStandardService(fromProvince, toProvince, fromDistrict, toDistrict)
	}
	return duration
}

func CalcDeliveryDurationStandardService(fromProvince, toProvince *location.Province, fromDistrict, toDistrict *location.District) (duration time.Duration) {
	switch {
	// Nội tỉnh
	case fromProvince.Code == toProvince.Code:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 72 * time.Hour
		case location.Suburban1:
			duration = 96 * time.Hour
		default:
			duration = 120 * time.Hour
		}
		return duration

	// Nội miền
	case fromProvince.Region == toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 120 * time.Hour
		case location.Suburban1:
			duration = 144 * time.Hour
		default:
			duration = 168 * time.Hour
		}
		return duration

	// Khác miền
	case fromProvince.Region != toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 168 * time.Hour
		case location.Suburban1:
			duration = 192 * time.Hour
		default:
			duration = 216 * time.Hour
		}
		return duration
	// Tỉnh thành khác, nội miền, khác miền, Nhanh (không hỗ trợ gói Chuẩn)
	default:
		duration = 216 * time.Hour
		return duration
	}
}

// TODO: move back to location?
const (
	HCMProvinceCode       = "79"
	BinhDuongProvinceCode = "74"
	DongNaiProvinceCode   = "75"
	VungTauProvinceCode   = "77"
)

var groupProvinceCodes = []string{BinhDuongProvinceCode, DongNaiProvinceCode, VungTauProvinceCode}

func CalcDeliveryDurationFastService(fromProvince, toProvince *location.Province, fromDistrict, toDistrict *location.District) (duration time.Duration) {
	switch {
	// Nội tỉnh
	case fromProvince.Code == toProvince.Code:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 24 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 72 * time.Hour
		}
		return duration

	// HCM <=> Binh Duong, Dong Nai, Ba Ria Vung Tau
	case fromProvince.Code == HCMProvinceCode && cm.StringsContain(groupProvinceCodes, toProvince.Code) ||
		cm.StringsContain(groupProvinceCodes, fromProvince.Code) && toProvince.Code == HCMProvinceCode:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 72 * time.Hour
		}
		return duration

	// Nội miền
	case fromProvince.Region == toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 120 * time.Hour
		}
		return duration

	// HCM <=> HN; DN <=> HCM, HN
	case fromProvince.Region != toProvince.Region && fromProvince.Special && toProvince.Special:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 48 * time.Hour
		default:
			duration = 72 * time.Hour
		}
		return duration

	// Khác miền
	case fromProvince.Region != toProvince.Region:
		switch toDistrict.UrbanType {
		case location.Urban:
			duration = 48 * time.Hour
		case location.Suburban1:
			duration = 72 * time.Hour
		default:
			duration = 120 * time.Hour
		}
		return duration
	// Tỉnh thành khác, nội miền, khác miền, Nhanh (không hỗ trợ gói Chuẩn)
	default:
		duration = 120 * time.Hour
		return duration
	}
}
