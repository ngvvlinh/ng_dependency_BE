package update_ghn_order

import (
	"context"
	"flag"
	"fmt"
	"time"

	"etop.vn/backend/cmd/etop-server/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/integration/ghn"
)

var (
	ll  = l.New()
	cfg config.Config
	db  cmsql.Database
	ctx context.Context
)

type M map[string]interface{}

func main() {
	cc.InitFlags()
	flag.Parse()

	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cfg.Postgres.Port = 6432
	cfg.Postgres.Database = "etop_09_01_2018"
	cm.SetEnvironment(cfg.Env)

	if db, err = cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	{
		// ffms, err := scanFulfillments()
		// if err != nil {
		// 	ll.Fatal("Error", l.Error(err))
		// }
		// count, updated := len(ffms), 0
		// for _, ffm := range ffms {
		// 	err := updateEtopDiscount(ffm)
		// 	if err != nil {
		// 		ll.S.Fatalf("Update etop discount fail %v", ffm.ID)
		// 	}
		// 	updated++
		// }
		// ll.S.Infof("Updated %v/%v", updated, count)
	}
	{
		ffms, err := scanFulfillmentExtendeds()
		if err != nil {
			ll.Fatal("Error", l.Error(err))
		}
		fmt.Println("// Total fulfillment :", len(ffms))
		fmt.Printf("[\n")
		for _, ffm := range ffms {
			fmt.Printf(`{
				"ffm_id": %v,
				"ffm_external_shipping_code": "%v",
				"shop_name": "%v",
				"shipping_fee_shop": %v,
				"actual_fee": %v,
				"etop_discount": %v,
				"shipping_state": "%v",
				"created_at": "%v"
			`, ffm.ID, ffm.ExternalShippingCode, ffm.Shop.GetShopName(), ffm.ShippingFeeShop,
				ffm.ActualPrice, ffm.EtopDiscount, ffm.ShippingState, ffm.CreatedAt)

			fmt.Print("},\n")
		}
		fmt.Printf("]\n")
	}
}

func updateEtopDiscount(ffm *model.Fulfillment) error {
	total := 0
	for _, feeLine := range ffm.ProviderShippingFeeLines {
		if feeLine.ShippingFeeType == "main" || feeLine.Cost < 0 {
			total += feeLine.Cost
		}
	}
	discount := total - ffm.ActualPrice
	err := db.
		Table("fulfillment").
		Where("id = ?", ffm.ID).
		ShouldUpdateMap(M{
			"etop_discount": discount,
		})
	return err
}

func scanFulfillments() (ffms model.Fulfillments, err error) {
	err = db.
		Where("external_shipping_code is not null AND actual_price is not null").
		Where("shipping_state not in (?, ?)", model.StateClosed, model.StateCancelled).
		OrderBy("id").
		Find(&ffms)
	return
}

func scanFulfillmentExtendeds() (ffms []*model.FulfillmentExtended, err error) {
	err = db.Table("fulfillment").Where("f.external_shipping_code is not null AND f.actual_price is not null").
		Where("f.shipping_state not in (?, ?)", model.StateClosed, model.StateCancelled).Find((*model.FulfillmentExtendeds)(&ffms))
	return
}

func getInfoAndUpdate() {
	count, updated := 0, 0
	var waitingToFinishList []string
	ffms, err := scanFulfillments()
	if err != nil {
		ll.Fatal("Error", l.Error(err))
	}
	if len(ffms) == 0 {
		ll.S.Infof("Done: updated %v/%v", updated, count)
		return
	}
	count += len(ffms)
	ll.S.Infof("Total ffms: %v", count)
	ctx := context.Background()

	for _, ffm := range ffms {
		// get order GHN info to update service fee
		ghnCmd := &ghn.RequestGetOrderCommand{
			Request: &ghn.OrderCodeRequest{
				OrderCode: ffm.ExternalShippingCode,
			},
		}
		if ghnErr := bus.Dispatch(ctx, ghnCmd); ghnErr != nil {
			ll.Fatal("Error: can not get ghn info", l.Error(ghnErr))
		}
		ghnOrder := ghnCmd.Result

		state := ghn.State(ghnOrder.CurrentStatus)
		if state == ghn.StateWaitingToFinish {
			waitingToFinishList = append(waitingToFinishList, ffm.ExternalShippingCode)
		}

		msgFakeCallback := &ghn.CallbackOrder{
			CoDAmount:            ghnOrder.CoDAmount,
			CurrentStatus:        ghnOrder.CurrentStatus,
			CurrentWarehouseName: ghnOrder.CurrentWarehouseName,
			CustomerID:           ghnOrder.CustomerID,
			CustomerName:         ghnOrder.CustomerName,
			CustomerPhone:        ghnOrder.CustomerName,
			ExternalCode:         ghnOrder.ExternalCode,
			Note:                 ghnOrder.Note,
			OrderCode:            ghnOrder.OrderCode,
			ReturnInfo:           ghnOrder.ReturnInfo,
			ServiceName:          ghnOrder.ServiceName,
			Weight:               ghnOrder.Weight,
		}
		updateFfm := ghn.CalcUpdateFulfillment(ffm, msgFakeCallback, ghnOrder)
		if err := db.Table("fulfillment").
			Where("id = ?", updateFfm.ID).
			ShouldUpdate(updateFfm); err != nil {
			ll.Fatal("Error: Update ffm error", l.Error(err))
		}
		updated++
	}
	ll.S.Infof("Updated %v/%v", updated, count)
	ll.S.Info("================================")
	ll.S.Infof("Please manually recheck ffms (waitingToFinishList) :: %v", waitingToFinishList)
	ll.S.Info("================================")
}

func syncOrderLogsGHN(ffm *model.Fulfillment) (_err error) {
	ctx := context.Background()
	ghnCmd := &ghn.RequestGetOrderLogsCommand{
		Request: &ghn.OrderLogsRequest{
			Condition: &ghn.OrderLogsCondition{
				OrderCode: ffm.ExternalShippingCode,
			},
		},
	}
	if ghnErr := bus.Dispatch(ctx, ghnCmd); ghnErr != nil {
		return ghnErr
	}
	ghnOrderLogs := ghnCmd.Result.Logs
	// order logs already sort by time asc
	lastSyncAt := ffm.LastSyncAt
	var lastTime time.Time
	for _, oLog := range ghnOrderLogs {
		createAt := oLog.CreateTime.ToTime()
		if createAt.After(lastSyncAt) {
			msgFakeCallback := oLog.OrderInfo.ToFakeCallbackOrder()
			ghnOrder := oLog.OrderInfo.ToGHNOrder()
			ffm = ghn.CalcUpdateFulfillment(ffm, msgFakeCallback, ghnOrder)
			lastTime = createAt
		}
	}
	// pp.Println("lastSyncAt ::", lastSyncAt, lastTime, lastTime.IsZero())
	if !lastTime.IsZero() && lastSyncAt.Before(lastTime) {
		ffm.LastSyncAt = lastTime
		if err := db.Table("fulfillment").
			Where("id = ?", ffm.ID).
			ShouldUpdate(ffm); err != nil {
			ll.Fatal("Error: Update ffm error", l.Error(err))
		}
	}
	return nil
}
