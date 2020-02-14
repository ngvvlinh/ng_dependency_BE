package main

import (
	"flag"

	"etop.vn/api/shopping/customering"
	"etop.vn/backend/cmd/etop-server/config"
	modelcustomering "etop.vn/backend/com/shopping/customering/model"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

type M map[string]interface{}

// var independentCustomers []dot.ID

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

	if err != nil {
		ll.S.Fatalf("Update receipt fail", err)
	}
	var fromCustomerID dot.ID
	var countOrder int
	var countReceipt int
	var countRefund int
	var countInventoryVoucher int
	for {
		customers, err := scanIndependentCustomer(fromCustomerID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(customers) == 0 {
			break
		}
		fromCustomerID = customers[len(customers)-1].ID
		var independentCustomersIDs []dot.ID
		for _, v := range customers {
			independentCustomersIDs = append(independentCustomersIDs, v.ID)
		}
		//func update receipt
		count, err := updateReceipts(independentCustomersIDs)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		countReceipt += count
		//func update order
		count, err = updateOrders(independentCustomersIDs)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		countOrder += count
		//func update refund
		count, err = updateRefunds(independentCustomersIDs)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		countRefund += count
		//func update inventory_voucher
		count, err = updateInventoryVouchers(independentCustomersIDs)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		countInventoryVoucher += count
	}
	ll.Info("Update order", l.Int("count", countOrder))
	ll.Info("Update receipt : %v", l.Int("count", countReceipt))
	ll.Info("Update refund : %v", l.Int("count", countRefund))
	ll.Info("Update inventory_voucher : %v", l.Int("count", countInventoryVoucher))
}

func updateReceipts(independentCustomersIDs []dot.ID) (int, error) {
	var update = M{}
	update["trader_id"] = customering.CustomerAnonymous
	count, err := db.
		Table("receipt").
		In("trader_id", independentCustomersIDs).
		UpdateMap(update)
	return count, err
}

func updateOrders(independentCustomersIDs []dot.ID) (int, error) {
	var update = M{}
	update["customer_id"] = customering.CustomerAnonymous
	count, err := db.
		Table("order").
		In("customer_id", independentCustomersIDs).
		UpdateMap(update)
	return count, err
}

func updateRefunds(independentCustomersIDs []dot.ID) (int, error) {
	var update = M{}
	update["customer_id"] = customering.CustomerAnonymous
	count, err := db.
		Table("refund").
		In("customer_id", independentCustomersIDs).
		UpdateMap(update)
	return count, err
}

func updateInventoryVouchers(independentCustomersIDs []dot.ID) (int, error) {
	var update = M{}
	update["trader_id"] = customering.CustomerAnonymous
	count, err := db.
		Table("inventory_voucher").
		In("trader_id", independentCustomersIDs).
		UpdateMap(update)
	return count, err
}

func scanIndependentCustomer(fromID dot.ID) (customers modelcustomering.ShopCustomers, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		Where("type = ?", "independent").
		OrderBy("id").
		Limit(1000).
		Find(&customers)
	return
}
