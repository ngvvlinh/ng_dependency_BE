package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/lib/pq"

	"o.o/api/main/identity"
	"o.o/backend/cmd/etop-etl/config"
	"o.o/backend/cmd/etop-etl/register"
	"o.o/backend/cmd/etop-etl/register/table_name"
	identityquery "o.o/backend/com/main/identity"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	servicelocation "o.o/backend/com/main/location"
	_ "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/projectpath"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/zexp/etl"
	accountmodel "o.o/backend/zexp/etl/main/account/model"
	accountusermodel "o.o/backend/zexp/etl/main/accountuser/model"
	addressmodel "o.o/backend/zexp/etl/main/address/model"
	customermodel "o.o/backend/zexp/etl/main/customer/model"
	fulfillmentmodel "o.o/backend/zexp/etl/main/fulfillment/model"
	inventoryvariantmodel "o.o/backend/zexp/etl/main/inventoryvariant/model"
	inventoryvouchermodel "o.o/backend/zexp/etl/main/inventoryvoucher/model"
	invitationmodel "o.o/backend/zexp/etl/main/invitation/model"
	moneytransactionshippingmodel "o.o/backend/zexp/etl/main/moneytransactionshipping/model"
	ordermodel "o.o/backend/zexp/etl/main/order/model"
	productshopcollectionmodel "o.o/backend/zexp/etl/main/productshopcollection/model"
	purchaseordermodel "o.o/backend/zexp/etl/main/purchaseorder/model"
	purchaserefundmodel "o.o/backend/zexp/etl/main/purchaserefund/model"
	receiptmodel "o.o/backend/zexp/etl/main/receipt/model"
	refundmodel "o.o/backend/zexp/etl/main/refund/model"
	shipnowfulfillmentmodel "o.o/backend/zexp/etl/main/shipnowfulfillment/model"
	shopmodel "o.o/backend/zexp/etl/main/shop/model"
	shopbrandmodel "o.o/backend/zexp/etl/main/shopbrand/model"
	shopcarriermodel "o.o/backend/zexp/etl/main/shopcarrier/model"
	shopcategorymodel "o.o/backend/zexp/etl/main/shopcategory/model"
	shopcollectionmodel "o.o/backend/zexp/etl/main/shopcollection/model"
	shopcustomergroupmodel "o.o/backend/zexp/etl/main/shopcustomergroup/model"
	shopcustomergroupcustomermodel "o.o/backend/zexp/etl/main/shopcustomergroupcustomer/model"
	shopledgermodel "o.o/backend/zexp/etl/main/shopledger/model"
	shopproductmodel "o.o/backend/zexp/etl/main/shopproduct/model"
	shopproductcollectionmodel "o.o/backend/zexp/etl/main/shopproductcollection/model"
	shopstocktakemodel "o.o/backend/zexp/etl/main/shopstocktake/model"
	shopsuppliermodel "o.o/backend/zexp/etl/main/shopsupplier/model"
	shoptradermodel "o.o/backend/zexp/etl/main/shoptrader/model"
	shoptraderaddressmodel "o.o/backend/zexp/etl/main/shoptraderaddress/model"
	shopvariantmodel "o.o/backend/zexp/etl/main/shopvariant/model"
	shopvariantsuppliermodel "o.o/backend/zexp/etl/main/shopvariantsupplier/model"
	usermodel "o.o/backend/zexp/etl/main/user/model"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll = l.New()
)

type ETLUtil struct {
	mapDBCfgs     map[string]config.Database
	mapDBs        map[string]*cmsql.Database
	mapTableNames map[string][]table_name.TableName
	identityQuery identity.QueryBus
	resetDB       bool
}

type accountUser struct {
	mUserIDs         map[dot.ID]bool
	mAccountIDs      map[dot.ID]bool
	latestUserRID    dot.ID
	latestAccountRID dot.ID
}

var mAccountUser map[string]*accountUser

type migrationFunc func(db *cmsql.Database)

var mTableNameAndModel = map[table_name.TableName]migrationFunc{
	table_name.User:                      func(db *cmsql.Database) { (&usermodel.User{}).Migration(db) },
	table_name.Account:                   func(db *cmsql.Database) { (&accountmodel.Account{}).Migration(db) },
	table_name.ShopCustomer:              func(db *cmsql.Database) { (&customermodel.ShopCustomer{}).Migration(db) },
	table_name.Order:                     func(db *cmsql.Database) { (&ordermodel.Order{}).Migration(db) },
	table_name.Shop:                      func(db *cmsql.Database) { (&shopmodel.Shop{}).Migration(db) },
	table_name.Fulfillment:               func(db *cmsql.Database) { (&fulfillmentmodel.Fulfillment{}).Migration(db) },
	table_name.ShopBrand:                 func(db *cmsql.Database) { (&shopbrandmodel.ShopBrand{}).Migration(db) },
	table_name.ShopProduct:               func(db *cmsql.Database) { (&shopproductmodel.ShopProduct{}).Migration(db) },
	table_name.AccountUser:               func(db *cmsql.Database) { (&accountusermodel.AccountUser{}).Migration(db) },
	table_name.Address:                   func(db *cmsql.Database) { (&addressmodel.Address{}).Migration(db) },
	table_name.InventoryVariant:          func(db *cmsql.Database) { (&inventoryvariantmodel.InventoryVariant{}).Migration(db) },
	table_name.InventoryVoucher:          func(db *cmsql.Database) { (&inventoryvouchermodel.InventoryVoucher{}).Migration(db) },
	table_name.Invitation:                func(db *cmsql.Database) { (&invitationmodel.Invitation{}).Migration(db) },
	table_name.MoneyTransactionShipping:  func(db *cmsql.Database) { (&moneytransactionshippingmodel.MoneyTransactionShipping{}).Migration(db) },
	table_name.ProductShopCollection:     func(db *cmsql.Database) { (&productshopcollectionmodel.ProductShopCollection{}).Migration(db) },
	table_name.PurchaseOrder:             func(db *cmsql.Database) { (&purchaseordermodel.PurchaseOrder{}).Migration(db) },
	table_name.PurchaseRefund:            func(db *cmsql.Database) { (&purchaserefundmodel.PurchaseRefund{}).Migration(db) },
	table_name.Receipt:                   func(db *cmsql.Database) { (&receiptmodel.Receipt{}).Migration(db) },
	table_name.Refund:                    func(db *cmsql.Database) { (&refundmodel.Refund{}).Migration(db) },
	table_name.ShipNowFufillment:         func(db *cmsql.Database) { (&shipnowfulfillmentmodel.ShipnowFulfillment{}).Migration(db) },
	table_name.ShopCarrier:               func(db *cmsql.Database) { (&shopcarriermodel.ShopCarrier{}).Migration(db) },
	table_name.ShopCategory:              func(db *cmsql.Database) { (&shopcategorymodel.ShopCategory{}).Migration(db) },
	table_name.ShopCollection:            func(db *cmsql.Database) { (&shopcollectionmodel.ShopCollection{}).Migration(db) },
	table_name.ShopCustomerGroup:         func(db *cmsql.Database) { (&shopcustomergroupmodel.ShopCustomerGroup{}).Migration(db) },
	table_name.ShopCustomerGroupCustomer: func(db *cmsql.Database) { (&shopcustomergroupcustomermodel.ShopCustomerGroupCustomer{}).Migration(db) },
	table_name.ShopLedger:                func(db *cmsql.Database) { (&shopledgermodel.ShopLedger{}).Migration(db) },
	table_name.ShopProductCollection:     func(db *cmsql.Database) { (&shopproductcollectionmodel.ShopProductCollection{}).Migration(db) },
	table_name.ShopStocktake:             func(db *cmsql.Database) { (&shopstocktakemodel.ShopStocktake{}).Migration(db) },
	table_name.ShopSupplier:              func(db *cmsql.Database) { (&shopsuppliermodel.ShopSupplier{}).Migration(db) },
	table_name.ShopTrader:                func(db *cmsql.Database) { (&shoptradermodel.ShopTrader{}).Migration(db) },
	table_name.ShopTraderAddress:         func(db *cmsql.Database) { (&shoptraderaddressmodel.ShopTraderAddress{}).Migration(db) },
	table_name.ShopVariant:               func(db *cmsql.Database) { (&shopvariantmodel.ShopVariant{}).Migration(db) },
	table_name.ShopVariantSupplier:       func(db *cmsql.Database) { (&shopvariantsuppliermodel.ShopVariantSupplier{}).Migration(db) },
}

func initDBs(mapDBCfgs map[string]config.Database, mapTableNames map[string][]table_name.TableName, resetDB bool) map[string]*cmsql.Database {
	mapDBs := make(map[string]*cmsql.Database)
	mAccountUser = make(map[string]*accountUser)
	for wlName, dbCfg := range mapDBCfgs {

		db, err := cmsql.Connect(dbCfg.Postgres)
		if err != nil {
			ll.Fatal("Unable to connect to Postgres", l.Error(err))
		}

		if wlName == drivers.ETop(cmenv.Env()).Key {
			sqlstore.New(db, nil, servicelocation.QueryMessageBus(servicelocation.New(nil)), nil) // TODO(vu): refactor this
		} else {
			if cmenv.IsDev() && resetDB {
				_, _ = db.Exec(`
				DO $$ DECLARE
					r RECORD;
				BEGIN
					FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
						EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
					END LOOP;
				END $$;`)
			}

			runAllScriptsDB(db, mapTableNames[wlName])

			for _, tableName := range dbCfg.Tables {
				tbName, ok := table_name.ParseTableName(tableName)
				if !ok {
					continue
				}
				mTableNameAndModel[tbName](db)
			}

			mAccountUser[wlName] = &accountUser{
				mUserIDs:    make(map[dot.ID]bool),
				mAccountIDs: make(map[dot.ID]bool),
			}
		}

		mapDBs[wlName] = db
	}
	return mapDBs
}

func convertTableNames(mapDBCfgs map[string]config.Database) map[string][]table_name.TableName {
	mapTableNames := make(map[string][]table_name.TableName)
	for wlKey, dbCfg := range mapDBCfgs {
		mapTableNames[wlKey] = table_name.ConvertStringsToTableNames(dbCfg.Tables)
	}
	return mapTableNames
}

func New(
	mapDBCfgs map[string]config.Database,
	resetDB bool,
) *ETLUtil {
	etlUtil := &ETLUtil{
		mapDBCfgs:     mapDBCfgs,
		mapTableNames: convertTableNames(mapDBCfgs),
		resetDB:       resetDB,
	}
	etlUtil.mapDBs = initDBs(mapDBCfgs, etlUtil.mapTableNames, resetDB)

	etlUtil.identityQuery = identityquery.QueryServiceMessageBus(identityquery.NewQueryService(etlUtil.mapDBs[drivers.ETop(cmenv.Env()).Key]))
	return etlUtil
}

func (s *ETLUtil) HandleETL(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, time.Local)
	if now.After(tomorrow) {
		tomorrow.Add(24 * time.Hour)
	}

	chanEvery5Minutes := make(chan bool, 1)
	chanEveryNight := make(chan bool, 1)

	go s.runEvery5Minutes(ctx, &chanEvery5Minutes)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				go s.runEvery5Minutes(ctx, &chanEvery5Minutes)

				if t.After(tomorrow) && t.Hour() == 2 {
					go s.runEveryNight(ctx, &chanEveryNight)
					tomorrow.Add(24 * time.Hour)
				}
			}
		}
	}()
}

func (s *ETLUtil) reloadETLEngine(ctx context.Context) *etl.ETLEngine {
	ng := etl.NewETLEngine(nil)

	srcDB := s.mapDBs[drivers.ETop(cmenv.Env()).Key]
	for _, driver := range drivers.Drivers(cmenv.Env()) {
		// ignore DB source (etop)
		if driver.ID == 0 {
			continue
		}

		// GET userID and accountIDs by WlPartnerID
		var userIDs, accountIDs []dot.ID
		newUserIDs, latestUserRID, err := scanUser(srcDB, mAccountUser[driver.Key].latestUserRID, driver.ID)
		if err != nil {
			ll.SendMessage(fmt.Sprintf("[Error][ETLUtil]: %s", err.Error()))
			continue
		}
		mAccountUser[driver.Key].latestUserRID = latestUserRID
		for _, userID := range newUserIDs {
			mAccountUser[driver.Key].mUserIDs[userID] = true
		}
		for userID := range mAccountUser[driver.Key].mUserIDs {
			userIDs = append(userIDs, userID)
		}

		newAccountIDs, deletedAccountIDs, latestAccountRID, err := scanShop(srcDB, mAccountUser[driver.Key].latestAccountRID, driver.ID)
		if err != nil {
			fmt.Println(err)
			ll.SendMessage(fmt.Sprintf("[Error][ETLUtil]: %s", err.Error()))
			continue
		}
		mAccountUser[driver.Key].latestAccountRID = latestAccountRID
		for _, accountID := range newAccountIDs {
			mAccountUser[driver.Key].mAccountIDs[accountID] = true
		}
		for _, accountID := range deletedAccountIDs {
			delete(mAccountUser[driver.Key].mAccountIDs, accountID)
		}
		for accountID := range mAccountUser[driver.Key].mAccountIDs {
			accountIDs = append(accountIDs, accountID)
		}

		dstDB := s.mapDBs[driver.Key]
		tableNames := s.mapTableNames[driver.Key]

		for _, tableName := range tableNames {
			funcz := register.GetRegisterFuncFromTableName(tableName)
			switch tableName {
			case table_name.User, table_name.AccountUser:
				funcz(ng, srcDB, dstDB, userIDs)
			default:
				funcz(ng, srcDB, dstDB, accountIDs)
			}
		}
	}

	return ng
}

func (s *ETLUtil) runEvery5Minutes(ctx context.Context, ch *chan bool) {
	go func(_ctx context.Context, _ch *chan bool) {
		if len(*_ch) != 0 {
			return
		}

		defer func() { <-*_ch }()
		*_ch <- true
		ng := s.reloadETLEngine(_ctx)
		ng.Run()
	}(ctx, ch)
}

func (s *ETLUtil) runEveryNight(ctx context.Context, ch *chan bool) {
	go func(_ctx context.Context, _ch *chan bool) {
		if len(*_ch) != 0 {
			return
		}
		defer func() { <-*_ch }()
		*_ch <- true
		ng := s.reloadETLEngine(_ctx)
		ng.RunEveryDay()
	}(ctx, ch)
}

func scanUser(db *cmsql.Database, fromRID, wlPartnerID dot.ID) (ids []dot.ID, latestRID dot.ID, err error) {
	for {
		rows, err := db.Table("user").
			Select("id, rid").
			Where("rid > ?", fromRID.Int64()).
			Where("wl_partner_id = ?", wlPartnerID).
			Limit(1000).
			OrderBy("rid ASC").
			Query()
		if err != nil {
			return nil, 0, err
		}

		var id, rid dot.ID
		count := 0
		for rows.Next() {
			err := rows.Scan(&id, &rid)
			if err != nil {
				return nil, 0, err
			}
			ids = append(ids, id)
			count += 1
			fromRID = rid
		}

		err = rows.Err()
		if err != nil {
			return nil, 0, err
		}

		if count == 0 {
			break
		}
	}
	return ids, fromRID, err
}

func scanShop(db *cmsql.Database, fromRID, wlPartnerID dot.ID) (ids, deletedIDs []dot.ID, latestRID dot.ID, err error) {
	for {
		rows, err := db.Table("shop").
			Select("id, rid, deleted_at").
			Where("rid > ?", fromRID.Int64()).
			Where("wl_partner_id = ?", wlPartnerID).
			OrderBy("rid ASC").
			Limit(1000).
			Query()
		if err != nil {
			return nil, nil, 0, err
		}

		var id, rid dot.ID
		var deletedAt pq.NullTime
		count := 0
		for rows.Next() {
			err := rows.Scan(&id, &rid, &deletedAt)
			if err != nil {
				return nil, nil, 0, err
			}
			if !deletedAt.Valid {
				ids = append(ids, id)
			} else {
				deletedIDs = append(deletedIDs, id)
			}
			count += 1
			fromRID = rid
		}

		err = rows.Err()
		if err != nil {
			return nil, nil, 0, err
		}

		if count == 0 {
			break
		}
	}
	return ids, deletedIDs, fromRID, err
}

func (s *ETLUtil) getUserIDsAndAccountIDs(ctx context.Context, partnerID dot.ID) (userIDs []dot.ID, accountIDs []dot.ID, err error) {
	getUsersByWLPartnerID := &identity.ListUsersByWLPartnerIDQuery{
		ID: partnerID,
	}
	if err := s.identityQuery.Dispatch(ctx, getUsersByWLPartnerID); err != nil {
		return nil, nil, err
	}

	for _, user := range getUsersByWLPartnerID.Result {
		userIDs = append(userIDs, user.ID)
	}

	getAccountsByUserID := &identitymodelx.GetAllAccountUsersQuery{
		UserIDs: userIDs,
	}
	if err := bus.Dispatch(ctx, getAccountsByUserID); err != nil {
		return nil, nil, err
	}
	for _, accountUser := range getAccountsByUserID.Result {
		accountIDs = append(accountIDs, accountUser.AccountID)
	}

	return
}

func runAllScriptsDB(db *cmsql.Database, tables []table_name.TableName) {
	mapTables := make(map[string]bool)

	for _, tableName := range tables {
		mapTables[tableName.Name()] = true
	}

	path := filepath.Join(projectpath.GetPath(), "zexp/etl")

	files, err := ioutil.ReadDir(filepath.Join(path, "db"))
	if err != nil {
		ll.Fatal(err.Error())
	}

	for _, f := range files {
		// Ignore tablenames was not found in config
		tableName := f.Name()[strings.Index(f.Name(), "_")+1 : strings.Index(f.Name(), ".sql")]
		if _, ok := mapTables[tableName]; !ok {
			continue
		}
		dstSQL, err := ioutil.ReadFile(filepath.Join(path, "db/"+f.Name()))
		if err != nil {
			ll.Fatal(err.Error())
		}
		db.MustExec(string(dstSQL))
	}
}
