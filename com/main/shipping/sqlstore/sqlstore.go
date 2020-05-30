package sqlstore

import (
	"context"

	ordertypes "o.o/api/main/ordering/types"
	"o.o/api/main/shipping"
	"o.o/api/meta"
	shippingstate "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/status5"
	addressconvert "o.o/backend/com/main/address/convert"
	"o.o/backend/com/main/shipping/convert"
	"o.o/backend/com/main/shipping/model"
	shippingmodely "o.o/backend/com/main/shipping/modely"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
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

func (s *FulfillmentStore) OptionalID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByID(id).Optional())
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

func (s *FulfillmentStore) OptionalShippingCode(code string) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShippingCode(code).Optional())
	return s
}

func (s *FulfillmentStore) ShippingCodes(codes []string) *FulfillmentStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "shipping_code", codes))
	return s
}

func (s *FulfillmentStore) ShopID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *FulfillmentStore) OptionalShopID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShopID(id).Optional())
	return s
}

func (s *FulfillmentStore) PartnerID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id))
	return s
}

func (s *FulfillmentStore) OptionalPartnerID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(id).Optional())
	return s
}

func (s *FulfillmentStore) MoneyTxShippingID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionID(id))
	return s
}

func (s *FulfillmentStore) MoneyTxShippingIDs(ids ...dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "money_transaction_id", ids))
	return s
}

func (s *FulfillmentStore) MoneyTxShippingExternalID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionShippingExternalID(id))
	return s
}

func (s *FulfillmentStore) OptionalMoneyTxShippingID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionID(id).Optional())
	return s
}

func (s *FulfillmentStore) OptionalMoneyTxShippingExternalID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByMoneyTransactionShippingExternalID(id).Optional())
	return s
}

func (s *FulfillmentStore) OrderID(id dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByOrderID(id))
	return s
}

func (s *FulfillmentStore) OrderIDs(ids ...dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, sq.In("order_id", ids))
	return s
}

func (s *FulfillmentStore) ConnectionIDs(ids ...dot.ID) *FulfillmentStore {
	s.preds = append(s.preds, sq.In("connection_id", ids))
	return s
}

func (s *FulfillmentStore) ShippingProvider(carrier shipping_provider.ShippingProvider) *FulfillmentStore {
	s.preds = append(s.preds, s.ft.ByShippingProvider(carrier))
	return s
}

func (s *FulfillmentStore) ShippingStates(states ...shippingstate.State) *FulfillmentStore {
	s.preds = append(s.preds, sq.PrefixedIn(&s.ft.prefix, "shipping_state", states))
	return s
}

func (s *FulfillmentStore) NotBelongToMoneyTx() *FulfillmentStore {
	preds := sq.And{
		sq.NewIsNullPart("money_transaction_id", true),
		sq.NewIsNullPart("cod_etop_transfered_at", true),
	}
	s.preds = append(s.preds, preds)
	return s
}

func (s *FulfillmentStore) StatusNotIn(statuses ...status5.Status) *FulfillmentStore {
	s.preds = append(s.preds, sq.NotIn("status", statuses))
	return s
}

func (s *FulfillmentStore) FilterForMoneyTx(isNoneCOD bool, states []shippingstate.State) *FulfillmentStore {
	var preds sq.WriterTo
	if isNoneCOD {
		amount := 0
		preds = sq.Or{
			sq.PrefixedIn(&s.ft.prefix, "shipping_state", states),
			sq.And{
				s.ft.ByShippingState(shippingstate.Delivered),
				s.ft.ByTotalCODAmountPtr(&amount),
			},
		}
	} else {
		preds = sq.PrefixedIn(&s.ft.prefix, "shipping_state", states)
	}
	s.preds = append(s.preds, preds)
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

func (s *FulfillmentStore) ListFfms() ([]*shipping.Fulfillment, error) {
	ffmDBs, err := s.ListFfmsDB()
	if err != nil {
		return nil, err
	}
	ffms := []*shipping.Fulfillment{}
	err = scheme.Convert(ffmDBs, &ffms)
	return ffms, err
}

func (s *FulfillmentStore) CreateFulfillmentDB(ffm *model.Fulfillment) (*model.Fulfillment, error) {
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

func (s *FulfillmentStore) CreateFulfillmentsDB(ffms []*model.Fulfillment) error {
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

func (s *FulfillmentStore) UpdateFulfillmentDB(ffm *model.Fulfillment) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	return s.query().Where(s.preds).ShouldUpdate(ffm)
}

func (s *FulfillmentStore) UpdateFulfillment(ffm *shipping.Fulfillment) error {
	var ffmDB model.Fulfillment
	if err := scheme.Convert(ffm, &ffmDB); err != nil {
		return err
	}
	return s.UpdateFulfillmentDB(&ffmDB)
}

func (s *FulfillmentStore) UpdateFulfillmentsDB(ffms []*model.Fulfillment) error {
	for _, ffm := range ffms {
		if err := ffm.BeforeInsert(); err != nil {
			return err
		}
		if err := s.query().Where(s.ft.ByID(ffm.ID)).ShouldUpdate(ffm); err != nil {
			return err
		}
	}
	return nil
}

func (s *FulfillmentStore) UpdateFulfillmentShippingState(args *shipping.UpdateFulfillmentShippingStateArgs, codAmount int) error {
	if args.ActualCompensationAmount.Valid {
		update := map[string]interface{}{
			"shipping_state":             args.ShippingState.String(),
			"updated_by":                 args.UpdatedBy,
			"actual_compensation_amount": args.ActualCompensationAmount.Apply(codAmount),
		}
		return s.query().Table("fulfillment").Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdateMap(update)
	}

	update := &model.Fulfillment{
		ShippingState: args.ShippingState,
		UpdatedBy:     args.UpdatedBy,
	}
	return s.query().Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdate(update)
}

func (s *FulfillmentStore) UpdateFulfillmentsMoneyTxID(args *shipping.UpdateFulfillmentsMoneyTxIDArgs) (updated int, _ error) {
	update := &model.Fulfillment{
		MoneyTransactionShippingExternalID: args.MoneyTxShippingExternalID,
		MoneyTransactionID:                 args.MoneyTxShippingID,
	}
	return s.IDs(args.FulfillmentIDs...).query().Where(s.preds).Update(update)
}

func (s *FulfillmentStore) RemoveFulfillmentsMoneyTxID(args *shipping.RemoveFulfillmentsMoneyTxIDArgs) (updated int, _ error) {
	count := 0
	if len(args.FulfillmentIDs) > 0 {
		s = s.IDs(args.FulfillmentIDs...)
		count++
	}
	if args.MoneyTxShippingID != 0 {
		s = s.MoneyTxShippingID(args.MoneyTxShippingID)
		count++
	}
	if args.MoneyTxShippingExternalID != 0 {
		s = s.MoneyTxShippingExternalID(args.MoneyTxShippingExternalID)
		count++
	}
	if count == 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Missing required fields").WithMetap("func", "RemoveFulfillmentsMoneyTxID")
	}
	return s.query().Table("fulfillment").Where(s.preds).UpdateMap(map[string]interface{}{
		"money_transaction_id":                   nil,
		"money_transaction_shipping_external_id": nil,
	})
}

// UpdateFulfillmentInfo
//
// Fullname, Phone là của người nhận
// Sẽ tách thành update thông tin người gửi, người nhận riêng. Cập nhật sau.
func (s *FulfillmentStore) UpdateFulfillmentInfo(args *shipping.UpdateFulfillmentInfoArgs, oldAddress *ordertypes.Address) error {
	if len(s.preds) == 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	update := &model.Fulfillment{
		AdminNote: args.AdminNote,
	}
	address := addressconvert.OrderAddressToModel(oldAddress)
	update.AddressTo = address.UpdateAddress(args.Phone.String, args.FullName.String)
	return s.query().Table("fulfillment").Where(s.preds).ShouldUpdate(update)
}

func (s *FulfillmentStore) UpdateFulfillmentsStatus(args *shipping.UpdateFulfillmentsStatusArgs) error {
	update := map[string]interface{}{}
	if args.Status.Valid {
		update["status"] = args.Status.Enum
	}
	if args.ShopConfirm.Valid {
		update["shop_confirm"] = args.ShopConfirm.Enum
	}
	if args.SyncStatus.Valid {
		update["sync_status"] = args.SyncStatus.Enum
	}
	return s.query().Table("fulfillment").Where(sq.In("id", args.FulfillmentIDs)).ShouldUpdateMap(update)
}

func (s *FulfillmentStore) CancelFulfillment(args *shipping.CancelFulfillmentArgs) error {
	update := &model.Fulfillment{
		ShopConfirm:   status3.N,
		CancelReason:  args.CancelReason,
		ShippingState: shippingstate.Cancelled,
	}
	return s.query().Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdate(update)
}

// FulfillmentExtended
func (s *FulfillmentStore) GetFulfillmentExtendedDB() (*shippingmodely.FulfillmentExtended, error) {
	query := s.query().Where(s.preds)
	ffm := &shippingmodely.FulfillmentExtended{}
	err := query.ShouldGet(ffm)
	return ffm, err
}

func (s *FulfillmentStore) GetFulfillmentExtended() (*shipping.FulfillmentExtended, error) {
	ffmDB, err := s.GetFulfillmentExtendedDB()
	if err != nil {
		return nil, err
	}
	var ffm shipping.FulfillmentExtended
	if err := scheme.Convert(ffmDB, ffm); err != nil {
		return nil, err
	}
	return &ffm, nil
}

func (s *FulfillmentStore) ListFulfillmentExtendedsDB() ([]*shippingmodely.FulfillmentExtended, error) {
	var ffms shippingmodely.FulfillmentExtendeds
	query := s.query().Where(s.preds...)
	if len(s.Paging.Sort) == 0 {
		s.Paging.Sort = append(s.Paging.Sort, "-created_at")
	}
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

func (s *FulfillmentStore) ListFulfillmentExtendeds() (res []*shipping.FulfillmentExtended, _ error) {
	ffmsDB, err := s.ListFulfillmentExtendedsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(ffmsDB, &res); err != nil {
		return nil, err
	}
	return
}

type ForceUpdateExternalShippingInfoArgs struct {
	FulfillmentID            dot.ID
	ExternalShippingNote     dot.NullString
	ExternalShippingSubState dot.NullString
}

func (s *FulfillmentStore) ForceUpdateExternalShippingInfo(args *ForceUpdateExternalShippingInfoArgs) error {
	update := map[string]interface{}{}
	if args.ExternalShippingNote.Valid {
		update["external_shipping_note"] = args.ExternalShippingNote.String
	}
	if args.ExternalShippingSubState.Valid {
		update["external_shipping_sub_state"] = args.ExternalShippingSubState.String
	}
	if len(update) == 0 {
		return nil
	}
	return s.query().Table("fulfillment").Where(s.ft.ByID(args.FulfillmentID)).ShouldUpdateMap(update)
}
