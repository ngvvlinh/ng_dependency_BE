package sqlstore

import (
	"context"
	"encoding/json"

	"etop.vn/api/external/payment"
	"etop.vn/api/main/etop"
	"etop.vn/backend/com/external/payment/payment/convert"
	"etop.vn/backend/com/external/payment/payment/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
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

func (s *PaymentStore) ID(id int64) *PaymentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *PaymentStore) ExternalTransactionID(id string) *PaymentStore {
	s.preds = append(s.preds, s.ft.ByExternalTransID(id))
	return s
}

func (s *PaymentStore) PaymentProvider(provider string) *PaymentStore {
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
	Status          etop.Status4
	State           payment.PaymentState
	PaymentProvider payment.PaymentProvider
	ExternalTransID string
	ExternalData    json.RawMessage
}

func (s *PaymentStore) CreatePayment(args *CreatePaymentArgs) (*payment.Payment, error) {
	id := cm.NewID()
	payment := &model.Payment{
		ID:              id,
		Amount:          args.Amount,
		Status:          int(args.Status),
		State:           string(args.State),
		PaymentProvider: string(args.PaymentProvider),
		ExternalTransID: args.ExternalTransID,
		ExternalData:    args.ExternalData,
	}
	if err := s.query().ShouldInsert(payment); err != nil {
		return nil, err
	}
	return s.ID(id).GetPayment()
}

type UpdateExternalPaymentInfoArgs struct {
	ID              int64
	Amount          int
	Status          etop.Status4
	State           payment.PaymentState
	ExternalData    json.RawMessage
	ExternalTransID string
}

func (s *PaymentStore) UpdateExternalPaymentInfo(args *UpdateExternalPaymentInfoArgs) (*payment.Payment, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Payment ID")
	}
	update := &model.Payment{
		Amount:          args.Amount,
		Status:          int(args.Status),
		State:           string(args.State),
		ExternalData:    args.ExternalData,
		ExternalTransID: args.ExternalTransID,
	}
	if err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetPayment()
}
