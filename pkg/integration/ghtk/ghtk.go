package ghtk

import (
	"context"
	"encoding/json"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"etop.vn/api/main/location"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/etop/model"
	ghtkClient "etop.vn/backend/pkg/integration/ghtk/client"
	"etop.vn/backend/pkg/integration/shipping"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

var ll = l.New()

func (c *Carrier) getClient(code byte) (*ghtkClient.Client, error) {
	client := c.clients[code]
	if client != nil {
		return client, nil
	}

	if cm.IsDev() {
		return nil, cm.Error(cm.InvalidArgument, "DEVELOPMENT: No client for GHTK", nil)
	}
	return nil, cm.Error(cm.InvalidArgument, "ghtk: invalid client code", nil)
}

func (c *Carrier) CalcShippingFee(ctx context.Context, cmd *CalcShippingFeeCommand) error {

	fromQuery := &location.GetLocationQuery{DistrictCode: cmd.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: cmd.ToDistrictCode}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province

	type Result struct {
		Code      byte
		Transport ghtkClient.TransportType
		Result    *ghtkClient.CalcShippingFeeResponse
		Error     error
	}
	var results []Result
	var wg sync.WaitGroup
	var m sync.Mutex

	wg.Add(len(c.clients) * 2)
	for code, c := range c.clients {
		go func(code byte, c *ghtkClient.Client) {
			defer wg.Done()
			if code == 'S' || code == 'D' {
				// backward-compatible
				// These accounts ('S' & 'D') will be deleted later
				return
			}
			req := *cmd.Request // clone the request to prevent race condition
			req.Transport = ghtkClient.TransportRoad
			resp, err := c.CalcShippingFee(ctx, &req)
			m.Lock()
			result := Result{code, ghtkClient.TransportRoad, resp, err}
			results = append(results, result)
			m.Unlock()
		}(code, c)
		go func(code byte, c *ghtkClient.Client) {
			defer wg.Done()
			if code == 'S' || code == 'D' {
				// backward-compatible
				// These accounts ('S' & 'D') will be deleted later
				return
			}
			req := *cmd.Request // clone the request to prevent race condition
			// trường hợp nội tỉnh: có gói nhanh
			// trường hợp nội vùng: bỏ qua gói nhanh
			if fromProvince.Code != toProvince.Code && fromProvince.Region == toProvince.Region {
				return
			}
			req.Transport = ghtkClient.TransportFly
			resp, err := c.CalcShippingFee(ctx, &req)
			m.Lock()
			result := Result{code, ghtkClient.TransportFly, resp, err}
			results = append(results, result)
			m.Unlock()
		}(code, c)
	}

	wg.Wait()
	if len(results) == 0 {
		return cm.Error(cm.ExternalServiceError, "Lỗi từ Giaohangtietkiem: không thể lấy thông tin gói cước dịch vụ", nil).
			WithMeta("reason", "timeout")
	}
	// Sort result for stable service id generating. This must run before generating service id,
	sort.Slice(results, func(i, j int) bool {
		if results[i].Code < results[j].Code {
			return true
		}
		if results[i].Code > results[j].Code {
			return false
		}
		return results[i].Transport < results[j].Transport
	})
	generator := newServiceIDGenerator(cmd.ArbitraryID)

	now := time.Now()
	expectedPickTime := shipping.CalcPickTime(model.TypeGHTK, now)
	var res []*model.AvailableShippingService
	for _, result := range results {
		// always generate service id, even if the result is error
		providerServiceID, err := generator.GenerateServiceID(result.Code, result.Transport)
		if err != nil {
			return err
		}
		if result.Error != nil {
			continue
		}

		expectedDeliveryDuration := CalcDeliveryDuration(result.Transport, fromProvince, toProvince)
		expectedDeliveryTime := expectedPickTime.Add(expectedDeliveryDuration)
		resItem := result.Result.Fee.ToShippingService(providerServiceID,
			result.Transport, expectedPickTime, expectedDeliveryTime)
		res = append(res, resItem)
	}
	res = shipping.CalcServicesTime(model.TypeGHTK, fromDistrict, toDistrict, res)
	cmd.Result = res
	return nil
}

func (c *Carrier) CalcSingleShippingFee(ctx context.Context, cmd *CalcSingleShippingFeeCommand) error {

	fromQuery := &location.GetLocationQuery{DistrictCode: cmd.FromDistrictCode}
	toQuery := &location.GetLocationQuery{DistrictCode: cmd.ToDistrictCode}
	if err := c.location.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return err
	}
	fromDistrict, fromProvince := fromQuery.Result.District, fromQuery.Result.Province
	toDistrict, toProvince := toQuery.Result.District, toQuery.Result.Province

	clientCode, transport, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(clientCode)
	if err != nil {
		return err
	}
	resp, err := client.CalcShippingFee(ctx, cmd.Request)
	if err != nil {
		return err
	}
	now := time.Now()
	expectedPickTime := shipping.CalcPickTime(model.TypeGHTK, now)
	expectedDeliveryDuration := CalcDeliveryDuration(transport, fromProvince, toProvince)
	expectedDeliveryTime := expectedPickTime.Add(expectedDeliveryDuration)
	service := resp.Fee.ToShippingService(cmd.ServiceID,
		transport, expectedPickTime, expectedDeliveryTime)
	cmd.Result = shipping.CalcServiceTime(model.TypeGHTK, fromDistrict, toDistrict, service)

	return nil
}

func (c *Carrier) CreateOrder(ctx context.Context, cmd *CreateOrderCommand) error {
	clientCode, transport, err := ParseServiceID(cmd.ServiceID)
	if err != nil {
		return err
	}

	client, err := c.getClient(clientCode)
	if err != nil {
		return err
	}

	// detect transport from ServiceID
	cmd.Request.Order.Transport = transport
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
	cmd.Result, err = client.GetOrder(ctx, cmd.LabelID, cmd.PartnerID)
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
	cmd.Result, err = client.CancelOrder(ctx, cmd.LabelID, "")
	return err
}

func CalcUpdateFulfillment(ffm *shipmodel.Fulfillment, msg *ghtkClient.CallbackOrder, ghtkOrder *ghtkClient.OrderInfo) *shipmodel.Fulfillment {
	if !shipping.CanUpdateFulfillmentFromWebhook(ffm) {
		return ffm
	}

	now := time.Now()
	data, _ := json.Marshal(ghtkOrder)
	var statusID int
	if msg == nil {
		statusID, _ = strconv.Atoi(ghtkOrder.Status.String())
	} else {
		statusID = int(msg.StatusID)
	}
	stateID := ghtkClient.StateID(statusID)
	update := &shipmodel.Fulfillment{
		ID:                        ffm.ID,
		ExternalShippingUpdatedAt: now,
		ExternalShippingData:      data,
		ExternalShippingState:     ghtkClient.StateMapping[stateID],
		ExternalShippingStatus:    stateID.ToStatus5(),
		ExternalShippingCode:      ghtkOrder.LabelID.String(),
		ShippingState:             stateID.ToModel(),
		EtopDiscount:              ffm.EtopDiscount,
		ShippingStatus:            stateID.ToStatus5(),
	}

	// make sure can not update ffm's shipping fee when it belong to a money traction
	if ffm.MoneyTransactionShippingExternalID == 0 {
		update.ProviderShippingFeeLines = CalcAndConvertShippingFeeLines(ghtkOrder)
		var shippingFeeShopLines []*model.ShippingFeeLine
		shippingFeeShopLines = model.GetShippingFeeShopLines(update.ProviderShippingFeeLines, ffm.EtopPriceRule, &ffm.EtopAdjustedShippingFeeMain)
		shippingFeeShop := 0
		for _, line := range shippingFeeShopLines {
			shippingFeeShop += int(line.Cost)
		}
		update.ShippingFeeShopLines = shippingFeeShopLines
		update.ShippingFeeShop = shipmodel.CalcShopShippingFee(shippingFeeShop, ffm)
	}

	// Only update status4 if the current status is not ending status
	newStatus := stateID.ToStatus5()
	// UpdateInfo ClosedAt
	if newStatus == model.S5Negative || newStatus == model.S5NegSuper || newStatus == model.S5Positive {
		if ffm.ExternalShippingClosedAt.IsZero() {
			update.ClosedAt = now
		}
		if ffm.ClosedAt.IsZero() {
			update.ClosedAt = now
		}
	}
	return update
}

func CalcAndConvertShippingFeeLines(order *ghtkClient.OrderInfo) []*model.ShippingFeeLine {
	var res []*model.ShippingFeeLine
	insuranceFee := int(order.Insurance)
	fee := int(order.ShipMoney)
	shippingFeeMain := fee - insuranceFee

	// shipping fee
	res = append(res, &model.ShippingFeeLine{
		ShippingFeeType:      model.ShippingFeeTypeMain,
		Cost:                 shippingFeeMain,
		ExternalShippingCode: order.LabelID.String(),
	})
	// insurance fee
	if insuranceFee > 0 {
		res = append(res, &model.ShippingFeeLine{
			ShippingFeeType:      model.ShippingFeeTypeInsurance,
			Cost:                 insuranceFee,
			ExternalShippingCode: order.LabelID.String(),
		})
	}
	return res
}

type serviceIDGenerator struct {
	rd *rand.Rand
}

func newServiceIDGenerator(seed int64) serviceIDGenerator {
	src := rand.NewSource(seed)
	rd := rand.New(src)
	return serviceIDGenerator{rd}
}

// GenerateServiceID generate new service id for using with ghtk. The generated
// id is always 8 character in length.
func (c serviceIDGenerator) GenerateServiceID(clientCode byte, transport ghtkClient.TransportType) (string, error) {
	n := c.rd.Uint64()
	v := gencode.Alphabet32.EncodeReverse(n, 8)
	v = v[:8]

	switch clientCode {
	case 'D':
		v[1] = 'D'
		v[2] = blacklist(v[2], 'S', 'T', 'D')
		v[3] = blacklist(v[3], 'S', 'T', 'D')
	case 'S':
		v[2] = 'S'
		v[1] = blacklist(v[1], 'D', 'T', 'S')
		v[3] = blacklist(v[3], 'D', 'T', 'S')
	case 'T':
		v[3] = 'T'
		v[1] = blacklist(v[1], 'D', 'S', 'T')
		v[2] = blacklist(v[2], 'D', 'S', 'T')
	default:
		return "", cm.Errorf(cm.Internal, nil, "invalid code")
	}

	switch transport {
	case ghtkClient.TransportRoad:
		v[5] = 'R'
		v[6] = blacklist(v[6], 'R', 'F')
	case ghtkClient.TransportFly:
		v[6] = 'F'
		v[5] = blacklist(v[5], 'R', 'F')
	default:
		return "", cm.Errorf(cm.Internal, nil, "invalid transport")
	}

	return string(v), nil
}

func DecodeShippingServiceName(code string) (name string, ok bool) {
	if len(code) != 8 {
		return "", false
	}
	switch {
	case code[5] == 'R': // road
		return model.ShippingServiceNameStandard, true
	case code[6] == 'F': // fly
		return model.ShippingServiceNameFaster, true
	}
	return "", false
}

func init() {
	model.GetShippingServiceRegistry().RegisterNameFunc(model.TypeGHTK, DecodeShippingServiceName)
}

func (c *Carrier) ParseServiceCode(code string) (serviceName string, ok bool) {
	return DecodeShippingServiceName(code)
}

func blacklist(current byte, blacks ...byte) byte {
	for _, b := range blacks {
		if current == b {
			// return an arbitrary character which does not collide with blacklist values
			return 'J'
		}
	}
	return current
}

func ParseServiceID(id string) (clientCode byte, transport ghtkClient.TransportType, err error) {
	if id == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "missing service id")
		return
	}

	// old id
	if strings.Contains(id, ".") {
		parts := strings.Split(id, ".")
		switch parts[0] {
		case "GHTKPublic":
			clientCode = 'D'
		case "GHTKSamePrice35":
			clientCode = 'S'
		default:
			err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
			return
		}
		return clientCode, ghtkClient.TransportType(parts[1]), nil
	}

	if len(id) != 8 {
		err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
		return
	}

	// TODO: refactor code
	if id[1] == 'D' {
		clientCode = 'D'
	}
	if id[2] == 'S' {
		// make sure that we don't overwrite the client code
		if clientCode != 0 {
			err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
		}
		clientCode = 'S'
	}
	if id[3] == 'T' {
		// make sure that we don't overwrite the client code
		if clientCode != 0 {
			err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
		}
		clientCode = 'T'
	}
	if id[5] == 'R' {
		transport = ghtkClient.TransportRoad
	}
	if id[6] == 'F' {
		if transport != "" {
			err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
		}
		transport = ghtkClient.TransportFly
	}
	if clientCode == 0 || transport == "" {
		err = cm.Errorf(cm.InvalidArgument, nil, "invalid service id")
	}
	return
}

func CalcDeliveryDuration(transport ghtkClient.TransportType, from, to *location.Province) time.Duration {
	switch {
	// Nội tỉnh
	case from.Code == to.Code:
		return 6 * time.Hour

	// HN, HCM, ĐN, Nội miền
	case from.Region == to.Region && from.Extra.Special:
		return 24 * time.Hour

	// HN, HCM, ĐN, Khác miền, Đặc biệt
	case from.Region != to.Region && from.Extra.Special && to.Extra.Special:
		if transport == ghtkClient.TransportFly {
			return 24 * time.Hour
		} else {
			return 96 * time.Hour
		}

	// HN, HCM, ĐH, Khác miền
	case from.Region != to.Region && from.Extra.Special:
		if transport == ghtkClient.TransportFly {
			return 48 * time.Hour
		} else {
			return 120 * time.Hour
		}

	// Tỉnh thành khác, nội miền, khác miền, Nhanh (không hỗ trợ gói Chuẩn)
	default:
		return 48 * time.Hour
	}
}

func SyncOrders(ffms []*shipmodel.Fulfillment) ([]*shipmodel.Fulfillment, error) {
	rate := time.Second / 30
	burstLimit := 30
	ctx := context.Background()
	tick := time.NewTicker(rate)
	defer tick.Stop()
	throttle := make(chan time.Time, burstLimit)
	go func() {
		for t := range tick.C {
			select {
			case throttle <- t:
			default:
			}
		}
	}()
	ch := make(chan error, burstLimit)
	ll.Info("Length GHTK SyncOrders", l.Int("len", len(ffms)))
	var _ffms []*shipmodel.Fulfillment
	count := 0
	for _, ffm := range ffms {
		<-throttle
		count++
		if count > burstLimit {
			time.Sleep(20 * time.Second)
			count = 0
		}
		go func(ffm *shipmodel.Fulfillment) (_err error) {
			defer func() {
				ch <- _err
			}()
			// get order info to update service fee
			ghtkCmd := &GetOrderCommand{
				ServiceID: ffm.ProviderServiceID,
				LabelID:   ffm.ShippingCode,
			}
			if ghtkErr := bus.Dispatch(ctx, ghtkCmd); ghtkErr != nil {
				ll.Error("GHTK get order error :: ", l.String("shipping_code", ffm.ShippingCode), l.Error(ghtkErr))
				return ghtkErr
			}
			updateFfm := CalcUpdateFulfillment(ffm, nil, &ghtkCmd.Result.Order)
			_ffms = append(_ffms, updateFfm)
			return nil
		}(ffm)
	}
	var successCount, errorCount int
	for i, l := 0, len(ffms); i < l; i++ {
		err := <-ch
		if err == nil {
			successCount++
		} else {
			errorCount++
		}
	}
	ll.S.Infof("Sync fulfillments GHTK info success: %v/%v, errors %v/%v", successCount, len(ffms), errorCount, len(ffms))
	return _ffms, nil
}
