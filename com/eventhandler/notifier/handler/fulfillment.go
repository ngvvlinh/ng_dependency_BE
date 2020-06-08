package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping_provider"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	ghtkclient "o.o/backend/pkg/integration/shipping/ghtk/client"
	"o.o/common/l"
)

var acceptNotifyStates = []string{shipping.Returning.String(), shipping.Returned.String(), shipping.Undeliverable.String()}

func HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	var history shipmodel.FulfillmentHistory
	if ok, err := historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	id := history.ID().Int64().Apply(0)
	var ffm shipmodel.Fulfillment
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

func prepareNotiFfmCommands(history shipmodel.FulfillmentHistory, ffm *shipmodel.Fulfillment) []*notifiermodel.CreateNotificationArgs {
	var res []*notifiermodel.CreateNotificationArgs
	externalShippingNote := history.ExternalShippingNote().String()
	externalSubState := history.ExternalShippingSubState().String()
	if (externalShippingNote.Valid && ffm.ExternalShippingNote != "") || (externalSubState.Valid && ffm.ExternalShippingSubState != "") {
		cmd := templateFfmChangedNote(ffm)
		res = append(res, cmd)
	}
	if history.ShippingFeeShop().Int().Valid {
		cmd := templateFfmChangedFee(ffm)
		res = append(res, cmd)
	}
	if history.ShippingState().String().Valid {
		cmd := templateFfmChangedStatus(ffm)
		res = append(res, cmd)
	}
	return res
}

func templateFfmChangedNote(ffm *shipmodel.Fulfillment) *notifiermodel.CreateNotificationArgs {
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
	return &notifiermodel.CreateNotificationArgs{
		AccountID:        ffm.ShopID,
		Title:            title,
		Message:          content,
		Entity:           notifiermodel.NotiFulfillment,
		EntityID:         ffm.ID,
		SendNotification: sendNotification,
	}
}

func templateFfmChangedFee(ffm *shipmodel.Fulfillment) *notifiermodel.CreateNotificationArgs {
	title := fmt.Sprintf("Thay đổi phí vận chuyển - %v %v - %v", Uppercase(ffm.ShippingProvider), ffm.ShippingCode, ffm.AddressTo.FullName)
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)
	content := fmt.Sprintf("Cước phí thay đổi thành %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.ShippingFeeShop, ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	return &notifiermodel.CreateNotificationArgs{
		AccountID:        ffm.ShopID,
		Title:            title,
		Message:          content,
		Entity:           notifiermodel.NotiFulfillment,
		EntityID:         ffm.ID,
		SendNotification: true,
	}
}

func templateFfmChangedStatus(ffm *shipmodel.Fulfillment) *notifiermodel.CreateNotificationArgs {
	content := ""
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)
	switch ffm.ShippingState {
	case shipping.Picking, shipping.Holding:
		content = fmt.Sprintf("Dự kiến giao vào %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", cm.FormatDateVN(ffm.ExpectedDeliveryAt), ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case shipping.Delivering, shipping.Delivered, shipping.Returned:
		content = fmt.Sprintf("Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case shipping.Returning:
		content = fmt.Sprintf("Dự kiến trả hàng trong vòng 3-5 ngày tới. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case shipping.Undeliverable:
		compensationAmount := ffm.ActualCompensationAmount
		if compensationAmount == 0 {
			compensationAmount = ffm.BasketValue
		}
		content = fmt.Sprintf("Giá trị bồi hoàn %vđ. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", cm.FormatCurrency(compensationAmount), ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	case shipping.Cancelled:
		var order = new(ordermodel.Order)
		_ = x.Table("order").Where("id = ?", ffm.OrderID).ShouldGet(order)
		cancelReason := ffm.CancelReason
		if cancelReason == "" {
			cancelReason = order.CancelReason
		}
		content = fmt.Sprintf("Lý do hủy: %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", cancelReason, ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
	default:
	}

	title := fmt.Sprintf("%v - %v %v - %v", shippingsharemodel.ShippingStateMap[ffm.ShippingState], Uppercase(ffm.ShippingProvider), ffm.ShippingCode, ffm.AddressTo.FullName)

	sendNotification := false
	if cm.StringsContain(acceptNotifyStates, ffm.ShippingState.String()) {
		sendNotification = true
	}
	return &notifiermodel.CreateNotificationArgs{
		AccountID:        ffm.ShopID,
		Title:            title,
		Message:          content,
		Entity:           notifiermodel.NotiFulfillment,
		EntityID:         ffm.ID,
		SendNotification: sendNotification,
	}
}

func Uppercase(provider shipping_provider.ShippingProvider) string {
	return strings.ToUpper(provider.String())
}