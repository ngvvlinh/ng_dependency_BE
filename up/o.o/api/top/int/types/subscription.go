package types

import (
	"time"

	"o.o/api/top/types/common"
	"o.o/api/top/types/etc/invoice_type"
	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/subject_referral"
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

type GetSubrProductsRequest struct {
	Type subscription_product_type.ProductSubscriptionType `json:"type"`
}

func (m *GetSubrProductsRequest) String() string { return jsonx.MustMarshalToString(m) }

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

type GetSubrPlansRequest struct {
	ProductID dot.ID `json:"product_id"`
}

func (m *GetSubrPlansRequest) String() string { return jsonx.MustMarshalToString(m) }

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
	CancelAtPeriodEnd    dot.NullBool        `json:"cancel_at_period_end"`
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

type Invoice struct {
	ID            dot.ID                           `json:"id"`
	AccountID     dot.ID                           `json:"account_id"`
	TotalAmount   int                              `json:"total_amount"`
	Description   string                           `json:"description"`
	PaymentID     dot.ID                           `json:"payment_id"`
	Status        status4.Status                   `json:"status"`
	PaymentStatus status4.Status                   `json:"payment_status"`
	Customer      *SubrCustomer                    `json:"customer"`
	Lines         []*InvoiceLine                   `json:"lines"`
	CreatedAt     dot.Time                         `json:"created_at"`
	UpdatedAt     dot.Time                         `json:"updated_at"`
	ReferralType  subject_referral.SubjectReferral `json:"referral_type"`
	ReferralIDs   []dot.ID                         `json:"referral_ids"`
	Classify      service_classify.ServiceClassify `json:"classify"`
	Type          invoice_type.InvoiceType         `json:"type"`
}

func (m *Invoice) String() string { return jsonx.MustMarshalToString(m) }

type InvoiceLine struct {
	ID           dot.ID                           `json:"id"`
	LineAmount   int                              `json:"line_amount"`
	Price        int                              `json:"price"`
	Quantity     int                              `json:"quantity"`
	Description  string                           `json:"description"`
	InvoiceID    dot.ID                           `json:"invoice_id"`
	ReferralType subject_referral.SubjectReferral `json:"referral_type"`
	ReferralID   dot.ID                           `json:"referral_id"`
	CreatedAt    dot.Time                         `json:"created_at"`
	UpdatedAt    dot.Time                         `json:"updated_at"`
}

func (m *InvoiceLine) String() string { return jsonx.MustMarshalToString(m) }

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
	// Bỏ trống field này nếu shop tạo subscription
	AccountID dot.ID `json:"account_id"`
	// Hủy subscription khi hết hạn
	CancelAtPeriodEnd bool `json:"cancel_at_period_end"`
	// thời điểm phát sinh hóa đơn (trường hợp tự động gia hạn)
	BillingCycleAnchorAt dot.Time            `json:"billing_cycle_anchor_at"`
	Lines                []*SubscriptionLine `json:"lines"`
	Customer             *SubrCustomer       `json:"customer"`
}

func (m *CreateSubscriptionRequest) String() string { return jsonx.MustMarshalToString(m) }

type UpdateSubscriptionInfoRequest struct {
	ID        dot.ID `json:"id"`
	AccountID dot.ID `json:"account_id"`
	// Hủy subscription khi hết hạn
	CancelAtPeriodEnd dot.NullBool `json:"cancel_at_period_end"`
	// thời điểm phát sinh hóa đơn (trường hợp tự động gia hạn)
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

type CreateInvoiceForSubscriptionRequest struct {
	AccountID      dot.ID                           `json:"account_id"`
	SubscriptionID dot.ID                           `json:"subscription_id"`
	TotalAmount    int                              `json:"total_amount"`
	Description    string                           `json:"description"`
	Customer       *SubrCustomer                    `json:"customer"`
	Classify       service_classify.ServiceClassify `json:"classify"`
}

func (m *CreateInvoiceForSubscriptionRequest) String() string { return jsonx.MustMarshalToString(m) }

type ManualPaymentInvoiceRequest struct {
	InvoiceID   dot.ID `json:"invoice_id"`
	AccountID   dot.ID `json:"account_id"`
	TotalAmount int    `json:"total_amount"`
}

func (m *ManualPaymentInvoiceRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInvoicesRequest struct {
	Paging *common.CursorPaging   `json:"paging"`
	Filter *GetShopInvoicesFilter `json:"filter"`
}

func (m *GetInvoicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetInvoicesResponse struct {
	Invoices []*Invoice             `json:"invoices"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *GetInvoicesResponse) String() string { return jsonx.MustMarshalToString(m) }

type GetShopInvoicesRequest struct {
	Paging *common.CursorPaging   `json:"paging"`
	Filter *GetShopInvoicesFilter `json:"filter"`
}

func (m *GetShopInvoicesRequest) String() string { return jsonx.MustMarshalToString(m) }

type GetShopInvoicesFilter struct {
	AccountID dot.ID                           `json:"account_id"`
	RefID     dot.ID                           `json:"ref_id"`
	RefType   subject_referral.SubjectReferral `json:"ref_type"`
	DateFrom  time.Time                        `json:"date_from"`
	DateTo    time.Time                        `json:"date_to"`
	Type      invoice_type.InvoiceType         `json:"type"`
}

func (m *GetShopInvoicesFilter) String() string { return jsonx.MustMarshalToString(m) }

type GetShopInvoicesResponse struct {
	Invoices []*Invoice             `json:"invoices"`
	Paging   *common.CursorPageInfo `json:"paging"`
}

func (m *GetShopInvoicesResponse) String() string { return jsonx.MustMarshalToString(m) }
