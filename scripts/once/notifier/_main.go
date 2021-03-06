package main

import (
	"context"
	"flag"

	"o.o/backend/cmd/etop-server/config"
	notimodel "o.o/backend/com/eventhandler/notifier/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/common/l"
)

var (
	ll         = l.New()
	cfg        config.Config
	db         *cmsql.Database
	dbNotifier *cmsql.Database
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

	if db, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	} else {
		sqlstore.Init(db)
	}
	if dbNotifier, err = cmsql.Connect(cfg.PostgresNotifier); err != nil {
		ll.Fatal("Error while connecting notifier database", l.Error(err))
	}

	var devices []*notimodel.Device
	if err = dbNotifier.Table("device").Find((*notimodel.Devices)(&devices)); err != nil {
		ll.Fatal("Can not get list devices", l.Error(err))
	}
	maxGoroutines := 8
	ch := make(chan int, maxGoroutines)
	chUpdate := make(chan error, maxGoroutines)
	for i, device := range devices {
		ch <- i
		go func(device *notimodel.Device) (_err error) {
			defer func() {
				<-ch
				chUpdate <- _err
			}()
			cmdUser := &identitymodelx.GetAccountUserQuery{
				AccountID:       device.AccountID,
				FindByAccountID: true,
			}
			if _err = bus.Dispatch(context.Background(), cmdUser); _err != nil {
				ll.Debug("can not get account user", l.Error(err))
				return _err
			}
			return dbNotifier.Table("device").Where("id = ?", device.ID).ShouldUpdateMap(M{"user_id": cmdUser.Result.UserID})
		}(device)
	}
	errCount := 0
	successCount := 0
	for i, n := 0, len(devices); i < n; i++ {
		err := <-chUpdate
		if err != nil {
			errCount++
		} else {
			successCount++
		}
	}
	ll.S.Infof("UpdateInfo device: successs %v/%v, error %v/%v", successCount, len(devices), errCount, len(devices))
}
