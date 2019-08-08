package handler

import (
	"context"
	"fmt"

	txmodel "etop.vn/backend/com/main/moneytx/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/mq"
	etopmodel "etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/notifier/model"
	"etop.vn/backend/pkg/pgevent"
	"etop.vn/common/l"
)

func HandleMoneyTransactionShippingEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	var history txmodel.MoneyTransactionShippingHistory
	if ok, err := x.Where("rid = ?", event.RID).Get(&history); err != nil {
		return mq.CodeStop, err
	} else if !ok {
		ll.Warn("money_transaction_shipping not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}
	id := *history.ID().Int64()
	var mts txmodel.MoneyTransactionShipping
	if ok, err := x.Where("id = ?", id).Get(&mts); err != nil {
		return mq.CodeStop, err
	} else if !ok {
		return mq.CodeIgnore, nil
	}
	cmds := prepareMtsNotiCommands(event, history, mts)
	if err := CreateNotifications(ctx, cmds); err != nil {
		return mq.CodeRetry, err
	}

	return mq.CodeOK, nil
}

func prepareMtsNotiCommands(event *pgevent.PgEvent, history txmodel.MoneyTransactionShippingHistory, mts txmodel.MoneyTransactionShipping) []*model.CreateNotificationArgs {
	var res []*model.CreateNotificationArgs
	if event.Op == pgevent.OpInsert {
		cmd := templateMtsCreated(mts)
		res = append(res, cmd)
	}

	if history.Status().Int() != nil && mts.Status == etopmodel.S3Positive {
		cmd := templateMtsConfirmed(mts)
		res = append(res, cmd)
	}
	return res
}

func templateMtsCreated(mts txmodel.MoneyTransactionShipping) *model.CreateNotificationArgs {
	title := fmt.Sprintf("Phiên đối soát mới %v: %vđ", mts.Code, cm.FormatCurrency(mts.TotalAmount))
	content := fmt.Sprintf("Gồm %v đơn hàng, tổng thu hộ: %vđ, giá trị phiên %vđ.", cm.FormatCurrency(mts.TotalOrders), cm.FormatCurrency(mts.TotalCOD), cm.FormatCurrency(mts.TotalAmount))
	return &model.CreateNotificationArgs{
		AccountID:        mts.ShopID,
		Title:            title,
		Message:          content,
		Entity:           model.NotiMoneyTransactionShipping,
		EntityID:         mts.ID,
		SendNotification: true,
	}
}

func templateMtsConfirmed(mts txmodel.MoneyTransactionShipping) *model.CreateNotificationArgs {
	title := fmt.Sprintf("Đã chuyển tiền phiên %v: %vđ", mts.Code, cm.FormatCurrency(mts.TotalAmount))
	content := fmt.Sprintf("Đã chuyển tiền vào %v. Gồm %v đơn hàng, tổng thu hộ: %vđ, giá trị phiên %vđ.", cm.FormatDateTimeVN(mts.ConfirmedAt), cm.FormatCurrency(mts.TotalOrders), cm.FormatCurrency(mts.TotalCOD), cm.FormatCurrency(mts.TotalAmount))
	return &model.CreateNotificationArgs{
		AccountID:        mts.ShopID,
		Title:            title,
		Message:          content,
		Entity:           model.NotiMoneyTransactionShipping,
		EntityID:         mts.ID,
		SendNotification: true,
	}
}
