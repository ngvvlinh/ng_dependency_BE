package query

import (
	"context"

	"o.o/api/main/connectioning"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipping/carrier"
	"o.o/backend/com/main/shipping/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	ghnclientv2 "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	ghndriverv2 "o.o/backend/pkg/integration/shipping/ghn/driverv2"
	"o.o/capi/dot"
)

var _ shipping.QueryService = &QueryService{}

type QueryService struct {
	shimentManager *carrier.ShipmentManager
	connectionQS   connectioning.QueryBus

	store sqlstore.FulfillmentStoreFactory
}

func NewQueryService(
	db com.MainDB,
	shipmentManager *carrier.ShipmentManager, connectionQS connectioning.QueryBus,
) *QueryService {
	return &QueryService{
		shimentManager: shipmentManager,
		connectionQS:   connectionQS,
		store:          sqlstore.NewFulfillmentStore(db),
	}
}

func QueryServiceMessageBus(q *QueryService) shipping.QueryBus {
	b := bus.New()
	return shipping.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GetFulfillmentByIDOrShippingCode(ctx context.Context, args *shipping.GetFulfillmentByIDOrShippingCodeArgs) (*shipping.Fulfillment, error) {
	if args.ID == 0 && args.ShippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	query := q.store(ctx)
	if args.ID != 0 {
		query = query.ID(args.ID)
	}
	if args.ShippingCode != "" {
		query = query.ShippingCode(args.ShippingCode)
	}
	if len(args.ConnectionIDs) > 0 {
		query = query.ConnectionIDs(args.ConnectionIDs...)
	}
	return query.GetFulfillment()
}

func (q *QueryService) ListFulfillmentsByIDs(ctx context.Context, IDs []dot.ID, shopID dot.ID) ([]*shipping.Fulfillment, error) {
	if shopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shopID")
	}
	fulfillments, err := q.store(ctx).IDs(IDs...).ShopID(shopID).ListFfms()
	if err != nil {
		return nil, err
	}
	return fulfillments, nil
}

func (q *QueryService) ListFulfillmentsByMoneyTx(ctx context.Context, args *shipping.ListFullfillmentsByMoneyTxArgs) ([]*shipping.Fulfillment, error) {
	if len(args.MoneyTxShippingIDs) == 0 && args.MoneyTxShippingExternalID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing moneyTx ID")
	}

	query := q.store(ctx).OptionalMoneyTxShippingExternalID(args.MoneyTxShippingExternalID)
	query.WithPaging(meta.Paging{
		Limit: 10000,
	})
	if len(args.MoneyTxShippingIDs) > 0 {
		query = query.MoneyTxShippingIDs(args.MoneyTxShippingIDs...)
	}

	return query.ListFfms()
}

func (q *QueryService) GetFulfillmentExtended(ctx context.Context, id dot.ID, shippingCode string) (*shipping.FulfillmentExtended, error) {
	if id == 0 && shippingCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing id or shipping_code")
	}
	return q.store(ctx).OptionalID(id).OptionalShippingCode(shippingCode).GetFulfillmentExtended()
}

func (q *QueryService) ListFulfillmentExtendedsByIDs(ctx context.Context, ids []dot.ID, shopID dot.ID) ([]*shipping.FulfillmentExtended, error) {
	return q.store(ctx).IDs(ids...).OptionalShopID(shopID).ListFulfillmentExtendeds()
}

func (q *QueryService) ListFulfillmentExtendedsByMoneyTxShippingID(ctx context.Context, shopID dot.ID, moneyTxShippingID dot.ID) ([]*shipping.FulfillmentExtended, error) {
	return q.store(ctx).ShopID(shopID).MoneyTxShippingID(moneyTxShippingID).ListFulfillmentExtendeds()
}

func (q *QueryService) ListFulfillmentsByShippingCodes(ctx context.Context, codes []string) ([]*shipping.Fulfillment, error) {
	if len(codes) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping codes")
	}
	return q.store(ctx).ShippingCodes(codes).ListFfms()
}

func (q *QueryService) ListFulfillmentsForMoneyTx(ctx context.Context, args *shipping.ListFulfillmentForMoneyTxArgs) ([]*shipping.Fulfillment, error) {
	if args.ShippingProvider == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shipping provider")
	}
	if args.ShippingStates == nil && !args.IsNoneCOD.Valid {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing required arguments")
	}
	// Chỉ lấy những ffms chưa đối soát và chưa nằm trong phiên thanh toán nào
	query := q.store(ctx).ShippingProvider(args.ShippingProvider).NotBelongToMoneyTx()

	query = query.FilterForMoneyTx(args.IsNoneCOD.Bool, args.ShippingStates)

	return query.ListFfms()
}

func (q *QueryService) ListCustomerReturnRates(
	ctx context.Context, args *shipping.ListCustomerReturnRatesArgs,
) (customerRateExtendeds []*shipping.CustomerReturnRateExtended, _ error) {
	mConnections := make(map[dot.ID]*connectioning.Connection)

	listConnectionsQuery := &connectioning.ListConnectionsQuery{
		IDs:               args.ConnectionIDs,
		ConnectionType:    connection_type.Shipping,
		ConnectionSubtype: connection_type.ConnectionSubtypeShipment,
		Status:            status3.WrapStatus(status3.P),
	}
	if err := q.connectionQS.Dispatch(ctx, listConnectionsQuery); err != nil {
		return nil, err
	}

	for _, connection := range listConnectionsQuery.Result {
		mConnections[connection.ID] = connection
	}
	for _, connection := range mConnections {
		// handle fabo
		if connection.ConnectionProvider == connection_type.ConnectionProviderGHN &&
			connection.ConnectionMethod == connection_type.ConnectionMethodDirect && connection.Version == "v2" {
			shipmentDriver, err := q.shimentManager.GetDriverByEtopAffiliateAccount(ctx, connection.ID)
			if err != nil {
				return nil, err
			}
			ghnDriver := shipmentDriver.(*ghndriverv2.GHNDriver)
			ghnClient := ghnDriver.GetClient()

			etlCustomerRateReq := &ghnclientv2.CustomerReturnRateRequest{
				Phone: args.Phone,
			}
			customerReturnRateResp, err := ghnClient.CustomerReturnRate(ctx, etlCustomerRateReq)
			if err != nil {
				return nil, err
			}

			customerRateExtendeds = append(customerRateExtendeds, &shipping.CustomerReturnRateExtended{
				Connection: connection,
				CustomerReturnRate: &shipping.CustomerReturnRate{
					Level:     customerReturnRateResp.Level.String(),
					LevelCode: customerReturnRateResp.LevelCode.String(),
					Rate:      float64(customerReturnRateResp.Rate),
				},
			})
			return customerRateExtendeds, nil
		} else {
			// TODO(ngoc)
		}
	}

	return
}
