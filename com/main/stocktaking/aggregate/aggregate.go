package aggregate

import (
	"context"
	"time"

	stocktake "o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/stocktake_type"
	"o.o/backend/com/main/stocktaking/convert"
	"o.o/backend/com/main/stocktaking/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
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
	if args.Type != stocktake_type.Balance && args.Type != stocktake_type.Discard {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Loại quản lí tồn kho không đúng")
	}
	var stockTake = &stocktake.ShopStocktake{}
	err := scheme.Convert(args, stockTake)
	if err != nil {
		return nil, err
	}
	if args.Type == stocktake_type.Discard {
		for _, v := range args.Lines {
			if v.NewQuantity > v.OldQuantity || v.NewQuantity < 0 {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Số lượng xuất hủy không đúng")
			}
		}
	}
	InventoryMaxCode, err := q.StocktakeStore(ctx).ShopID(args.ShopID).GetStocktakeMaximumCodeNorm()
	var maxCodeNorm int
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
	stockTake.Code = convert.GenerateCode(codeNorm)
	stockTake.CodeNorm = codeNorm
	err = q.StocktakeStore(ctx).Create(stockTake)
	return stockTake, err
}

func (q *StocktakeAggregate) UpdateStocktake(ctx context.Context, args *stocktake.UpdateStocktakeRequest) (*stocktake.ShopStocktake, error) {
	if args.ShopID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing shop_id in request", nil)
	}
	_stockTake, err := q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	err = scheme.Convert(args, _stockTake)
	if err != nil {
		return nil, err
	}
	_stockTake.UpdatedAt = time.Now()
	err = q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateAll(_stockTake)
	return _stockTake, err
}

func (q *StocktakeAggregate) ConfirmStocktake(ctx context.Context, args *stocktake.ConfirmStocktakeRequest) (*stocktake.ShopStocktake, error) {
	stocktakeDB, err := q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	if stocktakeDB.Status != status3.Z {
		return nil, cm.Error(cm.InvalidArgument, "Stocktake đã được xác nhận hoặc hủy bỏ, Vui lòng kiểm tra lại", nil)
	}
	stocktakeDB.ConfirmedAt = time.Now()
	stocktakeDB.Status = status3.P
	err = q.DB.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err = q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateAll(stocktakeDB)
		event := stocktake.StocktakeConfirmedEvent{
			StocktakeID:          stocktakeDB.ID,
			ShopID:               stocktakeDB.ShopID,
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

func (q *StocktakeAggregate) CancelStocktake(ctx context.Context, args *stocktake.CancelStocktakeRequest) (*stocktake.ShopStocktake, error) {
	stocktake, err := q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).GetShopStocktake()
	if err != nil {
		return nil, err
	}
	if stocktake.Status != status3.Z {
		return nil, cm.Error(cm.InvalidArgument, "Stocktake đã được xác nhận hoặc hủy bỏ, Vui lòng kiểm tra lại", nil)
	}
	stocktake.CancelledAt = time.Now()
	stocktake.Status = status3.N
	stocktake.CancelReason = args.CancelReason
	err = q.StocktakeStore(ctx).ShopID(args.ShopID).ID(args.ID).UpdateAll(stocktake)
	return stocktake, nil
}
