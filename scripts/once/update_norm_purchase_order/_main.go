package main

import (
	"flag"

	"o.o/backend/cmd/etop-server/config"
	purchaseorder "o.o/backend/com/main/purchaseorder/model"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
	"o.o/common/l"
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
	count, updated := 0, 0
	var arrayPurchaseOrders []*purchaseorder.PurchaseOrder
	for {
		purchaseOrders, err := scanPurchaseOrder(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(purchaseOrders) == 0 {
			break
		}
		fromID = purchaseOrders[len(purchaseOrders)-1].ID
		count += len(purchaseOrders)
		arrayPurchaseOrders = append(arrayPurchaseOrders, purchaseOrders...)
	}
	for _, purchaseOrder := range arrayPurchaseOrders {
		supplierFullNameNorm := validate.NormalizeSearch(purchaseOrder.Supplier.FullName)
		update := M{}

		if purchaseOrder.SupplierFullNameNorm == "" || purchaseOrder.SupplierFullNameNorm != supplierFullNameNorm {
			update["supplier_full_name_norm"] = supplierFullNameNorm
		}
		supplierPhoneNorm := validate.NormalizeSearchPhone(purchaseOrder.Supplier.Phone)
		if purchaseOrder.SupplierPhoneNorm == "" || purchaseOrder.SupplierFullNameNorm != supplierPhoneNorm {
			update["supplier_phone_norm"] = supplierPhoneNorm
		}
		if len(update) > 0 {
			err = db.
				Table("purchase_order").
				Where("id = ?", purchaseOrder.ID).
				ShouldUpdateMap(update)
			if err != nil {
				ll.S.Fatalf("can't update purchase_order id=%v", purchaseOrder.ID)
			}
			updated++
		}
	}
	ll.S.Infof("Updated %v/%v", updated, count)
}

func scanPurchaseOrder(fromID dot.ID) (purchaseOrders purchaseorder.PurchaseOrders, err error) {
	err = db.
		Where("id > ?", fromID).
		OrderBy("id").
		Limit(1000).
		Find(&purchaseOrders)
	return
}
