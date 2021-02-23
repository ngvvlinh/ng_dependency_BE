package invoice

import (
	"context"

	"o.o/api/external/payment"
	"o.o/api/subscripting/invoice"
	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	com "o.o/backend/com/main"
	"o.o/backend/com/subscripting/invoice/convert"
	"o.o/backend/com/subscripting/invoice/model"
	"o.o/backend/com/subscripting/invoice/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ invoice.Aggregate = &InvoiceAggregate{}
var scheme = conversion.Build(convert.RegisterConversions)

type InvoiceAggregate struct {
	db               *cmsql.Database
	eventBus         capi.EventBus
	invoiceStore     sqlstore.InvoiceStoreFactory
	invoiceLineStore sqlstore.InvoiceLineStoreFactory
	paymentAggr      payment.CommandBus
	subrQuery        subscription.QueryBus
	subrPlanQuery    subscriptionplan.QueryBus
}

func NewInvoiceAggregate(db com.MainDB, eventB capi.EventBus, paymentA payment.CommandBus, subrQuery subscription.QueryBus, subrPlanQuery subscriptionplan.QueryBus) *InvoiceAggregate {
	return &InvoiceAggregate{
		db:               db,
		eventBus:         eventB,
		invoiceStore:     sqlstore.NewInvoiceStore(db),
		invoiceLineStore: sqlstore.NewInvoiceLineStore(db),
		subrQuery:        subrQuery,
		subrPlanQuery:    subrPlanQuery,
		paymentAggr:      paymentA,
	}
}

func InvoiceAggregateMessageBus(a *InvoiceAggregate) invoice.CommandBus {
	b := bus.New()
	return invoice.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *InvoiceAggregate) CreateInvoice(ctx context.Context, args *invoice.CreateInvoiceArgs) (*invoice.InvoiceFtLine, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing account ID").WithMetap("func", "CreateInvoice")
	}
	if args.Customer == nil || args.Customer.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing customer").WithMetap("func", "CreateInvoice")
	}
	if len(args.Lines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing lines").WithMetap("func", "CreateInvoice")
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Total amount does not match").WithMetap("func", "CreateInvoice")
	}
	var referralIDs []dot.ID
	for _, line := range args.Lines {
		if err := verifyInvoiceLine(line, args.ReferralType); err != nil {
			return nil, err
		}
		if line.ReferralID != 0 {
			referralIDs = append(referralIDs, line.ReferralID)
		}
	}

	var inv invoice.Invoice
	if err := scheme.Convert(args, &inv); err != nil {
		return nil, err
	}

	invoiceID := cm.NewID()
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		inv.ID = invoiceID
		inv.ReferralIDs = referralIDs
		if err := a.invoiceStore(ctx).CreateInvoice(&inv); err != nil {
			return err
		}
		for _, line := range args.Lines {
			if line.ID == 0 {
				line.ID = cm.NewID()
			}
			line.InvoiceID = inv.ID
			if err := a.invoiceLineStore(ctx).CreateInvoiceLine(line); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.invoiceStore(ctx).ID(invoiceID).GetInvoiceFtLine()
}

func verifyInvoiceLine(line *invoice.InvoiceLine, refType subject_referral.SubjectReferral) error {
	if line.Quantity <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Quantity does not valid in invoice line").WithMeta("func", "verifyInvoiceLine")
	}
	if refType != line.ReferralType {
		return cm.Errorf(cm.InvalidArgument, nil, "Referral type does not match").WithMeta("func", "verifyInvoiceLine")
	}
	return nil
}

func (a *InvoiceAggregate) CreateInvoiceBySubrID(ctx context.Context, args *invoice.CreateInvoiceBySubrIDArgs) (*invoice.InvoiceFtLine, error) {
	if args.SubscriptionID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subscription ID").WithMetap("func", "CreateInvoiceBySubrID")
	}
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Mising account ID").WithMetap("func", "CreateInvoiceBySubrID")
	}
	query := &subscription.GetSubscriptionByIDQuery{
		ID:        args.SubscriptionID,
		AccountID: args.AccountID,
	}
	if err := a.subrQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	inv := query.Result
	invArgs := &invoice.CreateInvoiceArgs{
		AccountID:    inv.AccountID,
		TotalAmount:  0,   // filled below
		Lines:        nil, // filled below
		Description:  args.Description,
		Customer:     inv.Customer,
		ReferralType: subject_referral.Subscription,
	}
	if args.Customer != nil {
		invArgs.Customer = args.Customer
	}

	var invLines = make([]*invoice.InvoiceLine, len(inv.Lines))
	totalAmount := 0
	for i, line := range inv.Lines {
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
		bLine := &invoice.InvoiceLine{
			ID:           cm.NewID(),
			LineAmount:   lineAmount,
			Price:        plan.Price,
			Quantity:     line.Quantity,
			Description:  "",
			ReferralType: subject_referral.Subscription,
			ReferralID:   args.SubscriptionID,
		}
		invLines[i] = bLine
	}
	if args.TotalAmount != 0 && totalAmount != args.TotalAmount {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Total amount does not match").WithMetap("func", "CreateInvoiceBySubrID").WithMetap("expected total_amount", totalAmount)
	}
	invArgs.Lines = invLines
	invArgs.TotalAmount = totalAmount
	return a.CreateInvoice(ctx, invArgs)
}

func (a *InvoiceAggregate) UpdateInvoicePaymentInfo(ctx context.Context, args *invoice.UpdateInvoicePaymentInfoArgs) error {
	update := &model.Invoice{
		PaymentID:     args.PaymentID,
		PaymentStatus: args.PaymentStatus,
	}
	if args.PaymentStatus == status4.P {
		update.Status = status4.P
	}
	return a.invoiceStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateInvoiceDB(update)
}

func (a *InvoiceAggregate) UpdateInvoiceStatus(ctx context.Context, args *invoice.UpdateInvoiceStatusArgs) error {
	if !args.Status.Valid {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing invoice status")
	}
	update := &model.Invoice{
		Status: args.Status.Enum,
	}
	return a.invoiceStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateInvoiceDB(update)
}

func (a *InvoiceAggregate) DeleteInvoice(ctx context.Context, args *invoice.DeleteInvoiceArgs) error {
	inv, err := a.invoiceStore(ctx).ID(args.ID).AccountID(args.AccountID).GetInvoiceFtLine()
	if err != nil {
		return err
	}
	if !args.ForceDelete && inv.Status == status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "This invoice was completed. Can not delete")
	}
	_, err = a.invoiceStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).SoftDelete()
	if err != nil {
		return err
	}

	event := &invoice.InvoiceDeletedEvent{InvoinceID: args.ID}
	return a.eventBus.Publish(ctx, event)
}

/*
	ManualPaymentInvoice: thanh toán thủ công invoice (admin gọi)
	- Tạo payment
	- Gắn payment_id vào invoice
	- Cập nhật trạng thái thanh toán invoice
*/
func (a *InvoiceAggregate) ManualPaymentInvoice(ctx context.Context, args *invoice.ManualPaymentInvoiceArgs) error {
	inv, err := a.invoiceStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).GetInvoiceFtLine()
	if err != nil {
		return err
	}
	if inv.Status == status4.N {
		return cm.Errorf(cm.FailedPrecondition, nil, "Invoice was cancelled")
	}
	if inv.Status == status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Invoice was completed")
	}
	if inv.TotalAmount != args.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total amount does not match")
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		// Create manual payment and assign payment ID to invoice
		cmd := &payment.CreatePaymentCommand{
			Amount:          args.TotalAmount,
			Status:          status4.P,
			State:           payment_state.Success,
			PaymentProvider: payment_provider.Manual,
		}
		if err = a.paymentAggr.Dispatch(ctx, cmd); err != nil {
			return err
		}
		paymentID := cmd.Result.ID

		update := &invoice.UpdateInvoicePaymentInfoArgs{
			ID:            args.ID,
			AccountID:     args.AccountID,
			PaymentID:     paymentID,
			PaymentStatus: status4.P,
		}
		err = a.UpdateInvoicePaymentInfo(ctx, update)
		if err != nil {
			return err
		}
		event := &invoice.InvoicePaidEvent{
			ID:        args.ID,
			AccountID: args.AccountID,
		}
		return a.eventBus.Publish(ctx, event)
	})
}

func (a *InvoiceAggregate) PaymentInvoice(ctx context.Context, args *invoice.PaymentInvoiceArgs) error {
	inv, err := a.invoiceStore(ctx).ID(args.InvoiceID).OptionalAccountID(args.AccountID).GetInvoiceFtLine()
	if err != nil {
		return err
	}
	if inv.Status == status4.N {
		return cm.Errorf(cm.FailedPrecondition, nil, "Invoice was cancelled")
	}
	if inv.Status == status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Invoice was completed")
	}
	if inv.TotalAmount != args.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Total amount does not match")
	}

	// check credit balance
	event := &invoice.InvoicePayingEvent{
		PaymentMethod:   args.PaymentMethod,
		ServiceClassify: args.ServiceClassify,
		OwnerID:         args.OwnerID,
		TotalAmount:     args.TotalAmount,
	}
	if err = a.eventBus.Publish(ctx, event); err != nil {
		return err
	}

	return a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		switch args.PaymentMethod {
		case payment_method.Manual,
			payment_method.Balance:
			cmd := &payment.CreatePaymentCommand{
				Amount:          args.TotalAmount,
				Status:          status4.P,
				State:           payment_state.Success,
				PaymentProvider: payment_provider.Manual,
			}
			if err = a.paymentAggr.Dispatch(ctx, cmd); err != nil {
				return err
			}
			paymentID := cmd.Result.ID

			update := &invoice.UpdateInvoicePaymentInfoArgs{
				ID:            args.InvoiceID,
				AccountID:     args.AccountID,
				PaymentID:     paymentID,
				PaymentStatus: status4.P,
			}
			err = a.UpdateInvoicePaymentInfo(ctx, update)
			if err != nil {
				return err
			}
			event2 := &invoice.InvoicePaidEvent{
				ID:              args.InvoiceID,
				AccountID:       args.AccountID,
				PaymentMethod:   args.PaymentMethod,
				ServiceClassify: args.ServiceClassify,
			}
			return a.eventBus.Publish(ctx, event2)
		default:
			return cm.Errorf(cm.InvalidArgument, nil, "Phương thức thanh toán không hợp lệ")
		}

	})
}
