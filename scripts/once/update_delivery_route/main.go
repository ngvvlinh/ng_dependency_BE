package main

import (
	"flag"

	"etop.vn/api/main/location"
	"etop.vn/backend/cmd/etop-server/config"
	servicelocation "etop.vn/backend/com/main/location"
	shipnowmodel "etop.vn/backend/com/main/shipnow/model"
	shipping "etop.vn/backend/com/main/shipping/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmenv"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/sql/cmsql"
	etopmodel "etop.vn/backend/pkg/etop/model"
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
	locationBus := servicelocation.New().MessageBus()
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
	// shipnow fulfillment
	{
		maxGoroutines := 8
		ch := make(chan dot.ID, maxGoroutines)
		chInsert := make(chan error, maxGoroutines)

		var fromID dot.ID
		count, updated := 0, 0
		for {
			shipnowffms, err := scanShipnowFulfillment(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(shipnowffms) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromID = shipnowffms[len(shipnowffms)-1].ID
			count += len(shipnowffms)

			for _, _shipnowffm := range shipnowffms {
				shipnowffm := _shipnowffm
				ch <- shipnowffm.ID
				go func(shopID dot.ID) (_err error) {
					defer func() {
						<-ch
						chInsert <- _err
					}()

					_, err = db.
						Table("shipnow_fulfillment").
						Where("id = ?", shipnowffm.ID).
						UpdateMap(map[string]interface{}{
							"address_to_province_code": shipnowffm.DeliveryPoints[0].ShippingAddress.ProvinceCode,
							"address_to_district_code": shipnowffm.DeliveryPoints[0].ShippingAddress.DistrictCode,
						})
					if err != nil {
						ll.S.Fatalf("can't update shipnow_fulfillment id=%v", shipnowffm.ID)
						return err
					}
					return nil
				}(shipnowffm.ID)
			}
			for i := 0; i < len(shipnowffms); i++ {
				err := <-chInsert
				if err == nil {
					updated++
				} else {
					errs = append(errs, err)
				}
			}
			ll.S.Infof("Updated shipnow fulfillment %v/%v", updated, count)
		}
	}

	// fulfillment
	{
		var fromID dot.ID
		count, updated := 0, 0
		ctx := bus.Ctx()
		for {
			ffms, err := scanFulfillment(fromID)
			if err != nil {
				ll.Fatal("Error", l.Error(err))
			}
			if len(ffms) == 0 {
				ll.S.Infof("Done: updated %v/%v", updated, count)
				break
			}
			fromID = ffms[len(ffms)-1].ID
			count += len(ffms)

			var ids = make(map[string][]dot.ID)

			for _, _ffm := range ffms {
				ffm := _ffm
				deliveryRoute := etopmodel.RouteNationWide
				if ffm.AddressFrom.ProvinceCode == ffm.AddressTo.ProvinceCode {
					deliveryRoute = etopmodel.RouteSameProvince
				} else {
					queryFrom := location.GetLocationQuery{
						ProvinceCode: ffm.AddressFrom.ProvinceCode,
					}
					err = locationBus.Dispatch(ctx, &queryFrom)
					if err != nil {
						_err := cm.Errorf(cm.Internal, err, "fullfilment_id= %v  can't get province from code %v", ffm.ID, queryFrom.ProvinceCode)
						errs = append(errs, _err)
						continue
					}
					queryTo := location.GetLocationQuery{
						ProvinceCode: ffm.AddressTo.ProvinceCode,
					}
					err = locationBus.Dispatch(ctx, &queryTo)
					if err != nil {
						_err := cm.Errorf(cm.Internal, err, "fullfilment_id= %v  can't get province to code %v", ffm.ID, queryFrom.ProvinceCode)
						errs = append(errs, _err)
						continue
					}
					if queryFrom.Result.Province.Region == queryTo.Result.Province.Region {
						deliveryRoute = etopmodel.RouteSameRegion
					}
				}

				ids[string(deliveryRoute)] = append(ids[string(deliveryRoute)], ffm.ID)
			}
			for key, value := range ids {
				_, err = db.
					Table("fulfillment").
					In("id", value).
					UpdateMap(map[string]interface{}{
						"delivery_route": key,
					})
				if err != nil {
					errs = append(errs, err)
				} else {
					updated += len(ids[key])
				}
			}
		}
		for _, err = range errs {
			ll.S.Errorf("Error fulfillment %v", err)
		}
		ll.S.Infof("Updated fulfillment %v/%v", updated, count)
	}
}

func scanFulfillment(fromID dot.ID) (ffm shipping.Fulfillments, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).Where("delivery_route is null").
		OrderBy("id").
		Limit(1000).
		Find(&ffm)
	return ffm, err
}

func scanShipnowFulfillment(fromID dot.ID) (shipnowffm shipnowmodel.ShipnowFulfillments, err error) {
	err = db.
		Where("id > ?", fromID.Int64()).Where("address_to_province_code is null").
		OrderBy("id").
		Limit(1000).
		Find(&shipnowffm)
	return shipnowffm, err
}
