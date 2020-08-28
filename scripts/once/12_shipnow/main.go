package main

import (
	"o.o/api/main/connectioning"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/accountshipnow/model"
	connectioningaggregate "o.o/backend/com/main/connectioning/aggregate"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	ll  = l.New()
	cfg config.Config
	db  *cmsql.Database
)

func main() {
	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}

	postgres := cfg.Databases.Postgres
	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}

	connectionAggr := connectioningaggregate.NewConnectionAggregate(db, bus.New())
	connectionCommandBus := connectioningaggregate.ConnectionAggregateMessageBus(connectionAggr)

	var fromID dot.ID = 0
	count, created, errCount := 0, 0, 0
	ctx := bus.Ctx()
	for {
		xAccounts, err := scanXAccountAhamoves(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(xAccounts) == 0 {
			ll.S.Infof("Done: created %v/%v", created, count)
			ll.S.Infof("Error: %v/%v", errCount, count)
			break
		}

		fromID = xAccounts[len(xAccounts)-1].ID
		count += len(xAccounts)

		for _, xAccount := range xAccounts {
			cmd := &connectioning.CreateShopConnectionCommand{
				OwnerID:      xAccount.OwnerID,
				ConnectionID: connectioning.DefaultDirectAhamoveConnectionID,
				Token:        xAccount.ExternalToken,
				ExternalData: &connectioning.ShopConnectionExternalData{
					Identifier: xAccount.Phone,
				},
			}
			if err := connectionCommandBus.Dispatch(ctx, cmd); err != nil {
				ll.S.Errorf("Không tạo được shop connection", l.Error(err))
				errCount += 1
				continue
			}
			created += 1
		}
	}

}

func scanXAccountAhamoves(fromID dot.ID) (xAccounts model.ExternalAccountAhamoves, err error) {
	err = db.Table("external_account_ahamove").
		Where("external_token is not null").
		Where("id > ?", fromID.Int64()).
		Limit(1000).
		OrderBy("owner_id, created_at DESC").
		Find(&xAccounts)
	return
}
