package handler

import (
	"context"
	"fmt"

	"o.o/api/top/types/etc/status3"
	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/com/eventhandler/pgevent"
	txmodel "o.o/backend/com/main/moneytx/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
)

func HandleMoneyTransactionShippingEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	var history txmodel.MoneyTransactionShippingHistory
	if ok, err := historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("money_transaction_shipping not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	id := history.ID().Int64().Apply(0)
	var mts txmodel.MoneyTransactionShipping
	if ok, err := x.Where("id = ?", id).Get(&mts); err != nil {
		return mq.CodeStop, err
	} else if !ok {
		return mq.CodeIgnore, nil
	}
	cmds := prepareMtsNotiCommands(ctx, event, history, mts)
	if err := createNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}

	return mq.CodeOK, nil
}

func prepareMtsNotiCommands(ctx context.Context, event *pgevent.PgEvent, history txmodel.MoneyTransactionShippingHistory, mts txmodel.MoneyTransactionShipping) []*notifiermodel.CreateNotificationArgs {
	userIDs, err := filterRecipient(ctx, mts.ShopID, notifyTopicRolesMap[TopicMoneyTransactionShipping])
	if err != nil || userIDs == nil {
		return nil
	}

	var res []*notifiermodel.CreateNotificationArgs
	if event.Op == pgevent.OpInsert {
		cmds := templateMtsCreated(mts, userIDs)
		res = append(res, cmds...)
	}

	if history.Status().Int().Valid && mts.Status == status3.P {
		cmds := templateMtsConfirmed(mts, userIDs)
		res = append(res, cmds...)
	}
	return res
}

func templateMtsCreated(mts txmodel.MoneyTransactionShipping, userIDs []dot.ID) []*notifiermodel.CreateNotificationArgs {
	title := fmt.Sprintf("Phiên đối soát mới %v: %vđ", mts.Code, cm.FormatCurrency(mts.TotalAmount))
	content := fmt.Sprintf("Gồm %v đơn hàng, tổng thu hộ: %vđ, giá trị phiên %vđ.", cm.FormatCurrency(mts.TotalOrders), cm.FormatCurrency(mts.TotalCOD), cm.FormatCurrency(mts.TotalAmount))
	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     mts.ShopID,
		Title:      title,
		Message:    content,
		SendNotify: true,
		Entity:     notifiermodel.NotiMoneyTransactionShipping,
		EntityID:   mts.ID,
		TopicType:  TopicMoneyTransactionShipping,
	}
	return buildNotifyCmds(args)
}

func templateMtsConfirmed(mts txmodel.MoneyTransactionShipping, userIDs []dot.ID) []*notifiermodel.CreateNotificationArgs {
	title := fmt.Sprintf("Đã chuyển tiền phiên %v: %vđ", mts.Code, cm.FormatCurrency(mts.TotalAmount))
	content := fmt.Sprintf("Đã chuyển tiền vào %v. Gồm %v đơn hàng, tổng thu hộ: %vđ, giá trị phiên %vđ.", cm.FormatDateTimeVN(mts.ConfirmedAt), cm.FormatCurrency(mts.TotalOrders), cm.FormatCurrency(mts.TotalCOD), cm.FormatCurrency(mts.TotalAmount))
	args := &buildNotifyCmdArgs{
		UserIDs:    userIDs,
		ShopID:     mts.ShopID,
		Title:      title,
		Message:    content,
		SendNotify: true,
		Entity:     notifiermodel.NotiMoneyTransactionShipping,
		EntityID:   mts.ID,
		TopicType:  TopicMoneyTransactionShipping,
	}
	return buildNotifyCmds(args)
}
