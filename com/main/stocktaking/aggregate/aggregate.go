package aggregate

import (
	"context"
	"time"

	"etop.vn/api/main/etop"

	"etop.vn/capi"

	stocktake "etop.vn/api/main/stocktaking"
	"etop.vn/backend/com/main/stocktaking/convert"
	"etop.vn/backend/com/main/stocktaking/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/conversion"
)

var _ stocktake.Aggregate = &StocktakeAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type StocktakeAggregate struct {
	StocktakeStore sqlstore.ShopStocktakeFactory
	EventBus       capi.EventBus
	DB             *cmsql.Database
}

func NewAggregateStocktake(db *cmsql.Database, eventBus capi.EventBus) *StocktakeAggregate {
	return &StocktakeAggregate{
		StocktakeStore: sqlstore.NewStocktakeStore(db),
		EventBus:       eventBus,
		DB:             db,
	}
}

func (q *StocktakeAggregate) MessageBus() stocktake.CommandBus {
	b := bus.New()
	return stocktake.NewAggregateHandler(q).RegisterHandlers(b)
}

func (q *StocktakeAggregate) CreateStocktake(ctx context.Context, args *stocktake.CreateStocktakeRequest) (*stocktake.ShopStocktake, error) {
	if args.ShopID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing shop_id in request", nil)
	}
	var stockTake = &stocktake.ShopStocktake{}
	err := scheme.Convert(args, stockTake)
	if err != nil {
		return nil, err
	}
	InventoryMaxCode, err := q.StocktakeStore(ctx).ShopID(args.ShopID).GetStocktakeMaximumCodeNorm()
	var maxCodeNorm int32
	switch cm.ErrorCode(err) {
	case cm.NoError:
		maxCodeNorm = InventoryMaxCode.CodeNorm
	case cm.NotFound:
		// no-op
	default:
		return nil, err
	}
	if maxCodeNorm >= convert.MaxCodeNorm {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng nhập mã")
	}
	codeNorm := maxCodeNorm + 1
	stockTake.Code = convert.GenerateCode(int(codeNorm))
	stockTake.CodeNorm = codeNorm
	err = q.StocktakeStore(ctx).Create(stockTake)
	return stockTake, err
}

func (q *StocktakeAggregate) UpdateStocktake(ctx context.Context, args *stocktake.UpdateStocktakeRequest) (*stocktake.ShopStocktake, error) {
	if args.ShopID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing shop_id in request", nil)
	}
	stockTake, err := q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ShopID).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(args, stockTake)
	if err != nil {
		return nil, err
	}
	err = q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateAll(stockTake)
	return stockTake, err
}

func (q *StocktakeAggregate) ConfirmStocktake(ctx context.Context, args *stocktake.ConfirmStocktakeRequest) (*stocktake.ShopStocktake, error) {
	stocktakeDB, err := q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	if stocktakeDB.Status != etop.S3Zero {
		return nil, cm.Error(cm.InvalidArgument, "Stocktake đã được xác nhận hoặc hủy bỏ, Vui lòng kiểm tra lại", nil)
	}
	stocktakeDB.ConfirmedAt = time.Now()
	stocktakeDB.Status = etop.S3Positive
	err = q.DB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err = q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateAll(stocktakeDB)
		event := stocktake.StocktakeConfirmedEvent{
			Stocktake:            stocktakeDB,
			Overstock:            args.OverStock,
			ConfirmedBy:          args.ConfirmedBy,
			AutoInventoryVoucher: args.AutoInventoryVoucher,
		}
		err = q.EventBus.Publish(ctx, &event)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return stocktakeDB, nil
}

func (q *StocktakeAggregate) CancelStocktake(ctx context.Context, id int64, shopID int64) (*stocktake.ShopStocktake, error) {
	stocktakeDB, err := q.StocktakeStore(ctx).ShopID(shopID).ID(id).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	if stocktakeDB.Status != etop.S3Zero {
		return nil, cm.Error(cm.InvalidArgument, "Stocktake đã được xác nhận hoặc hủy bỏ, Vui lòng kiểm tra lại", nil)
	}
	stocktakeDB.CancelledAt = time.Now()
	stocktakeDB.Status = etop.S3Negative
	err = q.StocktakeStore(ctx).ShopID(shopID).ID(id).UpdateAll(stocktakeDB)
	return stocktakeDB, nil
}