package query

import (
	"context"

	"o.o/api/main/shipping"
	"o.o/api/meta"
	"o.o/backend/com/main/shipping/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ shipping.QueryService = &QueryService{}

type QueryService struct {
	store sqlstore.FulfillmentStoreFactory
}

func NewQueryService(db *cmsql.Database) *QueryService {
	return &QueryService{
		store: sqlstore.NewFulfillmentStore(db),
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
