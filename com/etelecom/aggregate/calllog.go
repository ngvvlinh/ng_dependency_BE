package aggregate

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/etelecom/call_log_direction"
	contacting "o.o/api/main/contact"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
	"o.o/common/l"
)

func (a *EtelecomAggregate) CreateCallLogFromCDR(
	ctx context.Context, args *etelecom.CreateCallLogFromCDRArgs,
) (*etelecom.CallLog, error) {
	if err := args.Validate(); err != nil {
		return nil, err
	}

	var extensionNumber, phoneNumber string
	if args.Direction == call_log_direction.In {
		extensionNumber = args.Callee
		phoneNumber = args.Caller
	} else {
		extensionNumber = args.Caller
		phoneNumber = args.Callee
	}

	var callLog etelecom.CallLog
	if err := scheme.Convert(args, &callLog); err != nil {
		return nil, err
	}

	// find hotlines
	hotlineIDs, err := a.getHotlineIDs(ctx, args.ConnectionID, args.OwnerID)
	if err != nil {
		return nil, err
	}
	if len(hotlineIDs) == 0 {
		return nil, cm.Errorf(cm.NotFound, nil, "Not found hotline")
	}

	extension, err := a.extensionStore(ctx).HotlineIDs(hotlineIDs...).ExtensionNumber(extensionNumber).GetExtension()
	if err != nil {
		return nil, err
	}

	getContactsQuery := &contacting.GetContactsByPhoneQuery{
		ShopID: extension.AccountID,
		Phone:  phoneNumber,
	}
	if err := a.contactQuery.Dispatch(ctx, getContactsQuery); err != nil {
		return nil, err
	}
	contacts := getContactsQuery.Result
	if len(contacts) > 0 {
		callLog.ContactID = contacts[0].ID
	}

	callLog.ID = cm.NewID()
	callLog.ExtensionID = extension.ID
	callLog.AccountID = extension.AccountID
	callLog.HotlineID = extension.HotlineID

	callLogResult, err := a.callLogStore(ctx).CreateCallLog(&callLog)
	if err != nil {
		return nil, err
	}

	event := &etelecom.CallLogCreatedEvent{
		ID:         callLogResult.ID,
		Direction:  callLogResult.Direction,
		Callee:     callLogResult.Callee,
		Duration:   callLogResult.Duration,
		CallStatus: callLogResult.CallStatus,
		HotlineID:  callLogResult.HotlineID,
	}
	if _err := a.eventBus.Publish(ctx, event); _err != nil {
		ll.Error("Error event call log created :: ", l.Error(_err))
	}
	return callLogResult, nil
}

func (a *EtelecomAggregate) getHotlineIDs(ctx context.Context, connectionID, ownerID dot.ID) (hotlineIDs []dot.ID, _ error) {
	hotlines, err := a.hotlineStore(ctx).OptionalOwnerID(ownerID).ConnectionID(connectionID).ListHotlines()
	if err != nil {
		return nil, err
	}
	for _, hl := range hotlines {
		hotlineIDs = append(hotlineIDs, hl.ID)
	}
	return hotlineIDs, nil
}

func (a *EtelecomAggregate) UpdateCallLogPostage(ctx context.Context, args *etelecom.UpdateCallLogPostageArgs) error {
	if args.ID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing call log ID")
	}
	update := &etelecom.CallLog{
		DurationPostage: args.DurationForPostage,
		Postage:         args.Postage,
	}
	return a.callLogStore(ctx).ID(args.ID).UpdateCallLog(update)
}
