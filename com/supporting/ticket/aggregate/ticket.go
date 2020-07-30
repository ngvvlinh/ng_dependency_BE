package aggregate

import (
	"context"
	"time"

	"o.o/api/main/moneytx"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	"o.o/api/supporting/ticket"
	"o.o/api/top/types/etc/shipping_provider"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_state"
	com "o.o/backend/com/main"
	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/com/supporting/ticket/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

var _ ticket.Aggregate = &TicketAggregate{}

type TicketAggregate struct {
	TicketStore        sqlstore.TicketStoreFactory
	TicketLabelStore   sqlstore.TicketLabelStoreFactory
	EventBus           capi.EventBus
	db                 *cmsql.Database
	MoneyTxQuery       moneytx.QueryBus
	TicketCommentStore sqlstore.TicketCommentStoreFactory
	ShippingQuery      shipping.QueryBus
	OrderQuery         ordering.QueryBus
	RedisStore         redis.Store
}

func NewTicketAggregate(
	eventBus capi.EventBus,
	db com.MainDB,
	moneyTxQ moneytx.QueryBus,
	shippingQ shipping.QueryBus,
	orderQ ordering.QueryBus,
	redisStore redis.Store,
	// carrierTicket ticket.
) *TicketAggregate {
	return &TicketAggregate{
		TicketStore:        sqlstore.NewTicketStore(db),
		TicketCommentStore: sqlstore.NewTicketCommentStore(db),
		TicketLabelStore:   sqlstore.NewTicketLabelStore(db),
		MoneyTxQuery:       moneyTxQ,
		EventBus:           eventBus,
		ShippingQuery:      shippingQ,
		db:                 db,
		OrderQuery:         orderQ,
		RedisStore:         redisStore,
	}
}

func TicketAggregateMessageBus(q *TicketAggregate) ticket.CommandBus {
	b := bus.New()
	return ticket.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a TicketAggregate) CreateTicket(ctx context.Context, args *ticket.CreateTicketArgs) (*ticket.Ticket, error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	var ticketCore = &ticket.Ticket{}
	err := scheme.Convert(args, ticketCore)
	if err != nil {
		return nil, err
	}
	// check reference
	var refCode = ""
	if args.RefID != 0 {
		switch args.RefType {
		case ticket_ref_type.FFM:
			query := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
				ID: args.RefID,
			}
			if err := a.ShippingQuery.Dispatch(ctx, query); err != nil {
				if cm.ErrorCode(err) == cm.NotFound {
					return nil, cm.Errorf(cm.NotFound, err, "Không tìm thấy đơn giao hàng")
				}
				return nil, err
			}
			refCode = query.Result.ShippingCode
			if query.Result.ShippingProvider == shipping_provider.GHN {
				// Send ticket to ghn
				err := a.createTicketThirdParty(ctx, ticketCore)
				if err != nil {
					return nil, err
				}
			}
		case ticket_ref_type.MoneyTransaction:
			query := &moneytx.GetMoneyTxShippingByIDQuery{
				MoneyTxShippingID: args.RefID,
				ShopID:            args.AccountID,
			}
			if err := a.MoneyTxQuery.Dispatch(ctx, query); err != nil {
				if cm.ErrorCode(err) == cm.NotFound {
					return nil, cm.Errorf(cm.NotFound, err, "Không tìm thấy phiên chuyển tiền")
				}
				return nil, err
			}
			refCode = query.Result.Code
		case ticket_ref_type.OrderTrading:
			queryOrder := &ordering.GetOrderByIDQuery{
				ID: args.RefID,
			}
			err := a.OrderQuery.Dispatch(ctx, queryOrder)
			if err != nil {
				if cm.ErrorCode(err) == cm.NotFound {
					return nil, cm.Errorf(cm.NotFound, err, "Không tìm thấy đơn hàng")
				}
				return nil, err
			}
			refCode = queryOrder.Result.Code
		default:
			//no-op(other)
		}
		if args.RefCode != "" && args.RefCode != refCode {
			return nil, cm.Errorf(cm.NotFound, nil, "ref_code không đúng")
		}
		ticketCore.RefCode = refCode
	}
	labels, err := a.listTicketLabels(ctx)
	if err != nil {
		return nil, err
	}
	// get all father label_ids of all labels
	for _, v := range args.LabelIDs {

		listLabelIDs, ok := getListLabelFatherID(v, labels)
		if !ok {
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Label không tồn tại")
		}
		for _, labelID := range listLabelIDs {
			if !cm.IDsContain(ticketCore.LabelIDs, labelID) {
				ticketCore.LabelIDs = append(ticketCore.LabelIDs, labelID)
			}
		}

	}

	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		err = a.TicketStore(ctx).Create(ticketCore)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(ticketCore.ID).GetTicket()
}

func (a TicketAggregate) createTicketThirdParty(ctx context.Context, args *ticket.Ticket) error {
	//TODO(Nam)
	//CreateTicket
	// sửa dụng tài khoản của connection
	// bỏ shipping_Provider, ffm ->coonection ->ghn -> driver

	// shop_connection -> aggre

	// khởi driver -> ghn

	// tao ticket
	panic("implement me")
}

func (a TicketAggregate) UpdateTicketInfo(ctx context.Context, args *ticket.UpdateTicketInfoArgs) (*ticket.Ticket, error) {
	//TODO maybe not use
	panic("implement me")
}

func (a TicketAggregate) ConfirmTicket(ctx context.Context, args *ticket.ConfirmTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if args.ConfirmBy == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConfirmID")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket đã đóng")
	}
	if !args.IsLeader {
		isPermission := false
		for _, v := range ticketCore.AssignedUserIDs {
			if v == args.ConfirmBy {
				isPermission = true
				break
			}
		}
		if !isPermission {
			return nil, cm.Errorf(cm.PermissionDenied, nil, "Ticket không thuộc sự quản lí của bạn")
		}
	}
	var ticketModel = &model.Ticket{
		ConfirmedAt: time.Now(),
		ConfirmedBy: args.ConfirmBy,
		Note:        args.Note,
		UpdatedBy:   args.ConfirmBy,
		Status:      status5.S,
		UpdatedAt:   time.Now(),
		State:       ticket_state.Processing,
	}
	err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel)
	if err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a TicketAggregate) CloseTicket(ctx context.Context, args *ticket.CloseTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if args.ClosedBy == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ConfirmID")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	if !args.IsLeader && args.ClosedBy != ticketCore.ConfirmedBy {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "Ticket không thuộc sự quản lí của bạn")
	}
	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket đã đóng")
	}
	var ticketModel = &model.Ticket{
		ClosedAt:  time.Now(),
		ClosedBy:  args.ClosedBy,
		Note:      args.Note,
		UpdatedBy: args.ClosedBy,
		UpdatedAt: time.Now(),
		State:     args.State,
	}
	switch args.State {
	case ticket_state.Success:
		ticketModel.Status = status5.P
	case ticket_state.Fail:
		ticketModel.Status = status5.NS
	case ticket_state.Ignore, ticket_state.Cancel:
		ticketModel.Status = status5.N
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "state đóng ticket không hợp lệ")
	}
	err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel)
	if err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a TicketAggregate) ReopenTicket(ctx context.Context, args *ticket.ReopenTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	var state = ticket_state.New
	var status = status5.Z
	if len(ticketCore.AssignedUserIDs) > 0 {
		state = ticket_state.Received
		status = status5.S
	}
	switch ticketCore.Status {
	case status5.N, status5.NS, status5.P:
		// no-op
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket chưa được close không thể mở lại.")
	}
	var ticketModel = &model.Ticket{
		Note:   args.Note,
		State:  state,
		Status: status,
	}
	err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel)
	if err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a TicketAggregate) AssignTicket(ctx context.Context, args *ticket.AssignedTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	assignedUserIDs := ticketCore.AssignedUserIDs
	if !args.IsLeader {
		// if not leader user will add themselves
		for _, v := range assignedUserIDs {
			if v == args.UpdatedBy {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Bạn đã được thêm vào ticket này rồi.")
			}
		}
		assignedUserIDs = append(assignedUserIDs, args.UpdatedBy)
	} else {
		assignedUserIDs = args.AssignedUserIDs
	}
	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket đã đóng")
	}
	var state = ticket_state.New
	var status = status5.Z
	if len(assignedUserIDs) > 0 {
		state = ticket_state.Received
		status = status5.S
	}
	ticketCore.State = state
	ticketCore.UpdatedBy = args.UpdatedBy
	ticketCore.UpdatedAt = time.Now()
	ticketCore.AssignedUserIDs = assignedUserIDs
	ticketCore.Status = status
	// follow requirement, we have case update status 2 -> 0. so have to use update all
	err = a.TicketStore(ctx).ID(args.ID).UpdateTicketALL(ticketCore)
	if err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a TicketAggregate) UnassignTicket(ctx context.Context, args *ticket.UnssignTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	if ticketCore.Status != status5.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Bạn không thể bỏ chỉ định trên ticket khác trạng thái mới.")
	}
	var assignedUserIDs []dot.ID
	isExisted := false
	for _, v := range ticketCore.AssignedUserIDs {
		if v == args.UpdatedBy {
			isExisted = true
			continue
		}
		assignedUserIDs = append(assignedUserIDs, v)
	}
	if !isExisted {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Bạn chưa được thêm vào ticket này.")
	}
	var state = ticket_state.New
	var status = status5.Z
	if len(assignedUserIDs) > 0 {
		state = ticket_state.Received
		status = status5.S
	}
	ticketCore.State = state
	ticketCore.UpdatedBy = args.UpdatedBy
	ticketCore.UpdatedAt = time.Now()
	ticketCore.AssignedUserIDs = assignedUserIDs
	ticketCore.Status = status
	// follow requirement, we have case update status 2 -> 0. so have to use update all
	err = a.TicketStore(ctx).ID(args.ID).UpdateTicketALL(ticketCore)
	if err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a TicketAggregate) listTicketLabels(ctx context.Context) ([]*ticket.TicketLabel, error) {
	var labels []*ticket.TicketLabel
	err := a.RedisStore.Get("ticket_labels", &labels)
	switch err {
	case redis.ErrNil:
		// no-op
	case nil:
		return labels, nil
	default:
		return nil, err
	}
	labels, err = a.TicketLabelStore(ctx).ListTicketLabels()
	if err != nil {
		return nil, err
	}
	err = a.SetTicketLabels(labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (a TicketAggregate) SetTicketLabels(labels []*ticket.TicketLabel) error {
	labels = MakeTreeLabel(labels)
	return a.RedisStore.SetWithTTL("ticket_labels", labels, 1*24*60*60)
}
