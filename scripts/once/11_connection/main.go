package main

import (
	"o.o/api/main/shipping"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/shipping/model"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
	"o.o/common/l"
)

var (
	cfg config.Config
	ll  = l.New()
	db  *cmsql.Database
)

func main() {
	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	postgres := cfg.Databases.Postgres
	postgres.Port = 25434
	postgres.Username = "etop"
	postgres.Database = "etopv1_2"
	postgres.Password = "ZMp42845Rd5Lwbvj"

	if db, err = cmsql.Connect(postgres); err != nil {
		ll.Fatal("Error while connection database", l.Error(err))
	}

	var fromID dot.ID = 0
	count, updated, errCount := 0, 0, 0
	for {
		ffms, err := scanFulfillments(fromID)
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		if len(ffms) == 0 {
			ll.S.Infof("Done: updated %v/%v", updated, count)
			ll.S.Infof("Error %v/%v", errCount, count)
			break
		}

		fromID = ffms[len(ffms)-1].ID
		count += len(ffms)

		ffmIDsByCarrier := make(map[shipping_provider.ShippingProvider][]dot.ID)
		for _, ffm := range ffms {
			ffmIDsByCarrier[ffm.ShippingProvider] = append(ffmIDsByCarrier[ffm.ShippingProvider], ffm.ID)
		}

		for carrier, ids := range ffmIDsByCarrier {
			connectionID := shipping.GetConnectionID(0, carrier)
			if connectionID == 0 {
				continue
			}
			update := map[string]interface{}{
				"connection_id":     connectionID,
				"connection_method": connection_type.ConnectionMethodBuiltin,
			}
			if err := db.Table("fulfillment").
				In("id", ids).
				ShouldUpdateMap(update); err != nil {
				errCount += len(ids)
			}
			updated += len(ids)
		}
	}
}

func scanFulfillments(fromID dot.ID) (ffms model.Fulfillments, err error) {
	err = db.Where("id > ?", fromID.String()).
		Where("connection_id IS NULL").
		OrderBy("id").
		Limit(1000).
		Find(&ffms)
	return
}
