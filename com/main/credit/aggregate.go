package credit

import (
	"context"
	"time"

	"o.o/api/main/credit"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/credit_type"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/credit/convert"
	"o.o/backend/com/main/credit/model"
	"o.o/backend/com/main/credit/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/common/l"
)

var ll = l.New()
var _ credit.Aggregate = &CreditAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)
var zeroTime = time.Unix(0, 0)

type CreditAggregate struct {
	dbTx          cmsql.Transactioner
	CreditStore   sqlstore.CreditFactory
	eventBus      capi.EventBus
	identityQuery identity.QueryBus
}

func NewAggregateCredit(
	bus capi.EventBus,
	db com.MainDB,
	identityQ identity.QueryBus,
) *CreditAggregate {
	return &CreditAggregate{
		dbTx:          (*cmsql.Database)(db),
		identityQuery: identityQ,
		eventBus:      bus,
		CreditStore:   sqlstore.NewCreditStore(db),
	}
}

func CreditAggregateMessageBus(q *CreditAggregate) credit.CommandBus {
	b := bus.New()
	return credit.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a CreditAggregate) CreateCredit(ctx context.Context, args *credit.CreateCreditArgs) (*credit.CreditExtended, error) {
	switch args.Type {
	case credit_type.Shop:
		if args.ShopID == 0 {
			return nil, cm.Error(cm.InvalidArgument, "Missing shop_id", nil)
		}
	default:
		return nil, cm.Error(cm.InvalidArgument, "Type does not support", nil)
	}
	if args.Amount == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing amount", nil)
	}

	creditClassify := args.Classify
	shopCreditAmount, err := a.CreditStore(ctx).ShopID(args.ShopID).Classify(creditClassify).SumCredit()
	if err != nil {
		return nil, err
	}
	if shopCreditAmount+args.Amount < 0 {
		return nil, cm.Error(cm.InvalidArgument, "Shop balance is not enough", nil)
	}
	var result = &credit.Credit{}
	err = scheme.Convert(args, result)
	if err != nil {
		return nil, err
	}

	getAccountQuery := &identity.GetShopByIDQuery{
		ID: args.ShopID,
	}
	if err = a.identityQuery.Dispatch(ctx, getAccountQuery); err != nil {
		return nil, err
	}

	err = a.dbTx.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err = a.CreditStore(ctx).Create(result)
		if err != nil {
			return err
		}

		event := &credit.CreditCreatedEvent{
			CreditID: result.ID,
			ShopID:   result.ShopID,
		}
		return a.eventBus.Publish(ctx, event)
	})
	if err != nil {
		return nil, err
	}

	return &credit.CreditExtended{
		Credit: result,
		Shop:   getAccountQuery.Result,
	}, nil
}

func (a CreditAggregate) ConfirmCredit(ctx context.Context, args *credit.ConfirmCreditArgs) (*credit.CreditExtended, error) {
	if args.ID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	query := a.CreditStore(ctx).ID(args.ID)
	if args.ShopID != 0 {
		query = query.ShopID(args.ShopID)
	}
	creditValue, err := query.Get()
	if err != nil {
		return nil, err
	}

	creditClassify := creditValue.Classify
	shopCreditAmount, err := a.CreditStore(ctx).ShopID(creditValue.ShopID).Classify(creditClassify).SumCredit()
	if err != nil {
		return nil, err
	}
	if creditValue.Status == status3.P {
		return nil, cm.Error(cm.FailedPrecondition, "This credit has already confirmed", nil)
	}
	if creditValue.Status != status3.Z {
		return nil, cm.Error(cm.FailedPrecondition, "Can not confirm this credit", nil)
	}
	if creditValue.PaidAt.IsZero() || creditValue.PaidAt.Equal(zeroTime) {
		return nil, cm.Error(cm.FailedPrecondition, "Missing paid at", nil)
	}
	if shopCreditAmount+creditValue.Amount < 0 {
		return nil, cm.Error(cm.InvalidArgument, "Shop balance is not enough", nil)
	}

	getShopQuery := &identity.GetShopByIDQuery{
		ID: creditValue.ShopID,
	}
	if err = a.identityQuery.Dispatch(ctx, getShopQuery); err != nil {
		return nil, err
	}
	err = a.dbTx.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		var updateModel = &model.Credit{
			Status: status3.P,
		}
		if err = query.UpdateCreditDB(updateModel); err != nil {
			return err
		}
		creditValue.Status = status3.P

		event := &credit.CreditConfirmedEvent{
			CreditID: args.ID,
			ShopID:   creditValue.ShopID,
		}
		return a.eventBus.Publish(ctx, event)
	})
	if err != nil {
		return nil, err
	}
	return &credit.CreditExtended{
		Credit: creditValue,
		Shop:   getShopQuery.Result,
	}, nil
}

func (a CreditAggregate) DeleteCredit(ctx context.Context, args *credit.DeleteCreditArgs) (int, error) {
	if args.ID == 0 {
		return 0, cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	query := a.CreditStore(ctx).ID(args.ID)
	if args.ShopID != 0 {
		query = query.ShopID(args.ShopID)
	}
	creditValue, err := query.Get()
	if err != nil {
		return 0, err
	}

	if creditValue.Status == status3.P {
		return 0, cm.Error(cm.FailedPrecondition, "This credit has already confirmed", nil)
	}
	if creditValue.Status != status3.Z {
		return 0, cm.Error(cm.FailedPrecondition, "Can not delete this credit", nil)
	}
	return a.CreditStore(ctx).ID(creditValue.ID).ShopID(creditValue.ShopID).SoftDelete()
}
