package handler

import (
	"context"
	"fmt"
	"strconv"

	"o.o/api/main/connectioning"
	cmtype "o.o/api/top/types/common"
	"o.o/api/top/types/etc/shipping"
	shippingsubstate "o.o/api/top/types/etc/shipping/substate"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	ordermodel "o.o/backend/com/main/ordering/model"
	shipmodel "o.o/backend/com/main/shipping/model"
	shippingsharemodel "o.o/backend/com/main/shipping/sharemodel"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
)

var acceptNotifyStates = []string{
	shipping.Returning.String(),
	shipping.Returned.String(),
	shipping.Undeliverable.String(),
}

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

	op := event.Op
	cmds := prepareNotifyFfmCommands(ctx, op, history, &ffm)
	if err := createNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}

	return mq.CodeOK, nil
}

func prepareNotifyFfmCommands(ctx context.Context, op pgevent.TGOP, history shipmodel.FulfillmentHistory, ffm *shipmodel.Fulfillment) []*notifiermodel.CreateNotificationArgs {
	userIDs, err := filterRecipient(ctx, ffm.ShopID, notifyTopicRolesMap[TopicFulfillment])
	if err != nil || userIDs == nil {
		return nil
	}

	connection, err := connectionStore(ctx).ID(ffm.ConnectionID).GetConnection()
	if err != nil {
		ll.SendMessagef("get connection got error, %v", err)
		return nil
	}

	var res []*notifiermodel.CreateNotificationArgs
	if history.ShippingFeeShop().Int().Valid {
		cmds := templateFfmChangedFee(connection, op, userIDs, ffm, history)
		if len(cmds) > 0 {
			res = append(res, cmds...)
		}
	}
	if history.ShippingState().String().Valid || history.ShippingSubstate().String().Valid {
		cmds := templateFfmChangedStatus(connection, userIDs, ffm, history)
		if len(cmds) > 0 {
			res = append(res, cmds...)
		}
	}
	return res
}

func templateFfmChangedFee(
	connection *connectioning.Connection,
	op pgevent.TGOP,
	userIDs []dot.ID,
	ffm *shipmodel.Fulfillment,
	history shipmodel.FulfillmentHistory,
) []*notifiermodel.CreateNotificationArgs {
	if op == pgevent.OpInsert {
		return nil
	}

	title := fmt.Sprintf("Thay đổi phí vận chuyển - %v %v - %v", connection.Name, ffm.ShippingCode, ffm.AddressTo.FullName)
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)
	content := fmt.Sprintf("Cước phí thay đổi thành %v. Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", history.ShippingFeeShop().Int64().Int64, ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)

	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     ffm.ShopID,
		Title:      title,
		Message:    content,
		SendNotify: true,
		Entity:     notifiermodel.NotiFulfillment,
		EntityID:   ffm.ID,
		Meta:       cmtype.Empty{},
		TopicType:  TopicFulfillment,
	}
	return buildNotifyCmds(args)
}

func templateFfmChangedStatus(
	connection *connectioning.Connection,
	userIDs []dot.ID,
	ffm *shipmodel.Fulfillment,
	history shipmodel.FulfillmentHistory,
) []*notifiermodel.CreateNotificationArgs {
	shippingState := history.ShippingState().String().String
	shippingSubState := history.ShippingSubstate().String().String
	content := ""
	totalCODAmount := cm.FormatCurrency(ffm.TotalCODAmount)

	// ưu tiên cho substate
	if shippingSubState != "" {
		substate, ok := shippingsubstate.ParseSubstate(shippingSubState)
		if !ok {
			return nil
		}
		title := fmt.Sprintf("%v - %v %v - %v", substate.GetLabelRefName(), connection.Name, ffm.ShippingCode, ffm.AddressTo.FullName)
		if ffm.ExternalShippingNote.Valid {
			content, _ = strconv.Unquote("\"" + ffm.ExternalShippingNote.String + "\"")
		} else {
			content = fmt.Sprintf("Đơn hàng thuộc người nhận %v, %v, %v. Thu hộ %vđ", ffm.AddressTo.FullName, ffm.AddressTo.Phone, ffm.AddressTo.Province, totalCODAmount)
		}

		args := &buildNotifyCmdArgs{
			UserIDs:    userIDs,
			ShopID:     ffm.ShopID,
			Title:      title,
			Message:    content,
			SendNotify: true,
			Entity:     notifiermodel.NotiFulfillment,
			EntityID:   ffm.ID,
			Meta:       cmtype.Empty{},
			TopicType:  TopicFulfillment,
		}
		return buildNotifyCmds(args)
	}

	if shippingState == "" {
		return nil
	}

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

	title := fmt.Sprintf("%v - %v %v - %v", shippingsharemodel.ShippingStateMap[ffm.ShippingState], connection.Name, ffm.ShippingCode, ffm.AddressTo.FullName)

	sendNotification := false
	if cm.StringsContain(acceptNotifyStates, ffm.ShippingState.String()) {
		sendNotification = true
	}
	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     ffm.ShopID,
		Title:      title,
		Message:    content,
		SendNotify: sendNotification,
		Entity:     notifiermodel.NotiFulfillment,
		EntityID:   ffm.ID,
		Meta:       cmtype.Empty{},
		TopicType:  TopicFulfillment,
	}
	return buildNotifyCmds(args)
}
