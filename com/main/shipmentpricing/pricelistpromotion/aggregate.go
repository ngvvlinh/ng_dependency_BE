package pricelistpromotion

import (
	"context"
	"time"

	"o.o/api/main/shipmentpricing/pricelistpromotion"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion/convert"
	"o.o/backend/com/main/shipmentpricing/pricelistpromotion/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var _ pricelistpromotion.Aggregate = &Aggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type Aggregate struct {
	db                      *cmsql.Database
	priceListPromotionStore sqlstore.PriceListStorePromotionFactory
}

func NewAggregate(db com.MainDB) *Aggregate {
	return &Aggregate{
		db:                      db,
		priceListPromotionStore: sqlstore.NewPriceListStorePromotion(db),
	}
}

func AggregateMessageBus(a *Aggregate) pricelistpromotion.CommandBus {
	b := bus.New()
	return pricelistpromotion.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *Aggregate) CreatePriceListPromotion(ctx context.Context, args *pricelistpromotion.CreatePriceListPromotionArgs) (*pricelistpromotion.ShipmentPriceListPromotion, error) {
	if args.ConnectionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing connection ID")
	}
	if args.PriceListID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing price list ID")
	}
	if args.DateFrom.IsZero() || args.DateTo.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "DateFrom, DateTo does not valid")
	}
	if args.DateFrom.After(args.DateTo) {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Date does not valid")
	}

	var promotion pricelistpromotion.ShipmentPriceListPromotion
	if err := scheme.Convert(args, &promotion); err != nil {
		return nil, err
	}
	return a.priceListPromotionStore(ctx).CreateShipmentPriceListPromotion(&promotion)
}

func (a *Aggregate) UpdatePriceListPromotion(ctx context.Context, args *pricelistpromotion.UpdatePriceListPromotionArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing price list promotion ID")
	}
	if args.DateTo.Sub(args.DateFrom) < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Date does not valid")
	}
	var promotion pricelistpromotion.ShipmentPriceListPromotion
	if err := scheme.Convert(args, &promotion); err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err := a.priceListPromotionStore(ctx).UpdateShipmentPriceListPromotion(&promotion); err != nil {
			return err
		}
		if !args.Status.Valid || args.Status.Enum != status3.P {
			return nil
		}

		_promotion, err := a.priceListPromotionStore(ctx).ID(args.ID).GetShipmentPriceListPromotion()
		if err != nil {
			return err
		}
		return checkValidPromotion(_promotion)
	})
}

func checkValidPromotion(promotion *pricelistpromotion.ShipmentPriceListPromotion) error {
	now := time.Now()
	if promotion.DateTo.Sub(promotion.DateFrom) < 0 ||
		now.After(promotion.DateTo) {
		return cm.Errorf(cm.FailedPrecondition, nil, "Date does not valid.")
	}

	rules := promotion.AppliedRules
	if rules != nil {
		if !rules.UserCreatedDate.IsValid() {
			return cm.Errorf(cm.FailedPrecondition, nil, "User created date does not valid")
		}
		if !rules.ShopCreatedDate.IsValid() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Shop created date does not valid")
		}
	}
	return nil
}

func (a *Aggregate) DeletePriceListPromotion(ctx context.Context, id dot.ID) error {
	_, err := a.priceListPromotionStore(ctx).ID(id).SoftDelete()
	return err
}
