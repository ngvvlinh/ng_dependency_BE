package main

import (
	"flag"

	"o.o/backend/cmd/etop-server/config"
	identitymodel "o.o/backend/com/main/identity/model"
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

	postgres := cfg.Databases.Postgres

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	var fromID dot.ID
	count, updated := 0, 0
	for {
		users, err := scanUser(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(users) == 0 {
			break
		}
		fromID = users[len(users)-1].ID
		count += len(users)
		for _, user := range users {
			userFullNameNorm := validate.NormalizeSearch(user.FullName)
			update := M{}
			if user.FullNameNorm == "" || user.FullNameNorm != userFullNameNorm {
				update["full_name_norm"] = userFullNameNorm
			}
			if len(update) > 0 {
				err = db.
					Table("user").
					Where("id = ?", user.ID).
					ShouldUpdateMap(update)
				if err != nil {
					ll.S.Fatalf("can't update user id=%v", user.ID)
				}
				updated++
			}
		}
	}
	ll.S.Infof("Updated %v/%v", updated, count)
}

func scanUser(fromID dot.ID) (User identitymodel.Users, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).
		OrderBy("id").
		Limit(1000).
		Find(&User)
	return
}
