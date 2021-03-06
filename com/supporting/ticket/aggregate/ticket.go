package aggregate

import (
	"context"
	"fmt"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/main/contact"
	"o.o/api/main/identity"
	"o.o/api/main/moneytx"
	"o.o/api/main/ordering"
	"o.o/api/main/shipping"
	"o.o/api/supporting/ticket"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/status5"
	"o.o/api/top/types/etc/ticket/ticket_ref_type"
	"o.o/api/top/types/etc/ticket/ticket_state"
	"o.o/api/top/types/etc/ticket/ticket_type"
	com "o.o/backend/com/main"
	"o.o/backend/com/supporting/ticket/model"
	"o.o/backend/com/supporting/ticket/provider"
	"o.o/backend/com/supporting/ticket/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi"
	"o.o/capi/dot"
)

const (
	ticketLabelsVersion = "v1.6"
)

var _ ticket.Aggregate = &TicketAggregate{}

type TicketAggregate struct {
	TicketStore                         sqlstore.TicketStoreFactory
	TicketCommentStore                  sqlstore.TicketCommentStoreFactory
	TicketLabelStore                    sqlstore.TicketLabelStoreFactory
	TicketLabelExternalStore            sqlstore.TicketLabelExternalStoreFactory
	TicketLabelTicketLabelExternalStore sqlstore.TicketLabelTicketLabelExternalsStoreFactory
	EventBus                            capi.EventBus
	db                                  *cmsql.Database
	MoneyTxQuery                        moneytx.QueryBus
	ShippingQuery                       shipping.QueryBus
	OrderQuery                          ordering.QueryBus
	IdentityQuery                       identity.QueryBus
	TicketManager                       *provider.TicketManager
	ConnectionQuery                     connectioning.QueryBus
	ContactQuery                        contact.QueryBus
	RedisStore                          redis.Store
}

func NewTicketAggregate(
	eventBus capi.EventBus,
	db com.MainDB,
	moneyTxQ moneytx.QueryBus,
	shippingQ shipping.QueryBus,
	orderQ ordering.QueryBus,
	identityQ identity.QueryBus,
	ticketManager *provider.TicketManager,
	connectionQ connectioning.QueryBus,
	contactQ contact.QueryBus,
	redisStore redis.Store,
) *TicketAggregate {
	return &TicketAggregate{
		TicketStore:                         sqlstore.NewTicketStore(db),
		TicketCommentStore:                  sqlstore.NewTicketCommentStore(db),
		TicketLabelStore:                    sqlstore.NewTicketLabelStore(db),
		TicketLabelExternalStore:            sqlstore.NewTicketLabelExternalStore(db),
		TicketLabelTicketLabelExternalStore: sqlstore.NewTicketLabelTicketLabelExternalStore(db),
		MoneyTxQuery:                        moneyTxQ,
		EventBus:                            eventBus,
		ShippingQuery:                       shippingQ,
		db:                                  db,
		OrderQuery:                          orderQ,
		IdentityQuery:                       identityQ,
		ConnectionQuery:                     connectionQ,
		TicketManager:                       ticketManager,
		ContactQuery:                        contactQ,
		RedisStore:                          redisStore,
	}
}

func TicketAggregateMessageBus(q *TicketAggregate) ticket.CommandBus {
	b := bus.New()
	return ticket.NewAggregateHandler(q).RegisterHandlers(b)
}

func (a *TicketAggregate) CreateTicket(ctx context.Context, args *ticket.CreateTicketArgs) (_ *ticket.Ticket, err error) {
	if args.AccountID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing AccountID")
	}
	if args.Type != ticket_type.System && args.Type != ticket_type.Internal {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Unsupported type %v", args.Type)
	}

	ticketCore := &ticket.Ticket{}
	if err := scheme.Convert(args, ticketCore); err != nil {
		return nil, err
	}

	if len(args.LabelIDs) != 0 {
		var labels []*ticket.TicketLabel
		labels, err = a.listTicketLabels(ctx, listTicketLabelsArgs{
			Type:   args.Type,
			ShopID: args.AccountID,
		})
		if err != nil {
			return nil, err
		}

		// get all father label_ids of all labels
		for _, labelID := range args.LabelIDs {
			listLabelIDs := getListLabelFatherID(labelID, labels)
			if len(listLabelIDs) == 0 {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "Label kh??ng t???n t???i (id = %d)", labelID)
			}
			for _, _labelID := range listLabelIDs {
				if !cm.IDsContain(ticketCore.LabelIDs, _labelID) {
					ticketCore.LabelIDs = append(ticketCore.LabelIDs, _labelID)
				}
			}
		}
	}

	// get reference ticket
	if ticketCore.RefTicketID.Valid {
		if _, err := a.TicketStore(ctx).ID(ticketCore.RefTicketID.ID).GetTicket(); err != nil {
			return nil, err
		}
	}

	// check reference code
	if args.RefID != 0 {
		getReferenceItemArgs := &ticket.GetReferenceItemArgs{
			RefID:     args.RefID,
			RefType:   args.RefType,
			AccountID: args.AccountID,
		}
		res, err := a.checkRefItem(ctx, getReferenceItemArgs)
		if err != nil {
			return nil, err
		}
		refCode := res.RefCode
		connectionID := res.ConnectionID

		// check ref_code
		if args.RefCode != "" && args.RefCode != refCode {
			return nil, cm.Errorf(cm.NotFound, nil, "ref_code kh??ng ????ng")
		}
		ticketCore.RefCode = refCode
		ticketCore.ConnectionID = connectionID
	}

	if err = a.TicketStore(ctx).Create(ticketCore); err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(ticketCore.ID).GetTicket()
}

func (a *TicketAggregate) UpdateTicketInfo(ctx context.Context, args *ticket.UpdateTicketInfoArgs) (*pbcm.UpdatedResponse, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thi???u th??ng tin ticket ID")
	}
	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket ???? ????ng")
	}
	if args.RefID != 0 {
		getReferenceItemArgs := &ticket.GetReferenceItemArgs{
			RefID:     args.RefID,
			RefType:   args.RefType,
			AccountID: args.AccountID,
		}
		_, err = a.checkRefItem(ctx, getReferenceItemArgs)
		if err != nil {
			return nil, err
		}
	}
	ticket := &ticket.Ticket{
		Title:       args.Title,
		Description: args.Description,
		RefType:     args.RefType,
		RefID:       args.RefID,
		LabelIDs:    args.Labels,
	}
	if err := a.TicketStore(ctx).ID(args.ID).OptionalAccountID(args.AccountID).UpdateTicket(ticket); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (a *TicketAggregate) ConfirmTicket(ctx context.Context, args *ticket.ConfirmTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if args.ConfirmBy == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing confirm_by")
	}

	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}
	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket ???? ????ng")
	}

	// system: leader c?? th??? confirm m???i ticket
	// nh???ng ng?????i ???????c assign v??o ticker c?? th??? confirm ticket
	if ticketCore.Type == ticket_type.System && !args.IsLeader {
		isPermission := false
		for _, v := range ticketCore.AssignedUserIDs {
			if v == args.ConfirmBy {
				isPermission = true
				break
			}
		}
		if !isPermission {
			return nil, cm.Errorf(cm.PermissionDenied, nil, "Ticket kh??ng thu???c s??? qu???n l?? c???a b???n")
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
	if err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel); err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a *TicketAggregate) CloseTicket(ctx context.Context, args *ticket.CloseTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}
	if args.ClosedBy == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing closed_by")
	}

	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}

	// system: ch??? leader ho???c ng?????i confirm m???i ???????c close ticket
	// internal: kh??ng ph??n quy???n ch??? n??y
	if ticketCore.Type == ticket_type.System && !args.IsLeader && args.ClosedBy != ticketCore.ConfirmedBy {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "Ticket kh??ng thu???c s??? qu???n l?? c???a b???n")
	}

	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket ???? ????ng")
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
	case ticket_state.Success, ticket_state.Fail,
		ticket_state.Ignore, ticket_state.Cancel:
		// no-op
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "state ????ng ticket kh??ng h???p l???")
	}
	ticketModel.Status = ticketModel.State.ToStatus5()
	if err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel); err != nil {
		return nil, err
	}
	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a *TicketAggregate) ReopenTicket(ctx context.Context, args *ticket.ReopenTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}

	// khi reopen th?? tr???ng th??i ticket s??? l?? new
	// n???u c?? ng?????i ???????c assign v??o ticket th?? tr???ng th??i s??? l?? received
	var state = ticket_state.New
	if len(ticketCore.AssignedUserIDs) > 0 {
		state = ticket_state.Received
	}

	switch ticketCore.Status {
	case status5.N, status5.NS, status5.P:
		// no-op
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket ch??a ???????c close kh??ng th??? m??? l???i.")
	}

	var ticketModel = &model.Ticket{
		Note:  args.Note,
		State: state,
	}
	if err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel); err != nil {
		return nil, err
	}

	if err = a.TicketStore(ctx).ID(args.ID).UpdateTicketStatus(state.ToStatus5()); err != nil {
		return nil, err
	}

	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a *TicketAggregate) AssignTicket(ctx context.Context, args *ticket.AssignedTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}

	assignedUserIDs := ticketCore.AssignedUserIDs

	// system: if not leader user will add themselves
	if ticketCore.Type == ticket_type.System && !args.IsLeader {
		for _, v := range assignedUserIDs {
			if v == args.UpdatedBy {
				return nil, cm.Errorf(cm.InvalidArgument, nil, "B???n ???? ???????c th??m v??o ticket n??y r???i.")
			}
		}
		assignedUserIDs = append(assignedUserIDs, args.UpdatedBy)
	} else {
		assignedUserIDs = args.AssignedUserIDs
	}

	if ticketCore.Status != status5.Z && ticketCore.Status != status5.S {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Ticket ???? ????ng")
	}

	ticketModel := &model.Ticket{
		UpdatedBy:       args.UpdatedBy,
		UpdatedAt:       time.Now(),
		AssignedUserIDs: assignedUserIDs,
	}

	// Khi assign ticket m???i t???o cho 1 ng?????i: chuy???n tr???ng th??i t??? new -> received
	// C??n l???i th?? gi??? nguy??n tr???ng th??i c??.
	if ticketCore.State == ticket_state.New {
		ticketModel.State = ticket_state.Received
		ticketModel.Status = ticket_state.Received.ToStatus5()
	}

	if err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel); err != nil {
		return nil, err
	}

	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

func (a *TicketAggregate) UnassignTicket(ctx context.Context, args *ticket.UnassignTicketArgs) (*ticket.Ticket, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing ID")
	}

	ticketCore, err := a.TicketStore(ctx).ID(args.ID).GetTicket()
	if err != nil {
		return nil, err
	}

	if ticketCore.Status != status5.Z {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "B???n kh??ng th??? b??? ch??? ?????nh tr??n ticket kh??c tr???ng th??i m???i.")
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
		return nil, cm.Errorf(cm.InvalidArgument, nil, "B???n ch??a ???????c th??m v??o ticket n??y.")
	}

	var state = ticket_state.New
	if len(assignedUserIDs) > 0 {
		state = ticket_state.Received
	}

	ticketModel := &model.Ticket{
		State:           state,
		UpdatedBy:       args.UpdatedBy,
		UpdatedAt:       time.Now(),
		AssignedUserIDs: assignedUserIDs,
		Status:          state.ToStatus5(),
	}
	if err = a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(ticketModel); err != nil {
		return nil, err
	}

	return a.TicketStore(ctx).ID(args.ID).GetTicket()
}

type listTicketLabelsArgs struct {
	Type   ticket_type.TicketType
	ShopID dot.ID
}

func (a *TicketAggregate) listTicketLabels(ctx context.Context, args listTicketLabelsArgs) ([]*ticket.TicketLabel, error) {
	var labels []*ticket.TicketLabel
	redisKey := generateTicketLabelKey(ctx, args.ShopID)
	if args.Type == ticket_type.System {
		redisKey = generateTicketLabelKey(ctx, 0)
	}
	err := a.RedisStore.Get(redisKey, &labels)
	switch err {
	case redis.ErrNil:
		// no-op
	case nil:
		return labels, nil
	default:
		return nil, err
	}
	query := a.TicketLabelStore(ctx).Type(args.Type)

	if args.Type == ticket_type.Internal {
		query = query.ShopID(args.ShopID)
	}
	labels, err = query.ListTicketLabels()
	if err != nil {
		return nil, err
	}

	if err := a.SetTicketLabels(ctx, args.ShopID, args.Type, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

func (a *TicketAggregate) SetTicketLabels(
	ctx context.Context, shopID dot.ID,
	typ ticket_type.TicketType, labels *[]*ticket.TicketLabel) error {
	*labels = MakeTreeLabel(*labels)
	if typ == ticket_type.System {
		return a.RedisStore.SetWithTTL(generateTicketLabelKey(ctx, 0), labels, 1*24*60*60)
	}
	return a.RedisStore.SetWithTTL(generateTicketLabelKey(ctx, shopID), labels, 1*24*60*60)
}

func (a *TicketAggregate) UpdateTicketRefTicketID(ctx context.Context, args *ticket.UpdateTicketRefTicketIDArgs) (*pbcm.UpdatedResponse, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Thi???u th??ng tin ticket ID")
	}
	if !args.RefTicketID.Valid {
		return &pbcm.UpdatedResponse{Updated: 0}, nil
	}
	update := &model.Ticket{
		RefTicketID: args.RefTicketID,
	}
	if err := a.TicketStore(ctx).ID(args.ID).UpdateTicketDB(update); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func generateTicketLabelKey(ctx context.Context, shopID dot.ID) string {
	return fmt.Sprintf("ticket_labels:%s:wl%s:sh%s", ticketLabelsVersion, wl.GetWLPartnerID(ctx), shopID.String())
}

func (a *TicketAggregate) checkRefItem(ctx context.Context, args *ticket.GetReferenceItemArgs) (*ticket.GetReferenceItemResponse, error) {
	var refCode = ""
	var connectionID dot.ID
	switch args.RefType {
	case ticket_ref_type.FFM:
		getFfmQuery := &shipping.GetFulfillmentByIDOrShippingCodeQuery{
			ID: args.RefID,
		}
		if err := a.ShippingQuery.Dispatch(ctx, getFfmQuery); err != nil {
			if cm.ErrorCode(err) == cm.NotFound {
				return nil, cm.Errorf(cm.NotFound, err, "Kh??ng t??m th???y ????n giao h??ng")
			}
			return nil, err
		}
		refCode = getFfmQuery.Result.ShippingCode
		connectionID = getFfmQuery.Result.ConnectionID
	case ticket_ref_type.MoneyTransaction:
		// type system
		getMoneyTxQuery := &moneytx.GetMoneyTxShippingByIDQuery{
			MoneyTxShippingID: args.RefID,
			ShopID:            args.AccountID,
		}
		if err := a.MoneyTxQuery.Dispatch(ctx, getMoneyTxQuery); err != nil {
			if cm.ErrorCode(err) == cm.NotFound {
				return nil, cm.Errorf(cm.NotFound, err, "Kh??ng t??m th???y phi??n chuy???n ti???n")
			}
			return nil, err
		}
		refCode = getMoneyTxQuery.Result.Code
	case ticket_ref_type.OrderTrading:
		// type system
		getOrderQuery := &ordering.GetOrderByIDQuery{
			ID: args.RefID,
		}
		err := a.OrderQuery.Dispatch(ctx, getOrderQuery)
		if err != nil {
			if cm.ErrorCode(err) == cm.NotFound {
				return nil, cm.Errorf(cm.NotFound, err, "Kh??ng t??m th???y ????n h??ng")
			}
			return nil, err
		}
		refCode = getOrderQuery.Result.Code
	case ticket_ref_type.Contact:
		getContactQuery := &contact.GetContactByIDQuery{
			ID:     args.RefID,
			ShopID: args.AccountID,
		}
		if err := a.ContactQuery.Dispatch(ctx, getContactQuery); err != nil {
			if cm.ErrorCode(err) == cm.NotFound {
				return nil, cm.Errorf(cm.NotFound, err, "Kh??ng t??m th???y li??n l???c")
			}
			return nil, err
		}

	default:
		//no-op(other)
	}
	return &ticket.GetReferenceItemResponse{
		ConnectionID: connectionID,
		RefCode:      refCode,
	}, nil
}
