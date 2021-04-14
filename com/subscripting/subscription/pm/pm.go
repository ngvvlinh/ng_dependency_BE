package pm

import (
	"context"
	"time"

	"o.o/api/main/credit"
	"o.o/api/main/identity"
	"o.o/api/main/invoicing"
	"o.o/api/main/transaction"
	"o.o/api/subscripting/subscription"
	"o.o/api/subscripting/subscriptionplan"
	subscriptingtypes "o.o/api/subscripting/types"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/api/webserver"
	subrplanutils "o.o/backend/com/subscripting/subscriptionplan/utils"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
)

type ProcessManager struct {
	invoiceQuery    invoicing.QueryBus
	invoiceAggr     invoicing.CommandBus
	subrQuery       subscription.QueryBus
	subrAggr        subscription.CommandBus
	subrPlanQuery   subscriptionplan.QueryBus
	identityQuery   identity.QueryBus
	creditAggr      credit.CommandBus
	transactionAggr transaction.CommandBus
}

func New(
	eventBus bus.EventRegistry,
	invoiceQuery invoicing.QueryBus,
	invoiceAggr invoicing.CommandBus,
	subrQuery subscription.QueryBus,
	subrAggr subscription.CommandBus,
	subrPlanQuery subscriptionplan.QueryBus,
	identityQuery identity.QueryBus,
	creditA credit.CommandBus,
	transactionA transaction.CommandBus,
) *ProcessManager {
	p := &ProcessManager{
		invoiceQuery:    invoiceQuery,
		invoiceAggr:     invoiceAggr,
		subrQuery:       subrQuery,
		subrAggr:        subrAggr,
		subrPlanQuery:   subrPlanQuery,
		identityQuery:   identityQuery,
		creditAggr:      creditA,
		transactionAggr: transactionA,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvoicePaid)
	eventBus.AddEventListener(m.WsWebsiteCreated)
}

func (m *ProcessManager) InvoicePaid(ctx context.Context, event *invoicing.InvoicePaidEvent) (_err error) {
	queryInv := &invoicing.GetInvoiceByIDQuery{
		ID:        event.ID,
		AccountID: event.AccountID,
	}
	if err := m.invoiceQuery.Dispatch(ctx, queryInv); err != nil {
		return err
	}
	inv := queryInv.Result
	if inv.Status != status4.P {
		return cm.Errorf(cm.FailedPrecondition, nil, "Invoice was not paid").WithMetap("event", "InvoicePaidEvent")
	}

	switch inv.ReferralType {
	case subject_referral.Subscription:
		return m.extendSubscription(ctx, inv)
	default:
		return nil
	}
}

func (m *ProcessManager) extendSubscription(ctx context.Context, inv *invoicing.InvoiceFtLine) (_err error) {
	for _, invLine := range inv.Lines {
		querySubr := &subscription.GetSubscriptionByIDQuery{
			ID:        invLine.ReferralID,
			AccountID: inv.AccountID,
		}
		if err := m.subrQuery.Dispatch(ctx, querySubr); err != nil {
			return err
		}
		subr := querySubr.Result

		// assume that only has 1 line in subscription_line
		var periodStartAt, periodEndAt time.Time
		for _, subrLine := range subr.Lines {
			if subrLine.PlanID == 0 {
				continue
			}
			queryPlan := &subscriptionplan.GetSubrPlanByIDQuery{
				ID: subrLine.PlanID,
			}
			if err := m.subrPlanQuery.Dispatch(ctx, queryPlan); err != nil {
				return err
			}
			plan := queryPlan.Result

			if subr.CurrentPeriodEndAt.IsZero() {
				periodStartAt = inv.CreatedAt
			} else {
				periodStartAt = subr.CurrentPeriodEndAt
			}
			duration, err := subrplanutils.CalcPlanDuration(plan.Interval, plan.IntervalCount)
			if err != nil {
				return err
			}
			periodEndAt = periodStartAt.Add(duration)
		}
		update := &subscription.UpdateSubscriptionPeriodCommand{
			ID:                   subr.ID,
			AccountID:            subr.AccountID,
			CurrentPeriodStartAt: periodStartAt,
			CurrentPeriodEndAt:   periodEndAt,
		}
		if subr.CurrentPeriodStartAt.IsZero() {
			update.StartedAt = periodStartAt
		}
		if err := m.subrAggr.Dispatch(ctx, update); err != nil {
			return err
		}
	}
	return nil
}

/*
	WsWebsiteCreatedEvent

	Create subscription with 1 month free plan
*/

func (m *ProcessManager) WsWebsiteCreated(ctx context.Context, event *webserver.WsWebsiteCreatedEvent) error {
	query := &identity.GetShopByIDQuery{
		ID: event.ShopID,
	}
	if err := m.identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}
	shop := query.Result

	userQuery := &identity.GetUserByIDQuery{
		UserID: shop.OwnerID,
	}
	if err := m.identityQuery.Dispatch(ctx, userQuery); err != nil {
		return err
	}
	user := userQuery.Result

	planQuery := &subscriptionplan.GetFreeSubrPlanByProductTypeQuery{
		ProductType: subscription_product_type.Ecomify,
	}
	if err := m.subrPlanQuery.Dispatch(ctx, planQuery); err != nil {
		if cm.ErrorCode(err) == cm.NotFound {
			return cm.Errorf(cm.NotFound, nil, "Vui lòng tạo gói dùng thử cho ecomify")
		}
		return err
	}
	plan := planQuery.Result

	cmd := &subscription.CreateSubscriptionCommand{
		AccountID: event.ShopID,
		Lines: []*subscription.SubscriptionLine{
			{
				PlanID:   plan.ID,
				Quantity: 1,
			},
		},
		Customer: &subscriptingtypes.CustomerInfo{
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}
	if err := m.subrAggr.Dispatch(ctx, cmd); err != nil {
		return err
	}
	subr := cmd.Result

	// create subr bill & confirm to active trial product
	billCmd := &invoicing.CreateInvoiceBySubrIDCommand{
		SubscriptionID: subr.ID,
		AccountID:      shop.ID,
		TotalAmount:    0,
		Description:    "Trial",
		Classify:       service_classify.Shipping,
	}
	if err := m.invoiceAggr.Dispatch(ctx, billCmd); err != nil {
		return err
	}

	billConfirmCmd := &invoicing.PaymentInvoiceCommand{
		InvoiceID:     billCmd.Result.ID,
		AccountID:     shop.ID,
		TotalAmount:   0,
		PaymentMethod: payment_method.Manual,
	}
	if err := m.invoiceAggr.Dispatch(ctx, billConfirmCmd); err != nil {
		return err
	}

	return nil
}
