package main

import (
	"flag"
	"sync"

	"o.o/api/top/types/etc/service_classify"
	"o.o/api/top/types/etc/subject_referral"
	"o.o/api/top/types/etc/transaction_type"
	"o.o/backend/cmd/etop-server/config"
	creditmodel "o.o/backend/com/main/credit/model"
	transactionaggr "o.o/backend/com/main/transaction/aggregate"
	transactionmodel "o.o/backend/com/main/transaction/model"
	"o.o/backend/pkg/common/bus"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	cfg config.Config
	db  *cmsql.Database
	ll  = l.New()
	mu  sync.Mutex
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	postgres := cfg.Databases.Postgres
	postgres.Host = "localhost"
	postgres.Port = 5432
	postgres.Username = "postgres"
	postgres.Password = "postgres"
	postgres.Database = "etop_dev"
	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	trxnAggr := transactionaggr.NewAggregate(db)

	var fromID = dot.ID(0)
	var wg sync.WaitGroup
	var count, updatedCount, errorCount = 0, 0, 0
	var ctx = bus.Ctx()

	for {
		credits, _err := scanCredits(fromID)
		if _err != nil {
			ll.Fatal("Fail to scan credit")
		}
		if len(credits) == 0 {
			break
		}
		fromID = credits[len(credits)-1].ID
		count += len(credits)
		for _, credit := range credits {
			go func(c *creditmodel.Credit) {
				wg.Add(1)
				// create transaction
				trxn := &transactionmodel.Transaction{
					Name:         "",
					ID:           c.ID,
					Amount:       c.Amount,
					AccountID:    c.ShopID,
					Status:       c.Status,
					Type:         transaction_type.Credit,
					Classify:     service_classify.ServiceClassify(c.Classify),
					Note:         "",
					ReferralType: subject_referral.Credit,
					ReferralIDs:  []dot.ID{c.ID},
					CreatedAt:    c.PaidAt,
					UpdatedAt:    c.UpdatedAt,
				}
				mu.Lock()
				defer func() {
					wg.Done()
					mu.Unlock()
				}()
				if !c.DeletedAt.IsZero() {
					return
				}
				if err = trxnAggr.ForceCreateTransaction(ctx, trxn); err != nil {
					errorCount++
					ll.Error("Create transaction error", l.Error(err))
				} else {
					updatedCount++
				}
			}(credit)
		}
	}
	wg.Wait()
	ll.Info("Create transaction finished")
	ll.S.Infof("Created: %v/%v", updatedCount, count)
	ll.S.Infof("Error: %v/%v", errorCount, count)
}

func scanCredits(fromID dot.ID) (credits creditmodel.Credits, err error) {
	err = db.Where("id > ?", fromID.String()).
		Where("deleted_at IS NULL").
		OrderBy("id").
		Limit(1000).
		Find(&credits)
	return
}
