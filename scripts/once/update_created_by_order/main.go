package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	ordermodel "etop.vn/backend/com/main/ordering/model"
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
	var errs []error
	count, updated := 0, 0

	var orderID dot.ID
	orderID = 1
	for {
		orders, err := scanOrderID(orderID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(orders) == 0 {
			ll.S.Infof("Done: updated order %v/%v", updated, count)
			break
		}
		count += len(orders)
		err = UpdateOrderCreatedBy(orderID, orders[len(orders)-1].ID)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		orderID = orders[len(orders)-1].ID
		updated += len(orders)
		ll.S.Infof("Updated order %v/%v", updated, count)
	}
	for _, err = range errs {
		ll.S.Errorf("Error: %v", err)
	}
}

func scanOrderID(orderID dot.ID) (orders ordermodel.Orders, err error) {
	err = db.
		Where("id > ?", orderID.Int64()).
		Where("order_source_type = 'etop_pos' or order_source_type = 'ts_app'").
		Where("created_by is null").
		OrderBy("id").
		Limit(1000).
		Find(&orders)
	return orders, err
}

func UpdateOrderCreatedBy(idFrom dot.ID, idTo dot.ID) (err error) {
	if _, err = db.SQL(`
		update "order" set
		created_by = shop.owner_id
		from shop
		where
		(
		"order".order_source_type = 'etop_pos' 
			or "order".order_source_type = 'ts_app') and 
		"order".shop_id = shop.id and
		created_by is null and
		"order".id >= ?
		and "order".id <= ?
    `, idFrom, idTo).Exec(); err != nil {
		return err
	}
	return nil
}
