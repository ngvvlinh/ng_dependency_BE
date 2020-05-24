package subscriptionbill

import (
	"context"
	"time"

	"o.o/api/external/payment"
	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionbill"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/backend/com/subscripting/subscriptionbill/convert"
	"o.o/backend/com/subscripting/subscriptionbill/model"
	"o.o/backend/com/subscripting/subscriptionbill/sqlstore"
	subrplanutils "o.o/backend/com/subscripting/subscriptionplan/utils"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ subscriptionbill.Aggregate = &SubrBillAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type SubrBillAggregate struct {
	db                cmsql.Transactioner
	eventBus          capi.EventBus
	subrBillStore     sqlstore.SubrBillStoreFactory
	subrBillLineStore sqlstore.SubrBillLineStoreFactory
	paymentAggr       payment.CommandBus
	subrQuery         subscription.QueryBus
	subrPlanQuery     subscriptionplan.QueryBus
}

func NewSubrBillAggregate(db *cmsql.Database, eventB capi.EventBus, paymentA payment.CommandBus, subrQuery subscription.QueryBus, subrPlanQuery subscriptionplan.QueryBus) *SubrBillAggregate {
	return &SubrBillAggregate{
		db:                db,
		eventBus:          eventB,
		subrBillStore:     sqlstore.NewSubrBillStore(db),
		subrBillLineStore: sqlstore.NewSubrBillLineStore(db),
		subrQuery:         subrQuery,
		subrPlanQuery:     subrPlanQuery,
		paymentAggr:       paymentA,
	}
}

func SubrBillAggregateMessageBus(a *SubrBillAggregate) subscriptionbill.CommandBus {
	b := bus.New()
	return subscriptionbill.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *SubrBillAggregate) CreateSubscriptionBill(ctx context.Context, args *subscriptionbill.CreateSubscriptionBillArgs) (*subscriptionbill.SubscriptionBillFtLine, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account ID").WithMetap("func", "CreateSubscriptionBill")
	}
	if args.SubscriptionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID").WithMetap("func", "CreateSubscriptionBill")
	}
	if args.Customer == nil || args.Customer.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing customer").WithMetap("func", "CreateSubscriptionBill")
	}
	if len(args.Lines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing lines").WithMetap("func", "CreateSubscriptionBill")
	}
	amount := 0
	for _, line := range args.Lines {
		lineAmount := line.Quantity * line.Price
		amount += lineAmount
		if line.LineAmount != lineAmount {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Line amount does not match")
		}
	}
	if amount != args.TotalAmount {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Total amount does not match").WithMetap("func", "CreateSubscriptionBill")
	}
	for _, line := range args.Lines {
		if err := verifySubrBillLine(line); err != nil {
			return nil, err
		}
	}

	var subrBill subscriptionbill.SubscriptionBill
	if err := scheme.Convert(args, &subrBill); err != nil {
		return nil, err
	}

	subrBillID := cm.NewID()
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		subrBill.ID = subrBillID
		if err := a.subrBillStore(ctx).CreateSubrBill(&subrBill); err != nil {
			return err
		}
		for _, line := range args.Lines {
			if line.ID == 0 {
				line.ID = cm.NewID()
			}
			line.SubscriptionBillID = subrBill.ID
			if err := a.subrBillLineStore(ctx).CreateSubrBillLine(line); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.subrBillStore(ctx).ID(subrBillID).GetSubrBillFtLine()
}

func verifySubrBillLine(line *subscriptionbill.SubscriptionBillLine) error {
	if line.SubscriptionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID in subscription bill line")
	}
	if line.Quantity <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Quantity does not valid in subscription bill line")
	}
	if line.PeriodStartAt.IsZero() {
		return cm.Errorf(cm.InvalidArgument, nil, "Period start at is not valid in subscription bill line")
	}
	if line.PeriodEndAt.IsZero() {
		return cm.Errorf(cm.InvalidArgument, nil, "Period end at is not valid in subscription bill line")
	}
	if line.PeriodEndAt.Sub(line.PeriodStartAt) < 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Period is not valid in subscription bill line")
	}
	return nil
}

func (a *SubrBillAggregate) CreateSubscriptionBillBySubrID(ctx context.Context, args *subscriptionbill.CreateSubscriptionBillBySubrIDArgs) (*subscriptionbill.SubscriptionBillFtLine, error) {
	if args.SubscriptionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID").WithMetap("func", "CreateSubscriptionBillBySubrID")
	}
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising account ID").WithMetap("func", "CreateSubscriptionBillBySubrID")
	}
	query := &subscription.GetSubscriptionByIDQuery{
		ID:        args.SubscriptionID,
		AccountID: args.AccountID,
	}
	if err := a.subrQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	subr := query.Result
	bill := &subscriptionbill.CreateSubscriptionBillArgs{
		AccountID:      subr.AccountID,
		SubscriptionID: subr.ID,
		TotalAmount:    0,   // filled below
		Lines:          nil, // filled below
		Description:    args.Description,
		Customer:       subr.Customer,
	}
	if args.Customer != nil {
		bill.Customer = args.Customer
	}

	var billLines = make([]*subscriptionbill.SubscriptionBillLine, len(subr.Lines))
	totalAmount := 0
	for i, line := range subr.Lines {
		if line.PlanID == 0 {
			continue
		}
		queryPlan := &subscriptionplan.GetSubrPlanByIDQuery{
			ID: line.PlanID,
		}
		if err := a.subrPlanQuery.Dispatch(ctx, queryPlan); err != nil {
			return nil, err
		}
		plan := queryPlan.Result
		lineAmount := line.Quantity * plan.Price
		totalAmount += lineAmount
		now := time.Now()
		bLine := &subscriptionbill.SubscriptionBillLine{
			ID:             cm.NewID(),
			LineAmount:     lineAmount,
			Price:          plan.Price,
			Quantity:       line.Quantity,
			Description:    "",
			PeriodStartAt:  time.Time{},
			PeriodEndAt:    time.Time{},
			SubscriptionID: subr.ID,
		}
		if subr.CurrentPeriodEndAt.IsZero() {
			bLine.PeriodStartAt = now
		} else {
			bLine.PeriodStartAt = subr.CurrentPeriodEndAt
		}
		duration, err := subrplanutils.CalcPlanDuration(plan.Interval, plan.IntervalCount)
		if err != nil {
			return nil, err
		}
		bLine.PeriodEndAt = bLine.PeriodStartAt.Add(duration)
		billLines[i] = bLine
	}
	if totalAmount != args.TotalAmount {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Total amount does not match").WithMetap("func", "CreateSubscriptionBillBySubrID").WithMetap("expected total_amount", totalAmount)
	}
	bill.Lines = billLines
	bill.TotalAmount = totalAmount
	return a.CreateSubscriptionBill(ctx, bill)
}

func (a *SubrBillAggregate) UpdateSubscriptionBillPaymentInfo(ctx context.Context, args *subscriptionbill.UpdateSubscriptionBillPaymentInfoArgs) error {
	update := &model.SubscriptionBill{
		PaymentID:     args.PaymentID,
		PaymentStatus: args.PaymentStatus,
	}
	if args.PaymentStatus == status4.P {
		update.Status = status4.P
	}
	return a.subrBillStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateSubrBillDB(update)
}

func (a *SubrBillAggregate) UpdateSubscriptionBillStatus(ctx context.Context, args *subscriptionbill.UpdateSubscriptionBillStatusArgs) error {
	if !args.Status.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing subscription bill status")
	}
	update := &model.SubscriptionBill{
		Status: args.Status.Enum,
	}
	return a.subrBillStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateSubrBillDB(update)
}

func (a *SubrBillAggregate) DeleteSubsciptionBill(ctx context.Context, id dot.ID, accountID dot.ID) error {
	bill, err := a.subrBillStore(ctx).ID(id).AccountID(accountID).GetSubrBillFtLine()
	if err != nil {
		return err
	}
	if bill.Status == status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "This bill was completed. Can not delete")
	}
	_, err = a.subrBillStore(ctx).ID(id).OptionalAccountID(accountID).SoftDelete()
	return err
}

/*
	ManualPaymentSubscriptionBill: thanh toán thủ công bill (admin gọi)
	- Tạo payment
	- Gắn payment_id vào bill
	- Cập nhật trạng thái thanh toán bill
*/
func (a *SubrBillAggregate) ManualPaymentSubscriptionBill(ctx context.Context, args *subscriptionbill.ManualPaymentSubrBillArgs) error {
	bill, err := a.subrBillStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).GetSubrBillFtLine()
	if err != nil {
		return err
	}
	if bill.Status == status4.N {
		return cm.Errorf(cm.FailedPrecondition, nil, "Subscription bill was cancelled")
	}
	if bill.Status == status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Subscription bill was completed")
	}
	if bill.TotalAmount != args.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total amount does not match")
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// Create manual payment and assign payment ID to bill
		cmd := &payment.CreatePaymentCommand{
			Amount:          args.TotalAmount,
			Status:          status4.P,
			State:           payment_state.Success,
			PaymentProvider: payment_provider.Manual,
		}
		if err := a.paymentAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
		paymentID := cmd.Result.ID

		update := &subscriptionbill.UpdateSubscriptionBillPaymentInfoArgs{
			ID:            args.ID,
			AccountID:     args.AccountID,
			PaymentID:     paymentID,
			PaymentStatus: status4.P,
		}
		err = a.UpdateSubscriptionBillPaymentInfo(ctx, update)
		if err != nil {
			return err
		}
		event := &subscriptionbill.SubscriptionBillPaidEvent{
			ID:        args.ID,
			AccountID: args.AccountID,
		}
		return a.eventBus.Publish(ctx, event)
	})
}
