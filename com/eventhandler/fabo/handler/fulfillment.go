package handler

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"o.o/api/fabo/fbmessaging"
	"o.o/api/fabo/fbmessaging/fb_customer_conversation_type"
	"o.o/api/fabo/fbpaging"
	"o.o/api/fabo/fbusering"
	"o.o/api/main/identity"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	shipping_state "o.o/api/top/types/etc/shipping"
	"o.o/api/top/types/etc/shipping/substate"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/com/fabo/pkg/fbclient"
	"o.o/backend/com/fabo/pkg/fbclient/model"
	fulfillmentmodel "o.o/backend/com/main/shipping/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/mq"
	"o.o/capi/dot"
	"o.o/common/l"
)

var validStates = []shipping_state.State{shipping_state.Created, shipping_state.Picking, shipping_state.Holding, shipping_state.Delivering, shipping_state.Delivered, shipping_state.Returning}

func (h *Handler) HandleFulfillmentEvent(ctx context.Context, event *pgevent.PgEvent) (mq.Code, error) {
	ll.Info("HandleFulfillmentEvent", l.Object("pgevent", event))
	var history fulfillmentmodel.FulfillmentHistory
	if ok, err := h.historyStore(ctx).GetHistory(&history, event.RID); err != nil {
		return mq.CodeStop, nil
	} else if !ok {
		ll.Warn("Fulfillment not found", l.Int64("rid", event.RID))
		return mq.CodeIgnore, nil
	}

	id := history.ID().ID().Apply(0)
	state := history.ShippingState().String().String
	substate := history.ShippingSubstate().String().String

	getFFmQuery := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
		ID: id,
	}
	if err := h.shippingQuery.Dispatch(ctx, getFFmQuery); err != nil {
		ll.Warn("fulfillment not found", l.Int64("rid", event.RID), l.ID("id", id))
		return mq.CodeIgnore, nil
	}
	ffm := getFFmQuery.Result

	code, err := h.sendMessageWhenChangeShippingStateAndSubstate(ctx, event, ffm, state, substate)
	if err != nil || (code != mq.CodeOK && code != mq.CodeIgnore) {
		return code, err
	}

	return mq.CodeOK, nil
}

func (h *Handler) sendMessageWhenChangeShippingStateAndSubstate(ctx context.Context, event *pgevent.PgEvent, ffm *shipping.Fulfillment, historyState, historySubstate string) (mq.Code, error) {
	// validate historyState and historySubstate
	if historyState == "" && historySubstate == "" {
		return mq.CodeIgnore, nil
	}

	if historyState != "" {
		var check bool
		for _, validState := range validStates {
			if validState.String() == historyState {
				check = true
				break
			}
		}
		if !check {
			return mq.CodeIgnore, nil
		}
	}

	// Get order by order_id into ffm
	getOrderQuery := &ordering.GetOrderByIDQuery{
		ID: ffm.OrderID,
	}
	if err := h.orderQuery.Dispatch(ctx, getOrderQuery); err != nil {
		ll.Warn("order not found", l.Int64("rid", event.RID), l.ID("id", ffm.OrderID))
		return mq.CodeIgnore, nil
	}
	order := getOrderQuery.Result
	if order.CustomerID == 0 {
		return mq.CodeIgnore, nil
	}

	// Get shop
	getShopQuery := &identity.GetShopByIDQuery{
		ID: ffm.ShopID,
	}
	if err := h.indentityQuery.Dispatch(ctx, getShopQuery); err != nil {
		ll.Warn("shop not found", l.Int64("rid", event.RID), l.ID("id", ffm.ShopID))
	}
	shop := getShopQuery.Result

	// list FbUsers by customerID
	listFbExternalUserIDsByCustomerIDQuery := &fbusering.ListFbExternalUserIDsByShopCustomerIDsQuery{
		CustomerIDs: []dot.ID{order.CustomerID},
	}
	if err := h.fbuserQuery.Dispatch(ctx, listFbExternalUserIDsByCustomerIDQuery); err != nil {
		ll.Warn("ListFbExternalUserIDs error", l.String("err", err.Error()))
		return mq.CodeIgnore, nil
	}
	fbExternalUserIDs := listFbExternalUserIDsByCustomerIDQuery.Result
	if len(fbExternalUserIDs) == 0 {
		return mq.CodeIgnore, nil
	}

	// list conversations by fb_user_ids (above)
	listFbCustomerConversationsQuery := &fbmessaging.ListFbCustomerConversationsByExternalUserIDsQuery{
		ExtUserIDs:       fbExternalUserIDs,
		ConversationType: fb_customer_conversation_type.Message.Wrap(),
	}
	if err := h.fbMessagingQuery.Dispatch(ctx, listFbCustomerConversationsQuery); err != nil {
		ll.Warn("listFbCustomerConversations error", l.String("err", err.Error()))
		return mq.CodeIgnore, nil
	}

	// generate message depends on historySubstate and historyState
	message, err := h.generateMessage(ffm, shop, historyState, historySubstate)
	if err != nil {
		ll.Warn("generateMessage error", l.String("err", err.Error()))
		return mq.CodeIgnore, nil
	}

	// Chỉ gửi tin nhắn fbUser đã có customerConversation trên hệ thống
	for _, customerConversation := range listFbCustomerConversationsQuery.Result {
		getFbExternalPageInternalQuery := &fbpaging.GetFbExternalPageInternalByExternalIDQuery{
			ExternalID: customerConversation.ExternalPageID,
		}
		if err := h.fbPagingQuery.Dispatch(ctx, getFbExternalPageInternalQuery); err != nil {
			ll.Warn("FbExternalPageInternal not found", l.String("err", err.Error()), l.String("external_page_id", customerConversation.ExternalPageID))
			return mq.CodeIgnore, nil
		}
		accessToken := getFbExternalPageInternalQuery.Result.Token

		// ignore error when send message
		_, _ = h.fbClient.CallAPISendMessage(&fbclient.SendMessageRequest{
			AccessToken: accessToken,
			SendMessageArgs: &model.SendMessageArgs{
				Recipient: &model.RecipientSendMessageRequest{
					ID: customerConversation.ExternalUserID,
				},
				Message: &model.MessageSendMessageRequest{Text: message},
				Tag:     string(fbclient.POST_PURCHASE_UPDATE),
			},
			PageID: customerConversation.ExternalPageID,
		})
	}
	return mq.CodeOK, nil
}

func (h *Handler) generateMessage(ffm *shipping.Fulfillment, shop *identity.Shop, historyState, historySubstate string) (message string, err error) {
	shippingNote := ffm.ShippingNote
	orderCode := ffm.ShippingCode
	codAmount := ffm.CODAmount

	if historySubstate != "" {
		substate, ok := substate.ParseSubstate(historySubstate)
		if !ok {
			return "", cm.Errorf(cm.FailedPrecondition, nil, "unsupported shipping substate %v", historyState)
		}
		message = templateForSubstate(substate, shippingNote, orderCode, codAmount)
	} else {
		state, ok := shipping_state.ParseState(historyState)
		if !ok {
			return "", cm.Errorf(cm.FailedPrecondition, nil, "unsupported shipping state %v", historyState)
		}
		message, err = templateForState(state, ffm, shop)
	}
	return
}

func templateForState(state shipping_state.State, ffm *shipping.Fulfillment, shop *identity.Shop) (string, error) {
	shopName := shop.Name
	customerName := ffm.AddressTo.FullName
	customerPhone := ffm.AddressTo.Phone
	customerAddress := ffm.AddressTo.Address1
	if ffm.AddressTo.Ward != "" {
		customerAddress += ", " + ffm.AddressTo.Ward
	}
	if ffm.AddressTo.District != "" {
		customerAddress += ", " + ffm.AddressTo.District
	}
	if ffm.AddressTo.Province != "" {
		customerAddress += ", " + ffm.AddressTo.Province
	}
	codAmount := ffm.CODAmount
	orderCode := ffm.ShippingCode
	switch state {
	case shipping_state.Created:
		tmpl := `Cảm ơn bạn đã mua hàng từ %s
Chi tiết đơn:
• Mã đơn: %s
• Tên người nhận: %s
• SDT: %s
• Địa chỉ: %s
• Thu hộ: %sđ
• Theo dõi: https://donhang.ghn.vn/?order_code=%s
Đơn hàng được chốt trên ứng dụng Faboshop!`
		return fmt.Sprintf(tmpl, shop.Name, orderCode, customerName, customerPhone, customerAddress, formatPrice(codAmount), orderCode), nil
	case shipping_state.Picking:
		tmpl := `Cập nhật trạng thái:
Nhà vận chuyển đang đến lấy hàng gửi cho %s
• Đơn hàng: %s - Thu hộ: %sđ
• Theo dõi: https://donhang.ghn.vn/?order_code=%s`
		return fmt.Sprintf(tmpl, customerName, orderCode, formatPrice(codAmount), orderCode), nil
	case shipping_state.Holding:
		tmpl := `Cập nhật trạng thái:
Đã bàn giao hàng cho nhà vận chuyển và vận chuyển đến %s
• Đơn hàng: %s - Thu hộ: %sđ
• Theo dõi: https://donhang.ghn.vn/?order_code=%s`
		return fmt.Sprintf(tmpl, customerAddress, orderCode, formatPrice(codAmount), orderCode), nil
	case shipping_state.Delivering:
		tmpl := `Cập nhật trạng thái:
Đơn hàng đang được giao đến %s. Shipper sẽ nhanh chóng giao hàng, %s vui lòng chờ điện thoại.
• Đơn hàng: %s - Thu hộ: %sđ
• Theo dõi: https://donhang.ghn.vn/?order_code=%s`
		return fmt.Sprintf(tmpl, customerName, customerName, orderCode, formatPrice(codAmount), orderCode), nil
	case shipping_state.Delivered:
		tmpl := `Cập nhật trạng thái:
Đơn hàng đã hoàn thành, xin cảm ơn bạn đã tin tưởng sử dụng dịch vụ của Shop %s.
• Đơn hàng: %s - Thu hộ: %sđ
• Theo dõi: https://donhang.ghn.vn/?order_code=%s`
		return fmt.Sprintf(tmpl, shopName, orderCode, formatPrice(codAmount), orderCode), nil
	case shipping_state.Returning:
		tmpl := `Cập nhật trạng thái:
Đơn hàng giao không thành công, nhà vận chuyển đã tiến hành hoàn trả hàng về Shop.
• Đơn hàng: %s - Thu hộ: %sđ
• Theo dõi: https://donhang.ghn.vn/?order_code=%s`
		return fmt.Sprintf(tmpl, orderCode, formatPrice(codAmount), orderCode), nil
	default:
		return "", cm.Errorf(cm.FailedPrecondition, nil, "unsupported shipping state %v", state)
	}
}

func templateForSubstate(substate substate.Substate, shippingNote, orderCode string, codAmount int) string {
	title := substate.GetLabelRefName()
	if shippingNote != "" {
		title += " - " + shippingNote
	}
	tmpl := `Cập nhật trạng thái:
%s
• Đơn hàng: %s
• Thu hộ: %sđ
Theo dõi: https://donhang.ghn.vn/?order_code=%s`
	return fmt.Sprintf(tmpl, strings.ToUpper(title), orderCode, formatPrice(codAmount), orderCode)
}

func formatPrice(n int) string {
	p := message.NewPrinter(language.Vietnamese)
	return p.Sprint(n)
}
