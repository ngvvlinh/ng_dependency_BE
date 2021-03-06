package invoicing

import (
	"context"

	"o.o/api/external/payment"
	"o.o/api/main/invoicing"
	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionplan"
	"o.o/api/top/types/etc/invoice_type"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/invoicing/convert"
	"o.o/backend/com/main/invoicing/model"
	"o.o/backend/com/main/invoicing/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ invoicing.Aggregate = &InvoiceAggregate{}
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

func InvoiceAggregateMessageBus(a *InvoiceAggregate) invoicing.CommandBus {
	b := bus.New()
	return invoicing.NewAggregateHandler(a).RegisterHandlers(b)
}

func (a *InvoiceAggregate) CreateInvoice(ctx context.Context, args *invoicing.CreateInvoiceArgs) (*invoicing.InvoiceFtLine, error) {
	if err := args.Validate(); err != nil {
		return nil, cm.Errorf(cm.ErrorCode(err), err, err.Error()).WithMetap("func", "CreateInvoice")
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

	var inv invoicing.Invoice
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

func verifyInvoiceLine(line *invoicing.InvoiceLine, refType subject_referral.SubjectReferral) error {
	if line.Quantity <= 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Quantity does not valid in invoice line").WithMeta("func", "verifyInvoiceLine")
	}
	if refType != line.ReferralType {
		return cm.Errorf(cm.InvalidArgument, nil, "Referral type does not match").WithMeta("func", "verifyInvoiceLine")
	}
	return nil
}

func (a *InvoiceAggregate) CreateInvoiceBySubrID(ctx context.Context, args *invoicing.CreateInvoiceBySubrIDArgs) (*invoicing.InvoiceFtLine, error) {
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
	invArgs := &invoicing.CreateInvoiceArgs{
		AccountID:    inv.AccountID,
		TotalAmount:  0,   // filled below
		Lines:        nil, // filled below
		Description:  args.Description,
		Customer:     inv.Customer,
		ReferralType: subject_referral.Subscription,
		Classify:     args.Classify,
		Type:         invoice_type.Out,
	}
	if args.Customer != nil {
		invArgs.Customer = args.Customer
	}

	var invLines = make([]*invoicing.InvoiceLine, len(inv.Lines))
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
		bLine := &invoicing.InvoiceLine{
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

func (a *InvoiceAggregate) UpdateInvoicePaymentInfo(ctx context.Context, args *invoicing.UpdateInvoicePaymentInfoArgs) error {
	update := &model.Invoice{
		PaymentID:     args.PaymentID,
		PaymentStatus: args.PaymentStatus,
	}
	if args.PaymentStatus == status4.P {
		update.Status = status4.P
	}
	if err := a.invoiceStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateInvoiceDB(update); err != nil {
		return err
	}
	inv, err := a.invoiceStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).GetInvoice()
	if err != nil {
		return err
	}
	if update.Status == status4.P {
		paidEvent := &invoicing.InvoicePaidEvent{
			ID:              args.ID,
			PaymentID:       args.PaymentID,
			AccountID:       args.AccountID,
			ServiceClassify: inv.Classify.Wrap(),
		}
		return a.eventBus.Publish(ctx, paidEvent)
	}
	return nil
}

func (a *InvoiceAggregate) DeleteInvoice(ctx context.Context, args *invoicing.DeleteInvoiceArgs) error {
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

	event := &invoicing.InvoiceDeletedEvent{InvoinceID: args.ID}
	return a.eventBus.Publish(ctx, event)
}

// Handle case payment invoice by balance or manual (admin confirm)
// - T???o payment
// - G???n payment_id v??o invoice
// - C???p nh???t tr???ng th??i thanh to??n invoice
func (a *InvoiceAggregate) PaymentInvoice(ctx context.Context, args *invoicing.PaymentInvoiceArgs) error {
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
	event := &invoicing.InvoicePayingEvent{
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
				ShopID:          args.AccountID,
				Amount:          args.TotalAmount,
				Status:          status4.P,
				State:           payment_state.Success,
				PaymentProvider: payment_provider.Manual,
			}
			if err = a.paymentAggr.Dispatch(ctx, cmd); err != nil {
				return err
			}
			paymentID := cmd.Result.ID

			update := &invoicing.UpdateInvoicePaymentInfoArgs{
				ID:            args.InvoiceID,
				AccountID:     args.AccountID,
				PaymentID:     paymentID,
				PaymentStatus: status4.P,
			}
			return a.UpdateInvoicePaymentInfo(ctx, update)
		default:
			return cm.Errorf(cm.InvalidArgument, nil, "Ph????ng th???c thanh to??n kh??ng h???p l???")
		}
	})
}
