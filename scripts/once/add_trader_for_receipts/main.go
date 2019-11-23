package main

import (
	"encoding/json"
	"flag"

	"etop.vn/api/shopping/tradering"
	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/com/main/receipting/model"
	carriermodel "etop.vn/backend/com/shopping/carrying/model"
	customermodel "etop.vn/backend/com/shopping/customering/model"
	suppliermodel "etop.vn/backend/com/shopping/suppliering/model"
	tradermodel "etop.vn/backend/com/shopping/tradering/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cm.SetEnvironment(cfg.Env)

	postgres := cfg.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	mapReceipt := make(map[int64]*model.Receipt)
	{
		fromID := int64(0)

		for {
			receipts, err := scanReceipts(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(receipts) == 0 {
				break
			}
			for _, receipt := range receipts {
				if receipt.Trader == nil {
					mapReceipt[receipt.ID] = receipt
				}
			}
			fromID = receipts[len(receipts)-1].ID
		}
	}

	{
		count, errCount, updatedCount := len(mapReceipt), 0, 0
		maxGoroutines := 8
		ch := make(chan int64, maxGoroutines)
		chInsert := make(chan error, maxGoroutines)
		for receiptID, receipt := range mapReceipt {
			ch <- receiptID
			go func(id, traderID int64) (_err error) {
				defer func() {
					<-ch
					chInsert <- _err
				}()
				var trader *model.Trader
				trader, _ = getTrader(traderID)
				if _err != nil {
					return _err
				}
				_ = updateTraderForReceipt(id, trader)
				return _err
			}(receiptID, receipt.TraderID)
		}
		for i, n := 0, len(mapReceipt); i < n; i++ {
			err := <-chInsert
			if err != nil {
				errCount++
			} else {
				updatedCount++
			}
		}
		ll.S.Infof("Update trader for receipts: success %v/%v, error %v/%v", updatedCount, count, errCount, count)
	}
}

func scanReceipts(fromID int64) (receipts model.Receipts, err error) {
	err = db.
		Where("id > ?", fromID).
		Where("deleted_at is null").
		OrderBy("id").
		Limit(1000).
		Find(&receipts)
	return
}

func getTrader(traderID int64) (trader *model.Trader, err error) {
	var tradersModel tradermodel.ShopTraders
	err = db.
		Where("id = ?", traderID).
		Find(&tradersModel)
	if err != nil {
		return nil, err
	}

	if len(tradersModel) == 0 {
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "not found")
	}
	traderModel := tradersModel[0]

	switch traderModel.Type {
	case tradering.CustomerType:
		var customersModel customermodel.ShopCustomers
		err = db.
			Where("id = ?", traderID).
			Limit(1).
			Find(&customersModel)
		if err != nil {
			return nil, err
		}

		if len(customersModel) == 0 {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "not found")
		}
		customerModel := customersModel[0]
		trader = &model.Trader{
			ID:       traderID,
			Type:     tradering.CustomerType,
			FullName: customerModel.FullName,
			Phone:    customerModel.Phone,
		}
	case tradering.CarrierType:
		var carriersModel carriermodel.ShopCarriers
		err = db.
			Where("id = ?", traderID).
			Limit(1).
			Find(&carriersModel)
		if err != nil {
			return nil, err
		}

		if len(carriersModel) == 0 {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "not found")
		}
		carrierModel := carriersModel[0]
		trader = &model.Trader{
			ID:       traderID,
			Type:     tradering.CarrierType,
			FullName: carrierModel.FullName,
		}
	case tradering.SupplierType:
		var suppliersModel suppliermodel.ShopSuppliers
		err = db.
			Where("id = ?", traderID).
			Limit(1).
			Find(&suppliersModel)
		supplierModel := suppliersModel[0]
		if err != nil {
			return nil, err
		}

		if len(suppliersModel) == 0 {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "not found")
		}
		trader = &model.Trader{
			ID:       traderID,
			Type:     tradering.SupplierType,
			FullName: supplierModel.FullName,
		}
	}

	return
}

func updateTraderForReceipt(ID int64, trader *model.Trader) (err error) {
	traderJSON, err := json.Marshal(trader)
	if err != nil {
		return err
	}
	_, err = db.
		Table("receipt").
		Where("id=?", ID).
		UpdateMap(map[string]interface{}{
			"trader": string(traderJSON),
		})
	return err
}
