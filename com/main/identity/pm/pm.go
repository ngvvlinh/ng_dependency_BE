package pm

import (
	"context"
	"time"

	"o.o/api/etelecom"
	"o.o/api/main/address"
	"o.o/api/main/identity"
	"o.o/api/main/invitation"
	"o.o/api/top/types/etc/account_type"
	"o.o/api/top/types/etc/address_type"
	"o.o/backend/com/main/authorization/convert"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/sqlstore"
)

type ProcessManager struct {
	identityQuery    identity.QueryBus
	invitationQuery  invitation.QueryBus
	addressQuery     address.QueryBus
	addressAggr      address.CommandBus
	identityAggr     identity.CommandBus
	AccountUserStore sqlstore.AccountUserStoreInterface
}

func New(
	eventBus bus.EventRegistry,
	identityQ identity.QueryBus,
	identityAggr identity.CommandBus,
	invitationQ invitation.QueryBus,
	addressQuery address.QueryBus,
	addressAggr address.CommandBus,
	AccountUserStore sqlstore.AccountUserStoreInterface,
) *ProcessManager {
	p := &ProcessManager{
		identityQuery:    identityQ,
		invitationQuery:  invitationQ,
		addressQuery:     addressQuery,
		addressAggr:      addressAggr,
		identityAggr:     identityAggr,
		AccountUserStore: AccountUserStore,
	}
	p.registerEventHandlers(eventBus)
	return p
}

func (m *ProcessManager) registerEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.InvitationAccepted)
	eventBus.AddEventListener(m.AddressCreated)
	eventBus.AddEventListener(m.DefaultAddressUpdated)
	eventBus.AddEventListener(m.AccountDeleting)
	eventBus.AddEventListener(m.ExtensionAssigned)
	eventBus.AddEventListener(m.UserOfExtensionRemoved)
	eventBus.AddEventListener(m.UserUpdated)
}

func (m *ProcessManager) UserOfExtensionRemoved(ctx context.Context, event *etelecom.RemovedUserOfExtensionEvent) error {
	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		UserID:    event.UserID,
		AccountID: event.AccountID,
	}
	if err := m.AccountUserStore.GetAccountUser(ctx, getAccountUserQuery); err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Kh??ng t??m th???y nh??n vi??n")
	}

	cmd := &identity.UpdateExtensionNumberNormCommand{
		AccountID:       event.AccountID,
		UserID:          event.UserID,
		ExtensionNumber: "",
	}
	return m.identityAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) ExtensionAssigned(ctx context.Context, event *etelecom.AssignedExtensionEvent) error {
	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		UserID:    event.UserID,
		AccountID: event.AccountID,
	}
	if err := m.AccountUserStore.GetAccountUser(ctx, getAccountUserQuery); err != nil {
		return cm.Errorf(cm.ErrorCode(err), err, "Kh??ng t??m th???y nh??n vi??n")
	}
	cmd := &identity.UpdateExtensionNumberNormCommand{
		AccountID:       event.AccountID,
		UserID:          event.UserID,
		ExtensionNumber: event.ExtensionNumber,
	}
	return m.identityAggr.Dispatch(ctx, cmd)
}

func (m *ProcessManager) InvitationAccepted(ctx context.Context, event *invitation.InvitationAcceptedEvent) error {
	getInvitationQuery := &invitation.GetInvitationQuery{
		ID: event.ID,
	}
	if err := m.invitationQuery.Dispatch(ctx, getInvitationQuery); err != nil {
		return err
	}
	currInvitation := getInvitationQuery.Result

	getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
		Email: currInvitation.Email,
		Phone: currInvitation.Phone,
	}
	if err := m.identityQuery.Dispatch(ctx, getUserQuery); err != nil {
		return err
	}
	currUser := getUserQuery.Result
	userID := currUser.ID

	if currInvitation.Email != "" {
		if currUser.EmailVerifiedAt.IsZero() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Thao t??c th???t b???i, email ch??a ???????c x??c nh???n")
		}
	} else {
		if currUser.PhoneVerifiedAt.IsZero() {
			return cm.Errorf(cm.FailedPrecondition, nil, "Thao t??c th???t b???i, phone ch??a ???????c x??c nh???n")
		}
	}

	getAccountUserQuery := &identitymodelx.GetAccountUserQuery{
		UserID:    userID,
		AccountID: currInvitation.AccountID,
	}
	err := m.AccountUserStore.GetAccountUser(ctx, getAccountUserQuery)
	switch cm.ErrorCode(err) {
	case cm.NotFound:
		err := m.createAccountUserWithRoles(ctx, currInvitation, currUser)
		if err != nil {
			return err
		}
	case cm.NoError:
		return cm.Errorf(cm.Internal, nil, "unexpected (account_user exists)")
	default:
		return err
	}
	return nil
}

func (m *ProcessManager) createAccountUserWithRoles(
	ctx context.Context, currInvitation *invitation.Invitation, currUser *identity.User,
) error {
	createAccountUserCmd := &identitymodelx.CreateAccountUserCommand{
		AccountUser: &identitymodel.AccountUser{
			AccountID: currInvitation.AccountID,
			UserID:    currUser.ID,
			Status:    0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Permission: identitymodel.Permission{
				Roles: convert.ConvertRolesToStrings(currInvitation.Roles),
			},
			FullName:  cm.Coalesce(currInvitation.FullName, currUser.FullName),
			ShortName: cm.Coalesce(currInvitation.ShortName, currUser.ShortName),
			Position:  currInvitation.Position,
		},
	}
	if err := m.AccountUserStore.CreateAccountUser(ctx, createAccountUserCmd); err != nil {
		return err
	}
	return nil
}

/**
 *	This handle event address was created
 *	Update ship_from_address_id if this field of shop not exist
 */
func (m *ProcessManager) AddressCreated(ctx context.Context, event *address.AddressCreatedEvent) error {
	if event == nil {
		return nil
	}
	// accept only type ship from
	if event.Type != address_type.Shipfrom {
		return nil
	}

	cmd := &identity.GetShopByIDQuery{
		ID: event.AccountID,
	}
	if err := m.identityQuery.Dispatch(ctx, cmd); err != nil {
		return err
	}
	if cmd.Result == nil || cmd.Result.ShipFromAddressID != 0 {
		return nil
	}

	// get address info
	cmdAddress := &address.GetAddressByIDQuery{
		ID: event.ID,
	}
	if err := m.addressQuery.Dispatch(ctx, cmdAddress); err != nil {
		return err
	}
	addressInfo := cmdAddress.Result

	if addressInfo.IsDefault == true { // update ShopFromAddressID of shop
		if err := m.identityAggr.Dispatch(ctx, &identity.UpdateShipFromAddressIDCommand{
			ID:                addressInfo.AccountID,
			ShipFromAddressID: addressInfo.ID,
		}); err != nil {
			return err
		}
		return nil
	}

	if err := m.addressAggr.Dispatch(ctx, &address.UpdateDefaultAddressCommand{
		ShopID:    event.AccountID,
		AddressID: addressInfo.ID,
		Type:      address_type.Shipfrom.String(),
	}); err != nil {
		return err
	}

	return nil
}

func (m *ProcessManager) DefaultAddressUpdated(ctx context.Context, event *address.AddressDefaultUpdatedEvent) error {
	if event.ID == 0 || event.ShipFromAddressID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Th??ng tin update ship_from_address_id kh??ng h???p l???")
	}

	cmd := &identity.GetShopByIDQuery{
		ID: event.ID,
	}

	if err := m.identityQuery.Dispatch(ctx, cmd); err != nil {
		return err
	}
	cmdUpdateShipFromAddressID := &identity.UpdateShipFromAddressIDCommand{
		ID:                event.ID,
		ShipFromAddressID: event.ShipFromAddressID,
	}

	if err := m.identityAggr.Dispatch(ctx, cmdUpdateShipFromAddressID); err != nil {
		return err
	}

	return nil
}

func (m *ProcessManager) AccountDeleting(ctx context.Context, event *identity.AccountDeletingEvent) error {
	// shop ch??? ???????c delete khi:
	// 1. Xo?? h???t nh??n vi??n
	// 2. Xo?? h???t invite c??n ch??a ??c ch???p nh???n
	if event.AccountType != account_type.Shop {
		return cm.Errorf(cm.InvalidArgument, nil, "Does not support delete this account (type: %v, id: %v)", event.AccountType.Name(), event.AccountID)
	}

	queryAccountUsers := &identity.ListAccountUsersQuery{
		AccountID: event.AccountID,
	}
	if err := m.identityQuery.Dispatch(ctx, queryAccountUsers); err != nil {
		return err
	}
	if len(queryAccountUsers.Result.AccountUsers) >= 2 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Vui l??ng g??? c??c nh??n vi??n ??ang c?? quy???n qu???n tr???")
	}

	queryInvitations := &invitation.ListInvitationsNotAcceptedYetByAccountIDQuery{
		AccountID: event.AccountID,
	}
	if err := m.invitationQuery.Dispatch(ctx, queryInvitations); err != nil {
		return err
	}
	if len(queryInvitations.Result) > 0 {
		return cm.Errorf(cm.FailedPrecondition, nil, "Vui l??ng x??a c??c l???i m???i qu???n tr??? ch??a ???????c ch???p nh???n")
	}

	return nil
}

func (m *ProcessManager) UserUpdated(ctx context.Context, event *identity.UserUpdatedEvent) error {
	update := &identity.UpdateAccountUserInfoCommand{
		AccountID: event.AccountID,
		UserID:    event.UserID,
		FullName:  event.FullName,
		Phone:     event.Phone,
	}
	if err := m.identityAggr.Dispatch(ctx, update); err != nil {
		return err
	}
	return nil
}
