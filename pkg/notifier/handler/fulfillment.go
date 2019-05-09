package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	ghtkclient "etop.vn/backend/pkg/integration/ghtk/client"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/mq"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/notifier/model"
	"etop.vn/backend/pkg/pgevent"
)

var acceptNotifyStates = []string{string(etopmodel.StateReturning), string(etopmodel.StateReturned), string(etopmodel.StateUndeliverable)}

func HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	var history etopmodel.FulfillmentHistory
	if ok, err := x.Where("rid = ?", event.RID).Get(&history); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	id := *history.ID().Int64()
	var ffm etopmodel.Fulfillment
	if ok, err := x.Where("id = ?", id).Get(&ffm); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID), l.Int64("id", id))
		return mq.CodeIgnore, nil
	}

	cmds := prepareNotiFfmCommands(history, &ffm)
	if err := CreateNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}

	return mq.CodeOK, nil
}

func prepareNotiFfmCommands(history etopmodel.FulfillmentHistory, ffm *etopmodel.Fulfillment) []*model.CreateNotificationArgs {
	var res []*model.CreateNotificationArgs
	externalShippingNote := history.ExternalShippingNote().String()
	externalSubState := history.ExternalShippingSubState().String()
	if (externalShippingNote != nil && ffm.ExternalShippingNote != "") || (externalSubState != nil && ffm.ExternalShippingSubState != "") {
		cmd := templateFfmChangedNote(ffm)
		res = append(res, cmd)
	}
	if history.ShippingFeeShop().Int32() != nil {
		cmd := templateFfmChangedFee(ffm)
		res = append(res, cmd)
	}
	if history.ShippingState().String() != nil {
		cmd := templateFfmChangedStatus(ffm)
		res = append(res, cmd)
	}
	return res
}

func templateFfmChangedNote(ffm *etopmodel.Fulfillment) *model.CreateNotificationArgs {
	title, content := "", ""
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)
	subState := ffm.ExternalShippingSubState
	if subState == "" {
		subState = "Cập nhật"
	}
	title = fmt.Sprintf("%v - %v %v - %v", subState, Uppercase(ffm.ShippingProvider), ffm.ShippingCode, ffm.AddressTo.FullName)
	if ffm.ExternalShippingNote != "" {
		content, _ = strconv.Unquote("\"" + ffm.ExternalShippingNote + "\"")
	} else {
		content = fmt.Sprintf("Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	}

	sendNotification := true
	if subState == ghtkclient.SubStateMapping[ghtkclient.StateIDShipperDelivered] {
		sendNotification = false
	}
	return &model.CreateNotificationArgs{
		AccountID:        ffm.ShopID,
		Title:            title,
		Message:          content,
		Entity:           model.NotiFulfillment,
		EntityID:         ffm.ID,
		SendNotification: sendNotification,
	}
}

func templateFfmChangedFee(ffm *etopmodel.Fulfillment) *model.CreateNotificationArgs {
	title := fmt.Sprintf("Thay đổi phí vận chuyển - %v %v - %v", Uppercase(ffm.ShippingProvider), ffm.ShippingCode, ffm.AddressTo.FullName)
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)
	content := fmt.Sprintf("Cước phí thay đổi thành %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.ShippingFeeShop, ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	return &model.CreateNotificationArgs{
		AccountID:        ffm.ShopID,
		Title:            title,
		Message:          content,
		Entity:           model.NotiFulfillment,
		EntityID:         ffm.ID,
		SendNotification: true,
	}
}

func templateFfmChangedStatus(ffm *etopmodel.Fulfillment) *model.CreateNotificationArgs {
	content := ""
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)
	switch ffm.ShippingState {
	case etopmodel.StatePicking, etopmodel.StateHolding:
		content = fmt.Sprintf("Dự kiến giao vào %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", cm.FormatDateVN(ffm.ExpectedDeliveryAt), ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case etopmodel.StateDelivering, etopmodel.StateDelivered, etopmodel.StateReturned:
		content = fmt.Sprintf("Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case etopmodel.StateReturning:
		content = fmt.Sprintf("Dự kiến trả hàng trong vòng 3-5 ngày tới. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case etopmodel.StateUndeliverable:
		compensationAmount := ffm.ActualCompensationAmount
		if compensationAmount == 0 {
			compensationAmount = ffm.BasketValue
		}
		content = fmt.Sprintf("Giá trị bồi hoàn %vđ. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", cm.FormatCurrency(compensationAmount), ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case etopmodel.StateCancelled:
		var order = new(etopmodel.Order)
		_ = x.Table("order").Where("id = ?", ffm.OrderID).ShouldGet(order)
		cancelReason := ffm.CancelReason
		if cancelReason == "" && order != nil {
			cancelReason = order.CancelReason
		}
		content = fmt.Sprintf("Lý do hủy: %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", cancelReason, ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	default:
	}

	title := fmt.Sprintf("%v - %v %v - %v", etopmodel.ShippingStateMap[ffm.ShippingState], Uppercase(ffm.ShippingProvider), ffm.ShippingCode, ffm.AddressTo.FullName)

	sendNotification := false
	if cm.StringsContain(acceptNotifyStates, string(ffm.ShippingState)) {
		sendNotification = true
	}
	return &model.CreateNotificationArgs{
		AccountID:        ffm.ShopID,
		Title:            title,
		Message:          content,
		Entity:           model.NotiFulfillment,
		EntityID:         ffm.ID,
		SendNotification: sendNotification,
	}
}

func Uppercase(provider etopmodel.ShippingProvider) string {
	return strings.ToUpper(string(provider))
}
