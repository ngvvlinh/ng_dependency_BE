package crm

import (
	"context"
	"fmt"
	"time"

	notimodel "etop.vn/backend/com/handler/notifier/model"
	shipmodel "etop.vn/backend/com/main/shipping/model"
	"etop.vn/backend/com/main/shipping/modelx"
	shipmodelx "etop.vn/backend/com/main/shipping/modelx"
	pbcm "etop.vn/backend/pb/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/integration/shipping/ghn"
	ghnclient "etop.vn/backend/pkg/integration/shipping/ghn/client"
	wrapcrm "etop.vn/backend/wrapper/services/crm"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("crm",
		VersionInfo,
		RefreshFulfillmentFromCarrier,
		SendNotification,
	)
}

var ghnCarrier *ghn.Carrier

func Init(ghn *ghn.Carrier) {
	ghnCarrier = ghn
}

func VersionInfo(ctx context.Context, q *wrapcrm.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service:   "etop-crm",
		Version:   "0.1",
		UpdatedAt: nil,
	}
	return nil
}

func RefreshFulfillmentFromCarrier(ctx context.Context, r *wrapcrm.RefreshFulfillmentFromCarrierEndpoint) error {
	query := &shipmodelx.GetFulfillmentQuery{
		ShippingCode: r.ShippingCode,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return err
	}
	ffm := query.Result
	var ffmUpdate *shipmodel.Fulfillment
	var err error
	switch ffm.ShippingProvider {
	case model.TypeGHN:
		ghnCmd := &ghn.RequestGetOrderCommand{
			ServiceID: ffm.ProviderServiceID,
			Request: &ghnclient.OrderCodeRequest{
				OrderCode: ffm.ShippingCode,
			},
			Result: nil,
		}
		if err := ghnCarrier.GetOrder(ctx, ghnCmd); err != nil {
			return err
		}
		ffmUpdate, err = ghnCarrier.CalcRefreshFulfillmentInfo(ctx, ffm, ghnCmd.Result)
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "This feature is not available for this carrier (%v)", ffm.ShippingProvider)
	}

	if err != nil {
		return err
	}
	t0 := time.Now()
	ffmUpdate.LastSyncAt = t0
	update := &modelx.UpdateFulfillmentCommand{
		Fulfillment: ffmUpdate,
	}
	if err := bus.Dispatch(ctx, update); err != nil {
		return err
	}
	r.Result = &pbcm.UpdatedResponse{
		Updated: 1,
	}
	return nil
}

func SendNotification(ctx context.Context, r *wrapcrm.SendNotificationEndpoint) error {
	cmd := &notimodel.CreateNotificationsArgs{
		AccountIDs:       []int64{r.AccountId},
		Title:            r.Title,
		Message:          r.Message,
		EntityID:         r.EntityId,
		Entity:           notimodel.NotiEntity(r.Entity.ToModel()),
		SendNotification: true,
		MetaData:         r.MetaData.GetData(),
	}
	_, _, err := sqlstore.CreateNotifications(ctx, cmd)
	if err != nil {
		return err
	}

	r.Result = pbcm.Message("ok", fmt.Sprintf(
		"Create successful"))
	return nil
}
