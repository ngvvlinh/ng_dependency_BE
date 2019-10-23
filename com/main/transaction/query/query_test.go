package query

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"etop.vn/api/main/transaction"
	"etop.vn/api/meta"
	transactionmodel "etop.vn/backend/com/main/transaction/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	. "etop.vn/backend/pkg/common/testing"
	"etop.vn/backend/pkg/etop/model"
)

var (
	db        *cmsql.Database
	tranID    = int64(123)
	accountID = int64(123456)
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
			Status:    int(model.S3Positive),
			Type:      string(transaction.TransactionTypeAffiliate),
			Metadata: &transactionmodel.TransactionMetadata{
				ReferralType: string(transaction.ReferralTypeOrder),
				ReferralIDs:  []int64{555555},
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

		Convey("List Succress", func() {
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
