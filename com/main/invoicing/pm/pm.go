package pm

import (
	"context"

	"o.o/api/external/payment"
	"o.o/api/main/credit"
	"o.o/api/main/identity"
	"o.o/api/main/invoicing"
	"o.o/api/main/transaction"
	"o.o/api/subscripting/types"
	"o.o/api/top/types/etc/invoice_type"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/subject_referral"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

type ProcessManager struct {
	transactionQuery transaction.QueryBus
	creditQ          credit.QueryBus
	invoiceAggr      invoicing.CommandBus
	invoiceQ         invoicing.QueryBus
	identityQ        identity.QueryBus
}

func New(
	eventBus bus.EventRegistry,
	transactionQ transaction.QueryBus,
	creditQ credit.QueryBus,
	invoiceAggr invoicing.CommandBus,
	invoiceQ invoicing.QueryBus,
	identityQ identity.QueryBus,
) *ProcessManager {
	p := &ProcessManager{
		transactionQuery: transactionQ,
		creditQ:          creditQ,
		invoiceAggr:      invoiceAggr,
		invoiceQ:         invoiceQ,
		identityQ:        identityQ,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvoicePaying)
	eventBus.AddEventListener(m.HandleCreditCreated)
	eventBus.AddEventListener(m.HandleCreditConfirmed)
	eventBus.AddEventListener(m.HandlePaymentStatusUpdated)
}

func (m *ProcessManager) InvoicePaying(ctx context.Context, event *invoicing.InvoicePayingEvent) error {
	if event.PaymentMethod != payment_method.Balance {
		return nil
	}
	if !event.ServiceClassify.Valid {
		return wrapError(cm.Internal, nil, "Missing credit classify")
	}
	if event.OwnerID == 0 {
		return wrapError(cm.Internal, nil, "Missing ownerID")
	}
	var balance int

	query := &transaction.GetBalanceUserQuery{
		UserID:   event.OwnerID,
		Classify: event.ServiceClassify.Enum,
	}
	if err := m.transactionQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	balance = query.Result.ActualBalance

	if balance < event.TotalAmount {
		return cm.Errorf(cm.FailedPrecondition, nil, "Số dư không đủ để thanh toán")
	}
	return nil
}

func wrapError(code xerrors.Code, err error, msg string) error {
	if err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Lỗi từ invoice process manager: %v", err.Error())
	}
	return cm.Errorf(code, err, "Lỗi từ invoice process manager: %v", msg)
}

func (m *ProcessManager) HandleCreditCreated(ctx context.Context, event *credit.CreditCreatedEvent) error {
	_, err := m.createInvoiceByCredit(ctx, event.CreditID, event.ShopID)
	return err
}

func (m *ProcessManager) createInvoiceByCredit(ctx context.Context, creditID, shopID dot.ID) (*invoicing.InvoiceFtLine, error) {
	getCreditQuery := &credit.GetCreditQuery{
		ID:     creditID,
		ShopID: shopID,
	}
	if err := m.creditQ.Dispatch(ctx, getCreditQuery); err != nil {
		return nil, err
	}
	_credit := getCreditQuery.Result

	// get customer info from shopOwner
	getShopQuery := &identity.GetShopByIDQuery{
		ID: shopID,
	}
	if err := m.identityQ.Dispatch(ctx, getShopQuery); err != nil {
		return nil, err
	}
	_shop := getShopQuery.Result
	getUserQuery := &identity.GetUserByIDQuery{
		UserID: _shop.OwnerID,
	}
	if err := m.identityQ.Dispatch(ctx, getUserQuery); err != nil {
		return nil, err
	}
	_user := getUserQuery.Result

	// create invoice
	amount := _credit.Amount
	invoiceType := invoice_type.In
	if amount < 0 {
		invoiceType = invoice_type.Out
		amount *= -1
	}
	cmd := &invoicing.CreateInvoiceCommand{
		AccountID:   _credit.ShopID,
		TotalAmount: amount,
		Customer: &types.CustomerInfo{
			FullName: _user.FullName,
			Email:    _user.Email,
			Phone:    _user.Phone,
		},
		ReferralType: subject_referral.Credit,
		Lines: []*invoicing.InvoiceLine{
			{
				LineAmount:   amount,
				Price:        amount,
				Quantity:     1,
				ReferralType: subject_referral.Credit,
				ReferralID:   _credit.ID,
			},
		},
		Classify: service_classify.ServiceClassify(_credit.Classify),
		Type:     invoiceType,
	}

	if err := m.invoiceAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return cmd.Result, nil
}

func (m *ProcessManager) HandleCreditConfirmed(ctx context.Context, event *credit.CreditConfirmedEvent) error {
	queryInv := &invoicing.GetInvoiceByReferralQuery{
		ReferralType: subject_referral.Credit,
		ReferralIDs:  []dot.ID{event.CreditID},
	}
	if err := m.invoiceQ.Dispatch(ctx, queryInv); err != nil {
		return err
	}
	inv := queryInv.Result

	cmdPaymentInvoice := &invoicing.PaymentInvoiceCommand{
		InvoiceID:       inv.ID,
		AccountID:       event.ShopID,
		TotalAmount:     inv.TotalAmount,
		PaymentMethod:   payment_method.Manual,
		ServiceClassify: inv.Classify.Wrap(),
	}
	return m.invoiceAggr.Dispatch(ctx, cmdPaymentInvoice)
}

func (m *ProcessManager) HandlePaymentStatusUpdated(ctx context.Context, event *payment.PaymentStatusUpdatedEvent) error {
	getInvoiceQuery := &invoicing.GetInvoiceByPaymentIDQuery{
		PaymentID: event.ID,
	}
	err := m.invoiceQ.Dispatch(ctx, getInvoiceQuery)
	switch cm.ErrorCode(err) {
	case cm.NoError:
		// continue
	case cm.NotFound:
		return nil
	default:
		return err
	}
	inv := getInvoiceQuery.Result

	update := &invoicing.UpdateInvoicePaymentInfoCommand{
		ID:            inv.ID,
		AccountID:     inv.AccountID,
		PaymentStatus: event.PaymentStatus,
	}
	return m.invoiceAggr.Dispatch(ctx, update)
}
