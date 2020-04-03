package aggregate

import (
	"context"
	"strconv"
	"strings"

	"etop.vn/api/main/moneytx"
	"etop.vn/api/main/shipping"
	"etop.vn/api/meta"
	"etop.vn/api/top/types/etc/connection_type"
	shippingstate "etop.vn/api/top/types/etc/shipping"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi/dot"
)

var acceptStates = []string{
	shippingstate.Returned.String(), shippingstate.Returning.String(), shippingstate.Delivered.String(), shippingstate.Undeliverable.String(),
}

func (a *MoneyTxAggregate) CreateMoneyTxShippingExternal(ctx context.Context, args *moneytx.CreateMoneyTxShippingExternalArgs) (*moneytx.MoneyTransactionShippingExternalExtended, error) {
	if args.Provider == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vui lòng chọn nhà vận chuyển")
	}
	if len(args.Lines) == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Vận đơn không được rỗng")
	}
	totalCOD := 0
	totalOrders := 0
	for _, line := range args.Lines {
		totalCOD += line.ExternalTotalCOD
		totalOrders++
	}

	externalTxID := cm.NewID()
	err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		subCode := args.Provider.String()
		subCode = strings.ToUpper(subCode)
		code, errCode := sqlstore.GenerateCode(ctx, tx, model.CodeTypeMoneyTransactionExternal, subCode)
		if errCode != nil {
			return errCode
		}
		externalTx := &moneytx.MoneyTransactionShippingExternal{
			ID:             externalTxID,
			Code:           code,
			TotalCOD:       totalCOD,
			TotalOrders:    totalOrders,
			ExternalPaidAt: args.ExternalPaidAt,
			Provider:       args.Provider,
			BankAccount:    args.BankAccount,
			Note:           args.Note,
			InvoiceNumber:  args.InvoiceNumber,
		}
		if err := a.moneyTxShippingExternalStore(ctx).CreateMoneyTxShippingExternal(externalTx); err != nil {
			return err
		}

		ffmIDs := make([]dot.ID, 0, len(args.Lines))
		for _, line := range args.Lines {
			createCmd := &moneytx.CreateMoneyTxShippingExternalLineArgs{
				ExternalCode:                       line.ExternalCode,
				ExternalTotalCOD:                   line.ExternalTotalCOD,
				ExternalCreatedAt:                  line.ExternalCreatedAt,
				ExternalClosedAt:                   line.ExternalClosedAt,
				EtopFulfillmentIDRaw:               line.EtopFulfillmentIDRaw,
				ExternalCustomer:                   line.ExternalCustomer,
				ExternalAddress:                    line.ExternalAddress,
				MoneyTransactionShippingExternalID: externalTx.ID,
				ExternalTotalShippingFee:           line.ExternalTotalShippingFee,
			}
			externalLine, err := a.CreateMoneyTxShippingExternalLine(ctx, createCmd)
			if err != nil {
				return err
			}
			if externalLine.EtopFulfillmentID != 0 && externalLine.ImportError == nil {
				ffmIDs = append(ffmIDs, externalLine.EtopFulfillmentID)
			}
		}

		event := &moneytx.MoneyTransactionShippingExternalCreatedEvent{
			MoneyTxShippingExternalID: externalTxID,
			FulfillementIDs:           ffmIDs,
		}
		if err := a.eventBus.Publish(ctx, event); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return a.moneyTxShippingExternalStore(ctx).ID(externalTxID).GetMoneyTxShippingExternalExtended()
}

func (a *MoneyTxAggregate) CreateMoneyTxShippingExternalLine(ctx context.Context, args *moneytx.CreateMoneyTxShippingExternalLineArgs) (*moneytx.MoneyTransactionShippingExternalLine, error) {
	if args.MoneyTransactionShippingExternalID == 0 {
		return nil, cm.Error(cm.InvalidArgument, "Missing MoneyTransactionShippingExternalID", nil)
	}
	line := &moneytx.MoneyTransactionShippingExternalLine{
		ID:                                 cm.NewID(),
		ExternalCode:                       args.ExternalCode,
		ExternalTotalCOD:                   args.ExternalTotalCOD,
		ExternalCreatedAt:                  args.ExternalCreatedAt,
		ExternalClosedAt:                   args.ExternalClosedAt,
		ExternalCustomer:                   args.ExternalCustomer,
		ExternalAddress:                    args.ExternalAddress,
		EtopFulfillmentIDRaw:               args.EtopFulfillmentIDRaw,
		MoneyTransactionShippingExternalID: args.MoneyTransactionShippingExternalID,
		ExternalTotalShippingFee:           args.ExternalTotalShippingFee,
	}
	if line.ExternalCode == "" {
		line.ImportError = &meta.Error{
			Code: "ffm_id_empty",
			Msg:  "Thiếu mã vận đơn",
		}
	} else {
		query := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
			ShippingCode: line.ExternalCode,
			Result:       nil,
		}
		if err := a.shippingQuery.Dispatch(ctx, query); err != nil {
			line.ImportError = &meta.Error{
				Code: "ffm_not_found",
				Msg:  "Không tìm thấy vận đơn trên Etop",
			}
		} else {
			ffm := query.Result
			line.EtopFulfillmentID = ffm.ID
			if ffm.MoneyTransactionShippingExternalID != 0 {
				line.ImportError = &meta.Error{
					Code: "ffm_exist_money_transaction_shipping_external",
					Msg:  "Vận đơn nằm trong phiên thanh toán nhà vận chuyển khác: " + strconv.Itoa(int(ffm.MoneyTransactionShippingExternalID)),
				}
			} else if !cm.StringsContain(acceptStates, ffm.ShippingState.String()) {
				line.ImportError = &meta.Error{
					Code: "ffm_not_done",
					Msg:  "Vận đơn chưa hoàn thành trên Etop",
				}
			} else if ffm.ShippingState == shippingstate.Delivered && ffm.TotalCODAmount != line.ExternalTotalCOD {
				line.ImportError = &meta.Error{
					Code: "ffm_not_balance",
					Msg:  "Giá trị vận đơn không đúng",
					Meta: map[string]string{
						"Etop":     strconv.Itoa(ffm.TotalCODAmount),
						"Provider": strconv.Itoa(line.ExternalTotalCOD),
					},
				}
			} else if ffm.ShippingState == shippingstate.Undeliverable && line.ExternalTotalCOD != ffm.ActualCompensationAmount {
				line.ImportError = &meta.Error{
					Code: "ffm_not_balance",
					Msg:  "Giá trị bồi hoàn không đúng",
				}
			} else if ffm.MoneyTransactionID != 0 {
				line.ImportError = &meta.Error{
					Code: "ffm_exist_money_transaction",
					Msg:  "Vận đơn nằm trong phiên thanh toán khác: " + strconv.Itoa(int(ffm.MoneyTransactionID)),
				}
			} else if line.ExternalTotalShippingFee != 0 && line.ExternalTotalShippingFee != ffm.ShippingFeeShop {
				line.ImportError = &meta.Error{
					Code: "ffm_shipping_fee_not_match",
					Msg:  "Tổng tiền cước không đúng.",
					Meta: map[string]string{
						"Etop":     strconv.Itoa(ffm.ShippingFeeShop),
						"Provider": strconv.Itoa(line.ExternalTotalShippingFee),
					},
				}
			} else if ffm.ConnectionMethod != connection_type.ConnectionMethodBuiltin {
				if ffm.ShippingType == 0 {
					// backward compatible
					// remove later
					// no error
				} else {
					line.ImportError = &meta.Error{
						Code: "ffm_not_in_etop",
						Msg:  "Vận đơn không được đối soát bởi Etop",
					}
				}
			}
		}
	}
	if err := a.moneyTxShippingExternalStore(ctx).CreateMoneyTxShippingExternalLine(line); err != nil {
		return nil, err
	}

	return line, nil
}

func (a *MoneyTxAggregate) UpdateMoneyTxShippingExternalInfo(ctx context.Context, args *moneytx.UpdateMoneyTxShippingExternalInfoArgs) (*moneytx.MoneyTransactionShippingExternalExtended, error) {
	return nil, nil
}

func (a *MoneyTxAggregate) ConfirmMoneyTxShippingExternals(ctx context.Context, ids []dot.ID) (updated int, _ error) {
	return 0, nil
}

func (a *MoneyTxAggregate) RemoveMoneyTxShippingExternalLines(ctx context.Context, args *moneytx.RemoveMoneyTxShippingExternalLinesArgs) (*moneytx.MoneyTransactionShippingExternalExtended, error) {
	panic("implement me")
}

func (a *MoneyTxAggregate) DeleteMoneyTxShippingExternal(ctx context.Context, ID dot.ID) (deleted int, _ error) {
	panic("implement me")
}
