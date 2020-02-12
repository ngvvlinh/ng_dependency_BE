package main

import (
	"flag"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/com/main/inventory/model"
	stocktaking "etop.vn/backend/com/main/stocktaking/model"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq/core"
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
	var fromIDInventory dot.ID
	var fromIDStocktake dot.ID
	countStoctake, updatedStocktake, countInventory, updatedInventory := 0, 0, 0, 0
	var arrInventoryVouchers []*model.InventoryVoucher
	var arrStocktakes []*stocktaking.ShopStocktake
	for {
		inventoryVouchers, err := scanInventoryVoucher(fromIDInventory)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}

		if len(inventoryVouchers) == 0 {
			break
		}
		fromIDInventory = inventoryVouchers[len(inventoryVouchers)-1].ID
		countInventory += len(inventoryVouchers)

		arrInventoryVouchers = append(arrInventoryVouchers, inventoryVouchers...)

	}
	for {
		stocktakes, err := scanStocktakes(fromIDStocktake)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(stocktakes) == 0 {
			break
		}
		fromIDStocktake = stocktakes[len(stocktakes)-1].ID
		countStoctake += len(stocktakes)
		arrStocktakes = append(arrStocktakes, stocktakes...)
	}
	for _, inventoryVoucher := range arrInventoryVouchers {
		var productIDsInventory []dot.ID
		for _, value := range inventoryVoucher.Lines {
			if value.ProductID != 0 {
				productIDsInventory = append(productIDsInventory, value.ProductID)
			}
		}
		updateIventory := M{}
		if inventoryVoucher.ProductIDs == nil && len(productIDsInventory) != 0 {
			updateIventory["product_ids"] = core.ArrayScanner(productIDsInventory)
		}
		if len(updateIventory) > 0 {
			err = db.
				Table("inventory_voucher").
				Where("id = ?", inventoryVoucher.ID).
				ShouldUpdateMap(updateIventory)
			if err != nil {
				ll.S.Fatalf("can't update inventory_voucher id=%v", inventoryVoucher.ID)
			}
			updatedInventory++
		}
	}
	for _, stocktake := range arrStocktakes {
		var productIDsStocktake []dot.ID
		for _, value := range stocktake.Lines {
			if value.ProductID != 0 {
				productIDsStocktake = append(productIDsStocktake, value.ProductID)
			}
		}
		updateStocktake := M{}
		if stocktake.ProductIDs == nil && len(productIDsStocktake) != 0 {
			updateStocktake["product_ids"] = core.ArrayScanner(productIDsStocktake)
		}
		if len(updateStocktake) > 0 {
			err = db.
				Table("shop_stocktake").
				Where("id = ?", stocktake.ID).
				ShouldUpdateMap(updateStocktake)
			if err != nil {
				ll.S.Fatalf("can't update stoctake id=%v", stocktake.ID)
			}
			updatedStocktake++
		}
	}
	ll.S.Infof("Updated Inventory: %v/%v", updatedInventory, countInventory)
	ll.S.Infof("Updated Stoctake: %v/%v", updatedStocktake, countStoctake)
}

func scanInventoryVoucher(fromID dot.ID) (vouchers model.InventoryVouchers, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&vouchers)
	return
}

func scanStocktakes(fromID dot.ID) (stocktakes stocktaking.ShopStocktakes, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&stocktakes)
	return
}
