package aggregate

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"o.o/api/main/transaction"
	"o.o/api/top/types/etc/status3"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	transactionmodel "o.o/backend/com/main/transaction/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	. "o.o/backend/pkg/common/testing"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll        = l.New()
	db        *cmsql.Database
	tranID    = dot.ID(123)
	accountID = dot.ID(123456)
	amount    = 25000
)

func init() {
	postgres := cc.DefaultPostgres()
	db = cmsql.MustConnect(postgres)
	db.MustExec(`
		DROP TABLE IF EXISTS transaction;
		CREATE TABLE transaction (
			id BIGINT,
			amount INT,
			name TEXT,
			account_id INT,
			status INT2,
			type TEXT,
			classify TEXT,
			note TEXT,
			metadata JSONB,
			referral_type TEXT,
			referral_ids INT8[],
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE
		);
	`)
}

func TestTransactionAggregate(t *testing.T) {
	Convey("Aggregate", t, func() {
		_tran := &transactionmodel.Transaction{
			ID:           tranID,
			Amount:       amount,
			AccountID:    accountID,
			Status:       status3.Z,
			ReferralType: subject_referral.Invoice,
			ReferralIDs:  []dot.ID{555555},
			Note:         "note",
		}
		Aggr := AggregateMessageBus(NewAggregate(db))
		ctx := context.Background()
		_, err := db.Insert(_tran)
		So(err, ShouldBeNil)

		Reset(func() {
			db.MustExec("truncate transaction")
		})

		Convey("Create Success", func() {
			cmd := &transaction.CreateTransactionCommand{
				Amount:       amount,
				AccountID:    accountID,
				Type:         transaction_type.Credit,
				Note:         "123456",
				ReferralType: subject_referral.Invoice,
				ReferralIDs:  []dot.ID{555555},
			}
			err = Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			tran := cmd.Result
			So(tran.AccountID, ShouldEqual, accountID)
		})

		Convey("Confirm Missing Transaction ID", func() {
			cmd := &transaction.ConfirmTransactionCommand{}
			err = Aggr.Dispatch(ctx, cmd)
			So(err, ShouldCMError, cm.InvalidArgument, "Missing TransactionID")
		})

		Convey("Confirm Missing Account ID", func() {
			cmd := &transaction.ConfirmTransactionCommand{
				TrxnID: tranID,
			}
			err = Aggr.Dispatch(ctx, cmd)
			So(err, ShouldCMError, cm.InvalidArgument, "Missing AccountID")
		})

		Convey("Confirm Success", func() {
			cmd := &transaction.ConfirmTransactionCommand{
				TrxnID:    tranID,
				AccountID: accountID,
			}
			err = Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			So(cmd.Result.Status, ShouldEqual, status3.P)

			Convey("Confirm Fail Precondition: Status = 1", func() {
				cmd = &transaction.ConfirmTransactionCommand{
					TrxnID:    tranID,
					AccountID: accountID,
				}
				err = Aggr.Dispatch(ctx, cmd)
				So(err, ShouldCMError, cm.FailedPrecondition, "Can not confirm this transaction")
			})
		})

		Convey("Cancel Success", func() {
			cmd := &transaction.CancelTransactionCommand{
				TrxnID:    tranID,
				AccountID: accountID,
			}
			err = Aggr.Dispatch(ctx, cmd)
			So(err, ShouldBeNil)
			So(cmd.Result.Status, ShouldEqual, status3.N)

			Convey("Cancel Fail Precondition: Status = -1", func() {
				cmd = &transaction.CancelTransactionCommand{
					TrxnID:    tranID,
					AccountID: accountID,
				}
				err = Aggr.Dispatch(ctx, cmd)
				So(err, ShouldCMError, cm.FailedPrecondition, "Can not cancel this transaction")
			})
		})
	})
}
