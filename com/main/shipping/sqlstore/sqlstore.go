package sqlstore

import (
	"context"

	"etop.vn/api/main/shipping"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/status5"
	"etop.vn/backend/com/main/shipping/convert"
	"etop.vn/backend/com/main/shipping/model"
	shippingsharemodel "etop.vn/backend/com/main/shipping/sharemodel"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

type FulfillmentStoreFactory func(context.Context) *FulfillmentStore

func NewFulfillmentStore(db *cmsql.Database) FulfillmentStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *FulfillmentStore {
		return &FulfillmentStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type FulfillmentStore struct {
	ft FulfillmentFilters

	query cmsql.QueryFactory
	preds []interface{}
	sqlstore.Paging

	includeDeleted bool
}

func (s *FulfillmentStore) WithPaging(paging meta.Paging) *FulfillmentStore {
	s.Paging.WithPaging(paging)
	return s
}

func (s *FulfillmentStore) ID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *FulfillmentStore) IDs(ids ...dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "id", ids))
	return s
}

func (s *FulfillmentStore) ShippingCode(code string) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code))
	return s
}

func (s *FulfillmentStore) IDOrShippingCode(id dot.ID, shippingCode string) *FulfillmentStore {
	s.preds = append(s.preds, sq.Once{
		s.ft.ByID(id),
		s.ft.ByShippingCode(shippingCode),
	})
	return s
}

func (s *FulfillmentStore) ShopID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *FulfillmentStore) PartnerID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *FulfillmentStore) OrderID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *FulfillmentStore) GetFfmDB() (*model.Fulfillment, error) {
	var ffm model.Fulfillment
	err := s.query().Where(s.preds...).ShouldGet(&ffm)
	return &ffm, err
}

func (s *FulfillmentStore) GetFulfillment() (*shipping.Fulfillment, error) {
	ffmDB, err := s.GetFfmDB()
	if err != nil {
		return nil, err
	}
	ffm := &shipping.Fulfillment{}
	err = scheme.Convert(ffmDB, ffm)
	return ffm, err
}

func (s *FulfillmentStore) ListFfmsDB() ([]*model.Fulfillment, error) {
	var ffms model.Fulfillments
	query := s.query().Where(s.preds...)
	query, err := sqlstore.LimitSort(query, &s.Paging, SortFulfillment)
	if err != nil {
		return nil, err
	}
	err = query.Find(&ffms)
	if err != nil {
		return nil, err
	}
	s.Paging.Apply(ffms)
	return ffms, nil
}

func (s *FulfillmentStore) CreateFulfillmentDB(ctx context.Context, ffm *model.Fulfillment) (*model.Fulfillment, error) {
	if ffm.ID == 0 {
		ffm.ID = cm.NewID()
	}
	if err := ffm.BeforeInsert(); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(ffm); err != nil {
		return nil, err
	}
	return s.ID(ffm.ID).GetFfmDB()
}

func (s *FulfillmentStore) CreateFulfillmentsDB(ctx context.Context, ffms []*model.Fulfillment) error {
	for _, ffm := range ffms {
		if err := ffm.BeforeInsert(); err != nil {
			return err
		}
		if err := s.query().ShouldInsert(ffm); err != nil {
			return err
		}
	}
	return nil
}

func (s *FulfillmentStore) UpdateFulfillmentDB(ctx context.Context, ffm *model.Fulfillment) (*model.Fulfillment, error) {
	if err := s.query().Where(s.ft.ByID(ffm.ID)).Where("status not in (?, ?, ?)", status5.N, status5.NS, status5.P).ShouldUpdate(ffm); err != nil {
		return nil, err
	}
	return s.ID(ffm.ID).GetFfmDB()
}

func (s *FulfillmentStore) UpdateFulfillmentsDB(ctx context.Context, ffms []*model.Fulfillment) error {
	for _, ffm := range ffms {
		if err := ffm.BeforeInsert(); err != nil {
			return err
		}
		if err := s.query().Where(s.ft.ByID(ffm.ID)).Where("status not in (?, ?, ?)", status5.N, status5.NS, status5.P).ShouldUpdate(ffm); err != nil {
			return err
		}
	}
	return nil
}

func (s *FulfillmentStore) UpdateFulfillmentShippingState(args *shipping.UpdateFulfillmentShippingStateArgs) error {
	if args.ActualCompensationAmount.Valid {
		update := map[string]interface{}{
			"shipping_state":             args.ShippingState.String(),
			"actual_compensation_amount": args.ActualCompensationAmount.Apply(0),
		}
		return s.query().Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdateMap(update)
	}

	update := &model.Fulfillment{
		ShippingState: args.ShippingState,
	}
	return s.query().Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdate(update)
}

func (s *FulfillmentStore) UpdateFulfillmentShippingFees(args *shipping.UpdateFulfillmentShippingFeesArgs) error {
	var lines []*shippingsharemodel.ShippingFeeLine
	if err := scheme.Convert(args.ShippingFeeLines, lines); err != nil {
		return err
	}

	update := &model.Fulfillment{
		ShippingFeeShopLines: lines,
		ShippingFeeShop:      shippingsharemodel.GetTotalShippingFee(lines),
	}
	return s.query().Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdate(update)
}

func (s *FulfillmentStore) UpdateFulfillmentsMoneyTxShippingExternalID(args *shipping.UpdateFulfillmentsMoneyTxShippingExternalIDArgs) error {
	update := &model.Fulfillment{
		MoneyTransactionShippingExternalID: args.MoneyTxShippingExternalID,
	}
	return s.IDs(args.FulfillmentIDs...).query().Where(s.preds).ShouldUpdate(update)
}
