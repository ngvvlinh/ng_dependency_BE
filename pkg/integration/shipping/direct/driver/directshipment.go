package driver

import (
	"context"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/location"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/status5"
	carriertypes "o.o/backend/com/main/shipping/carrier/types"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	directclient "o.o/backend/pkg/integration/shipping/direct/client"
)

var (
	defautDrivers = []string{
		"shipping/shipment/direct/partner",
	}
)

var _ carriertypes.ShipmentCarrier = &DirectShipmentDriver{}

type DirectShipmentDriver struct {
	conn       *connectioning.Connection
	client     *directclient.Client
	locationQS location.QueryBus
}

func New(locationQS location.QueryBus, cfg directclient.PartnerAccountCfg) (*DirectShipmentDriver, error) {
	client, err := directclient.New(cfg)
	if err != nil {
		return nil, err
	}
	return &DirectShipmentDriver{
		conn:       cfg.Connection,
		client:     client,
		locationQS: locationQS,
	}, nil
}

func (d *DirectShipmentDriver) Ping(context.Context) error {
	if err := d.client.Ping(); err != nil {
		return cm.Errorf(cm.ExternalServiceError, err, "Can not init %v client", d.conn.Name)
	}
	return nil
}

func (d *DirectShipmentDriver) GetAffiliateID() string {
	return ""
}

func (d *DirectShipmentDriver) CreateFulfillment(
	ctx context.Context,
	ffm *shipmodel.Fulfillment,
	args *carriertypes.GetShippingServicesArgs,
	service *shippingsharemodel.AvailableShippingService) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	note := carriertypes.GetShippingProviderNote(ffm)

	fromQuery := &location.GetLocationQuery{
		DistrictCode: ffm.AddressFrom.DistrictCode,
		WardCode:     ffm.AddressFrom.WardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: ffm.AddressTo.DistrictCode,
		WardCode:     ffm.AddressTo.WardCode,
	}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromAddress := fromQuery.Result
	toAddress := toQuery.Result

	var lines []*directclient.ItemLine
	for _, line := range ffm.Lines {
		lines = append(lines, &directclient.ItemLine{
			Name:     line.ProductName,
			Price:    line.ListPrice,
			Quantity: line.Quantity,
		})
	}

	cmd := &directclient.CreateFulfillmentRequest{
		PickupAddress: directclient.SimpleAddress{
			FullName: ffm.AddressFrom.GetFullName(),
			Phone:    ffm.AddressFrom.Phone,
			Email:    ffm.AddressFrom.Email,
			Address1: ffm.AddressFrom.Address1,
			Address2: ffm.AddressFrom.Address2,
			Province: fromAddress.Province.Name,
			District: fromAddress.District.Name,
			Ward:     fromAddress.Ward.Name,
		},
		ShippingAddress: directclient.SimpleAddress{
			FullName: ffm.AddressTo.GetFullName(),
			Phone:    ffm.AddressTo.Phone,
			Email:    ffm.AddressTo.Email,
			Address1: ffm.AddressTo.Address1,
			Address2: ffm.AddressTo.Address2,
			Province: toAddress.Province.Name,
			District: toAddress.District.Name,
			Ward:     toAddress.Ward.Name,
		},
		Lines:               lines,
		TotalWeight:         args.ChargeableWeight,
		BasketValue:         ffm.BasketValue,
		TotalCODAmount:      ffm.TotalCODAmount,
		CODAmount:           ffm.TotalCODAmount,
		ShippingNote:        note,
		IncludeInsurance:    ffm.IncludeInsurance,
		ShippingServiceCode: service.ProviderServiceID,
		ShippingFee:         service.ServiceFee,
	}

	// `service.ServiceFee`: giá đã makeup của TopShip
	// Khi tạo đơn qua NVC, gửi giá gốc của NVC lên để bên họ kiểm tra giá trước khi tạo
	if service.ShipmentPriceInfo != nil {
		cmd.ShippingFee = service.ShipmentPriceInfo.OriginFee
	}

	resp, err := d.client.CreateFulfillment(ctx, cmd)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	updateFfm := &shipmodel.Fulfillment{
		ID:                ffm.ID,
		ProviderServiceID: service.ProviderServiceID,
		Status:            status5.S, // Now processing

		ShippingFeeShop: resp.ShippingFee.Int(),

		ShippingCode:              resp.ShippingCode.String(),
		ExternalShippingName:      service.Name,
		ExternalShippingID:        resp.FulfillmentID.String(),
		ExternalShippingCode:      resp.ShippingCode.String(),
		ExternalShippingCreatedAt: now,
		ExternalShippingUpdatedAt: now,
		ShippingCreatedAt:         now,
		ExternalShippingFee:       resp.ShippingFee.Int(),
		ShippingState:             shippingstate.Created,
		SyncStatus:                status4.P,
		SyncStates: &shippingsharemodel.FulfillmentSyncStates{
			SyncAt:    now,
			TrySyncAt: now,
		},
		ExpectedPickAt:           service.ExpectedPickAt,
		ExpectedDeliveryAt:       service.ExpectedDeliveryAt,
		ProviderShippingFeeLines: toShippingFeeLines(resp.ShippingFeeLines),
		ShippingFeeShopLines:     toShippingFeeLines(resp.ShippingFeeLines),
	}
	return updateFfm, nil
}

func (d *DirectShipmentDriver) UpdateFulfillment(context.Context, *shipmodel.Fulfillment) (ffmToUpdate *shipmodel.Fulfillment, _ error) {
	return nil, cm.Errorf(cm.ExternalServiceError, nil, "This carrier does not support this method")
}

func (d *DirectShipmentDriver) CancelFulfillment(ctx context.Context, ffm *shipmodel.Fulfillment) error {
	externalFfmID := ffm.ExternalShippingID
	shippingCode := ffm.ExternalShippingCode
	cmd := &directclient.CancelFulfillmentRequest{
		FulfillmentID: externalFfmID,
		ShippingCode:  shippingCode,
	}
	return d.client.CancelFulfillment(ctx, cmd)
}

func (d *DirectShipmentDriver) GetShippingServices(ctx context.Context, args *carriertypes.GetShippingServicesArgs) ([]*shippingsharemodel.AvailableShippingService, error) {
	// if args.FromWardCode == "" {
	// 	return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ gửi hàng - phường/xã không được để trống!")
	// }
	// if args.ToWardCode == "" {
	// 	return nil, cm.Errorf(cm.InvalidArgument, nil, "Địa chỉ nhận hàng - phường/xã không được để trống!")
	// }

	fromQuery := &location.GetLocationQuery{
		DistrictCode: args.FromDistrictCode,
		WardCode:     args.FromWardCode,
	}
	toQuery := &location.GetLocationQuery{
		DistrictCode: args.ToDistrictCode,
		WardCode:     args.ToWardCode,
	}
	if err := d.locationQS.DispatchAll(ctx, fromQuery, toQuery); err != nil {
		return nil, err
	}
	fromDistrict, fromProvince, fromWard := fromQuery.Result.District, fromQuery.Result.Province, fromQuery.Result.Ward
	toDistrict, toProvince, toWard := toQuery.Result.District, toQuery.Result.Province, toQuery.Result.Ward

	cmd := &directclient.GetShippingServicesRequest{
		BasketValue: args.BasketValue,
		TotalWeight: args.ChargeableWeight,
		PickupAddress: directclient.SimpleAddress{
			Province: fromProvince.Name,
			District: fromDistrict.Name,
		},
		ShippingAddress: directclient.SimpleAddress{
			Province: toProvince.Name,
			District: toDistrict.Name,
		},
		IncludeInsurance: args.IncludeInsurance,
		TotalCODAmount:   args.CODAmount,
		CODAmount:        args.CODAmount,
	}
	if fromWard != nil {
		cmd.PickupAddress.Ward = fromWard.Name
	}
	if toWard != nil {
		cmd.ShippingAddress.Ward = toWard.Name
	}
	carrierServices, err := d.client.GetShippingServices(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return toAvailableShippingServices(carrierServices), nil
}

func (d *DirectShipmentDriver) GetServiceName(code string) (serviceName string, ok bool) {
	return "", false
}

func (d *DirectShipmentDriver) GetMaxValueFreeInsuranceFee() int {
	return 0
}

func (d *DirectShipmentDriver) SignIn(ctx context.Context, args *carriertypes.SignInArgs) (*carriertypes.AccountResponse, error) {
	cmd := &directclient.SignInRequest{
		Email:    args.Email,
		Password: args.Password,
	}
	resp, err := d.client.SignIn(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &carriertypes.AccountResponse{
		Token:  resp.Token.String(),
		UserID: resp.UserID.String(),
	}, nil
}

func (d *DirectShipmentDriver) SignUp(ctx context.Context, args *carriertypes.SignUpArgs) (*carriertypes.AccountResponse, error) {
	cmd := &directclient.SignUpRequest{
		Phone:    args.Phone,
		Email:    args.Email,
		Password: args.Password,
		Name:     args.Name,
	}
	resp, err := d.client.SignUp(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &carriertypes.AccountResponse{
		Token:  resp.Token.String(),
		UserID: resp.UserID.String(),
	}, nil
}

func (d *DirectShipmentDriver) ParseServiceID(code string) (serviceID string, err error) {
	// Giữ nguyên serviceID của đối tác
	// Có thể bổ sung rule sau
	return code, nil
}
