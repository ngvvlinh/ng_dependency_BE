package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/backend/cmd/etop-etl/config"
	"etop.vn/backend/cmd/etop-etl/register"
	"etop.vn/backend/cmd/etop-etl/register/table_name"
	identityquery "etop.vn/backend/com/main/identity"
	identitymodelx "etop.vn/backend/com/main/identity/modelx"
	_ "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel/drivers"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/backend/pkg/common/extservice/telebot"
	"etop.vn/backend/pkg/common/projectpath"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/zexp/etl"
	"etop.vn/capi/dot"
	"etop.vn/common/l"
)

var (
	ll = l.New()
)

type ETLUtil struct {
	bot           *telebot.Channel
	mapDBCfgs     map[string]config.Database
	mapDBs        map[string]*cmsql.Database
	mapTableNames map[string][]table_name.TableName
	identityQuery identity.QueryBus
	resetDB       bool
}

func initDBs(mapDBCfgs map[string]config.Database) map[string]*cmsql.Database {
	mapDBs := make(map[string]*cmsql.Database)
	for wlName, dbCfg := range mapDBCfgs {
		db, err := cmsql.Connect(dbCfg.Postgres)
		if err != nil {
			ll.Fatal("Unable to connect to Postgres", l.Error(err))
		}
		if wlName == drivers.ETop(cmenv.Env()).Key {
			sqlstore.Init(db)

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
	bot *telebot.Channel, resetDB bool,
) *ETLUtil {
	etlUtil := &ETLUtil{
		bot:           bot,
		mapDBCfgs:     mapDBCfgs,
		mapDBs:        initDBs(mapDBCfgs),
		mapTableNames: convertTableNames(mapDBCfgs),
		resetDB:       resetDB,
	}
	etlUtil.identityQuery = identityquery.NewQueryService(etlUtil.mapDBs[drivers.ETop(cmenv.Env()).Key]).MessageBus()
	return etlUtil
}

func (s *ETLUtil) HandleETL(ctx context.Context) {
	for wlKey, _ := range s.mapDBCfgs {
		db := s.mapDBs[wlKey]

		if drivers.ETop(cmenv.Env()).Key == wlKey {
			continue
		}

		if cmenv.IsDev() && s.resetDB {
			_, _ = db.Exec(`
				DO $$ DECLARE
					r RECORD;
				BEGIN
					FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
						EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
					END LOOP;
				END $$;`)
		}

		runAllScriptsDB(db, s.mapTableNames[wlKey])
	}

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
		userIDs, accountIDs, err := s.getUserIDsAndAccountIDs(ctx, driver.ID)
		if err != nil {
			s.bot.SendMessage(fmt.Sprintf("[Error][ETLUtil]: %s", err.Error()))
			continue
		}

		dstDB := s.mapDBs[driver.Key]
		tableNames := s.mapTableNames[driver.Key]

		for _, tableName := range tableNames {
			funcz := register.GetRegisterFuncFromTableName(tableName)
			switch tableName {
			case table_name.User:
				funcz(ng, srcDB, dstDB, userIDs)
			default:
				funcz(ng, srcDB, dstDB, accountIDs)
			}
		}
	}

	return ng
}

func (s *ETLUtil) runEvery5Minutes(ctx context.Context, ch *chan bool) {
	go func(ctx context.Context, ch *chan bool) {
		if len(*ch) != 0 {
			return
		}

		defer func() { <-*ch }()
		*ch <- true
		ng := s.reloadETLEngine(ctx)
		ng.Run()
	}(ctx, ch)
}

func (s *ETLUtil) runEveryNight(ctx context.Context, ch *chan bool) {
	go func(_ctx context.Context, ch *chan bool) {
		if len(*ch) != 0 {
			return
		}
		defer func() { <-*ch }()
		*ch <- true
		ng := s.reloadETLEngine(_ctx)
		ng.RunEveryDay()
	}(ctx, ch)
}

func checkDBExists(db *cmsql.Database, databaseName string) (exists bool, err error) {
	statement := fmt.Sprintf(`SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s');`, databaseName)

	row := db.QueryRow(statement)
	err = row.Scan(&exists)
	return
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
