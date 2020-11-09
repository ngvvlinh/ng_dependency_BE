package aggregate

import (
	"context"
	"sort"

	"o.o/api/supporting/ticket"
	"o.o/backend/com/supporting/ticket/convert"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

var scheme = conversion.Build(convert.RegisterConversions)

func (a TicketAggregate) CreateTicketLabel(ctx context.Context, args *ticket.CreateTicketLabelArgs) (*ticket.TicketLabel, error) {
	var ticketLabelCore = &ticket.TicketLabel{}
	err := scheme.Convert(args, ticketLabelCore)
	if err != nil {
		return nil, err
	}
	if args.ParentID != 0 {
		_, err = a.TicketLabelStore(ctx).ID(args.ParentID).GetTicketLabel()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Không tìm thấy parent.")
		case cm.NoError:
			// no-op
		default:
			return nil, err
		}
	}

	if err := a.validateTicketLabelBeforeCreateOrUpdate(ctx, ticketLabelCore); err != nil {
		return nil, err
	}
	if err := a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err = a.TicketLabelStore(ctx).Create(ticketLabelCore); err != nil {
			return err
		}

		result, err := a.TicketLabelStore(ctx).ListTicketLabels()
		if err != nil {
			return err
		}

		if err := a.SetTicketLabels(&result); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return a.TicketLabelStore(ctx).ID(ticketLabelCore.ID).GetTicketLabel()
}

func (a TicketAggregate) UpdateTicketLabel(ctx context.Context, args *ticket.UpdateTicketLabelArgs) (*ticket.TicketLabel, error) {
	ticketLabelCore, err := a.TicketLabelStore(ctx).ID(args.ID).GetTicketLabel()
	if err != nil {
		return nil, err
	}

	//TODO check parent is not child
	if args.ParentID.Valid && args.ParentID.ID != 0 && args.ParentID.ID != ticketLabelCore.ParentID {
		_, err = a.TicketLabelStore(ctx).ID(args.ParentID.ID).GetTicketLabel()
		switch cm.ErrorCode(err) {
		case cm.NotFound:
			return nil, cm.Errorf(cm.InvalidArgument, nil, "Không tìm thấy label cha.")
		case cm.NoError:
			// no-op
		default:
			return nil, err
		}
	}
	if err = scheme.Convert(args, ticketLabelCore); err != nil {
		return nil, err
	}
	if err = a.validateTicketLabelBeforeCreateOrUpdate(ctx, ticketLabelCore); err != nil {
		return nil, err
	}

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if err = a.TicketLabelStore(ctx).ID(args.ID).UpdateTicketLabel(ticketLabelCore); err != nil {
			return err
		}

		result, err := a.TicketLabelStore(ctx).ListTicketLabels()
		if err != nil {
			return err
		}

		return a.SetTicketLabels(&result)
	})
	if err != nil {
		return nil, err
	}

	result, err := a.TicketLabelStore(ctx).ListTicketLabels()
	if err != nil {
		return nil, err
	}

	if err = a.SetTicketLabels(&result); err != nil {
		return nil, err
	}

	return ticketLabelCore, nil

}

func (a TicketAggregate) DeleteTicketLabel(ctx context.Context, args *ticket.DeleteTicketLabelArgs) (int, error) {
	if _, err := a.TicketLabelStore(ctx).ID(args.ID).GetTicketLabel(); err != nil {
		return 0, err
	}

	labels, err := a.listTicketLabels(ctx)
	if err != nil {
		return 0, err
	}

	label := getLabel(args.ID, labels)
	if !args.DeleteChild && label != nil && len(label.Children) > 0 {
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Label có chứa label con")
	}

	ids := getListLabelChildIDs(label)
	ids = append(ids, args.ID)

	err = a.db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		if _, err = a.TicketLabelStore(ctx).IDs(ids...).SoftDelete(); err != nil {
			return err
		}

		// delete all ticket_label_ticket_label_external by ticket_label_id
		if _, err := a.TicketLabelTicketLabelExternalStore(ctx).TicketLabelIDs(ids).SoftDelete(); err != nil {
			return err
		}

		result, err := a.TicketLabelStore(ctx).ListTicketLabels()
		if err != nil {
			return err
		}

		return a.SetTicketLabels(&result)
	})
	if err != nil {
		return 0, err
	}

	return len(ids), nil
}

func (a TicketAggregate) validateTicketLabelBeforeCreateOrUpdate(ctx context.Context, label *ticket.TicketLabel) error {
	if label.Name == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Tên label không hợp lệ")
	}
	if label.Code == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Mã label không hợp lệ")
	}
	ticketCore, err := a.TicketLabelStore(ctx).Code(label.Code).GetTicketLabel()
	switch cm.ErrorCode(err) {
	case cm.NoError:
		if ticketCore.ID != label.ID {
			return cm.Errorf(cm.InvalidArgument, nil, "Mã label đã tồn tại")
		}
	case cm.NotFound:
		// no-op
	default:
		return err
	}

	return nil
}

func getListLabelChildIDs(l *ticket.TicketLabel) []dot.ID {
	if l == nil {
		return []dot.ID{}
	}
	var result []dot.ID
	result = append(result, l.ID)
	for _, v := range l.Children {
		result = append(result, getListLabelChildIDs(v)...)
	}
	return result
}

func getLabel(id dot.ID, a []*ticket.TicketLabel) *ticket.TicketLabel {
	if len(a) == 0 {
		return nil
	}
	for _, v := range a {
		if v.ID == id {
			return v
		}
	}
	for _, v := range a {
		result := getLabel(id, v.Children)
		if result != nil {
			return result
		}
	}
	return nil
}

func MakeTreeLabel(ticketLabels []*ticket.TicketLabel) []*ticket.TicketLabel {
	sort.Slice(ticketLabels, func(i, j int) bool {
		return ticketLabels[i].ParentID < ticketLabels[j].ParentID
	})
	var result []*ticket.TicketLabel
	var ok bool
	for _, ticketLabel := range ticketLabels {
		if ticketLabel.ParentID == 0 {
			result = append(result, ticketLabel)
			continue
		}
		result, ok = addToTreeLabel(ticketLabel, result)
		if ok {
			// error or something else
		}
	}
	return result
}

func addToTreeLabel(child *ticket.TicketLabel, father []*ticket.TicketLabel) ([]*ticket.TicketLabel, bool) {
	if len(father) == 0 {
		return father, false
	}
	for k, v := range father {
		if child.ParentID == v.ID {
			father[k].Children = append(father[k].Children, child)
			return father, true
		}
	}
	for k, _ := range father {
		var found bool
		father[k].Children, found = addToTreeLabel(child, father[k].Children)
		if found {
			return father, true
		}
	}
	return father, false
}

// get all label_ids from this id to father have id = 0
func getListLabelFatherID(id dot.ID, labels []*ticket.TicketLabel) []dot.ID {
	if len(labels) == 0 {
		return nil
	}

	for _, label := range labels {
		if label.ID == id {
			return []dot.ID{label.ID}
		}
		ids := getListLabelFatherID(id, label.Children)
		if len(ids) != 0 {
			return append(ids, label.ID)
		}
	}
	return nil
}
