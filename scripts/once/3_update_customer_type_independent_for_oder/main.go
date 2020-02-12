package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	ordering "etop.vn/backend/com/main/ordering/model"
	customering "etop.vn/backend/com/shopping/customering/model"
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
	var customerID dot.ID
	var orderID dot.ID
	count, updated := 0, 0
	//var arrOrders []*ordering.Order
	var arrCustomers []*customering.ShopCustomer
	mapShopCustomer := make(map[dot.ID]dot.ID)
	for {
		customers, err := getCustomerTypeIndependent(customerID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(customers) == 0 {
			break
		}
		customerID = customers[len(customers)-1].ID
		arrCustomers = append(arrCustomers, customers...)
	}

	for _, customer := range arrCustomers {
		mapShopCustomer[customer.ShopID] = customer.ID
	}
	for {
		orders, err := scanOrders(orderID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		count += len(orders)
		if len(orders) == 0 {
			break
		}
		orderID = orders[len(orders)-1].ID
		for _, order := range orders {
			update := M{}
			if order.CustomerID != 0 && mapShopCustomer[order.ShopID] == order.CustomerID {
				update["customer_id"] = 1000570234798935240
			}
			if len(update) > 0 {
				err = db.
					Table("order").
					Where("id = ?", order.ID).
					ShouldUpdateMap(update)
				if err != nil {
					ll.S.Fatalf("can't update order id=%v", order.ID)
				}
				updated++
			}
		}
	}
	ll.S.Infof("Updated order: %v/%v", updated, count)
}

func scanOrders(fromID dot.ID) (orders ordering.Orders, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&orders)
	return
}

func getCustomerTypeIndependent(fromID dot.ID) (customers customering.ShopCustomers, err error) {
	err = db.
		Where("type = 'independent' AND id >?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&customers)
	return
}
