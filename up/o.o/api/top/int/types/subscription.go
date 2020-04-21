package types

import (
	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subscription_plan_interval"
	"o.o/api/top/types/etc/subscription_product_type"
	"o.o/capi/dot"
	"o.o/common/jsonx"
)

type SubscriptionProduct struct {
	ID          dot.ID                                            `json:"id"`
	Name        string                                            `json:"name"`
	Type        subscription_product_type.ProductSubscriptionType `json:"type"`
	Description string                                            `json:"description"`
	ImageURL    string                                            `json:"image_url"`
	Status      status3.Status                                    `json:"status"`
	CreatedAt   dot.Time                                          `json:"created_at"`
	UpdatedAt   dot.Time                                          `json:"updated_at"`
}

func (m *SubscriptionProduct) String() string { return jsonx.MustMarshalToString(m) }

type CreateSubrProductRequest struct {
	Name        string                                            `json:"name"`
	Description string                                            `json:"description"`
	ImageURL    string                                            `json:"image_url"`
	Type        subscription_product_type.ProductSubscriptionType `json:"type"`
}

func (m *CreateSubrProductRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSubrProductsResponse struct {
	SubscriptionProducts []*SubscriptionProduct `json:"subscription_products"`
}

func (m *GetSubrProductsResponse) String() string { return jsonx.MustMarshalToString(m) }

type SubscriptionPlan struct {
	ID            dot.ID                                              `json:"id"`
	Name          string                                              `json:"name"`
	Price         int                                                 `json:"price"`
	Status        status3.Status                                      `json:"status"`
	Description   string                                              `json:"description"`
	ProductID     dot.ID                                              `json:"product_id"`
	Interval      subscription_plan_interval.SubscriptionPlanInterval `json:"interval"`
	IntervalCount int                                                 `json:"interval_count"`
	CreatedAt     dot.Time                                            `json:"created_at"`
	UpdatedAt     dot.Time                                            `json:"updated_at"`
}

func (m *SubscriptionPlan) String() string { return jsonx.MustMarshalToString(m) }

type CreateSubrPlanRequest struct {
	Name          string                                              `json:"name"`
	Price         int                                                 `json:"price"`
	Description   string                                              `json:"description"`
	ProductID     dot.ID                                              `json:"product_id"`
	Interval      subscription_plan_interval.SubscriptionPlanInterval `json:"interval"`
	IntervalCount int                                                 `json:"interval_count"`
}

func (m *CreateSubrPlanRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateSubrPlanRequest struct {
	ID            dot.ID                                              `json:"id"`
	Name          string                                              `json:"name"`
	Price         int                                                 `json:"price"`
	Description   string                                              `json:"description"`
	ProductID     dot.ID                                              `json:"product_id"`
	Interval      subscription_plan_interval.SubscriptionPlanInterval `json:"interval"`
	IntervalCount int                                                 `json:"interval_count"`
}

func (m *UpdateSubrPlanRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSubrPlansResponse struct {
	SubscriptionPlans []*SubscriptionPlan `json:"subscription_plans"`
}

func (m *GetSubrPlansResponse) String() string { return jsonx.MustMarshalToString(m) }

type Subscription struct {
	ID                   dot.ID              `json:"id"`
	AccountID            dot.ID              `json:"account_id"`
	CancelAtPeriodEnd    bool                `json:"cancel_at_period_end"`
	CurrentPeriodStartAt dot.Time            `json:"current_period_start_at"`
	CurrentPeriodEndAt   dot.Time            `json:"current_period_end_at"`
	Status               status3.Status      `json:"status"`
	BillingCycleAnchorAt dot.Time            `json:"billing_cycle_anchor_at"`
	StartedAt            dot.Time            `json:"started_at"`
	Lines                []*SubscriptionLine `json:"lines"`
	Customer             *SubrCustomer       `json:"customer"`
	CreatedAt            dot.Time            `json:"created_at"`
	UpdatedAt            dot.Time            `json:"updated_at"`
}

func (m *Subscription) String() string { return jsonx.MustMarshalToString(m) }

type SubrCustomer struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (m *SubrCustomer) String() string { return jsonx.MustMarshalToString(m) }

type SubscriptionLine struct {
	ID             dot.ID   `json:"id"`
	PlanID         dot.ID   `json:"plan_id"`
	SubscriptionID dot.ID   `json:"subscription_id"`
	Quantity       int      `json:"quantity"`
	CreatedAt      dot.Time `json:"created_at"`
	UpdatedAt      dot.Time `json:"updated_at"`
}

func (m *SubscriptionLine) String() string { return jsonx.MustMarshalToString(m) }

type SubscriptionBill struct {
	ID             dot.ID          `json:"id"`
	AccountID      dot.ID          `json:"account_id"`
	SubscriptionID dot.ID          `json:"subscription_id"`
	TotalAmount    int             `json:"total_amount"`
	Description    string          `json:"description"`
	PaymentID      dot.ID          `json:"payment_id"`
	Status         status4.Status  `json:"status"`
	PaymentStatus  status4.Status  `json:"payment_status"`
	Customer       *SubrCustomer   `json:"customer"`
	Lines          []*SubrBillLine `json:"lines"`
	CreatedAt      dot.Time        `json:"created_at"`
	UpdatedAt      dot.Time        `json:"updated_at"`
}

func (m *SubscriptionBill) String() string { return jsonx.MustMarshalToString(m) }

type SubrBillLine struct {
	ID                 dot.ID   `json:"id"`
	LineAmount         int      `json:"line_amount"`
	Price              int      `json:"price"`
	Quantity           int      `json:"quantity"`
	Description        string   `json:"description"`
	PeriodStartAt      dot.Time `json:"period_start_at"`
	PeriodEndAt        dot.Time `json:"period_end_at"`
	SubscriptionID     dot.ID   `json:"subscription_id"`
	SubscriptionBillID dot.ID   `json:"subscription_bill_id"`
	SubscriptionLineID dot.ID   `json:"subscription_line_id"`
	CreatedAt          dot.Time `json:"created_at"`
	UpdatedAt          dot.Time `json:"updated_at"`
}

func (m *SubrBillLine) String() string { return jsonx.MustMarshalToString(m) }

type GetSubscriptionsRequest struct {
	AccountID dot.ID           `json:"account_id"`
	Paging    *common.Paging   `json:"paging"`
	Filters   []*common.Filter `json:"filters"`
}

func (m *GetSubscriptionsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSubscriptionsResponse struct {
	Subscriptions []*Subscription  `json:"subscriptions"`
	Paging        *common.PageInfo `json:"paging"`
}

func (m *GetSubscriptionsResponse) String() string { return jsonx.MustMarshalToString(m) }

type CreateSubscriptionRequest struct {
	AccountID            dot.ID              `json:"account_id"`
	CancelAtPeriodEnd    bool                `json:"cancel_at_period_end"`
	BillingCycleAnchorAt dot.Time            `json:"billing_cycle_anchor_at"`
	Lines                []*SubscriptionLine `json:"lines"`
	Customer             *SubrCustomer       `json:"customer"`
}

func (m *CreateSubscriptionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateSubscriptionInfoRequest struct {
	ID                   dot.ID              `json:"id"`
	AccountID            dot.ID              `json:"account_id"`
	CancelAtPeriodEnd    bool                `json:"cancel_at_period_end"`
	BillingCycleAnchorAt dot.Time            `json:"billing_cycle_anchor_at"`
	Customer             *SubrCustomer       `json:"customer"`
	Lines                []*SubscriptionLine `json:"lines"`
}

func (m *UpdateSubscriptionInfoRequest) String() string { return jsonx.MustMarshalToString(m) }

type SubscriptionIDRequest struct {
	ID        dot.ID `json:"id"`
	AccountID dot.ID `json:"account_id"`
}

func (m *SubscriptionIDRequest) String() string { return jsonx.MustMarshalToString(m) }

type CreateSubscriptionBillRequest struct {
	AccountID      dot.ID        `json:"account_id"`
	SubscriptionID dot.ID        `json:"subscription_id"`
	TotalAmount    int           `json:"total_amount"`
	Description    string        `json:"description"`
	Customer       *SubrCustomer `json:"customer"`
}

func (m *CreateSubscriptionBillRequest) String() string { return jsonx.MustMarshalToString(m) }

type ManualPaymentSubscriptionBillRequest struct {
	SubscriptionBillID dot.ID `json:"subscription_bill_id"`
	AccountID          dot.ID `json:"account_id"`
	TotalAmount        int    `json:"total_amount"`
}

func (m *ManualPaymentSubscriptionBillRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSubscriptionBillsRequest struct {
	AccountID dot.ID           `json:"account_id"`
	Paging    *common.Paging   `json:"paging"`
	Filters   []*common.Filter `json:"filters"`
}

func (m *GetSubscriptionBillsRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetSubscriptionBillsResponse struct {
	SubscriptionBills []*SubscriptionBill `json:"subscription_bills"`
	Paging            *common.PageInfo    `json:"paging"`
}

func (m *GetSubscriptionBillsResponse) String() string { return jsonx.MustMarshalToString(m) }
