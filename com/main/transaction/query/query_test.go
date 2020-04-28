package query

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"o.o/api/main/transaction"
	"o.o/api/meta"
	"o.o/api/top/types/etc/status3"
	transactionmodel "o.o/backend/com/main/transaction/model"
	cm "o.o/backend/pkg/common"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	. "o.o/backend/pkg/common/testing"
	"o.o/capi/dot"
)

var (
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
			id INT,
			amount INT,
			account_id INT,
			status INT2,
			type TEXT,
			metadata JSONB,
			note TEXT,
			created_at TIMESTAMP WITH TIME ZONE,
			updated_at TIMESTAMP WITH TIME ZONE
		);
	`)
}

func TestTransactionQueryService(t *testing.T) {
	Convey("Transaction Aggregate", t, func() {
		Reset(func() {
			db.MustExec("truncate transaction")
		})
		_tran := &transactionmodel.Transaction{
			ID:        tranID,
			Amount:    amount,
			AccountID: accountID,
			Status:    status3.P,
			Type:      transaction.TransactionTypeAffiliate,
			Metadata: &transactionmodel.TransactionMetadata{
				ReferralType: transaction.ReferralTypeOrder.String(),
				ReferralIDs:  []dot.ID{555555},
			},
		}
		QS := NewQueryService(db).MessageBus()
		ctx := context.Background()
		_, err := db.Insert(_tran)
		So(err, ShouldBeNil)

		Convey("Get Success", func() {
			query := &transaction.GetTransactionByIDQuery{
				TrxnID:    tranID,
				AccountID: accountID,
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			tran := query.Result
			So(tran.ID, ShouldEqual, tranID)
		})

		Convey("List Success", func() {
			query := &transaction.ListTransactionsQuery{
				AccountID: accountID,
				Paging: meta.Paging{
					Offset: 0,
					Limit:  1,
					Sort:   []string{"-created_at"},
				},
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			trans := query.Result.Transactions
			So(query.Result.Count, ShouldEqual, 1)
			So(trans[0].ID, ShouldEqual, tranID)
		})

		Convey("GetBalance Success", func() {
			query := &transaction.GetBalanceQuery{
				AccountID:       accountID,
				TransactionType: transaction.TransactionTypeAffiliate,
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldBeNil)
			So(query.Result, ShouldEqual, amount)
		})

		Convey("GetBalance Missing AccountID", func() {
			query := &transaction.GetBalanceQuery{
				TransactionType: transaction.TransactionTypeAffiliate,
			}
			err := QS.Dispatch(ctx, query)
			So(err, ShouldCMError, cm.InvalidArgument, "Missing AccountID")
		})
	})
}
