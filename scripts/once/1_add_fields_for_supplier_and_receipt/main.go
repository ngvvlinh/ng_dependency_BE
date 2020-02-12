package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	receipting "etop.vn/backend/com/main/receipting/model"
	suppliering "etop.vn/backend/com/shopping/suppliering/model"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

type M map[string]interface{}

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cmenv.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	var fromID dot.ID
	var fromIDSupplier dot.ID
	count, updated := 0, 0
	countSupplier, updatedSupplier := 0, 0
	var arrayReceipt []*receipting.Receipt
	var arrSupplers []*suppliering.ShopSupplier
	for {
		receipts, err := scanReceipt(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		suppliers, err := scanSupplier(fromIDSupplier)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(receipts) == 0 && len(suppliers) == 0 {
			break
		}
		fromID = receipts[len(receipts)-1].ID
		fromIDSupplier = suppliers[len(suppliers)-1].ID
		count += len(receipts)
		countSupplier += len(suppliers)
		arrayReceipt = append(arrayReceipt, receipts...)
		arrSupplers = append(arrSupplers, suppliers...)
	}
	for _, receipt := range arrayReceipt {
		update := M{}
		if receipt.Trader != nil {
			update["trader_type"] = receipt.Trader.Type
		}
		traderPhoneNorm := validate.NormalizeSearchPhone(receipt.Trader.Phone)
		if receipt.TraderPhoneNorm == "" || receipt.TraderPhoneNorm != traderPhoneNorm {
			update["trader_phone_norm"] = traderPhoneNorm
		}
		if len(update) > 0 {
			err = db.
				Table("receipt").
				Where("id = ?", receipt.ID).
				ShouldUpdateMap(update)
			if err != nil {
				ll.S.Fatalf("can't update receipt id=%v", receipt.ID)
			}
			updated++
		}
	}
	for _, supplier := range arrSupplers {
		update := M{}
		companyNameNorm := validate.NormalizedSearchToTsVector(validate.NormalizeSearch(supplier.CompanyName))
		if supplier.CompanyNameNorm == "" || supplier.CompanyNameNorm != companyNameNorm {
			update["company_name_norm"] = companyNameNorm
		}
		if len(update) > 0 {
			err = db.
				Table("shop_supplier").
				Where("id = ?", supplier.ID).
				ShouldUpdateMap(update)
			if err != nil {
				ll.S.Fatalf("can't update supplier id=%v", supplier.ID)
			}
			updatedSupplier++
		}
	}
	ll.S.Infof("Updated receipt: %v/%v", updated, count)
	ll.S.Infof("Updated supplier: %v/%v", updatedSupplier, countSupplier)
}

func scanReceipt(fromID dot.ID) (receipts receipting.Receipts, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&receipts)
	return
}
func scanSupplier(fromID dot.ID) (suppliers suppliering.ShopSuppliers, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&suppliers)
	return
}
