package sqlstore

import (
	"context"
	"encoding/json"

	"o.o/api/external/payment"
	"o.o/api/top/types/etc/payment_provider"
	"o.o/api/top/types/etc/payment_state"
	"o.o/api/top/types/etc/status4"
	"o.o/backend/com/external/payment/payment/convert"
	"o.o/backend/com/external/payment/payment/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

type PaymentStoreFactory func(context.Context) *PaymentStore

func NewPaymentStore(db *cmsql.Database) PaymentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *PaymentStore {
		return &PaymentStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type PaymentStore struct {
	query cmsql.QueryFactory
	ft    PaymentFilters
	preds []interface{}
}

func (s *PaymentStore) ID(id dot.ID) *PaymentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PaymentStore) ExternalTransactionID(id string) *PaymentStore {
	s.preds = append(s.preds, s.ft.ByExternalTransID(id))
	return s
}

func (s *PaymentStore) PaymentProvider(provider payment_provider.PaymentProvider) *PaymentStore {
	s.preds = append(s.preds, s.ft.ByPaymentProvider(provider))
	return s
}

func (s *PaymentStore) GetPaymentDB() (*model.Payment, error) {
	var payment model.Payment
	err := s.query().Where(s.preds...).ShouldGet(&payment)
	return &payment, err
}

func (s *PaymentStore) GetPayment() (*payment.Payment, error) {
	payment, err := s.GetPaymentDB()
	if err != nil {
		return nil, err
	}
	return convert.Payment(payment), nil
}

type CreatePaymentArgs struct {
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	PaymentProvider payment_provider.PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
}

func (s *PaymentStore) CreatePayment(args *CreatePaymentArgs) (*payment.Payment, error) {
	id := cm.NewID()
	payment := &model.Payment{
		ID:              id,
		Amount:          args.Amount,
		Status:          args.Status,
		State:           args.State,
		PaymentProvider: args.PaymentProvider,
		ExternalTransID: args.ExternalTransID,
		ExternalData:    args.ExternalData,
	}
	if err := s.query().ShouldInsert(payment); err != nil {
		return nil, err
	}
	return s.ID(id).GetPayment()
}

type UpdateExternalPaymentInfoArgs struct {
	ID              dot.ID
	Amount          int
	Status          status4.Status
	State           payment_state.PaymentState
	ExternalData    json.RawMessage
	ExternalTransID string
}

func (s *PaymentStore) UpdateExternalPaymentInfo(args *UpdateExternalPaymentInfoArgs) (*payment.Payment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Payment ID")
	}
	update := &model.Payment{
		Amount:          args.Amount,
		Status:          args.Status,
		State:           args.State,
		ExternalData:    args.ExternalData,
		ExternalTransID: args.ExternalTransID,
	}
	if err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetPayment()
}
